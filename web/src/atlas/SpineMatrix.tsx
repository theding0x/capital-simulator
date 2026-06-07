import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import "./spine.css";
import {
  CIRCUIT_NODES,
  cellEntries,
  cellKey,
  cellTiers,
  chapterId,
  volRoman,
  type InventoryEntry,
  type InventoryNode,
} from "./inventory";
import { ChapterDrawer } from "./ChapterDrawer";
import { Gloss } from "./Gloss";
import { GLOSS_MAP, type GlossContent } from "./glossContent";

interface SpineMatrixProps {
  /** Optionally pre-select a chapter on mount (deep-link). */
  initialChapterId?: string;
}

interface CellPos {
  volume: 1 | 2 | 3;
  node: InventoryNode;
}

const VOLS: { v: 1 | 2 | 3; name: string; depth: string }[] = [
  { v: 1, name: "I",   depth: "production"  },
  { v: 2, name: "II",  depth: "circulation" },
  { v: 3, name: "III", depth: "totality"    },
];

/** Gloss content for a T3 chapter: a canonical Marx quote for the curated
 *  fetishism/concept nodes, otherwise the chapter's own qualitative note from
 *  the inventory — flagged `note` so it renders as an editorial note, never
 *  dressed as a verbatim Marx quotation with a fabricated citation. */
function glossFor(entry: InventoryEntry): GlossContent & { note?: boolean } {
  const canon = GLOSS_MAP[chapterId(entry)];
  if (canon) return canon;
  return {
    label: "representation note",
    quote: entry.note,
    citation: "",
    note: true,
  };
}

