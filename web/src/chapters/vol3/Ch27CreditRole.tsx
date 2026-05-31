import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { StockCompany, CooperativeFactory } from "../../types";
import "./Ch27CreditRole.css";

// Marx's three roles of credit (Vol. III Ch. 27).
const CREDIT_ROLES = [
  {
    role: 1,
    label: "I — Equalisation of the Profit Rate",
    description:
      "Credit enables capital to flow rapidly across spheres of production without " +
      "physical cash settlement, accelerating equalisation of the rate of profit " +
      "to the general average.",
  },
  {
    role: 2,
    label: "II — Reduction of Circulation Costs",
    description:
      "Credit substitutes bills of exchange, bank-notes, and deposit transfers for " +
      "metallic coin, reducing the quantity of money required in circulation. " +
      "Scotland example: £27m of deposits circulate on only £3m of currency — " +
      "saving £24m of metallic money.",
  },
  {
    role: 3,
    label: "III — Joint-Stock Companies and Co-operatives",
    description:
      "Credit gives rise to joint-stock companies — the abolition of private capital " +
      "within capitalism. The functioning capitalist becomes a salaried manager; the " +
      "owner a pure money-capitalist drawing interest. Co-operative factories are the " +
      "positive, worker-owned sprout of the new mode of production.",
  },
];

// Seed exemplars used in the UI.
const STOCK_EXEMPLARS = [
  {
    label: "United Alkali Trust",
    fixed: 500_000_000,
    floating: 100_000_000,
    total: 600_000_000,
    note: "£5m fixed + £1m floating = £6m total (pence). Manager earns wages, not profit.",
  },
];

const CIRCULATION_EXEMPLAR = {
  metallic: 2_700_000_000,
  actual: 300_000_000,
  saving: 2_400_000_000,
  note: "Scottish banking: £27m deposits circulate on £3m currency — £24m saved (pence).",
};

function pToGBP(pence: number): string {
  return `£${(pence / 100).toLocaleString("en-GB", { maximumFractionDigits: 0 })}`;
}

