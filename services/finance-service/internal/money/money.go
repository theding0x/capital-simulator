// Package money holds finance-service's canonical money-unit conversions: the
// single Pence↔LabourMinute scale every package in the service shares, so rent
// and average-profit magnitudes stay commensurable. It is the narrow, do-it-now
// slice of ADR-0001 (docs/adr/0001-value-money-composition.md), whose deferred
// repo-wide pkg/money will eventually own this conversion for every service
// (#411, #369). Keep this package free of domain logic — conversion constants
// only.
package money

// LabourMinutesPerPound is the value of £1 in labour-minutes: £1 = one working
// day = 8 hours = 480 labour-minutes. The rent sub-system (Ch. 37–46: the
// surplus-profit ≙ differential-rent identity and capitalised land-prices) is
// built on this scale; avgprofit's labour-power index now shares it instead of
// its former divergent 3600 (£1 = a 60-hour week), which disagreed by 7.5×.
const LabourMinutesPerPound = 480

// PencePerPound is the £→pence scale for money-magnitudes recorded in pence
// (e.g. the Ch. 38 surplus-profit): £1 = 100 pence.
const PencePerPound = 100
