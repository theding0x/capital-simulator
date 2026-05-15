import { useState } from "react";
import { api } from "../api";
import type { ComputePiecePriceResult, PieceWage } from "../types";
import { usePounds } from "../CurrencyContext";
import "./Ch21PieceWages.css";

export function Ch21PieceWages() {
  const fmt = usePounds();
  const [dailyWage, setDailyWage] = useState(144);
  const [dayValueProduct, setDayValueProduct] = useState(288);
  const [normalOutput, setNormalOutput] = useState(24);
  const [piecesProduced, setPiecesProduced] = useState(24);
  const [quality, setQuality] = useState<"accepted" | "rejected">("accepted");
  const [result, setResult] = useState<ComputePiecePriceResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [contractAgentID, setContractAgentID] = useState("");
  const [contractPrice, setContractPrice] = useState(6);
  const [contractOutput, setContractOutput] = useState(24);
  const [contract, setContract] = useState<PieceWage | null>(null);
  const [contractError, setContractError] = useState<string | null>(null);

  async function handleCompute(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const res = await api.computePiecePrice({
        daily_wage_farthings: dailyWage,
        day_value_product_farthings: dayValueProduct,
        normal_output: normalOutput,
        pieces_produced: piecesProduced,
        quality_outcome: quality,
      });
      setResult(res);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    }
  }

  async function handleCreateContract(e: React.FormEvent) {
    e.preventDefault();
    setContractError(null);
    try {
      const res = await api.createPieceWage(contractAgentID, {
        price_pence: contractPrice,
        normal_output: contractOutput,
      });
      setContract(res);
    } catch (err) {
      setContractError(err instanceof Error ? err.message : String(err));
    }
  }

  const livePiecePrice =
    normalOutput > 0 ? Math.floor(dailyWage / normalOutput) : 0;
  const livePieceValue =
    normalOutput > 0 ? Math.floor(dayValueProduct / normalOutput) : 0;

  return (
    <div className="ch21-piece-wage">
      <section className="ch21-section">
        <h2>Piece-Wage Calculator</h2>
        <p className="ch21-explainer">
          The piece-wage is the time-wage in disguise. Payment per piece sets
          the unit price so the average worker earns the daily wage by
          producing the socially normal output. The piece value always exceeds
          the piece price — the difference is the unpaid surplus labour the
          capitalist extracts.
        </p>

        <div className="ch21-live-hint">
          Live piece price:{" "}
          <strong>{livePiecePrice} farthings</strong> ({dailyWage}f ÷{" "}
          {normalOutput} pieces)
          {" · "}
          Piece value:{" "}
          <strong>{livePieceValue} farthings</strong> ({dayValueProduct}f ÷{" "}
          {normalOutput} pieces)
        </div>

        <form onSubmit={handleCompute} className="ch21-form">
          <label>
            Daily wage (farthings, ¼d.)
            <input
              type="number"
              min={1}
              value={dailyWage}
              onChange={(e) => setDailyWage(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Day value product (farthings)
            <input
              type="number"
              min={1}
              value={dayValueProduct}
              onChange={(e) => setDayValueProduct(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Normal output (pieces per day)
            <input
              type="number"
              min={1}
              value={normalOutput}
              onChange={(e) => setNormalOutput(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Pieces produced (session)
            <input
              type="number"
              min={0}
              value={piecesProduced}
              onChange={(e) => setPiecesProduced(Number(e.target.value))}
            />
          </label>
          <label>
            Quality outcome
            <select
              value={quality}
              onChange={(e) =>
                setQuality(e.target.value as "accepted" | "rejected")
              }
            >
              <option value="accepted">Accepted</option>
              <option value="rejected">Rejected</option>
            </select>
          </label>
          <button type="submit">Compute</button>
        </form>
        {error && <p className="ch21-error">{error}</p>}
      </section>

      {result && (
        <section className="ch21-section ch21-results">
          <h2>Results</h2>

          <div className="ch21-cards">
            <div className="ch21-card">
              <h3>Piece Price vs. Piece Value</h3>
              <p className="ch21-card-subtitle">
                The key illusion of the piece-wage form
              </p>
              <dl>
                <dt>Piece price (wage side)</dt>
                <dd>{result.piece_price} farthings</dd>
                <dt>Piece value (value side)</dt>
                <dd>{result.piece_value} farthings</dd>
                <dt>Surplus per piece</dt>
                <dd className="ch21-surplus">
                  {result.piece_value - result.piece_price} farthings
                </dd>
              </dl>
              <p className="ch21-note">
                The surplus is invisible to the worker: piece-wages appear to
                pay for product, not time, hiding that the working day still
                contains unpaid labour.
              </p>
            </div>

            <div className="ch21-card ch21-card-earnings">
              <h3>Session Earnings</h3>
              <p className="ch21-card-subtitle">
                {quality === "accepted"
                  ? `${piecesProduced} pieces accepted`
                  : `${piecesProduced} pieces rejected — earns nothing`}
              </p>
              <dl>
                <dt>Pieces produced</dt>
                <dd>{piecesProduced}</dd>
                <dt>Quality outcome</dt>
                <dd>{quality}</dd>
                <dt>Actual earnings</dt>
                <dd className="ch21-total">{result.actual_earnings} farthings</dd>
              </dl>
              <p className="ch21-note">
                Rejected pieces earn nothing — quality control is enforced
                directly by the wage form, not by a separate discipline system.
              </p>
            </div>
          </div>
        </section>
      )}
      <section className="ch21-section">
        <h2>Piece-Wage Contract</h2>
        <p className="ch21-explainer">
          Register a piece-wage contract for a worker. The implied daily wage
          shows what the worker earns when producing the normal output — the
          connection between piece-rate and time-wage made explicit.
        </p>
        <form onSubmit={handleCreateContract} className="ch21-form">
          <label>
            Agent ID
            <input
              type="text"
              value={contractAgentID}
              onChange={(e) => setContractAgentID(e.target.value)}
              placeholder="e.g. worker-001"
              required
            />
          </label>
          <label>
            Price per piece (pence)
            <input
              type="number"
              min={1}
              value={contractPrice}
              onChange={(e) => setContractPrice(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Normal output (pieces per day)
            <input
              type="number"
              min={1}
              value={contractOutput}
              onChange={(e) => setContractOutput(Number(e.target.value))}
              required
            />
          </label>
          <button type="submit">Create Contract</button>
        </form>
        {contractError && <p className="ch21-error">{contractError}</p>}
        {contract && (
          <dl className="ch21-contract-result">
            <dt>Piece price</dt>
            <dd>{fmt(contract.price_pence)}</dd>
            <dt>Normal output</dt>
            <dd>{contract.normal_output} pieces</dd>
            <dt>Implied daily wage</dt>
            <dd className="ch21-total">{fmt(contract.implied_daily_wage)}</dd>
          </dl>
        )}
      </section>
    </div>
  );
}
