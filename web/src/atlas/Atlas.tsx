import { useState } from "react";
import "./atlas.css";
import { useSnapshot } from "./useSnapshot";
import { CircuitField } from "./CircuitField";
import { Abode } from "./Abode";
import { VitalSigns } from "./VitalSigns";
import { TickHeartbeat } from "./TickHeartbeat";

/** The Observatory: the whole circuit of capital, in motion — and beneath it,
 *  the hidden abode of production. */
export default function Atlas() {
  const { snapshot, stale } = useSnapshot();
  const [speed, setSpeed] = useState(1);
  const [descended, setDescended] = useState(false);

  return (
    <div className="atlas-shell">
      <header className="atlas-top">
        <span className="brand">Capital Simulator — Atlas</span>
        <nav className="nav">
          <a href="#/">Atlas</a>
          <a href="#/chapters">Chapters</a>
        </nav>
      </header>

      <div className="atlas-body">
        <aside className="atlas-rail">
          {snapshot ? (
            <VitalSigns vitals={snapshot.aggregate} capitalCount={snapshot.capitals.length} />
          ) : (
            <p style={{ color: "var(--ink-muted)", fontSize: 13 }}>Loading the field…</p>
          )}
          {snapshot && (
            <button
              className={"abode-threshold" + (descended ? " open" : "")}
              data-testid="threshold"
              onClick={() => setDescended((d) => !d)}
            >
              {descended ? "↑ Ascend to the surface" : "↓ Descend into production"}
            </button>
          )}
        </aside>

        {snapshot ? (
          descended ? (
            <div className="atlas-field-wrap abode-wrap">
              <Abode abode={snapshot.abode} />
            </div>
          ) : (
            <CircuitField snapshot={snapshot} speed={speed} />
          )
        ) : (
          <div className="atlas-field-wrap" />
        )}
      </div>

      <footer className="atlas-bottom">
        <TickHeartbeat
          tick={snapshot?.tick ?? 0}
          running={snapshot?.running ?? false}
          stale={stale}
          speed={speed}
          onSpeed={setSpeed}
        />
      </footer>
    </div>
  );
}
