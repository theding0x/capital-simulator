import { useState } from "react";
import { api } from "../api";
import type { WageForm } from "../types";
import "./Ch19WageForm.css";

export function Ch19WageForm() {
  const [agentID, setAgentID] = useState("");
  const [dailyPence, setDailyPence] = useState(36);
  const [workingDayHours, setWorkingDayHours] = useState(12);
  const [lpvDailyPence, setLpvDailyPence] = useState(36);
  const [necessaryMinutes, setNecessaryMinutes] = useState(360);
  const [result, setResult] = useState<WageForm | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const wf = await api.createWageForm({
        agent_id: agentID,
        daily_pence: dailyPence,
        working_day_minutes: workingDayHours * 60,
        lpv_daily_pence: lpvDailyPence,
        necessary_minutes: necessaryMinutes,
      });
      setResult(wf);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    }
  }

  return (
    <div className="ch19-wage-form">
      <section className="ch19-section">
        <h2>Enter Wage Form</h2>
        <p className="ch19-explainer">
          The wage form conceals the split between necessary and surplus labour.
          Enter the apparent wage and the real value of labour-power to reveal
          what the wage hides.
        </p>
        <form onSubmit={handleSubmit} className="ch19-form">
          <label>
            Agent ID
            <input
              type="text"
              value={agentID}
              onChange={(e) => setAgentID(e.target.value)}
              required
            />
          </label>
          <label>
            Daily wage (pence)
            <input
              type="number"
              min={1}
              value={dailyPence}
              onChange={(e) => setDailyPence(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Working day (hours)
            <input
              type="number"
              min={1}
              max={24}
              value={workingDayHours}
              onChange={(e) => setWorkingDayHours(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Value of labour-power (pence/day)
            <input
              type="number"
              min={1}
              value={lpvDailyPence}
              onChange={(e) => setLpvDailyPence(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Necessary labour (minutes)
            <input
              type="number"
              min={1}
              value={necessaryMinutes}
              onChange={(e) => setNecessaryMinutes(Number(e.target.value))}
              required
            />
          </label>
          <button type="submit">Reveal Wage Form</button>
        </form>
        {error && <p className="ch19-error">{error}</p>}
      </section>

      {result && (
        <section className="ch19-section ch19-results">
          <h2>Results</h2>

          <div className="ch19-cards">
            <div className="ch19-card">
              <h3>Appearance</h3>
              <p className="ch19-card-subtitle">How the wage presents itself</p>
              <dl>
                <dt>Hourly wage</dt>
                <dd>{result.hourly_wage}p/hr</dd>
                <dt>Paid hours</dt>
                <dd>{result.appearance.paid_hours} h</dd>
                <dt>Unpaid hours visible</dt>
                <dd>{result.appearance.unpaid_hours} h</dd>
              </dl>
              <p className="ch19-note">
                The wage form makes all {result.appearance.paid_hours} hours
                appear paid — surplus labour is invisible.
              </p>
            </div>

            <div className="ch19-card ch19-card-reveal">
              <h3>Decomposition</h3>
              <p className="ch19-card-subtitle">The real paid/unpaid split</p>
              <dl>
                <dt>Necessary labour (paid)</dt>
                <dd>{result.decomposition.paid_minutes} min</dd>
                <dt>Surplus labour (unpaid)</dt>
                <dd>{result.decomposition.unpaid_minutes} min</dd>
              </dl>
              <div className="ch19-bar">
                <div
                  className="ch19-bar-paid"
                  style={{
                    width: `${
                      (result.decomposition.paid_minutes /
                        (result.decomposition.paid_minutes +
                          result.decomposition.unpaid_minutes)) *
                      100
                    }%`,
                  }}
                  title="Paid (necessary) labour"
                />
                <div
                  className="ch19-bar-unpaid"
                  style={{
                    width: `${
                      (result.decomposition.unpaid_minutes /
                        (result.decomposition.paid_minutes +
                          result.decomposition.unpaid_minutes)) *
                      100
                    }%`,
                  }}
                  title="Unpaid (surplus) labour"
                />
              </div>
              <p className="ch19-note">
                Only {result.decomposition.paid_minutes} min reproduces
                labour-power; {result.decomposition.unpaid_minutes} min is
                surplus extracted by capital.
              </p>
            </div>
          </div>
        </section>
      )}
    </div>
  );
}
