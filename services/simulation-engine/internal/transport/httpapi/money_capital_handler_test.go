package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// newMCHandler returns a Handler wired to a fresh in-memory store for Ch. 18 tests.
func newMCHandler(t *testing.T) (*httpapi.Handler, *store.Memory) {
	t.Helper()
	mem := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{MoneyCapital: mem})
	return h, mem
}

// --- POST /v1/reproduction/apportionments ---

// TestHandleCreateApportionment verifies 201 on a valid request, 400 on
// malformed JSON, and 400 on an invalid (sum-mismatch) apportionment.
func TestHandleCreateApportionment(t *testing.T) {
	t.Parallel()

	t.Run("valid apportionment returns 201", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		// 400000 + 350000 + 200000 + 50000 = 1000000 == total
		body := `{
			"period": "1865",
			"total_social_money_pence": 1000000,
			"dept_i_reserve_pence": 400000,
			"dept_ii_reserve_pence": 350000,
			"wage_rotation_pence": 200000,
			"idle_hoard_pence": 50000
		}`
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/apportionments", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateApportionment(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var resp map[string]any
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp["id"] == "" || resp["id"] == nil {
			t.Fatal("expected id in response")
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/apportionments", bytes.NewBufferString(`{bad json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateApportionment(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("sum mismatch returns 400", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		// 400000 + 300000 + 200000 + 50000 = 950000 ≠ 1000000
		body := `{
			"period": "1865",
			"total_social_money_pence": 1000000,
			"dept_i_reserve_pence": 400000,
			"dept_ii_reserve_pence": 300000,
			"wage_rotation_pence": 200000,
			"idle_hoard_pence": 50000
		}`
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/apportionments", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateApportionment(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
		}
	})
}

// --- POST /v1/reproduction/department-reserves ---

// TestHandleCreateDepartmentReserve verifies 201 on a valid request and 400 on
// malformed JSON.
func TestHandleCreateDepartmentReserve(t *testing.T) {
	t.Parallel()

	t.Run("valid Department I reserve returns 201", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		body := `{
			"period": "1865",
			"department": "department_i",
			"reserve_purpose": "means_of_production",
			"reserve_pence": 400000
		}`
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/department-reserves", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateDepartmentReserve(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/department-reserves", bytes.NewBufferString(`{bad`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateDepartmentReserve(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

// --- POST /v1/reproduction/circulating-money-masses ---

// TestHandleCreateCirculatingMoneyMass verifies that the server computes
// EffectiveCirculatingValuePence from the supplied stock and velocity and that
// malformed JSON returns 400.
func TestHandleCreateCirculatingMoneyMass(t *testing.T) {
	t.Parallel()

	t.Run("server computes effective circulation from stock and velocity", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		// money_stock_pence=500000, velocity=20000 → effective=1000000
		body := `{
			"period": "1865",
			"money_stock_pence": 500000,
			"velocity_per_year_basis_points": 20000
		}`
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/circulating-money-masses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateCirculatingMoneyMass(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
		var resp map[string]any
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		// The server must compute effective_circulating_value_pence = 500000 * 20000 / 10000 = 1000000
		effective, ok := resp["effective_circulating_value_pence"].(float64)
		if !ok {
			t.Fatalf("expected effective_circulating_value_pence in response, got %v", resp)
		}
		if int64(effective) != 1_000_000 {
			t.Fatalf("effective_circulating_value_pence: want 1000000, got %d", int64(effective))
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/circulating-money-masses", bytes.NewBufferString(`{bad`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateCirculatingMoneyMass(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

// --- POST /v1/reproduction/wage-rotation-funds ---

// TestHandleCreateWageRotationFund verifies 201 on valid input and 400 on
// malformed JSON.
func TestHandleCreateWageRotationFund(t *testing.T) {
	t.Parallel()

	t.Run("valid weekly fund returns 201", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		body := `{
			"period": "1865",
			"fund_pence": 200000,
			"wage_cycle_frequency": 52
		}`
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/wage-rotation-funds", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateWageRotationFund(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/wage-rotation-funds", bytes.NewBufferString(`{bad`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateWageRotationFund(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

// --- POST /v1/reproduction/inter-department-settlements ---

// TestHandleCreateInterDepartmentSettlement verifies 201 on a valid balanced
// I→II settlement and 400 on malformed JSON.
func TestHandleCreateInterDepartmentSettlement(t *testing.T) {
	t.Parallel()

	t.Run("valid I→II settlement returns 201", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		// Under simple reproduction the inter-department exchange nets to zero:
		// Department I supplies means of production to Department II.
		body := `{
			"period": "1865",
			"from_department": "department_i",
			"to_department": "department_ii",
			"settled_pence": 2000000,
			"settlement_purpose": "means_of_production"
		}`
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/inter-department-settlements", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateInterDepartmentSettlement(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/inter-department-settlements", bytes.NewBufferString(`{bad`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateInterDepartmentSettlement(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})
}

// --- GET /v1/reproduction/apportionments/{id}/balance-check ---

// TestHandleGetApportionmentBalanceCheck verifies 200 with balanced=true for a
// correctly apportioned record and 404 for a missing ID.
func TestHandleGetApportionmentBalanceCheck(t *testing.T) {
	t.Parallel()

	t.Run("balanced apportionment returns 200 with balanced=true", func(t *testing.T) {
		t.Parallel()
		h, mem := newMCHandler(t)
		ctx := t.Context()

		// Create a balanced record directly via the store to seed the state.
		a := repro.MoneySupplyApportionment{
			ID:                    repro.NewMoneySupplyApportionmentID(),
			Period:                "1865",
			TotalSocialMoneyPence: 1_000_000,
			DeptIReservePence:     400_000,
			DeptIIReservePence:    350_000,
			WageRotationPence:     200_000,
			IdleHoardPence:        50_000,
		}
		created, err := mem.CreateMoneySupplyApportionment(ctx, a)
		if err != nil {
			t.Fatalf("seed apportionment: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet,
			"/v1/reproduction/apportionments/"+string(created.ID)+"/balance-check", nil)
		req.SetPathValue("id", string(created.ID))
		w := httptest.NewRecorder()
		h.GetApportionmentBalanceCheck(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp map[string]any
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		balanced, ok := resp["balanced"].(bool)
		if !ok {
			t.Fatalf("expected bool 'balanced' field, got %T: %v", resp["balanced"], resp["balanced"])
		}
		if !balanced {
			t.Fatal("expected balanced=true")
		}
		id, ok := resp["apportionment_id"].(string)
		if !ok || id == "" {
			t.Fatalf("expected non-empty apportionment_id, got %v", resp["apportionment_id"])
		}
	})

	t.Run("non-existent ID returns 404", func(t *testing.T) {
		t.Parallel()
		h, _ := newMCHandler(t)
		req := httptest.NewRequest(http.MethodGet,
			"/v1/reproduction/apportionments/nonexistent/balance-check", nil)
		req.SetPathValue("id", "nonexistent")
		w := httptest.NewRecorder()
		h.GetApportionmentBalanceCheck(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}
