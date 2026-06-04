import { useEffect, useRef, useState } from "react";
import { api } from "../api";
import type { ObservatorySnapshot } from "../types";

const POLL_MS = 2000;

export interface SnapshotState {
  snapshot: ObservatorySnapshot | null;
  error: string | null;
  stale: boolean;
}

/** Polls the observatory snapshot every 2s; holds last-good on error (stale=true). */
export function useSnapshot(): SnapshotState {
  const [snapshot, setSnapshot] = useState<ObservatorySnapshot | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [stale, setStale] = useState(false);
  const timer = useRef<number | null>(null);

  useEffect(() => {
    let active = true;
    async function poll() {
      try {
        const snap = await api.getObservatorySnapshot();
        if (!active) return;
        setSnapshot(snap);
        setError(null);
        setStale(false);
      } catch (e) {
        if (!active) return;
        setError(e instanceof Error ? e.message : "snapshot failed");
        setStale(true);
      }
    }
    void poll();
    timer.current = window.setInterval(() => void poll(), POLL_MS);
    return () => {
      active = false;
      if (timer.current !== null) window.clearInterval(timer.current);
    };
  }, []);

  return { snapshot, error, stale };
}
