import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  RentHistoricalForm,
  HistoricalRentTransition,
  SmallPeasantProduction,
  RentFormStage,
} from "../../types";
import "./Ch47Genesis.css";

const STAGE_LABELS: Record<RentFormStage, string> = {
  1: "1 — Labour-rent",
  2: "2 — Product-rent",
  3: "3 — Money-rent",
  4: "4 — Capitalist ground-rent",
};

export function Ch47Genesis() {
  const [forms, setForms] = useState<RentHistoricalForm[]>([]);
  const [transitions, setTransitions] = useState<HistoricalRentTransition[]>([]);
  const [peasants, setPeasants] = useState<SmallPeasantProduction[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Rent historical form fields
  const [fStage, setFStage] = useState<RentFormStage>(1);
  const [fName, setFName] = useState("");
  const [fCoercion, setFCoercion] = useState("");
  const [fSurplus, setFSurplus] = useState("");
  const [fRegion, setFRegion] = useState("");
  const [fPeriod, setFPeriod] = useState("");
  const [fError, setFError] = useState<string | null>(null);
  const [lastF, setLastF] = useState<RentHistoricalForm | null>(null);

  // Historical rent transition fields
  const [tFromStage, setTFromStage] = useState<RentFormStage>(1);
  const [tToStage, setTToStage] = useState<RentFormStage>(2);
  const [tPreconditions, setTPreconditions] = useState("");
  const [tDrivingForce, setTDrivingForce] = useState("");
  const [tError, setTError] = useState<string | null>(null);
  const [lastT, setLastT] = useState<HistoricalRentTransition | null>(null);

  // Small peasant production fields
  const [pParcelId, setPParcelId] = useState("5eed000000000000003700");
  const [pRent, setPRent] = useState(2);
  const [pProfit, setPProfit] = useState(1);
  const [pWages, setPWages] = useState(3);
  const [pConflated, setPConflated] = useState(true);
  const [pError, setPError] = useState<string | null>(null);
  const [lastP, setLastP] = useState<SmallPeasantProduction | null>(null);

  async function refreshLists() {
    try {
      const [fs, ts, ps] = await Promise.all([
        financeApi.listRentHistoricalForms(),
        financeApi.listHistoricalRentTransitions(),
        financeApi.listSmallPeasantProductions(),
      ]);
      setForms(fs);
      setTransitions(ts);
      setPeasants(ps);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateForm() {
    setFError(null);
    setLastF(null);
    try {
      const result = await financeApi.createRentHistoricalForm({
        stage: fStage,
        name: fName,
        coercion_kind: fCoercion,
        surplus_form: fSurplus,
        example_region: fRegion,
        example_period: fPeriod,
      });
      setLastF(result);
      await refreshLists();
    } catch (e) {
      setFError(String(e));
    }
  }

  async function handleCreateTransition() {
    setTError(null);
    setLastT(null);
    try {
      const result = await financeApi.createHistoricalRentTransition({
        from_stage: tFromStage,
        to_stage: tToStage,
        preconditions: tPreconditions,
        driving_force: tDrivingForce,
      });
      setLastT(result);
      await refreshLists();
    } catch (e) {
      setTError(String(e));
    }
  }

  async function handleCreatePeasant() {
    setPError(null);
    setLastP(null);
    try {
      const result = await financeApi.createSmallPeasantProduction({
        parcel_id: pParcelId,
        rent_component_bp: pRent,
        profit_component_bp: pProfit,
        wages_component_bp: pWages,
        is_conflated: pConflated,
      });
      setLastP(result);
      await refreshLists();
    } catch (e) {
      setPError(String(e));
    }
  }

  return (
    <div className="v3-ch47">
      <section className="v3-ch47-intro">
        <p>
          Chapter&nbsp;47 traces ground-rent back through its pre-capitalist history. Rent does not
          begin as a monetary category — it passes through three successive forms before capitalist
          ground-rent proper emerges.
        </p>
        <p>
          <strong>Labour-rent</strong> is the most transparent: the direct producer works part of
          the week on the lord's land and part on their own. Surplus-labour is performed directly
          for the landlord, visible to all — no mystification is possible. Its basis is
          extra-economic coercion: the serf is legally bound to the land and the person of the lord.
        </p>
        <p>
          <strong>Product-rent</strong> replaces the direct labour obligation with a share of the
          product itself. The producer surrenders a portion of the harvest. The personal bond begins
          to loosen; the direct producer has a stronger incentive to improve cultivation, since they
          retain any surplus above the quota.
        </p>
        <p>
          <strong>Money-rent</strong> is the crucial transitional form. The producer now sells their
          product on the market and hands over a money sum to the landlord. This dissolves the
          personal bond: the producer relates to the landlord only through money. It creates the
          conditions for the producer to become either an independent peasant proprietor or a
          landless wage-labourer — in either case the capitalist mode of production can emerge.
        </p>
        <p>
          Fully developed <strong>capitalist ground-rent</strong> requires the three-class
          separation: landowner, agricultural capitalist (who rents the land), and wage-labourer
          (who works it). The rent is now a deduction from surplus-value rather than surplus-labour
          or surplus-product.
        </p>
        <p>
          The <strong>small peasant proprietor</strong> represents a transitional figure in whom all
          three class functions are conflated. They are simultaneously labourer (their own wages),
          capitalist (their own profit), and landowner (their own rent). The three revenues — wages,
          profit, rent — cannot be analytically separated; they appear as one undifferentiated income.
          Marx's example is France: rent&nbsp;2 + profit&nbsp;1 + wages&nbsp;3 = income&nbsp;6.
        </p>
        <p className="v3-ch47-note">
          Rates in basis points (bp). 10,000&nbsp;bp&nbsp;= 100%.
        </p>
      </section>

      {/* Section 1 — Rent historical forms */}
      <section className="v3-ch47-section" aria-label="Rent historical forms">
        <h4>1. Rent historical forms</h4>
        <p className="v3-ch47-subintro">
          Record a historical form of ground-rent. Stage&nbsp;1 = labour-rent (extra-economic
          coercion; surplus-labour performed directly); Stage&nbsp;2 = product-rent (share of
          harvest); Stage&nbsp;3 = money-rent (money sum; dissolves personal bond); Stage&nbsp;4 =
          capitalist ground-rent (deduction from surplus-value; three-class separation complete).
        </p>
        <div className="v3-ch47-form">
          <label>
            Stage
            <select
              value={fStage}
              onChange={(e) => setFStage(Number(e.target.value) as RentFormStage)}
            >
              {(Object.keys(STAGE_LABELS) as unknown as RentFormStage[]).map((k) => (
                <option key={k} value={k}>{STAGE_LABELS[k]}</option>
              ))}
            </select>
          </label>
          <label>
            Name
            <input value={fName} onChange={(e) => setFName(e.target.value)} placeholder="e.g. Corvée labour" />
          </label>
          <label>
            Coercion kind
            <input value={fCoercion} onChange={(e) => setFCoercion(e.target.value)} placeholder="e.g. extra-economic / legal / market" />
          </label>
          <label>
            Surplus form
            <input value={fSurplus} onChange={(e) => setFSurplus(e.target.value)} placeholder="e.g. surplus-labour / surplus-product / surplus-value" />
          </label>
          <label>
            Example region
            <input value={fRegion} onChange={(e) => setFRegion(e.target.value)} placeholder="e.g. Danubian principalities" />
          </label>
          <label>
            Example period
            <input value={fPeriod} onChange={(e) => setFPeriod(e.target.value)} placeholder="e.g. 9th–14th century" />
          </label>
        </div>
        <button className="v3-ch47-btn" onClick={handleCreateForm}>
          Record rent form
        </button>
        {fError && <p className="v3-ch47-error">{fError}</p>}
        {lastF && (
          <div className="v3-ch47-preview">
            <strong>
              Saved — {STAGE_LABELS[lastF.stage]}: {lastF.name}
            </strong>
            <span className="v3-ch47-hint">ID: {lastF.id}</span>
          </div>
        )}
        {forms.length > 0 && (
          <div className="v3-ch47-table-wrap">
            <table className="v3-ch47-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Stage</th>
                  <th>Name</th>
                  <th>Coercion</th>
                  <th>Surplus form</th>
                  <th>Region</th>
                  <th>Period</th>
                </tr>
              </thead>
              <tbody>
                {forms.map((f) => (
                  <tr key={f.id}>
                    <td title={f.id}>{f.id.slice(0, 8)}&hellip;</td>
                    <td>{STAGE_LABELS[f.stage]}</td>
                    <td>{f.name}</td>
                    <td>{f.coercion_kind}</td>
                    <td>{f.surplus_form}</td>
                    <td>{f.example_region}</td>
                    <td>{f.example_period}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch47-error">{listError}</p>}
      </section>

      {/* Section 2 — Historical rent transitions */}
      <section className="v3-ch47-section" aria-label="Historical rent transitions">
        <h4>2. Historical rent transitions</h4>
        <p className="v3-ch47-subintro">
          Record a historical transition between rent forms. Transitions are progressive:
          from_stage &lt; to_stage. Each transition is driven by the dissolution of the
          preconditions for the earlier form.
        </p>
        <div className="v3-ch47-form">
          <label>
            From stage
            <select
              value={tFromStage}
              onChange={(e) => setTFromStage(Number(e.target.value) as RentFormStage)}
            >
              {(Object.keys(STAGE_LABELS) as unknown as RentFormStage[]).map((k) => (
                <option key={k} value={k}>{STAGE_LABELS[k]}</option>
              ))}
            </select>
          </label>
          <label>
            To stage
            <select
              value={tToStage}
              onChange={(e) => setTToStage(Number(e.target.value) as RentFormStage)}
            >
              {(Object.keys(STAGE_LABELS) as unknown as RentFormStage[]).map((k) => (
                <option key={k} value={k}>{STAGE_LABELS[k]}</option>
              ))}
            </select>
          </label>
          <label>
            Preconditions
            <input
              value={tPreconditions}
              onChange={(e) => setTPreconditions(e.target.value)}
              placeholder="e.g. serf tied to land; no commodity market"
            />
          </label>
          <label>
            Driving force
            <input
              value={tDrivingForce}
              onChange={(e) => setTDrivingForce(e.target.value)}
              placeholder="e.g. development of commodity production"
            />
          </label>
        </div>
        <button className="v3-ch47-btn" onClick={handleCreateTransition}>
          Record transition
        </button>
        {tError && <p className="v3-ch47-error">{tError}</p>}
        {lastT && (
          <div className="v3-ch47-preview">
            <strong>
              Saved — {STAGE_LABELS[lastT.from_stage]} &rarr; {STAGE_LABELS[lastT.to_stage]}
            </strong>
            <span className="v3-ch47-hint">ID: {lastT.id}</span>
          </div>
        )}
        {transitions.length > 0 && (
          <div className="v3-ch47-table-wrap">
            <table className="v3-ch47-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>From stage</th>
                  <th>To stage</th>
                  <th>Preconditions</th>
                  <th>Driving force</th>
                </tr>
              </thead>
              <tbody>
                {transitions.map((t) => (
                  <tr key={t.id}>
                    <td title={t.id}>{t.id.slice(0, 8)}&hellip;</td>
                    <td>{STAGE_LABELS[t.from_stage]}</td>
                    <td>{STAGE_LABELS[t.to_stage]}</td>
                    <td>{t.preconditions}</td>
                    <td>{t.driving_force}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3 — Small peasant production */}
      <section className="v3-ch47-section" aria-label="Small peasant production">
        <h4>3. Small peasant production</h4>
        <p className="v3-ch47-subintro">
          The small peasant proprietor is simultaneously labourer, capitalist, and landowner — rent,
          profit, and wages are conflated into one undifferentiated income. The server computes{" "}
          <code>total_income_bp = rent + profit + wages</code>. Seed fixture (France): rent&nbsp;2 +
          profit&nbsp;1 + wages&nbsp;3 = income&nbsp;6, conflated.
        </p>
        <div className="v3-ch47-form">
          <label>
            Parcel ID
            <input value={pParcelId} onChange={(e) => setPParcelId(e.target.value)} />
          </label>
          <label>
            Rent component (bp)
            <input
              type="number"
              min={0}
              value={pRent}
              onChange={(e) => setPRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Profit component (bp)
            <input
              type="number"
              min={0}
              value={pProfit}
              onChange={(e) => setPProfit(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Wages component (bp)
            <input
              type="number"
              min={0}
              value={pWages}
              onChange={(e) => setPWages(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch47-hint">
              computed total income = {pRent + pProfit + pWages}&nbsp;bp
            </span>
          </label>
          <label className="v3-ch47-checkbox-label">
            <input
              type="checkbox"
              checked={pConflated}
              onChange={(e) => setPConflated(e.target.checked)}
            />
            Is conflated (rent/profit/wages undifferentiated)
          </label>
        </div>
        <button className="v3-ch47-btn" onClick={handleCreatePeasant}>
          Record peasant production
        </button>
        {pError && <p className="v3-ch47-error">{pError}</p>}
        {lastP && (
          <div className="v3-ch47-preview">
            <strong>
              Saved — parcel&nbsp;{lastP.parcel_id.slice(0, 8)}&hellip;,
              total&nbsp;income&nbsp;{lastP.total_income_bp}&nbsp;bp,
              conflated:&nbsp;{lastP.is_conflated ? "yes" : "no"}
            </strong>
            <span className="v3-ch47-hint">ID: {lastP.id}</span>
          </div>
        )}
        {peasants.length > 0 && (
          <div className="v3-ch47-table-wrap">
            <table className="v3-ch47-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Parcel</th>
                  <th>Rent (bp)</th>
                  <th>Profit (bp)</th>
                  <th>Wages (bp)</th>
                  <th>Total income (bp)</th>
                  <th>Conflated</th>
                </tr>
              </thead>
              <tbody>
                {peasants.map((p) => (
                  <tr key={p.id}>
                    <td title={p.id}>{p.id.slice(0, 8)}&hellip;</td>
                    <td title={p.parcel_id}>{p.parcel_id.slice(0, 8)}&hellip;</td>
                    <td>{p.rent_component_bp}</td>
                    <td>{p.profit_component_bp}</td>
                    <td>{p.wages_component_bp}</td>
                    <td>{p.total_income_bp}</td>
                    <td>{p.is_conflated ? "yes" : "no"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
