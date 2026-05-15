import { useEffect, useState } from "react";
import { api } from "./api";
import type { Commodity, Owner } from "./types";
import { CHAPTERS } from "./chapters/registry";
import { Sidebar } from "./components/Sidebar";
import { ChapterShell } from "./components/ChapterShell";
import { CurrencyProvider, CurrencyToggle } from "./CurrencyContext";

const defaultChapterId =
  CHAPTERS.filter((c) => c.status === "done").at(-1)?.id ?? "ch01";

export default function App() {
  const [activeChapterId, setActiveChapterId] = useState(defaultChapterId);
  const [commodities, setCommodities] = useState<Commodity[]>([]);
  const [owners, setOwners] = useState<Owner[]>([]);

  async function refreshShared() {
    try {
      const [comms, ownerList] = await Promise.all([
        api.listCommodities(),
        api.listOwners(),
      ]);
      setCommodities(comms);
      setOwners(ownerList);
    } catch {
      // chapter panels handle their own error display
    }
  }

  useEffect(() => { refreshShared(); }, []);

  const activeChapter = CHAPTERS.find((c) => c.id === activeChapterId);

  return (
    <CurrencyProvider>
      <div className="app-shell">
        <header className="topbar">
          <span className="topbar-logo">Capital Simulator</span>
          {activeChapter && (
            <span className="topbar-chapter">
              Vol. I &middot; Ch.{String(activeChapter.number).padStart(2, "0")} &mdash;{" "}
              {activeChapter.title}
            </span>
          )}
          <CurrencyToggle />
        </header>
        <Sidebar activeChapterId={activeChapterId} onSelect={setActiveChapterId} />
        <ChapterShell
          activeChapterId={activeChapterId}
          commodities={commodities}
          owners={owners}
          onSharedChanged={refreshShared}
        />
      </div>
    </CurrencyProvider>
  );
}
