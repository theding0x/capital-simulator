package circulation_test

import (
	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

// Compile-time assertion: circulation.Pence is an alias of money.Pence
// (ADR 0001). Assigning across the two types fails to build if the alias is
// ever replaced by a distinct named type.
var _ money.Pence = circulation.Pence(0)
