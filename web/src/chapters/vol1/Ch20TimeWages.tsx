import { useState } from "react";
import { api } from "../../api";
import type { WorkingSession } from "../../types";
import { useCurrency, usePounds } from "../../CurrencyContext";
import "./Ch20TimeWages.css";

function TimeWageInsight({
  dailyPence,
  workingDayHours,
  onHoursChange,
}: {
  dailyPence: number;
  workingDayHours: number;
  onHoursChange: (h: number) => void;
}) {
  const probePrice =
    workingDayHours > 0 ? dailyPence / workingDayHours : 0;
  return (
    <section className="v1-ch20-insight">
      <h2 className="v1-ch20-insight-h2">
        Why the time-wage falls as the day lengthens
      </h2>
      <p className="v1-ch20-insight-prose">
        The daily price of labour-power is set first; the hourly price is a
        derivative — the day's pay divided by hours worked. Hold the daily wage
        constant and stretch the working day from 8h to 12h to 16h: the hourly
        price drops by a third, then by half. The capitalist appears to pay the
        same wage; the worker's labour is in fact purchased at a progressively
        lower unit price.
      </p>
      <div className="v1-ch20-probe">
        <label className="v1-ch20-probe-label" htmlFor="v1-ch20-probe-slider">
          Extension probe — drag to lengthen the working day
        </label>
        <input
          id="v1-ch20-probe-slider"
          className="v1-ch20-probe-slider"
          type="range"
          min={1}
          max={24}
          step={1}
          value={workingDayHours}
          onChange={(e) => onHoursChange(Number(e.target.value))}
          aria-label="Working day length in hours"
        />
        <div className="v1-ch20-probe-readout">
          <span className="v1-ch20-probe-hours">{workingDayHours}h</span>
          <span className="v1-ch20-probe-arrow">→</span>
          <span className="v1-ch20-probe-price">
            {probePrice.toFixed(2)}d. / hour
          </span>
        </div>
      </div>
    </section>
  );
}

function PriceOscillationCoda() {
  return (
    <aside className="v1-ch20-coda">
      <p className="v1-ch20-coda-quote">
        “The day-wages, week-wages, etc., may remain the same, although the
        price of labour falls. ... The actual price of labour falls below its
        nominal price.”
        <span className="v1-ch20-coda-cite">
          — Marx, Capital Vol. I, Ch. 20
        </span>
      </p>
    </aside>
  );
}

