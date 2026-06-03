import { useEffect, useState } from "react";
import { api } from "../../api";
import type {
  ComputePiecePriceResult,
  PiecePricePoint,
  PieceWage,
  SubContract,
} from "../../types";
import { usePounds } from "../../CurrencyContext";
import "./Ch21PieceWages.css";

function PieceWageInsight() {
  return (
    <section className="v1-ch21-insight">
      <h2 className="v1-ch21-insight-h2">
        Piece-price vs. piece-value — the wage form's central illusion
      </h2>
      <p className="v1-ch21-insight-prose">
        Two prices attach to every piece in this chapter. The <em>piece price</em>
        is the wage split — daily pay divided by normal output. The
        <em> piece value</em> is the new value the piece adds — the day's value
        product divided by normal output. The capitalist pays the first and
        pockets the difference. The piece-wage form hides this gap by appearing
        to pay for the product itself, as if no time had been bought at all.
      </p>
    </section>
  );
}

// The seed Spinning Mill worker (agent migration 00019/00031). Its piece-wage
// baseline is 6 farthings at 24 pieces (144f implied daily wage).
const SPINNING_MILL_AGENT_ID = "5eed000000000000001900ff";

function PiecePriceSeries({ points }: { points: PiecePricePoint[] }) {
  if (points.length === 0) {
    return (
      <p className="v1-ch21-series-empty">
        No productivity ticks recorded yet — apply a productivity rise to begin
        the series.
      </p>
    );
  }
  const W = 320;
  const H = 140;
  const pad = 24;
  const maxPrice = Math.max(...points.map((p) => p.price_pence), 1);
  const x = (i: number) =>
    points.length === 1 ? W / 2 : pad + (i / (points.length - 1)) * (W - 2 * pad);
  const y = (price: number) => pad + (1 - price / maxPrice) * (H - 2 * pad);
  return (
    <div className="v1-ch21-series">
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="v1-ch21-series-svg"
        role="img"
        aria-label="Price per piece falling as productivity rises"
      >
        <polyline
          className="v1-ch21-series-line"
          fill="none"
          points={points.map((p, i) => `${x(i)},${y(p.price_pence)}`).join(" ")}
        />
        {points.map((p, i) => (
          <circle
            key={p.id}
            cx={x(i)}
            cy={y(p.price_pence)}
            r={3.5}
            className="v1-ch21-series-dot"
          >
            <title>{`${p.productivity_factor}× productivity → ${p.price_pence}f per piece`}</title>
          </circle>
        ))}
      </svg>
      <table className="v1-ch21-series-table">
        <thead>
          <tr>
            <th>Productivity</th>
            <th>Price / piece</th>
            <th>Normal output</th>
            <th>Daily wage</th>
          </tr>
        </thead>
        <tbody>
          {points.map((p) => (
            <tr key={p.id}>
              <td>{p.productivity_factor}×</td>
              <td>{p.price_pence}f</td>
              <td>{p.normal_output}</td>
              <td>{p.price_pence * p.normal_output}f</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function DynamicPiecePrice() {
  const [points, setPoints] = useState<PiecePricePoint[]>([]);
  const [factor, setFactor] = useState(2);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    let active = true;
    api
      .listPiecePricePoints(SPINNING_MILL_AGENT_ID)
      .then((pts) => {
        if (active) setPoints(pts);
      })
      .catch((err) => {
        if (active) setError(err instanceof Error ? err.message : String(err));
      });
    return () => {
      active = false;
    };
  }, []);

  async function handleReprice(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await api.repricePieceWage(SPINNING_MILL_AGENT_ID, {
        productivity_factor: factor,
      });
      setPoints(res.points);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  const latest = points.length > 0 ? points[points.length - 1] : null;

  return (
    <section className="ch21-section v1-ch21-dynamic">
      <h2>Dynamic piece-price — productivity rises, the price falls</h2>
      <p className="ch21-explainer">
        “The price of labour … falls in the same proportion as the
        productiveness of labour rises.” The daily wage is unchanged; as the
        socially-normal output grows, the price per piece falls in inverse
        proportion. The simulation-engine tick scheduler drives this from the
        Ch. 15 factory productivity factor — here you can apply a rise by hand.
        At 2× productivity the seed Spinning Mill price halves, 6f → 3f.
      </p>
      <PiecePriceSeries points={points} />
      {latest && (
        <p className="v1-ch21-series-current">
          Current: <strong>{latest.price_pence}f</strong> per piece at{" "}
          {latest.productivity_factor}× productivity — {latest.normal_output}{" "}
          pieces/day, daily wage{" "}
          <strong>{latest.price_pence * latest.normal_output}f</strong> held
          constant.
        </p>
      )}
      <form onSubmit={handleReprice} className="ch21-form v1-ch21-dynamic-form">
        <label>
          Productivity factor (× baseline)
          <input
            type="number"
            min={0.1}
            step={0.1}
            value={factor}
            onChange={(e) => setFactor(Number(e.target.value))}
            required
          />
        </label>
        <button type="submit" disabled={loading}>
          {loading ? "Applying…" : "Apply productivity rise"}
        </button>
      </form>
      {error && <p className="ch21-error">{error}</p>}
    </section>
  );
}

function PiecePriceValueBars({
  piecePrice,
  pieceValue,
}: {
  piecePrice: number;
  pieceValue: number;
}) {
  const max = Math.max(piecePrice, pieceValue, 1);
  const pricePct = (piecePrice / max) * 100;
  const valuePct = (pieceValue / max) * 100;
  return (
    <div className="v1-ch21-bars" aria-label="Piece price versus piece value">
      <div className="v1-ch21-bar-row">
        <span className="v1-ch21-bar-label">Piece price</span>
        <div className="v1-ch21-bar-track">
          <div
            className="v1-ch21-bar-fill v1-ch21-bar-price"
            style={{ width: `${pricePct}%` }}
          />
        </div>
        <span className="v1-ch21-bar-value">{piecePrice}f</span>
      </div>
      <div className="v1-ch21-bar-row">
        <span className="v1-ch21-bar-label">Piece value</span>
        <div className="v1-ch21-bar-track">
          <div
            className="v1-ch21-bar-fill v1-ch21-bar-value-fill"
            style={{ width: `${valuePct}%` }}
          />
        </div>
        <span className="v1-ch21-bar-value">{pieceValue}f</span>
      </div>
      <p className="v1-ch21-bars-note">
        Gap = unpaid surplus per piece ({pieceValue - piecePrice}f).
      </p>
    </div>
  );
}

function PieceWageCoda() {
  return (
    <aside className="v1-ch21-coda">
      <p className="v1-ch21-coda-quote">
        “The form of piece-wages is just as irrational as that of time-wages.
        ... It is in fact a converted form of the time-wage, just as the
        time-wage is a converted form of the value or price of labour-power.”
        <span className="v1-ch21-coda-cite">
          — Marx, Capital Vol. I, Ch. 21
        </span>
      </p>
    </aside>
  );
}

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

  const [subHeadID, setSubHeadID] = useState("");
  const [subAssistantsCsv, setSubAssistantsCsv] = useState("");
  const [subPieceRate, setSubPieceRate] = useState(8);
  const [subAssistantRate, setSubAssistantRate] = useState(5);
  const [subContract, setSubContract] = useState<SubContract | null>(null);
  const [subError, setSubError] = useState<string | null>(null);

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

  async function handleCreateSubContract(e: React.FormEvent) {
    e.preventDefault();
    setSubError(null);
    try {
      const assistantIds = subAssistantsCsv
        .split(",")
        .map((s) => s.trim())
        .filter(Boolean);
      const res = await api.createSubContract({
        head_labourer_id: subHeadID,
        assistant_ids: assistantIds,
        piece_rate_pence: subPieceRate,
        assistant_rate_pence: subAssistantRate,
      });
      setSubContract(res);
    } catch (err) {
      setSubError(err instanceof Error ? err.message : String(err));
    }
  }

  const livePiecePrice =
    normalOutput > 0 ? Math.floor(dailyWage / normalOutput) : 0;
  const livePieceValue =
    normalOutput > 0 ? Math.floor(dayValueProduct / normalOutput) : 0;

  return (
    <div className="ch21-piece-wage">
      <PieceWageInsight />
      <DynamicPiecePrice />
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
              <PiecePriceValueBars
                piecePrice={result.piece_price}
                pieceValue={result.piece_value}
              />
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
      <section className="ch21-section">
        <h2>Sub-Contract — head labourer and assistants</h2>
        <p className="ch21-explainer">
          §4 of the chapter: a head labourer takes the piece-rate from the
          capitalist, then employs assistants at a lower internal rate. The
          spread — head rate minus assistant rate — is exploitation
          intermediated by another worker, masking the wage relation behind
          one that looks like skilled-versus-helper.
        </p>
        <form onSubmit={handleCreateSubContract} className="ch21-form">
          <label>
            Head labourer ID
            <input
              type="text"
              value={subHeadID}
              onChange={(e) => setSubHeadID(e.target.value)}
              placeholder="e.g. head-001"
              required
            />
          </label>
          <label>
            Assistant IDs (comma-separated)
            <input
              type="text"
              value={subAssistantsCsv}
              onChange={(e) => setSubAssistantsCsv(e.target.value)}
              placeholder="assistant-001, assistant-002"
              required
            />
          </label>
          <label>
            Piece rate to head (pence)
            <input
              type="number"
              min={1}
              value={subPieceRate}
              onChange={(e) => setSubPieceRate(Number(e.target.value))}
              required
            />
          </label>
          <label>
            Assistant rate (pence)
            <input
              type="number"
              min={0}
              value={subAssistantRate}
              onChange={(e) => setSubAssistantRate(Number(e.target.value))}
              required
            />
          </label>
          <button type="submit">Register sub-contract</button>
        </form>
        {subError && <p className="ch21-error">{subError}</p>}
        {subContract && (
          <dl className="ch21-contract-result v1-ch21-sub-result">
            <dt>Head rate</dt>
            <dd>{fmt(subContract.piece_rate_pence)}</dd>
            <dt>Assistant rate</dt>
            <dd>{fmt(subContract.assistant_rate_pence)}</dd>
            <dt>Spread (per piece)</dt>
            <dd className="ch21-total">{fmt(subContract.spread)}</dd>
            <dt>Assistants</dt>
            <dd>{subContract.assistant_ids.length}</dd>
          </dl>
        )}
      </section>
      <PieceWageCoda />
    </div>
  );
}
