package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

func TestListProductionRelationBasesEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListProductionRelationBases(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/production-bases", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Fatalf("status=%d body=%s, want 200 items:[]", rec.Code, rec.Body.String())
	}
}

func TestListDistributionRelationsEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListDistributionRelations(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/distribution-relations", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Fatalf("status=%d body=%s, want 200 items:[]", rec.Code, rec.Body.String())
	}
}

func TestListHistoricalTransiencesEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListHistoricalTransiences(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/historical-transiences", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Fatalf("status=%d body=%s, want 200 items:[]", rec.Code, rec.Body.String())
	}
}

func TestCreateProductiveForceContradictionCreated(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	body := `{"productive_force":"socialised production","social_form":"private property","contradiction_nature":"form vs force","crisis_signal":"crises of 1825/1847/1857"}`
	rec := httptest.NewRecorder()
	h.CreateProductiveForceContradiction(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/productive-force-contradictions", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateProductiveForceContradictionBadJSON(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.CreateProductiveForceContradiction(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/productive-force-contradictions", strings.NewReader("{bad")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestCreateProductiveForceContradictionEmptyRejected(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.CreateProductiveForceContradiction(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/productive-force-contradictions", strings.NewReader(`{"productive_force":"","social_form":""}`)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}
