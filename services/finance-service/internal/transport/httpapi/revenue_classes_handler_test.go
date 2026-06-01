package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

func TestListSocialClassesEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListSocialClasses(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/classes", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Fatalf("status=%d body=%s, want 200 items:[]", rec.Code, rec.Body.String())
	}
}

func TestGetSocialClassRoundTrip(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	saved, err := mem.CreateSocialClass(context.Background(), revenue.SocialClass{Name: "wage-labourers", PrimaryRevenueForm: "wages"})
	if err != nil {
		t.Fatalf("seed class: %v", err)
	}
	if _, err := mem.CreateClassIncomeSource(context.Background(), revenue.ClassIncomeSource{ClassID: string(saved.ID), RevenueForm: "wages", SurplusValueComponentBP: 0}); err != nil {
		t.Fatalf("seed source: %v", err)
	}
	h := New(mem, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/classes/"+string(saved.ID), nil)
	req.SetPathValue("id", string(saved.ID))
	rec := httptest.NewRecorder()
	h.GetSocialClass(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var got socialClassDetailResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if string(got.Class.ID) != string(saved.ID) {
		t.Errorf("class id = %q, want %q", got.Class.ID, saved.ID)
	}
	if len(got.IncomeSources) != 1 {
		t.Errorf("income sources = %d, want 1", len(got.IncomeSources))
	}
}

func TestGetSocialClassNotFound(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/classes/nope", nil)
	req.SetPathValue("id", "nope")
	rec := httptest.NewRecorder()
	h.GetSocialClass(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestListClassAmbiguitiesEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListClassAmbiguities(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/class-ambiguities", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Fatalf("status=%d body=%s, want 200 items:[]", rec.Code, rec.Body.String())
	}
}
