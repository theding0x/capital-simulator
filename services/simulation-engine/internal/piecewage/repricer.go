// Package piecewage holds the HTTP client the PiecePriceTicker uses to push the
// current factory productivity factor into agent-service. agent-service reprices
// every piece-wage from its baseline and appends a price point when the factor
// changes — Marx's Vol. I Ch. 21 law that piece-prices fall as productivity
// rises, deferred to the tick scheduler by issue #219. It mirrors the sibling
// productivity package, which reaches the other way for relative-surplus inputs.
package piecewage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Repricer reprices every piece-wage in agent-service at a productivity factor.
type Repricer struct {
	AgentBaseURL string
	Client       *http.Client
}

// New returns a Repricer with a sensible default client.
func New(agentURL string) *Repricer {
	return &Repricer{
		AgentBaseURL: agentURL,
		Client:       &http.Client{Timeout: 5 * time.Second},
	}
}

type pieceWage struct {
	AgentID string `json:"agent_id"`
}

type repriceResult struct {
	Appended bool `json:"appended"`
}

// Reprice lists every piece-wage and reprices each at the given productivity
// factor, returning the number of piece-wages whose price actually changed this
// pass (a point was appended). It continues past a single agent's failure and
// returns the joined error so one bad agent does not starve the rest.
func (rp *Repricer) Reprice(ctx context.Context, productivityFactor float64) (int, error) {
	wages, err := rp.list(ctx)
	if err != nil {
		return 0, err
	}
	var advanced int
	var errs []error
	for _, wage := range wages {
		appended, err := rp.repriceOne(ctx, wage.AgentID, productivityFactor)
		if err != nil {
			errs = append(errs, fmt.Errorf("reprice %s: %w", wage.AgentID, err))
			continue
		}
		if appended {
			advanced++
		}
	}
	return advanced, errors.Join(errs...)
}

func (rp *Repricer) list(ctx context.Context) ([]pieceWage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rp.AgentBaseURL+"/v1/piece-wages", nil)
	if err != nil {
		return nil, err
	}
	resp, err := rp.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list piece-wages: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("list piece-wages: status %d: %s", resp.StatusCode, string(body))
	}
	var wages []pieceWage
	if err := json.NewDecoder(resp.Body).Decode(&wages); err != nil {
		return nil, fmt.Errorf("list piece-wages: decode: %w", err)
	}
	return wages, nil
}

func (rp *Repricer) repriceOne(ctx context.Context, agentID string, factor float64) (bool, error) {
	body, err := json.Marshal(map[string]float64{"productivity_factor": factor})
	if err != nil {
		return false, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		rp.AgentBaseURL+"/v1/piece-wages/"+agentID+"/reprice", bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := rp.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return false, fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}
	var result repriceResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("decode: %w", err)
	}
	return result.Appended, nil
}
