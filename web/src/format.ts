export function fmtMinutes(m: number): string {
  if (!Number.isFinite(m)) return "-";
  if (m < 60) return `${m}m`;
  const h = Math.floor(m / 60);
  const rem = m % 60;
  return rem === 0 ? `${h}h` : `${h}h ${rem}m`;
}

export function fmtQty(n: number): string {
  if (!Number.isFinite(n)) return "-";
  if (Number.isInteger(n)) return String(n);
  return n.toFixed(4).replace(/\.?0+$/, "");
}
