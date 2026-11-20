import { useEffect, useRef, useState, type FormEvent } from "react";
import { api, type CreatedOrganizer, type CreatedUser, type Event } from "./api";
import { clearPersisted, usePersisted } from "./store";
import { ErrorPaths } from "./errorpaths";
import { Gate } from "./gate";
import { Tickets, type TicketRow } from "./tickets";
import { Card, Copyable, Empty, ErrorLine } from "./ui";

export default function App() {
  const [users, setUsers] = usePersisted<CreatedUser[]>("users", []);
  const [orgs, setOrgs] = usePersisted<CreatedOrganizer[]>("organizers", []);
  const [events, setEvents] = usePersisted<Event[]>("events", []);
  const [tickets, setTickets] = usePersisted<TicketRow[]>("tickets", []);
  // A pass the holder's view handed to the gate (at: so the same pass twice still lands).
  const [handed, setHanded] = useState<{ json: string; at: number } | null>(null);
  const gateRef = useRef<HTMLDivElement>(null);
  // Demo mode: offer every action, so the API's refusals can be shown.
  const [refused, setRefused] = usePersisted("show-refused", false);

  const clear = () => {
    if (!window.confirm("Forget every user, organizer, event and ticket this page has created? (They stay in the service.)")) return;
    clearPersisted();
    setUsers([]);
    setOrgs([]);
    setEvents([]);
    setTickets([]);
  };

  return (
    <div className="page">
      <header className="top">
        <img src="/favicon.svg" alt="" width="40" height="40" />
        <div>
          <h1>Ticketing demo</h1>
          <p>Tickets as ERC-721 tokens on a local Besu chain. Each form calls the service's REST API.</p>
        </div>
        <label className="toggle">
          <input type="checkbox" checked={refused} onChange={(e) => setRefused(e.target.checked)} />
          Show refused actions
        </label>
        <button type="button" className="ghost" onClick={clear}>
          Clear
        </button>
      </header>
      {refused && <ErrorPaths />}

      <Users users={users} onCreated={(u) => setUsers((prev) => [u, ...prev])} />
      <Organizers orgs={orgs} onCreated={(o) => setOrgs((prev) => [o, ...prev])} />
      <Events
        orgs={orgs}
        events={events}
        onCreated={(e) => setEvents((prev) => [e, ...prev])}
        onUpdated={(e) => setEvents((prev) => prev.map((x) => (x.id === e.id ? e : x)))}
      />
      <Tickets
        users={users}
        orgs={orgs}
        events={events}
        rows={tickets}
        setRows={setTickets}
        refused={refused}
        onHandToGate={(json) => {
          setHanded({ json, at: Date.now() });
          gateRef.current?.scrollIntoView({ behavior: "smooth", block: "start" });
        }}
      />
      <div ref={gateRef}>
        <Gate orgs={orgs} events={events} setRows={setTickets} incoming={handed} refused={refused} />
      </div>
    </div>
  );
}

function Users({ users, onCreated }: { users: CreatedUser[]; onCreated: (u: CreatedUser) => void }) {
  const [name, setName] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>(null);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      onCreated(await api.createUser(name.trim()));
      setName("");
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card
      step="1"
      title="Users"
      intro={
        <>
          <code>POST /users</code>. The service creates a custodial wallet for each user: a key it holds,
          encrypted, in Postgres. The token is shown once.
        </>
      }
    >
      <form onSubmit={submit} className="row">
        <label>
          Name
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="alice" required maxLength={200} />
        </label>
        <button type="submit" disabled={busy || !name.trim()}>
          {busy ? "Creating…" : "Create user"}
        </button>
      </form>
      <ErrorLine error={error} />
      {users.length === 0 ? (
        <Empty>No users yet.</Empty>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>API token</th>
                <th>Postgres id</th>
                <th>Wallet address</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id}>
                  <td>{u.subject}</td>
                  <td><Copyable value={u.token} /></td>
                  <td><Copyable value={u.id} /></td>
                  <td><Copyable value={u.wallet_address} /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Card>
  );
}

