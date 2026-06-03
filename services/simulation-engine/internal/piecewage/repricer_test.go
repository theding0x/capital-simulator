package piecewage_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/piecewage"
)

// The repricer lists every piece-wage and reprices each at the pushed factor,
// counting only the agents whose price actually changed (appended = true).
func TestRepricer_RepricesEachWageAndCountsAppended(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var gotFactors []float64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/piece-wages":
			_ = json.NewEncoder(w).Encode([]map[string]string{
				{"agent_id": "a1"},
				{"agent_id": "a2"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/piece-wages/a1/reprice":
			var body struct {
				ProductivityFactor float64 `json:"productivity_factor"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			mu.Lock()
			gotFactors = append(gotFactors, body.ProductivityFactor)
			mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]bool{"appended": true})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/piece-wages/a2/reprice":
			_ = json.NewEncoder(w).Encode(map[string]bool{"appended": false})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	advanced, err := piecewage.New(srv.URL).Reprice(context.Background(), 2.0)
	if err != nil {
		t.Fatalf("reprice: %v", err)
	}
	if advanced != 1 {
		t.Errorf("advanced = %d, want 1 (only a1 appended)", advanced)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(gotFactors) != 1 || gotFactors[0] != 2.0 {
		t.Errorf("forwarded factors = %v, want [2.0]", gotFactors)
	}
}

func TestRepricer_ListErrorReturned(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	if _, err := piecewage.New(srv.URL).Reprice(context.Background(), 2.0); err == nil {
		t.Fatal("expected error when listing piece-wages fails")
	}
}
