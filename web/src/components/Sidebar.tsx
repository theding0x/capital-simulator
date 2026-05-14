import { CHAPTERS } from "../chapters/registry";
import type { ChapterDef } from "../chapters/registry";

interface SidebarProps {
  activeChapterId: string;
  onSelect: (id: string) => void;
}

export function Sidebar({ activeChapterId, onSelect }: SidebarProps) {
  const parts: { name: string; chapters: ChapterDef[] }[] = [];
  for (const ch of CHAPTERS) {
    let part = parts.find((p) => p.name === ch.part);
    if (!part) {
      part = { name: ch.part, chapters: [] };
      parts.push(part);
    }
    part.chapters.push(ch);
  }

  return (
    <nav className="sidebar">
      {parts.map((part) => (
        <div key={part.name}>
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
              <button key={ch.id} className={cls} onClick={() => onSelect(ch.id)}>
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
    </nav>
  );
}
