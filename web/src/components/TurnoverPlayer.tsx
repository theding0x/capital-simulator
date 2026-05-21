import { useEffect, useState } from "react";
import type { Dispatch, SetStateAction } from "react";
import type { CircuitNode } from "../chapters/registry";

/**
 * The animation cycle: M → M-C → P → C′ → C′—M′ → M′ → (loop)
 * These six nodes mirror the visual nodes in CircuitSpine; the player
 * advances `selectedNode` through them on a timer when playing.
 */
const ANIMATION_SEQUENCE: CircuitNode[] = [
  "M",
  "M-C",
  "P",
  "C-prime",
  "C-M-prime",
  "M-prime",
];

/** Base interval at ×1 speed. ×10 makes each step 200ms. */
const BASE_INTERVAL_MS = 2000;

type Speed = 1 | 2 | 5 | 10;

interface TurnoverPlayerProps {
  setSelectedNode: Dispatch<SetStateAction<CircuitNode | null>>;
}

export function TurnoverPlayer({ setSelectedNode }: TurnoverPlayerProps) {
  const [playing, setPlaying] = useState(false);
  const [speed, setSpeed] = useState<Speed>(1);

  useEffect(() => {
    if (!playing) return;
    const interval = BASE_INTERVAL_MS / speed;
    const tid = setInterval(() => {
      setSelectedNode((prev) => {
        const currentIndex = prev ? ANIMATION_SEQUENCE.indexOf(prev) : -1;
        // If prev is null OR a non-motion meta-node, start at M.
        const safeIndex = currentIndex < 0 ? -1 : currentIndex;
        const next = ANIMATION_SEQUENCE[(safeIndex + 1) % ANIMATION_SEQUENCE.length];
        return next;
      });
    }, interval);
    return () => clearInterval(tid);
  }, [playing, speed, setSelectedNode]);

  function handlePlayPause() {
    setPlaying((p) => !p);
  }

  function handleRewind() {
    setPlaying(false);
    setSelectedNode(null);
  }

  return (
    <div className="turnover-player" aria-label="Turnover controls">
      <button
        type="button"
        className="turnover-btn"
        onClick={handlePlayPause}
        aria-label={playing ? "Pause turnover" : "Play turnover"}
        title={playing ? "Pause" : "Play"}
      >
        {playing ? "⏸" : "▶"}
      </button>
      <button
        type="button"
        className="turnover-btn"
        onClick={handleRewind}
        aria-label="Reset circuit to beginning"
        title="Reset"
      >
        ⏮
      </button>
      <span className="turnover-speed-group" role="group" aria-label="Animation speed">
        {([1, 2, 5, 10] as const).map((s) => (
          <button
            key={s}
            type="button"
            className={
              "turnover-speed" + (speed === s ? " turnover-speed--active" : "")
            }
            onClick={() => setSpeed(s)}
            aria-pressed={speed === s}
          >
            ×{s}
          </button>
        ))}
      </span>
    </div>
  );
}
