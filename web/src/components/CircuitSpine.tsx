import { Fragment } from "react";
import type { ReactNode } from "react";
import type { CircuitNode } from "../chapters/registry";

/**
 * The six motion-nodes of the circuit, in order of motion.
 *
 *   M  ──▶  M—C  ──▶  P  ──▶  C′  ──▶  C′—M′  ──▶  M′
 *   ╰─────────────────── ΔM (surplus) ─────────────────╯
 *
 * `delta-M`, `whole`, and `historical` are meta-nodes used only by
 * Sidebar filtering; they are not visualised on the spine.
 */
const SPINE_NODES: { node: CircuitNode; glyph: string; aria: string }[] = [
  { node: "M",         glyph: "M",      aria: "M — money capital advanced" },
  { node: "M-C",       glyph: "M—C",    aria: "M—C — purchase of labour-power and means of production" },
  { node: "P",         glyph: "P",      aria: "P — production and valorization" },
  { node: "C-prime",   glyph: "C′",     aria: "C′ — finished commodity carrying surplus-value" },
  { node: "C-M-prime", glyph: "C′—M′",  aria: "C′—M′ — sale, realising value in money" },
  { node: "M-prime",   glyph: "M′",     aria: "M′ — money capital returned" },
];

interface CircuitSpineProps {
  selectedNode: CircuitNode | null;
  onSelectNode: (node: CircuitNode) => void;
  /** Slot for the TurnoverPlayer controls. */
  children?: ReactNode;
}

export function CircuitSpine({ selectedNode, onSelectNode, children }: CircuitSpineProps) {
  return (
    <section className="circuit-spine" aria-label="The circuit of capital">
      <div className="circuit-spine-flow">
        {SPINE_NODES.map((n, i) => {
          const isActive = selectedNode === n.node;
          return (
            <Fragment key={n.node}>
              <button
                type="button"
                className={
                  "circuit-spine-node" +
                  (isActive ? " circuit-spine-node--active" : "")
                }
                onClick={() => onSelectNode(n.node)}
                aria-label={n.aria}
                aria-pressed={isActive}
                title={n.aria}
              >
                <span className="circuit-spine-node-glyph">{n.glyph}</span>
              </button>
              {i < SPINE_NODES.length - 1 && (
                <span className="circuit-spine-arrow" aria-hidden="true">──▶</span>
              )}
            </Fragment>
          );
        })}
      </div>
      <div className="circuit-spine-return" aria-hidden="true">
        <span className="circuit-spine-return-curve">╰</span>
        <span className="circuit-spine-return-line">────────────── ΔM (surplus-value) ──────────────</span>
        <span className="circuit-spine-return-curve">╯</span>
      </div>
      {children && <div className="circuit-spine-controls">{children}</div>}
    </section>
  );
}
