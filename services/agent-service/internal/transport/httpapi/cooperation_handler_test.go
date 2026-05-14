package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

// newAgentTestServer wires the in-memory store into a test mux and returns
// the live URL. Reused by every handler test in this package.
func newAgentTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	// Cooperation routes.
	mux.HandleFunc("POST /v1/cooperations", h.CreateCooperation)
	mux.HandleFunc("GET /v1/cooperations", h.ListCooperations)
	mux.HandleFunc("GET /v1/cooperations/{id}", h.GetCooperation)
	mux.HandleFunc("POST /v1/cooperations/{id}/collective-working-day", h.ComputeCollectiveWorkingDay)
	mux.HandleFunc("POST /v1/cooperations/{id}/average-social-labour", h.ComputeAverageSocialLabour)
	mux.HandleFunc("POST /v1/cooperations/minimum-capital", h.ComputeMinimumCapital)
	// Manufacture routes.
	mux.HandleFunc("POST /v1/manufactures", h.CreateManufacture)
	mux.HandleFunc("GET /v1/manufactures", h.ListManufactures)
	mux.HandleFunc("GET /v1/manufactures/{id}", h.GetManufacture)
	mux.HandleFunc("POST /v1/manufactures/{id}/proportional-group-size", h.ProportionalGroupSize)
	mux.HandleFunc("POST /v1/manufactures/{id}/scale", h.ScaleManufacture)
	mux.HandleFunc("GET /v1/manufactures/{id}/collective-labourer", h.GetCollectiveLabourer)
	mux.HandleFunc("GET /v1/manufactures/{id}/minimum-capital", h.GetManufactureMinimumCapital)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// seedCapitalist creates a capitalist agent directly via the store; the
// cooperation/manufacture handlers expect an agent.AgentID payload.
func seedCapitalist(t *testing.T, st *store.Memory, name string) agent.AgentID {
	t.Helper()
	created, err := st.Create(context.Background(), agent.Agent{
		Name:         name,
		Class:        agent.CapitalistClass,
		MoneyBalance: agent.Pence(100_000),
	})
	if err != nil {
		t.Fatalf("seed capitalist: %v", err)
	}
	return agent.AgentID(created.ID)
}

// Capital Vol. I, Ch. 13, Burke fn. 1: "any given five men will, in their
// total, afford a proportion of labour equal to any other five." Building
// a 5-worker cooperation must surface the productivity_factor field that
// drives the Part IV bridge endpoint.
func TestCreateCooperation_BurkePlatoon(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Robert Owen")
	body := `{
		"name": "spinning floor",
		"capitalist_id": "` + string(cap) + `",
		"members": [
			{"worker_id":"w1","working_day_minutes":720},
			{"worker_id":"w2","working_day_minutes":720},
			{"worker_id":"w3","working_day_minutes":720},
			{"worker_id":"w4","working_day_minutes":720},
			{"worker_id":"w5","working_day_minutes":720,"supervisory":true}
		]
	}`
	res, err := http.Post(ts.URL+"/v1/cooperations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var got cooperationResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Size != 5 {
		t.Fatalf("Size = %d, want 5", got.Size)
	}
	if got.CollectiveWorkingDay != 3600 {
		t.Fatalf("CollectiveWorkingDay = %d, want 3600 (5 × 720)", got.CollectiveWorkingDay)
	}
	if got.AverageSocialLabour != 720 {
		t.Fatalf("AverageSocialLabour = %d, want 720", got.AverageSocialLabour)
	}
	if got.ProductivityFactor <= 1.0 {
		t.Fatalf("ProductivityFactor = %v, want > 1.0 (5-worker cooperation bonus)", got.ProductivityFactor)
	}
	if len(got.Supervisors) != 1 {
		t.Fatalf("Supervisors = %d, want 1", len(got.Supervisors))
	}
}

