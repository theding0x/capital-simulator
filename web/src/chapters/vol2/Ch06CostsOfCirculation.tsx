// Vol. II Ch. 6 — The Costs of Circulation
// §I: Faux frais of production (purchase/sale time, book-keeping, money as faux frais)
// §II: Costs of storage (commodity supply, voluntary vs involuntary)
// §III: Costs of transport (value-adding; LabourCost + MeansOfTransport + Surplus)
import { useEffect, useState } from "react";
import { circulationCostsApi } from "../../api";
import type {
  CirculationCost,
  SystemFauxFraisResult,
  TransportTariff,
  TransportLeg,
} from "../../types";
import "./Ch06CostsOfCirculation.css";

function NaturePill({ nature }: { nature: string }) {
  return (
    <span className={`nature-pill ${nature}`}>
      {nature === "value_adding" ? "value-adding" : "faux frais"}
    </span>
  );
}

function penceStr(p: number) {
  return `${p.toLocaleString()}p`;
}

export function Ch06CostsOfCirculation() {
  const [costs, setCosts] = useState<CirculationCost[]>([]);
  const [fauxFrais, setFauxFrais] = useState<SystemFauxFraisResult | null>(null);
  const [tariffs, setTariffs] = useState<TransportTariff[]>([]);
  const [legs, setLegs] = useState<TransportLeg[]>([]);

  useEffect(() => {
    circulationCostsApi.list().then((r) => setCosts(r.items)).catch(() => {});
    circulationCostsApi.systemFauxFrais().then(setFauxFrais).catch(() => {});

    Promise.all([
      fetch("/v1/transport-tariffs").then((r) => r.json()).catch(() => ({ items: [] })),
      fetch("/v1/transport-legs").then((r) => r.json()).catch(() => ({ items: [] })),
    ]).then(([tariffRes, legRes]) => {
      if (Array.isArray(tariffRes.items)) setTariffs(tariffRes.items);
      if (Array.isArray(legRes.items)) setLegs(legRes.items);
    });
  }, []);

  const vpTotal = costs
    .filter((c) => c.nature === "value_preserving")
    .reduce((s, c) => s + c.pence, 0);
  const vaTotal = costs
    .filter((c) => c.nature === "value_adding")
    .reduce((s, c) => s + c.pence, 0);
  const total = vpTotal + vaTotal;

  return (
    <div className="ch06-panel">
      {/* §I — Faux frais (non-value-adding) */}
      <section className="ch06-section">
        <h2 className="ch06-h2">§I — Costs of Circulation (Faux Frais)</h2>
        <p className="small muted">
          Purchase/sale time, book-keeping, and money reserve are pure deductions from
          surplus-value — they preserve value but add none.
        </p>
        {costs.length === 0 ? (
          <p className="cost-empty">No circulation costs recorded.</p>
        ) : (
          <div className="cost-list">
            {costs.map((c) => (
              <div key={c.id} className="cost-row">
                <span className="cost-kind">{c.kind}</span>
                <span className="cost-ic">{c.industrial_capital_id || "—"}</span>
                <NaturePill nature={c.nature} />
                <span className="cost-pence">{penceStr(c.pence)}</span>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Aggregate */}
      {costs.length > 0 && (
        <section className="ch06-section">
          <h2 className="ch06-h2">Aggregate</h2>
          <div className="aggregate-block">
            <div className="agg-card">
              <span className="agg-label">Faux frais (value-preserving)</span>
              <span className="agg-value">{penceStr(vpTotal)}</span>
            </div>
            <div className="agg-card">
              <span className="agg-label">Transport (value-adding)</span>
              <span className="agg-value">{penceStr(vaTotal)}</span>
            </div>
            <div className="agg-card">
              <span className="agg-label">Total circulation costs</span>
              <span className="agg-value">{penceStr(total)}</span>
            </div>
          </div>
        </section>
      )}

      {/* §I(c) — Money as faux frais */}
      {fauxFrais && fauxFrais.items.length > 0 && (
        <section className="ch06-section">
          <h2 className="ch06-h2">§I(c) — Money Reserve (System Faux Frais)</h2>
          <p className="small muted">
            Money hoarded as a reserve never circulates as capital — it is a pure societal
            cost, like wear on coin.
          </p>
          <div className="faux-frais-list">
            {fauxFrais.items.map((item) => (
              <div key={item.id} className="faux-frais-row">
                <span className="ff-ic">{item.industrial_capital_id || "national reserve"}</span>
                <span className="ff-pence">{penceStr(item.pence)}</span>
                <span className="ff-annual">↻ {penceStr(item.annual_replacement_pence)}/yr</span>
              </div>
            ))}
          </div>
          <div className="aggregate-block">
            <div className="agg-card">
              <span className="agg-label">Total reserved</span>
              <span className="agg-value">{penceStr(fauxFrais.total_reserve_pence)}</span>
            </div>
            <div className="agg-card">
              <span className="agg-label">Annual replacement</span>
              <span className="agg-value">
                {penceStr(fauxFrais.total_annual_replacement_pence)}
              </span>
            </div>
          </div>
        </section>
      )}

      {/* §III — Transport */}
      <section className="ch06-section">
        <h2 className="ch06-h2">§III — Transport (Value-Adding)</h2>
        <p className="small muted">
          Transport changes the use-value's location — a real productive act. Labour + means
          of transport + surplus = value added to the commodity.
        </p>

        {tariffs.length > 0 && (
          <>
            <h3 className="small" style={{ margin: 0, fontWeight: 600 }}>Tariffs</h3>
            <div className="transport-grid">
              {tariffs.map((t) => {
                const multiplier =
                  Math.max(t.fragility_multiplier_basis_points, t.breakage_risk_multiplier_basis_points);
                const effectiveRate =
                  multiplier === 0
                    ? t.base_pence_per_ton_mile
                    : Math.round((t.base_pence_per_ton_mile * multiplier) / 10000);
                return (
                  <div key={t.id} className="tariff-card">
                    <span className="tariff-commodity">{t.commodity_id}</span>
                    <span className="tariff-rate">
                      {effectiveRate}p/ton-mile
                    </span>
                    {multiplier > 0 && (
                      <span className="small muted">
                        base {t.base_pence_per_ton_mile}p × {multiplier / 10000}×
                      </span>
                    )}
                  </div>
                );
              })}
            </div>
          </>
        )}

        {legs.length > 0 && (
          <>
            <h3 className="small" style={{ margin: 0, fontWeight: 600 }}>Transport Legs</h3>
            <div className="transport-grid">
              {legs.map((leg) => {
                const addedValue =
                  leg.labour_cost_pence + leg.means_of_transport_pence + leg.surplus_pence;
                return (
                  <div key={leg.id} className="leg-card">
                    <span className="leg-route">
                      {leg.origin_id} → {leg.destination_id}
                    </span>
                    <span className="leg-muted">{leg.commodity_id}</span>
                    <span className="leg-muted">
                      {(leg.distance_meters / 1609).toFixed(1)} miles ·{" "}
                      {(leg.weight_grams / 1_000_000).toFixed(2)} t
                    </span>
                    <span className="leg-added-value">
                      Added value: {penceStr(addedValue)}
                    </span>
                  </div>
                );
              })}
            </div>
          </>
        )}

        {tariffs.length === 0 && legs.length === 0 && (
          <p className="cost-empty">No transport data recorded.</p>
        )}
      </section>
    </div>
  );
}
