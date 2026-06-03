package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
	"github.com/theding0x/capital-simulator/services/market-service/internal/transport/httpapi"
)

type stubNotifier struct {
	mu      sync.Mutex
	called  int
	lastEx  market.Exchange
	failErr error
}

func (s *stubNotifier) NotifyExchange(_ context.Context, e market.Exchange) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.called++
	s.lastEx = e
	return s.failErr
}

func postExchange(t *testing.T, h *httpapi.Handler) *httptest.ResponseRecorder {
	t.Helper()
	body := map[string]any{
		"giver_id":              "weaver",
		"receiver_id":           "tailor",
		"giver_commodity_id":    "linen",
		"giver_qty":             20.0,
		"receiver_commodity_id": "money",
		"receiver_qty":          1.0,
		"realised_value":        1200,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/exchanges", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateExchange(rr, req)
	return rr
}

// A recorded exchange notifies agent-service exactly once, with the created
// exchange (id assigned) — this closes the counterparty's circuit leg (#216).
func TestCreateExchange_NotifiesAgent(t *testing.T) {
	t.Parallel()
	h := httpapi.New(store.NewMemory(), nil)
	stub := &stubNotifier{}
	h.Notifier = stub

	rr := postExchange(t, h)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}
	if stub.called != 1 {
		t.Fatalf("notifier called %d times, want 1", stub.called)
	}
	if stub.lastEx.GiverID != "weaver" || stub.lastEx.ReceiverID != "tailor" {
		t.Errorf("notified exchange parties = %s/%s", stub.lastEx.GiverID, stub.lastEx.ReceiverID)
	}
	if stub.lastEx.ID.IsZero() {
		t.Error("notified exchange should carry the assigned id")
	}
}

// A notifier failure is best-effort: the exchange is still recorded (201).
func TestCreateExchange_NotifyFailureDoesNotFailExchange(t *testing.T) {
	t.Parallel()
	h := httpapi.New(store.NewMemory(), nil)
	h.Notifier = &stubNotifier{failErr: errors.New("agent down")}

	rr := postExchange(t, h)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 despite notify failure; body: %s", rr.Code, rr.Body.String())
	}
}
