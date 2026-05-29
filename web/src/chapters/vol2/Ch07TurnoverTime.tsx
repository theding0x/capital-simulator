// Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers
import { useEffect, useMemo, useState } from "react";
import { turnoversApi } from "../../api";
import type { Turnover, TurnoverNumber } from "../../types";
import "./Ch07TurnoverTime.css";

function formatN(bp: number): string {
  return (bp / 10000).toFixed(2) + "/year";
}

function formatMinutes(m: number): string {
  if (m < 60) return `${m} min`;
  if (m < 1440) return `${(m / 60).toFixed(1)} h`;
  if (m < 10080) return `${(m / 1440).toFixed(1)} days`;
  if (m < 525600) return `${(m / 10080).toFixed(1)} weeks`;
  return `${(m / 525600).toFixed(2)} years`;
}

interface TurnoverRow {
  turnover: Turnover;
  number: TurnoverNumber | null;
}

export function Ch07TurnoverTime() {
  const [rows, setRows] = useState<TurnoverRow[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    turnoversApi
      .list()
      .then(async (res) => {
        const items = res.items ?? [];
        const withNumbers = await Promise.all(
          items.map(async (t) => {
            try {
              const number = await turnoversApi.getNumber(t.id);
              return { turnover: t, number };
            } catch {
              return { turnover: t, number: null };
            }
          }),
        );
        setRows(withNumbers);
      })
      .finally(() => setLoading(false));
  }, []);

  const chartRows = useMemo(
    () =>
      rows
        .filter((r): r is TurnoverRow & { number: TurnoverNumber } => r.number !== null)
        .map((r) => ({
          id: r.turnover.id,
          lens: r.turnover.lens,
          n: r.number.basis_points / 10000,
          minutes: r.turnover.turnover_time_minutes,
        }))
        .sort((a, b) => b.n - a.n),
    [rows],
  );

  const maxN = chartRows.length > 0 ? chartRows[0].n : 1;
  const slowThreshold = 4;

  return (
    <>
      <NFormulaInsight />
      {loading && <p>Loading turnovers…</p>}
      {!loading && rows.length === 0 && <p>No turnovers recorded.</p>}
      {!loading && rows.length > 0 && (
        <div className="ch07-turnover">
          {chartRows.length > 0 && (
            <section className="v2-ch07-n-chart">
              <h3 className="v2-ch07-n-chart-h3">Turnovers per year (n) across fixtures</h3>
              {chartRows.map((c) => {
                const widthPct = (c.n / maxN) * 100;
                const slow = c.n < slowThreshold;
                return (
                  <div key={c.id} className="v2-ch07-n-row">
                    <span className="v2-ch07-n-label">
                      {c.lens === "productive" ? "Productive" : "Money"}
                      <span className="v2-ch07-n-label-sub">
                        u = {formatMinutes(c.minutes)}
                      </span>
                    </span>
                    <div className="v2-ch07-n-track">
                      <div
                        className={
                          "v2-ch07-n-fill" + (slow ? " v2-ch07-n-fill--slow" : "")
                        }
                        style={{ width: `${widthPct}%` }}
                      />
                    </div>
                    <span className="v2-ch07-n-value">{c.n.toFixed(2)}</span>
                  </div>
                );
              })}
            </section>
          )}

          <div className="ch07-fixtures">
            {rows.map(({ turnover, number }) => (
              <div key={turnover.id} className="ch07-fixture-card">
                <h4>{turnover.lens === "productive" ? "Productive Circuit" : "Money Circuit"}</h4>

                <div className="ch07-stat">
                  <span className="ch07-stat-label">Turnover time</span>
                  <span className="ch07-stat-value">{formatMinutes(turnover.turnover_time_minutes)}</span>
                </div>

                <div className="ch07-stat">
                  <span className="ch07-stat-label">Cycles recorded</span>
                  <span className="ch07-stat-value">{turnover.cycles.length}</span>
                </div>

                {number !== null && (
                  <>
                    <div className="ch07-stat">
                      <span className="ch07-stat-label">Turnovers/year (n)</span>
                      <span className="ch07-stat-value ch07-number-highlight">{formatN(number.basis_points)}</span>
                    </div>
                    <div className="ch07-stat">
                      <span className="ch07-stat-label">Basis points</span>
                      <span className="ch07-stat-value">{number.basis_points.toLocaleString()}</span>
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
      <HeterogeneityCoda />
    </>
  );
}

function NFormulaInsight() {
  return (
    <section className="v2-ch07-insight">
      <h2 className="v2-ch07-insight-h2">n = U / u</h2>
      <span className="v2-ch07-insight-formula">
        n = U / u · U = 525,600 min (1 year) · BP = 10,000 × n
      </span>
      <p className="v2-ch07-insight-prose">
        The <em>number of turnovers</em> per year is one over the turnover
        time. The chapter's central observation is that <code>n</code> is
        wildly heterogeneous across capitals — a cotton spinner can renew its
        circuit dozens of times a year while a wine-maturer might take more
        than four years for one cycle. The bar chart below makes that spread
        visible.
      </p>
    </section>
  );
}

function HeterogeneityCoda() {
  return (
    <aside className="v2-ch07-coda">
      <p className="v2-ch07-coda-quote">
        “The different periods of turnover of the various individual capitals
        invested in different branches of industry … fix the relations between
        the magnitude of the capital advanced and the mass of surplus-value
        annually produced.”
        <span className="v2-ch07-coda-cite">
          — Marx, Capital Vol. II, Ch. 7
        </span>
      </p>
    </aside>
  );
}
