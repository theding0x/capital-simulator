import { useCallback, useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { fmtHoursLong } from "../../format";
import type { LabourScenarioResult } from "../../types";

// Ch. 17 partitions a single working day under three independent
// magnitudes — duration, intensity, productivity — and reads the rate
// of surplus-value off the result. Marx isolates each factor in §§1-3
// then combines them in §4. The presets below replay each section so
// the user can step through the chapter's textual examples without
// hand-computing the partition each time.
type Preset = {
  id: string;
  label: string;
  section: string;
  workingDay: number;
  necessary: number;
  intensity: number;
  productivity: number;
  caption: string;
};

const PRESETS: Preset[] = [
  {
    id: "s1",
    label: "§1 — baseline 12h, productivity rise",
    section: "§1",
    workingDay: 720,
    necessary: 360,
    intensity: 1.0,
    productivity: 1.333,
    caption:
      "Necessary labour fell from 480 min (4 sh) to 360 min (3 sh): productivity rose by 4/3. Surplus rises to match. Daily value unchanged at 720.",
  },
  {
    id: "s2",
    label: "§2 — intensified 12h day",
    section: "§2",
    workingDay: 720,
    necessary: 360,
    intensity: 1.333,
    productivity: 1.0,
    caption:
      "Same 12-hour day, but each minute packs 4/3 the normal labour. Daily value rises to 960 min — Law 1 no longer holds.",
  },
  {
    id: "s3",
    label: "§3 — working day lengthened 12h → 14h",
    section: "§3",
    workingDay: 840,
    necessary: 360,
    intensity: 1.0,
    productivity: 1.0,
    caption:
      "Necessary labour fixed at 360 min; the extra 120 min becomes surplus. Absolute surplus-value mechanism (Ch. 16 §1) seen here as a partition shift.",
  },
  {
    id: "s4a",
    label: "§4A — LP value rises, day lengthens to compensate",
    section: "§4A",
    workingDay: 960,
    necessary: 480,
    intensity: 1.0,
    productivity: 1.0,
    caption:
      "Subsistence demands 4 sh (480 min) and the day stretches to 16 h. Absolute surplus rises 33⅓% (360 → 480 min) even as the worker's daily value climbs.",
  },
];

export function Ch17MagnitudeChanges() {
  return <MagnitudeCalculator />;
}

function MagnitudeCalculator() {
  const [workingDay, setWorkingDay] = useState(PRESETS[0].workingDay);
  const [necessary, setNecessary] = useState(PRESETS[0].necessary);
  const [intensity, setIntensity] = useState(PRESETS[0].intensity);
  const [productivity, setProductivity] = useState(PRESETS[0].productivity);
  const [activePreset, setActivePreset] = useState<string>(PRESETS[0].id);
  const [result, setResult] = useState<LabourScenarioResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  const recompute = useCallback(async () => {
    setErr(null);
    try {
      const r = await api.computeLabourScenario({
        working_day_minutes: workingDay,
        necessary_labour_minutes: necessary,
        intensity_factor: intensity,
        productivity_factor: productivity,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
      setResult(null);
    }
  }, [workingDay, necessary, intensity, productivity]);

  useEffect(() => {
    void recompute();
  }, [recompute]);

  function applyPreset(p: Preset) {
    setActivePreset(p.id);
    setWorkingDay(p.workingDay);
    setNecessary(p.necessary);
    setIntensity(p.intensity);
    setProductivity(p.productivity);
  }

  function submit(e: FormEvent) {
    e.preventDefault();
    void recompute();
  }

  const active = PRESETS.find((p) => p.id === activePreset);
  const surplusPct = result && result.necessary_labour_minutes > 0
    ? (result.rate_of_surplus_value * 100).toFixed(1)
    : "—";

  return (
    <section className="card">
      <h2>Magnitude calculator (§§1-4)</h2>
      <p className="description">
        Three independent magnitudes — duration, intensity, productivity
        — partition the working day into necessary and surplus labour.
        The §1 law says daily value stays constant at normal intensity;
        the §2 law says intensity scales it; §3 lengthens the day
        directly; §4 combines all three. Cycle through the presets to
        see each law in isolation.
      </p>
      <div className="form-grid" style={{ marginBottom: "0.75rem" }}>
        {PRESETS.map((p) => (
          <button
            key={p.id}
            type="button"
            className={`btn ${activePreset === p.id ? "btn-primary" : ""}`}
            onClick={() => applyPreset(p)}
          >
            <strong style={{ marginRight: "0.4rem" }}>{p.section}</strong>
            {p.label.replace(/^§\d[A-D]?\s—\s/, "")}
          </button>
        ))}
      </div>
      {active && (
        <p className="small muted" style={{ marginTop: 0, marginBottom: "1rem" }}>
          {active.caption}
        </p>
      )}
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Working day (min)</span>
          <input
            type="number"
            min={60}
            max={1440}
            value={workingDay}
            onChange={(e) => {
              setActivePreset("");
              setWorkingDay(Number(e.target.value));
            }}
          />
        </label>
        <label>
          <span>Necessary labour (min)</span>
          <input
            type="number"
            min={1}
            max={workingDay - 1}
            value={necessary}
            onChange={(e) => {
              setActivePreset("");
              setNecessary(Number(e.target.value));
            }}
          />
        </label>
        <label>
          <span>Intensity factor (1.0 = normal)</span>
          <input
            type="number"
            min={0.1}
            step={0.05}
            value={intensity}
            onChange={(e) => {
              setActivePreset("");
              setIntensity(Number(e.target.value));
            }}
          />
        </label>
        <label>
          <span>Productivity factor (1.0 = unchanged)</span>
          <input
            type="number"
            min={0.1}
            step={0.05}
            value={productivity}
            onChange={(e) => {
              setActivePreset("");
              setProductivity(Number(e.target.value));
            }}
          />
        </label>
      </form>
      {err && <p className="error">{err}</p>}
      {result && (
        <table className="data-table" style={{ marginTop: "1rem" }}>
          <tbody>
            <tr>
              <td>Daily value created</td>
              <td>
                <strong>{fmtHoursLong(result.daily_value_minutes)}</strong>{" "}
                <span className="muted">({result.daily_value_minutes} min)</span>
              </td>
            </tr>
            <tr>
              <td>Necessary labour (= value of labour-power)</td>
              <td>
                {fmtHoursLong(result.necessary_labour_minutes)}{" "}
                <span className="muted">({result.necessary_labour_minutes} min)</span>
              </td>
            </tr>
            <tr>
              <td>Surplus labour</td>
              <td>
                <strong>{fmtHoursLong(result.surplus_labour_minutes)}</strong>{" "}
                <span className="muted">({result.surplus_labour_minutes} min)</span>
              </td>
            </tr>
            <tr>
              <td>Rate of surplus-value (s/v)</td>
              <td>
                <strong>{surplusPct}%</strong>
              </td>
            </tr>
            <tr>
              <td>§1 Law 1 — daily value constant?</td>
              <td>
                {result.law_constant_daily_value ? (
                  <span>holds (normal intensity)</span>
                ) : (
                  <span className="muted">
                    fails — intensity ≠ 1.0 lifts the daily value
                  </span>
                )}
              </td>
            </tr>
          </tbody>
        </table>
      )}
      <p className="small muted" style={{ marginTop: "1rem" }}>
        §1, closing: "The value of labour-power, and the surplus-value,
        vary in opposite directions." Cycle §1 ↔ baseline to see this
        directly: the value-of-labour-power column moves down as the
        surplus column moves up, and the daily-value total holds steady.
      </p>
    </section>
  );
}
