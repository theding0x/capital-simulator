import { useMemo } from "react";
import type { RentSnapshot } from "../types";
import { formatBP, formatPence } from "./animation";

interface RentTerrainProps {
  rent: RentSnapshot;
}

const W = 680;
const H = 258;
const PAD_L = 16;
const PAD_R = 16;
const PAD_T = 20;
const PAD_B = 44;

/** Vol III · Ch. 37–47 — the ground-rent terrain. Each soil is a band: equal
 *  capital yields an equal price of production (lead) plus the absolute rent the
 *  monopoly of landed property exacts on every soil (red) — together the worst
 *  soil's regulating market price, the waterline. More fertile soils clear more
 *  output at that price, and the surplus profit rises above the waterline as
 *  differential rent (gold). */
export function RentTerrain({ rent }: RentTerrainProps) {
  const geom = useMemo(() => {
    const grades = rent.grades ?? [];
    if (grades.length === 0) return null;

    const innerW = W - PAD_L - PAD_R;
    const innerH = H - PAD_T - PAD_B;
    const baseY = PAD_T + innerH;

    // Each soil's realized price = price of production + absolute rent (to the
    // waterline) + differential rent (above it). Tallest band sets the scale.
    const priceProd = rent.waterline_pence - (grades[0]?.absolute_rent_pence ?? 0);
    const stackOf = (g: (typeof grades)[number]) =>
      priceProd + g.absolute_rent_pence + g.differential_rent_pence;
    const maxStack = Math.max(...grades.map(stackOf), 1);
    const scale = (v: number) => (v / maxStack) * innerH;

    const n = grades.length;
    const slot = innerW / n;
    const barW = Math.min(slot * 0.62, 96);

    const bands = grades.map((g, i) => {
      const cx = PAD_L + slot * (i + 0.5);
      const x = cx - barW / 2;
      const hBase = scale(priceProd);
      // Floor the absolute-rent band so the monopoly levy stays legible on every
      // soil — including the worst, whose whole rent it is.
      const hAbs = g.absolute_rent_pence > 0 ? Math.max(scale(g.absolute_rent_pence), 4) : 0;
      const hDiff = scale(g.differential_rent_pence);
      // Stack from the baseline up: price of production, absolute, differential.
      const yBase = baseY - hBase;
      const yAbs = yBase - hAbs;
      const yDiff = yAbs - hDiff;
      return { g, cx, x, barW, hBase, hAbs, hDiff, yBase, yAbs, yDiff };
    });

    // The waterline = the regulating worst-soil price = price of production + AR.
    const waterY = baseY - scale(rent.waterline_pence);
    return { bands, baseY, waterY, innerW };
  }, [rent]);

  if (!geom) {
    return <p className="chart-empty">The totality has not yet run &mdash; start the engine.</p>;
  }

  return (
    <div className="rent-terrain">
      <svg
        viewBox={`0 0 ${W} ${H}`}
        width="100%"
        preserveAspectRatio="none"
        role="img"
        aria-label="ground-rent terrain: soil grades over the worst-soil waterline; fertile soils rise above it as differential rent"
      >
        {/* Soil bands */}
        {geom.bands.map(({ g, x, barW, hBase, hAbs, hDiff, yBase, yAbs, yDiff }) => (
          <g key={g.id}>
            {/* price of production — equal capital across soils */}
            <rect x={x} y={yBase} width={barW} height={hBase} fill="rgba(74,90,138,0.42)" />
            {/* absolute rent — the monopoly levy on every soil */}
            {hAbs > 0 && (
              <rect x={x} y={yAbs} width={barW} height={hAbs} fill="rgba(192,57,43,0.5)" />
            )}
            {/* differential rent — surplus profit above the waterline */}
            {hDiff > 0 && (
              <rect x={x} y={yDiff} width={barW} height={hDiff} fill="#c8a240" />
            )}
            {/* grade label */}
            <text x={x + barW / 2} y={geom.baseY + 14} textAnchor="middle" className="rt-axis">
              {g.id}
            </text>
            <text x={x + barW / 2} y={geom.baseY + 26} textAnchor="middle" className="rt-axis dim">
              {g.yield_units}&times; yield
            </text>
            {g.regulating && (
              <text x={x + barW / 2} y={geom.baseY + 38} textAnchor="middle" className="rt-axis reg">
                regulates
              </text>
            )}
          </g>
        ))}

        {/* The waterline — the regulating worst-soil market price */}
        <line
          x1={PAD_L}
          x2={W - PAD_R}
          y1={geom.waterY}
          y2={geom.waterY}
          stroke="#e8c660"
          strokeWidth="1.5"
          strokeDasharray="5 4"
        />
        <text x={W - PAD_R} y={geom.waterY - 6} textAnchor="end" className="rt-waterlbl">
          regulating price &middot; worst soil
        </text>
      </svg>

      <div className="rent-legend">
        <span className="lg lead">price of production</span>
        <span className="lg red">absolute rent &mdash; the monopoly levy</span>
        <span className="lg gold">differential rent &mdash; surplus profit</span>
        <span className="lg dim">&minus;&minus;&minus; waterline = worst-soil price</span>
      </div>

      <div className="rent-grades">
        {(rent.grades ?? []).map((g) => (
          <div className={"rent-grade" + (g.regulating ? " reg" : "")} key={g.id}>
            <span className="rg-id">{g.id}</span>
            <span className="rg-name">
              {g.name}
              {g.regulating && <span className="rg-reg"> &middot; regulates</span>}
            </span>
            <span className="rg-rent">{formatPence(g.total_rent_pence)}</span>
            <span className="rg-land">land {formatPence(g.land_price_pence)}</span>
          </div>
        ))}
      </div>

      <div className="rent-foot">
        <span>
          Differential <strong className="g">{formatPence(rent.differential_pence)}</strong>
        </span>
        <span>
          Absolute <strong className="r">{formatPence(rent.absolute_pence)}</strong>
        </span>
        <span>
          Total rent <strong className="g">{formatPence(rent.total_rent_pence)}</strong>
        </span>
        <span>
          Land price <strong>{formatPence(rent.land_price_pence)}</strong> @ {formatBP(rent.interest_rate_bp)}
        </span>
      </div>
    </div>
  );
}
