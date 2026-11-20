// A small typed client for the ticketing API (internal/api/openapi.yaml).
// Requests go to /api/..., which the Vite dev server proxies to the service.

export type Address = string;

export interface CreatedUser {
  id: string;
  subject: string;
  wallet_address: Address;
  token: string;
}

export interface CreatedOrganizer {
  id: string;
  name: string;
  box_office_address: Address;
  gate_address: Address;
  payout_address: Address;
  token: string;
}

export interface Event {
  id: string;
  organizer_id: string;
  name: string;
  symbol: string;
  starts_at: string;
  max_supply: number;
  face_value: string;
  max_resale_bps: number;
  status: "deploying" | "on_sale" | "failed";
  contract_address?: Address | null;
  tickets_issued: number;
  tickets_sold: number;
  deploy_operation_id: string;
}

export interface Ticket {
  id: string;
  event_id: string;
  token_id: string;
  status: "pending_mint" | "active" | "checked_in" | "failed";
  owner_address?: Address | null;
  listing_price?: string | null;
  mint_operation_id: string;
  checkin_operation_id?: string | null;
}

/** A listing, cancellation or purchase, saved and on its way. */
export interface ResaleAction {
  operation_id: string;
  ticket: Ticket;
}

export interface Operation {
  id: string;
  status: "signed" | "submitted" | "mined" | "reverted" | "failed";
  tx_hash: string;
  block_number?: number | null;
  revert_reason?: string | null;
  last_error?: string | null;
}

/**
 * A holder's signed check-in pass (EIP-712), valid five minutes. What the QR
 * code at the door holds; the gate sends it back as the check-in body.
 */
export interface CheckInPass {
  ticket_id: string;
  token_id: string;
  contract_address: Address;
  chain_id: number;
  deadline: number;
  expires_at: string;
  signature: string;
}

export const isFinal =(s: Operation["status"]) => s === "mined" || s === "reverted" || s === "failed";

export interface NewEvent {
  name: string;
  symbol: string;
  starts_at: string;
  max_supply: number;
  face_value: string;
  max_resale_bps: number;
}

/** An RFC 9457 problem from the API, e.g. 409 symbol_taken. */
export class ApiError extends Error {
  constructor(
    readonly status: number,
    readonly code: string,
    readonly title: string,
    readonly detail?: string,
  ) {
    super(detail ? `${title}: ${detail}` : title);
  }
}

async function call<T>(
  method: string,
  path: string,
  opts: { body?: unknown; token?: string; idempotent?: boolean } = {},
): Promise<T> {
  const headers: Record<string, string> = {};
  if (opts.body !== undefined) headers["Content-Type"] = "application/json";
  if (opts.token) headers["Authorization"] = `Bearer ${opts.token}`;
  // Writes that send a transaction need a key; a fresh one per form submit.
  if (opts.idempotent) headers["Idempotency-Key"] = crypto.randomUUID();

  let res: Response;
  try {
    res = await fetch(`/api${path}`, {
      method,
      headers,
      body: opts.body === undefined ? undefined : JSON.stringify(opts.body),
    });
  } catch {
    throw new ApiError(0, "unreachable", "Can't reach the service", "Is it running on port 8080?");
  }
  const text = await res.text();
  const data = text ? JSON.parse(text) : undefined;
  if (!res.ok) {
    if (data && typeof data === "object" && "code" in data) {
      throw new ApiError(res.status, data.code, data.title, data.detail);
    }
    // The proxy answers 502/504 itself when the service is down.
    throw new ApiError(res.status, "http_error", `HTTP ${res.status}`, "Is the service running on port 8080?");
  }
  return data as T;
}

export const api = {
  createUser: (subject: string) => call<CreatedUser>("POST", "/users", { body: { subject } }),

  createOrganizer: (name: string, adminToken: string) =>
    call<CreatedOrganizer>("POST", "/organizers", { body: { name }, token: adminToken }),

  createEvent: (ev: NewEvent, organizerToken: string) =>
    call<Event>("POST", "/events", { body: ev, token: organizerToken, idempotent: true }),

  getEvent: (id: string) => call<Event>("GET", `/events/${id}`),

  // The box office mints the next ticket number to the buyer's wallet.
  buyTicket: (eventId: string, userToken: string) =>
    // The id may be pasted by hand: escape it so it stays one path segment.
    call<Ticket>("POST", `/events/${encodeURIComponent(eventId)}/tickets`, { token: userToken, idempotent: true }),

  getTicket: (id: string) => call<Ticket>("GET", `/tickets/${id}`),

  // Resale: each is a transaction from the user's own wallet.
  listTicket: (id: string, price: string, holderToken: string) =>
    call<ResaleAction>("PUT", `/tickets/${id}/listing`, { body: { price }, token: holderToken, idempotent: true }),

  cancelListing: (id: string, holderToken: string) =>
    call<ResaleAction>("DELETE", `/tickets/${id}/listing`, { token: holderToken, idempotent: true }),

  buyResale: (id: string, buyerToken: string, maxPrice: string) =>
    call<ResaleAction>("POST", `/tickets/${id}/purchase`, {
      body: { max_price: maxPrice },
      token: buyerToken,
      idempotent: true,
    }),

  // Check-in: the holder signs a pass (nothing goes on-chain), the gate sends it.
  createPass: (ticketId: string, holderToken: string) =>
    call<CheckInPass>("POST", `/tickets/${ticketId}/pass`, { token: holderToken }),

  checkIn: (ticketId: string, pass: CheckInPass, organizerToken: string) =>
    call<Ticket>("POST", `/tickets/${ticketId}/check-in`, { body: pass, token: organizerToken, idempotent: true }),

  // Readable only by the caller who started the operation.
  getOperation: (id: string, token: string) => call<Operation>("GET", `/operations/${id}`, { token }),
};
