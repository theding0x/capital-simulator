package market

import "errors"

// MoneyRequired implements the quantity-of-money formula from Capital Vol. I,
// Ch. 3 §2.b:
//
//	M = ΣP / V
//
// where ΣP is the sum of the prices of commodities to be circulated in a
// given period and V is the average number of times each money-unit changes
// hands (its velocity). The result is the total quantity of money required
// for that period's circulation, in pence.
//
// "Other things remaining the same, the quantity of money functioning as the
// circulating medium varies directly with the amount of commodities circulating
// at given prices, and inversely with the velocity of circulation of the same
// pieces of coin." (Ch. 3 §2.b)
func MoneyRequired(sumOfPrices, velocity int64) (Pence, error) {
	if sumOfPrices < 0 {
		return 0, errors.New("market: sum_of_prices must be non-negative")
	}
	if velocity <= 0 {
		return 0, errors.New("market: velocity must be positive")
	}
	return Pence(sumOfPrices / velocity), nil
}
