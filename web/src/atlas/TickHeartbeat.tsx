import { useEffect, useMemo, useRef, useState } from "react";
import { api } from "../api";
import type { EngineTick } from "../types";

interface TransportProps {
  tick: number;
  running: boolean;
  onToggle: () => void;
  speed: number;
  onSpeed: (s: number) => void;
  reduced: boolean;
  onReduced: (v: boolean) => void;
}

/** The console: play/pause, ×1/2/5/10 speed, an ECG of recent engine ticks, a
 *  reduced-motion toggle, and the turn counter. */
export function Transport({
  tick,
  running,
  onToggle,
  speed,
  onSpeed,
  reduced,
  onReduced,
}: TransportProps) {
  const [ticks, setTicks] = useState<EngineTick[]>([]);
  const timer = useRef<number | null>(null);
  useEffect(() => {
    let active = true;
    const poll = async () => {
      try {
        const t = await api.listEngineTicks(60);
        if (active) setTicks(t);
      } catch {
        /* keep last */
      }
    };
    void poll();
    timer.current = window.setInterval(() => void poll(), 2000);
    return () => {
      active = false;
      if (timer.current !== null) window.clearInterval(timer.current);
    };
  }, []);
  const pts = useMemo(() => {
    if (!ticks.length) return "0,14 240,14";
    const ordered = [...ticks].reverse(); // oldest → newest
    const max = Math.max(1, ...ordered.map((t) => t.entities_advanced));
    const n = ordered.length;
    return ordered
      .map(
        (t, i) =>
          `${((i / Math.max(1, n - 1)) * 240).toFixed(1)},${(26 - (t.entities_advanced / max) * 24).toFixed(1)}`
      )
      .join(" ");
  }, [ticks]);
  return (
    <div className="transport">
      <button className="tp-btn" onClick={onToggle} aria-label={running ? "Pause" : "Play"}>
        {running ? "❚❚" : "▶"}
      </button>
      <div className="tp-speeds" role="group" aria-label="Speed">
        {[1, 2, 5, 10].map((s) => (
          <button
            key={s}
            className={"tp-speed" + (speed === s ? " active" : "")}
            onClick={() => onSpeed(s)}
          >
            ×{s}
          </button>
        ))}
      </div>
      <svg className="tp-ecg" viewBox="0 0 240 28" preserveAspectRatio="none">
        <polyline points={pts} fill="none" stroke="#c8a240" strokeWidth="1.4" opacity="0.85" />
      </svg>
      <button
        className={"tp-reduced" + (reduced ? " on" : "")}
        onClick={() => onReduced(!reduced)}
        title="Reduce motion"
      >
        {reduced ? "motion off" : "motion on"}
      </button>
      <span className="tp-turn">
        turn <b>{tick}</b>
      </span>
    </div>
  );
}
