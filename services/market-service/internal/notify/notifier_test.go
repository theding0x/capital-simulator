package notify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// The notifier POSTs the settlement to agent-service's /v1/circuit-legs with
// the exchange's parties, commodity, and realised value intact.
func TestNotifyExchange_PostsSettlement(t *testing.T) {
	t.Parallel()
	var got settlementPayload
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	n := New(srv.URL)
	ex := market.Exchange{
		ID:                  "exch-1",
		GiverID:             "weaver",
		ReceiverID:          "tailor",
		GiverCommodityID:    "linen",
		ReceiverCommodityID: "money",
		RealisedValue:       1200,
	}
	if err := n.NotifyExchange(context.Background(), ex); err != nil {
		t.Fatalf("NotifyExchange: %v", err)
	}
	if path != "/v1/circuit-legs" {
		t.Errorf("path = %s, want /v1/circuit-legs", path)
	}
	if got.GiverID != "weaver" || got.ReceiverID != "tailor" {
		t.Errorf("parties = %s/%s, want weaver/tailor", got.GiverID, got.ReceiverID)
	}
	if got.GiverCommodityID != "linen" || got.RealisedValue != 1200 {
		t.Errorf("payload = %+v", got)
	}
}

func TestNotifyExchange_ErrorOnBadStatus(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	n := New(srv.URL)
	err := n.NotifyExchange(context.Background(), market.Exchange{
		ID: "e", GiverID: "a", ReceiverID: "b", GiverCommodityID: "c", RealisedValue: 1,
	})
	if err == nil {
		t.Fatal("expected error on 500 status")
	}
}
