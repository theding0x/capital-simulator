import type { AbodeReadout } from "../types";
import { formatPence, formatBP, formatMinutes } from "./animation";
import { ImmiserationChart } from "./ImmiserationChart";
import { ReserveArmy } from "./ReserveArmy";
import { Levers } from "./Levers";

/** The social working day (hero): necessary (paid, v) vs surplus (unpaid, s); the
 *  ratio s/v is the rate of exploitation. */
function WorkingDay({ abode }: { abode: AbodeReadout }) {
  const day = abode.necessary_labour_minutes + abode.surplus_labour_minutes || 1;
  const necPct = (abode.necessary_labour_minutes / day) * 100;
  const surPct = 100 - necPct;
  return (
    <section className="wd">
      <div className="wd-head">
        <div className="wd-eyebrow">The social working day</div>
        <div className="wd-sv">
          <span className="wd-sv-num">{formatBP(abode.rate_of_exploitation_bp)}</span>
          <span className="wd-sv-lbl">rate of exploitation · s / v</span>
        </div>
      </div>
      <div
        className="wd-bar"
        data-testid="workingday"
        role="img"
        aria-label={`necessary ${formatMinutes(abode.necessary_labour_minutes)}, surplus ${formatMinutes(abode.surplus_labour_minutes)}`}
      >
        <div className="wd-nec" style={{ width: `${necPct}%` }}>
          <span className="wd-seg-k">necessary</span>
          <span className="wd-seg-v">{formatMinutes(abode.necessary_labour_minutes)}</span>
        </div>
        <div className="wd-sur" style={{ width: `${surPct}%` }}>
          <span className="wd-seg-k">surplus · unpaid</span>
          <span className="wd-seg-v">{formatMinutes(abode.surplus_labour_minutes)}</span>
        </div>
      </div>
      <div className="wd-foot">
        <span>
          <span className="sw nec"></span> paid — reproduces the value of labour-power
        </span>
        <span>
          <span className="sw sur"></span> unpaid — pumped out as surplus-value
        </span>
      </div>
    </section>
  );
}

/** Demoted stat tiles: living labour Σv, surplus extracted Σs, organic composition c/v. */
function StatRow({ abode }: { abode: AbodeReadout }) {
  const tiles = [
    {
      k: "Living labour · Σv",
      v: formatPence(abode.total_variable_pence),
      sub: `${abode.employed_count} employed`,
      tone: "lead",
    },
    {
      k: "Surplus extracted · Σs",
      v: formatPence(abode.total_surplus_pence),
      sub: "rises to the surface as capital",
      tone: "gold",
    },
    {
      k: "Organic composition · c/v",
      v: formatBP(abode.organic_composition_bp),
      sub: "dead labour dominating living",
      tone: "",
    },
  ];
  return (
    <div className="stat-row">
      {tiles.map((t) => (
        <div className={"stat" + (t.tone ? " " + t.tone : "")} key={t.k}>
          <div className="stat-k">{t.k}</div>
          <div className="stat-v">{t.v}</div>
          <div className="stat-sub">{t.sub}</div>
        </div>
      ))}
    </div>
  );
}

/** The hidden abode: working-day hero, demoted tiles, the immiseration chart
 *  (promoted), the reserve army, the levers. */
export function Abode({ abode }: { abode: AbodeReadout }) {
  return (
    <div className="abode-inner" data-testid="abode">
      <WorkingDay abode={abode} />
      <StatRow abode={abode} />
      <section className="chart-card">
        <div className="chart-eyebrow">The general law in motion</div>
        <p className="chart-gloss">
          Accumulation widens the rate of exploitation, swells the reserve army, and presses the
          wage to its floor — the immiseration of the producer, run in real time.
        </p>
        <ImmiserationChart series={abode.law_series} />
      </section>
      <ReserveArmy abode={abode} />
      <Levers abode={abode} />
    </div>
  );
}
