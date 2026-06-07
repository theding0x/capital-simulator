import { createContext, useContext, useState } from "react";
import type { ReactNode } from "react";

import { fmtPounds, fmtPoundsModern } from "./format";

interface CurrencyContextShape {
  modern: boolean;
  toggle: () => void;
}

const STORE_KEY = "atlas.currency.modern";

function loadModern(): boolean {
  try {
    return localStorage.getItem(STORE_KEY) === "true";
  } catch {
    return false;
  }
}

const CurrencyContext = createContext<CurrencyContextShape>({
  modern: false,
  toggle: () => {},
});

export function CurrencyProvider({ children }: { children: ReactNode }) {
  const [modern, setModern] = useState(loadModern);
  const toggle = () =>
    setModern((m) => {
      const next = !m;
      try {
        localStorage.setItem(STORE_KEY, String(next));
      } catch {
        /* storage unavailable — non-fatal */
      }
      return next;
    });
  return (
    <CurrencyContext.Provider value={{ modern, toggle }}>
      {children}
    </CurrencyContext.Provider>
  );
}

export function useCurrency(): CurrencyContextShape {
  return useContext(CurrencyContext);
}

export function usePounds(): (pence: number) => string {
  const { modern } = useCurrency();
  return modern ? fmtPoundsModern : fmtPounds;
}

export function CurrencyToggle() {
  const { modern, toggle } = useCurrency();
  return (
    <button
      className={`currency-toggle${modern ? " currency-toggle--active" : ""}`}
      onClick={toggle}
      title={
        modern
          ? "Showing 2025 prices — click for 1860s historical"
          : "Showing 1860s prices — click for 2025 equivalent"
      }
    >
      {modern ? "2025 £" : "1860s £"}
    </button>
  );
}
