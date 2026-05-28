import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../../api";
import type { CircuitProof, ExchangeSimulation } from "../../types";
import { usePounds } from "../../CurrencyContext";
import "./Ch05Contradictions.css";

interface Ch05Props {
  onSharedChanged: () => void;
}

interface Candidates {
  cheating: boolean;
  nonEquivalents: boolean;
  usury: boolean;
}

const NO_CANDIDATES: Candidates = {
  cheating: false,
  nonEquivalents: false,
  usury: false,
};

export function Ch05Contradictions({ onSharedChanged: _onSharedChanged }: Ch05Props) {
  const [candidates, setCandidates] = useState<Candidates>(NO_CANDIDATES);

  function markUsury() {
    setCandidates((c) => ({ ...c, usury: true }));
  }

  function markExchangeRedistribution(aBefore: number, bBefore: number) {
    // Equal-value exchange that still shows redistribution → "cheating"
    // (one party walks away with more than they came in with). Non-equivalent
    // input also shows redistribution → "exchange at non-equivalents".
    setCandidates((c) => ({
      ...c,
      cheating: c.cheating || aBefore === bBefore,
      nonEquivalents: c.nonEquivalents || aBefore !== bBefore,
    }));
  }

  return (
    <>
      <InsightBox candidates={candidates} />
      <CircuitProbePanel onUsuryProven={markUsury} />
      <ExchangeSimulationPanel onRedistribution={markExchangeRedistribution} />
      <Coda />
    </>
  );
}

function InsightBox({ candidates }: { candidates: Candidates }) {
  const items: Array<{ key: keyof Candidates; label: string; gloss: string }> = [
    {
      key: "cheating",
      label: "Cheating (one party defrauds the other)",
      gloss:
        "If A overcharges B, A's gain is exactly B's loss. The sum of values in society is unchanged.",
    },
    {
      key: "nonEquivalents",
      label: "Exchange at non-equivalents",
      gloss:
        "Trading unequal values shifts wealth between parties; it does not enlarge the social total.",
    },
    {
      key: "usury",
      label: "Usurer's capital (M—M′)",
      gloss:
        "Money lent and money returned, with no commodity between them, produces ∆M only by redistribution.",
    },
  ];

  return (
    <section className="v1-ch05-insight">
      <h2 className="v1-ch05-insight-h2">Where can ∆M &gt; 0 come from?</h2>
      <p className="v1-ch05-insight-prose">
        Marx works through three candidate sources of new value in circulation
        and falsifies each. Run the probes below; every candidate that comes
        back as <em>redistribution</em> is ticked off as disproven.
      </p>
      <ul className="v1-ch05-candidate-list">
        {items.map((it) => (
          <li
            key={it.key}
            className={
              "v1-ch05-candidate" +
              (candidates[it.key] ? " v1-ch05-candidate--checked" : "")
            }
          >
            <span className="v1-ch05-candidate-check" aria-hidden="true">
              {candidates[it.key] ? "✓" : "·"}
            </span>
            <span className="v1-ch05-candidate-body">
              <span className="v1-ch05-candidate-label">{it.label}</span>
              <span className="v1-ch05-candidate-gloss">{it.gloss}</span>
            </span>
          </li>
        ))}
      </ul>
    </section>
  );
}

