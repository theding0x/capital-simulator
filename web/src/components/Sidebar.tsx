import { CHAPTERS } from "../chapters/registry";
import type { ChapterDef, CircuitNode, Volume } from "../chapters/registry";

interface SidebarProps {
  activeChapterId: string;
  selectedNode: CircuitNode | null;
  onSelect: (id: string) => void;
  onClearFilter: () => void;
}

const VOLUME_TITLE: Record<Volume, string> = {
  1: "Volume I — The Process of Production of Capital",
  2: "Volume II — The Process of Circulation of Capital",
  3: "Volume III — The Process of Capitalist Production as a Whole",
};

/** Human-readable label for a circuit node, used in the filter banner. */
const NODE_LABEL: Record<CircuitNode, string> = {
  "M":         "M (money capital)",
  "M-C":       "M—C (purchase)",
  "P":         "P (production)",
  "C-prime":   "C′ (finished commodity)",
  "C-M-prime": "C′—M′ (sale)",
  "M-prime":   "M′ (money returned)",
  "delta-M":   "ΔM (surplus-value)",
  "whole":     "the whole circuit",
  "historical":"historical preconditions",
};

type PartGroup = { name: string; chapters: ChapterDef[] };
type VolumeGroup = { volume: Volume; parts: PartGroup[] };

function groupByVolumeAndPart(chapters: ChapterDef[]): VolumeGroup[] {
  const volumes: VolumeGroup[] = [];
  for (const ch of chapters) {
    let vol = volumes.find((v) => v.volume === ch.volume);
    if (!vol) {
      vol = { volume: ch.volume, parts: [] };
      volumes.push(vol);
    }
    let part = vol.parts.find((p) => p.name === ch.part);
    if (!part) {
      part = { name: ch.part, chapters: [] };
      vol.parts.push(part);
    }
    part.chapters.push(ch);
  }
  return volumes;
}

export function Sidebar({
  activeChapterId,
  selectedNode,
  onSelect,
  onClearFilter,
}: SidebarProps) {
  const filtered = selectedNode
    ? CHAPTERS.filter((c) => c.circuitNode.includes(selectedNode))
    : CHAPTERS;
  const volumes = groupByVolumeAndPart(filtered);

  return (
    <nav className="sidebar" aria-label="Chapters">
      {selectedNode && (
        <div className="sidebar-filter">
          <span className="sidebar-filter-label">
            Filtered to <strong>{NODE_LABEL[selectedNode]}</strong>
          </span>
          <button
            type="button"
            className="sidebar-filter-clear"
            onClick={onClearFilter}
          >
            show all
          </button>
        </div>
      )}
      {volumes.length === 0 && (
        <div className="sidebar-empty">
          No chapters illuminate this moment of the circuit yet.
        </div>
      )}
      {volumes.map((vol) => (
        <div key={vol.volume} className="sidebar-volume">
          <div className="sidebar-volume-title">{VOLUME_TITLE[vol.volume]}</div>
          {vol.parts.map((part) => (
            <div key={part.name} className="sidebar-part-group">
              <div className="sidebar-part">{part.name}</div>
              {part.chapters.map((ch) => {
                const isActive = ch.id === activeChapterId;
                const cls = [
                  "sidebar-item",
                  isActive
                    ? "sidebar-item--active"
                    : ch.status === "done"
                    ? "sidebar-item--done"
                    : "sidebar-item--pending",
                ].join(" ");
                return (
                  <button
                    key={ch.id}
                    type="button"
                    className={cls}
                    onClick={() => onSelect(ch.id)}
                  >
                    <span className="sidebar-num">
                      {String(ch.number).padStart(2, "0")}
                    </span>
                    <span className="sidebar-title">{ch.title}</span>
                    {ch.status === "done" && !isActive && (
                      <span className="sidebar-icon">✓</span>
                    )}
                    {isActive && <span className="sidebar-icon">›</span>}
                  </button>
                );
              })}
            </div>
          ))}
        </div>
      ))}
    </nav>
  );
}
