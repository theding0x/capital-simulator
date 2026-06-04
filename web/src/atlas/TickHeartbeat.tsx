import { useEffect, useRef, useState } from "react";
import { api } from "../api";
import type { EngineTick } from "../types";

interface TickHeartbeatProps {
  tick: number;
  running: boolean;
  stale: boolean;
  speed: number;
  onSpeed: (s: number) => void;
}

const SPEEDS = [1, 2, 5, 10];

/** Transport (play/pause/speed), tick counter, and an ECG of recent engine ticks. */
export function TickHeartbeat({ tick, running, stale, speed, onSpeed }: TickHeartbeatProps) {
  const [ticks, setTicks] = useState<EngineTick[]>([]);
  const [localRunning, setLocalRunning] = useState(running);
  const timer = useRef<number | null>(null);

  useEffect(() => setLocalRunning(running), [running]);

  useEffect(() => {
    let active = true;
    async function poll() {
      try {
        const t = await api.listEngineTicks(60);
        if (active) setTicks(t);
      } catch {
        /* heartbeat keeps last-good */
      }
    }
    void poll();
    timer.current = window.setInterval(() => void poll(), 2000);
    return () => {
      active = false;
      if (timer.current !== null) window.clearInterval(timer.current);
    };
  }, []);

  async function toggle() {
    const next = !localRunning;
    setLocalRunning(next); // optimistic
    try {
      if (next) await api.startEngine();
      else await api.stopEngine();
    } catch {
      setLocalRunning(!next); // revert on failure
    }
  }

  const points = ecgPoints(ticks);

  return (
    <>
      <button className="atlas-btn" onClick={() => void toggle()} aria-label={localRunning ? "Pause" : "Play"}>
        {localRunning ? "⏸" : "▶"}
      </button>
      <span role="group" aria-label="Animation speed">
        {SPEEDS.map((s) => (
          <span key={s} className={"atlas-speed" + (speed === s ? " active" : "")}
            onClick={() => onSpeed(s)}>×{s}</span>
        ))}
      </span>
      <svg className="atlas-ecg" viewBox="0 0 240 28" preserveAspectRatio="none" data-testid="atlas-ecg">
        <polyline points={points} fill="none" stroke="#c8a240" strokeWidth="1.4" opacity="0.85" />
      </svg>
      <span className="atlas-turn">turn {tick}</span>
      {stale && <span className="atlas-stale">⚠ reconnecting</span>}
    </>
  );
}

/** Map recent ticks' entities_advanced to an ECG polyline over a 240×28 box. */
function ecgPoints(ticks: EngineTick[]): string {
  if (ticks.length === 0) return "0,14 240,14";
  const ordered = [...ticks].reverse(); // oldest → newest
  const max = Math.max(1, ...ordered.map((t) => t.entities_advanced));
  const n = ordered.length;
  return ordered
    .map((t, i) => {
      const x = (i / Math.max(1, n - 1)) * 240;
      const y = 26 - (t.entities_advanced / max) * 24;
      return `${x.toFixed(1)},${y.toFixed(1)}`;
    })
    .join(" ");
}
