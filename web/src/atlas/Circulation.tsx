import type { CirculationSnapshot, DeptCVS } from "../types";
import { formatPence } from "./animation";

interface DeptBarProps {
  d: DeptCVS;
}

function DeptBar({ d }: DeptBarProps) {
  const total = d.c + d.v + d.s || 1;
  const segs: [string, number, string, string][] = [
    ["c", d.c, "var(--lead)", "#eef1ff"],
    ["v", d.v, "var(--lead-hover)", "#eef1ff"],
    ["s", d.s, "var(--gold-bright)", "#1a1305"],
  ];
  return (
    <div
      className="rep-deptbar"
      role="img"
      aria-label={`c ${d.c}, v ${d.v}, s ${d.s}`}
    >
      {segs.map(([k, val, bg, color]) => (
        <div
          key={k}
          className="rep-seg"
          style={{ width: (val / total) * 100 + "%", background: bg, color }}
        >
          <span className="rep-seg-k">{k}</span>
        </div>
      ))}
    </div>
  );
}

interface CirculationProps {
  circulation: CirculationSnapshot;
  reduced: boolean;
}

/** Vol II · the two-department reproduction scheme. */
export function Circulation({ circulation, reduced }: CirculationProps) {
  const c = circulation;
  const dI = c.departments.I;
  const dII = c.departments.II;
  return (
    <div className="strat-inner">
      <section className="rep-card">
        <div className="rep-head">
          <div>
            <div className="strat-eyebrow">
              Vol II &middot; Ch. 20&ndash;21 &mdash; reproduction of the total social capital
            </div>
            <h3 className="strat-h">The two departments</h3>
          </div>
          <div className="rep-mode">
            {c.extended ? "extended scale — accumulating" : "simple reproduction"}
          </div>
        </div>

        <div className="rep-grid">
          <div className="rep-dept">
            <div className="rep-dept-lbl">
              <b>Dept I</b> &middot; means of production{" "}
              <span className="rep-tot">{formatPence(dI.c + dI.v + dI.s)}</span>
            </div>
            <DeptBar d={dI} />
          </div>

          <div className="rep-exchange" aria-hidden="true">
            <span className="rep-x-lbl">I(v+s) &#x21C4; II(c)</span>
            <div className="rep-x-track">
              {!reduced && (
                <>
                  <span className="rep-x-dot"></span>
                  <span className="rep-x-dot d2"></span>
                  <span className="rep-x-dot d3"></span>
                </>
              )}
            </div>
            <span className="rep-x-val">
              {formatPence(c.keystone_exchange_pence)} &#x21C4; {formatPence(c.dept_ii_constant_pence)}
            </span>
          </div>

          <div className="rep-dept">
            <div className="rep-dept-lbl">
              <b>Dept II</b> &middot; means of consumption{" "}
              <span className="rep-tot">{formatPence(dII.c + dII.v + dII.s)}</span>
            </div>
            <DeptBar d={dII} />
          </div>
        </div>

        <div className="rep-foot">
          <span>
            <span className="sw lead"></span> c constant &middot;{" "}
            <span className="sw leadh"></span> v variable &middot;{" "}
            <span className="sw gold"></span> s surplus
          </span>
          <span className="rep-note">
            {c.extended
              ? "Dept I leads — surplus reinvested, c climbing relative to v. The rising organic composition is born here, two volumes before the abode shows its effects."
              : "The keystone exchange closes the loop: Dept I’s wages + surplus buy Dept II’s consumer goods; Dept II’s constant capital buys Dept I’s machines."}
          </span>
        </div>
      </section>
    </div>
  );
}
