package circulation

// FieldCapital is the per-capital projection consumed by the Atlas Observatory
// snapshot. It carries the latest StageDistribution (the orbit's M/P/C arcs),
// the capital's cost-price (c+v) and surplus (s) from its latest
// SupplyDemandImbalance, and its run status. Read-model only: no persistence,
// no ID constructor.
type FieldCapital struct {
	ID              IndustrialCapitalID
	TotalPence      Pence
	MoneyPence      Pence
	ProductionPence Pence
	CommodityPence  Pence
	CostPricePence  Pence // c + v (SupplyDemandImbalance.DemandPence)
	SurplusPence    Pence // s     (SupplyDemandImbalance.ExcessPence)
	Status          IndustrialCapitalStatus
}
