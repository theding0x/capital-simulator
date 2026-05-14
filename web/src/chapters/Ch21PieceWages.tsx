import { useState } from "react";
import { api } from "../api";
import type { ComputePiecePriceResult } from "../types";
import "./Ch21PieceWages.css";

export function Ch21PieceWages() {
  const [dailyWage, setDailyWage] = useState(144);
  const [dayValueProduct, setDayValueProduct] = useState(288);
  const [normalOutput, setNormalOutput] = useState(24);
  const [piecesProduced, setPiecesProduced] = useState(24);
  const [quality, setQuality] = useState<"accepted" | "rejected">("accepted");
  const [result, setResult] = useState<ComputePiecePriceResult | null>(null);
  const [error, setError] = useState<string | null>(null);

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
    </div>
  );
}
