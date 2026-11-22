import { useState, type FormEvent } from "react";
import { api, type CreatedOrganizer, type CreatedUser, type Event, type Ticket } from "./api";
import { ErrorLine } from "./ui";

// Adding what was created elsewhere (curl, the Python scripts, another
// browser), so the page can show it and act with it.
//
// Events and tickets are public: their id is enough. Users and organizers are
// not readable from the API, and a token is shown only when the account is
// created, so for those the page takes what creation printed: the JSON curl
// returned, or create_user.py / create_organizer.py's labeled lines.

const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
const ADDRESS = /^0x[0-9a-f]{40}$/i;

function checkId(id: string): string {
  const v = id.trim();
  if (!UUID.test(v)) throw new Error("That id is not a UUID.");
  return v.toLowerCase();
}

/** Replaces the item with the same id, or puts it first. */
export function upsert<T>(list: T[], item: T, id: (x: T) => string): T[] {
  return list.some((x) => id(x) === id(item)) ? list.map((x) => (id(x) === id(item) ? item : x)) : [item, ...list];
}

// The scripts' labels, and the JSON field each stands for.
const LABELS: Record<string, string> = {
  user: "id",
  organizer: "id",
  name: "name",
  wallet: "wallet_address",
  "box office": "box_office_address",
  gate: "gate_address",
  payout: "payout_address",
  token: "token",
};

/**
 * Reads a creation response: JSON from curl, or the scripts' lines such as
 * "wallet  0x…" and "box office  0x…  (mints tickets)". Also notes which kind
 * the script said it created, from its first label.
 */
function parseCreated(text: string): { fields: Record<string, string>; kind?: "user" | "organizer" } {
  const t = text.trim();
  if (t.startsWith("{")) {
    let obj: unknown;
    try {
      obj = JSON.parse(t);
    } catch {
      throw new Error("That looks like JSON, but it doesn't parse.");
    }
    const fields: Record<string, string> = {};
    for (const [k, v] of Object.entries(obj as Record<string, unknown>)) if (typeof v === "string") fields[k] = v;
    if ("subject" in fields) fields.name = fields.subject;
    const kind = "wallet_address" in fields ? "user" : "box_office_address" in fields ? "organizer" : undefined;
    return { fields, kind };
  }
  const fields: Record<string, string> = {};
  let kind: "user" | "organizer" | undefined;
  for (const line of t.split("\n")) {
    const m = line.trim().match(/^([a-z ]+?)\s{2,}(.+)$/i);
    if (!m) continue;
    const label = m[1].toLowerCase();
    if (!(label in LABELS)) continue;
    // A name may have spaces; anything else is one word, then a comment.
    fields[LABELS[label]] = label === "name" ? m[2].trim() : m[2].trim().split(/\s+/)[0];
    if (label === "user" || label === "organizer") kind = label;
  }
  return { fields, kind };
}

function need(text: string, fields: Record<string, string>, names: string[], what: string) {
  if (/^tk_\S+$/.test(text.trim())) {
    throw new Error(
      `That is only a token. Paste everything the ${what} was created with: the script's output without --token-only, or curl's JSON.`,
    );
  }
  const missing = names.filter((n) => !fields[n]);
  if (missing.length > 0) throw new Error(`Missing ${missing.join(", ")}. Paste everything the ${what} was created with.`);
}

const wrongKind = (want: string, got?: string) => {
  if (got && got !== want) throw new Error(`That is ${got === "organizer" ? "an organizer" : "a user"}, not ${want === "organizer" ? "an organizer" : "a user"}.`);
};

export function pastedUser(text: string): CreatedUser {
  const { fields, kind } = parseCreated(text);
  wrongKind("user", kind);
  need(text, fields, ["id", "name", "wallet_address", "token"], "user");
  if (!ADDRESS.test(fields.wallet_address)) throw new Error("The wallet address is not a 0x address.");
  return {
    id: checkId(fields.id),
    subject: fields.name,
    wallet_address: fields.wallet_address.toLowerCase(),
    token: fields.token,
  };
}

export function pastedOrganizer(text: string): CreatedOrganizer {
  const { fields, kind } = parseCreated(text);
  wrongKind("organizer", kind);
  const addrs = ["box_office_address", "gate_address", "payout_address"];
  need(text, fields, ["id", "name", ...addrs, "token"], "organizer");
  for (const a of addrs) if (!ADDRESS.test(fields[a])) throw new Error(`The ${a.replace(/_/g, " ")} is not a 0x address.`);
  return {
    id: checkId(fields.id),
    name: fields.name,
    box_office_address: fields.box_office_address.toLowerCase(),
    gate_address: fields.gate_address.toLowerCase(),
    payout_address: fields.payout_address.toLowerCase(),
    token: fields.token,
  };
}

export const existingEvent = (id: string): Promise<Event> => api.getEvent(checkId(id));

export const existingTicket = (id: string): Promise<Ticket> => api.getTicket(checkId(id));

/**
 * A collapsed "Add … created elsewhere" form: an id, or with `paste`, a box
 * for what creation printed. `add` looks it up or parses it, and stores it.
 */
export function AddExisting({
  what,
  paste,
  add,
}: {
  what: string;
  paste?: string;
  add: (value: string) => void | Promise<void>;
}) {
  const [value, setValue] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>(null);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await add(value);
      setValue("");
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  };

  return (
    <details className="more">
      <summary>Add {what} created elsewhere (curl, scripts)</summary>
      <form onSubmit={submit} className="row">
        {paste ? (
          <label className="wide">
            {paste}
            <textarea
              value={value}
              onChange={(e) => setValue(e.target.value)}
              rows={5}
              spellCheck={false}
              autoComplete="off"
              required
            />
          </label>
        ) : (
          <label>
            Id
            <input value={value} onChange={(e) => setValue(e.target.value)} placeholder="uuid" spellCheck={false} required />
          </label>
        )}
        <button type="submit" disabled={busy || !value.trim()}>
          {busy ? "Looking up…" : "Add"}
        </button>
      </form>
      <ErrorLine error={error} />
    </details>
  );
}
