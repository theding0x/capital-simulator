package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
)

func newSurplusValueTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h := New(nil, Deps{})
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/surplus-value/absolute", h.ComputeAbsoluteSurplusValue)
	mux.HandleFunc("POST /v1/surplus-value/relative", h.ComputeRelativeSurplusValue)
	mux.HandleFunc("GET /v1/surplus-value/rate", h.GetSurplusValueRate)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// Ch. 10 §1 absolute fixture: prolonging working-day II (9 h = 6 h necessary +
// 3 h surplus) by 3 h gives working-day III (12 h = 6 + 6). Necessary labour is
// fixed and the whole extension lands in surplus labour. The step itself is
// derived: it uses Marx's two days, but he performs no such prolongation.
func TestComputeAbsoluteSurplusValue_Para1(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	body := `{"working_day_minutes":540,"necessary_labour_minutes":360,"extension_minutes":180}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/absolute", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got absoluteSurplusValueResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if want := (workingDayDTO{Total: 540, NecessaryLabour: 360, SurplusLabour: 180}); got.OldWorkingDay != want {
		t.Fatalf("old day = %+v, want %+v", got.OldWorkingDay, want)
	}
	if want := (workingDayDTO{Total: 720, NecessaryLabour: 360, SurplusLabour: 360}); got.NewWorkingDay != want {
		t.Fatalf("new day = %+v, want %+v (necessary labour fixed)", got.NewWorkingDay, want)
	}
	if got.AbsoluteSurplusValue.LabourMinutes != 180 {
		t.Fatalf("ASV = %d, want 180", got.AbsoluteSurplusValue.LabourMinutes)
	}
	if got.AbsoluteSurplusValue.Origin != "absolute" {
		t.Fatalf("Origin = %q, want absolute", got.AbsoluteSurplusValue.Origin)
	}
}

func TestComputeAbsoluteSurplusValue_RejectsNonPositiveExtension(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":600,"extension_minutes":0}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/absolute", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

// Ch. 12 relative fixture: the 12 h day (10 h necessary + 2 h surplus) with
// labour-power falling from 5s to 3s. Ch. 12 gives no social productivity
// factor for these figures; 5/3 is derived from his before/after figures
// (10 h → 6 h).
func TestComputeRelativeSurplusValue_Para1(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	body := fmt.Sprintf(`{"working_day_minutes":720,"necessary_labour_minutes":600,"productivity_factor":%v}`, 5.0/3.0)
	res, err := http.Post(ts.URL+"/v1/surplus-value/relative", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got relativeSurplusValueResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if want := (workingDayDTO{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}); got.OldWorkingDay != want {
		t.Fatalf("old day = %+v, want %+v", got.OldWorkingDay, want)
	}
	if got.NewWorkingDay.NecessaryLabour != 360 {
		t.Fatalf("new NecessaryLabour = %d, want 360", got.NewWorkingDay.NecessaryLabour)
	}
	if got.NewWorkingDay.SurplusLabour != 360 {
		t.Fatalf("new SurplusLabour = %d, want 360", got.NewWorkingDay.SurplusLabour)
	}
	if got.RelativeSurplusValue.LabourMinutes != 240 {
		t.Fatalf("RSV = %d, want 240", got.RelativeSurplusValue.LabourMinutes)
	}
	if got.RelativeSurplusValue.Origin != "relative" {
		t.Fatalf("Origin = %q, want relative", got.RelativeSurplusValue.Origin)
	}
}

// Ch. 16 Mill critique: rate of profit (4%) differs from rate of surplus-value
// (20%) for £400 c + £100 v + £20 s. The endpoint takes labour-minutes, so the
// pounds cross at the MELT (£1 = 2,400): s £20 = 48,000, v £100 = 240,000 and
// c+v £500 = 1,200,000.
func TestGetSurplusValueRate_Para1_MillCritique(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	c, v, s := money.FromLSD(400, 0, 0), money.FromLSD(100, 0, 0), money.FromLSD(20, 0, 0)
	url := fmt.Sprintf("%s/v1/surplus-value/rate?surplus_labour=%d&necessary_labour=%d&total_capital=%d",
		ts.URL, s.ToLabourMinutes(), v.ToLabourMinutes(), (c + v).ToLabourMinutes())
	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got surplusValueRateResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.SurplusLabourMinutes != 48000 || got.NecessaryLabourMinutes != 240000 || got.TotalCapitalAdvanced != 1200000 {
		t.Fatalf("echoed magnitudes = %d / %d / %d, want 48000 / 240000 / 1200000",
			got.SurplusLabourMinutes, got.NecessaryLabourMinutes, got.TotalCapitalAdvanced)
	}
	if got.SurplusValueMinutes != 48000 {
		t.Fatalf("SurplusValueMinutes = %d, want 48000 (defaults from surplus_labour)", got.SurplusValueMinutes)
	}
	if got.RateOfSurplusValue != 0.20 {
		t.Fatalf("RateOfSurplusValue = %v, want 0.20", got.RateOfSurplusValue)
	}
	if got.RateOfProfit != 0.04 {
		t.Fatalf("RateOfProfit = %v, want 0.04", got.RateOfProfit)
	}
	if !got.MillCritiqueHolds {
		t.Fatal("MillCritiqueHolds must be true when total_capital > necessary_labour")
	}
}

func TestGetSurplusValueRate_MissingParam(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	res, err := http.Get(ts.URL + "/v1/surplus-value/rate?surplus_labour=20&necessary_labour=100")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}
