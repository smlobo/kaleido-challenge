import { useEffect, useState, type FormEvent } from "react";
import { api, isFinal, type CheckInPass, type CreatedOrganizer, type CreatedUser, type Event, type Operation, type Ticket } from "./api";
import { AddExisting, existingTicket } from "./existing";
import { PassView } from "./pass";
import { Card, Copyable, Empty, ErrorLine } from "./ui";

export type Kind = "mint" | "list" | "cancel" | "buy" | "checkin";

/** A ticket as the page tracks it: the API's view, plus the transaction in flight. */
export interface TicketRow {
  ticket: Ticket;
  /** The transaction in flight for this ticket, polled until it is final. */
  pending?: { kind: Kind; opId: string; token: string; status: Operation["status"] };
  /** How the last transaction ended. */
  last?: { kind: Kind; status: Operation["status"]; detail?: string | null; block?: number | null };
}

const KIND_LABEL: Record<Kind, string> = { mint: "mint", list: "list", cancel: "cancel listing", buy: "resale", checkin: "check-in" };

/** Prices are in base units; 6 decimals, like USDC, so 110000000 = 110.00. */
const formatPrice = (base: string) => {
  try {
    const n = BigInt(base);
    return `${n / 1_000_000n}.${(n % 1_000_000n).toString().padStart(6, "0").slice(0, 2)}`;
  } catch {
    return base;
  }
};

/** The highest listing price the contract accepts: face value × max_resale_bps / 10000. */
const resaleCap = (ev?: Event) =>
  ev ? ((BigInt(ev.face_value) * BigInt(ev.max_resale_bps)) / 10_000n).toString() : "";

// Buy-form choices beyond the page's own users and on-sale events, shown with refused actions.
const PASTE = "__paste__";
const NO_TOKEN = "__none__";
const ORG = "org:";

