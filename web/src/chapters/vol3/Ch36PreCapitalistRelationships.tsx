import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { UsurersCapital } from "../../types";
import "./Ch36PreCapitalistRelationships.css";

function stageLabel(stage: number): string {
  switch (stage) {
    case 1: return "Antiquity";
    case 2: return "Mediaeval";
    case 3: return "Transition";
    case 4: return "Subordinated";
    default: return `Stage ${stage}`;
  }
}

function formatRate(bp: number): string {
  return `${bp} bp (${(bp / 100).toFixed(0)}%)`;
}

export function Ch36PreCapitalistRelationships() {
  const [items, setItems] = useState<UsurersCapital[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Form state
  const [name, setName] = useState("Roman Usurer — Republican Period");
  const [stage, setStage] = useState(1);
  const [wageLabour, setWageLabour] = useState(false);
  const [borrowersText, setBorrowersText] = useState("Roman peasants\nPlebeian smallholders");
  const [interestRateBp, setInterestRateBp] = useState(4800);
  const [subordinatedTo, setSubordinatedTo] = useState("");
  const [description, setDescription] = useState(
    "Ancient Roman usury; interest rate ~48%; operates on simple commodity circulation without wage-labour; drains petty producers without transforming their mode of production",
  );
  const [createError, setCreateError] = useState<string | null>(null);
  const [lastCreated, setLastCreated] = useState<UsurersCapital | null>(null);

  async function refreshList() {
    try {
      const result = await financeApi.listUsurersCapitals();
      setItems(result);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refreshList().catch(() => {});
  }, []);

  async function handleCreate() {
    setCreateError(null);
    setLastCreated(null);
    const borrowers = borrowersText
      .split("\n")
      .map((s) => s.trim())
      .filter((s) => s.length > 0);
    try {
      const result = await financeApi.createUsurersCapital({
        name,
        stage,
        wage_labour: wageLabour,
        borrowers,
        interest_rate_bp: interestRateBp,
        subordinated_to: subordinatedTo,
        description,
      });
      setLastCreated(result);
      await refreshList();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch36">
      <section className="v3-ch36-intro">
        <p>
          Usurer&rsquo;s capital is the <strong>oldest antediluvian form</strong> of
          interest-bearing capital. Long before the capitalist mode of production
          existed, usury flourished in ancient Rome, in the mediaeval communes, and
          throughout the era of petty production. Its precondition is only{" "}
          <strong>commodity circulation and money</strong> — it does not require
          wage-labour or the production of surplus-value in the capitalist sense.
          The usurer lends to the peasant, the petty producer, and the noble;
          he extracts interest without transforming the underlying mode of production.
        </p>
        <p>
          This is what makes usury&rsquo;s capital <strong>parasitic</strong> in the
          pre-capitalist world: it drains the surplus of small producers and
          aristocrats alike while leaving their social relations of production
          intact. In ancient Rome, interest rates of 48&nbsp;% or more were common;
          in the Middle Ages, the Church&rsquo;s formal prohibition of usury coexisted
          with rates of 20&ndash;30&nbsp;% extracted through disguised instruments.
          Usury ruined the independent peasantry and the slave-owner without
          creating a proletariat disciplined by the wage-relation.
        </p>
        <p>
          As capitalism matures, usurer&rsquo;s capital is{" "}
          <strong>subordinated and disciplined</strong>. The modern credit system —
          the bank, the bill of exchange, the joint-stock company — concentrates
          loanable money and transmits it at a rate of interest broadly governed
          by the <strong>general rate of profit</strong>. The independent usurer
          gives way to the banker; extortionate rates are competed down; and
          interest-bearing capital becomes a mere department of industrial capital,
          advancing money on its behalf and appropriating a share of surplus-value
          rather than ruining the producer outright. This chapter (Ch.&nbsp;36)
          completes <strong>Part&nbsp;V</strong> (Ch.&nbsp;21&ndash;36) of
          Volume&nbsp;III.
        </p>
      </section>

      {/* Create form */}
      <section className="v3-ch36-section" aria-label="Record a usurer's capital">
        <h4>Record a usurer&rsquo;s capital</h4>
        <p className="v3-ch36-subintro">
          Enter a historical instance of interest-bearing capital in its
          pre-capitalist or transitional form.
        </p>
        <div className="v3-ch36-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Stage
            <select value={stage} onChange={(e) => setStage(Number(e.target.value))}>
              <option value={1}>1 — Antiquity</option>
              <option value={2}>2 — Mediaeval</option>
              <option value={3}>3 — Transition</option>
              <option value={4}>4 — Subordinated</option>
            </select>
          </label>
          <label style={{ flexDirection: "row", alignItems: "center", gap: "0.5rem" }}>
            <input
              type="checkbox"
              checked={wageLabour}
              onChange={(e) => setWageLabour(e.target.checked)}
              style={{ width: "auto" }}
            />
            Wage-labour present (true only in the Subordinated stage)
          </label>
          <label>
            Borrowers (one debtor class per line)
            <textarea
              value={borrowersText}
              onChange={(e) => setBorrowersText(e.target.value)}
              rows={3}
            />
            <span className="v3-ch36-hint">
              e.g. &ldquo;Roman peasants&rdquo;, &ldquo;Mediaeval nobles&rdquo;, &ldquo;Petty commodity producers&rdquo;
            </span>
          </label>
          <label>
            Interest rate (basis points, e.g. 4800 = 48%)
            <input
              type="number"
              min={0}
              value={interestRateBp}
              onChange={(e) => setInterestRateBp(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch36-hint">{formatRate(interestRateBp)}</span>
          </label>
          <label>
            Subordinated to (leave blank unless stage = Subordinated)
            <input
              value={subordinatedTo}
              onChange={(e) => setSubordinatedTo(e.target.value)}
              placeholder="e.g. industrial_capital"
            />
          </label>
          <label>
            Description
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </label>
        </div>
        <button className="v3-ch36-btn" onClick={handleCreate}>
          Save usurer&rsquo;s capital
        </button>
        {createError && <p className="v3-ch36-error">{createError}</p>}
        {lastCreated && (
          <div className="v3-ch36-preview">
            Saved &mdash; {lastCreated.name} ({stageLabel(lastCreated.stage)},&nbsp;
            {formatRate(lastCreated.interest_rate_bp)})
          </div>
        )}
      </section>

      {/* List */}
      <section className="v3-ch36-section" aria-label="Stored instances of usurer's capital">
        <h4>Stored instances</h4>
        {listError && <p className="v3-ch36-error">{listError}</p>}
        {items.length === 0 ? (
          <p className="v3-ch36-subintro">No instances stored yet.</p>
        ) : (
          <div className="v3-ch36-table-wrap">
            <table className="v3-ch36-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Stage</th>
                  <th>Wage-labour</th>
                  <th>Borrowers</th>
                  <th>Interest rate</th>
                  <th>Subordinated to</th>
                </tr>
              </thead>
              <tbody>
                {items.map((item) => (
                  <tr key={item.id}>
                    <td>{item.name}</td>
                    <td>{stageLabel(item.stage)}</td>
                    <td>
                      {item.wage_labour ? (
                        <span className="v3-ch36-wage-yes">yes</span>
                      ) : (
                        <span className="v3-ch36-wage-no">no</span>
                      )}
                    </td>
                    <td>{item.borrowers.join(", ")}</td>
                    <td>{formatRate(item.interest_rate_bp)}</td>
                    <td>{item.subordinated_to || "—"}</td>
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