function CircuitProbePanel({ onUsuryProven }: { onUsuryProven: () => void }) {
  const fmt = usePounds();
  const [m, setM] = useState(10000);
  const [commodityId, setCommodityId] = useState("");
  const [mPrime, setMPrime] = useState(10500);
  const [result, setResult] = useState<CircuitProof | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.computeCircuit({
        m,
        commodity_id: commodityId || undefined,
        m_prime: mPrime,
      });
      setResult(r);
      if (r.origin === "redistribution" && !r.commodity_id) {
        onUsuryProven();
      }
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Circuit Probe</h2>
      <p className="description">
        Test whether M—C—M′ can produce surplus-value through circulation
        alone. Leave the commodity blank to probe M—M′ (usurer's capital).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>M advanced (pence)</span>
          <input
            type="number"
            value={m}
            onChange={(e) => setM(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Commodity ID (blank = M—M′)</span>
          <input
            value={commodityId}
            onChange={(e) => setCommodityId(e.target.value)}
          />
        </label>
        <label>
          <span>M′ returned (pence)</span>
          <input
            type="number"
            value={mPrime}
            onChange={(e) => setMPrime(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && <CircuitResult result={result} fmt={fmt} />}
    </section>
  );
}

function CircuitResult({
  result,
  fmt,
}: {
  result: CircuitProof;
  fmt: (p: number) => string;
}) {
  const delta = result.surplus_value;
  return (
    <div className="v1-ch05-circuit-result">
      <div className="v1-ch05-circuit-result-row">
        <span className="v1-ch05-circuit-formula">
          {fmt(result.m)}
          {result.commodity_id
            ? ` → C (${result.commodity_id.slice(0, 8)}…) → `
            : " → "}
          {fmt(result.m_prime)}
          {result.origin === "redistribution" && delta > 0 && (
            <span className="v1-ch05-circuit-prime">′</span>
          )}
        </span>
        <span
          className={
            "v1-ch05-origin-tag v1-ch05-origin-tag--" + result.origin
          }
        >
          {result.origin}
        </span>
      </div>
      <span className="v1-ch05-circuit-delta">
        ∆M = {delta >= 0 ? "+" : ""}
        {fmt(delta)}
      </span>
      <p className="v1-ch05-circuit-gloss">
        {result.origin === "redistribution"
          ? result.commodity_id
            ? "Value moved between buyer and seller but was not created. The circuit's apparent surplus is the counterparty's loss."
            : "Pure M—M′. Lending money for more money is redistribution between debtor and creditor — no new value."
          : "Exchange at value: M′ = M. No surplus arose from circulation."}
      </p>
    </div>
  );
}

function ExchangeSimulationPanel({
  onRedistribution,
}: {
  onRedistribution: (aBefore: number, bBefore: number) => void;
}) {
  const fmt = usePounds();
  const [aValue, setAValue] = useState(5000);
  const [bValue, setBValue] = useState(5000);
  const [result, setResult] = useState<ExchangeSimulation | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.simulateExchange({ a_value: aValue, b_value: bValue });
      setResult(r);
      if (r.origin === "redistribution") {
        onRedistribution(r.a_before, r.b_before);
      }
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Exchange Simulation</h2>
      <p className="description">
        Prove the invariant <code>a_before + b_before == a_after + b_after</code>.
        Bilateral exchange conserves total social value: whatever one party gains,
        the other loses.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>A holds (pence)</span>
          <input
            type="number"
            value={aValue}
            onChange={(e) => setAValue(Number(e.target.value))}
            min={0}
          />
        </label>
        <label>
          <span>B holds (pence)</span>
          <input
            type="number"
            value={bValue}
            onChange={(e) => setBValue(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Simulate</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && <ConservationDiagram result={result} fmt={fmt} />}
    </section>
  );
}

function ConservationDiagram({
  result,
  fmt,
}: {
  result: ExchangeSimulation;
  fmt: (p: number) => string;
}) {
  const totalBefore = result.a_before + result.b_before;
  const totalAfter = result.a_after + result.b_after;
  // Scale to whichever total is larger so both rows show the conservation
  // invariant: identical total bar length before vs. after.
  const scale = Math.max(totalBefore, totalAfter, 1);
  const pct = (n: number) => `${(n / scale) * 100}%`;

  const aDelta = result.a_after - result.a_before;
  const bDelta = result.b_after - result.b_before;
  const deltaClass = (d: number) =>
    d > 0
      ? "v1-ch05-delta-value--positive"
      : d < 0
      ? "v1-ch05-delta-value--negative"
      : "v1-ch05-delta-value--zero";
  const deltaSign = (d: number) => (d > 0 ? "+" : "");

  return (
    <>
      <div className="v1-ch05-conservation" role="img" aria-label="Value conservation diagram">
        <div className="v1-ch05-conservation-row">
          <span className="v1-ch05-conservation-label">Before</span>
          <div className="v1-ch05-conservation-track">
            <div
              className="v1-ch05-conservation-seg v1-ch05-conservation-seg--a"
              style={{ flexBasis: pct(result.a_before) }}
              title={`A: ${fmt(result.a_before)}`}
            >
              A
            </div>
            <div
              className="v1-ch05-conservation-seg v1-ch05-conservation-seg--b"
              style={{ flexBasis: pct(result.b_before) }}
              title={`B: ${fmt(result.b_before)}`}
            >
              B
            </div>
          </div>
          <span className="v1-ch05-conservation-total">{fmt(totalBefore)}</span>
        </div>
        <div className="v1-ch05-conservation-row">
          <span className="v1-ch05-conservation-label">After</span>
          <div className="v1-ch05-conservation-track">
            <div
              className="v1-ch05-conservation-seg v1-ch05-conservation-seg--a"
              style={{ flexBasis: pct(result.a_after) }}
              title={`A: ${fmt(result.a_after)}`}
            >
              A
            </div>
            <div
              className="v1-ch05-conservation-seg v1-ch05-conservation-seg--b"
              style={{ flexBasis: pct(result.b_after) }}
              title={`B: ${fmt(result.b_after)}`}
            >
              B
            </div>
          </div>
          <span className="v1-ch05-conservation-total">{fmt(totalAfter)}</span>
        </div>
      </div>
      <div className="v1-ch05-deltas">
        <div className="v1-ch05-delta">
          <span className="v1-ch05-delta-party">∆A</span>
          <span className={deltaClass(aDelta)}>
            {deltaSign(aDelta)}
            {fmt(aDelta)}
          </span>
        </div>
        <div className="v1-ch05-delta">
          <span className="v1-ch05-delta-party">∆B</span>
          <span className={deltaClass(bDelta)}>
            {deltaSign(bDelta)}
            {fmt(bDelta)}
          </span>
        </div>
      </div>
      <p className="v1-ch05-conservation-note">
        <span
          className={
            "v1-ch05-origin-tag v1-ch05-origin-tag--" + result.origin
          }
        >
          {result.origin}
        </span>{" "}
        {result.origin === "redistribution"
          ? "A's gain equals B's loss. Total social value is unchanged."
          : "Values are equal: exchange conserves value perfectly."}
      </p>
    </>
  );
}

function Coda() {
  return (
    <aside className="v1-ch05-coda">
      <p className="v1-ch05-coda-quote">
        “Capital cannot therefore originate in circulation, and it is equally
        impossible for it to originate apart from circulation. It must have
        its origin both in circulation and yet not in circulation.”
        <span className="v1-ch05-coda-cite">
          — Marx, Capital Vol. I, Ch. 5
        </span>
      </p>
    </aside>
  );
}