export function Tickets({
  users,
  orgs,
  events,
  rows,
  setRows,
  onExisting,
  onHandToGate,
  refused,
}: {
  users: CreatedUser[];
  orgs: CreatedOrganizer[];
  events: Event[];
  rows: TicketRow[];
  setRows: (f: (prev: TicketRow[]) => TicketRow[]) => void;
  onExisting: (t: Ticket) => Promise<void>;
  onHandToGate: (passJson: string) => void;
  refused: boolean;
}) {
  const onSale = events.filter((e) => e.status === "on_sale");
  const [eventId, setEventId] = useState("");
  const [buyerId, setBuyerId] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const [shown, setShown] = useState<{ pass: CheckInPass; title: string; holder: string } | null>(null);

  const [pastedEvent, setPastedEvent] = useState("");
  const [pastedToken, setPastedToken] = useState("");

  // With refused actions shown, any event (deploying, failed, or a pasted id)
  // and any token (an organizer's, a pasted one, or none) can be tried.
  const eventChoices = refused ? events : onSale;
  const event = eventChoices.find((e) => e.id === eventId) ?? eventChoices[0];
  // What each dropdown shows: the choice made, else its first option. The
  // paste options are first when there is nothing else to pick.
  const eventSel = refused && (eventId === PASTE || !event) ? PASTE : (event?.id ?? "");
  const special = (v: string) => v === PASTE || v === NO_TOKEN || orgs.some((o) => ORG + o.id === v);
  const buyer = users.find((u) => u.id === buyerId) ?? users[0];
  const buyerSel = refused && (special(buyerId) || !buyer) ? (special(buyerId) ? buyerId : orgs[0] ? ORG + orgs[0].id : PASTE) : (buyer?.id ?? "");
  const buyOrg = orgs.find((o) => ORG + o.id === buyerSel);

  const targetEventId = eventSel === PASTE ? pastedEvent.trim() : event?.id;
  const buyerToken =
    buyerSel === PASTE ? pastedToken.trim() : buyerSel === NO_TOKEN ? "" : (buyOrg?.token ?? buyer?.token);

  const update = (id: string, f: (r: TicketRow) => TicketRow) =>
    setRows((prev) => prev.map((r) => (r.ticket.id === id ? f(r) : r)));

  // Poll every transaction in flight until it is final, then reload its ticket.
  const inFlight = rows
    .filter((r) => r.pending)
    .map((r) => `${r.ticket.id}:${r.pending!.opId}`)
    .join(",");
  useEffect(() => {
    if (!inFlight) return;
    const t = setInterval(() => {
      for (const r of rows.filter((r) => r.pending)) {
        const p = r.pending!;
        api
          .getOperation(p.opId, p.token)
          .then(async (op) => {
            if (!isFinal(op.status)) {
              if (op.status !== p.status) update(r.ticket.id, (x) => ({ ...x, pending: { ...p, status: op.status } }));
              return;
            }
            const ticket = await api.getTicket(r.ticket.id);
            update(r.ticket.id, () => ({
              ticket,
              last: { kind: p.kind, status: op.status, detail: op.revert_reason ?? op.last_error, block: op.block_number },
            }));
          })
          .catch(() => {
            // try again next tick
          });
      }
    }, 1500);
    return () => clearInterval(t);
  }, [inFlight]); // eslint-disable-line react-hooks/exhaustive-deps

  // A ticket added from elsewhere while its mint was in flight has no
  // operation to poll (only its buyer may read that): re-read the ticket.
  const minting = rows
    .filter((r) => !r.pending && r.ticket.status === "pending_mint")
    .map((r) => r.ticket.id)
    .join(",");
  useEffect(() => {
    if (!minting) return;
    const t = setInterval(() => {
      for (const id of minting.split(",")) {
        api
          .getTicket(id)
          .then((ticket) => ticket.status !== "pending_mint" && update(id, (r) => ({ ...r, ticket })))
          .catch(() => {
            // try again next tick
          });
      }
    }, 2000);
    return () => clearInterval(t);
  }, [minting]); // eslint-disable-line react-hooks/exhaustive-deps

  const buy = async (e: FormEvent) => {
    e.preventDefault();
    if (!targetEventId || buyerToken === undefined) return;
    setBusy(true);
    setError(null);
    try {
      const ticket = await api.buyTicket(targetEventId, buyerToken);
      setRows((prev) => [
        { ticket, pending: { kind: "mint", opId: ticket.mint_operation_id, token: buyerToken, status: "signed" } },
        ...prev,
      ]);
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  };

  /** Starts a resale action and tracks its transaction. */
  const act = async (row: TicketRow, kind: Kind, token: string, run: () => ReturnType<typeof api.listTicket>) => {
    const res = await run();
    update(row.ticket.id, (r) => ({
      ...r,
      ticket: res.ticket,
      pending: { kind, opId: res.operation_id, token, status: "signed" },
    }));
  };

  return (
    <Card
      step="4"
      title="Tickets"
      intro={
        <>
          <code>POST /events/{"{id}"}/tickets</code> mints a ticket to the buyer's wallet, from the organizer's
          box office. Resale is three calls, each a transaction from the user's own wallet:{" "}
          <code>PUT …/listing</code>, <code>DELETE …/listing</code> and <code>POST …/purchase</code>. The contract
          checks who acts and keeps the price under the event's cap. No money moves yet.
        </>
      }
    >
      <form onSubmit={buy} className="row">
        <label>
          Event
          <select
            value={eventSel}
            onChange={(e) => setEventId(e.target.value)}
            disabled={!refused && onSale.length === 0}
          >
            {eventChoices.length === 0 && !refused && <option value="">No event on sale yet</option>}
            {eventChoices.map((ev) => (
              <option key={ev.id} value={ev.id}>
                {ev.name}
                {ev.status !== "on_sale" ? ` (${ev.status})` : ""}
              </option>
            ))}
            {refused && <option value={PASTE}>Paste an event id…</option>}
          </select>
        </label>
        {eventSel === PASTE && (
          <label>
            Event id
            <input
              value={pastedEvent}
              onChange={(e) => setPastedEvent(e.target.value)}
              placeholder="any uuid, or not one"
              spellCheck={false}
            />
          </label>
        )}
        <label>
          Buyer
          <select
            value={buyerSel}
            onChange={(e) => setBuyerId(e.target.value)}
            disabled={!refused && users.length === 0}
          >
            {users.length === 0 && !refused && <option value="">Create a user first</option>}
            {users.map((u) => (
              <option key={u.id} value={u.id}>
                {u.subject}
              </option>
            ))}
            {refused &&
              orgs.map((o) => (
                <option key={o.id} value={ORG + o.id}>
                  {o.name} (organizer's token)
                </option>
              ))}
            {refused && <option value={PASTE}>Paste a token…</option>}
            {refused && <option value={NO_TOKEN}>No token</option>}
          </select>
        </label>
        {buyerSel === PASTE && (
          <label>
            Buyer token
            <input
              value={pastedToken}
              onChange={(e) => setPastedToken(e.target.value)}
              placeholder="any token, or not one"
              autoComplete="off"
              spellCheck={false}
            />
          </label>
        )}
        <button type="submit" disabled={busy || !targetEventId || buyerToken === undefined}>
          {busy ? "Buying…" : "Buy ticket"}
        </button>
      </form>
      <ErrorLine error={error} />
      <AddExisting what="a ticket" add={async (id) => onExisting(await existingTicket(id))} />
      {rows.length === 0 ? (
        <Empty>No tickets yet.</Empty>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Ticket</th>
                <th>Status</th>
                <th>Owner</th>
                <th>Listed at</th>
                <th>Last transaction</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                <TicketLine
                  key={r.ticket.id}
                  row={r}
                  users={users}
                  events={events}
                  act={act}
                  showPass={(pass, title, holder) => setShown({ pass, title, holder })}
                  refused={refused}
                />
              ))}
            </tbody>
          </table>
        </div>
      )}
      {shown && (
        <PassView
          {...shown}
          onClose={() => setShown(null)}
          onHandToGate={() => {
            onHandToGate(JSON.stringify(shown.pass));
            setShown(null);
          }}
        />
      )}
    </Card>
  );
}

