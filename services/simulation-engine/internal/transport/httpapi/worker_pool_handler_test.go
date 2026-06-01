package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// tickWithBody POSTs a tick with an optional JSON body (empty string = no body)
// and decodes the response.
func tickWithBody(t *testing.T, h *Handler, id, body string) reproductionTickResponse {
	t.Helper()
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(http.MethodPost, "/v1/reproduction/simple/schemes/"+id+"/tick", nil)
	} else {
		req = httptest.NewRequest(http.MethodPost,
			"/v1/reproduction/simple/schemes/"+id+"/tick", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.AdvanceReproductionTick(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("tick: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp reproductionTickResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode tick: %v", err)
	}
	return resp
}

// marxScheme1865 creates the canonical scheme + both departments (pence):
// I.v=1,000,000 + II.v=500,000 ⇒ wage bill 1,500,000.
func marxScheme1865(t *testing.T, h *Handler, period string) string {
	t.Helper()
	id := createSimpleScheme(t, h, period)
	addSimpleDept(t, h, id, `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`)
	addSimpleDept(t, h, id, `{"department":"department_ii","constant_pence":2000000,"variable_pence":500000,"surplus_pence":500000,"total_pence":3000000}`)
	return id
}

// TestTick_LabourForceShortfall is the issue #220 acceptance: a tick whose wage
// bill cannot cover the worker pool's subsistence is flagged labour-force-shortfall.
func TestTick_LabourForceShortfall(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)
	id := marxScheme1865(t, h, "1865-labour")

	// 10 workers × 200,000 = 2,000,000 required > 1,500,000 wage bill.
	resp := tickWithBody(t, h, id, `{"worker_pool_size":10,"subsistence_basket_pence":200000}`)

	if resp.WageBillPence != 1500000 {
		t.Fatalf("wage bill: got %d, want 1500000 (I.v + II.v)", resp.WageBillPence)
	}
	if resp.WorkerPoolSize != 10 {
		t.Fatalf("worker pool: got %d, want 10", resp.WorkerPoolSize)
	}
	if resp.SubsistenceCovered {
		t.Fatal("expected subsistence NOT covered (2,000,000 required > 1,500,000 wage bill)")
	}
	if resp.Status != "labour-force-shortfall" {
		t.Fatalf("status: got %q, want labour-force-shortfall", resp.Status)
	}
}

// TestTick_LabourForceCovered: a wage bill that covers subsistence is reproduced.
func TestTick_LabourForceCovered(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)
	id := marxScheme1865(t, h, "1866-labour")

	// 10 workers × 100,000 = 1,000,000 required <= 1,500,000 wage bill.
	resp := tickWithBody(t, h, id, `{"worker_pool_size":10,"subsistence_basket_pence":100000}`)
	if !resp.SubsistenceCovered || resp.Status != "reproduced" {
		t.Fatalf("expected covered/reproduced, got covered=%v status=%q", resp.SubsistenceCovered, resp.Status)
	}
}

// TestTick_NoBodyUntracked: a tick with no body skips the labour-force check
// (backward compatible) and is reproduced.
func TestTick_NoBodyUntracked(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)
	id := marxScheme1865(t, h, "1867-labour")

	resp := tickWithBody(t, h, id, "")
	if resp.WorkerPoolSize != 0 || !resp.SubsistenceCovered || resp.Status != "reproduced" {
		t.Fatalf("no-body tick should be untracked/reproduced, got pool=%d covered=%v status=%q",
			resp.WorkerPoolSize, resp.SubsistenceCovered, resp.Status)
	}
}
