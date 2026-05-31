// Package store defines the persistence boundary for finance-service.
// The Store interface is the only seam between the HTTP/domain layer and
// the underlying database (MySQL in production, in-memory in tests).
//
// finance-service is the home for Capital Vol. III — profit, the average
// rate of profit, commercial capital, interest-bearing capital, credit,
// rent, and the revenue distribution. Each Vol. III chapter PR adds the
// methods it needs to this interface and provides the matching Memory +
// MySQL implementations.
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// Sentinel errors that callers can branch on with errors.Is.
var (
	ErrNotFound      = errors.New("finance: not found")
	ErrAlreadyExists = errors.New("finance: already exists")
)

// Store is the persistence contract for finance-service.
//
// Vol. III Ch. 1 ("Cost-Price and Profit") opens the interface with cost-price
// persistence. Later chapters append their own methods here and the matching
// implementations in memory.go + mysql.go, following the per-chapter pattern
// documented in CLAUDE.md.
type Store interface {
	// CreateCostPrice persists a computed cost-price, assigning an ID and
	// created-at timestamp when absent. It returns ErrAlreadyExists if the ID
	// collides.
	CreateCostPrice(ctx context.Context, cp profit.CostPrice) (profit.CostPrice, error)
	// GetCostPrice returns the cost-price with the given ID, or ErrNotFound.
	GetCostPrice(ctx context.Context, id profit.CostPriceID) (profit.CostPrice, error)
	// ListCostPrices returns all stored cost-prices, newest first.
	ListCostPrices(ctx context.Context) ([]profit.CostPrice, error)

	// CreateProfitRate persists a rate-of-profit analysis (Vol. III Ch. 2),
	// assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateProfitRate(ctx context.Context, a profit.ProfitRateAnalysis) (profit.ProfitRateAnalysis, error)
	// GetProfitRate returns the analysis with the given ID, or ErrNotFound.
	GetProfitRate(ctx context.Context, id profit.ProfitRateID) (profit.ProfitRateAnalysis, error)
	// ListProfitRates returns all stored analyses, newest first.
	ListProfitRates(ctx context.Context) ([]profit.ProfitRateAnalysis, error)

	// CreateVariation persists a variation analysis (Vol. III Ch. 3) — the
	// movement of the rate of profit between two decompositions — assigning an
	// ID and created-at timestamp when absent. It returns ErrAlreadyExists if
	// the ID collides.
	CreateVariation(ctx context.Context, a profit.VariationAnalysis) (profit.VariationAnalysis, error)
	// GetVariation returns the variation analysis with the given ID, or ErrNotFound.
	GetVariation(ctx context.Context, id profit.VariationAnalysisID) (profit.VariationAnalysis, error)
	// ListVariations returns all stored variation analyses, newest first.
	ListVariations(ctx context.Context) ([]profit.VariationAnalysis, error)

	// CreateTurnoverAnalysis persists a turnover analysis (Vol. III Ch. 4) — the
	// annual rate of profit p' = s'·n·(v/C) corrected for the number of turnovers
	// — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateTurnoverAnalysis(ctx context.Context, a profit.TurnoverAnalysis) (profit.TurnoverAnalysis, error)
	// GetTurnoverAnalysis returns the turnover analysis with the given ID, or ErrNotFound.
	GetTurnoverAnalysis(ctx context.Context, id profit.TurnoverAnalysisID) (profit.TurnoverAnalysis, error)
	// ListTurnoverAnalyses returns all stored turnover analyses, newest first.
	ListTurnoverAnalyses(ctx context.Context) ([]profit.TurnoverAnalysis, error)

	// CreateEconomyAnalysis persists an economy analysis (Vol. III Ch. 5) — a
	// saving in constant capital and the rise it produces in the rate of profit
	// s/(c+v) — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateEconomyAnalysis(ctx context.Context, a profit.EconomyAnalysis) (profit.EconomyAnalysis, error)
	// GetEconomyAnalysis returns the economy analysis with the given ID, or ErrNotFound.
	GetEconomyAnalysis(ctx context.Context, id profit.EconomyAnalysisID) (profit.EconomyAnalysis, error)
	// ListEconomyAnalyses returns all stored economy analyses, newest first.
	ListEconomyAnalyses(ctx context.Context) ([]profit.EconomyAnalysis, error)

	// CreatePriceFluctuationAnalysis persists a price-fluctuation analysis (Vol.
	// III Ch. 6) — the movement in the rate of profit produced by a re-pricing of
	// the raw-material element of constant capital — assigning an ID and
	// created-at timestamp when absent. It returns ErrAlreadyExists if the ID
	// collides.
	CreatePriceFluctuationAnalysis(ctx context.Context, a profit.PriceFluctuationAnalysis) (profit.PriceFluctuationAnalysis, error)
	// GetPriceFluctuationAnalysis returns the analysis with the given ID, or ErrNotFound.
	GetPriceFluctuationAnalysis(ctx context.Context, id profit.PriceFluctuationAnalysisID) (profit.PriceFluctuationAnalysis, error)
	// ListPriceFluctuationAnalyses returns all stored analyses, newest first.
	ListPriceFluctuationAnalyses(ctx context.Context) ([]profit.PriceFluctuationAnalysis, error)

	// CreateCompositionEffect persists an organic-composition comparison (Vol. III
	// Ch. 7) — two capitals of equal s and v but different c, and the gap in their
	// rates of profit — assigning an ID and created-at timestamp when absent. It
	// returns ErrAlreadyExists if the ID collides.
	CreateCompositionEffect(ctx context.Context, a profit.CompositionEffectAnalysis) (profit.CompositionEffectAnalysis, error)
	// GetCompositionEffect returns the comparison with the given ID, or ErrNotFound.
	GetCompositionEffect(ctx context.Context, id profit.CompositionEffectID) (profit.CompositionEffectAnalysis, error)
	// ListCompositionEffects returns all stored comparisons, newest first.
	ListCompositionEffects(ctx context.Context) ([]profit.CompositionEffectAnalysis, error)

	// CreateMagnitudeChange persists a capital-magnitude change (Vol. III Ch. 7) —
	// a change in the size of capital and what becomes of its rate of profit —
	// assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateMagnitudeChange(ctx context.Context, a profit.MagnitudeChangeAnalysis) (profit.MagnitudeChangeAnalysis, error)
	// GetMagnitudeChange returns the change with the given ID, or ErrNotFound.
	GetMagnitudeChange(ctx context.Context, id profit.MagnitudeChangeID) (profit.MagnitudeChangeAnalysis, error)
	// ListMagnitudeChanges returns all stored changes, newest first.
	ListMagnitudeChanges(ctx context.Context) ([]profit.MagnitudeChangeAnalysis, error)

	// CreateProductionSphere persists a production-sphere analysis (Vol. III Ch. 8)
	// — a branch of production with its organic composition and individual rate of
	// profit — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateProductionSphere(ctx context.Context, s avgprofit.ProductionSphere) (avgprofit.ProductionSphere, error)
	// GetProductionSphere returns the sphere with the given ID, or ErrNotFound.
	GetProductionSphere(ctx context.Context, id avgprofit.ProductionSphereID) (avgprofit.ProductionSphere, error)
	// ListProductionSpheres returns all stored spheres, newest first.
	ListProductionSpheres(ctx context.Context) ([]avgprofit.ProductionSphere, error)

	// CreateGeneralProfitRate persists a general-rate computation (Vol. III Ch. 9)
	// — the capital-weighted general rate over a set of production spheres —
	// assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateGeneralProfitRate(ctx context.Context, g avgprofit.GeneralProfitRate) (avgprofit.GeneralProfitRate, error)
	// GetGeneralProfitRate returns the general rate with the given ID, or ErrNotFound.
	GetGeneralProfitRate(ctx context.Context, id avgprofit.GeneralProfitRateID) (avgprofit.GeneralProfitRate, error)

	// CreatePriceOfProduction persists a price-of-production record (Vol. III Ch. 9)
	// — a sphere's cost-price, general rate, average profit, price, and deviation
	// from value — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreatePriceOfProduction(ctx context.Context, p avgprofit.PriceOfProduction) (avgprofit.PriceOfProduction, error)
	// GetPriceOfProduction returns the record with the given ID, or ErrNotFound.
	GetPriceOfProduction(ctx context.Context, id avgprofit.PriceOfProductionID) (avgprofit.PriceOfProduction, error)
	// ListPricesOfProduction returns all stored records, newest first.
	ListPricesOfProduction(ctx context.Context) ([]avgprofit.PriceOfProduction, error)

	// CreateMarketValue persists a market-value record (Vol. III Ch. 10) —
	// the value set by bulk production conditions in a branch — assigning an ID
	// and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateMarketValue(ctx context.Context, v avgprofit.MarketValue) (avgprofit.MarketValue, error)
	// GetMarketValue returns the record with the given ID, or ErrNotFound.
	GetMarketValue(ctx context.Context, id avgprofit.MarketValueID) (avgprofit.MarketValue, error)

	// CreateSurplusProfit persists a surplus-profit record (Vol. III Ch. 10) —
	// the signed extra profit of a below-market-value firm — assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateSurplusProfit(ctx context.Context, s avgprofit.SurplusProfit) (avgprofit.SurplusProfit, error)
	// GetSurplusProfit returns the record with the given ID, or ErrNotFound.
	GetSurplusProfit(ctx context.Context, id avgprofit.SurplusProfitID) (avgprofit.SurplusProfit, error)

	// CreateEqualisation persists an equalisation record (Vol. III Ch. 10) —
	// a sphere's competitive convergence toward the general rate — assigning an ID
	// and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateEqualisation(ctx context.Context, e avgprofit.Equalisation) (avgprofit.Equalisation, error)
	// GetEqualisation returns the record with the given ID, or ErrNotFound.
	GetEqualisation(ctx context.Context, id avgprofit.EqualisationID) (avgprofit.Equalisation, error)

	// CreateWageEffectAnalysis persists a wage-effect analysis (Vol. III Ch. 11) —
	// the shift in prices of production across spheres caused by a general wage
	// fluctuation — assigning an ID and created-at timestamp when absent. Returns
	// ErrAlreadyExists on ID collision.
	CreateWageEffectAnalysis(ctx context.Context, a avgprofit.WageEffectAnalysis) (avgprofit.WageEffectAnalysis, error)
	// GetWageEffectAnalysis returns the analysis with the given ID, or ErrNotFound.
	GetWageEffectAnalysis(ctx context.Context, id avgprofit.WageEffectAnalysisID) (avgprofit.WageEffectAnalysis, error)
	// ListWageEffectAnalyses returns all stored analyses, newest first.
	ListWageEffectAnalyses(ctx context.Context) ([]avgprofit.WageEffectAnalysis, error)

	// CreatePriceOfProductionChange persists a price-of-production change (Vol. III
	// Ch. 12) — a movement in a sphere's price of production and which of the two
	// admissible causes produced it — assigning an ID and created-at timestamp when
	// absent. Returns ErrAlreadyExists on ID collision.
	CreatePriceOfProductionChange(ctx context.Context, c avgprofit.PriceOfProductionChange) (avgprofit.PriceOfProductionChange, error)
	// GetPriceOfProductionChange returns the change with the given ID, or ErrNotFound.
	GetPriceOfProductionChange(ctx context.Context, id avgprofit.PriceOfProductionChangeID) (avgprofit.PriceOfProductionChange, error)

	// CreateCompositionTrajectory persists a composition trajectory (Vol. III Ch. 13)
	// — the time-series of rising c/v with constant s′ and the falling rate of profit
	// it produces — assigning an ID and created-at timestamp when absent. The derived
	// ProfitRates slice is recomputed from Periods on store. Returns ErrAlreadyExists
	// on ID collision.
	CreateCompositionTrajectory(ctx context.Context, t tendency.CompositionTrajectory) (tendency.CompositionTrajectory, error)
	// GetCompositionTrajectory returns the trajectory with the given ID, or ErrNotFound.
	GetCompositionTrajectory(ctx context.Context, id tendency.CompositionTrajectoryID) (tendency.CompositionTrajectory, error)
	// ListCompositionTrajectories returns all stored trajectories, newest first.
	ListCompositionTrajectories(ctx context.Context) ([]tendency.CompositionTrajectory, error)

	// CreateRateMassContradiction persists a rate-mass contradiction (Vol. III Ch. 13)
	// — the falling rate of profit against the (possibly growing) mass of profit —
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists
	// on ID collision.
	CreateRateMassContradiction(ctx context.Context, r tendency.RateMassContradiction) (tendency.RateMassContradiction, error)
	// GetRateMassContradiction returns the record with the given ID, or ErrNotFound.
	GetRateMassContradiction(ctx context.Context, id tendency.RateMassContradictionID) (tendency.RateMassContradiction, error)
	// ListRateMassContradictions returns all stored records, newest first.
	ListRateMassContradictions(ctx context.Context) ([]tendency.RateMassContradiction, error)

	// CreateCounteractingForce persists a counteracting force (Vol. III Ch. 14),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists
	// on ID collision.
	CreateCounteractingForce(ctx context.Context, f tendency.CounteractingForce) (tendency.CounteractingForce, error)
	// GetCounteractingForce returns the force with the given ID, or ErrNotFound.
	GetCounteractingForce(ctx context.Context, id tendency.CounteractingForceID) (tendency.CounteractingForce, error)

	// CreateCounteractingScenario persists a counteracting scenario (Vol. III Ch. 14),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists
	// on ID collision.
	CreateCounteractingScenario(ctx context.Context, s tendency.CounteractingScenario) (tendency.CounteractingScenario, error)
	// GetCounteractingScenario returns the scenario with the given ID, or ErrNotFound.
	GetCounteractingScenario(ctx context.Context, id tendency.CounteractingScenarioID) (tendency.CounteractingScenario, error)
	// ListCounteractingScenarios returns all stored scenarios, newest first.
	ListCounteractingScenarios(ctx context.Context) ([]tendency.CounteractingScenario, error)

	// CreateCrisis persists a crisis record (Vol. III Ch. 15), assigning an ID
	// and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCrisis(ctx context.Context, c tendency.Crisis) (tendency.Crisis, error)
	// GetCrisis returns the crisis with the given ID, or ErrNotFound.
	GetCrisis(ctx context.Context, id tendency.CrisisID) (tendency.Crisis, error)
	// ListCrises returns all stored crises, newest first.
	ListCrises(ctx context.Context) ([]tendency.Crisis, error)

	// CreateInternalContradiction persists an internal contradiction (Vol. III Ch. 15),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateInternalContradiction(ctx context.Context, c tendency.InternalContradiction) (tendency.InternalContradiction, error)
	// GetInternalContradiction returns the contradiction with the given ID, or ErrNotFound.
	GetInternalContradiction(ctx context.Context, id tendency.InternalContradictionID) (tendency.InternalContradiction, error)
	// ListInternalContradictions returns all stored contradictions, newest first.
	ListInternalContradictions(ctx context.Context) ([]tendency.InternalContradiction, error)

	// CreateCommercialCapital persists a commercial-capital record (Vol. III Ch. 16 — opens Part IV),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCommercialCapital(ctx context.Context, cc merchant.CommercialCapital) (merchant.CommercialCapital, error)
	// GetCommercialCapital returns the commercial-capital record with the given ID, or ErrNotFound.
	GetCommercialCapital(ctx context.Context, id merchant.CommercialCapitalID) (merchant.CommercialCapital, error)
	// ListCommercialCapitals returns all stored commercial-capital records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCommercialCapitals(ctx context.Context) ([]merchant.CommercialCapital, error)

	// CreateCommercialProfit persists a commercial-profit record (Vol. III Ch. 17 — the adjusted
	// general rate of profit once commercial capital is included in the denominator), assigning an
	// ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCommercialProfit(ctx context.Context, cp merchant.CommercialProfit) (merchant.CommercialProfit, error)
	// GetCommercialProfit returns the commercial-profit record with the given ID, or ErrNotFound.
	GetCommercialProfit(ctx context.Context, id merchant.CommercialProfitID) (merchant.CommercialProfit, error)
	// ListCommercialProfits returns all stored commercial-profit records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCommercialProfits(ctx context.Context) ([]merchant.CommercialProfit, error)

	// CreateMerchantTurnover persists a merchant-turnover record (Vol. III Ch. 18 — annual profit
	// invariant to turnover; per-unit markup = profit / (turnover × units)), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateMerchantTurnover(ctx context.Context, mt merchant.MerchantTurnover) (merchant.MerchantTurnover, error)
	// GetMerchantTurnover returns the merchant-turnover record with the given ID, or ErrNotFound.
	GetMerchantTurnover(ctx context.Context, id merchant.MerchantTurnoverID) (merchant.MerchantTurnover, error)
	// ListMerchantTurnovers returns all stored merchant-turnover records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListMerchantTurnovers(ctx context.Context) ([]merchant.MerchantTurnover, error)

	// CreateMoneyDealingCapital persists a money-dealing-capital record (Vol. III Ch. 19 —
	// technical monetary operations: receipt/payment/exchange/safekeeping/bookkeeping;
	// MoneyDealingProfit = roundHalfUp(money_advanced×rate); NewValueCreated == 0), assigning an
	// ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateMoneyDealingCapital(ctx context.Context, m merchant.MoneyDealingCapital) (merchant.MoneyDealingCapital, error)
	// GetMoneyDealingCapital returns the money-dealing-capital record with the given ID, or ErrNotFound.
	GetMoneyDealingCapital(ctx context.Context, id merchant.MoneyDealingCapitalID) (merchant.MoneyDealingCapital, error)
	// ListMoneyDealingCapitals returns all stored money-dealing-capital records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListMoneyDealingCapitals(ctx context.Context) ([]merchant.MoneyDealingCapital, error)

	// CreateHistoricalMerchantCapital persists a historical-merchant-capital record (Vol. III Ch. 20 —
	// historical forms from Venice/Genoa/Dutch carrying trade through subordination under industrial
	// capital; DevelopmentStage 1–3; SubordinationIndex bp [0,10000]; COMPLETES Part IV), assigning
	// an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateHistoricalMerchantCapital(ctx context.Context, m merchant.HistoricalMerchantCapital) (merchant.HistoricalMerchantCapital, error)
	// GetHistoricalMerchantCapital returns the historical-merchant-capital record with the given ID, or ErrNotFound.
	GetHistoricalMerchantCapital(ctx context.Context, id merchant.HistoricalMerchantCapitalID) (merchant.HistoricalMerchantCapital, error)
	// ListHistoricalMerchantCapitals returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListHistoricalMerchantCapitals(ctx context.Context) ([]merchant.HistoricalMerchantCapital, error)

	// CreateInterestBearingCapital persists an interest-bearing-capital record (Vol. III Ch. 21 —
	// opens Part V; M—M′ circuit; InterestEarned=roundHalfUp(money_advanced×rate_bp, 10000);
	// MoneyReturned=MoneyAdvanced+InterestEarned; NewValueCreated==0), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateInterestBearingCapital(ctx context.Context, ibc credit.InterestBearingCapital) (credit.InterestBearingCapital, error)
	// GetInterestBearingCapital returns the interest-bearing-capital record with the given ID, or ErrNotFound.
	GetInterestBearingCapital(ctx context.Context, id credit.InterestBearingCapitalID) (credit.InterestBearingCapital, error)
	// ListInterestBearingCapitals returns all stored interest-bearing-capital records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListInterestBearingCapitals(ctx context.Context) ([]credit.InterestBearingCapital, error)

	// CreateRateOfInterest persists a rate-of-interest observation (Vol. III Ch. 22 —
	// observed market interest rate in bp, average profit rate in bp, industrial cycle phase;
	// RateBP>=0; RateBP<AverageProfitRateBP), assigning an ID and created-at timestamp when
	// absent. Returns ErrAlreadyExists on ID collision.
	CreateRateOfInterest(ctx context.Context, r credit.RateOfInterest) (credit.RateOfInterest, error)
	// GetRateOfInterest returns the rate-of-interest record with the given ID, or ErrNotFound.
	GetRateOfInterest(ctx context.Context, id credit.RateOfInterestID) (credit.RateOfInterest, error)
	// ListRatesOfInterest returns all stored rate-of-interest records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListRatesOfInterest(ctx context.Context) ([]credit.RateOfInterest, error)

	// CreateProfitDivision persists a profit-division record (Vol. III Ch. 23 —
	// splits total profit into interest + profit of enterprise; TotalProfitBP>0;
	// InterestBP>=0; InterestBP<TotalProfitBP; ProfitOfEnterpriseBP=Total-Interest),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateProfitDivision(ctx context.Context, pd credit.ProfitDivision) (credit.ProfitDivision, error)
	// GetProfitDivision returns the profit-division record with the given ID, or ErrNotFound.
	GetProfitDivision(ctx context.Context, id credit.ProfitDivisionID) (credit.ProfitDivision, error)
	// ListProfitDivisions returns all stored profit-division records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListProfitDivisions(ctx context.Context) ([]credit.ProfitDivision, error)

	// CreateCompoundInterestSchedule persists a compound-interest schedule (Vol. III Ch. 24 —
	// integer compound iteration; FinalValue derived; NewValueCreated==0; Principal/RateBP/PeriodYears>0),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCompoundInterestSchedule(ctx context.Context, s credit.CompoundInterestSchedule) (credit.CompoundInterestSchedule, error)
	// GetCompoundInterestSchedule returns the schedule with the given ID, or ErrNotFound.
	GetCompoundInterestSchedule(ctx context.Context, id credit.CompoundInterestID) (credit.CompoundInterestSchedule, error)
	// ListCompoundInterestSchedules returns all stored schedules, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCompoundInterestSchedules(ctx context.Context) ([]credit.CompoundInterestSchedule, error)

	// CreateBillOfExchange persists a bill of exchange (Vol. III Ch. 25 —
	// commercial credit; FaceValue>0; MaturityDays>0), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateBillOfExchange(ctx context.Context, b credit.BillOfExchange) (credit.BillOfExchange, error)
	// GetBillOfExchange returns the bill with the given ID, or ErrNotFound.
	GetBillOfExchange(ctx context.Context, id credit.BillOfExchangeID) (credit.BillOfExchange, error)
	// ListBillsOfExchange returns all stored bills, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListBillsOfExchange(ctx context.Context) ([]credit.BillOfExchange, error)

	// CreateFictitiousCapital persists a fictitious-capital record (Vol. III Ch. 25 —
	// capitalised income stream; AnnualIncome>0; RateBP>0; CapitalisedValue derived;
	// KindPublicDebt has RealCapitalExists==false), assigning an ID and created-at
	// timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateFictitiousCapital(ctx context.Context, fc credit.FictitiousCapital) (credit.FictitiousCapital, error)
	// GetFictitiousCapital returns the record with the given ID, or ErrNotFound.
	GetFictitiousCapital(ctx context.Context, id credit.FictitiousCapitalID) (credit.FictitiousCapital, error)
	// ListFictitiousCapitals returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListFictitiousCapitals(ctx context.Context) ([]credit.FictitiousCapital, error)

	// CreateMoneyCapitalAccumulation persists a money-capital-accumulation record
	// (Vol. III Ch. 26 — accumulation of money-capital and its influence on the
	// interest rate across industrial cycle phases; InterestRateBP>=0;
	// IdleCapital.Amount>=0; LoanableSupply.Amount>=0), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateMoneyCapitalAccumulation(ctx context.Context, a credit.MoneyCapitalAccumulation) (credit.MoneyCapitalAccumulation, error)
	// GetMoneyCapitalAccumulation returns the record with the given ID, or ErrNotFound.
	GetMoneyCapitalAccumulation(ctx context.Context, id credit.MoneyCapitalAccumulationID) (credit.MoneyCapitalAccumulation, error)
	// ListMoneyCapitalAccumulations returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListMoneyCapitalAccumulations(ctx context.Context) ([]credit.MoneyCapitalAccumulation, error)

	// CreateStockCompany persists a joint-stock company (Vol. III Ch. 27 —
	// TotalCapital==FixedCapital+FloatingCapital; Manager.WageMinutes>0), assigning
	// an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateStockCompany(ctx context.Context, sc credit.StockCompany) (credit.StockCompany, error)
	// GetStockCompany returns the stock company with the given ID, or ErrNotFound.
	GetStockCompany(ctx context.Context, id credit.StockCompanyID) (credit.StockCompany, error)
	// ListStockCompanies returns all stored stock companies, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListStockCompanies(ctx context.Context) ([]credit.StockCompany, error)

	// CreateCooperativeFactory persists a worker-owned cooperative factory
	// (Vol. III Ch. 27 — WorkerOwned==true; positive transitional form), assigning
	// an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCooperativeFactory(ctx context.Context, cf credit.CooperativeFactory) (credit.CooperativeFactory, error)
	// GetCooperativeFactory returns the cooperative factory with the given ID, or ErrNotFound.
	GetCooperativeFactory(ctx context.Context, id credit.CooperativeFactoryID) (credit.CooperativeFactory, error)
	// ListCooperativeFactories returns all stored cooperative factories, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCooperativeFactories(ctx context.Context) ([]credit.CooperativeFactory, error)

	// CreateCurrencyObservation persists a currency-observation record
	// (Vol. III Ch. 28 — medium of circulation split into coin function and
	// capital-transfer function; TotalCurrencyOutstanding>0;
	// CoinFunctionAmount+CapitalTransferAmount<=Total; ReserveFund.Amount>=0),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCurrencyObservation(ctx context.Context, o credit.CurrencyObservation) (credit.CurrencyObservation, error)
	// GetCurrencyObservation returns the record with the given ID, or ErrNotFound.
	GetCurrencyObservation(ctx context.Context, id credit.CurrencyObservationID) (credit.CurrencyObservation, error)
	// ListCurrencyObservations returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCurrencyObservations(ctx context.Context) ([]credit.CurrencyObservation, error)

	// CreateBankCapital persists a bank-capital decomposition (Vol. III Ch. 29 —
	// cash + securities; TotalCapital==CashAmount+SecuritiesAmount; Components is a
	// list of cash/bill/bond/stock/mortgage parts), assigning an ID and created-at
	// timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateBankCapital(ctx context.Context, bc credit.BankCapital) (credit.BankCapital, error)
	// GetBankCapital returns the bank-capital record with the given ID, or ErrNotFound.
	GetBankCapital(ctx context.Context, id credit.BankCapitalID) (credit.BankCapital, error)
	// ListBankCapitals returns all stored bank-capital records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListBankCapitals(ctx context.Context) ([]credit.BankCapital, error)

	// CreateFictitiousCapitalValuation persists a security valuation (Vol. III Ch. 29 —
	// capitalised income; MarketValue==roundHalfUp(AnnualIncome*10000, InterestRateBP);
	// InterestRateBP>0; market value moves inversely with the rate), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateFictitiousCapitalValuation(ctx context.Context, v credit.FictitiousCapitalValuation) (credit.FictitiousCapitalValuation, error)
	// GetFictitiousCapitalValuation returns the valuation with the given ID, or ErrNotFound.
	GetFictitiousCapitalValuation(ctx context.Context, id credit.FictitiousCapitalValuationID) (credit.FictitiousCapitalValuation, error)
	// ListFictitiousCapitalValuations returns all stored valuations, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListFictitiousCapitalValuations(ctx context.Context) ([]credit.FictitiousCapitalValuation, error)

	// CreateRealCapitalAccumulation persists a real-capital-accumulation observation
	// (Vol. III Ch. 30 — whether money-capital accumulation corresponds to real-capital
	// accumulation across the industrial cycle; Phase.IsValid(); CreditLimit.ReserveCapital>=0;
	// the two growth rates may diverge), assigning an ID and created-at timestamp when absent.
	// Returns ErrAlreadyExists on ID collision.
	CreateRealCapitalAccumulation(ctx context.Context, a credit.RealCapitalAccumulation) (credit.RealCapitalAccumulation, error)
	// GetRealCapitalAccumulation returns the record with the given ID, or ErrNotFound.
	GetRealCapitalAccumulation(ctx context.Context, id credit.RealCapitalAccumulationID) (credit.RealCapitalAccumulation, error)
	// ListRealCapitalAccumulations returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListRealCapitalAccumulations(ctx context.Context) ([]credit.RealCapitalAccumulation, error)

	// CreateFloatingCapital persists a floating-capital record (Vol. III Ch. 31 —
	// Money-Capital and Real Capital, II; Amount>=0; Source ∈ {1,2}; embedded
	// LoanCapitalVelocity with TimesLent>=1 and LoanCapitalCreated==MoneyAmount×TimesLent;
	// embedded CirculationReserveSaving with Amount>=0), assigning an ID and created-at
	// timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateFloatingCapital(ctx context.Context, f credit.FloatingCapital) (credit.FloatingCapital, error)
	// GetFloatingCapital returns the record with the given ID, or ErrNotFound.
	GetFloatingCapital(ctx context.Context, id credit.FloatingCapitalID) (credit.FloatingCapital, error)
	// ListFloatingCapitals returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListFloatingCapitals(ctx context.Context) ([]credit.FloatingCapital, error)

	// CreateCapitalRelease persists a capital-release record (Vol. III Ch. 32 —
	// Money-Capital and Real Capital, III (Conclusion); ReleasedAmount>=0;
	// Cause ∈ {1,2,3,4}; embedded LoanCapitalOverstatement with OverstatementPct>=0;
	// a release reflects real accumulation only when Reinvested), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCapitalRelease(ctx context.Context, c credit.CapitalRelease) (credit.CapitalRelease, error)
	// GetCapitalRelease returns the record with the given ID, or ErrNotFound.
	GetCapitalRelease(ctx context.Context, id credit.CapitalReleaseID) (credit.CapitalRelease, error)
	// ListCapitalReleases returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCapitalReleases(ctx context.Context) ([]credit.CapitalRelease, error)

	// CreateClearingHouseSettlement persists a clearing-house settlement record
	// (Vol. III Ch. 33 — The Medium of Circulation in the Credit System;
	// MoneyUsed<=TotalClaims; NetBalance==TotalClaims-MoneyUsed), assigning an ID
	// and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateClearingHouseSettlement(ctx context.Context, s credit.ClearingHouseSettlement) (credit.ClearingHouseSettlement, error)
	// GetClearingHouseSettlement returns the record with the given ID, or ErrNotFound.
	GetClearingHouseSettlement(ctx context.Context, id credit.ClearingHouseSettlementID) (credit.ClearingHouseSettlement, error)
	// ListClearingHouseSettlements returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListClearingHouseSettlements(ctx context.Context) ([]credit.ClearingHouseSettlement, error)

	// CreateNoteIssueConstraint persists a note-issue-constraint record
	// (Vol. III Ch. 34 — The Currency Principle and the Bank Acts of 1844 and 1845;
	// Regime.IsValid(); MaxNotes==GoldBacking+UnbackedLimit), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateNoteIssueConstraint(ctx context.Context, n credit.NoteIssueConstraint) (credit.NoteIssueConstraint, error)
	// GetNoteIssueConstraint returns the record with the given ID, or ErrNotFound.
	GetNoteIssueConstraint(ctx context.Context, id credit.NoteIssueConstraintID) (credit.NoteIssueConstraint, error)
	// ListNoteIssueConstraints returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListNoteIssueConstraints(ctx context.Context) ([]credit.NoteIssueConstraint, error)

	// CreateGoldReserve persists a gold-reserve record (Vol. III Ch. 35 —
	// Precious Metal and Rate of Exchange; Amount>=0), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateGoldReserve(ctx context.Context, g credit.GoldReserve) (credit.GoldReserve, error)
	// GetGoldReserve returns the record with the given ID, or ErrNotFound.
	GetGoldReserve(ctx context.Context, id credit.GoldReserveID) (credit.GoldReserve, error)
	// ListGoldReserves returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListGoldReserves(ctx context.Context) ([]credit.GoldReserve, error)

	// CreateRateOfExchange persists a rate-of-exchange record (Vol. III Ch. 35 —
	// deviation from mint par in basis points; sign unrestricted), assigning an ID and
	// created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateRateOfExchange(ctx context.Context, r credit.RateOfExchange) (credit.RateOfExchange, error)
	// GetRateOfExchange returns the record with the given ID, or ErrNotFound.
	GetRateOfExchange(ctx context.Context, id credit.RateOfExchangeID) (credit.RateOfExchange, error)
	// ListRatesOfExchange returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListRatesOfExchange(ctx context.Context) ([]credit.RateOfExchange, error)

	// CreateUsurersCapital persists a usurers-capital record (Vol. III Ch. 36 —
	// Pre-Capitalist Relationships; Stage.IsValid(); WageLabour==false for Antiquity/Mediaeval;
	// InterestRateBP>=0; StageSubordinated requires SubordinatedTo=="industrial_capital"),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateUsurersCapital(ctx context.Context, u credit.UsurersCapital) (credit.UsurersCapital, error)
	// GetUsurersCapital returns the record with the given ID, or ErrNotFound.
	GetUsurersCapital(ctx context.Context, id credit.UsurersCapitalID) (credit.UsurersCapital, error)
	// ListUsurersCapitals returns all stored records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListUsurersCapitals(ctx context.Context) ([]credit.UsurersCapital, error)

	// CreateLandParcel persists a land-parcel record (Vol. III Ch. 37 — opens Part VI),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateLandParcel(ctx context.Context, p rent.LandParcel) (rent.LandParcel, error)
	// GetLandParcel returns the land-parcel with the given ID, or ErrNotFound.
	GetLandParcel(ctx context.Context, id rent.LandParcelID) (rent.LandParcel, error)
	// ListLandParcels returns all stored land-parcels, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListLandParcels(ctx context.Context) ([]rent.LandParcel, error)

	// CreateLandowner persists a landowner record (Vol. III Ch. 37),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateLandowner(ctx context.Context, lo rent.Landowner) (rent.Landowner, error)
	// GetLandowner returns the landowner with the given ID, or ErrNotFound.
	GetLandowner(ctx context.Context, id rent.LandownerID) (rent.Landowner, error)

	// CreateAgriculturalCapitalist persists an agricultural-capitalist record (Vol. III Ch. 37),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateAgriculturalCapitalist(ctx context.Context, ac rent.AgriculturalCapitalist) (rent.AgriculturalCapitalist, error)
	// GetAgriculturalCapitalist returns the agricultural capitalist with the given ID, or ErrNotFound.
	GetAgriculturalCapitalist(ctx context.Context, id rent.AgriculturalCapitalistID) (rent.AgriculturalCapitalist, error)

	// CreateGroundRent persists a ground-rent record (Vol. III Ch. 37),
	// assigning an ID and created-at timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateGroundRent(ctx context.Context, g rent.GroundRent) (rent.GroundRent, error)
	// ListGroundRents returns all stored ground-rent records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListGroundRents(ctx context.Context) ([]rent.GroundRent, error)

	// CreatePoPSurplusProfit persists a price-of-production surplus-profit record
	// (Vol. III Ch. 38 — gap between general and individual price of production;
	// SurplusProfitBP derived server-side), assigning an ID and created-at timestamp
	// when absent. Returns ErrAlreadyExists on ID collision.
	CreatePoPSurplusProfit(ctx context.Context, s rent.PriceOfProductionSurplusProfit) (rent.PriceOfProductionSurplusProfit, error)
	// ListPoPSurplusProfits returns all stored PoP surplus-profit records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListPoPSurplusProfits(ctx context.Context) ([]rent.PriceOfProductionSurplusProfit, error)

	// CreateMonopolisedNaturalForce persists a monopolised-natural-force record
	// (Vol. III Ch. 38 — waterfall, mine, fertile soil whose monopolisation enables
	// differential rent), assigning an ID and created-at timestamp when absent.
	// Returns ErrAlreadyExists on ID collision.
	CreateMonopolisedNaturalForce(ctx context.Context, f rent.MonopolisedNaturalForce) (rent.MonopolisedNaturalForce, error)
	// ListMonopolisedNaturalForces returns all stored monopolised-natural-force records,
	// newest first (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListMonopolisedNaturalForces(ctx context.Context) ([]rent.MonopolisedNaturalForce, error)

	// CreateDifferentialRent persists a differential-rent record (Vol. III Ch. 38 —
	// surplus-profit converted into rent), assigning an ID and created-at timestamp
	// when absent. Returns ErrAlreadyExists on ID collision.
	CreateDifferentialRent(ctx context.Context, dr rent.DifferentialRent) (rent.DifferentialRent, error)
	// ListDifferentialRents returns all stored differential-rent records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListDifferentialRents(ctx context.Context) ([]rent.DifferentialRent, error)

	// CreateCapitalisedRentPrice persists a capitalised-rent-price record (Vol. III Ch. 38 —
	// annual rent capitalised at the interest rate gives the land purchase price;
	// CapitalisedPriceLabourMinutes derived server-side), assigning an ID and created-at
	// timestamp when absent. Returns ErrAlreadyExists on ID collision.
	CreateCapitalisedRentPrice(ctx context.Context, crp rent.CapitalisedRentPrice) (rent.CapitalisedRentPrice, error)
	// ListCapitalisedRentPrices returns all stored capitalised-rent-price records, newest first
	// (created_at DESC, id ASC). Never returns nil — an empty store returns an empty slice.
	ListCapitalisedRentPrices(ctx context.Context) ([]rent.CapitalisedRentPrice, error)
}
