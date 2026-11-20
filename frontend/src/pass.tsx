import QRCode from "qrcode";
import { useEffect, useRef, useState } from "react";
import type { CheckInPass } from "./api";
import { Copyable } from "./ui";

/** Seconds until a unix deadline, ticking once a second. */
export function useSecondsLeft(deadline: number) {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const t = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(t);
  }, []);
  return Math.max(0, Math.floor(deadline - now / 1000));
}

export const mmss = (s: number) => `${Math.floor(s / 60)}:${String(s % 60).padStart(2, "0")}`;

/**
 * The holder's phone at the door: the signed pass as a QR code, with how
 * long it stays valid. The QR code holds the pass exactly as the API returned
 * it, as JSON; the gate sends it back unchanged.
 */
export function PassView({
  pass,
  title,
  holder,
  onHandToGate,
  onClose,
}: {
  pass: CheckInPass;
  title: string;
  holder: string;
  onHandToGate: () => void;
  onClose: () => void;
}) {
  const canvas = useRef<HTMLCanvasElement>(null);
  const json = JSON.stringify(pass);
  const left = useSecondsLeft(pass.deadline);

  useEffect(() => {
    if (canvas.current) {
      QRCode.toCanvas(canvas.current, json, { width: 280, margin: 2, errorCorrectionLevel: "M" }).catch(() => {});
    }
  }, [json]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  return (
    <div className="overlay" role="dialog" aria-modal="true" aria-label={`Check-in pass for ${title}`} onClick={onClose}>
      <div className="pass" onClick={(e) => e.stopPropagation()}>
        <div className="pass-head">
          <div>
            <h3>{title}</h3>
            <p className="muted">Holder: {holder}</p>
          </div>
          <button type="button" className="ghost small" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>
        <canvas ref={canvas} className={left === 0 ? "qr expired" : "qr"} />
        <p className={left === 0 ? "countdown expired" : left < 60 ? "countdown soon" : "countdown"}>
          {left === 0 ? "Expired: make a new pass" : `Valid for ${mmss(left)}`}
        </p>
        <p className="muted small-print">
          Signed by the holder's wallet (EIP-712) for ticket #{pass.token_id} at this event's contract. Nothing was
          sent to the chain; the gate's check-in carries it.
        </p>
        <div className="row center">
          <Copyable value={json} />
          <button type="button" onClick={onHandToGate} disabled={left === 0}>
            Hand to the gate
          </button>
        </div>
      </div>
    </div>
  );
}