function Organizers({ orgs, onCreated }: { orgs: CreatedOrganizer[]; onCreated: (o: CreatedOrganizer) => void }) {
  const [name, setName] = useState("");
  const [adminToken, setAdminToken] = usePersisted("admin-token", "");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>(null);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      onCreated(await api.createOrganizer(name.trim(), adminToken.trim()));
      setName("");
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card
      step="2"
      title="Organizers"
      intro={
        <>
          <code>POST /organizers</code>, with the platform admin's token (the service's <code>ADMIN_TOKEN</code>).
          Each organizer gets its own keys: the box office mints tickets, the gate checks them in.
        </>
      }
    >
      <form onSubmit={submit} className="row">
        <label>
          Name
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Rock Fest Ltd" required maxLength={200} />
        </label>
        <label>
          Admin token
          <input
            value={adminToken}
            onChange={(e) => setAdminToken(e.target.value)}
            placeholder="ADMIN_TOKEN"
            type="password"
            autoComplete="off"
            required
          />
        </label>
        <button type="submit" disabled={busy || !name.trim() || !adminToken.trim()}>
          {busy ? "Creating…" : "Create organizer"}
        </button>
      </form>
      <ErrorLine error={error} />
      {orgs.length === 0 ? (
        <Empty>No organizers yet.</Empty>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>API token</th>
                <th>Postgres id</th>
                <th>Box office address</th>
                <th>Gate address</th>
              </tr>
            </thead>
            <tbody>
              {orgs.map((o) => (
                <tr key={o.id}>
                  <td>{o.name}</td>
                  <td><Copyable value={o.token} /></td>
                  <td><Copyable value={o.id} /></td>
                  <td><Copyable value={o.box_office_address} /></td>
                  <td><Copyable value={o.gate_address} /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Card>
  );
}

// Mirrors create_event.py: capital letters and digits from the name.
const defaultSymbol = (name: string) => name.toUpperCase().replace(/[^A-Z0-9]/g, "").slice(0, 11) || "EVENT";

const PASTED = "__pasted__";

function Events({
  orgs,
  events,
  onCreated,
  onUpdated,
}: {
  orgs: CreatedOrganizer[];
  events: Event[];
  onCreated: (e: Event) => void;
  onUpdated: (e: Event) => void;
}) {
  const [name, setName] = useState("");
  const [orgChoice, setOrgChoice] = useState<string>(orgs[0]?.id ?? PASTED);
  const [pastedToken, setPastedToken] = useState("");
  const [symbol, setSymbol] = useState("");
  const [startsInHours, setStartsInHours] = useState(24);
  const [supply, setSupply] = useState(100);
  const [faceValue, setFaceValue] = useState("100000000");
  const [resaleBps, setResaleBps] = useState(11000);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>(null);

  // A new organizer created above becomes the default choice.
  useEffect(() => {
    if (orgs.length > 0 && (orgChoice === PASTED ? !pastedToken : !orgs.some((o) => o.id === orgChoice))) {
      setOrgChoice(orgs[0].id);
    }
  }, [orgs]); // eslint-disable-line react-hooks/exhaustive-deps

  const token = orgChoice === PASTED ? pastedToken.trim() : (orgs.find((o) => o.id === orgChoice)?.token ?? "");

  // Poll each deploying event until its contract is mined (or the deploy fails).
  const deploying = events.filter((e) => e.status === "deploying").map((e) => e.id).join(",");
  useEffect(() => {
    if (!deploying) return;
    const t = setInterval(async () => {
      for (const id of deploying.split(",")) {
        try {
          const ev = await api.getEvent(id);
          if (ev.status !== "deploying") onUpdated(ev);
        } catch {
          // try again next tick
        }
      }
    }, 2000);
    return () => clearInterval(t);
  }, [deploying]); // eslint-disable-line react-hooks/exhaustive-deps

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      const ev = await api.createEvent(
        {
          name: name.trim(),
          symbol: symbol.trim() || defaultSymbol(name),
          starts_at: new Date(Date.now() + startsInHours * 3600_000).toISOString(),
          max_supply: supply,
          face_value: faceValue.trim(),
          max_resale_bps: resaleBps,
        },
        token,
      );
      onCreated(ev);
      setName("");
      setSymbol("");
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  };

  const orgName = (id: string) => orgs.find((o) => o.id === id)?.name ?? `${id.slice(0, 8)}…`;

  return (
    <Card
      step="3"
      title="Events"
      intro={
        <>
          <code>POST /events</code>, with an organizer's token. The service deploys a <code>TicketNFT</code>{" "}
          contract for the event and answers at once; the contract address appears when the deploy is mined,
          a block or two later.
        </>
      }
    >
      <form onSubmit={submit} className="row">
        <label>
          Name
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Rock Fest 2027" required maxLength={100} />
        </label>
        <label>
          Organizer
          <select value={orgChoice} onChange={(e) => setOrgChoice(e.target.value)}>
            {orgs.map((o) => (
              <option key={o.id} value={o.id}>
                {o.name}
              </option>
            ))}
            <option value={PASTED}>Paste a token…</option>
          </select>
        </label>
        {orgChoice === PASTED && (
          <label>
            Organizer token
            <input
              value={pastedToken}
              onChange={(e) => setPastedToken(e.target.value)}
              type="password"
              autoComplete="off"
              required
            />
          </label>
        )}
        <button type="submit" disabled={busy || !name.trim() || !token}>
          {busy ? "Deploying…" : "Create event"}
        </button>
      </form>
      <details className="more">
        <summary>More settings</summary>
        <div className="row">
          <label>
            Symbol
            <input
              value={symbol}
              onChange={(e) => setSymbol(e.target.value.toUpperCase())}
              placeholder={name ? defaultSymbol(name) : "ROCK2027"}
              pattern="[A-Z0-9]{1,11}"
              maxLength={11}
            />
          </label>
          <label>
            Starts in (hours)
            <input type="number" min={1} value={startsInHours} onChange={(e) => setStartsInHours(Number(e.target.value))} />
          </label>
          <label>
            Supply
            <input type="number" min={1} value={supply} onChange={(e) => setSupply(Number(e.target.value))} />
          </label>
          <label>
            Face value (base units)
            <input value={faceValue} onChange={(e) => setFaceValue(e.target.value)} inputMode="numeric" pattern="[0-9]+" />
          </label>
          <label>
            Resale cap (bps)
            <input type="number" min={0} value={resaleBps} onChange={(e) => setResaleBps(Number(e.target.value))} />
          </label>
        </div>
      </details>
      <ErrorLine error={error} />
      {events.length === 0 ? (
        <Empty>No events yet.</Empty>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Organizer</th>
                <th>Status</th>
                <th>Postgres id</th>
                <th>Contract address</th>
              </tr>
            </thead>
            <tbody>
              {events.map((ev) => (
                <tr key={ev.id}>
                  <td>
                    {ev.name} <span className="muted">({ev.symbol})</span>
                  </td>
                  <td>{orgName(ev.organizer_id)}</td>
                  <td>
                    <span className={`pill ${ev.status}`}>{ev.status.replace("_", " ")}</span>
                  </td>
                  <td><Copyable value={ev.id} /></td>
                  <td>{ev.contract_address ? <Copyable value={ev.contract_address} /> : <span className="muted">waiting for the deploy…</span>}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Card>
  );
}