export function Ch20TimeWages() {
  const { modern } = useCurrency();
  const fmt = usePounds();
  const [agentID, setAgentID] = useState("");
  const [dailyPence, setDailyPence] = useState(36);
  const [workingDayHours, setWorkingDayHours] = useState(12);
  const [overtimeHours, setOvertimeHours] = useState(0);
  const [overtimeRate, setOvertimeRate] = useState(0);
  const [wagePeriod, setWagePeriod] = useState<"daily" | "weekly">("daily");
  const [result, setResult] = useState<WorkingSession | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const session = await api.createWorkingSession({
        agent_id: agentID,
        daily_labour_power_value: dailyPence,
        working_day_minutes: workingDayHours * 60,
        overtime_hours: overtimeHours,
        overtime_rate_pence: overtimeRate,
        wage_period: wagePeriod,
      });
      setResult(session);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    }
  }

  const livePrice =
    workingDayHours > 0
      ? (dailyPence / workingDayHours).toFixed(4)
      : "—";

  return (
    <div className="ch20-time-wage">
      <TimeWageInsight
        dailyPence={dailyPence}
        workingDayHours={workingDayHours}
        onHoursChange={setWorkingDayHours}
      />
      <section className="ch20-section">
        <h2>Enter Working Session</h2>
        <p className="ch20-explainer">
          The time-wage expresses the price of labour as a price per unit of
          time. As the working day lengthens without a change in the daily value
          of labour-power, the hourly price falls — even when nominal pay stays
          the same.
        </p>

        <div className="ch20-live-hint">
          Live hourly price: <strong>{livePrice}d.</strong> ({dailyPence}p ÷{" "}
          {workingDayHours}h)
        </div>

        <form onSubmit={handleSubmit} className="ch20-form">
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
            Daily value of labour-power (pence)
            <input
              type="number"
              min={1}
              value={dailyPence}
              onChange={(e) => setDailyPence(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Normal working day (hours)
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
            Overtime hours
            <input
              type="number"
              min={0}
              value={overtimeHours}
              onChange={(e) => setOvertimeHours(Number(e.target.value))}
            />
          </label>
          <label>
            Overtime rate (pence/hour)
            <input
              type="number"
              min={0}
              value={overtimeRate}
              onChange={(e) => setOvertimeRate(Number(e.target.value))}
            />
          </label>
          <label>
            Wage period
            <select
              value={wagePeriod}
              onChange={(e) =>
                setWagePeriod(e.target.value as "daily" | "weekly")
              }
            >
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
            </select>
          </label>
          <button type="submit">Compute Time-Wage</button>
        </form>
        {error && <p className="ch20-error">{error}</p>}
      </section>

      {result && (
        <section className="ch20-section ch20-results">
          <h2>Results</h2>

          <div className="ch20-cards">
            <div className="ch20-card">
              <h3>Hourly Price of Labour</h3>
              {!modern && (
                <p className="ch20-card-subtitle">Exact rational fraction</p>
              )}
              <dl>
                {!modern && (
                  <>
                    <dt>Fraction</dt>
                    <dd>
                      {result.hourly_price.numerator}/{result.hourly_price.denominator}
                    </dd>
                  </>
                )}
                <dt>Decimal</dt>
                <dd>
                  {modern
                    ? fmt(result.hourly_price.as_float)
                    : `${result.hourly_price.as_float.toFixed(4)}d.`}
                </dd>
                <dt>Daily value</dt>
                <dd>
                  {modern
                    ? fmt(result.daily_labour_power_value.pence)
                    : `${result.daily_labour_power_value.pence}p`}
                </dd>
                <dt>Normal hours</dt>
                <dd>{result.working_day_minutes.minutes / 60}h</dd>
              </dl>
              <p className="ch20-note">
                A longer day lowers the hourly price even when total daily pay is
                unchanged — the worker's labour becomes cheaper per hour.
              </p>
            </div>

            <div className="ch20-card ch20-card-wage">
              <h3>Nominal Wage</h3>
              <p className="ch20-card-subtitle">Total money received</p>
              <dl>
                <dt>Normal pay</dt>
                <dd>
                  {modern
                    ? fmt(result.daily_labour_power_value.pence)
                    : `${result.daily_labour_power_value.pence}p`}
                </dd>
                {result.overtime_hours.hours > 0 && (
                  <>
                    <dt>
                      Overtime ({result.overtime_hours.hours}h &times;{" "}
                      {modern
                        ? fmt(result.overtime_rate_pence.pence)
                        : `${result.overtime_rate_pence.pence}d`})
                    </dt>
                    <dd>
                      {modern
                        ? fmt(result.overtime_hours.hours * result.overtime_rate_pence.pence)
                        : `${result.overtime_hours.hours * result.overtime_rate_pence.pence}p`}
                    </dd>
                  </>
                )}
                <dt>Total nominal wage</dt>
                <dd className="ch20-total">
                  {modern
                    ? fmt(result.nominal_wage.pence)
                    : `${result.nominal_wage.pence}p`}
                </dd>
              </dl>
              <p className="ch20-note">
                Overtime earns more gross but the capitalist still extracts
                unpaid labour: the overtime rate pays only a fraction of the
                value actually produced in the extra hour.
              </p>
            </div>
          </div>
        </section>
      )}
      <PriceOscillationCoda />
    </div>
  );
}
