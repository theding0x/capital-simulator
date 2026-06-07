import { useEffect, useState } from "react";
import { api } from "./api";
import type { Commodity, Owner } from "./types";
import { CHAPTERS } from "./chapters/registry";
import type { CircuitNode } from "./chapters/registry";
import { Sidebar } from "./components/Sidebar";
import { ChapterShell } from "./components/ChapterShell";
import { CircuitSpine } from "./components/CircuitSpine";
import { TurnoverPlayer } from "./components/TurnoverPlayer";
import { CurrencyProvider, CurrencyToggle } from "./CurrencyContext";

const defaultChapterId =
  CHAPTERS.filter((c) => c.status === "done").at(-1)?.id ?? "v1-ch01";

export default function App() {
  const [activeChapterId, setActiveChapterId] = useState(defaultChapterId);
  const [selectedNode, setSelectedNode] = useState<CircuitNode | null>(null);
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

  useEffect(() => {
    refreshShared();
  }, []);

  function handleSelectNode(node: CircuitNode) {
    setSelectedNode((prev) => (prev === node ? null : node));
  }

  return (
    <CurrencyProvider>
      <div className="app-shell">
        <header className="topbar">
          <span className="topbar-logo">Capital Simulator</span>
          <a href="#/" style={{ color: "var(--ink-muted)", textDecoration: "none", fontSize: 13, marginLeft: 12 }}>← Atlas</a>
          <span className="topbar-subtitle">— the self-expansion of value —</span>
          <CurrencyToggle />
        </header>
        <CircuitSpine selectedNode={selectedNode} onSelectNode={handleSelectNode}>
          <TurnoverPlayer setSelectedNode={setSelectedNode} />
        </CircuitSpine>
        <Sidebar
          activeChapterId={activeChapterId}
          selectedNode={selectedNode}
          onSelect={setActiveChapterId}
          onClearFilter={() => setSelectedNode(null)}
        />
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
