import type { DistributionSnapshot } from "../types";
import { formatBP, formatPence } from "./animation";
import { TRPFChart } from "./TRPFChart";

interface TotalityProps {
  distribution: DistributionSnapshot;
}

/** Vol III · the totality — the tendential fall + distribution + the trinity fetish. */
export function Totality({ distribution }: TotalityProps) {
  const d = distribution;
  const tiles = [
    {
      k: "General rate · p̅′",
      v: formatBP(d.general_rate_bp),
      sub: "ΣS ÷ ΣC — the centre of gravity",
      tone: "gold",
    },
    {
      k: "Rate of interest",
      v: formatBP(d.interest_rate_bp),
      sub: "a share of average profit",
      tone: "lead",
    },
    {
      k: "Fictitious capital",
      v: formatPence(d.fictitious_pence),
      sub: "capitalised future surplus",
      tone: "",
    },
  ] as const;

  return (
    <div className="strat-inner">
      <section className="trpf-card">
        <div className="strat-eyebrow">
          Vol III &middot; Ch. 13&ndash;15 &mdash; the tendential fall of the rate of profit
        </div>
        <p className="strat-gloss">
          Rising organic composition with a constant rate of exploitation drives{" "}
          <span className="g">p&#x0305;&#x2032;</span> down even as the{" "}
          <span className="l">mass</span> of profit climbs &mdash; the rate&ndash;mass
          contradiction. Each crisis violently devalues constant capital and restores the
          rate.
        </p>
        <TRPFChart series={d.dist_series} />
      </section>

      <div className="stat-row">
        {tiles.map((t) => (
          <div className={"stat" + (t.tone ? " " + t.tone : "")} key={t.k}>
            <div className="stat-k">{t.k}</div>
            <div className="stat-v">{t.v}</div>
            <div className="stat-sub">{t.sub}</div>
          </div>
        ))}
      </div>

      <section className="trinity-card">
        <div className="strat-eyebrow" style={{ color: "var(--gold)" }}>
          Vol III &middot; Ch. 48 &mdash; the trinity formula
        </div>
        <div className="trinity-row">
          <div className="trinity-cell">
            <span className="tc-src">Capital</span>
            <span className="tc-arrow">&rarr;</span>
            <span className="tc-rev">interest</span>
          </div>
          <div className="trinity-cell">
            <span className="tc-src">Land</span>
            <span className="tc-arrow">&rarr;</span>
            <span className="tc-rev">rent</span>
          </div>
          <div className="trinity-cell">
            <span className="tc-src">Labour</span>
            <span className="tc-arrow">&rarr;</span>
            <span className="tc-rev">wages</span>
          </div>
        </div>
        <p className="trinity-note">
          <em>
            &ldquo;An enchanted, perverted, topsy-turvy world.&rdquo;
          </em>{" "}
          The culminating mystification &mdash; each revenue appears to spring from its
          own source, when all of it derives from surplus-labour alone.
        </p>
        <p className="trinity-note" style={{ marginTop: "8px", color: "var(--ink-dim)", fontSize: "0.74rem" }}>
          Note: rent ({formatPence(d.trinity.rent_pence)}) is not yet modelled in this simulation &mdash; ground-rent (Vol III Part VI) is represented as zero.
        </p>
      </section>
    </div>
  );
}
