import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import { usePounds } from "../CurrencyContext";
import type {
  HistoricalStage,
  CapitalOrigin,
  ColonialTransfer,
  NationalDebt,
  IndustrialCapitalGenesis,
} from "../types";

export function Ch31IndustrialCapitalist() {
  const fmt = usePounds();

  const [stages, setStages] = useState<HistoricalStage[]>([]);
  const [selectedStageId, setSelectedStageId] = useState("");
  const [genesis, setGenesis] = useState<IndustrialCapitalGenesis | null>(null);
  const [genesisError, setGenesisError] = useState("");

  // capital origin form
  const [coSource, setCoSource] = useState("Liverpool slave trade");
  const [coAmount, setCoAmount] = useState("15000");
  const [coPeriod, setCoPeriod] = useState("1730");
  const [coError, setCoError] = useState("");
  const [coCreating, setCoCreating] = useState(false);
  const [capitalOrigins, setCapitalOrigins] = useState<CapitalOrigin[]>([]);

  // colonial transfer form
  const [ctFrom, setCtFrom] = useState("Americas");
  const [ctTo, setCtTo] = useState("England");
  const [ctValue, setCtValue] = useState("96000");
  const [ctMethod, setCtMethod] = useState("plunder");
  const [ctError, setCtError] = useState("");
  const [ctCreating, setCtCreating] = useState(false);
  const [transfers, setTransfers] = useState<ColonialTransfer[]>([]);

  // national debt form
  const [ndAmount, setNdAmount] = useState("2400000");
  const [ndRate, setNdRate] = useState("800");
  const [ndCreditor, setNdCreditor] = useState("private-bankers");
  const [ndError, setNdError] = useState("");
  const [ndCreating, setNdCreating] = useState(false);
  const [debts, setDebts] = useState<NationalDebt[]>([]);

  useEffect(() => {
    api.listHistoricalStages().then((list) => {
      setStages(list);
      if (list.length > 0) {
        setSelectedStageId(list[0].id);
      }
    });
  }, []);

  useEffect(() => {
    if (!selectedStageId) return;
    loadGenesis(selectedStageId);
    setCapitalOrigins([]);
    setTransfers([]);
    setDebts([]);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedStageId]);

  async function loadGenesis(id: string) {
    setGenesisError("");
    try {
      const g = await api.getIndustrialCapitalGenesis(id);
      setGenesis(g);
    } catch (e) {
      setGenesisError(String(e));
    }
  }

  async function handleCreateCapitalOrigin(e: FormEvent) {
    e.preventDefault();
    if (!selectedStageId) return;
    setCoError("");
    setCoCreating(true);
    try {
      const co = await api.createCapitalOrigin(selectedStageId, {
        source: coSource,
        amount_pence: parseInt(coAmount, 10),
        period: coPeriod,
      });
      setCapitalOrigins((prev) => [...prev, co]);
      await loadGenesis(selectedStageId);
    } catch (e) {
      setCoError(String(e));
    } finally {
      setCoCreating(false);
    }
  }

  async function handleCreateColonialTransfer(e: FormEvent) {
    e.preventDefault();
    if (!selectedStageId) return;
    setCtError("");
    setCtCreating(true);
    try {
      const ct = await api.createColonialTransfer(selectedStageId, {
        from: ctFrom,
        to: ctTo,
        value_pence: parseInt(ctValue, 10),
        method: ctMethod,
      });
      setTransfers((prev) => [...prev, ct]);
      await loadGenesis(selectedStageId);
    } catch (e) {
      setCtError(String(e));
    } finally {
      setCtCreating(false);
    }
  }

  async function handleCreateNationalDebt(e: FormEvent) {
    e.preventDefault();
    if (!selectedStageId) return;
    setNdError("");
    setNdCreating(true);
    try {
      const nd = await api.createNationalDebt(selectedStageId, {
        amount_pence: parseInt(ndAmount, 10),
        interest_rate_bps: parseInt(ndRate, 10),
        creditor_class: ndCreditor,
      });
      setDebts((prev) => [...prev, nd]);
      await loadGenesis(selectedStageId);
    } catch (e) {
      setNdError(String(e));
    } finally {
      setNdCreating(false);
    }
  }

  return (
    <div className="chapter-panel">
      <section className="panel-section">
        <h2>Industrial Capital Genesis</h2>
        <p className="muted small">
          Select a historical stage to see the aggregate capital formed through
          slave trade profits, colonial plunder, and similar primitive
          accumulation. National debts and protection systems are structural
          levers and are listed separately but excluded from the capital total.
        </p>
        <label>
          Historical Stage
          <select
            value={selectedStageId}
            onChange={(e) => setSelectedStageId(e.target.value)}
          >
            {stages.map((s) => (
              <option key={s.id} value={s.id}>
                {s.name}
              </option>
            ))}
          </select>
        </label>
        {genesisError && <p className="error-msg">{genesisError}</p>}
        {genesis && (
          <dl className="result-list">
            <dt>Total Capital Formed</dt>
            <dd>{fmt(genesis.total_capital_formed_pence)}</dd>
            <dt>Capital Origins</dt>
            <dd>{genesis.origins.length}</dd>
            <dt>Colonial Transfers</dt>
            <dd>{genesis.colonial_transfers.length}</dd>
            <dt>National Debts</dt>
            <dd>{genesis.national_debts.length}</dd>
            <dt>Protection Systems</dt>
            <dd>{genesis.protection_systems.length}</dd>
          </dl>
        )}
        {genesis && genesis.protection_systems.length > 0 && (
          <div style={{ marginTop: "0.75rem" }}>
            <strong>Protection Systems</strong>
            <ul className="result-list" style={{ marginTop: "0.25rem" }}>
              {genesis.protection_systems.map((ps) => (
                <li key={ps.id}>
                  {ps.beneficiary} — {(ps.tariff_rate_bps / 100).toFixed(1)}% tariff
                  ({ps.period_start}&#8211;{ps.period_end})
                </li>
              ))}
            </ul>
          </div>
        )}
      </section>

      <section className="panel-section">
        <h2>Register Capital Origin</h2>
        <p className="muted small">
          Record a source of capital formation: slave trade profits, state
          monopoly revenues, prize money from colonial wars, etc. Marx&apos;s
          fixture: Liverpool slave trade 1730&#8211;1792.
        </p>
        <form onSubmit={handleCreateCapitalOrigin} className="form-grid">
          <label>
            Source
            <input
              value={coSource}
              onChange={(e) => setCoSource(e.target.value)}
              required
            />
          </label>
          <label>
            Amount (halfpennies)
            <input
              type="number"
              min="1"
              value={coAmount}
              onChange={(e) => setCoAmount(e.target.value)}
              required
            />
          </label>
          <label>
            Period
            <input
              value={coPeriod}
              onChange={(e) => setCoPeriod(e.target.value)}
              placeholder="e.g. 1730"
              required
            />
          </label>
          {coError && <p className="error-msg">{coError}</p>}
          <button
            type="submit"
            className="btn-primary"
            disabled={coCreating || !selectedStageId}
          >
            {coCreating ? "Registering…" : "Register Origin"}
          </button>
        </form>
        {capitalOrigins.length > 0 && (
          <ul className="result-list" style={{ marginTop: "0.75rem" }}>
            {capitalOrigins.map((co) => (
              <li key={co.id}>
                <strong>{co.source}</strong> ({co.period}) &#8212; {fmt(co.amount_pence)}
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="panel-section">
        <h2>Record Colonial Transfer</h2>
        <p className="muted small">
          Record wealth extracted from a colony and transferred to the
          metropolis. Marx&apos;s fixture: colonial plunder from the Americas
          funding English industrial capital.
        </p>
        <form onSubmit={handleCreateColonialTransfer} className="form-grid">
          <label>
            From (colony/region)
            <input
              value={ctFrom}
              onChange={(e) => setCtFrom(e.target.value)}
              required
            />
          </label>
          <label>
            To (metropolis)
            <input
              value={ctTo}
              onChange={(e) => setCtTo(e.target.value)}
              required
            />
          </label>
          <label>
            Value (halfpennies)
            <input
              type="number"
              min="1"
              value={ctValue}
              onChange={(e) => setCtValue(e.target.value)}
              required
            />
          </label>
          <label>
            Method
            <input
              value={ctMethod}
              onChange={(e) => setCtMethod(e.target.value)}
              placeholder="e.g. plunder, tribute, monopoly"
              required
            />
          </label>
          {ctError && <p className="error-msg">{ctError}</p>}
          <button
            type="submit"
            className="btn-primary"
            disabled={ctCreating || !selectedStageId}
          >
            {ctCreating ? "Recording…" : "Record Transfer"}
          </button>
        </form>
        {transfers.length > 0 && (
          <ul className="result-list" style={{ marginTop: "0.75rem" }}>
            {transfers.map((ct) => (
              <li key={ct.id}>
                {ct.from} &#8594; {ct.to} via <em>{ct.method}</em> &#8212; {fmt(ct.value_pence)}
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="panel-section">
        <h2>Record National Debt</h2>
        <p className="muted small">
          National debt as a weapon of primitive accumulation: the state borrows
          from private bankers and services the debt by taxing workers.
          Marx&apos;s fixture: Bank of England founding 1694 &#8212; &#163;1.2 million at 8%.
        </p>
        <form onSubmit={handleCreateNationalDebt} className="form-grid">
          <label>
            Amount (halfpennies)
            <input
              type="number"
              min="1"
              value={ndAmount}
              onChange={(e) => setNdAmount(e.target.value)}
              required
            />
          </label>
          <label>
            Interest Rate (basis points)
            <input
              type="number"
              min="1"
              value={ndRate}
              onChange={(e) => setNdRate(e.target.value)}
              placeholder="e.g. 800 = 8%"
              required
            />
          </label>
          <label>
            Creditor Class
            <input
              value={ndCreditor}
              onChange={(e) => setNdCreditor(e.target.value)}
              placeholder="e.g. private-bankers"
              required
            />
          </label>
          {ndError && <p className="error-msg">{ndError}</p>}
          <button
            type="submit"
            className="btn-primary"
            disabled={ndCreating || !selectedStageId}
          >
            {ndCreating ? "Recording…" : "Record National Debt"}
          </button>
        </form>
        {debts.length > 0 && (
          <ul className="result-list" style={{ marginTop: "0.75rem" }}>
            {debts.map((nd) => (
              <li key={nd.id}>
                {fmt(nd.amount_pence)} at {(nd.interest_rate_bps / 100).toFixed(1)}% &#8212;{" "}
                {nd.creditor_class}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
