export interface CapitalCompositionValues {
  constant: number;
  variable: number;
  surplusRate: number;
  periods: number;
}

interface Preset {
  label: string;
  values: CapitalCompositionValues;
}

interface Props {
  values: CapitalCompositionValues;
  maxPeriods?: number;
  presets: Preset[];
  onChange: (next: CapitalCompositionValues) => void;
}

export function CapitalCompositionForm({ values, maxPeriods = 20, presets, onChange }: Props) {
  return (
    <>
      {presets.length > 0 && (
        <div className="preset-row">
          {presets.map((p) => (
            <button
              key={p.label}
              type="button"
              className="preset-btn"
              onClick={() => onChange(p.values)}
            >
              {p.label}
            </button>
          ))}
        </div>
      )}
      <div className="form-row">
        <label className="form-label">
          Constant capital c (pence)
          <input
            type="number"
            min={0}
            className="form-input"
            value={values.constant}
            onChange={(e) => onChange({ ...values, constant: Number(e.target.value) })}
            required
          />
        </label>
        <label className="form-label">
          Variable capital v (pence)
          <input
            type="number"
            min={1}
            className="form-input"
            value={values.variable}
            onChange={(e) => onChange({ ...values, variable: Number(e.target.value) })}
            required
          />
        </label>
      </div>
      <div className="form-row">
        <label className="form-label">
          Rate of surplus-value s′
          <input
            type="number"
            min={0}
            step={0.01}
            className="form-input"
            value={values.surplusRate}
            onChange={(e) => onChange({ ...values, surplusRate: Number(e.target.value) })}
            required
          />
        </label>
        <label className="form-label">
          Periods: {values.periods}
          <input
            type="range"
            min={1}
            max={maxPeriods}
            className="form-input"
            value={values.periods}
            onChange={(e) => onChange({ ...values, periods: Number(e.target.value) })}
          />
        </label>
      </div>
    </>
  );
}
