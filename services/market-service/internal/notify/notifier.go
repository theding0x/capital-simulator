// Package notify pushes completed exchanges to agent-service so the
// counterparty's circuit leg is closed (#216). It mirrors the simulation-engine
// productivity fetcher: a thin http.Client keyed off AGENT_SERVICE_URL,
// invoked best-effort from the exchange handler.
package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// AgentNotifier posts a completed Exchange to agent-service's settlement
// endpoint.
type AgentNotifier struct {
	AgentBaseURL string
	Client       *http.Client
}

// New returns a notifier with a sensible client timeout.
func New(agentURL string) *AgentNotifier {
	return &AgentNotifier{
		AgentBaseURL: agentURL,
		Client:       &http.Client{Timeout: 5 * time.Second},
	}
}

// settlementPayload is the JSON body agent-service's POST /v1/circuit-legs
// expects. The realised value is a labour-minute magnitude.
type settlementPayload struct {
	ExchangeID          string `json:"exchange_id"`
	GiverID             string `json:"giver_id"`
	ReceiverID          string `json:"receiver_id"`
	GiverCommodityID    string `json:"giver_commodity_id"`
	ReceiverCommodityID string `json:"receiver_commodity_id"`
	RealisedValue       int64  `json:"realised_value"`
}

// NotifyExchange posts the settlement to agent-service. Callers treat it as
// best-effort: a failure here must never fail the exchange that was already
// recorded.
func (n *AgentNotifier) NotifyExchange(ctx context.Context, e market.Exchange) error {
	payload := settlementPayload{
		ExchangeID:          string(e.ID),
		GiverID:             string(e.GiverID),
		ReceiverID:          string(e.ReceiverID),
		GiverCommodityID:    string(e.GiverCommodityID),
		ReceiverCommodityID: string(e.ReceiverCommodityID),
		RealisedValue:       int64(e.RealisedValue),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.AgentBaseURL+"/v1/circuit-legs", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.Client.Do(req)
	if err != nil {
		return fmt.Errorf("notify agent: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("notify agent: status %d", resp.StatusCode)
	}
	return nil
}
