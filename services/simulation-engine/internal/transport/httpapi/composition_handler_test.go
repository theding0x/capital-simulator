package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

func compositionHandler() *httpapi.Handler {
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{Composition: mem})
}

func TestCreateComponent_201(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	body := comp.CapitalComponent{
		Role:  comp.RoleInstrumentOfLabour,
		Kind:  comp.KindFixed,
		Pence: 1_000_000,
	}
	b, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/capital-components", bytes.NewReader(b))
	h.CreateComponent(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Location") == "" {
		t.Fatal("Location header must be set")
	}
}

func TestCreateComponent_400_BadRole(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	body := map[string]any{
		"role":  "wizard",
		"kind":  "fixed",
		"pence": 100,
	}
	b, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/capital-components", bytes.NewReader(b))
	h.CreateComponent(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for invalid role, got %d", rec.Code)
	}
}

func TestCreateComponent_400_RoleMismatch(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	// auxiliary_material must be circulating — fixed is Ramsay's error
	body := comp.CapitalComponent{
		Role:  comp.RoleAuxiliaryMaterial,
		Kind:  comp.KindFixed,
		Pence: 500,
	}
	b, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/capital-components", bytes.NewReader(b))
	h.CreateComponent(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for role-kind mismatch, got %d", rec.Code)
	}
}

func TestGetComponent_404(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/capital-components/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	h.GetComponent(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}

func TestListComponents_NonNull(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/capital-components", nil)
	h.ListComponents(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	var items []comp.CapitalComponent
	if err := json.NewDecoder(rec.Body).Decode(&items); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if items == nil {
		t.Fatal("list must not be null")
	}
}

func TestRegisterFixedItem_201(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	// Marx §I: "a machine worth £10,000 lasts for ten years"
	item := comp.FixedCapitalItem{
		Description:        "Ten-year spinning machine",
		PencePurchased:     1_000_000,
		ServiceLifeMinutes: 10 * 525_600,
		WearModel:          comp.WearLinear,
	}
	b, _ := json.Marshal(item)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/fixed-items", bytes.NewReader(b))
	h.RegisterFixedItem(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRegisterFixedItem_400_ZeroLife(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	item := comp.FixedCapitalItem{
		PencePurchased:     1_000_000,
		ServiceLifeMinutes: 0, // must be positive
	}
	b, _ := json.Marshal(item)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/fixed-items", bytes.NewReader(b))
	h.RegisterFixedItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for zero service life, got %d", rec.Code)
	}
}

func TestGetFixedItem_404(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/fixed-items/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	h.GetFixedItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}

func TestGetSinkingFund_AutoCreated(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	item := comp.FixedCapitalItem{
		PencePurchased:     500_000,
		ServiceLifeMinutes: 5 * 525_600,
	}
	b, _ := json.Marshal(item)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/fixed-items", bytes.NewReader(b))
	h.RegisterFixedItem(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: %d %s", rec.Code, rec.Body.String())
	}
	var created comp.FixedCapitalItem
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/v1/fixed-items/"+string(created.ID)+"/sinking-fund", nil)
	req2.SetPathValue("id", string(created.ID))
	h.GetSinkingFund(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("want 200 for sinking fund, got %d", rec2.Code)
	}
	var sf comp.SinkingFundForItem
	if err := json.NewDecoder(rec2.Body).Decode(&sf); err != nil {
		t.Fatalf("decode sf: %v", err)
	}
	if sf.TargetPence != 500_000 {
		t.Errorf("target_pence = %d, want 500000", sf.TargetPence)
	}
	if sf.AccumulatedPence != 0 {
		t.Errorf("accumulated_pence = %d, want 0 initially", sf.AccumulatedPence)
	}
}

func TestRecordWear_AccumulatesInSinkingFund(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	item := comp.FixedCapitalItem{
		PencePurchased:     1_000_000,
		ServiceLifeMinutes: 10 * 525_600,
	}
	b, _ := json.Marshal(item)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/fixed-items", bytes.NewReader(b))
	h.RegisterFixedItem(rec, req)
	var created comp.FixedCapitalItem
	_ = json.NewDecoder(rec.Body).Decode(&created)

	// Record one year of wear (100,000 pence = £1,000 per year on £10,000 machine)
	wear := comp.WearAndTear{
		Pence:                100_000,
		PeriodCoveredMinutes: 525_600,
		Cause:                comp.WearUse,
	}
	bw, _ := json.Marshal(wear)
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/v1/fixed-items/"+string(created.ID)+"/wear", bytes.NewReader(bw))
	req2.SetPathValue("id", string(created.ID))
	h.RecordWear(rec2, req2)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("record wear: %d %s", rec2.Code, rec2.Body.String())
	}

	rec3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/v1/fixed-items/"+string(created.ID)+"/sinking-fund", nil)
	req3.SetPathValue("id", string(created.ID))
	h.GetSinkingFund(rec3, req3)
	var sf comp.SinkingFundForItem
	_ = json.NewDecoder(rec3.Body).Decode(&sf)
	if sf.AccumulatedPence != 100_000 {
		t.Errorf("accumulated_pence = %d, want 100000 after one year of wear", sf.AccumulatedPence)
	}
}

func TestRecordRepair_Running_400_IfExtendsLife(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	item := comp.FixedCapitalItem{PencePurchased: 100_000, ServiceLifeMinutes: 525_600}
	b, _ := json.Marshal(item)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/fixed-items", bytes.NewReader(b))
	h.RegisterFixedItem(rec, req)
	var created comp.FixedCapitalItem
	_ = json.NewDecoder(rec.Body).Decode(&created)

	repair := comp.Repair{
		Kind:                        comp.RepairRunning,
		Pence:                       500,
		OccurredAt:                  time.Now(),
		ServiceLifeExtensionMinutes: 60, // must be 0 for running repair
	}
	br, _ := json.Marshal(repair)
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/v1/fixed-items/"+string(created.ID)+"/repairs", bytes.NewReader(br))
	req2.SetPathValue("id", string(created.ID))
	h.RecordRepair(rec2, req2)
	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for running repair with life extension, got %d", rec2.Code)
	}
}

func TestRecordCirculatingCycle_201(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	c := comp.CirculatingCycle{
		CapitalComponentID: comp.NewCapitalComponentID(),
		StartedAt:          time.Now().Add(-24 * time.Hour),
		EndedAt:            time.Now(),
		Pence:              50_000,
	}
	b, _ := json.Marshal(c)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/circulating-cycles", bytes.NewReader(b))
	h.RecordCirculatingCycle(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateComponent_400_BadJSON(t *testing.T) {
	t.Parallel()
	h := compositionHandler()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/capital-components", bytes.NewReader([]byte("not json")))
	h.CreateComponent(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for bad JSON, got %d", rec.Code)
	}
}
