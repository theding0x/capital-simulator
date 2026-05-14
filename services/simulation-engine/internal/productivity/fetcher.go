// Package productivity holds the fetcher implementation that resolves a
// Part IV productivity factor by source + id. For Ch. 13 (cooperation) and
// Ch. 14 (manufacture) it GETs the chapter response from agent-service and
// reads the `productivity_factor` field. For Ch. 15 (factory) it queries
// the local machinery store directly. The bridge endpoint in
// simulation-engine calls this to translate Marx's Ch. 13/14/15 mechanisms
// into Ch. 12's relative-surplus formula.
package productivity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// Fetcher resolves cooperation / manufacture / factory productivity
// factors. Cooperation and Manufacture are GET'd from agent-service;
// Factory is read from the local store.
type Fetcher struct {
	AgentBaseURL string
	Factories    store.FactoryStore
	Client       *http.Client
}

// New returns a fetcher with sensible defaults.
func New(agentURL string, factories store.FactoryStore) *Fetcher {
	return &Fetcher{
		AgentBaseURL: agentURL,
		Factories:    factories,
		Client:       &http.Client{Timeout: 5 * time.Second},
	}
}

// Fetch implements httpapi.ProductivityFetcher.
func (f *Fetcher) Fetch(ctx context.Context, source httpapi.ProductivitySource, id string) (float64, error) {
	switch source {
	case httpapi.SourceCooperation:
		return f.fetchAgentFactor(ctx, "/v1/cooperations/"+id, "cooperation")
	case httpapi.SourceManufacture:
		return f.fetchAgentFactor(ctx, "/v1/manufactures/"+id, "manufacture")
	case httpapi.SourceFactory:
		return f.fetchFactoryFactor(ctx, id)
	default:
		return 0, httpapi.ErrUnknownSource
	}
}

func (f *Fetcher) fetchAgentFactor(ctx context.Context, path, label string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.AgentBaseURL+path, nil)
	if err != nil {
		return 0, err
	}
	resp, err := f.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch %s: %w", label, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, errors.New("productivity: source record not found")
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("fetch %s: status %d: %s", label, resp.StatusCode, string(body))
	}
	var payload struct {
		ProductivityFactor float64 `json:"productivity_factor"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("fetch %s: decode: %w", label, err)
	}
	if payload.ProductivityFactor <= 0 {
		return 0, errors.New("productivity: source returned a non-positive factor")
	}
	return payload.ProductivityFactor, nil
}

func (f *Fetcher) fetchFactoryFactor(ctx context.Context, id string) (float64, error) {
	if f.Factories == nil {
		return 0, errors.New("productivity: factory store not configured")
	}
	factory, err := f.Factories.GetFactory(ctx, machinery.FactoryID(id))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, errors.New("productivity: source record not found")
		}
		return 0, err
	}
	return machinery.FactoryProductivityFactor(factory), nil
}
