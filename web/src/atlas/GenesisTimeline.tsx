import { useCallback, useEffect, useRef, useState } from "react";
import { GENESIS_EPISODES, type GenesisEpisode } from "./genesis";
import { findEntry, type InventoryEntry } from "./inventory";
import { ChapterDrawer } from "./ChapterDrawer";

/** Vol I · Ch. 26–33 (+ Vol III historical) — the genesis timeline floor. The
 *  walkable historical bedrock beneath the living world: primitive accumulation
 *  as a red thread of episodes. Hovering or focusing an episode reveals its
 *  gloss; selecting it opens the chapter panel in the shared ChapterDrawer. */
export function GenesisTimeline() {
  const [activeId, setActiveId] = useState<string>(GENESIS_EPISODES[0].id);
  const [drawerEntry, setDrawerEntry] = useState<InventoryEntry | null>(null);

  // Horizontal-scroll affordance: fade the edges that have more beyond them.
  const trackRef = useRef<HTMLDivElement>(null);
  const [edges, setEdges] = useState({ atStart: true, atEnd: true });
  const syncEdges = useCallback(() => {
    const el = trackRef.current;
    if (!el) return;
    const atStart = el.scrollLeft <= 2;
    const atEnd = el.scrollLeft + el.clientWidth >= el.scrollWidth - 2;
    // Only re-render when an edge actually flips — onScroll fires continuously.
    setEdges((prev) =>
      prev.atStart === atStart && prev.atEnd === atEnd ? prev : { atStart, atEnd }
    );
  }, []);
  useEffect(() => {
    syncEdges();
    window.addEventListener("resize", syncEdges);
    return () => window.removeEventListener("resize", syncEdges);
  }, [syncEdges]);

  const active =
    GENESIS_EPISODES.find((e) => e.id === activeId) ?? GENESIS_EPISODES[0];

  function open(ep: GenesisEpisode) {
    const entry = findEntry(ep.volume, ep.ch);
    if (entry) setDrawerEntry(entry);
  }

  return (
    <div className="genesis">
      <div
        className={
          "genesis-scroller" +
          (edges.atStart ? " at-start" : "") +
          (edges.atEnd ? " at-end" : "")
        }
      >
      <div
        className="genesis-track"
        ref={trackRef}
        onScroll={syncEdges}
        role="group"
        aria-label="Primitive accumulation — the historical timeline"
      >
        <div className="genesis-line" aria-hidden="true" />
        {GENESIS_EPISODES.map((ep) => {
          const on = ep.id === active.id;
          return (
            <button
              key={ep.id}
              className={"genesis-node " + ep.kind + (on ? " on" : "")}
              onMouseEnter={() => setActiveId(ep.id)}
              onFocus={() => setActiveId(ep.id)}
              onClick={() => open(ep)}
              aria-label={`${ep.era} — ${ep.title}. ${ep.citation}. Open the chapter.`}
            >
              <span
                className="gn-dot"
                aria-hidden="true"
                style={{ ["--mag" as string]: ep.magnitude }}
              />
              <span className="gn-era">{ep.era}</span>
              <span className="gn-title">{ep.title}</span>
            </button>
          );
        })}
      </div>
        <span className="genesis-scrollhint" aria-hidden="true">
          scroll &rarr;
        </span>
      </div>

      <div className="genesis-readout" aria-live="polite" aria-atomic="true">
        <span className="gr-cite">{active.citation}</span>
        <p className="gr-blurb">{active.blurb}</p>
      </div>

      <p className="genesis-foot">
        Walkable beneath the living world &mdash; the rosy dawn was written
        &ldquo;in letters of blood and fire.&rdquo;
      </p>

      {drawerEntry && (
        <ChapterDrawer entry={drawerEntry} onClose={() => setDrawerEntry(null)} />
      )}
    </div>
  );
}
