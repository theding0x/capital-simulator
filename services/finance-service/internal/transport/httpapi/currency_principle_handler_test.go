package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ── POST /v1/credit/note-issue-constraints ───────────────────────────────────

// TestCreateNoteIssueConstraint_BankAct1844 verifies the Bank Act 1844 founding
// fixture: goldBacking=1 400 000 000, unbackedLimit=1 400 000 000
// → 201, Location set, max_notes=2 800 000 000.
func TestCreateNoteIssueConstraint_BankAct1844(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Bank Act 1844 — founding regime","regime":2,"gold_backing":1400000000,"unbacked_limit":1400000000,"notes_beyond_max":0,"description":"Currency School doctrine institutionalised."}`
	res, err := http.Post(ts.URL+"/v1/credit/note-issue-constraints", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", res.StatusCode)
	}
	if res.Header.Get("Location") == "" {
		t.Error("expected Location header")
	}

	var n credit.NoteIssueConstraint
	if err := json.NewDecoder(res.Body).Decode(&n); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if n.ID.IsZero() {
		t.Error("ID must be set")
	}
	if n.MaxNotes != 2_800_000_000 {
		t.Errorf("max_notes: got %d, want 2800000000", n.MaxNotes)
	}
}

// TestCreateNoteIssueConstraint_MaxNotesComputedServerSide verifies that even
// when max_notes is not supplied in the request body the server derives it as
// gold_backing + unbacked_limit.
func TestCreateNoteIssueConstraint_MaxNotesComputedServerSide(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// Omit max_notes — it is not part of the request schema.
	body := `{"name":"Act 1844 — no explicit max_notes","regime":2,"gold_backing":1400000000,"unbacked_limit":1400000000,"notes_beyond_max":0,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/note-issue-constraints", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", res.StatusCode)
	}

	var n credit.NoteIssueConstraint
	if err := json.NewDecoder(res.Body).Decode(&n); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if n.MaxNotes != n.GoldBacking+n.UnbackedLimit {
		t.Errorf("max_notes invariant violated: got %d, want gold_backing+unbacked_limit=%d",
			n.MaxNotes, n.GoldBacking+n.UnbackedLimit)
	}
}

// TestCreateNoteIssueConstraint_BadJSON tests 400 on malformed JSON.
func TestCreateNoteIssueConstraint_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/note-issue-constraints", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateNoteIssueConstraint_InvalidRegime_Zero tests 422 when regime=0.
func TestCreateNoteIssueConstraint_InvalidRegime_Zero(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","regime":0,"gold_backing":1400000000,"unbacked_limit":1400000000,"notes_beyond_max":0,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/note-issue-constraints", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/credit/note-issue-constraints/{id} ───────────────────────────────

// TestGetNoteIssueConstraint_NotFound tests 404 for an unknown ID.
func TestGetNoteIssueConstraint_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/note-issue-constraints/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// ── GET /v1/credit/note-issue-constraints ────────────────────────────────────

// TestListNoteIssueConstraints_NeverNull verifies the list endpoint returns
// {"items":[]} not null when the store is empty.
func TestListNoteIssueConstraints_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/note-issue-constraints")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.NoteIssueConstraint `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// TestNoteIssueConstraint_PostThenList verifies that a created record appears
// in the list response with matching fields.
func TestNoteIssueConstraint_PostThenList(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Crisis of 1847 — Act in force","regime":2,"gold_backing":1400000000,"unbacked_limit":1400000000,"notes_beyond_max":0,"description":"Gold drain with high interest, not high prices."}`
	createRes, err := http.Post(ts.URL+"/v1/credit/note-issue-constraints", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", createRes.StatusCode)
	}

	var created credit.NoteIssueConstraint
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	listRes, err := http.Get(ts.URL + "/v1/credit/note-issue-constraints")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	defer listRes.Body.Close()
	if listRes.StatusCode != http.StatusOK {
		t.Fatalf("list status: got %d, want 200", listRes.StatusCode)
	}

	var payload struct {
		Items []credit.NoteIssueConstraint `json:"items"`
	}
	if err := json.NewDecoder(listRes.Body).Decode(&payload); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(payload.Items) == 0 {
		t.Fatal("items must contain the created record")
	}

	found := false
	for _, item := range payload.Items {
		if item.ID == created.ID {
			found = true
			if item.MaxNotes != created.MaxNotes {
				t.Errorf("list item MaxNotes: got %d, want %d", item.MaxNotes, created.MaxNotes)
			}
			if item.Regime != created.Regime {
				t.Errorf("list item Regime: got %d, want %d", item.Regime, created.Regime)
			}
			break
		}
	}
	if !found {
		t.Errorf("created record ID=%s not found in list", created.ID)
	}
}

// TestNoteIssueConstraint_PostThenGetByID verifies that a created record can
// be retrieved via its ID with all fields intact.
func TestNoteIssueConstraint_PostThenGetByID(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Crisis of 1857 — Act suspended","regime":3,"gold_backing":1400000000,"unbacked_limit":1400000000,"notes_beyond_max":48883000,"description":"Treasury letter suspends the Bank Act."}`
	createRes, err := http.Post(ts.URL+"/v1/credit/note-issue-constraints", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", createRes.StatusCode)
	}

	var created credit.NoteIssueConstraint
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/note-issue-constraints/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.NoteIssueConstraint
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	if fetched.Regime != created.Regime {
		t.Errorf("Regime mismatch: got %d, want %d", fetched.Regime, created.Regime)
	}
	if fetched.MaxNotes != 2_800_000_000 {
		t.Errorf("MaxNotes: got %d, want 2800000000", fetched.MaxNotes)
	}
	if fetched.NotesBeyondMax != 48_883_000 {
		t.Errorf("NotesBeyondMax: got %d, want 48883000", fetched.NotesBeyondMax)
	}
}
