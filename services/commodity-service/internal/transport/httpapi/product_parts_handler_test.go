package httpapi

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/store"
)

func newProductPartsServer(t *testing.T) *httptest.Server {
	t.Helper()
	mem := store.NewMemory()
	h := New(mem, mem, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/product-parts", h.ProductParts)
	mux.HandleFunc("POST /v1/working-day-revision", h.WorkingDayRevision)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func postRaw(t *testing.T, url, body string) *http.Response {
	t.Helper()
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	t.Cleanup(func() { _ = resp.Body.Close() })
	return resp
}

// Ch.9 §2: 20 lbs. of yarn = 24s. c + 3s. v + 3s. s splits 16 / 2 / 2 lbs.
func TestHTTP_ProductParts_Ch9Yarn(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	resp := postRaw(t, ts.URL+"/v1/product-parts",
		`{"constant":2880,"variable":360,"surplus":360,"total_product":20}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got struct {
		TotalProduct float64 `json:"total_product"`
		Parts        struct {
			Constant float64 `json:"constant"`
			Variable float64 `json:"variable"`
			Surplus  float64 `json:"surplus"`
		} `json:"parts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.TotalProduct != 20 {
		t.Errorf("total_product = %v, want 20", got.TotalProduct)
	}
	for name, c := range map[string][2]float64{
		"constant": {got.Parts.Constant, 16},
		"variable": {got.Parts.Variable, 2},
		"surplus":  {got.Parts.Surplus, 2},
	} {
		if math.Abs(c[0]-c[1]) > 1e-9 {
			t.Errorf("%s = %v lbs, want %v", name, c[0], c[1])
		}
	}
}

func TestHTTP_ProductParts_BadInput(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	for name, body := range map[string]string{
		"malformed JSON":    `{"constant":`,
		"empty account":     `{"constant":0,"variable":0,"surplus":0,"total_product":20}`,
		"negative product":  `{"constant":24,"variable":3,"surplus":3,"total_product":-1}`,
		"negative constant": `{"constant":-1,"variable":3,"surplus":3,"total_product":20}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if resp := postRaw(t, ts.URL+"/v1/product-parts", body); resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
		})
	}
}

// Ch.9 §3: Senior's 11½-hour day on a 5¾-hour necessary labour, lengthened to
// 13 hours, gives 7¼ hours of surplus labour and a rate of 126 2/23%.
func TestHTTP_WorkingDayRevision_Ch9SeniorThirteenHours(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	resp := postRaw(t, ts.URL+"/v1/working-day-revision", `{"necessary_labour":345,"working_day":780}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got commodity.RevisedWorkingDay
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.WorkingDay != 780 || got.NecessaryLabour != 345 || got.SurplusLabour != 435 {
		t.Errorf("got %+v, want day 780, necessary 345, surplus 435", got)
	}
	if got.RateBasisPoints != 12609 {
		t.Errorf("rate_basis_points = %d, want 12609 (126 2/23%%)", got.RateBasisPoints)
	}
	if math.Abs(float64(got.Rate)-29.0/23.0) > 1e-9 {
		t.Errorf("rate = %v, want 29/23", got.Rate)
	}
}

func TestHTTP_WorkingDayRevision_BadInput(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	for name, body := range map[string]string{
		"malformed JSON":       `{"working_day":`,
		"day below necessary":  `{"necessary_labour":345,"working_day":344}`,
		"no necessary labour":  `{"necessary_labour":0,"working_day":690}`,
		"negative working day": `{"necessary_labour":345,"working_day":-1}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if resp := postRaw(t, ts.URL+"/v1/working-day-revision", body); resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
		})
	}
}

// An omitted key decodes to zero and would give a plausible-looking answer
// (a missing surplus reads as "no surplus-value"), so every field is required.
func TestHTTP_ProductParts_MissingFieldRejected(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	for name, body := range map[string]string{
		"missing total_product": `{"constant":2880,"variable":360,"surplus":360}`,
		"null total_product":    `{"constant":2880,"variable":360,"surplus":360,"total_product":null}`,
		"missing surplus":       `{"constant":2880,"variable":360,"total_product":20}`,
		"missing constant":      `{"variable":360,"surplus":360,"total_product":20}`,
		"missing variable":      `{"constant":2880,"surplus":360,"total_product":20}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if resp := postRaw(t, ts.URL+"/v1/product-parts", body); resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
		})
	}
}

func TestHTTP_WorkingDayRevision_MissingFieldRejected(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	for name, body := range map[string]string{
		"missing working_day":      `{"necessary_labour":345}`,
		"missing necessary_labour": `{"working_day":690}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if resp := postRaw(t, ts.URL+"/v1/working-day-revision", body); resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
		})
	}
}

// A huge product must still produce a body: +Inf is not encodable and used to
// leave a 200 with an empty response.
func TestHTTP_ProductParts_HugeTotalStillEncodes(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	resp := postRaw(t, ts.URL+"/v1/product-parts",
		`{"constant":9223372036854775807,"variable":1,"surplus":1,"total_product":1e300}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got productPartsResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("response body did not decode: %v", err)
	}
}

// The working-day error must name the necessary labour, not variable capital.
func TestHTTP_WorkingDayRevision_NoNecessaryLabourMessage(t *testing.T) {
	t.Parallel()
	ts := newProductPartsServer(t)
	resp := postRaw(t, ts.URL+"/v1/working-day-revision", `{"necessary_labour":0,"working_day":690}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !strings.Contains(body.Error, "necessary_labour") || strings.Contains(body.Error, "variable capital") {
		t.Errorf("error = %q, want it to name necessary_labour, not variable capital", body.Error)
	}
}

// Guest safety: neither handler may reach a store. They run here with a nil
// handler store, so any store access would panic.
func TestHTTP_ProductPartsAndRevision_TouchNoStore(t *testing.T) {
	t.Parallel()
	h := New(nil, nil, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/product-parts", h.ProductParts)
	mux.HandleFunc("POST /v1/working-day-revision", h.WorkingDayRevision)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	if resp := postRaw(t, ts.URL+"/v1/product-parts",
		`{"constant":2880,"variable":360,"surplus":360,"total_product":20}`); resp.StatusCode != http.StatusOK {
		t.Errorf("product-parts status = %d, want 200", resp.StatusCode)
	}
	if resp := postRaw(t, ts.URL+"/v1/working-day-revision",
		`{"necessary_labour":345,"working_day":690}`); resp.StatusCode != http.StatusOK {
		t.Errorf("working-day-revision status = %d, want 200", resp.StatusCode)
	}
}
