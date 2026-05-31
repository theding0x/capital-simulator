import { useEffect, useState } from "react";
import { financeApi } from "../../api";
import type { NoteIssueConstraint } from "../../types";
import "./Ch34CurrencyPrinciple.css";

function penceToPounds(pence: number): string {
  return `£${(pence / 100).toLocaleString("en-GB", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
}

function regimeLabel(regime: number): string {
  switch (regime) {
    case 1: return "Pre-Act 1844";
    case 2: return "Bank Act 1844";
    case 3: return "Suspended";
    default: return `Regime ${regime}`;
  }
}

export function Ch34CurrencyPrinciple() {
  const [records, setRecords] = useState<NoteIssueConstraint[]>([]);
  const [listError, setListError] = useState<string | null>(null);

  // Form state
  const [name, setName] = useState("Bank of England Issue Dept. 1844");
  const [regime, setRegime] = useState(2);
  const [goldBacking, setGoldBacking] = useState(1_400_000_000);
  const [unbackedLimit, setUnbackedLimit] = useState(1_400_000_000);
  const [notesBeyondMax, setNotesBeyondMax] = useState(0);
  const [description, setDescription] = useState(
    "Peel's Act: notes backed 1:1 by gold above the £14m fiduciary limit",
  );
  const [createError, setCreateError] = useState<string | null>(null);
  const [lastCreated, setLastCreated] = useState<NoteIssueConstraint | null>(null);

  const maxNotes = goldBacking + unbackedLimit;

  async function refresh() {
    try {
      const items = await financeApi.listNoteIssueConstraints();
      setRecords(items);
    } catch (e) {
      setListError(String(e));
    }
  }

  useEffect(() => {
    refresh().catch(() => {});
  }, []);

  async function handleCreate() {
    setCreateError(null);
    setLastCreated(null);
    try {
      const result = await financeApi.createNoteIssueConstraint({
        name,
        regime,
        gold_backing: goldBacking,
        unbacked_limit: unbackedLimit,
        notes_beyond_max: notesBeyondMax,
        description,
      });
      setLastCreated(result);
      await refresh();
    } catch (e) {
      setCreateError(String(e));
    }
  }

  return (
    <div className="v3-ch34">
      <section className="v3-ch34-intro">
        <p>
          The <strong>Currency School</strong> (Overstone, Torrens, Norman) held that
          banknotes should behave exactly like gold coin: whenever gold drained abroad,
          the note circulation should contract by an equal amount, raising the interest
          rate and reversing the drain. This was Ricardo's quantity-theory transplanted
          from coin to paper. The <strong>Banking School</strong> (Tooke, Fullarton,
          Wilson) replied that notes issued against bills of exchange return automatically
          to the bank when the bills mature — the <em>law of reflux</em> — so no
          mechanical rule is necessary.
        </p>
        <p>
          <strong>Peel's Bank Act of 1844</strong> enacted the Currency Principle by
          splitting the Bank of England into two legally separate departments. The{" "}
          <strong>Issue Department</strong> could issue notes only up to a fixed
          fiduciary limit (£14&nbsp;million) beyond its gold reserve; every additional
          note required a pound of gold. The <strong>Banking Department</strong> then
          lent notes obtained from Issue. The partition made the note supply rigid: a
          gold drain contracted the circulation mechanically, exactly as the Currency
          School demanded.
        </p>
        <p>
          Marx shows that the Act <strong>intensified crises rather than
          preventing them</strong>. It had to be suspended in{" "}
          <strong>1847</strong> and again in <strong>1857</strong> — each time the
          Banking Department's reserve ran out before the drain stopped, and the
          government authorised the Bank to issue beyond the legal maximum. The Act
          confused cause and effect: the gold drain was a <em>symptom</em> of
          over-extension, not its cause; contracting the note supply only made the
          symptom lethal.
        </p>
        <p>
          <strong>Hubbard's refutation</strong> of the quantity theory was empirical.
          He showed that gold and note circulation moved <em>with</em> the rate of
          interest, not against the price level: when interest rose, gold flowed{" "}
          <em>in</em>, not out, because a high rate of return attracted foreign capital.
          The Currency School's chain (too many notes → rising prices → gold outflow)
          simply did not appear in the data.
        </p>
      </section>

      {/* Note-issue-constraint form */}
      <section className="v3-ch34-section" aria-label="Record a note-issue constraint">
        <h4>Record a note-issue constraint</h4>
        <p className="v3-ch34-subintro">
          Enter the gold backing and the fiduciary (unbacked) limit. The server
          derives Max&nbsp;Notes&nbsp;=&nbsp;GoldBacking&nbsp;+&nbsp;UnbackedLimit.
          Notes&nbsp;Beyond&nbsp;Max records any excess issued under a suspension.
        </p>
        <div className="v3-ch34-form">
          <label>
            Name
            <input value={name} onChange={(e) => setName(e.target.value)} />
          </label>
          <label>
            Regime
            <select value={regime} onChange={(e) => setRegime(Number(e.target.value))}>
              <option value={1}>1 — Pre-Act 1844</option>
              <option value={2}>2 — Bank Act 1844</option>
              <option value={3}>3 — Suspended</option>
            </select>
          </label>
          <label>
            Gold Backing (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={goldBacking}
              onChange={(e) => setGoldBacking(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch34-hint">{penceToPounds(goldBacking)}</span>
          </label>
          <label>
            Unbacked Limit (pence, &ge; 0)
            <input
              type="number"
              min={0}
              value={unbackedLimit}
              onChange={(e) => setUnbackedLimit(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch34-hint">{penceToPounds(unbackedLimit)}</span>
          </label>
          <label>
            Notes Beyond Max (pence, &ge; 0 — suspension excess)
            <input
              type="number"
              min={0}
              value={notesBeyondMax}
              onChange={(e) => setNotesBeyondMax(Math.max(0, Number(e.target.value)))}
            />
            <span className="v3-ch34-hint">{penceToPounds(notesBeyondMax)}</span>
          </label>
          <label>
            Description
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </label>
        </div>
        <div className="v3-ch34-preview">
          Max Notes = {penceToPounds(goldBacking)} + {penceToPounds(unbackedLimit)} ={" "}
          <strong>{penceToPounds(maxNotes)}</strong>
        </div>
        <button className="v3-ch34-btn" onClick={handleCreate}>
          Save constraint
        </button>
        {createError && <p className="v3-ch34-error">{createError}</p>}
        {lastCreated && (
          <div className="v3-ch34-preview">
            Saved — Max Notes: {penceToPounds(lastCreated.max_notes)} &nbsp;|&nbsp;
            Regime: {regimeLabel(lastCreated.regime)}
          </div>
        )}
      </section>

      {/* Stored records */}
      <section className="v3-ch34-section" aria-label="Stored note-issue constraints">
        <h4>Stored note-issue constraints</h4>
        {listError && <p className="v3-ch34-error">{listError}</p>}
        {records.length === 0 ? (
          <p className="v3-ch34-subintro">No constraints stored yet.</p>
        ) : (
          <div className="v3-ch34-table-wrap">
            <table className="v3-ch34-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Regime</th>
                  <th>Gold Backing</th>
                  <th>Unbacked Limit</th>
                  <th>Max Notes</th>
                  <th>Notes Beyond Max</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                {records.map((r) => (
                  <tr key={r.id}>
                    <td>{r.name}</td>
                    <td>
                      <span className={`v3-ch34-regime v3-ch34-regime-${r.regime}`}>
                        {regimeLabel(r.regime)}
                      </span>
                    </td>
                    <td>{penceToPounds(r.gold_backing)}</td>
                    <td>{penceToPounds(r.unbacked_limit)}</td>
                    <td>{penceToPounds(r.max_notes)}</td>
                    <td>{penceToPounds(r.notes_beyond_max)}</td>
                    <td>{r.description}</td>
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
