import { useState, type ReactNode } from "react";
import { ApiError } from "./api";

/** A value too long for a cell, shortened, with a copy button for the whole. */
export function Copyable({ value, short = true }: { value: string; short?: boolean }) {
  const [copied, setCopied] = useState(false);
  const shown = short && value.length > 18 ? `${value.slice(0, 10)}…${value.slice(-6)}` : value;
  const copy = async () => {
    if (await writeClipboard(value)) {
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } else {
      window.prompt("Copy this value:", value);
    }
  };
  return (
    <span className="copyable">
      <code title={value}>{shown}</code>
      <button type="button" className="copy" onClick={copy} aria-label="Copy to clipboard">
        {copied ? "Copied" : "Copy"}
      </button>
    </span>
  );
}

export function Card({ title, step, intro, children }: { title: string; step: string; intro: ReactNode; children: ReactNode }) {
  return (
    <section className="card">
      <div className="card-head">
        <span className="step">{step}</span>
        <h2>{title}</h2>
      </div>
      <p className="intro">{intro}</p>
      {children}
    </section>
  );
}

export function ErrorLine({ error }: { error: unknown }) {
  if (!error) return null;
  if (error instanceof ApiError) {
    return (
      <p className="error" role="alert">
        <span className="code">
          {error.status || "—"} {error.code}
        </span>{" "}
        {error.title}
        {error.detail ? `: ${error.detail}` : ""}
      </p>
    );
  }
  return (
    <p className="error" role="alert">
      {error instanceof Error ? error.message : String(error)}
    </p>
  );
}

export function Empty({ children }: { children: ReactNode }) {
  return <p className="empty">{children}</p>;
}

/**
 * Copies text; false if the browser wouldn't. navigator.clipboard exists only
 * on HTTPS or localhost; over plain http (the cluster's LAN address) the older
 * execCommand("copy") on a selected, off-screen textarea still works.
 */
async function writeClipboard(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch {
    // fall through to the older way
  }
  const area = document.createElement("textarea");
  area.value = text;
  area.setAttribute("readonly", "");
  area.style.position = "fixed";
  area.style.left = "-9999px";
  document.body.appendChild(area);
  area.select();
  try {
    return document.execCommand("copy");
  } catch {
    return false;
  } finally {
    area.remove();
  }
}
