import { useCallback, useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import { fmtHoursLong } from "../../format";
import type { LabourScenarioResult } from "../../types";
import "./Ch17MagnitudeChanges.css";

const BASELINE = { workingDay: 720, necessary: 480, intensity: 1.0, productivity: 1.0 };

const LAWS = [
  { id: "law1", num: "Law 1", name: "Productivity" },
  { id: "law2", num: "Law 2", name: "Intensity" },
  { id: "law3", num: "Law 3", name: "Duration" },
  { id: "law4", num: "Law 4", name: "Combination" },
] as const;

function diagnoseLaw(
  workingDay: number,
  necessary: number,
  intensity: number,
  productivity: number,
): string {
  const dayMoved = workingDay !== BASELINE.workingDay;
  const intensityMoved = Math.abs(intensity - BASELINE.intensity) > 0.001;
  const productivityMoved =
    Math.abs(productivity - BASELINE.productivity) > 0.001 ||
    necessary !== BASELINE.necessary;
  const moves = [intensityMoved, dayMoved, productivityMoved].filter(Boolean).length;
  if (moves === 0) return "";
  if (moves > 1) return "law4";
  if (productivityMoved) return "law1";
  if (intensityMoved) return "law2";
  if (dayMoved) return "law3";
  return "";
}

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
  return (
    <>
      <MagnitudeInsight />
      <MagnitudeCalculator />
    </>
  );
}

