package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// --- POST /v1/merchant/money-dealing ---

// TestCreateMoneyDealingCapital_Seed1901 posts the generic-dealer seed fixture
// (Vol. III Ch. 19, seed 1901): moneyAdvanced=500000, generalRateBP=1500,
// operationKinds=[1,2,3].
//
//	profit = roundHalfUp(500000*1500, 10000) = 75000; new_value_created = 0
//
// Expects 201 + Location header with the assigned ID.
func TestCreateMoneyDealingCapital_Seed1901(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":500000,"general_rate_bp":1500,"operation_kinds":[1,2,3]}`
	res, err := http.Post(ts.URL+"/v1/merchant/money-dealing", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/merchant/money-dealing/") {
		t.Errorf("Location = %q, want /v1/merchant/money-dealing/ prefix", loc)
	}

	var md merchant.MoneyDealingCapital
	if err := json.NewDecoder(res.Body).Decode(&md); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if md.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if md.CreatedAt.IsZero() {
		t.Error("CreatedAt not set by store")
	}
	if md.MoneyDealingProfit != 75000 {
		t.Errorf("money_dealing_profit = %d, want 75000", md.MoneyDealingProfit)
	}
	if md.NewValueCreated != 0 {
		t.Errorf("new_value_created = %d, want 0", md.NewValueCreated)
	}
	// Verify operation_kinds round-trips as [1,2,3].
	if len(md.OperationKinds) != 3 {
		t.Fatalf("len(operation_kinds) = %d, want 3", len(md.OperationKinds))
	}
	want := []merchant.MoneyOperationKind{
		merchant.KindReceipt, merchant.KindPayment, merchant.KindExchange,
	}
	for i, k := range want {
		if md.OperationKinds[i] != k {
			t.Errorf("operation_kinds[%d] = %d, want %d", i, md.OperationKinds[i], k)
		}
	}
}

// TestCreateMoneyDealingCapital_BadJSON returns 400 on a malformed body.
func TestCreateMoneyDealingCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/merchant/money-dealing", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateMoneyDealingCapital_InvalidKind returns 422 when operation_kinds
// contains 0 (handler maps ErrInvalidOperationKind → 422).
func TestCreateMoneyDealingCapital_InvalidKindZero(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":500000,"general_rate_bp":1500,"operation_kinds":[0]}`
	res, err := http.Post(ts.URL+"/v1/merchant/money-dealing", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for kind=0", res.StatusCode)
	}
}

// TestCreateMoneyDealingCapital_InvalidKindSix returns 422 when operation_kinds
// contains 6 (out of range).
func TestCreateMoneyDealingCapital_InvalidKindSix(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":500000,"general_rate_bp":1500,"operation_kinds":[6]}`
	res, err := http.Post(ts.URL+"/v1/merchant/money-dealing", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for kind=6", res.StatusCode)
	}
}

// TestCreateMoneyDealingCapital_ZeroMoneyAdvanced returns 422 when
// money_advanced=0 (handler maps ErrNonPositiveMoney → 422).
func TestCreateMoneyDealingCapital_ZeroMoneyAdvanced(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":0,"general_rate_bp":1500,"operation_kinds":[1]}`
	res, err := http.Post(ts.URL+"/v1/merchant/money-dealing", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for money_advanced=0", res.StatusCode)
	}
}

// --- GET /v1/merchant/money-dealing ---

// TestListMoneyDealingCapitals_NeverNull asserts the list endpoint returns 200
// and an items array that is non-null even when the store is empty.
func TestListMoneyDealingCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/money-dealing")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		t.Fatalf("decode raw: %v", err)
	}
	items, ok := raw["items"]
	if !ok {
		t.Fatal("response has no items key")
	}
	if string(items) == "null" {
		t.Errorf("items = null, want [] for an empty store")
	}
	var decoded []merchant.MoneyDealingCapital
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// --- GET /v1/merchant/money-dealing/{id} ---

// TestGetMoneyDealingCapital_NotFound returns 404 for an unknown id.
func TestGetMoneyDealingCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/money-dealing/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}
