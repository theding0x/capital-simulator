import type { AbodeReadout } from "../types";
import { formatBP, formatPence, formatMinutes } from "./animation";
import { GeneralLawTrend } from "./GeneralLawTrend";
import { Levers } from "./Levers";

interface AbodeProps {
  abode: AbodeReadout;
}

/** The hidden abode of production. We have left "the noisy sphere, where
 *  everything takes place on the surface and in view of all men" for the place
 *  where the class relation — surplus as unpaid labour, the reserve army — is
 *  laid bare. The working day divides into necessary (paid, v) and surplus
 *  (unpaid, s) labour; their ratio is the rate of exploitation s/v. */
export function Abode({ abode }: AbodeProps) {
  const day = abode.necessary_labour_minutes + abode.surplus_labour_minutes || 1;
  const necPct = (abode.necessary_labour_minutes / day) * 100;
  const surPct = 100 - necPct;
  return (
    <div className="abode" data-testid="abode">
      <div className="abode-head">
        <h2>The hidden abode of production</h2>
        <p className="abode-hint">No admittance except on business.</p>
      </div>

      <section className="abode-card">
        <div className="abode-card-k">The social working day · s/v {formatBP(abode.rate_of_exploitation_bp)}</div>
        <div className="workingday" data-testid="workingday">
          <div className="wd-nec" style={{ width: `${necPct}%` }}>
            <span>necessary · {formatMinutes(abode.necessary_labour_minutes)}</span>
          </div>
          <div className="wd-sur" style={{ width: `${surPct}%` }}>
            <span>surplus · {formatMinutes(abode.surplus_labour_minutes)}</span>
          </div>
        </div>
      </section>

      <div className="abode-grid">
        <section className="abode-card">
          <div className="abode-card-k">Living labour · Σv (paid)</div>
          <div className="abode-card-v">{formatPence(abode.total_variable_pence)}</div>
          <div className="abode-hint">{abode.employed_count} employed</div>
        </section>
        <section className="abode-card gold">
          <div className="abode-card-k">Surplus extracted · Σs (unpaid)</div>
          <div className="abode-card-v">{formatPence(abode.total_surplus_pence)}</div>
          <div className="abode-hint">rises to the surface as capital</div>
        </section>
        <section className="abode-card">
          <div className="abode-card-k">Industrial reserve army</div>
          <div className="abode-card-v">{abode.reserve_army_count}</div>
          <div className="abode-hint">wage pressure {formatBP(abode.reserve_army_pressure_bp)} · paid wage {formatPence(abode.wage_pence)}</div>
        </section>
        <section className="abode-card">
          <div className="abode-card-k">Organic composition · c/v</div>
          <div className="abode-card-v">{formatBP(abode.organic_composition_bp)}</div>
          <div className="abode-hint">dead labour dominating living</div>
        </section>
      </div>

      <section className="abode-card">
        <div className="abode-card-k">The general law in motion</div>
        <GeneralLawTrend series={abode.law_series} />
      </section>

      <Levers abode={abode} />
    </div>
  );
}
