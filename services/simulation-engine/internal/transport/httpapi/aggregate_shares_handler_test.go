package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newShareProjectionHandler wires one in-memory store to the three stores the
// Ch. 20/21 → Ch. 17 department-share projection touches (issue #215).
func newShareProjectionHandler(t *testing.T) *Handler {
	t.Helper()
	mem := store.NewMemory()
	return New(nil, Deps{
		SimpleReproduction:   mem,
		ExtendedReproduction: mem,
		SurplusCirculations:  mem,
	})
}

// aggregateShares reads the department-share split for a period via the
// social-aggregate endpoint.
func aggregateShares(t *testing.T, h *Handler, period string) (deptI, deptII int64) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/reproduction/social-aggregate?period="+period, nil)
	w := httptest.NewRecorder()
	h.GetSocialAggregate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GetSocialAggregate(%s): expected 200, got %d: %s", period, w.Code, w.Body.String())
	}
	var resp struct {
		DepartmentIShareBP  int64 `json:"department_i_share_bp"`
		DepartmentIIShareBP int64 `json:"department_ii_share_bp"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode aggregate: %v", err)
	}
	return resp.DepartmentIShareBP, resp.DepartmentIIShareBP
}

func createSimpleScheme(t *testing.T, h *Handler, period string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/simple/schemes",
		bytes.NewBufferString(`{"period":"`+period+`"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateSimpleReproductionScheme(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create simple scheme: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode scheme: %v", err)
	}
	id, _ := resp["id"].(string)
	if id == "" {
		t.Fatal("created scheme has empty id")
	}
	return id
}

func addSimpleDept(t *testing.T, h *Handler, id, deptJSON string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/simple/schemes/"+id+"/departments", bytes.NewBufferString(deptJSON))
	req.SetPathValue("id", id)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.AddDepartmentToScheme(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("add department: expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func tickSimpleScheme(t *testing.T, h *Handler, id string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/simple/schemes/"+id+"/tick", nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.AdvanceReproductionTick(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("tick simple scheme: expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func tickExtendedSchemeByID(t *testing.T, h *Handler, id string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/extended/schemes/"+id+"/tick", nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.handleTickExtendedScheme(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("tick extended scheme: expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

// TestSimpleReproductionTick_ProjectsDepartmentShares is the issue #215
// acceptance test for Ch. 20: ticking the canonical 1865 fixture populates the
// social-capital aggregate's department shares with 6667/3333 (sum 10000).
func TestSimpleReproductionTick_ProjectsDepartmentShares(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)

	if dI, dII := aggregateShares(t, h, "1865"); dI != 0 || dII != 0 {
		t.Fatalf("pre-tick shares should be 0/0, got %d/%d", dI, dII)
	}

	id := createSimpleScheme(t, h, "1865")
	addSimpleDept(t, h, id, `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`)
	addSimpleDept(t, h, id, `{"department":"department_ii","constant_pence":2000000,"variable_pence":500000,"surplus_pence":500000,"total_pence":3000000}`)

	tickSimpleScheme(t, h, id)

	dI, dII := aggregateShares(t, h, "1865")
	if dI != 6667 || dII != 3333 {
		t.Fatalf("post-tick shares: got %d/%d, want 6667/3333", dI, dII)
	}
	if dI+dII != 10000 {
		t.Fatalf("shares must sum to 10000, got %d", dI+dII)
	}
}

// TestExtendedReproductionTick_ProjectsDepartmentShares is the issue #215
// acceptance test for Ch. 21: an extended-reproduction tick populates the
// period's aggregate shares (non-zero, summing to 10000).
func TestExtendedReproductionTick_ProjectsDepartmentShares(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)

	id, code := createExtendedSchemeInner(t, h,
		`{"period":"1865-ext","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if code != http.StatusCreated || id == "" {
		t.Fatalf("create extended scheme: code %d, id %q", code, id)
	}
	if c := addExtendedDeptInner(t, h, id,
		`{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`); c != http.StatusCreated {
		t.Fatalf("add extended deptI: %d", c)
	}
	if c := addExtendedDeptInner(t, h, id,
		`{"department":"department_ii","constant_pence":1500000,"variable_pence":750000,"surplus_pence":750000,"total_pence":3000000}`); c != http.StatusCreated {
		t.Fatalf("add extended deptII: %d", c)
	}

	tickExtendedSchemeByID(t, h, id)

	dI, dII := aggregateShares(t, h, "1865-ext")
	if dI <= 0 || dII <= 0 {
		t.Fatalf("post-tick extended shares must be non-zero, got %d/%d", dI, dII)
	}
	if dI+dII != 10000 {
		t.Fatalf("shares must sum to 10000, got %d", dI+dII)
	}
}
