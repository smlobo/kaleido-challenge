import { useEffect, useState } from "react";

// What the page has created survives a reload, in this browser only. These
// are local demo tokens; "Clear" in the header forgets them all.
const PREFIX = "ticketing-demo:";

export function usePersisted<T>(key: string, initial: T): [T, (v: T | ((prev: T) => T)) => void] {
  const [value, setValue] = useState<T>(() => {
    try {
      const raw = localStorage.getItem(PREFIX + key);
      return raw ? (JSON.parse(raw) as T) : initial;
    } catch {
      return initial;
    }
  });
  useEffect(() => {
    try {
      localStorage.setItem(PREFIX + key, JSON.stringify(value));
    } catch {
      // private mode or storage full: keep working in memory
    }
  }, [key, value]);
  return [value, setValue];
}

export function clearPersisted() {
  try {
    Object.keys(localStorage)
      .filter((k) => k.startsWith(PREFIX))
      .forEach((k) => localStorage.removeItem(k));
  } catch {
    // nothing stored
  }
}
