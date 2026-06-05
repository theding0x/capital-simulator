import type { FieldCapital } from "../types";

/** Linear interpolation of `current` toward `target` by fraction `alpha` (0..1). */
export function ease(current: number, target: number, alpha: number): number {
  return current + (target - current) * alpha;
}

/** Arc fractions [money, production, commodity] summing to 1; all-money fallback. */
export function arcFractions(
  c: Pick<FieldCapital, "money_pence" | "production_pence" | "commodity_pence">
): [number, number, number] {
  const sum = c.money_pence + c.production_pence + c.commodity_pence;
  if (sum <= 0) return [1, 0, 0];
  return [c.money_pence / sum, c.production_pence / sum, c.commodity_pence / sum];
}

/** Deterministic per-capital spin tempo (seconds) from its id — the visual turnover. */
export function spinSeconds(id: string): number {
  let h = 0;
  for (let i = 0; i < id.length; i++) h = (h * 31 + id.charCodeAt(i)) >>> 0;
  return 3 + (h % 60) / 10; // 3.0–9.0s
}

/** Orbit pixel radius from total relative to the field max (area-proportional). */
export function orbitRadius(totalPence: number, maxTotal: number, min = 26, max = 74): number {
  if (maxTotal <= 0) return min;
  const t = Math.sqrt(Math.max(0, totalPence) / maxTotal);
  return Math.round(min + (max - min) * t);
}

/** £ with thousands separators (pence → pounds). */
export function formatPence(pence: number): string {
  return "£" + Math.round(pence / 100).toLocaleString("en-GB");
}

/** Basis points → percent string, e.g. 1615 → "16.2%". */
export function formatBP(bp: number): string {
  return (bp / 100).toFixed(1) + "%";
}

/** Laps per second for the dot: turnover_number scaled to a watchable base tempo. */
export function lapRateFor(turnoverNumber: number, speed: number): number {
  const n = turnoverNumber > 0 ? turnoverNumber : 1;
  // One turnover ≈ 8s at speed 1; higher turnover = proportionally faster.
  return (n / 8) * Math.max(1, speed);
}

/**
 * Map lap-progress p∈[0,1) to an angle (radians, 0 at top, clockwise) over arcs
 * sized [fm, fp, fc] (summing to 1), spending extra time in production (P) by
 * `lingerP`× so the dot visibly slows there. Pure + deterministic.
 */
export function pacedAngle(
  p: number,
  fm: number,
  fp: number,
  fc: number,
  lingerP = 2.4
): number {
  // Time weights: production over-weighted, then normalised to 1.
  const wm = fm;
  const wp = fp * lingerP;
  const wc = fc;
  const wsum = wm + wp + wc || 1;
  const tm = wm / wsum;
  const tp = wp / wsum;
  // Cumulative ANGLE boundaries (fraction of circle) for M, P, C.
  const aM = fm;
  const aP = fm + fp;
  let frac: number; // fraction around the circle, 0..1
  const q = ((p % 1) + 1) % 1;
  if (q < tm) {
    frac = (q / tm) * aM; // through M
  } else if (q < tm + tp) {
    frac = aM + ((q - tm) / tp) * (aP - aM); // through P (slow)
  } else {
    frac = aP + ((q - tm - tp) / (1 - tm - tp)) * (1 - aP); // through C
  }
  return frac * 2 * Math.PI; // 0 at top, clockwise (caller offsets)
}

/** Target CSS scale for an orbit of `totalPence` relative to the field max. */
export function targetScale(totalPence: number, maxTotal: number): number {
  if (maxTotal <= 0) return 1;
  return Math.max(0.4, Math.sqrt(Math.max(0, totalPence) / maxTotal));
}
