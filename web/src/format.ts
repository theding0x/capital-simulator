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

// fmtHoursLong renders LabourMinutes as "Xh" / "Xh Ym". Mirrors the
// Ch.13 / Ch.14 panels' previous minutesToHours helper.
export function fmtHoursLong(m: number): string {
  if (!Number.isFinite(m)) return "—";
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

// fmtPounds renders a pence integer as "£X.YZ". Negative values render
// with a leading minus sign.
export function fmtPounds(p: number): string {
  if (!Number.isFinite(p)) return "—";
  const sign = p < 0 ? "-" : "";
  const abs = Math.abs(Math.trunc(p));
  const pounds = Math.floor(abs / 100);
  const pence = abs % 100;
  return `${sign}£${pounds}.${pence.toString().padStart(2, "0")}`;
}

// fmtCompact renders a number with locale grouping. Used in the
// machinery panel for counts of units / labour-minutes.
export function fmtCompact(n: number): string {
  if (!Number.isFinite(n)) return "—";
  return n.toLocaleString("en-GB");
}

// fmtLabourMinutes pairs the raw count with an hours equivalent — the
// "1,200 lab-min (20.0 h)" format the Ch.15 panel uses for wear-and-tear.
export function fmtLabourMinutes(m: number): string {
  if (!Number.isFinite(m)) return "—";
  return `${fmtCompact(m)} lab-min (${(m / 60).toFixed(1)} h)`;
}