export function SpineMatrix({ initialChapterId }: SpineMatrixProps) {
  const cellMap = useMemo(() => cellEntries(), []);

  // Default selection: Vol 1 | P
  const [selPos, setSelPos] = useState<CellPos>({ volume: 1, node: "P" });
  const [drawerEntry, setDrawerEntry] = useState<InventoryEntry | null>(null);

  // roving tabindex: track focused [row, col] indices within non-empty cells
  const gridRef = useRef<HTMLDivElement>(null);

  // If initialChapterId is provided, open the first matching entry on mount.
  useEffect(() => {
    if (!initialChapterId) return;
    for (const entries of cellMap.values()) {
      const match = entries.find((e) => chapterId(e) === initialChapterId);
      if (match) {
        setSelPos({ volume: match.volume, node: match.nodes[0] });
        setDrawerEntry(match);
        break;
      }
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const selKey = cellKey(selPos.volume, selPos.node);
  const selEntries = cellMap.get(selKey) ?? [];
  const selVol = VOLS.find((v) => v.v === selPos.volume)!;

  // Cell refs matrix for roving tabindex: [volIndex][colIndex]
  const cellRefs = useRef<(HTMLDivElement | null)[][]>(
    VOLS.map(() => CIRCUIT_NODES.map(() => null))
  );

  const setCellRef = useCallback(
    (vi: number, ci: number) => (el: HTMLDivElement | null) => {
      cellRefs.current[vi][ci] = el;
    },
    []
  );

  // Roving tabindex: only non-empty cells are tab-stops
  const isNonEmpty = useCallback(
    (vi: number, ci: number): boolean => {
      const v = VOLS[vi].v;
      const node = CIRCUIT_NODES[ci];
      return (cellMap.get(cellKey(v, node))?.length ?? 0) > 0;
    },
    [cellMap]
  );

  const focusCell = useCallback((vi: number, ci: number) => {
    const el = cellRefs.current[vi]?.[ci];
    if (el) el.focus();
  }, []);

  const handleCellKeyDown = useCallback(
    (e: React.KeyboardEvent, vi: number, ci: number) => {
      let nvi = vi, nci = ci;
      if (e.key === "ArrowRight") {
        // move right within row, no wrap
        let next = ci + 1;
        while (next < CIRCUIT_NODES.length) {
          if (isNonEmpty(vi, next)) { nci = next; break; }
          next++;
        }
        e.preventDefault();
      } else if (e.key === "ArrowLeft") {
        let next = ci - 1;
        while (next >= 0) {
          if (isNonEmpty(vi, next)) { nci = next; break; }
          next--;
        }
        e.preventDefault();
      } else if (e.key === "ArrowDown") {
        let next = vi + 1;
        while (next < VOLS.length) {
          if (isNonEmpty(next, ci)) { nvi = next; break; }
          next++;
        }
        e.preventDefault();
      } else if (e.key === "ArrowUp") {
        let next = vi - 1;
        while (next >= 0) {
          if (isNonEmpty(next, ci)) { nvi = next; break; }
          next--;
        }
        e.preventDefault();
      } else if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        const v = VOLS[vi].v;
        const node = CIRCUIT_NODES[ci];
        const entries = cellMap.get(cellKey(v, node)) ?? [];
        if (entries.length > 0) {
          setSelPos({ volume: v, node });
        }
        return;
      } else {
        return;
      }
      if (nvi !== vi || nci !== ci) focusCell(nvi, nci);
    },
    [isNonEmpty, focusCell, cellMap]
  );

  function handleCellClick(v: 1 | 2 | 3, node: InventoryNode) {
    const entries = cellMap.get(cellKey(v, node)) ?? [];
    if (entries.length === 0) return;
    setSelPos({ volume: v, node });
  }

  function handleEntryClick(entry: InventoryEntry) {
    setDrawerEntry(entry);
  }

  return (
    <div className="spine-section">
      <div className="spine-frame" ref={gridRef}>
        {/* Column axis */}
        <div className="spine-axis-x" role="row" aria-hidden="true">
          {CIRCUIT_NODES.map((n) => (
            <div
              key={n}
              className={"spine-ax" + (n === "whole" || n === "historical" ? " cross" : "")}
            >
              {n}
            </div>
          ))}
        </div>

        {/* Grid rows */}
        <div role="grid" aria-label="Circuit × Volume — the spine matrix">
          {VOLS.map((vol, vi) => (
            <div key={vol.v} className="spine-row" role="row">
              <div className="spine-rowlbl" role="rowheader">
                <span className="vn">Vol {vol.name}</span>
                <span className="vd">{vol.depth}</span>
              </div>
              {CIRCUIT_NODES.map((node, ci) => {
                const k = cellKey(vol.v, node);
                const entries = cellMap.get(k) ?? [];
                const empty = entries.length === 0;
                const tiers = cellTiers(entries);
                const isSel = selKey === k;
                const tabIdx = !empty && isSel ? 0 : -1;

                return (
                  <div
                    key={node}
                    ref={setCellRef(vi, ci)}
                    role="gridcell"
                    aria-selected={isSel}
                    aria-disabled={empty}
                    aria-label={
                      empty
                        ? `Vol ${vol.name} ${node} — empty`
                        : `Vol ${vol.name} ${node} — ${entries.length} chapter${entries.length > 1 ? "s" : ""}, tiers ${tiers.join(", ")}`
                    }
                    tabIndex={empty ? -1 : tabIdx}
                    className={
                      "spine-cell" +
                      (empty ? " empty" : "") +
                      (isSel ? " sel" : "")
                    }
                    onClick={() => handleCellClick(vol.v, node)}
                    onKeyDown={(e) => handleCellKeyDown(e, vi, ci)}
                  >
                    {!empty && (
                      <span className="spine-cell-n">{entries.length} ch</span>
                    )}
                    {!empty && (
                      <span className="spine-dots" aria-hidden="true">
                        {tiers.map((t) => (
                          <span
                            key={t}
                            className={`spine-dot dot-${t}`}
                            aria-label={t}
                          />
                        ))}
                      </span>
                    )}
                  </div>
                );
              })}
            </div>
          ))}
        </div>

        {/* Readout */}
        <div className="spine-readout" aria-live="polite" aria-atomic="true">
          <div className="ro-eyebrow">
            Vol {selVol.name} · {selVol.depth}&nbsp;&times;&nbsp;circuit moment &ldquo;{selPos.node}&rdquo;
          </div>
          {selEntries.length === 0 ? (
            <div className="ro-empty">No chapters address this circuit moment at this depth.</div>
          ) : (
            <>
              <div className="ro-title">
                {selEntries.length} chapter{selEntries.length > 1 ? "s" : ""} open here
              </div>
              <div className="ro-chs">
                {selEntries.slice(0, 6).map((entry) => {
                  const hasT3 = entry.tiers.includes("T3");
                  const g = hasT3 ? glossFor(entry) : null;
                  return (
                    <div className="ro-ch-wrap" key={`${entry.volume}.${entry.ch}`}>
                      <button
                        className="ro-ch"
                        style={{ all: "unset", cursor: "pointer", display: "flex", gap: 9, alignItems: "baseline", width: "100%" }}
                        onClick={() => handleEntryClick(entry)}
                        aria-label={`Open Vol ${volRoman(entry.volume)} Ch.${entry.ch} — ${entry.title}`}
                      >
                        <span className="cn">
                          {volRoman(entry.volume)}&middot;{entry.ch}
                        </span>
                        <span className="ct">{entry.title}</span>
                        <span className="ch-tiers">
                          {entry.tiers.map((t) => (
                            <span key={t} className={`tier-pill ${t}`} aria-label={t}>
                              {t}
                            </span>
                          ))}
                        </span>
                      </button>
                      {g && (
                        <Gloss
                          label={g.label}
                          quote={g.quote}
                          citation={g.citation}
                          form={g.form}
                          note={g.note}
                        />
                      )}
                    </div>
                  );
                })}
                {selEntries.length > 6 && (
                  <div className="ro-ch">
                    <span className="cn" />
                    <span className="ct" style={{ color: "var(--ink-dim)", fontStyle: "italic" }}>
                      +{selEntries.length - 6} more
                    </span>
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </div>

      {/* Movement hints */}
      <div className="spine-moves" aria-label="Navigation guide">
        <div className="spine-move">
          <span className="ico" aria-hidden="true">&#8596;</span>
          <div>
            <div className="mt">Pan — the circuit</div>
            <div className="md">
              Move along the moments of{" "}
              <span style={{ fontFamily: "var(--font-mono)" }}>
                M &#8594; M-C &#8594; P &#8594; C&#8242; &#8594; C-M&#8242; &#8594; M&#8242; &#8594; &#916;M
              </span>
              . Same depth, different phase.
            </div>
          </div>
        </div>
        <div className="spine-move">
          <span className="ico" aria-hidden="true">&#8595;</span>
          <div>
            <div className="mt">Descend — the depth</div>
            <div className="md">
              Sink through the volumes. The same circuit, explained deeper:
              production &#8594; circulation &#8594; the totality.
            </div>
          </div>
        </div>
        <div className="spine-move">
          <span className="ico" aria-hidden="true">&#9678;</span>
          <div>
            <div className="mt">Open — the cell</div>
            <div className="md">
              Any cell opens to its chapter(s), each rendered at its tier.
            </div>
          </div>
        </div>
      </div>

      {/* Chapter drawer overlay */}
      {drawerEntry && (
        <ChapterDrawer
          entry={drawerEntry}
          onClose={() => setDrawerEntry(null)}
        />
      )}
    </div>
  );
}