export function Ch27CreditRole() {
  const [stockCompanies, setStockCompanies] = useState<StockCompany[]>([]);
  const [cooperatives, setCooperatives] = useState<CooperativeFactory[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Stock company create form.
  const [scName, setScName] = useState("United Alkali Trust");
  const [scFixed, setScFixed] = useState(500_000_000);
  const [scFloating, setScFloating] = useState(100_000_000);
  const [scTotal, setScTotal] = useState(600_000_000);
  const [scMgrName, setScMgrName] = useState("Salaried Director");
  const [scMgrTitle, setScMgrTitle] = useState("Managing Director");
  const [scMgrWage, setScMgrWage] = useState(480);
  const [scDesc, setScDesc] = useState("");
  const [scError, setScError] = useState<string | null>(null);

  // Cooperative create form.
  const [cfName, setCfName] = useState("Rochdale Pioneers");
  const [cfWorkers, setCfWorkers] = useState(28);
  const [cfCapital, setCfCapital] = useState(2800);
  const [cfDesc, setCfDesc] = useState("First modern workers cooperative (1844)");
  const [cfError, setCfError] = useState<string | null>(null);

  const scTotalValid = scFixed + scFloating === scTotal;

  async function refresh() {
    try {
      const [scs, cfs] = await Promise.all([
        financeApi.listStockCompanies(),
        financeApi.listCooperativeFactories(),
      ]);
      setStockCompanies(scs);
      setCooperatives(cfs);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  async function handleCreateStockCompany() {
    setScError(null);
    try {
      await financeApi.createStockCompany({
        name: scName,
        fixed_capital: scFixed,
        floating_capital: scFloating,
        total_capital: scTotal,
        manager_name: scMgrName,
        manager_title: scMgrTitle,
        manager_wage_minutes: scMgrWage,
        description: scDesc,
      });
      await refresh();
    } catch (e) {
      setScError(String(e));
    }
  }

  async function handleCreateCooperative() {
    setCfError(null);
    try {
      await financeApi.createCooperativeFactory({
        name: cfName,
        worker_owned: true,
        worker_count: cfWorkers,
        capital_advanced: cfCapital,
        description: cfDesc,
      });
      await refresh();
    } catch (e) {
      setCfError(String(e));
    }
  }

  return (
    <div className="v3-ch27">
      {/* Introduction */}
      <section className="v3-ch27-intro">
        <p>
          Credit in its developed form plays three historically decisive roles in capitalist
          production. Each role accelerates the circuit of capital, but also intensifies its
          contradictions.
        </p>
      </section>

      {/* Three roles */}
      <section className="v3-ch27-section" aria-label="Three roles of credit">
        <h4>The three roles of credit</h4>
        <div className="v3-ch27-roles">
          {CREDIT_ROLES.map((r) => (
            <div className="v3-ch27-role-card" key={r.role}>
              <div className="v3-ch27-role-label">{r.label}</div>
              <p className="v3-ch27-role-desc">{r.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Circulation saving exemplar */}
      <section className="v3-ch27-section" aria-label="Circulation cost saving exemplar">
        <h4>Role II exemplar — Scottish banking (Marx, Ch. 27)</h4>
        <div className="v3-ch27-saving-card">
          <div className="v3-ch27-saving-row">
            <span className="v3-ch27-saving-label">Metallic equivalent</span>
            <span className="v3-ch27-saving-val">{pToGBP(CIRCULATION_EXEMPLAR.metallic)}</span>
          </div>
          <div className="v3-ch27-saving-row">
            <span className="v3-ch27-saving-label">Actual currency in use</span>
            <span className="v3-ch27-saving-val">{pToGBP(CIRCULATION_EXEMPLAR.actual)}</span>
          </div>
          <div className="v3-ch27-saving-row v3-ch27-saving-row--total">
            <span className="v3-ch27-saving-label">Circulation cost saving</span>
            <span className="v3-ch27-saving-val v3-ch27-saving-ok">
              {pToGBP(CIRCULATION_EXEMPLAR.saving)}
            </span>
          </div>
          <p className="v3-ch27-saving-note">{CIRCULATION_EXEMPLAR.note}</p>
        </div>
      </section>

      {/* Stock company exemplars */}
      <section className="v3-ch27-section" aria-label="Stock company exemplars">
        <h4>Role III exemplars — joint-stock companies (Marx, Ch. 27)</h4>
        {STOCK_EXEMPLARS.map((ex) => (
          <div className="v3-ch27-exemplar-card" key={ex.label}>
            <div className="v3-ch27-exemplar-name">{ex.label}</div>
            <div className="v3-ch27-exemplar-meta">
              <span>Fixed: {pToGBP(ex.fixed)}</span>
              <span>Floating: {pToGBP(ex.floating)}</span>
              <span className="v3-ch27-total">Total: {pToGBP(ex.total)}</span>
            </div>
            <p className="v3-ch27-exemplar-note">{ex.note}</p>
          </div>
        ))}
      </section>

      {/* Create stock company */}
      <section className="v3-ch27-section" aria-label="Create a stock company">
        <h4>Record a joint-stock company</h4>
        <p className="v3-ch27-subintro">
          Total capital must equal fixed + floating capital (invariant enforced server-side).
        </p>
        <div className="v3-ch27-form">
          <label>
            Name
            <input value={scName} onChange={(e) => setScName(e.target.value)} />
          </label>
          <label>
            Fixed capital (pence)
            <input
              type="number"
              min={0}
              value={scFixed}
              onChange={(e) => {
                setScFixed(Number(e.target.value));
                setScTotal(Number(e.target.value) + scFloating);
              }}
            />
          </label>
          <label>
            Floating capital (pence)
            <input
              type="number"
              min={0}
              value={scFloating}
              onChange={(e) => {
                setScFloating(Number(e.target.value));
                setScTotal(scFixed + Number(e.target.value));
              }}
            />
          </label>
          <label>
            Total capital (pence)
            <input
              type="number"
              min={0}
              value={scTotal}
              onChange={(e) => setScTotal(Number(e.target.value))}
            />
            {!scTotalValid && (
              <span className="v3-ch27-warn">
                Must equal fixed ({scFixed}) + floating ({scFloating}) = {scFixed + scFloating}
              </span>
            )}
          </label>
          <label>
            Manager name
            <input value={scMgrName} onChange={(e) => setScMgrName(e.target.value)} />
          </label>
          <label>
            Manager title
            <input value={scMgrTitle} onChange={(e) => setScMgrTitle(e.target.value)} />
          </label>
          <label>
            Manager wage (minutes, &gt; 0)
            <input
              type="number"
              min={1}
              value={scMgrWage}
              onChange={(e) => setScMgrWage(Number(e.target.value))}
            />
            <span className="v3-ch27-hint">Wages of superintendence — NOT profit of enterprise.</span>
          </label>
          <label>
            Description
            <input value={scDesc} onChange={(e) => setScDesc(e.target.value)} />
          </label>
        </div>
        <button
          className="v3-ch27-btn"
          onClick={handleCreateStockCompany}
          disabled={!scTotalValid || scMgrWage <= 0}
        >
          Save stock company
        </button>
        {scError && <p className="v3-ch27-error">{scError}</p>}
      </section>

      {/* Stored stock companies */}
      <section className="v3-ch27-section" aria-label="Stored stock companies">
        <h4>Stored stock companies</h4>
        {listError && <p className="v3-ch27-error">{listError}</p>}
        <div className="v3-ch27-cards">
          {stockCompanies.length === 0 && (
            <p className="v3-ch27-subintro">No stock companies stored yet.</p>
          )}
          {stockCompanies.map((sc) => (
            <div className="v3-ch27-card" key={sc.id}>
              <div className="v3-ch27-card-name">{sc.name}</div>
              <div className="v3-ch27-card-meta">
                <span>Fixed: {pToGBP(sc.fixed_capital)}</span>
                <span>Floating: {pToGBP(sc.floating_capital)}</span>
                <span className="v3-ch27-total">Total: {pToGBP(sc.total_capital)}</span>
              </div>
              <div className="v3-ch27-card-mgr">
                Manager: {sc.manager.name} ({sc.manager.title}) &mdash;{" "}
                {sc.manager.wage_minutes} min wages
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Create cooperative factory */}
      <section className="v3-ch27-section" aria-label="Create a cooperative factory">
        <h4>Record a co-operative factory</h4>
        <p className="v3-ch27-subintro">
          Co-operative factories are worker-owned — the positive transitional form beyond
          capital. worker_owned is always true (enforced server-side).
        </p>
        <div className="v3-ch27-form">
          <label>
            Name
            <input value={cfName} onChange={(e) => setCfName(e.target.value)} />
          </label>
          <label>
            Number of worker-owners
            <input
              type="number"
              min={1}
              value={cfWorkers}
              onChange={(e) => setCfWorkers(Number(e.target.value))}
            />
          </label>
          <label>
            Capital advanced (pence)
            <input
              type="number"
              min={0}
              value={cfCapital}
              onChange={(e) => setCfCapital(Number(e.target.value))}
            />
          </label>
          <label>
            Description
            <input value={cfDesc} onChange={(e) => setCfDesc(e.target.value)} />
          </label>
        </div>
        <button className="v3-ch27-btn" onClick={handleCreateCooperative}>
          Save co-operative factory
        </button>
        {cfError && <p className="v3-ch27-error">{cfError}</p>}
      </section>

      {/* Stored cooperatives */}
      <section className="v3-ch27-section" aria-label="Stored co-operative factories">
        <h4>Stored co-operative factories</h4>
        <div className="v3-ch27-cards">
          {cooperatives.length === 0 && (
            <p className="v3-ch27-subintro">No co-operative factories stored yet.</p>
          )}
          {cooperatives.map((cf) => (
            <div className="v3-ch27-card v3-ch27-card--coop" key={cf.id}>
              <div className="v3-ch27-card-name">{cf.name}</div>
              <div className="v3-ch27-card-meta">
                <span className="v3-ch27-worker-tag">worker-owned</span>
                <span>{cf.worker_count} workers</span>
                <span>Capital: {pToGBP(cf.capital_advanced)}</span>
              </div>
              {cf.description && (
                <div className="v3-ch27-card-desc">{cf.description}</div>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
