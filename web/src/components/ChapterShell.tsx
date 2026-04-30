import { CHAPTERS } from "../chapters/registry";
import type { Commodity, Owner } from "../types";
import { Ch01Commodity } from "../chapters/Ch01Commodity";
import { Ch02Exchange } from "../chapters/Ch02Exchange";
import { Ch03Money } from "../chapters/Ch03Money";

interface ChapterShellProps {
  activeChapterId: string;
  commodities: Commodity[];
  owners: Owner[];
  onSharedChanged: () => void;
}

const QUOTES: Partial<Record<string, string>> = {
  ch01: "The wealth of societies in which the capitalist mode of production prevails appears as an immense collection of commodities.",
  ch02: "Commodities cannot go to market and make exchanges of their own account.",
  ch03: "Money is a crystal formed of necessity in the course of exchanges.",
};

export function ChapterShell({
  activeChapterId,
  commodities,
  owners,
  onSharedChanged,
}: ChapterShellProps) {
  const chapter = CHAPTERS.find((c) => c.id === activeChapterId);
  if (!chapter) return null;

  const quote = QUOTES[activeChapterId];

  return (
    <main className="chapter-main">
      <header className="chapter-header">
        <div className="chapter-header-num">
          Vol. I &middot; Chapter {String(chapter.number).padStart(2, "0")}
        </div>
        <h1 className="chapter-header-title">{chapter.title}</h1>
        {quote && <p className="chapter-header-quote">{quote}</p>}
      </header>

      <div className="chapter-body">
        {chapter.status === "pending" ? (
          <div className="chapter-placeholder">
            <p>Not yet implemented</p>
            <p className="small muted">Coming in a future chapter branch</p>
          </div>
        ) : activeChapterId === "ch01" ? (
          <Ch01Commodity commodities={commodities} onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch02" ? (
          <Ch02Exchange
            commodities={commodities}
            owners={owners}
            onSharedChanged={onSharedChanged}
          />
        ) : activeChapterId === "ch03" ? (
          <Ch03Money owners={owners} onSharedChanged={onSharedChanged} />
        ) : null}
      </div>
    </main>
  );
}
