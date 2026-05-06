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
  labour_power_value_minutes: number; // LabourMinutes; daily reproduction cost
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
  labour_power_value_minutes?: number; // optional; defaults to 0 on the server
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

// --- agent-service types (Ch. 7: The Labour-Process and Valorization) --------

export interface RawMaterial {
  commodity_id: string;
  quantity: number;
  snlt_per_unit: number; // LabourMinutes per unit
}

export interface Instrument {
  commodity_id: string;
  wear_per_run: number; // LabourMinutes transferred per run
}

export interface MeansOfProduction {
  raw_materials: RawMaterial[];
  instruments: Instrument[];
}

export interface LabourProcess {
  id: string;
  worker_id: string;
  capitalist_id: string;
  means: MeansOfProduction;
  duration: number; // LabourMinutes
  necessary_labour_minutes: number;
  product_kind: string;
  product_quantity: number;
  created_at: string;
}

export interface Product {
  commodity_kind: string;
  quantity: number;
  total_value: number; // LabourMinutes
}

export interface ValorizationSummary {
  necessary_labour: number;
  surplus_labour: number;
  surplus_value: number;
  product_value: number;
}

export interface RunLabourProcessResult {
  labour_process: LabourProcess;
  product: Product;
  valorization: ValorizationSummary;
}

export interface RunLabourProcessInput {
  worker_id: string;
  capitalist_id: string;
  means_of_production: MeansOfProduction;
  duration_minutes: number;
  product_kind: string;
  product_quantity: number;
}

// --- commodity-service types (Ch. 8: Constant & Variable Capital) ------------

export interface ConstantCapitalInput {
  original_value: number;       // LabourMinutes
  kind: "instrument" | "raw_material" | "auxiliary";
  service_life_days: number;    // 0 = consumed wholly in one cycle
}

export interface VariableCapitalInput {
  wage_value: number;           // LabourMinutes — value of labour-power
  working_day: number;          // LabourMinutes — total hours worked
}

export interface ProductValue {
  constant: number;             // c: transferred value
  variable: number;             // v: reproduced value
  surplus: number;              // s: surplus-value
}

export interface CapitalCompositionResult {
  constant: number;
  variable: number;
  composition_ratio: number;    // c/v
}

export interface DecomposeCapitalInput {
  constant_capitals: ConstantCapitalInput[];
  variable_capital: VariableCapitalInput;
}

export interface DecomposeCapitalResult {
  product_value: ProductValue;
  capital_composition: CapitalCompositionResult;
  composition_ratio: number;
}

// --- commodity-service types (Ch. 9: Rate of Surplus-Value) ------------------

export interface ProductionAccount {
  id: string;
  constant: number;             // c (LabourMinutes)
  variable: number;             // v (LabourMinutes)
  surplus: number;              // s (LabourMinutes)
  created_at: string;
}

export interface ProductionAccountResult extends ProductionAccount {
  rate_of_surplus_value: number;    // s/v as float
  value_product: number;            // v + s (LabourMinutes)
  expanded_capital: number;         // c + v + s (LabourMinutes)
  surplus_produce_fraction: number; // s/(v+s) as float
}

export interface CreateProductionAccountInput {
  constant: number;
  variable: number;
  surplus: number;
}

export interface RateOfSurplusValueInput {
  surplus: number;
  variable: number;
}

export interface RateOfSurplusValueResult {
  rate: number;
  surplus: number;
  variable: number;
  value_product: number;
  surplus_produce_fraction: number;
}

// --- agent-service types (Ch. 10: The Working-Day) ---------------------------

export interface WorkingDay {
  id: string;
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  created_at: string;
}

export interface WorkingDayResponse {
  working_day: WorkingDay;
  total_minutes: number;
  rate_of_surplus_value: number;
  exceeds_statutory?: boolean;
}

export interface ValidateWorkingDayResponse {
  total_minutes: number;
  rate_of_surplus_value: number;
  valid: boolean;
  error?: string;
}

export interface RelaySet {
  shift_kind: "day" | "night";
  working_day: WorkingDay;
  worker_ids: string[];
}

export interface RelaySchedule {
  id: string;
  sets: [RelaySet, RelaySet];
  created_at: string;
}

export interface CreateWorkingDayInput {
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  statutory_limit_minutes?: number;
}

export interface ValidateWorkingDayInput {
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  statutory_limit_minutes?: number;
}

export interface RelaySetInput {
  shift_kind: "day" | "night";
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  worker_ids: string[];
}

export interface CreateRelayScheduleInput {
  sets: [RelaySetInput, RelaySetInput];
}

// --- simulation-engine types (Ch. 11: Rate and Mass of Surplus-Value) ---------

export interface SurplusValueRate {
  surplus_labour: number;   // LabourMinutes
  necessary_labour: number; // LabourMinutes
}

export interface SurplusValueSnapshot {
  rate: SurplusValueRate;
  variable_capital: number;
  worker_count?: number;     // omitted when worker-count formula not used
  mass: number;              // primary result
  mass_by_rate: number;      // formula I: (s/v) × V
  mass_by_workers?: number;  // omitted when worker-count formula not used
}

export interface ComputeMassInput {
  surplus_labour: number;
  necessary_labour: number;
  variable_capital?: number;
  labour_power_value?: number;
  worker_count?: number;
}

export interface SurplusLimitsResponse {
  absolute_workday_limit: number; // always 1440
  minimum_capital?: number;
  labour_power_value?: number;
  worker_count?: number;
}
