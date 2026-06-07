package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/observatory"
)

// snapshotIntervalMS is the advisory client poll interval echoed in the snapshot.
// Advancement is driven by the client's `advance` query param, not the server.
const snapshotIntervalMS = 2000

// observatorySnapshotResponse is the GET /v1/observatory/snapshot body: the whole
// field of industrial capitals, the aggregate vital-signs, and the hidden abode,
// for one session's in-memory run. Consumed by the Atlas page.
type observatorySnapshotResponse struct {
	Tick         int64              `json:"tick"`
	Running      bool               `json:"running"`
	IntervalMS   int64              `json:"interval_ms"`
	Capitals     []fieldCapitalDTO  `json:"capitals"`
	Aggregate    aggregateVitalsDTO `json:"aggregate"`
	Abode        abodeDTO           `json:"abode"`
	Circulation  circulationDTO     `json:"circulation"`
	Distribution distributionDTO    `json:"distribution"`
}

type fieldCapitalDTO struct {
	ID              string `json:"id"`
	TotalPence      int64  `json:"total_pence"`
	MoneyPence      int64  `json:"money_pence"`
	ProductionPence int64  `json:"production_pence"`
	CommodityPence  int64  `json:"commodity_pence"`
	CostPricePence  int64  `json:"cost_price_pence"`
	SurplusPence    int64  `json:"surplus_pence"`
	Status          string `json:"status"`
	TurnoverNumber  int64  `json:"turnover_number"`
}

type aggregateVitalsDTO struct {
	TotalSocialCapitalPence int64 `json:"total_social_capital_pence"`
	CostPricePence          int64 `json:"cost_price_pence"`
	SurplusPence            int64 `json:"surplus_pence"`
	AvgRateOfProfitBP       int64 `json:"avg_rate_of_profit_bp"`
}