function TicketLine({
  row,
  users,
  events,
  act,
  showPass,
  refused,
}: {
  row: TicketRow;
  users: CreatedUser[];
  events: Event[];
  act: (row: TicketRow, kind: Kind, token: string, run: () => ReturnType<typeof api.listTicket>) => Promise<void>;
  showPass: (pass: CheckInPass, title: string, holder: string) => void;
  refused: boolean;
}) {
  const { ticket: t, pending, last } = row;
  const ev = events.find((e) => e.id === t.event_id);
  const owner = users.find((u) => t.owner_address && u.wallet_address.toLowerCase() === t.owner_address.toLowerCase());
  const others = users.filter((u) => u.id !== owner?.id);
  const cap = resaleCap(ev);
  const listed = !!t.listing_price;

  const [price, setPrice] = useState(cap);
  const [buyerId, setBuyerId] = useState("");
  const [actorId, setActorId] = useState("");
  const [maxPrice, setMaxPrice] = useState("");
  const [error, setError] = useState<unknown>(null);
  // The buyer picked in the row: another user by default, or the owner to show the refusal.
  const buyer = users.find((u) => u.id === buyerId) ?? others[0] ?? owner;
  // With refused actions shown, any user can try any action; the owner by default.
  const actor = users.find((u) => u.id === actorId) ?? owner ?? users[0];

  const run = async (f: () => Promise<void>) => {
    setError(null);
    try {
      await f();
    } catch (err) {
      setError(err);
    }
  };

  const canAct = t.status === "active" && !pending;
  const title = `${ev?.name ?? "Ticket"} #${t.token_id}`;

  const list = (u: CreatedUser) => run(() => act(row, "list", u.token, () => api.listTicket(t.id, price.trim(), u.token)));
  const cancel = (u: CreatedUser) => run(() => act(row, "cancel", u.token, () => api.cancelListing(t.id, u.token)));
  const pass = (u: CreatedUser) => run(async () => showPass(await api.createPass(t.id, u.token), title, u.subject));
  const buy = (u: CreatedUser, max: string) => run(() => act(row, "buy", u.token, () => api.buyResale(t.id, u.token, max)));

  const priceInput = (
    <input
      className="price"
      value={price}
      onChange={(e) => setPrice(e.target.value)}
      inputMode="numeric"
      aria-label="Asking price in base units"
      title={`At most ${cap} (the resale cap)`}
    />
  );

  return (
    <>
      <tr>
        <td>
          {ev?.name ?? "event"} <span className="muted">#{t.token_id}</span>
          <div><Copyable value={t.id} /></div>
        </td>
        <td>
          <span className={`pill ${t.status}`}>{t.status.replace("_", " ")}</span>
        </td>
        <td>
          {owner ? owner.subject : t.owner_address ? <Copyable value={t.owner_address} /> : <span className="muted">—</span>}
        </td>
        <td>
          {listed ? (
            <>
              {formatPrice(t.listing_price!)} <span className="muted">({t.listing_price})</span>
            </>
          ) : (
            <span className="muted">not listed</span>
          )}
        </td>
        <td>
          {pending ? (
            <span className={`pill ${pending.status}`}>
              {KIND_LABEL[pending.kind]} · {pending.status}…
            </span>
          ) : last ? (
            <span className={`pill ${last.status}`} title={last.detail ?? undefined}>
              {KIND_LABEL[last.kind]} · {last.status}
              {last.block ? ` in block ${last.block}` : ""}
            </span>
          ) : (
            <span className="muted">—</span>
          )}
        </td>
        {refused && actor ? (
          // Every action, whatever the ticket's state; the service decides.
          <td className="actions refused">
            <span className="action">
              <span className="muted">as</span>
              <select value={actor.id} onChange={(e) => setActorId(e.target.value)} aria-label="Act as user">
                {users.map((u) => (
                  <option key={u.id} value={u.id}>
                    {u.subject}
                    {u.id === owner?.id ? " (owner)" : ""}
                  </option>
                ))}
              </select>
            </span>
            <span className="action">
              {priceInput}
              <button type="button" className="small" onClick={() => list(actor)}>
                List
              </button>
            </span>
            <button type="button" className="small ghost" onClick={() => cancel(actor)}>
              Cancel
            </button>
            <span className="action">
              <input
                className="price"
                value={maxPrice}
                onChange={(e) => setMaxPrice(e.target.value)}
                placeholder={t.listing_price ?? "max price"}
                inputMode="numeric"
                aria-label="Most the buyer agrees to pay"
                title="max_price: the most the buyer agrees to pay (default: the listed price)"
              />
              <button type="button" className="small" onClick={() => buy(actor, maxPrice.trim() || t.listing_price || cap)}>
                Buy
              </button>
            </span>
            <button type="button" className="small ghost" onClick={() => pass(actor)}>
              Show pass
            </button>
          </td>
        ) : (
          <td className="actions">
            {canAct && owner && (
              <span className="action">
                {priceInput}
                <button type="button" className="small" onClick={() => list(owner)}>
                  {listed ? "Re-price" : "List"}
                </button>
              </span>
            )}
            {canAct && owner && listed && (
              <button type="button" className="small ghost" onClick={() => cancel(owner)}>
                Cancel
              </button>
            )}
            {canAct && owner && (
              <button type="button" className="small ghost" onClick={() => pass(owner)}>
                Show pass
              </button>
            )}
            {canAct && listed && buyer && (
              <span className="action">
                <select value={buyer.id} onChange={(e) => setBuyerId(e.target.value)} aria-label="Buyer">
                  {others.map((u) => (
                    <option key={u.id} value={u.id}>
                      {u.subject}
                    </option>
                  ))}
                  {owner && (
                    // The owner can try too: the service refuses with 422 buyer_is_seller.
                    <option value={owner.id}>{owner.subject} (owner)</option>
                  )}
                </select>
                <button type="button" className="small" onClick={() => buy(buyer, t.listing_price!)}>
                  Buy
                </button>
              </span>
            )}
            {!canAct && <span className="muted">{pending ? "waiting for the chain…" : "—"}</span>}
          </td>
        )}
      </tr>
      {error ? (
        <tr className="row-error">
          <td colSpan={6}>
            <ErrorLine error={error} />
          </td>
        </tr>
      ) : null}
    </>
  );
}
