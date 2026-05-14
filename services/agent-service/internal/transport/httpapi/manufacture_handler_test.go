package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Capital Vol. I, Ch. 14, §2 (Caslon type-foundry): "four founders ...
// two breakers ... one rubber". Each role rate × headcount = 8,000/hr.
// The manufacture's productive_power_factor must exceed the simple-
// cooperation baseline (i.e. > 1.0 by construction).
func TestCreateManufacture_TypeFoundryProportionality(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Henry Caslon")
	body := `{
		"name": "Caslon Type-Foundry",
		"capitalist_id": "` + string(cap) + `",
		"form": "serial",
		"origin": "splitting",
		"individual_working_day_minutes": 720,
		"roles": [
			{"name":"founder","skill_level":"skilled","output_rate_per_hour":2000,"head_count":4,"tool_name":"mould"},
			{"name":"breaker","skill_level":"unskilled","output_rate_per_hour":4000,"head_count":2,"tool_name":"iron"},
			{"name":"rubber","skill_level":"unskilled","output_rate_per_hour":8000,"head_count":1,"tool_name":"stone"}
		]
	}`
	res, err := http.Post(ts.URL+"/v1/manufactures", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var got manufactureResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.TotalWorkers != 7 {
		t.Fatalf("TotalWorkers = %d, want 7", got.TotalWorkers)
	}
	if got.ProductivityFactor <= 1.0 {
		t.Fatalf("ProductivityFactor = %v, want > 1.0", got.ProductivityFactor)
	}
	if got.ProductivePowerFactor != got.ProductivityFactor {
		t.Fatalf("ProductivePowerFactor (%v) and ProductivityFactor (%v) should match",
			got.ProductivePowerFactor, got.ProductivityFactor)
	}
	if got.CollectiveLabourer.IsParalysedIfAbsentOne != true {
		t.Fatalf("a serial manufacture must paralyse when any role is missing (§3)")
	}
}

func TestCreateManufacture_RejectsInvalidForm(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Matthew Boulton")
	body := `{
		"name": "x",
		"capitalist_id": "` + string(cap) + `",
		"form": "galaxy_brain",
		"origin": "splitting",
		"individual_working_day_minutes": 720,
		"roles": [{"name":"a","skill_level":"skilled","output_rate_per_hour":1,"head_count":1,"tool_name":""}]
	}`
	res, err := http.Post(ts.URL+"/v1/manufactures", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

func TestProportionalGroupSize_IronLaw(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Henry Caslon")
	createBody := `{
		"name": "Caslon Type-Foundry",
		"capitalist_id": "` + string(cap) + `",
		"form": "serial",
		"origin": "splitting",
		"individual_working_day_minutes": 720,
		"roles": [
			{"name":"founder","skill_level":"skilled","output_rate_per_hour":2000,"head_count":4,"tool_name":"mould"},
			{"name":"breaker","skill_level":"unskilled","output_rate_per_hour":4000,"head_count":2,"tool_name":"iron"},
			{"name":"rubber","skill_level":"unskilled","output_rate_per_hour":8000,"head_count":1,"tool_name":"stone"}
		]
	}`
	createRes, err := http.Post(ts.URL+"/v1/manufactures", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("POST create: %v", err)
	}
	var created manufactureResponse
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	createRes.Body.Close()

	// Target 16,000 units/hr — doubling the canonical fixture.
	probeBody := `{"target_output_rate": 16000}`
	probeRes, err := http.Post(
		ts.URL+"/v1/manufactures/"+string(created.ID)+"/proportional-group-size",
		"application/json", strings.NewReader(probeBody),
	)
	if err != nil {
		t.Fatalf("POST probe: %v", err)
	}
	defer probeRes.Body.Close()
	if probeRes.StatusCode != http.StatusOK {
		t.Fatalf("probe status = %d, want 200", probeRes.StatusCode)
	}
	var got proportionalGroupSizeResponse
	if err := json.NewDecoder(probeRes.Body).Decode(&got); err != nil {
		t.Fatalf("decode probe: %v", err)
	}
	// 16,000 / 2,000 = 8 founders, 16,000 / 4,000 = 4 breakers, 16,000 / 8,000 = 2 rubbers
	if got.Headcounts["founder"] != 8 {
		t.Fatalf("founders = %d, want 8", got.Headcounts["founder"])
	}
	if got.Headcounts["breaker"] != 4 {
		t.Fatalf("breakers = %d, want 4", got.Headcounts["breaker"])
	}
	if got.Headcounts["rubber"] != 2 {
		t.Fatalf("rubbers = %d, want 2", got.Headcounts["rubber"])
	}
}

func TestScaleManufacture_IntegerMultiple(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Henry Caslon")
	body := `{
		"name": "Caslon Type-Foundry",
		"capitalist_id": "` + string(cap) + `",
		"form": "serial",
		"origin": "splitting",
		"individual_working_day_minutes": 720,
		"roles": [
			{"name":"founder","skill_level":"skilled","output_rate_per_hour":2000,"head_count":4,"tool_name":"mould"},
			{"name":"breaker","skill_level":"unskilled","output_rate_per_hour":4000,"head_count":2,"tool_name":"iron"},
			{"name":"rubber","skill_level":"unskilled","output_rate_per_hour":8000,"head_count":1,"tool_name":"stone"}
		]
	}`
	cr, err := http.Post(ts.URL+"/v1/manufactures", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST create: %v", err)
	}
	var created manufactureResponse
	json.NewDecoder(cr.Body).Decode(&created) //nolint
	cr.Body.Close()

	sr, err := http.Post(
		ts.URL+"/v1/manufactures/"+string(created.ID)+"/scale",
		"application/json", bytes.NewReader([]byte(`{"multiplier":3}`)),
	)
	if err != nil {
		t.Fatalf("POST scale: %v", err)
	}
	defer sr.Body.Close()
	if sr.StatusCode != http.StatusOK {
		t.Fatalf("scale status = %d, want 200", sr.StatusCode)
	}
	var scaled manufactureResponse
	if err := json.NewDecoder(sr.Body).Decode(&scaled); err != nil {
		t.Fatalf("decode scale: %v", err)
	}
	if scaled.TotalWorkers != 21 { // 7 × 3
		t.Fatalf("scaled total = %d, want 21", scaled.TotalWorkers)
	}
}

func TestListAndGetCollectiveLabourer_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Henry Caslon")
	body := `{
		"name": "Caslon",
		"capitalist_id": "` + string(cap) + `",
		"form": "serial",
		"origin": "splitting",
		"individual_working_day_minutes": 720,
		"roles": [
			{"name":"founder","skill_level":"skilled","output_rate_per_hour":2000,"head_count":4,"tool_name":"mould"},
			{"name":"breaker","skill_level":"unskilled","output_rate_per_hour":4000,"head_count":2,"tool_name":"iron"},
			{"name":"rubber","skill_level":"unskilled","output_rate_per_hour":8000,"head_count":1,"tool_name":"stone"}
		]
	}`
	cr, err := http.Post(ts.URL+"/v1/manufactures", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST create: %v", err)
	}
	var created manufactureResponse
	json.NewDecoder(cr.Body).Decode(&created) //nolint
	cr.Body.Close()

	// GET /v1/manufactures (list)
	listRes, err := http.Get(ts.URL + "/v1/manufactures")
	if err != nil {
		t.Fatalf("GET list: %v", err)
	}
	defer listRes.Body.Close()
	if listRes.StatusCode != http.StatusOK {
		t.Fatalf("list status = %d, want 200", listRes.StatusCode)
	}

	// GET /v1/manufactures/{id}
	getRes, err := http.Get(ts.URL + "/v1/manufactures/" + string(created.ID))
	if err != nil {
		t.Fatalf("GET one: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}

	// GET /v1/manufactures/{id}/collective-labourer
	clRes, err := http.Get(ts.URL + "/v1/manufactures/" + string(created.ID) + "/collective-labourer")
	if err != nil {
		t.Fatalf("GET collective: %v", err)
	}
	defer clRes.Body.Close()
	if clRes.StatusCode != http.StatusOK {
		t.Fatalf("collective status = %d, want 200", clRes.StatusCode)
	}
	var cl collectiveLabourerResponse
	if err := json.NewDecoder(clRes.Body).Decode(&cl); err != nil {
		t.Fatalf("decode collective: %v", err)
	}
	if cl.TotalWorkers != 7 {
		t.Fatalf("TotalWorkers = %d, want 7", cl.TotalWorkers)
	}
}

func TestGetManufactureMinimumCapital_ExceedsCooperationBaseline(t *testing.T) {
	t.Parallel()
	ts, st := newAgentTestServer(t)
	cap := seedCapitalist(t, st, "Henry Caslon")
	body := `{
		"name": "Caslon Type-Foundry",
		"capitalist_id": "` + string(cap) + `",
		"form": "serial",
		"origin": "splitting",
		"individual_working_day_minutes": 720,
		"roles": [
			{"name":"founder","skill_level":"skilled","output_rate_per_hour":2000,"head_count":4,"tool_name":"mould"},
			{"name":"breaker","skill_level":"unskilled","output_rate_per_hour":4000,"head_count":2,"tool_name":"iron"},
			{"name":"rubber","skill_level":"unskilled","output_rate_per_hour":8000,"head_count":1,"tool_name":"stone"}
		]
	}`
	cr, err := http.Post(ts.URL+"/v1/manufactures", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST create: %v", err)
	}
	var created manufactureResponse
	json.NewDecoder(cr.Body).Decode(&created) //nolint
	cr.Body.Close()

	// raw_material_cost_factor 0.1 ≈ a tenth of a penny per LabourMinute
	// of collective output. With Caslon's role mix (3 unskilled to 4
	// skilled) the wage bill is below the skilled-rate baseline used by
	// CooperationBaseline; the §5 invariant only re-asserts itself once
	// raw materials flow at a non-trivial rate.
	mc, err := http.Get(ts.URL + "/v1/manufactures/" + string(created.ID) + "/minimum-capital?raw_material_cost_factor=0.1")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer mc.Body.Close()
	if mc.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", mc.StatusCode)
	}
	var got minimumCapitalManufactureResponse
	if err := json.NewDecoder(mc.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.MinimumCapitalPence <= got.CooperationBaseline {
		t.Fatalf("manufacture minimum (%d) must exceed cooperation baseline (%d) — §5",
			got.MinimumCapitalPence, got.CooperationBaseline)
	}
	if got.MinimumCapitalPence <= agent.Pence(0) {
		t.Fatal("MinimumCapitalPence must be positive")
	}
}
