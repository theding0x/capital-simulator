// Durable UI preferences for the Atlas page, persisted in localStorage so they
// survive a reload (unlike the simulation run, which resets). Defensive: bad or
// absent values fall back to defaults, and storage failures are non-fatal.
const KEY = "atlas.prefs.v1";

export interface AtlasPrefs {
  speed: number;
  reduced: boolean;
}

const DEFAULTS: AtlasPrefs = { speed: 1, reduced: false };

export function loadPrefs(): AtlasPrefs {
  try {
    const raw = localStorage.getItem(KEY);
    if (!raw) return { ...DEFAULTS };
    const p = JSON.parse(raw) as Partial<AtlasPrefs>;
    return {
      speed: typeof p.speed === "number" ? p.speed : DEFAULTS.speed,
      reduced: typeof p.reduced === "boolean" ? p.reduced : DEFAULTS.reduced,
    };
  } catch {
    return { ...DEFAULTS };
  }
}

export function savePrefs(p: AtlasPrefs): void {
  try {
    localStorage.setItem(KEY, JSON.stringify(p));
  } catch {
    /* storage unavailable (e.g. private mode) — non-fatal */
  }
}