type abodeDTO struct {
	TotalVariablePence     int64                 `json:"total_variable_pence"`
	TotalSurplusPence      int64                 `json:"total_surplus_pence"`
	RateOfExploitationBP   int64                 `json:"rate_of_exploitation_bp"`
	NecessaryLabourMinutes int64                 `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64                 `json:"surplus_labour_minutes"`
	OrganicCompositionBP   int64                 `json:"organic_composition_bp"`
	ReserveArmyCount       int64                 `json:"reserve_army_count"`
	ReserveArmyPressureBP  int64                 `json:"reserve_army_pressure_bp"`
	EmployedCount          int64                 `json:"employed_count"`
	WagePence              int64                 `json:"wage_pence"`
	TotalPopulation        int64                 `json:"total_population"`
	SurplusRateBaseBP      int64                 `json:"surplus_rate_base_bp"`
	BaseWagePence          int64                 `json:"base_wage_pence"`
	AccumulationRateBP     int64                 `json:"accumulation_rate_bp"`
	LawSeries              []generalLawPeriodDTO `json:"law_series"`
}

type generalLawPeriodDTO struct {
	Period               int64 `json:"period"`
	WagePence            int64 `json:"wage_pence"`
	RateOfExploitationBP int64 `json:"rate_of_exploitation_bp"`
	ReserveArmyCount     int64 `json:"reserve_army_count"`
	OrganicCompositionBP int64 `json:"organic_composition_bp"`
}

type deptCVSDTO struct {
	C int64 `json:"c"`
	V int64 `json:"v"`
	S int64 `json:"s"`
}

type circulationDTO struct {
	Departments struct {
		I  deptCVSDTO `json:"I"`
		II deptCVSDTO `json:"II"`
	} `json:"departments"`
	KeystoneExchangePence int64 `json:"keystone_exchange_pence"`
	DeptIIConstantPence   int64 `json:"dept_ii_constant_pence"`
	AccumulationRateBP    int64 `json:"accumulation_rate_bp"`
	Extended              bool  `json:"extended"`
}

type distSeriesPointDTO struct {
	Period          int64 `json:"period"`
	PprimeBP        int64 `json:"pprime_bp"`
	ProfitMassPence int64 `json:"profit_mass_pence"`
	CompositionBP   int64 `json:"composition_bp"`
	Crisis          bool  `json:"crisis"`
}

type trinityDTO struct {
	WagesPence  int64 `json:"wages_pence"`
	ProfitPence int64 `json:"profit_pence"`
	RentPence   int64 `json:"rent_pence"`
}

type soilGradeDTO struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	FertilityBP           int64  `json:"fertility_bp"`
	YieldUnits            int64  `json:"yield_units"`
	IndividualPricePence  int64  `json:"individual_price_pence"`
	DifferentialRentPence int64  `json:"differential_rent_pence"`
	AbsoluteRentPence     int64  `json:"absolute_rent_pence"`
	TotalRentPence        int64  `json:"total_rent_pence"`
	LandPricePence        int64  `json:"land_price_pence"`
	Regulating            bool   `json:"regulating"`
}

type rentDTO struct {
	Grades               []soilGradeDTO `json:"grades"`
	WaterlinePence       int64          `json:"waterline_pence"`
	CapitalPerGradePence int64          `json:"capital_per_grade_pence"`
	DifferentialPence    int64          `json:"differential_pence"`
	AbsolutePence        int64          `json:"absolute_pence"`
	TotalRentPence       int64          `json:"total_rent_pence"`
	LandPricePence       int64          `json:"land_price_pence"`
	InterestRateBP       int64          `json:"interest_rate_bp"`
}

type distributionDTO struct {
	GeneralRateBP   int64                `json:"general_rate_bp"`
	InterestRateBP  int64                `json:"interest_rate_bp"`
	LoanablePence   int64                `json:"loanable_pence"`
	FictitiousPence int64                `json:"fictitious_pence"`
	DistSeries      []distSeriesPointDTO `json:"dist_series"`
	Rent            rentDTO              `json:"rent"`
	Trinity         trinityDTO           `json:"trinity"`
}

// GetObservatorySnapshot handles GET /v1/observatory/snapshot?advance=N. It reads
// the X-Atlas-Session header, advances that session's in-memory run by N periods
// (default 1), and returns the projection. No store I/O; nothing is persisted.
func (h *Handler) GetObservatorySnapshot(w http.ResponseWriter, r *http.Request) {
	if h.Observatory == nil {
		h.writeServerError(w, errors.New("observatory not configured"))
		return
	}
	advance := parseAdvance(r.URL.Query().Get("advance"))
	run := h.Observatory.GetOrCreate(r.Header.Get("X-Atlas-Session"))
	run.Advance(advance)
	writeJSON(w, http.StatusOK, buildSnapshotResponse(run.Snapshot(), advance))
}

// parseAdvance returns the requested advance count: empty or invalid → 1, "0" →
// 0 (paused). Run.Advance clamps the upper bound.
func parseAdvance(q string) int {
	if q == "" {
		return 1
	}
	n, err := strconv.Atoi(q)
	if err != nil || n < 0 {
		return 1
	}
	return n
}

// buildSnapshotResponse maps a RunSnapshot to the wire DTO. Slices are always
// non-nil so the client never sees `null`.
func buildSnapshotResponse(snap observatory.RunSnapshot, advance int) observatorySnapshotResponse {
	resp := observatorySnapshotResponse{
		Tick:       snap.Tick,
		Running:    advance > 0,
		IntervalMS: snapshotIntervalMS,
		Capitals:   make([]fieldCapitalDTO, len(snap.Field)),
	}
	var sumTotal, sumCost, sumSurplus int64
	for i, fc := range snap.Field {
		resp.Capitals[i] = fieldCapitalDTO{
			ID:              string(fc.ID),
			TotalPence:      int64(fc.TotalPence),
			MoneyPence:      int64(fc.MoneyPence),
			ProductionPence: int64(fc.ProductionPence),
			CommodityPence:  int64(fc.CommodityPence),
			CostPricePence:  int64(fc.CostPricePence),
			SurplusPence:    int64(fc.SurplusPence),
			Status:          string(fc.Status),
			TurnoverNumber:  fc.TurnoverNumber,
		}
		sumTotal += int64(fc.TotalPence)
		sumCost += int64(fc.CostPricePence)
		sumSurplus += int64(fc.SurplusPence)
	}
	resp.Aggregate = aggregateVitalsDTO{
		TotalSocialCapitalPence: sumTotal,
		CostPricePence:          sumCost,
		SurplusPence:            sumSurplus,
		AvgRateOfProfitBP:       observatory.RateBP(sumSurplus, sumCost),
	}
	ar := snap.Readout
	law := make([]generalLawPeriodDTO, len(snap.Periods))
	for i, p := range snap.Periods {
		law[i] = generalLawPeriodDTO{
			Period:               p.Period,
			WagePence:            p.WagePence,
			RateOfExploitationBP: p.RateOfExploitationBP,
			ReserveArmyCount:     p.ReserveArmyCount,
			OrganicCompositionBP: p.OrganicCompositionBP,
		}
	}
	resp.Abode = abodeDTO{
		TotalVariablePence:     ar.TotalVariablePence,
		TotalSurplusPence:      ar.TotalSurplusPence,
		RateOfExploitationBP:   ar.RateOfExploitationBP,
		NecessaryLabourMinutes: ar.NecessaryLabourMinutes,
		SurplusLabourMinutes:   ar.SurplusLabourMinutes,
		OrganicCompositionBP:   ar.OrganicCompositionBP,
		ReserveArmyCount:       ar.ReserveArmyCount,
		ReserveArmyPressureBP:  ar.ReserveArmyPressureBP,
		EmployedCount:          ar.EmployedCount,
		WagePence:              ar.WagePence,
		TotalPopulation:        ar.TotalPopulation,
		SurplusRateBaseBP:      snap.Abode.SurplusRateBaseBP,
		BaseWagePence:          snap.Abode.BaseWagePence,
		AccumulationRateBP:     snap.Abode.AccumulationRateBP,
		LawSeries:              law,
	}
	resp.Circulation = mapCirculationDTO(snap.Circulation)
	resp.Distribution = mapDistributionDTO(snap.Distribution)
	return resp
}

func mapCirculationDTO(c observatory.CirculationSnapshot) circulationDTO {
	dto := circulationDTO{
		KeystoneExchangePence: c.KeystoneExchangePence,
		DeptIIConstantPence:   c.DeptIIConstantPence,
		AccumulationRateBP:    c.AccumulationRateBP,
		Extended:              c.Extended,
	}
	dto.Departments.I = deptCVSDTO{C: c.Departments.I.C, V: c.Departments.I.V, S: c.Departments.I.S}
	dto.Departments.II = deptCVSDTO{C: c.Departments.II.C, V: c.Departments.II.V, S: c.Departments.II.S}
	return dto
}

func mapDistributionDTO(d observatory.DistributionSnapshot) distributionDTO {
	series := make([]distSeriesPointDTO, len(d.DistSeries))
	for i, p := range d.DistSeries {
		series[i] = distSeriesPointDTO{
			Period:          p.Period,
			PprimeBP:        p.PprimeBP,
			ProfitMassPence: p.ProfitMassPence,
			CompositionBP:   p.CompositionBP,
			Crisis:          p.Crisis,
		}
	}
	grades := make([]soilGradeDTO, len(d.Rent.Grades))
	for i, g := range d.Rent.Grades {
		grades[i] = soilGradeDTO{
			ID:                    g.ID,
			Name:                  g.Name,
			FertilityBP:           g.FertilityBP,
			YieldUnits:            g.YieldUnits,
			IndividualPricePence:  g.IndividualPricePence,
			DifferentialRentPence: g.DifferentialRentPence,
			AbsoluteRentPence:     g.AbsoluteRentPence,
			TotalRentPence:        g.TotalRentPence,
			LandPricePence:        g.LandPricePence,
			Regulating:            g.Regulating,
		}
	}

	return distributionDTO{
		GeneralRateBP:   d.GeneralRateBP,
		InterestRateBP:  d.InterestRateBP,
		LoanablePence:   d.LoanablePence,
		FictitiousPence: d.FictitiousPence,
		DistSeries:      series,
		Rent: rentDTO{
			Grades:               grades,
			WaterlinePence:       d.Rent.WaterlinePence,
			CapitalPerGradePence: d.Rent.CapitalPerGradePence,
			DifferentialPence:    d.Rent.DifferentialPence,
			AbsolutePence:        d.Rent.AbsolutePence,
			TotalRentPence:       d.Rent.TotalRentPence,
			LandPricePence:       d.Rent.LandPricePence,
			InterestRateBP:       d.Rent.InterestRateBP,
		},
		Trinity: trinityDTO{
			WagesPence:  d.Trinity.WagesPence,
			ProfitPence: d.Trinity.ProfitPence,
			RentPence:   d.Trinity.RentPence,
		},
	}
}
