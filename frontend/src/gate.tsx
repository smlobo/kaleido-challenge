import jsQR from "jsqr";
import { useEffect, useRef, useState } from "react";
import { api, type CheckInPass, type CreatedOrganizer, type Event } from "./api";
import { mmss, useSecondsLeft } from "./pass";
import type { TicketRow } from "./tickets";
import { Card, ErrorLine } from "./ui";

/** Reads a pass from scanned or pasted text; null if it isn't one. */
function parsePass(text: string): CheckInPass | null {
  try {
    const p = JSON.parse(text);
    if (p && typeof p.ticket_id === "string" && typeof p.deadline === "number" && typeof p.signature === "string") {
      return p as CheckInPass;
    }
  } catch {
    // not JSON
  }
  return null;
}

/**
 * The door. The gate scans the holder's QR code (or takes it pasted), works
 * out which organizer's gate it is from the pass's contract address, and
 * checks the ticket in with that organizer's token. The contract burns the
 * ticket only if the gate role sends it and the ticket's holder signed the pass.
 */
export function Gate({
  orgs,
  events,
  setRows,
  incoming,
  refused,
}: {
  orgs: CreatedOrganizer[];
  events: Event[];
  setRows: (f: (prev: TicketRow[]) => TicketRow[]) => void;
  incoming: { json: string; at: number } | null;
  refused: boolean;
}) {
  const [text, setText] = useState("");
  const [scanning, setScanning] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const [done, setDone] = useState<string | null>(null);
  // With refused actions shown: which organizer's gate to use, whatever the event.
  const [gateId, setGateId] = useState("");

  // A pass handed over from the holder's view.
  useEffect(() => {
    if (incoming) {
      setText(incoming.json);
      setDone(null);
      setError(null);
    }
  }, [incoming]);

  const pass = parsePass(text.trim());
  const event = pass && events.find((e) => e.contract_address?.toLowerCase() === pass.contract_address?.toLowerCase());
  const rightOrg = event && orgs.find((o) => o.id === event.organizer_id);
  const org = (refused && orgs.find((o) => o.id === gateId)) || rightOrg;

  /** Changes the pass in the box, for showing what the service and the contract refuse. */
  const tamper = (f: (p: CheckInPass) => CheckInPass) => {
    if (pass) {
      setText(JSON.stringify(f({ ...pass })));
      setDone(null);
      setError(null);
    }
  };

  const checkIn = async () => {
    if (!pass || !org) return;
    setBusy(true);
    setError(null);
    setDone(null);
    try {
      const ticket = await api.checkIn(pass.ticket_id, pass, org.token);
      const pending = { kind: "checkin" as const, opId: ticket.checkin_operation_id!, token: org.token, status: "signed" as const };
      setRows((prev) =>
        prev.some((r) => r.ticket.id === ticket.id)
          ? prev.map((r) => (r.ticket.id === ticket.id ? { ...r, ticket, pending } : r))
          : [{ ticket, pending }, ...prev],
      );
      setDone(
        `${event?.name ?? "Ticket"} #${pass.token_id}: check-in sent. The ticket is burned when the block is mined; watch its row above.` +
          (refused ? " The pass stays in the box: check in again to try replaying it." : ""),
      );
      if (!refused) setText("");
    } catch (err) {
      setError(err);
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card
      step="5"
      title="Gate"
      intro={
        <>
          At the door the holder shows a pass (<strong>Show pass</strong> on a ticket above): a signature from their
          wallet, as a QR code. The gate scans it and sends <code>POST /tickets/{"{id}"}/check-in</code> with the
          organizer's token. On-chain, the gate key sends the transaction and the contract checks the holder signed
          the pass, then burns the ticket.
        </>
      }
    >
      <div className="row">
        <button type="button" onClick={() => setScanning((s) => !s)}>
          {scanning ? "Stop camera" : "Scan with camera"}
        </button>
        <ImageScan onScan={(t) => setText(t)} onError={setError} />
      </div>
      {scanning && (
        <Camera
          onScan={(t) => {
            setText(t);
            setScanning(false);
          }}
          onError={(e) => {
            setError(e);
            setScanning(false);
          }}
        />
      )}
      <label>
        Pass (scanned, handed over, or pasted)
        <textarea
          value={text}
          onChange={(e) => setText(e.target.value)}
          rows={3}
          placeholder='{"ticket_id":"…","deadline":…,"signature":"0x…"}'
          spellCheck={false}
        />
      </label>
      {text.trim() && !pass && <p className="error">That isn't a check-in pass.</p>}
      {refused && pass && (
        <div className="row tamper">
          <span className="muted">Try:</span>
          <button type="button" className="small ghost" onClick={() => tamper((p) => ({ ...p, deadline: Math.floor(Date.now() / 1000) - 60 }))}>
            Expire it
          </button>
          <button type="button" className="small ghost" onClick={() => tamper((p) => ({ ...p, deadline: p.deadline + 600 }))}>
            Extend the deadline
          </button>
          <button
            type="button"
            className="small ghost"
            onClick={() =>
              tamper((p) => {
                // Change one hex digit of s (in the middle): still well-formed, signed by nobody.
                const i = 80;
                const c = p.signature[i] === "0" ? "1" : "0";
                return { ...p, signature: p.signature.slice(0, i) + c + p.signature.slice(i + 1) };
              })
            }
          >
            Forge the signature
          </button>
          <label className="inline">
            Gate of
            <select value={gateId} onChange={(e) => setGateId(e.target.value)}>
              <option value="">the event's organizer</option>
              {orgs
                .filter((o) => o.id !== rightOrg?.id)
                .map((o) => (
                  <option key={o.id} value={o.id}>
                    {o.name} (wrong gate)
                  </option>
                ))}
            </select>
          </label>
        </div>
      )}
      {pass && (
        <PassSummary
          pass={pass}
          eventName={event?.name}
          gate={org?.name}
          onCheckIn={checkIn}
          busy={busy}
        />
      )}
      <ErrorLine error={error} />
      {done && <p className="ok">{done}</p>}
    </Card>
  );
}

function PassSummary({
  pass,
  eventName,
  gate,
  onCheckIn,
  busy,
}: {
  pass: CheckInPass;
  eventName?: string;
  gate?: string;
  onCheckIn: () => void;
  busy: boolean;
}) {
  const left = useSecondsLeft(pass.deadline);
  return (
    <div className="scanned">
      <dl>
        <dt>Ticket</dt>
        <dd>
          {eventName ?? <span className="muted">an event this page didn't create</span>} #{pass.token_id}
        </dd>
        <dt>Gate</dt>
        <dd>{gate ?? <span className="muted">no organizer on this page runs that event</span>}</dd>
        <dt>Valid</dt>
        <dd className={left === 0 ? "bad" : undefined}>{left === 0 ? "expired" : `${mmss(left)} left`}</dd>
      </dl>
      <button type="button" onClick={onCheckIn} disabled={busy || !gate}>
        {busy ? "Checking in…" : "Check in"}
      </button>
    </div>
  );
}

/** Reads QR codes from the camera until one decodes. */
function Camera({ onScan, onError }: { onScan: (text: string) => void; onError: (e: unknown) => void }) {
  const video = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    let stream: MediaStream | undefined;
    let frame = 0;
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d", { willReadFrequently: true });

    const tick = () => {
      const v = video.current;
      if (v && ctx && v.readyState === v.HAVE_ENOUGH_DATA) {
        canvas.width = v.videoWidth;
        canvas.height = v.videoHeight;
        ctx.drawImage(v, 0, 0);
        const code = jsQR(ctx.getImageData(0, 0, canvas.width, canvas.height).data, canvas.width, canvas.height);
        if (code?.data) {
          onScan(code.data);
          return;
        }
      }
      frame = requestAnimationFrame(tick);
    };

    navigator.mediaDevices
      ?.getUserMedia({ video: { facingMode: "environment" } })
      .then((s) => {
        stream = s;
        if (video.current) {
          video.current.srcObject = s;
          void video.current.play();
        }
        frame = requestAnimationFrame(tick);
      })
      .catch((e) => onError(new Error(`Camera unavailable: ${e instanceof Error ? e.message : e}`)));
    if (!navigator.mediaDevices) onError(new Error("No camera access in this browser."));

    return () => {
      cancelAnimationFrame(frame);
      stream?.getTracks().forEach((t) => t.stop());
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return <video ref={video} className="camera" muted playsInline />;
}

/** Reads a QR code from an image file, e.g. a screenshot of the pass. */
function ImageScan({ onScan, onError }: { onScan: (text: string) => void; onError: (e: unknown) => void }) {
  const read = async (file: File) => {
    try {
      const bmp = await createImageBitmap(file);
      const canvas = document.createElement("canvas");
      canvas.width = bmp.width;
      canvas.height = bmp.height;
      const ctx = canvas.getContext("2d")!;
      ctx.drawImage(bmp, 0, 0);
      const code = jsQR(ctx.getImageData(0, 0, bmp.width, bmp.height).data, bmp.width, bmp.height);
      if (!code?.data) throw new Error("No QR code found in that image.");
      onScan(code.data);
    } catch (e) {
      onError(e);
    }
  };
  return (
    <label className="file">
      <input type="file" accept="image/*" onChange={(e) => e.target.files?.[0] && read(e.target.files[0])} />
      <span className="button-like">Scan an image</span>
    </label>
  );
}
