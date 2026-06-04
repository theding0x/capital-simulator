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
