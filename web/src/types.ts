// TypeScript shapes that mirror the Go structs in commodity-service.
// Keep these in sync with services/commodity-service/internal/commodity/*.go.

export interface ConcreteLabour {
  kind: string;
  description?: string;
}

export interface UseValue {
  description: string;
  unit: string;
}

export interface Commodity {
  id: string;
  name: string;
  use_value: UseValue;
  concrete_labour: ConcreteLabour;
  snlt_per_unit: number; // labour-minutes per unit
  created_at: string;
  updated_at: string;
}

export interface ValueResponse {
  commodity: Commodity;
  quantity: number;
  value: number; // labour-minutes
  value_hours: number;
}

export interface ExchangeRatio {
  base: Commodity;
  quote: Commodity;
  base_qty: number;
  quote_qty: number;
  common_value: number;
}

export interface SimpleForm {
  relative: Commodity;
  equivalent: Commodity;
  relative_qty: number;
  equivalent_qty: number;
  common_value: number;
}

export interface LabourRelation {
  counterpart: Commodity;
  counterpart_qty: number;
  labour_time: number;
}

export interface SocialRelations {
  subject: Commodity;
  labour_per_unit: number;
  concrete_labour: ConcreteLabour;
  labour_relations: LabourRelation[];
  note: string;
}

// --- market-service types (Ch. 2: Exchange) ---------------------------------

