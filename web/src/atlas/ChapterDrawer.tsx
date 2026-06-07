import { useEffect, useRef } from "react";
import { CHAPTER_PANELS } from "../components/ChapterShell";
import { chapterId, volRoman, type InventoryEntry } from "./inventory";
import "./spine.css";

interface ChapterDrawerProps {
  entry: InventoryEntry;
  onClose: () => void;
}

/** Focus-trap: keep Tab/Shift+Tab within the drawer. */
function useFocusTrap(ref: React.RefObject<HTMLElement | null>, active: boolean) {
  useEffect(() => {
    if (!active || !ref.current) return;
    const el = ref.current;
    const focusable = (): NodeListOf<HTMLElement> =>
      el.querySelectorAll<HTMLElement>(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key !== "Tab") return;
      const els = Array.from(focusable());
      if (els.length === 0) return;
      const first = els[0], last = els[els.length - 1];
      if (e.shiftKey) {
        if (document.activeElement === first) { e.preventDefault(); last.focus(); }
      } else {
        if (document.activeElement === last) { e.preventDefault(); first.focus(); }
      }
    };
    el.addEventListener("keydown", handleKeyDown);
    return () => el.removeEventListener("keydown", handleKeyDown);
  }, [active, ref]);
}

export function ChapterDrawer({ entry, onClose }: ChapterDrawerProps) {
  const drawerRef = useRef<HTMLDivElement>(null);
  const closeRef = useRef<HTMLButtonElement>(null);
  // Track the element that had focus before the drawer opened, for restoration.
  const priorFocusRef = useRef<Element | null>(document.activeElement);

  useFocusTrap(drawerRef, true);

  // Move focus into drawer on mount; restore on close.
  useEffect(() => {
    closeRef.current?.focus();
    return () => {
      const prior = priorFocusRef.current;
      if (prior instanceof HTMLElement) prior.focus();
    };
  }, []);

  // Escape closes the drawer.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [onClose]);

  const id = chapterId(entry);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const Panel = CHAPTER_PANELS[id] as React.ComponentType<any> | undefined;

  const volLabel = volRoman(entry.volume);
  const chNum = String(entry.ch).padStart(2, "0");

  return (
    <>
      {/* Backdrop */}
      <div
        className="spine-overlay"
        role="presentation"
        aria-hidden="true"
        onClick={onClose}
      />

      {/* Drawer */}
      <div
        ref={drawerRef}
        className="chapter-drawer"
        role="dialog"
        aria-modal="true"
        aria-label={`Chapter: Vol ${volLabel} Ch.${chNum} — ${entry.title}`}
      >
        <div className="drawer-header">
          <div className="drawer-meta">
            <div className="drawer-eyebrow">
              Vol {volLabel} &middot; Ch. {chNum}
            </div>
            <h2 className="drawer-title">{entry.title}</h2>
            <div className="drawer-part">{entry.part}</div>
          </div>
          <button
            ref={closeRef}
            className="drawer-close"
            aria-label="Close chapter drawer"
            onClick={onClose}
          >
            &times;
          </button>
        </div>

        <div className="drawer-body">
          {Panel ? (
            /* Render the real chapter panel with the minimal shared props it needs.
               Panels handle their own data fetching and empty states. */
            <Panel
              commodities={[]}
              owners={[]}
              onSharedChanged={() => undefined}
            />
          ) : (
            /* Fallback: title + part + link when panel not yet wired. */
            <div className="drawer-fallback">
              <div className="drawer-fallback-eyebrow">
                Vol {volLabel} &middot; Ch. {chNum}
              </div>
              <h3 className="drawer-fallback-title">{entry.title}</h3>
              <p className="drawer-fallback-part">{entry.part}</p>
              <a
                className="drawer-fallback-link"
                href="#/chapters"
                onClick={onClose}
              >
                Open in the chapter dashboard &#8594;
              </a>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
