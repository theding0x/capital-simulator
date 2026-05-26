package httpapi_test

// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital:
// The Physiocrats and Adam Smith.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

func TestListEconomistAttributions_All(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, httpapi.Deps{EconomistAttributions: store.NewMemory()})

	req := httptest.NewRequest(http.MethodGet, "/v1/circulation/economist-attributions", nil)
	rec := httptest.NewRecorder()
	h.ListEconomistAttributions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var attrs []circulation.EconomistAttribution
	if err := json.NewDecoder(rec.Body).Decode(&attrs); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	// Seed data has 3 records (2 Quesnay + 1 Smith).
	if len(attrs) != 3 {
		t.Fatalf("expected 3 attributions, got %d", len(attrs))
	}
}

func TestListEconomistAttributions_FilterQuesnay(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, httpapi.Deps{EconomistAttributions: store.NewMemory()})

	req := httptest.NewRequest(http.MethodGet, "/v1/circulation/economist-attributions?theorist=Quesnay", nil)
	rec := httptest.NewRecorder()
	h.ListEconomistAttributions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var attrs []circulation.EconomistAttribution
	if err := json.NewDecoder(rec.Body).Decode(&attrs); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(attrs) != 2 {
		t.Fatalf("expected 2 Quesnay attributions, got %d", len(attrs))
	}
	for _, a := range attrs {
		if a.Theorist != "Quesnay" {
			t.Errorf("expected theorist Quesnay, got %q", a.Theorist)
		}
		if len(a.Errors) != 0 {
			t.Errorf("Quesnay attribution must carry no errors; got %v", a.Errors)
		}
	}
}

func TestListEconomistAttributions_FilterSmith(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, httpapi.Deps{EconomistAttributions: store.NewMemory()})

	req := httptest.NewRequest(http.MethodGet, "/v1/circulation/economist-attributions?theorist=Smith", nil)
	rec := httptest.NewRecorder()
	h.ListEconomistAttributions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var attrs []circulation.EconomistAttribution
	if err := json.NewDecoder(rec.Body).Decode(&attrs); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(attrs) != 1 {
		t.Fatalf("expected 1 Smith attribution, got %d", len(attrs))
	}
	smith := attrs[0]
	if smith.EditionYear != 1776 {
		t.Errorf("expected Smith edition 1776, got %d", smith.EditionYear)
	}
	if len(smith.Errors) != 3 {
		t.Errorf("expected 3 errors for Smith, got %d: %v", len(smith.Errors), smith.Errors)
	}
}

func TestListEconomistAttributions_UnknownTheorist_ReturnsEmpty(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, httpapi.Deps{EconomistAttributions: store.NewMemory()})

	req := httptest.NewRequest(http.MethodGet, "/v1/circulation/economist-attributions?theorist=Ricardo", nil)
	rec := httptest.NewRecorder()
	h.ListEconomistAttributions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for unknown theorist, got %d", rec.Code)
	}
	var attrs []circulation.EconomistAttribution
	if err := json.NewDecoder(rec.Body).Decode(&attrs); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if attrs == nil {
		t.Fatal("response must be non-nil empty slice, not null")
	}
	if len(attrs) != 0 {
		t.Fatalf("expected 0 attributions for Ricardo, got %d", len(attrs))
	}
}

func TestListEconomistAttributions_ResponseIsNonNull(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, httpapi.Deps{EconomistAttributions: store.NewMemory()})

	req := httptest.NewRequest(http.MethodGet, "/v1/circulation/economist-attributions?theorist=NoOne", nil)
	rec := httptest.NewRecorder()
	h.ListEconomistAttributions(rec, req)

	body := rec.Body.String()
	// JSON null would be "null\n"; an empty array is "[]\n".
	if body == "null\n" || body == "null" {
		t.Fatalf("response body must not be null; got %q", body)
	}
}