export interface Owner {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Offer {
  id: string;
  owner_id: string;
  commodity_id: string;
  quantity: number;
  seeks_kind: string;
  seeks_commodity_id?: string;
  created_at: string;
}

export interface Exchange {
  id: string;
  giver_id: string;
  receiver_id: string;
  giver_commodity_id: string;
  giver_qty: number;
  receiver_commodity_id: string;
  receiver_qty: number;
  realised_value: number; // labour-minutes
  created_at: string;
}

export interface UniversalEquivalent {
  commodity_id: string;
  set_at: string;
}

export interface MoneyCommodity {
  commodity_id: string;
  created_at: string;
}

export interface Price {
  commodity_id: string;
  money_commodity_id: string;
  amount: number;
  updated_at: string;
}

export interface CreateOwnerInput {
  name: string;
}

export interface CreateOfferInput {
  owner_id: string;
  commodity_id: string;
  quantity: number;
  seeks_kind: string;
  seeks_commodity_id?: string;
}

export interface CreateExchangeInput {
  giver_id: string;
  receiver_id: string;
  giver_commodity_id: string;
  giver_qty: number;
  receiver_commodity_id: string;
  receiver_qty: number;
  realised_value: number;
}

export interface ComputePriceInput {
  commodity_id: string;
  money_commodity_id: string;
  commodity_snlt: number;
  money_snlt: number;
  unit_qty: number;
}

export interface CreateCommodityInput {
  name: string;
  use_value: UseValue;
  concrete_labour: ConcreteLabour;
  snlt_per_unit: number;
}

export interface UpdateCommodityInput {
  name?: string;
  use_value_description?: string;
  use_value_unit?: string;
  concrete_labour_kind?: string;
  concrete_labour_description?: string;
  snlt_per_unit?: number;
}

// --- market-service types (Ch. 3: Money) ------------------------------------

export interface CircuitLeg {
  kind: string;
  commodity_id: string;
  money_id: string;
  owner_id: string;
  value: number;
}

export interface Circuit {
  id: string;
  sale_leg: CircuitLeg;
  purchase_leg: CircuitLeg;
  created_at: string;
}

export interface Hoard {
  id: string;
  owner_id: string;
  amount: number;
  created_at: string;
}

export interface PaymentObligation {
  id: string;
  creditor_id: string;
  debtor_id: string;
  amount: number;
  created_at: string;
  due_at: string;
  paid_at?: string;
}

export interface WorldMoneyTransfer {
  id: string;
  sender_id: string;
  receiver_id: string;
  gold_mg: number;
  created_at: string;
}

export interface CreateCircuitInput {
  sale_leg: CircuitLeg;
  purchase_leg: CircuitLeg;
}

export interface CreateHoardInput {
  owner_id: string;
  amount: number;
}

export interface CreatePaymentObligationInput {
  creditor_id: string;
  debtor_id: string;
  amount: number;
  due_at: string;
}

export interface CreateWorldMoneyTransferInput {
  sender_id: string;
  receiver_id: string;
  gold_mg: number;
}

// --- agent-service types (Ch. 4: The General Formula for Capital) -----------

export interface Agent {
  id: string;
  name: string;
  class: "capitalist" | "worker" | "miser" | "owner";
  money_balance: number; // Pence (pennies); divide by 100 for £
  labour_minutes: number;
  hoarding: boolean;
  created_at: string;
  updated_at: string;
}

export interface CapitalCircuit {
  id: string;
  agent_id: string;
  m_advanced: number; // Pence
  commodity_id: string;
  m_returned: number; // Pence
  surplus_value: number; // Pence; = m_returned - m_advanced
  circuit_type: "C-M-C" | "M-C-M-prime";
  created_at: string;
}

export interface CreateAgentInput {
  name: string;
  class: "capitalist" | "worker" | "miser" | "owner";
  money_balance: number;
  labour_minutes?: number;
}

export interface UpdateAgentInput {
  name?: string;
  money_balance?: number;
  labour_minutes?: number;
}

export interface CreateAgentCircuitInput {
  m_advanced: number;
  commodity_id: string;
  m_returned: number;
  circuit_type: "C-M-C" | "M-C-M-prime";
}

// --- agent-service types (Ch. 5: Contradictions in the General Formula) ----

export interface CircuitProof {
  m: number; // Pence
  commodity_id?: string;
  m_prime: number; // Pence
  surplus_value: number; // Pence
  origin: "equivalent" | "redistribution";
}

export interface ExchangeSimulation {
  a_before: number; // Pence
  b_before: number;
  a_after: number;
  b_after: number;
  origin: "equivalent" | "redistribution";
}

export interface ComputeCircuitInput {
  m: number;
  commodity_id?: string;
  m_prime: number;
}

export interface SimulateExchangeInput {
  a_value: number;
  b_value: number;
}

// --- agent-service types (Ch. 6: The Buying and Selling of Labour-Power) -----

export interface LabourPower {
  capacity_minutes_per_day: number; // labour-minutes
}

export interface SubsistenceItem {
  name: string;
  snlt_minutes: number; // labour-minutes
  essential: boolean;
}

export interface LabourWorker {
  id: string;
  kind: "worker";
  owns_labour_power: boolean;
  owns_commodities_to_sell: boolean;
  labour_power: LabourPower;
  created_at: string;
  updated_at: string;
}

export interface LabourCapitalist {
  id: string;
  kind: "capitalist";
  money_capital: number; // labour-minutes (capital in value form)
  created_at: string;
  updated_at: string;
}

export interface LabourPowerOffering {
  id: string;
  owner_id: string;
  capacity_minutes_per_day: number; // labour-minutes
  contract_days: number;
  asking_wage: number; // labour-minutes
  created_at: string;
}

export interface LabourPowerPurchase {
  id: string;
  seller_id: string;
  buyer_id: string;
  wage_minutes: number; // labour-minutes
  contract_days: number;
  created_at: string;
}

export interface CreateLabourWorkerInput {
  owns_labour_power: boolean;
  owns_commodities_to_sell: boolean;
  capacity_minutes_per_day: number;
}

export interface CreateLabourCapitalistInput {
  money_capital: number;
}

export interface CreateLabourPowerOfferingInput {
  owner_id: string;
  capacity_minutes_per_day: number;
  contract_days: number;
  asking_wage: number;
}

export interface CreateLabourPowerPurchaseInput {
  seller_id: string;
  buyer_id: string;
  wage_minutes: number;
  contract_days: number;
}
