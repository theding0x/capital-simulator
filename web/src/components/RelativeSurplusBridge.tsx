import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { fmtMinutes } from "../format";
import type { ProductivitySource, RelativeSurplusResult } from "../types";

interface Props {
  source: ProductivitySource;
  sourceId: string;
  factor: number;
  sourceLabel: string;
}

// RelativeSurplusBridge surfaces Marx's Part IV arc directly in the UI:
// a Ch.13 / Ch.14 / Ch.15 record produces a productivity factor; we feed
// it into Ch.12's relative-surplus-value formula and show the new
// working-day split. The endpoint round-trips to simulation-engine, which
// then re-fetches the factor from agent-service (or its own store) — so
// this is the genuine end-to-end path, not a client-side simulation.
export function RelativeSurplusBridge({ source, sourceId, factor, sourceLabel }: Props) {
  const [workingDayTotal, setWorkingDayTotal] = useState(720);
  const [currentLpv, setCurrentLpv] = useState(600);
  const [result, setResult] = useState<RelativeSurplusResult | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setBusy(true);
    try {
      const r = await api.relativeSurplusFromSource({
        working_day_total: workingDayTotal,
        current_lpv: currentLpv,
        source,
        source_id: sourceId,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div style={{ marginTop: "0.75rem", padding: "0.5rem 0.75rem", border: "1px solid var(--ink-muted)" }}>
      <p className="small" style={{ marginTop: 0 }}>
        <strong>Feed productivity into Ch. 12.</strong> {sourceLabel} produces a factor of{" "}
        <code>{factor.toFixed(3)}</code>. Apply it to a working day to see the
        relative-SV gain Marx claims is possible (§1).
      </p>
      <form onSubmit={submit} style={{ display: "flex", gap: "0.5rem", flexWrap: "wrap", alignItems: "end" }}>
        <label className="small">
          <span style={{ display: "block" }}>Working day (min)</span>
          <input
            type="number"
            min={1}
            value={workingDayTotal}
            onChange={(e) => setWorkingDayTotal(Number(e.target.value))}
            style={{ width: "6rem" }}
          />
        </label>
        <label className="small">
          <span style={{ display: "block" }}>Current LPV (min)</span>
          <input
            type="number"
            min={1}
            value={currentLpv}
            onChange={(e) => setCurrentLpv(Number(e.target.value))}
            style={{ width: "6rem" }}
          />
        </label>
        <button type="submit" className="primary" disabled={busy}>
          {busy ? "…" : "Apply"}
        </button>
        {err && <span className="error">{err}</span>}
      </form>
      {result && (
        <div className="small" style={{ marginTop: "0.5rem" }}>
          <div>
            New LPV: <strong>{fmtMinutes(result.new_labour_power_value)}</strong>
            {" "}(was {fmtMinutes(currentLpv)})
          </div>
          <div>
            Surplus labour: {fmtMinutes(result.old_working_day.surplus_labour)} →{" "}
            <strong>{fmtMinutes(result.new_working_day.surplus_labour)}</strong>
            {" "}(Δ {result.relative_surplus_value > 0 ? "+" : ""}{result.relative_surplus_value} min)
          </div>
        </div>
      )}
    </div>
  );
}
