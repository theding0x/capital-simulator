// Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers
import { useEffect, useState } from "react";
import { turnoversApi } from "../../api";
import type { Turnover, TurnoverNumber } from "../../types";
import "./Ch07TurnoverTime.css";

/** Format basis points as a decimal turnover number (e.g. 521428 bp → "52.14/year"). */
function formatN(bp: number): string {
  return (bp / 10000).toFixed(2) + "/year";
}

/** Format minutes as a human-readable duration. */
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

  if (loading) return <p>Loading turnovers…</p>;
  if (rows.length === 0) return <p>No turnovers recorded.</p>;

  return (
    <div className="ch07-turnover">
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

      <p style={{ fontSize: "0.72rem", color: "var(--ink-muted)", margin: 0 }}>
        n&nbsp;=&nbsp;U&nbsp;/&nbsp;u, where U&nbsp;=&nbsp;525,600 min (1 year) and u&nbsp;= average turnover time.
        BasisPoints&nbsp;=&nbsp;10,000&nbsp;×&nbsp;n.
      </p>
    </div>
  );
}