func TestCreateCooperation_RejectsNoMembers(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Manchester Mill Co.")
	body := `{"name":"empty","capitalist_id":"` + string(cap) + `","members":[]}`
	res, err := http.Post(ts.URL+"/v1/cooperations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

func TestComputeMinimumCapital_Para8_300vs10(t *testing.T) {
	t.Parallel()
	ts, _ := newAgentTestServer(t)
	// 300 × 72 = 21,600 pence (Marx's §8 reference).
	body := `{"worker_count":300,"daily_wage_pence":72}`
	res, err := http.Post(ts.URL+"/v1/cooperations/minimum-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got minimumCapitalResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.MinimumCapital != agent.Pence(21_600) {
		t.Fatalf("MinimumCapital = %d, want 21600", got.MinimumCapital)
	}
}

func TestComputeCollectiveWorkingDay_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Andrew Ure")
	body := `{
		"name": "12-men gang",
		"capitalist_id": "` + string(cap) + `",
		"members": [
			{"worker_id":"w01","working_day_minutes":720},{"worker_id":"w02","working_day_minutes":720},
			{"worker_id":"w03","working_day_minutes":720},{"worker_id":"w04","working_day_minutes":720},
			{"worker_id":"w05","working_day_minutes":720},{"worker_id":"w06","working_day_minutes":720},
			{"worker_id":"w07","working_day_minutes":720},{"worker_id":"w08","working_day_minutes":720},
			{"worker_id":"w09","working_day_minutes":720},{"worker_id":"w10","working_day_minutes":720},
			{"worker_id":"w11","working_day_minutes":720},{"worker_id":"w12","working_day_minutes":720}
		]
	}`
	res, err := http.Post(ts.URL+"/v1/cooperations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	var created cooperationResponse
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	res.Body.Close()

	cwdRes, err := http.Post(ts.URL+"/v1/cooperations/"+string(created.ID)+"/collective-working-day", "application/json", bytes.NewReader([]byte("{}")))
	if err != nil {
		t.Fatalf("POST collective: %v", err)
	}
	defer cwdRes.Body.Close()
	if cwdRes.StatusCode != http.StatusOK {
		t.Fatalf("collective status = %d, want 200", cwdRes.StatusCode)
	}
	var got collectiveWorkingDayResponse
	if err := json.NewDecoder(cwdRes.Body).Decode(&got); err != nil {
		t.Fatalf("decode collective: %v", err)
	}
	if got.CollectiveWorkingDayMinutes != 8640 {
		t.Fatalf("CollectiveWorkingDayMinutes = %d, want 8640 (12 × 720)", got.CollectiveWorkingDayMinutes)
	}
}

func TestGetCooperation_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newAgentTestServer(t)
	res, err := http.Get(ts.URL + "/v1/cooperations/no-such-id")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", res.StatusCode)
	}
}

// Round-trip coverage for the remaining cooperation routes.
func TestListAndAverageSocialLabour_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Andrew Ure")
	body := `{
		"name": "5-men gang",
		"capitalist_id": "` + string(cap) + `",
		"members": [
			{"worker_id":"w1","working_day_minutes":720},{"worker_id":"w2","working_day_minutes":720},
			{"worker_id":"w3","working_day_minutes":720},{"worker_id":"w4","working_day_minutes":720},
			{"worker_id":"w5","working_day_minutes":720}
		]
	}`
	cr, err := http.Post(ts.URL+"/v1/cooperations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	var created cooperationResponse
	json.NewDecoder(cr.Body).Decode(&created) //nolint
	cr.Body.Close()

	// GET /v1/cooperations
	listRes, err := http.Get(ts.URL + "/v1/cooperations")
	if err != nil {
		t.Fatalf("GET list: %v", err)
	}
	defer listRes.Body.Close()
	if listRes.StatusCode != http.StatusOK {
		t.Fatalf("list status = %d, want 200", listRes.StatusCode)
	}
	var list struct {
		Items []cooperationResponse `json:"items"`
	}
	if err := json.NewDecoder(listRes.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("list len = %d, want 1", len(list.Items))
	}

	// GET /v1/cooperations/{id}
	getRes, err := http.Get(ts.URL + "/v1/cooperations/" + string(created.ID))
	if err != nil {
		t.Fatalf("GET one: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}

	// POST /v1/cooperations/{id}/average-social-labour
	avgRes, err := http.Post(
		ts.URL+"/v1/cooperations/"+string(created.ID)+"/average-social-labour",
		"application/json", bytes.NewReader([]byte("{}")),
	)
	if err != nil {
		t.Fatalf("POST avg: %v", err)
	}
	defer avgRes.Body.Close()
	if avgRes.StatusCode != http.StatusOK {
		t.Fatalf("avg status = %d, want 200", avgRes.StatusCode)
	}
	var avg averageSocialLabourResponse
	if err := json.NewDecoder(avgRes.Body).Decode(&avg); err != nil {
		t.Fatalf("decode avg: %v", err)
	}
	if avg.AverageSocialLabourMinutes != 720 {
		t.Fatalf("AverageSocialLabour = %d, want 720", avg.AverageSocialLabourMinutes)
	}
	if !avg.MeetsBurkeMinimumPlatoon {
		t.Fatal("5 members must meet the Burke minimum platoon")
	}
}
