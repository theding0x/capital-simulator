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

// PENCE_PER_POUND names the scale a wire money field is stored on: how many
// of the field's units make £1. Money in this codebase is not stored on one
// scale, and the web deliberately does not rescale wire values; each call site
// names the scale of the field it formats.
//
//   POUNDS      1    amount is already in whole pounds (Vol. I Ch.23-25
//                    reproduction schemes).
//   DECIMAL     100  100 pence to the pound (decimal pence), e.g. agent
//                    balances and the circulation chapters.
//   PRE_DECIMAL 240  240 pence to the pound, i.e. pkg/money.Pence
//                    (£1 = 20s = 240d).
//   FARTHINGS   960  960 farthings to the pound (Vol. I Ch.21 piece-wage
//                    fields, despite their `_pence` JSON names).
//
// The Go side lives in pkg/money (PencePerPound = 240).
export const PENCE_PER_POUND = {
  POUNDS: 1,
  DECIMAL: 100,
  PRE_DECIMAL: 240,
  FARTHINGS: 960,
} as const;

export type PencePerPound = (typeof PENCE_PER_POUND)[keyof typeof PENCE_PER_POUND];

// fmtPounds renders a money amount stored at `perPound` units to the pound as
// "£X.YZ", e.g. 10000 at scale 100 is "£100.00" and 120000 at scale 240 is
// "£500.00". The sub-penny remainder is rounded to the nearest hundredth of a
// pound, half away from zero (matches Go money.roundDiv). A negative amount
// renders with a leading minus sign unless it rounds to zero. Fractional
// amounts are accepted. `perPound` is required so that every call site states
// the scale of the field it formats.
export function fmtPounds(amount: number, perPound: PencePerPound): string {
  if (!Number.isFinite(amount)) return "—";
  const abs = Math.abs(amount);
  const cents = Math.floor((200 * abs + perPound) / (2 * perPound));
  const sign = amount < 0 && cents > 0 ? "-" : "";
  const pounds = Math.floor(cents / 100);
  const rem = cents % 100;
  return `${sign}£${pounds}.${rem.toString().padStart(2, "0")}`;
}

// MODERN_MULTIPLIER converts 1860s pounds to 2025 GBP equivalents (a
// pound-to-pound ratio, independent of the pence scale of the stored field).
// Source: Bank of England CPI calculator — £100 (1860s) = £15,782.99 (2025).
export const MODERN_MULTIPLIER = 157.8299;

// fmtPoundsModern renders a money amount stored at `perPound` units to the
// pound as its 2025 purchasing-power equivalent: £15,782.99 for 10000 at
// scale 100 (= £100 in 1860s).
export function fmtPoundsModern(amount: number, perPound: PencePerPound): string {
  if (!Number.isFinite(amount)) return "—";
  const sign = amount < 0 ? "-" : "";
  const modernPounds = (Math.abs(amount) / perPound) * MODERN_MULTIPLIER;
  return `${sign}£${modernPounds.toLocaleString("en-GB", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
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
