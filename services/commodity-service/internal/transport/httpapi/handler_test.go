package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/store"
)

// newTestServer wires a Memory store + Handler and returns an *httptest.Server
// pointed at it. We register routes on a stdlib mux directly to avoid pulling
// the full pkg/httpx server (with its signal-handling + shutdown loop) into
// the test process.
func newTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	mem := store.NewMemory()
	h := New(mem, mem, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/commodities", h.Create)
	mux.HandleFunc("GET /v1/commodities", h.List)
	mux.HandleFunc("GET /v1/commodities/{id}", h.Get)
	mux.HandleFunc("PATCH /v1/commodities/{id}", h.Update)
	mux.HandleFunc("DELETE /v1/commodities/{id}", h.Delete)
	mux.HandleFunc("POST /v1/commodities/{id}/value", h.Value)
	mux.HandleFunc("GET /v1/commodities/{id}/value-form", h.ValueForm)
	mux.HandleFunc("GET /v1/commodities/{id}/social-relations", h.SocialRelations)
	mux.HandleFunc("POST /v1/exchange-ratio", h.ExchangeRatio)

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, mem
}

func mustCreate(t *testing.T, s *store.Memory, name, kind string, snlt commodity.LabourMinutes) commodity.Commodity {
	t.Helper()
	c, err := s.Create(context.Background(), commodity.Commodity{
		Name:           name,
		UseValue:       commodity.UseValue{Description: name + " for use", Unit: "units"},
		ConcreteLabour: commodity.ConcreteLabour{Kind: kind},
		SNLTPerUnit:    snlt,
	})
	if err != nil {
		t.Fatalf("seed %s: %v", name, err)
	}
	return c
}

func TestHTTP_CreateAndGet(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t)

	body := `{"name":"linen","use_value":{"description":"linen for clothing","unit":"yards"},"concrete_labour":{"kind":"weaving"},"snlt_per_unit":30}`
	res, err := http.Post(ts.URL+"/v1/commodities", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var created commodity.Commodity
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	res.Body.Close()
	if created.ID.IsZero() {
		t.Fatal("created.ID is zero")
	}

	res2, err := http.Get(ts.URL + "/v1/commodities/" + string(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	if res2.StatusCode != http.StatusOK {
		t.Errorf("get status = %d, want 200", res2.StatusCode)
	}
	res2.Body.Close()
}

func TestHTTP_Get_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t)
	res, err := http.Get(ts.URL + "/v1/commodities/missing")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
	res.Body.Close()
}

func TestHTTP_Value(t *testing.T) {
	t.Parallel()
	ts, mem := newTestServer(t)
	linen := mustCreate(t, mem, "linen", "weaving", 30)

	body := `{"quantity": 20}`
	res, err := http.Post(ts.URL+"/v1/commodities/"+string(linen.ID)+"/value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got valueResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if got.Value != 600 {
		t.Errorf("value = %d, want 600", got.Value)
	}
	if got.ValueHours != 10 {
		t.Errorf("value_hours = %v, want 10", got.ValueHours)
	}
}

func TestHTTP_ExchangeRatio(t *testing.T) {
	t.Parallel()
	ts, mem := newTestServer(t)
	linen := mustCreate(t, mem, "linen", "weaving", 30)
	coat := mustCreate(t, mem, "coat", "tailoring", 600)

	body, _ := json.Marshal(map[string]any{
		"base_id": linen.ID, "quote_id": coat.ID, "base_qty": 20,
	})
	res, err := http.Post(ts.URL+"/v1/exchange-ratio", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var ratio commodity.ExchangeRatio
	if err := json.NewDecoder(res.Body).Decode(&ratio); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if ratio.QuoteQty != 1 {
		t.Errorf("20 yards linen should = 1 coat; got %v", ratio.QuoteQty)
	}
}

func TestHTTP_SocialRelations(t *testing.T) {
	t.Parallel()
	ts, mem := newTestServer(t)
	linen := mustCreate(t, mem, "linen", "weaving", 30)
	mustCreate(t, mem, "coat", "tailoring", 600)
	mustCreate(t, mem, "iron", "smelting", 120)

	res, err := http.Get(ts.URL + "/v1/commodities/" + string(linen.ID) + "/social-relations")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var rel commodity.SocialRelations
	if err := json.NewDecoder(res.Body).Decode(&rel); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if rel.Subject.Name != "linen" {
		t.Errorf("subject.name = %v", rel.Subject.Name)
	}
	if len(rel.LabourRelations) != 2 {
		t.Errorf("relations = %d, want 2", len(rel.LabourRelations))
	}
}

func TestHTTP_ValueForm_Money(t *testing.T) {
	t.Parallel()
	ts, mem := newTestServer(t)
	linen := mustCreate(t, mem, "linen", "weaving", 30)
	gold := mustCreate(t, mem, "gold", "mining", 60)
	mustCreate(t, mem, "coat", "tailoring", 600)

	res, err := http.Get(ts.URL + "/v1/commodities/" + string(linen.ID) + "/value-form?kind=money&money_id=" + string(gold.ID))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	res.Body.Close()
}
