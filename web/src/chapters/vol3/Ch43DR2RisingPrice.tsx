import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type {
  RisingPriceVariant,
  RisingPriceDR2Outcome,
  RentSequence,
  CreateRentSequenceResponse,
  EngelsTableEntry,
} from "../../types";
import "./Ch43DR2RisingPrice.css";

const VARIANT_LABELS: Record<RisingPriceVariant, string> = {
  1: "Eventuality A — constant productivity",
  2: "Eventuality B — decreasing productivity",
  3: "Eventuality B — increasing productivity",
  4: "Eventuality A — total rent doubles",
  5: "Eventuality B — mixed result",
};

export function Ch43DR2RisingPrice() {
  const [engelsEntries, setEngelsEntries] = useState<EngelsTableEntry[]>([]);
  const [rentSequences, setRentSequences] = useState<RentSequence[]>([]);
  const [outcomes, setOutcomes] = useState<RisingPriceDR2Outcome[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Engels Table XI form
  const [tableNumber] = useState<number>(11);
  const [soilGradeLabel, setSoilGradeLabel] = useState("A");
  const [outputBushels, setOutputBushels] = useState(1);
  const [rentShillings, setRentShillings] = useState(0);
  const [capitalShillings, setCapitalShillings] = useState(50);
  const [engelsError, setEngelsError] = useState<string | null>(null);
  const [lastEngels, setLastEngels] = useState<EngelsTableEntry | null>(null);

  // Rent sequence form
  const [baseRent, setBaseRent] = useState(0);
  const [increment, setIncrement] = useState(12);
  const [numGrades, setNumGrades] = useState(5);
  const [sequenceError, setSequenceError] = useState<string | null>(null);
  const [lastSequenceResult, setLastSequenceResult] = useState<CreateRentSequenceResponse | null>(null);

  // Rising-price outcome form
  const [variant, setVariant] = useState<RisingPriceVariant>(1);
  const [totalCapitalShillings, setTotalCapitalShillings] = useState(250);
  const [totalRentShillings, setTotalRentShillings] = useState(120);
  const [totalOutputBushels, setTotalOutputBushels] = useState(70);
  const [pricePerBushelShillings, setPricePerBushelShillings] = useState(6);
  const [outcomeError, setOutcomeError] = useState<string | null>(null);
  const [lastOutcome, setLastOutcome] = useState<RisingPriceDR2Outcome | null>(null);

  async function refreshLists() {
    try {
      const [entries, seqs, outs] = await Promise.all([
        financeApi.listEngelsTableEntries(),
        financeApi.listRentSequences(),
        financeApi.listRisingPriceDR2Outcomes(),
      ]);
      setEngelsEntries(entries);
      setRentSequences(seqs);
      setOutcomes(outs);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshLists().catch(() => {});
  }, []);

  async function handleCreateEngelsEntry() {
    setEngelsError(null);
    setLastEngels(null);
    try {
      const result = await financeApi.createEngelsTableEntry({
        table_number: tableNumber,
        soil_grade_label: soilGradeLabel,
        output_bushels: outputBushels,
        rent_shillings: rentShillings,
        capital_shillings: capitalShillings,
      });
      setLastEngels(result);
      await refreshLists();
    } catch (e) {
      setEngelsError(String(e));
    }
  }

  async function handleCreateRentSequence() {
    setSequenceError(null);
    setLastSequenceResult(null);
    try {
      const result = await financeApi.createRentSequence({
        base_rent: baseRent,
        increment,
        num_grades: numGrades,
      });
      setLastSequenceResult(result);
      await refreshLists();
    } catch (e) {
      setSequenceError(String(e));
    }
  }

  async function handleCreateOutcome() {
    setOutcomeError(null);
    setLastOutcome(null);
    try {
      const result = await financeApi.createRisingPriceDR2Outcome({
        variant,
        total_capital_shillings: totalCapitalShillings,
        total_rent_shillings: totalRentShillings,
        total_output_bushels: totalOutputBushels,
        price_per_bushel_shillings: pricePerBushelShillings,
      });
      setLastOutcome(result);
      await refreshLists();
    } catch (e) {
      setOutcomeError(String(e));
    }
  }

  const totalEngelsRent = engelsEntries
    .filter((e) => e.table_number === 11)
    .reduce((acc, e) => acc + e.rent_shillings, 0);

  return (
    <div className="v3-ch43">
      <section className="v3-ch43-intro">
        <p>
          Chapter&nbsp;43 was left by Marx as a title only; Engels completed it from Marx&apos;s
          notes and tables. It examines DR&nbsp;II under a{" "}
          <strong>rising price of production</strong> — a case that occurs when additional capital
          invested on the same soil yields a <em>below-average</em> return, or when a new and
          inferior soil must be brought into cultivation to meet demand.
        </p>
        <p>
          The key result is that when capital doubles, total rent in most cases{" "}
          <strong>doubles or more</strong>. Engels analyses five eventualities (A and B variants)
          based on whether productivity of the additional capital is constant, decreasing, or
          increasing. Differential rent is an{" "}
          <strong>arithmetic progression</strong> of rent differences measured from the
          rent-free regulating soil. Units throughout are <em>shillings</em>&nbsp;(s) and{" "}
          <em>bushels</em>&nbsp;(bu), following Engels&apos;s own tables.
        </p>
      </section>

      {/* Section 1: Engels Table XI */}
      <section className="v3-ch43-section" aria-label="Engels Table XI">
        <h4>1. Engels Table&nbsp;XI</h4>
        <p className="v3-ch43-subintro">
          Engels&apos;s Table&nbsp;XI shows the five soil grades (A–E) under a rising price of
          production. Each grade produces a given output (bu), requires a given capital (s), and
          yields a rent (s) above the rent-free grade&nbsp;A. The total rent is the sum across
          all grades. Seed: grades A=0s, B=12s, C=24s, D=36s, E=48s (arithmetic progression
          increment 12s, total 120s). Add or browse entries below.
        </p>
        <div className="v3-ch43-form">
          <label>
            Soil grade label
            <input
              value={soilGradeLabel}
              onChange={(e) => setSoilGradeLabel(e.target.value)}
              placeholder="A"
            />
          </label>
          <label>
            Output (bu)
            <input
              type="number"
              min={0}
              value={outputBushels}
              onChange={(e) => setOutputBushels(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Capital (s)
            <input
              type="number"
              min={0}
              value={capitalShillings}
              onChange={(e) => setCapitalShillings(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Rent (s)
            <input
              type="number"
              min={0}
              value={rentShillings}
              onChange={(e) => setRentShillings(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch43-hint">
              Grade A = 0s (regulating, rent-free). Grades B–E follow the arithmetic progression.
            </span>
          </label>
        </div>
        <button className="v3-ch43-btn" onClick={handleCreateEngelsEntry}>
          Add Table&nbsp;XI entry
        </button>
        {engelsError && <p className="v3-ch43-error">{engelsError}</p>}
        {lastEngels && (
          <div className="v3-ch43-preview">
            <strong>
              Saved — grade {lastEngels.soil_grade_label}: {lastEngels.output_bushels}&nbsp;bu,
              capital {lastEngels.capital_shillings}s, rent {lastEngels.rent_shillings}s
            </strong>
            <span className="v3-ch43-hint">ID: {lastEngels.id}</span>
          </div>
        )}
        {engelsEntries.length > 0 && (
          <div className="v3-ch43-table-wrap">
            <table className="v3-ch43-table">
              <thead>
                <tr>
                  <th>Grade</th>
                  <th>Output (bu)</th>
                  <th>Capital (s)</th>
                  <th>Rent (s)</th>
                </tr>
              </thead>
              <tbody>
                {engelsEntries
                  .filter((e) => e.table_number === 11)
                  .map((e) => (
                    <tr key={e.id}>
                      <td>{e.soil_grade_label}</td>
                      <td>{e.output_bushels}</td>
                      <td>{e.capital_shillings}</td>
                      <td>{e.rent_shillings}</td>
                    </tr>
                  ))}
              </tbody>
              <tfoot>
                <tr>
                  <td colSpan={3}>Total rent (Table&nbsp;XI)</td>
                  <td>{totalEngelsRent}s</td>
                </tr>
              </tfoot>
            </table>
          </div>
        )}
        {listError && <p className="v3-ch43-error">{listError}</p>}
      </section>

      {/* Section 2: Rent sequences */}
      <section className="v3-ch43-section" aria-label="Rent sequences">
        <h4>2. Rent sequences</h4>
        <p className="v3-ch43-subintro">
          Differential rent under a rising price follows an arithmetic progression:
          rent[g]&nbsp;=&nbsp;base_rent&nbsp;+&nbsp;(g&minus;1)&nbsp;&times;&nbsp;increment.
          The Engels Table&nbsp;XI progression has base_rent=0, increment=12s, num_grades=5,
          giving the sequence [0, 12, 24, 36, 48] whose sum is 120s.
        </p>
        <div className="v3-ch43-form">
          <label>
            Base rent (s)
            <input
              type="number"
              min={0}
              value={baseRent}
              onChange={(e) => setBaseRent(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Increment per grade (s)
            <input
              type="number"
              min={0}
              value={increment}
              onChange={(e) => setIncrement(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Number of grades
            <input
              type="number"
              min={1}
              value={numGrades}
              onChange={(e) => setNumGrades(Math.max(1, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch43-btn" onClick={handleCreateRentSequence}>
          Compute &amp; save rent sequence
        </button>
        {sequenceError && <p className="v3-ch43-error">{sequenceError}</p>}
        {lastSequenceResult && (
          <div className="v3-ch43-preview">
            <strong>
              Sequence: {lastSequenceResult.sequence.join(", ")}
            </strong>
            <span>
              Sum: {lastSequenceResult.sequence.reduce((a, b) => a + b, 0)}s &mdash; base{" "}
              {lastSequenceResult.base_rent}s, increment {lastSequenceResult.increment}s,{" "}
              {lastSequenceResult.num_grades} grades
            </span>
            <span className="v3-ch43-hint">ID: {lastSequenceResult.id}</span>
          </div>
        )}
        {rentSequences.length > 0 && (
          <div className="v3-ch43-table-wrap">
            <table className="v3-ch43-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Base (s)</th>
                  <th>Increment (s)</th>
                  <th>Grades</th>
                </tr>
              </thead>
              <tbody>
                {rentSequences.map((s) => (
                  <tr key={s.id}>
                    <td title={s.id}>{s.id.slice(0, 8)}&hellip;</td>
                    <td>{s.base_rent}</td>
                    <td>{s.increment}</td>
                    <td>{s.num_grades}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Section 3: Rising-price DR II outcomes */}
      <section className="v3-ch43-section" aria-label="Rising-price DR II outcomes">
        <h4>3. Rising-price DR&nbsp;II outcomes</h4>
        <p className="v3-ch43-subintro">
          Record the aggregate outcome for each of Engels&apos;s five eventualities. Capital and
          rent are in shillings (s); output in bushels (bu). The seed fixture uses eventuality&nbsp;A
          (variant&nbsp;1) at the canonical Table&nbsp;XI values: 250s capital, 120s total rent,
          70&nbsp;bu output, price 6s/bu.
        </p>
        <div className="v3-ch43-form">
          <label>
            Eventuality variant
            <select
              value={variant}
              onChange={(e) => setVariant(Number(e.target.value) as RisingPriceVariant)}
            >
              <option value={1}>{VARIANT_LABELS[1]}</option>
              <option value={2}>{VARIANT_LABELS[2]}</option>
              <option value={3}>{VARIANT_LABELS[3]}</option>
              <option value={4}>{VARIANT_LABELS[4]}</option>
              <option value={5}>{VARIANT_LABELS[5]}</option>
            </select>
          </label>
          <label>
            Total capital (s)
            <input
              type="number"
              min={1}
              value={totalCapitalShillings}
              onChange={(e) => setTotalCapitalShillings(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Total rent (s)
            <input
              type="number"
              min={0}
              value={totalRentShillings}
              onChange={(e) => setTotalRentShillings(Math.max(0, Number(e.target.value)))}
            />
          </label>
          <label>
            Total output (bu)
            <input
              type="number"
              min={1}
              value={totalOutputBushels}
              onChange={(e) => setTotalOutputBushels(Math.max(1, Number(e.target.value)))}
            />
          </label>
          <label>
            Price per bushel (s)
            <input
              type="number"
              min={1}
              value={pricePerBushelShillings}
              onChange={(e) => setPricePerBushelShillings(Math.max(1, Number(e.target.value)))}
            />
          </label>
        </div>
        <button className="v3-ch43-btn" onClick={handleCreateOutcome}>
          Save rising-price DR&nbsp;II outcome
        </button>
        {outcomeError && <p className="v3-ch43-error">{outcomeError}</p>}
        {lastOutcome && (
          <div className="v3-ch43-preview">
            <strong>
              Saved — {VARIANT_LABELS[lastOutcome.variant]}: {lastOutcome.total_rent_shillings}s
              rent, {lastOutcome.total_output_bushels}&nbsp;bu @ {lastOutcome.price_per_bushel_shillings}s/bu
            </strong>
            <span className="v3-ch43-hint">ID: {lastOutcome.id}</span>
          </div>
        )}
        {outcomes.length > 0 && (
          <div className="v3-ch43-table-wrap">
            <table className="v3-ch43-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Variant</th>
                  <th>Capital (s)</th>
                  <th>Rent (s)</th>
                  <th>Output (bu)</th>
                  <th>Price (s/bu)</th>
                </tr>
              </thead>
              <tbody>
                {outcomes.map((o) => (
                  <tr key={o.id}>
                    <td title={o.id}>{o.id.slice(0, 8)}&hellip;</td>
                    <td>{VARIANT_LABELS[o.variant]}</td>
                    <td>{o.total_capital_shillings}</td>
                    <td>{o.total_rent_shillings}</td>
                    <td>{o.total_output_bushels}</td>
                    <td>{o.price_per_bushel_shillings}</td>
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