function MagnitudeInsight() {
  return (
    <section className="v1-ch17-insight">
      <h2 className="v1-ch17-insight-h2">Three independent magnitudes</h2>
      <p className="v1-ch17-insight-prose">
        The working-day partition responds to three knobs that can be turned
        independently. The §1 law (constant daily value) holds only along the
        productivity axis; moving the intensity axis breaks it.
      </p>
      <div className="v1-ch17-axes">
        <div className="v1-ch17-axis">
          <span className="v1-ch17-axis-tag">§3 Axis A</span>
          <span className="v1-ch17-axis-name">Duration</span>
        </div>
        <div className="v1-ch17-axis">
          <span className="v1-ch17-axis-tag">§2 Axis B</span>
          <span className="v1-ch17-axis-name">Intensity</span>
        </div>
        <div className="v1-ch17-axis">
          <span className="v1-ch17-axis-tag">§1 Axis C</span>
          <span className="v1-ch17-axis-name">Productivity</span>
        </div>
      </div>
    </section>
  );
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
      <div className="v1-ch17-presets">
        <span className="v1-ch17-presets-label">Marx fixture</span>
        {PRESETS.map((p) => (
          <button
            key={p.id}
            type="button"
            className={
              "v1-ch17-preset-button" +
              (activePreset === p.id ? " v1-ch17-preset-button--active" : "")
            }
            onClick={() => applyPreset(p)}
          >
            {p.section} {p.label.replace(/^§\d[A-D]?\s—\s/, "")}
          </button>
        ))}
      </div>
      {active && (
        <p className="small muted" style={{ marginTop: 0, marginBottom: "1rem" }}>
          {active.caption}
        </p>
      )}

      <div className="v1-ch17-sliders">
        <div className="v1-ch17-slider-row">
          <span
            className={
              "v1-ch17-slider-label" +
              (workingDay !== BASELINE.workingDay ? " v1-ch17-slider-label--moved" : "")
            }
          >
            Duration
          </span>
          <input
            type="range"
            className="v1-ch17-slider"
            min={120}
            max={1440}
            step={10}
            value={workingDay}
            onChange={(e) => {
              setActivePreset("");
              setWorkingDay(Number(e.target.value));
            }}
          />
          <span className="v1-ch17-slider-value">{fmtHoursLong(workingDay)}</span>
        </div>
        <div className="v1-ch17-slider-row">
          <span
            className={
              "v1-ch17-slider-label" +
              (Math.abs(intensity - BASELINE.intensity) > 0.001
                ? " v1-ch17-slider-label--moved"
                : "")
            }
          >
            Intensity
          </span>
          <input
            type="range"
            className="v1-ch17-slider"
            min={0.5}
            max={2.0}
            step={0.01}
            value={intensity}
            onChange={(e) => {
              setActivePreset("");
              setIntensity(Number(e.target.value));
            }}
          />
          <span className="v1-ch17-slider-value">×{intensity.toFixed(2)}</span>
        </div>
        <div className="v1-ch17-slider-row">
          <span
            className={
              "v1-ch17-slider-label" +
              (Math.abs(productivity - BASELINE.productivity) > 0.001
                ? " v1-ch17-slider-label--moved"
                : "")
            }
          >
            Productivity
          </span>
          <input
            type="range"
            className="v1-ch17-slider"
            min={0.5}
            max={2.0}
            step={0.01}
            value={productivity}
            onChange={(e) => {
              setActivePreset("");
              setProductivity(Number(e.target.value));
            }}
          />
          <span className="v1-ch17-slider-value">×{productivity.toFixed(2)}</span>
        </div>
        <div className="v1-ch17-slider-row">
          <span
            className={
              "v1-ch17-slider-label" +
              (necessary !== BASELINE.necessary ? " v1-ch17-slider-label--moved" : "")
            }
          >
            Necessary
          </span>
          <input
            type="range"
            className="v1-ch17-slider"
            min={60}
            max={workingDay - 1}
            step={10}
            value={necessary}
            onChange={(e) => {
              setActivePreset("");
              setNecessary(Number(e.target.value));
            }}
          />
          <span className="v1-ch17-slider-value">{fmtHoursLong(necessary)}</span>
        </div>
      </div>

      <form onSubmit={submit} style={{ display: "none" }} aria-hidden="true" />
      {err && <p className="error">{err}</p>}

      {result && (
        <>
          <div className="v1-ch17-result">
            {(() => {
              const total =
                result.necessary_labour_minutes + result.surplus_labour_minutes;
              const max = Math.max(total, 1);
              const pct = (n: number) => `${(n / max) * 100}%`;
              return (
                <div className="v1-ch17-bar">
                  <span className="v1-ch17-bar-label">Working day</span>
                  <div className="v1-ch17-bar-track">
                    <div
                      className="v1-ch17-bar-seg--necessary"
                      style={{ flexBasis: pct(result.necessary_labour_minutes) }}
                      title={`v: ${fmtHoursLong(result.necessary_labour_minutes)}`}
                    />
                    <div
                      className="v1-ch17-bar-seg--surplus"
                      style={{ flexBasis: pct(result.surplus_labour_minutes) }}
                      title={`s: ${fmtHoursLong(result.surplus_labour_minutes)}`}
                    />
                  </div>
                  <span className="v1-ch17-bar-total">
                    {fmtHoursLong(total)}
                  </span>
                </div>
              );
            })()}
          </div>
          <div className="v1-ch17-rate">
            <span className="v1-ch17-rate-label">s / v</span>
            <span className="v1-ch17-rate-value">{surplusPct}%</span>
            <span style={{ color: "var(--ink-muted)", fontSize: "0.75rem" }}>
              · daily value {fmtHoursLong(result.daily_value_minutes)}{" "}
              {result.law_constant_daily_value ? "(Law 1 holds)" : "(Law 1 breaks)"}
            </span>
          </div>
          {(() => {
            const activeLaw = diagnoseLaw(workingDay, necessary, intensity, productivity);
            return (
              <div className="v1-ch17-law-block">
                {LAWS.map((l) => (
                  <div
                    key={l.id}
                    className={
                      "v1-ch17-law" +
                      (activeLaw === l.id ? " v1-ch17-law--active" : "")
                    }
                  >
                    <span className="v1-ch17-law-tag">{l.num}</span>
                    <span className="v1-ch17-law-name">{l.name}</span>
                  </div>
                ))}
              </div>
            );
          })()}
        </>
      )}

      <p className="v1-ch17-closing">
        §1, closing: “The value of labour-power, and the surplus-value, vary in
        opposite directions.” Move the Productivity slider alone to watch the
        necessary segment shrink and the surplus segment grow with the total
        held fixed.
      </p>
    </section>
  );
}
