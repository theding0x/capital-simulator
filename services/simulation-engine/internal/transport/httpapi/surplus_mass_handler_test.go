package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postSurplusMass(t *testing.T, body string) int {
	t.Helper()
	h := New(nil, Deps{})
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/surplus/mass", h.ComputeMass)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	res, err := http.Post(ts.URL+"/v1/surplus/mass", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	return res.StatusCode
}

// §1 second law: 100 labourers at 3s each can never reach 24 × 100 working-hours.
func TestComputeMass_RejectsWorkdayAtAbsoluteLimit(t *testing.T) {
	t.Parallel()
	body := `{"necessary_labour":360,"surplus_labour":1080,"labour_power_value":3,"worker_count":100}`
	if got := postSurplusMass(t, body); got != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", got)
	}
}

func TestComputeMass_AcceptsWorkdayBelowAbsoluteLimit(t *testing.T) {
	t.Parallel()
	body := `{"necessary_labour":360,"surplus_labour":1079,"labour_power_value":3,"worker_count":100}`
	if got := postSurplusMass(t, body); got != http.StatusOK {
		t.Fatalf("status = %d, want 200", got)
	}
}
