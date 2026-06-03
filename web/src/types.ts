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

// agent-service (Ch. 4: coordinated multi-agent exchange settlement, #216).
// One closed leg landed on an owner's circuit by a market exchange — distinct
// from the market C-M-C CircuitLeg above. Value is a labour-minute magnitude.
export interface AgentCircuitLeg {
  id: string;
  agent_id: string;
  kind: "sale" | "purchase";
  counterparty_id: string;
  commodity_id: string;
  value_minutes: number;
  exchange_id: string;
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

// --- simulation-engine types (Ch. 12: Relative Surplus-Value) ----------------

export interface ProductionWorkingDay {
  total: number;             // LabourMinutes
  necessary_labour: number;  // LabourMinutes
  surplus_labour: number;    // LabourMinutes
}

export interface ShortenWorkingDayResponse {
  working_day: ProductionWorkingDay;
  relative_surplus_value: number; // LabourMinutes delta (new SL − old SL)
}

export interface ProductionRateResult {
  rate: number;              // s/v as float
  surplus_labour: number;    // LabourMinutes (echoed)
  necessary_labour: number;  // LabourMinutes (echoed)
}

export interface ExtraSurplusValueResult {
  extra_surplus_value: number; // LabourMinutes total
  per_unit: number;            // LabourMinutes per article
  individual_value: number;
  social_value: number;
  quantity: number;
}

export interface RecordWorkingDayInput {
  total: number;
  labour_power_value: number;
}

export interface ShortenWorkingDayInput {
  total: number;
  necessary_labour: number;
  surplus_labour: number;
  new_labour_power_value: number;
}

export interface ExtraSurplusValueInput {
  individual_value: number;
  social_value: number;
  quantity: number;
}

// --- agent-service types (Ch. 13: Co-operation) ------------------------------

export interface CooperationMember {
  worker_id: string;
  supervisory: boolean;
  working_day_minutes: number;
}

export interface Supervisor {
  worker_id: string;
}

export interface Cooperation {
  id: string;
  name: string;
  capitalist_id: string;
  members: CooperationMember[];
  created_at: string;
  size: number;
  collective_working_day_minutes: number;
  average_social_labour_minutes: number;
  collective_power_factor: number;
  productivity_factor: number;
  cooperative_origin: string;
  supervisors: Supervisor[];
}

export interface CollectiveWorkingDayResult {
  cooperation_id: string;
  size: number;
  individual_working_day_minutes: number;
  collective_working_day_minutes: number;
  collective_power_factor: number;
  collective_output_use_value_units: number;
  cooperative_origin: string;
}

export interface AverageSocialLabourResult {
  cooperation_id: string;
  size: number;
  collective_working_day_minutes: number;
  average_social_labour_minutes: number;
  meets_burke_minimum_platoon: boolean;
}

export interface MinimumCapitalInput {
  worker_count: number;
  daily_wage_pence: number;
}

export interface MinimumCapitalResult {
  worker_count: number;
  daily_wage_pence: number;
  minimum_capital_pence: number;
}

export interface CreateCooperationInput {
  name: string;
  capitalist_id: string;
  members: CooperationMember[];
}

// Ch. 14 — Division of Labour and Manufacture
export type ManufactureForm = "heterogeneous" | "serial";
export type ManufactureOrigin = "combination" | "splitting";
export type SkillLevel = "skilled" | "unskilled";

export interface DetailRole {
  name: string;
  skill_level: SkillLevel;
  output_rate_per_hour: number;
  head_count: number;
  tool_name: string;
}

export interface CollectiveLabourerSummary {
  roles: DetailRole[];
  total_workers: number;
  output_per_period_minutes: number;
  period_minutes: number;
  is_paralysed_if_absent_one: boolean;
}

export interface Manufacture {
  id: string;
  name: string;
  capitalist_id: string;
  form: ManufactureForm;
  origin: ManufactureOrigin;
  individual_working_day_minutes: number;
  roles: DetailRole[];
  created_at: string;
  total_workers: number;
  collective_labourer: CollectiveLabourerSummary;
  hierarchy: DetailRole[];
  productive_power_factor: number;
  productivity_factor: number;
  productive_power_minutes: number;
  cooperation_baseline_minutes: number;
}

export interface CreateManufactureInput {
  name: string;
  capitalist_id: string;
  form: ManufactureForm;
  origin: ManufactureOrigin;
  individual_working_day_minutes: number;
  roles: DetailRole[];
}

export interface ProportionalGroupSizeResult {
  manufacture_id: string;
  target_output_rate: number;
  headcounts: Record<string, number>;
}

export interface ManufactureMinimumCapitalResult {
  manufacture_id: string;
  raw_material_cost_factor: number;
  minimum_capital_pence: number;
  cooperation_baseline_pence: number;
}

// Ch. 15 — Machinery and Modern Industry
export type PrimeMoverKind = "steam" | "water" | "electric" | "animal";

export interface PrimeMover {
  kind: PrimeMoverKind;
  horsepower: number;
}

export interface Machine {
  id: string;
  name: string;
  motor_mechanism: string;
  transmitting_mechanism: string;
  working_tool: string;
  machine_value: number;
  lifespan_days: number;
  productive_power: number;
  hand_labour_per_unit: number;
  accumulated_wear: number;
  accumulated_depreciation: number;
  daily_wear_and_tear: number;
  value_transferred_per_unit: number;
  daily_labour_displaced: number;
  created_at?: string;
}

export interface CreateMachineInput {
  name: string;
  motor_mechanism?: string;
  transmitting_mechanism?: string;
  working_tool?: string;
  machine_value: number;
  lifespan_days: number;
  productive_power: number;
  hand_labour_per_unit?: number;
}

export interface MachineEntryInput {
  id?: string;
  name?: string;
  motor_mechanism?: string;
  transmitting_mechanism?: string;
  working_tool?: string;
  machine_value?: number;
  lifespan_days?: number;
  productive_power?: number;
  hand_labour_per_unit?: number;
}

export interface CreateFactoryInput {
  name: string;
  prime_mover: PrimeMover;
  machines: MachineEntryInput[];
  intensity_factor?: number;
}

export interface Factory {
  id: string;
  name: string;
  prime_mover: PrimeMover;
  machines: Machine[];
  tick_count: number;
  intensity_factor: number;
  total_productive_power: number;
  daily_value_transfer: number;
  productivity_factor: number;
  created_at?: string;
}

export interface FactoryTick {
  factory_id: string;
  sequence: number;
  value_transferred: number;
  units_produced: number;
  hand_labour_saved: number;
  occurred_at: string;
}

export interface MachineWearResponse {
  id: string;
  accumulated_wear: number;
  accumulated_depreciation: number;
  remaining_value: number;
}

export interface FactoryTickResult {
  factory: Factory;
  tick: FactoryTick;
  value_transferred: number;
  units_produced: number;
  hand_labour_saved: number;
}

// Ch. 16 — Absolute and Relative Surplus-Value.
export type SurplusValueOrigin = "absolute" | "relative";

export interface SurplusValueShape {
  labour_minutes: number;
  origin: SurplusValueOrigin;
}

export interface AbsoluteSurplusValueInput {
  working_day_minutes: number;
  necessary_labour_minutes: number;
  extension_minutes: number;
}

export interface AbsoluteSurplusValueResult {
  old_working_day: ProductionWorkingDayShape;
  new_working_day: ProductionWorkingDayShape;
  absolute_surplus_value: SurplusValueShape;
}

export interface RelativeSurplusValueInput {
  working_day_minutes: number;
  necessary_labour_minutes: number;
  productivity_factor: number;
}

export interface RelativeSurplusValueResult {
  old_working_day: ProductionWorkingDayShape;
  new_working_day: ProductionWorkingDayShape;
  relative_surplus_value: SurplusValueShape;
}

export interface SurplusValueRateResult {
  surplus_labour_minutes: number;
  necessary_labour_minutes: number;
  surplus_value_minutes: number;
  total_capital_advanced: number;
  rate_of_surplus_value: number;
  rate_of_profit: number;
  mill_critique_holds: boolean;
}

// Ch. 17 — Changes of magnitude in price of labour-power and surplus-value
export interface LabourScenarioInput {
  working_day_minutes: number;
  necessary_labour_minutes: number;
  intensity_factor: number;
  productivity_factor: number;
  // Optional Ch. 25 reserve-army overlay; omit (or leave zero) for full employment.
  reserve_army_magnitude?: number;
  reserve_army_workforce_total?: number;
}

export interface LabourScenarioResult {
  daily_value_minutes: number;
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  labour_power_value_minutes: number;
  rate_of_surplus_value: number;
  law_constant_daily_value: boolean;
  // Present only when a reserve-army overlay was supplied. compressed_wage_minutes
  // is the value of labour-power after the reserve army's downward pull.
  reserve_army_pressure_bp?: number;
  compressed_wage_minutes?: number;
}

// Ch. 18 — Various Formula for the Rate of Surplus-Value
export interface RatesOfSurplusValueInput {
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
}

export interface RatesOfSurplusValueResult {
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  working_day_minutes: number;
  formula_i: number;
  formula_ii: number;
  formula_iii: number;
}

// Part IV bridge — Ch.13/14/15 productivity → Ch.12 relative surplus-value.
export type ProductivitySource = "cooperation" | "manufacture" | "factory";

export interface RelativeSurplusInput {
  working_day_total: number;
  current_lpv: number;
  productivity_factor: number;
}

export interface RelativeSurplusFromSourceInput {
  working_day_total: number;
  current_lpv: number;
  source: ProductivitySource;
  source_id: string;
}

export interface ProductionWorkingDayShape {
  total: number;
  necessary_labour: number;
  surplus_labour: number;
}

export interface RelativeSurplusResult {
  productivity_factor: number;
  old_working_day: ProductionWorkingDayShape;
  new_working_day: ProductionWorkingDayShape;
  new_labour_power_value: number;
  relative_surplus_value: number;
  source?: { source: ProductivitySource; id: string };
}

// Ch. 19 — The Transformation of the Value of Labour-Power into Wages

export interface WageLabourValue {
  daily_pence: number;
  necessary_minutes: number;
}

export interface Wage {
  daily_pence: number;
  working_day_minutes: number;
}

export interface WageAppearance {
  paid_hours: number;
  unpaid_hours: number;
}

export interface LabourDecomposition {
  paid_minutes: number;
  unpaid_minutes: number;
}

export interface WageForm {
  id: string;
  agent_id: string;
  wage: Wage;
  labour_power_value: WageLabourValue;
  created_at: string;
  hourly_wage: number;
  appearance: WageAppearance;
  decomposition: LabourDecomposition;
}

export interface CreateWageFormInput {
  agent_id: string;
  daily_pence: number;
  working_day_minutes: number;
  lpv_daily_pence: number;
  necessary_minutes: number;
}

// Ch. 20 — Time-Wages

export interface HourlyPriceOfLabour {
  numerator: number;
  denominator: number;
  as_float: number;
}

export interface NominalWage {
  pence: number;
}

// Ch. 25 → Ch. 20: stateless nominal-wage compute with an optional reserve-army overlay.
export interface NominalWageComputeInput {
  daily_labour_power_value: number;
  working_day_minutes: number;
  overtime_hours?: number;
  overtime_rate_pence?: number;
  reserve_army_magnitude?: number;
  reserve_army_workforce_total?: number;
}

export interface NominalWageComputeResult {
  nominal_wage: NominalWage;
  reserve_army_pressure_bp?: number;
  compressed_wage_pence?: number;
}

export interface WorkingSession {
  id: string;
  agent_id: string;
  wage_form_id?: string;
  daily_labour_power_value: { pence: number };
  working_day_minutes: { minutes: number };
  overtime_hours: { hours: number };
  overtime_rate_pence: { pence: number };
  wage_period: "daily" | "weekly";
  created_at: string;
  hourly_price: HourlyPriceOfLabour;
  nominal_wage: NominalWage;
}

export interface CreateWorkingSessionInput {
  agent_id: string;
  wage_form_id?: string;
  daily_labour_power_value?: number;
  working_day_minutes: number;
  overtime_hours: number;
  overtime_rate_pence: number;
  wage_period: "daily" | "weekly";
}

export interface ComputeHourlyPriceInput {
  daily_pence: number;
  working_day_minutes: number;
}

// --- Ch. 21: Piece-Wages ---

export interface PieceWage {
  id: string;
  agent_id: string;
  wage_form_id?: string;
  price_pence: number;
  normal_output: number;
  implied_daily_wage: number;
  created_at: string;
}

export interface CreatePieceWageInput {
  wage_form_id?: string;
  price_pence?: number;
  normal_output: number;
}

export interface SubContract {
  id: string;
  head_labourer_id: string;
  assistant_ids: string[];
  piece_rate_pence: number;
  assistant_rate_pence: number;
  created_at: string;
  spread: number;
}

export interface CreateSubContractInput {
  head_labourer_id: string;
  assistant_ids: string[];
  piece_rate_pence: number;
  assistant_rate_pence: number;
}

export interface ComputePiecePriceInput {
  daily_wage_farthings: number;
  day_value_product_farthings: number;
  normal_output: number;
  pieces_produced?: number;
  quality_outcome?: "accepted" | "rejected";
}

export interface ComputePiecePriceResult {
  piece_price: number;
  piece_value: number;
  actual_earnings: number;
}

// --- agent-service types (Ch. 22: National Differences of Wages) ---

export interface NationalIntensity {
  country_code: string;
  factor: number;
}

export interface DayWage {
  country_code: string;
  nominal_pence: number;
  working_day_minutes: number;
}

export interface StandardisedWage {
  country_code: string;
  amount: number;
}

export interface RelativeLabourPrice {
  country_code: string;
  ratio: number;
}

export interface SpindleRatio {
  country_code: string;
  spindles_per_worker: number;
}

export interface WageComparison {
  countries: string[];
  day_wages: DayWage[];
  standardised_wages: StandardisedWage[];
  relative_prices: RelativeLabourPrice[];
  spindle_ratios: SpindleRatio[];
}

export interface RegisterIntensityInput {
  country_code: string;
  factor: number;
}

export interface RegisterDayWageInput {
  country_code: string;
  nominal_pence: number;
  working_day_minutes: number;
}

// --- simulation-engine types (Ch. 23: Simple Reproduction) -------------------

export interface SimpleReproductionInput {
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  periods: number;
}

export interface ReproductionCycleItem {
  period: number;
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  surplus_total: number;
  revenue: number;
  accumulated: number;
}

export interface SimpleReproductionResult {
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  periods: number;
  cycles: ReproductionCycleItem[];
  repayment_period: number;
}

export interface RepaymentPeriodInput {
  capital: number;
  annual_revenue: number;
}

export interface RepaymentPeriodResult {
  capital: number;
  annual_revenue: number;
  repayment_period: number;
}

// --- simulation-engine types (Ch. 24: The Transformation of Surplus-Value into Capital) ---

export interface AccumulationPeriod {
  period: number;
  constant_capital: number;
  variable_capital: number;
  surplus_produced: number;
  new_constant: number;
  new_variable: number;
  revenue: number;
}

export interface ExtendedReproductionInput {
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  accum_rate: number;
  composition_ratio: number;
  periods: number;
}

export interface ExtendedReproductionResult {
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  accum_rate: number;
  composition_ratio: number;
  periods: number;
  cycles: AccumulationPeriod[];
}

export interface SplitSurplusInput {
  surplus: number;
  accum_rate: number;
  composition_ratio: number;
}

export interface SplitSurplusResult {
  surplus: number;
  accum_rate: number;
  composition_ratio: number;
  new_constant: number;
  new_variable: number;
  revenue: number;
}

// --- simulation-engine types (Ch. 25: The General Law of Capitalist Accumulation) ---

export interface OrganicCompositionInput {
  constant_capital: number;
  variable_capital: number;
}

export interface OrganicCompositionResult {
  constant_capital: number;
  variable_capital: number;
  ratio: number;
}

export interface LabourDemandInput {
  total_capital: number;
  organic_composition_ratio: number;
  wage_pence: number;
}

export interface LabourDemandResult {
  total_capital: number;
  organic_composition_ratio: number;
  wage_pence: number;
  workers: number;
}

export interface ReserveArmyInput {
  worker_supply: number;
  workers_demanded: number;
}

export interface ReserveArmyResult {
  worker_supply: number;
  workers_demanded: number;
  reserve_army_size: number;
  relative_proportion: number;
}

export interface GeneralLawSnapshot {
  period: number;
  constant_capital: number;
  variable_capital: number;
  organic_composition: number;
  workers: number;
  reserve_army_size: number;
  relative_proportion: number;
}

export interface GeneralLawScenarioInput {
  name: string;
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  accumulation_rate: number;
  productivity_growth: number;
  wage_pence: number;
  worker_supply: number;
  periods: number;
}

export interface GeneralLawScenarioResult {
  id: string;
  name: string;
  constant_capital: number;
  variable_capital: number;
  surplus_rate: number;
  accumulation_rate: number;
  productivity_growth: number;
  wage_pence: number;
  worker_supply: number;
  periods: number;
  created_at: string;
  series: GeneralLawSnapshot[];
}

// --- simulation-engine (Ch. 26: The Secret of Primitive Accumulation) -------

export interface PrimitiveAccumulation {
  period: string;
  method: string;
  labourers_expropriated: number;
  capital_formed: number;
}

export interface HistoricalStageInput {
  name: string;
  description: string;
  primitive_accumulations: PrimitiveAccumulation[];
}

export interface HistoricalStage {
  id: string;
  name: string;
  description: string;
  primitive_accumulations: PrimitiveAccumulation[];
  total_capital_formed: number;
  total_labourers_expropriated: number;
  created_at: string;
}

export interface SeedScenarioInput {
  pre_capitalist_workers: number;
  composition_ratio: number;
  surplus_rate: number;
  periods: number;
}

export interface SeededCapitalStock {
  constant_capital: number;
  variable_capital: number;
}

export interface ProducerSeparation {
  pre_capitalist_workers: number;
  displaced_workers: number;
  free_proletarians: number;
}

export interface SeededReproductionCycle {
  period: number;
  capital: SeededCapitalStock;
  surplus_rate: number;
  fund: {
    total: number;
    revenue: number;
    accumulated: number;
  };
}

export interface SeedScenarioResult {
  stage_id: string;
  stage_name: string;
  seeded_capital: SeededCapitalStock;
  separation: ProducerSeparation;
  surplus_rate: number;
  periods: number;
  cycles: SeededReproductionCycle[];
}

// Ch. 27 — Expropriation of the Agricultural Population

export interface EnclosureEventInput {
  period: string;
  acres_enclosed: number;
  population_displaced: number;
  beneficiary: string;
}

export interface EnclosureEvent {
  id: string;
  period: string;
  acres_enclosed: number;
  population_displaced: number;
  beneficiary: string;
  created_at: string;
}

// Ch. 28 — Bloody Legislation Against the Expropriated

export interface WageStatuteInput {
  period: string;
  jurisdiction: string;
  max_wage_pence: number;
  min_wage_pence: number;
  enforcement_penalty: string;
}

export interface WageStatute {
  id: string;
  historical_stage_id: string;
  period: string;
  jurisdiction: string;
  max_wage_pence: number;
  min_wage_pence: number;
  enforcement_penalty: string;
  created_at: string;
}

export interface VagrancyLawInput {
  period: string;
  jurisdiction: string;
  punishment: string;
  target_population: string;
}

export interface VagrancyLaw {
  id: string;
  historical_stage_id: string;
  period: string;
  jurisdiction: string;
  punishment: string;
  target_population: string;
  created_at: string;
}

export interface StatutoryWage {
  acted_wage_pence: number;
  market_wage_pence: number;
  deviation: number;
}

export interface LabourDisciplineRegime {
  period: string;
  mechanisms: string[];
  wage_statutes: WageStatute[];
  vagrancy_laws: VagrancyLaw[];
}

// Ch. 29 — Genesis of the Capitalist Farmer

export type TenantForm = "bailiff" | "metayer" | "capitalist-farmer";

export interface FarmTenureInput {
  form: TenantForm;
  lease_period_years: number;
  rent_pence: number;
  capital_advanced_pence: number;
  revenue_pence: number;
  wage_costs_pence: number;
}

export interface FarmingSurplus {
  revenue: number;
  nominal_rent: number;
  wage_costs: number;
  profit: number;
}

export interface FarmTenure {
  id: string;
  historical_stage_id: string;
  form: TenantForm;
  lease_period_years: number;
  rent_pence: number;
  capital_advanced_pence: number;
  revenue_pence: number;
  wage_costs_pence: number;
  surplus: FarmingSurplus;
  created_at: string;
}

export interface RealRentInput {
  nominal_rent_pence: number;
  depreciation_factor: number;
}

export interface RealRentResult {
  nominal_rent_pence: number;
  real_rent_pence: number;
  depreciation_factor: number;
}

// Ch. 30 — Reaction of the Agricultural Revolution on Industry

export interface DomesticIndustryInput {
  name: string;
  households_engaged: number;
  annual_output_pence: number;
}

export interface DomesticIndustry {
  id: string;
  historical_stage_id: string;
  name: string;
  households_engaged: number;
  annual_output_pence: number;
  created_at: string;
}

export interface MarketFormationInput {
  name: string;
  households_engaged: number;
  annual_output_pence: number;
  expropriated: number;
}

export interface MarketFormation {
  period: string;
  labourers_proletarianised: number;
  domestic_industry_destroyed: boolean;
  new_wage_labourers: number;
  market_size_pence: number;
}

export interface HomeMarketSize {
  commoditised_output_pence: number;
  wage_labourers: number;
}

// Ch. 31 — Genesis of the Industrial Capitalist

export interface CapitalOriginInput {
  source: string;
  amount_pence: number;
  period: string;
}

export interface CapitalOrigin {
  id: string;
  historical_stage_id: string;
  source: string;
  amount_pence: number;
  period: string;
  created_at: string;
}

export interface ColonialTransferInput {
  from: string;
  to: string;
  value_pence: number;
  method: string;
}

export interface ColonialTransfer {
  id: string;
  historical_stage_id: string;
  from: string;
  to: string;
  value_pence: number;
  method: string;
  created_at: string;
}

export interface NationalDebtInput {
  amount_pence: number;
  interest_rate_bps: number;
  creditor_class: string;
}

export interface NationalDebt {
  id: string;
  historical_stage_id: string;
  amount_pence: number;
  interest_rate_bps: number;
  creditor_class: string;
  created_at: string;
}

export interface ProtectionSystem {
  id: string;
  historical_stage_id: string;
  tariff_rate_bps: number;
  beneficiary: string;
  period_start: string;
  period_end: string;
  created_at: string;
}

export interface IndustrialCapitalGenesis {
  historical_stage_id: string;
  origins: CapitalOrigin[];
  colonial_transfers: ColonialTransfer[];
  national_debts: NationalDebt[];
  protection_systems: ProtectionSystem[];
  total_capital_formed_pence: number;
}

// Ch. 32 — The Historical Tendency of Capitalist Accumulation.

export interface Negation {
  stage: string;
  description: string;
}

export interface NegationOfNegationResponse {
  stages: Negation[];
  negation: Negation;
}

export interface CentralisationStep {
  step_index: number;
  firms_absorbed: number;
  capital_concentrated_pence: number;
}

export interface AccumulationTrajectory {
  id: string;
  name: string;
  initial_firms: number;
  initial_capital_pence: number;
  steps: CentralisationStep[];
  final_firms: number;
  final_capital_pence: number;
  reserve_army_size: number;
  created_at?: string;
}

export interface RunCentralisationInput {
  name?: string;
  firms: number;
  total_capital_pence: number;
  wage_labourers: number;
  steps: number;
  absorption_rate: number;
  capital_growth_rate: number;
  persist?: boolean;
}

// Ch. 33 — The Modern Theory of Colonisation.

export interface ColonialLabourMarket {
  id: string;
  colony: string;
  free_labourers: number;
  annual_wage_pence: number;
  land_available: boolean;
  wakefield_scheme_applied: boolean;
  independence_years: number;
  surplus_labour_extractable: boolean;
  created_at?: string;
}

export interface CreateColonialLabourMarketInput {
  colony: string;
  free_labourers: number;
  annual_wage_pence: number;
  land_available: boolean;
}

export interface SufficientPriceInput {
  price_per_acre_pence: number;
  desired_acres: number;
}

export interface RegulateColonialMarketInput {
  sufficient_price: SufficientPriceInput;
}

export interface ComputeIndependenceInput {
  worker_id: string;
  sufficient_price: SufficientPriceInput;
  years_limit?: number;
}

export interface WageWorkerIndependence {
  worker_id: string;
  annual_savings_pence: number;
  years_worked: number;
  savings_pence: number;
  target_ransom_pence: number;
  became_landowner: boolean;
}

// Vol. II Ch. 2 — The Circuit of Productive Capital
export interface LatentMoneyCapital {
  productive_circuit_id: string;
  accumulated: number;
  threshold: number;
  is_armed: boolean;
}

export interface ReserveFund {
  productive_circuit_id: string;
  balance: number;
}

export interface RevenueCircuit {
  productive_circuit_id: string;
  amount: number;
  spent_at: string;
}

export interface CapitalisationStep {
  productive_circuit_id: string;
  amount_injected: number;
  delta_constant_pence: number;
  delta_variable_pence: number;
  occurred_at: string;
}

export interface ReserveDraw {
  productive_circuit_id: string;
  drawn_pence: number;
  reason: string;
  occurred_at: string;
}

export interface ProductiveCircuit {
  id: string;
  agent_id?: string;
  constant_pence: number;
  variable_pence: number;
  total_advance: number;
  mode: string;
  min_capitalisation_increment_pence: number;
  latent_money_capital: LatentMoneyCapital;
  reserve_fund: ReserveFund;
  revenue_circuits: RevenueCircuit[];
  capitalisation_steps: CapitalisationStep[];
  reserve_draws: ReserveDraw[];
  created_at?: string;
}

export interface CreateProductiveCircuitInput {
  agent_id?: string;
  constant_pence: number;
  variable_pence: number;
  mode: "simple" | "extended" | "mixed";
  min_capitalisation_increment_pence: number;
}

// Vol. II Ch. 3 — The Circuit of Commodity-Capital
export interface OpeningCommodityCapital {
  constant_pence: number;
  variable_pence: number;
  surplus_pence: number;
  value_original_pence: number;
  total: number;
  pounds_total: number;
}

export interface CommodityAugmented {
  commodity_circuit_id: string;
  constant_pence: number;
  variable_pence: number;
  surplus_pence: number;
  capitalised_pence: number;
  total: number;
  pounds_total: number;
  is_extended: boolean;
  closed_at: string;
}

export interface ValueComposition {
  constant_share: number;
  variable_share: number;
  surplus_share: number;
  total: number;
}

export interface SuccessivePartialSale {
  commodity_circuit_id: string;
  quantity: number;
  realised_pence: number;
  decomposition: ValueComposition;
  sold_at: string;
}

export interface MeansOfProductionSource {
  commodity_circuit_id: string;
  source_commodity_circuit_id?: string;
  source_kind: "producer_circuit" | "merchant_holding" | "import";
}

export interface CommodityCircuit {
  id: string;
  agent_id?: string;
  initial: OpeningCommodityCapital;
  mode: "simple" | "extended" | "mixed";
  is_first_investment: boolean;
  social_capital_lens: boolean;
  partial_sales: SuccessivePartialSale[];
  mp_sources: MeansOfProductionSource[];
  terminal?: CommodityAugmented;
  created_at: string;
}

export interface CreateCommodityCircuitInput {
  agent_id?: string;
  constant_pence: number;
  variable_pence: number;
  surplus_pence: number;
  pounds_total: number;
  mode: "simple" | "extended" | "mixed";
  is_first_investment?: boolean;
  social_capital_lens?: boolean;
}

export interface CommodityCircuitAggregate {
  circuit_count: number;
  closed_count: number;
  total_initial_pence: number;
  total_terminal_pence: number;
}

// Vol. II Ch. 1 — The Circuit of Money-Capital (M—C(Lp+Mp)…P…C'—M')

export type CircuitMoment = "M" | "M-C" | "P" | "C-prime" | "C-M-prime" | "M-prime";

export interface MoneyAdvance {
  amount: number;       // Pence
  agent_id?: string;
  advanced_at: string;  // RFC3339
}

export interface LabourLeg {
  amount: number;           // Pence
  labour_power_hours: number;
}

export interface MeansLeg {
  amount: number;             // Pence
  means_capacity_hours: number;
}

export interface PurchasePhase {
  money_circuit_id: string;
  labour_leg: LabourLeg;
  means_leg: MeansLeg;
  total: number; // Pence
}

export interface ProductiveState {
  money_circuit_id: string;
  entered_at: string;    // RFC3339
  constant_pence: number;
  variable_pence: number;
  total_advance: number; // Pence
}

export interface MoneyCircuitCommodityCapital {
  money_circuit_id: string;
  commodity_id?: string;
  value_original: number; // Pence
  value_surplus: number;  // Pence
  total: number;          // Pence
}

export interface Realisation {
  money_circuit_id: string;
  sold_at: string;          // RFC3339
  realised_pence: number;
  surplus_realised: number; // Pence
}

export interface MoneyCircuit {
  id: string;
  advance: MoneyAdvance;
  moment: CircuitMoment;
  purchase?: PurchasePhase;
  productive?: ProductiveState;
  commodity?: MoneyCircuitCommodityCapital;
  realisation?: Realisation;
  created_at: string; // RFC3339
}

export interface CreateMoneyCircuitInput {
  amount: number;        // Pence
  agent_id?: string;
  advanced_at?: string;  // RFC3339 — omit to use server time
}

// Vol. II Ch. 5 — The Time of Circulation

export interface ProductionTime {
  labour_time_nanos: number;
  labour_interruption_nanos: number;
  latent_nanos: number;
  natural_process_nanos: number;
}

export interface CirculationTime {
  selling_time_nanos: number;
  buying_time_nanos: number;
}

export interface TurnoverTime {
  id: string;
  industrial_capital_id: string;
  production: ProductionTime;
  circulation: CirculationTime;
  circulation_open: boolean;
  created_at: string; // RFC3339
  updated_at: string; // RFC3339
}

export interface CreateTurnoverTimeInput {
  industrial_capital_id: string;
}

export type SellingOutcome = "sold" | "spoiled" | "partial" | "pending";

export interface SellingPhase {
  id: string;
  turnover_time_id: string;
  industrial_capital_id: string;
  commodity_circuit_id: string;
  opened_at: string;      // RFC3339
  closed_at?: string;     // RFC3339
  outcome: SellingOutcome;
}

export interface BuyingPhase {
  id: string;
  turnover_time_id: string;
  industrial_capital_id: string;
  money_circuit_id: string;
  opened_at: string;      // RFC3339
  closed_at?: string;     // RFC3339
  market_location: string;
}

export type NaturalProcessKind = "ripening" | "fermentation" | "tanning" | "drying" | "other";

export interface NaturalProcessSpan {
  id: string;
  turnover_time_id: string;
  industrial_capital_id: string;
  process: NaturalProcessKind;
  duration_nanos: number;
  started_at: string; // RFC3339
}

export interface LatentProductiveCapital {
  id: string;
  turnover_time_id: string;
  industrial_capital_id: string;
  pence: number;
  held_at: string;                // RFC3339
  entered_production_at?: string; // RFC3339
}

export interface Perishability {
  id: string;
  commodity_id: string;
  window_nanos: number;
}

export interface MarketSeparation {
  id: string;
  industrial_capital_id: string;
  selling_market_id: string;
  buying_market_id: string;
}

export interface ActiveFractionResponse {
  turnover_time_id: string;
  production_time_nanos: number;
  total_nanos: number;
  active_fraction_basis_points: number;
}

// Vol. II Ch. 4 — The Three Formulas of the Circuit

export type CapitalStage = "money" | "production" | "commodity";
export type IndustrialCapitalStatus = "active" | "halted";
export type EconomyMode = "natural" | "money" | "credit";
export type StageBlockReason = "unsold_commodity" | "missing_means_of_production" | "delayed_money_return";
export type ValueRevolutionAffecting = "means_of_production" | "labour_power" | "output";

export interface StageDistribution {
  id: string;
  industrial_capital_id: string;
  at: string;
  money_pence: number;
  production_pence: number;
  commodity_pence: number;
}

export interface StageBlock {
  id: string;
  industrial_capital_id: string;
  stage: CapitalStage;
  reason: StageBlockReason;
  opened_at: string;
  closed_at?: string;
}

export interface IndustrialCapital {
  id: string;
  agent_id?: string;
  money_circuit_id?: string;
  productive_circuit_id?: string;
  commodity_circuit_id?: string;
  total_pence: number;
  economy_mode: EconomyMode;
  stagnation_tolerance_ticks: number;
  status: IndustrialCapitalStatus;
  latest_distribution?: StageDistribution;
  open_blocks: StageBlock[];
  created_at: string;
  updated_at: string;
}

export interface CapitalPart {
  id: string;
  industrial_capital_id: string;
  pence: number;
  stage: CapitalStage;
  entered_stage_at: string;
}

export interface ValueRevolutionEvent {
  id: string;
  industrial_capital_id: string;
  affecting: ValueRevolutionAffecting;
  direction_pence: number;
  occurred_at: string;
}

export interface MoneyCapitalSetFree {
  id: string;
  value_revolution_event_id: string;
  industrial_capital_id: string;
  amount_pence: number;
}

export interface MoneyCapitalTiedUp {
  id: string;
  value_revolution_event_id: string;
  industrial_capital_id: string;
  amount_pence: number;
}

export interface ValueRevolutionResult {
  event: ValueRevolutionEvent;
  set_free?: MoneyCapitalSetFree;
  tied_up?: MoneyCapitalTiedUp;
}

export interface MetamorphosisInterlock {
  id: string;
  buyer_industrial_capital_id: string;
  seller_industrial_capital_id: string;
  seller_mp_source_origin: string;
  pence: number;
  occurred_at: string;
}

export interface SupplyDemandImbalance {
  id: string;
  industrial_capital_id: string;
  period: string;
  demand_pence: number;
  supply_pence: number;
  excess_pence: number;
}

export interface AggregateSupplyDemand {
  period: string;
  demand_pence: number;
  supply_pence: number;
  excess_pence: number;
}

export interface SinkingFund {
  id: string;
  industrial_capital_id: string;
  fixed_capital_pence: number;
  lifetime_years: number;
  annual_payment_pence: number;
  accumulated_pence: number;
}

// Vol. II Ch. 6 — Costs of Circulation

export type CirculationCostKind =
  | "purchase_sale_time"
  | "book_keeping"
  | "money_as_faux_frais"
  | "supply_formation"
  | "commodity_supply"
  | "transportation";

export type CostNature = "value_preserving" | "value_adding";

export interface CirculationCost {
  id: string;
  industrial_capital_id: string;
  kind: CirculationCostKind;
  nature: CostNature;
  pence: number;
  incurred_at: string;
}

export interface CirculationCostList {
  items: CirculationCost[];
}

export interface CirculationAgent {
  id: string;
  industrial_capital_id: string;
  agent_id: string;
  role: string;
  wage_rate_pence: number;
  labour_minutes_necessary: number;
  labour_minutes_surplus: number;
}

export interface MoneyAsFauxFrais {
  id: string;
  industrial_capital_id: string;
  pence: number;
  annual_replacement_pence: number;
}

export interface StorageCost {
  id: string;
  commodity_supply_id: string;
  building_pence: number;
  labour_pence: number;
  preservation_labour_pence: number;
}

export interface CommoditySupply {
  id: string;
  industrial_capital_id: string;
  commodity_id: string;
  pence: number;
  held_since: string;
  is_voluntary: boolean;
  storage_costs: StorageCost[];
}

export interface TransportTariff {
  id: string;
  commodity_id: string;
  base_pence_per_ton_mile: number;
  fragility_multiplier_basis_points: number;
  breakage_risk_multiplier_basis_points: number;
}

export interface TransportLeg {
  id: string;
  commodity_id: string;
  quantity: number;
  origin_id: string;
  destination_id: string;
  distance_meters: number;
  weight_grams: number;
  volume_cubic_millimeters: number;
  labour_cost_pence: number;
  means_of_transport_pence: number;
  surplus_pence: number;
}

export interface AggregateCirculationCostsResult {
  industrial_capital_id: string;
  period: string;
  value_preserving_pence: number;
  value_adding_pence: number;
  total_pence: number;
}

export interface SystemFauxFraisItem {
  id: string;
  industrial_capital_id: string;
  pence: number;
  annual_replacement_pence: number;
}

export interface SystemFauxFraisResult {
  period: string;
  items: SystemFauxFraisItem[];
  total_pence: number;
  total_annual_replacement_pence: number;
}

// Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers
export type CircuitLens = "money" | "productive";

export interface Turnover {
  id: string;
  industrial_capital_id: string;
  lens: CircuitLens;
  turnover_time_minutes: number;
  cycles: string[];
}

export interface TurnoverCycle {
  id: string;
  turnover_id: string;
  started_at: string;
  ended_at: string;
  advance_pence: number;
  returned_pence: number;
  production_minutes: number;
  circulation_minutes: number;
}

export interface TurnoverNumber {
  id: string;
  turnover_id: string;
  year_reference_minutes: number;
  turnover_time_minutes: number;
  basis_points: number;
}

// Vol. II Ch. 8 — Fixed Capital and Circulating Capital

export type CapitalKind = "fixed" | "circulating";
export type RoleInProcess =
  | "instrument_of_labour"
  | "raw_material"
  | "auxiliary_material"
  | "labour_power";
export type WearCause = "use" | "atmospheric" | "moral_depreciation";
export type WearModel = "linear" | "quadratic";
export type RepairKind = "running" | "large";
export type ReinvestmentKind = "intensive" | "extensive";

export interface CapitalComponent {
  id: string;
  industrial_capital_id: string;
  kind: CapitalKind;
  role: RoleInProcess;
  pence: number;
  entered_at: string;
}

export interface FixedCapitalItem {
  id: string;
  capital_component_id: string;
  description: string;
  pence_purchased: number;
  service_life_minutes: number;
  wear_model: WearModel;
  purchased_at: string;
  retired_at?: string;
}

export interface FixedCapitalSubcomponent {
  id: string;
  fixed_capital_item_id: string;
  name: string;
  pence: number;
  service_life_minutes: number;
  wear_model: WearModel;
  last_replaced_at: string;
}

export interface WearAndTear {
  id: string;
  fixed_capital_item_id: string;
  pence: number;
  period_covered_minutes: number;
  calculated_at: string;
  cause: WearCause;
}

export interface SinkingFundForItem {
  id: string;
  fixed_capital_item_id: string;
  accumulated_pence: number;
  target_pence: number;
}

export interface Repair {
  id: string;
  fixed_capital_item_id: string;
  kind: RepairKind;
  pence: number;
  occurred_at: string;
  service_life_extension_minutes: number;
}

export interface CirculatingCycle {
  id: string;
  capital_component_id: string;
  started_at: string;
  ended_at: string;
  pence: number;
}

export interface SinkingFundReinvestment {
  id: string;
  fixed_capital_item_id: string;
  pence: number;
  kind: ReinvestmentKind;
  occurred_at: string;
}

// --- Vol. II Ch. 9 — The Aggregate Turnover of Advanced Capital --------------

export type TurnoverDifferenceKind = "real" | "apparent";

export type CrisisCyclePhase =
  | "depression"
  | "medium_activity"
  | "precipitancy"
  | "crisis";

export interface ComponentTurnoverContribution {
  id: string;
  aggregate_turnover_id: string;
  capital_component_id: string;
  kind: CapitalKind;
  advanced_pence: number;
  turnover_number_basis_points: number;
  annual_return_pence: number;
  difference_kind: TurnoverDifferenceKind;
}

export interface CrisisCyclePhaseRecord {
  id: string;
  lifetime_cycle_id: string;
  phase: CrisisCyclePhase;
  observed_at: string;
  notes?: string;
}

export interface LifetimeCycle {
  id: string;
  aggregate_turnover_id: string;
  duration_minutes: number;
  started_at: string;
  completed_cycles: number;
  crisis_phases?: CrisisCyclePhaseRecord[];
}

export interface AggregateTurnover {
  id: string;
  industrial_capital_id: string;
  lens: CircuitLens;
  advanced_total_pence: number;
  annual_turned_over_pence: number;
  aggregate_number_basis_points: number;
  contributions: ComponentTurnoverContribution[];
  lifetime_cycle?: LifetimeCycle;
  latest_crisis_phase?: CrisisCyclePhaseRecord;
}

// Vol. II Ch. 10–11 — Theories of Fixed and Circulating Capital:
// Physiocrats, Adam Smith, and Ricardo.

export type KnownEconomistError =
  // Ch. 10–11: fixed/circulating capital theories
  | "error_smith_conflation"
  | "error_smith_circulation_capital_conflation"
  | "error_smith_revenue_in_capital"
  | "error_ricardo_durability_collapse"
  | "error_ricardo_conflation"
  | "error_ricardo_fixed_capital_price_explanation"
  | "error_ricardo_no_aggregate_turnover"
  | "error_ricardo_no_value_revolution"
  // Ch. 19: former presentations of reproduction
  | "error_quesnay_productive_sterile_division"
  | "error_quesnay_missing_value_theory"
  | "error_smith_revenue_dogma"
  | "error_smith_missing_constant_replacement"
  | "error_smith_labour_regress"
  | "error_ricardo_incomplete_reproduction";

export interface EconomistAttribution {
  id: string;
  concept: string;
  theorist: string;
  edition_year: number;
  anticipates: string;
  errors: KnownEconomistError[];
}

// Vol. II Ch. 12 — The Working Period
export type WorkingPeriodMode = "discrete" | "connected";
export type WorkingPeriodShorteningKind =
  | "cooperation"
  | "machinery"
  | "parallel_labour"
  | "bakewell_breeding";

export interface WorkingPeriodInterruption {
  id: string;
  working_period_id: string;
  interrupted_at: string;
  resumed_at?: string;
  deterioration_pence: number;
  waste_of_means_of_production_pence: number;
}

export interface WorkingPeriodShortening {
  id: string;
  working_period_id: string;
  kind: WorkingPeriodShorteningKind;
  old_working_day_count: number;
  new_working_day_count: number;
  additional_fixed_capital_pence: number;
  additional_labour_count: number;
  occurred_at: string;
}

export interface WorkingPeriodCreditFinancing {
  id: string;
  working_period_id: string;
  own_capital_pence: number;
  borrowed_capital_pence: number;
  mortgage_provider_id?: string;
  created_at: string;
}

export interface WorkingPeriod {
  id: string;
  industrial_capital_id: string;
  commodity_id: string;
  working_day_count: number;
  mode: WorkingPeriodMode;
  working_day_length_minutes: number;
  created_at: string;
  interruptions: WorkingPeriodInterruption[];
  shortenings: WorkingPeriodShortening[];
  credit_financing?: WorkingPeriodCreditFinancing;
}

export interface NaturalWorkingPeriodConstraint {
  id: string;
  commodity_id: string;
  min_working_days: number;
  reason: string;
  created_at: string;
}

export interface CreateWorkingPeriodInput {
  industrial_capital_id: string;
  commodity_id: string;
  working_day_count: number;
  mode: WorkingPeriodMode;
  working_day_length_minutes: number;
}

export interface CreateWorkingPeriodInterruptionInput {
  deterioration_pence: number;
  waste_of_means_of_production_pence: number;
}

export interface CreateWorkingPeriodShorteningInput {
  kind: WorkingPeriodShorteningKind;
  old_working_day_count: number;
  new_working_day_count: number;
  additional_fixed_capital_pence: number;
  additional_labour_count: number;
}

export interface CreateWorkingPeriodCreditFinancingInput {
  own_capital_pence: number;
  borrowed_capital_pence: number;
  mortgage_provider_id?: string;
}

export interface CreateNaturalConstraintInput {
  commodity_id: string;
  min_working_days: number;
  reason: string;
}

// --- simulation-engine types (Vol. II Ch. 13: Time of Production) ----------

export interface ProductionLabourGap {
  id: string;
  working_period_id: string;
  production_time_nanos: number;
  labour_time_nanos: number;
  gap_nanos: number;
  created_at: string;
}

export interface NaturalProcessActivation {
  id: string;
  working_period_id: string;
  kind: string;
  started_at: string;
  ended_at?: string;
  requires_labour_maintenance: boolean;
}

export interface NaturalSubject {
  id: string;
  working_period_id: string;
  commodity_id: string;
  description: string;
  is_living: boolean;
  created_at: string;
}

export interface NaturalProcessIndustry {
  id: string;
  industry_name: string;
  typical_production_days_count: number;
  typical_labour_days_count: number;
  ratio_basis_points: number;
  created_at: string;
}

export interface CreateProductionLabourGapInput {
  production_time_nanos: number;
  labour_time_nanos: number;
  gap_nanos: number;
}

export interface CreateNaturalProcessActivationInput {
  kind: string;
  started_at: string;
  ended_at?: string;
  requires_labour_maintenance: boolean;
}

export interface CreateNaturalSubjectInput {
  commodity_id: string;
  description: string;
  is_living: boolean;
}

export interface CreateIndustryBenchmarkInput {
  industry_name: string;
  typical_production_days_count: number;
  typical_labour_days_count: number;
  ratio_basis_points: number;
}


// Vol. II Ch. 14 — The Time of Circulation (refined)

export type CirculationSpeedReason = "telegraph" | "steamship" | "telephone" | "railway" | "other";

export interface SellingPhaseDetailed {
  id: string;
  selling_phase_id: string;
  turnover_time_id: string;
  industrial_capital_id: string;
  market_id: string;
  market_distance_meters: number;
  communication_lag_minutes: number;
  buyer_search_minutes: number;
  bargaining_minutes: number;
  legal_settlement_minutes: number;
  total_minutes: number;
  created_at: string;
}

export interface BuyingPhaseDetailed {
  id: string;
  buying_phase_id: string;
  turnover_time_id: string;
  industrial_capital_id: string;
  market_id: string;
  market_distance_meters: number;
  communication_lag_minutes: number;
  supplier_search_minutes: number;
  order_processing_minutes: number;
  legal_settlement_minutes: number;
  total_minutes: number;
  created_at: string;
}

export interface DistanceLagRelation {
  id: string;
  market_id: string;
  distance_meters: number;
  lag_minutes_per_km: number;
  created_at: string;
}

export interface CirculationSpeedImprovement {
  id: string;
  market_id: string;
  old_lag_minutes_per_km: number;
  new_lag_minutes_per_km: number;
  reason: CirculationSpeedReason;
  effective_at: string;
  created_at: string;
}

export interface AnnualSurplusPenalty {
  id?: string;
  industrial_capital_id: string;
  period: string;
  ideal_annual_rate_basis_points: number;
  actual_annual_rate_basis_points: number;
  penalty_basis_points: number;
}

export interface CreateSellingPhaseDetailedInput {
  selling_phase_id: string;
  industrial_capital_id: string;
  market_id: string;
  market_distance_meters: number;
  communication_lag_minutes: number;
  buyer_search_minutes: number;
  bargaining_minutes: number;
  legal_settlement_minutes: number;
  total_minutes: number;
  distance_lag_relation_id?: string;
}

export interface CreateDistanceLagRelationInput {
  market_id: string;
  distance_meters: number;
  lag_minutes_per_km: number;
}

export interface CreateCirculationSpeedImprovementInput {
  market_id: string;
  old_lag_minutes_per_km: number;
  new_lag_minutes_per_km: number;
  reason: CirculationSpeedReason;
  effective_at: string;
}

// ── Vol. II Ch. 15 — The Effects of a Change of Prices ──────────────────────

export interface ValueRevolutionEventV2 {
  id: string;
  industrial_capital_id: string;
  affecting: ValueRevolutionAffecting;
  direction_pence: number;
  occurred_at: string;
  basis_points_change: number;
  revalued_capital_pence: number;
}

export interface RevaluedCapital {
  id: string;
  value_revolution_event_id: string;
  industrial_capital_id: string;
  old_advance_pence: number;
  new_advance_pence: number;
  delta_pence: number;
}

export interface RealisedValueDelta {
  id: string;
  value_revolution_event_id: string;
  industrial_capital_id: string;
  expected_realisation_pence: number;
  actual_realisation_pence: number;
  delta_pence: number;
}

export interface InventoryRevaluation {
  id: string;
  industrial_capital_id: string;
  commodity_id: string;
  quantity_held: number;
  old_unit_pence: number;
  new_unit_pence: number;
  delta_pence: number;
}

export interface SpeculativeHold {
  id: string;
  industrial_capital_id: string;
  commodity_id: string;
  held_since: string;
  anticipated_direction: number;
}

export interface PriceCaseRecord {
  value_revolution_event_id: string;
  fall_case?: string;
  rise_case?: string;
}

export interface PriceRevolutionRecord {
  event: ValueRevolutionEventV2;
  revalued_capital?: RevaluedCapital;
  realised_delta?: RealisedValueDelta;
  inventory_revaluations: InventoryRevaluation[];
  fall_case?: string;
  rise_case?: string;
}

export interface CompoundCapitalAdjustment {
  id: string;
  industrial_capital_id: string;
  period: string;
  revaluation_delta_pence: number;
  realisation_delta_pence: number;
  inventory_delta_pence: number;
  net_effect_pence: number;
}

export interface CreatePriceRevolutionInput {
  industrial_capital_id: string;
  affecting: ValueRevolutionAffecting;
  direction_pence: number;
  occurred_at?: string;
  basis_points_change: number;
  old_advance_pence?: number;
  expected_realisation_pence?: number;
  actual_realisation_pence?: number;
}

export interface RecordInventoryRevaluationInput {
  industrial_capital_id: string;
  commodity_id: string;
  quantity_held: number;
  old_unit_pence: number;
  new_unit_pence: number;
}

export interface RecordSpeculativeHoldInput {
  industrial_capital_id: string;
  commodity_id: string;
  held_since?: string;
  anticipated_direction: number;
}

// Vol. II Ch. 16 — The Turnover of Variable Capital

export interface VariableCapitalAdvance {
  id: string;
  industrial_capital_id: string;
  pence: number;
  turnover_period_minutes: number;
  turnovers_per_year_basis_points: number;
}

export interface PerCircuitSurplusRate {
  id: string;
  variable_capital_advance_id: string;
  surplus_pence: number;
  variable_pence: number;
  basis_points: number;
}

export interface AnnualSurplusRate {
  id: string;
  variable_capital_advance_id: string;
  period: string;
  n_basis_points: number;
  per_circuit_basis_points: number;
  annual_basis_points: number;
}

export interface AnnualSurplusRateContrast {
  id: string;
  capital_a_id: string;
  capital_b_id: string;
  capital_a_annual_basis_points: number;
  capital_b_annual_basis_points: number;
  ratio_basis_points: number;
}

export interface VariableCapitalReproduction {
  variable_capital_advance_id: string;
  year_reference_minutes: number;
  total_variable_paid_annual_pence: number;
  total_surplus_annual_pence: number;
}

export interface CreateVariableCapitalAdvanceInput {
  industrial_capital_id: string;
  pence: number;
  turnover_period_minutes: number;
  turnovers_per_year_basis_points: number;
}

// Vol. II Ch. 17 — The Circulation of Surplus-Value
export interface RealisationSourceEntry {
  surplus_circulation_id: string;
  source: string;
  pence: number;
  deferred_realisation_pence: number;
}

export interface SurplusCirculation {
  id: string;
  period: string;
  total_surplus_pence: number;
  realised_surplus_pence: number;
  unrealised_surplus_pence: number;
  realisation_sources: RealisationSourceEntry[];
}

export interface SocialCapitalAggregate {
  period: string;
  total_advanced_pence: number;
  total_annual_output_pence: number;
  total_annual_surplus_pence: number;
  department_i_share_bp: number;
  department_ii_share_bp: number;
}

export interface RealisationPuzzle {
  surplus_circulation_id: string;
  puzzle_statement: string;
  resolved_in_chapter: number;
}

// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction

export interface MoneySupplyApportionment {
  id: string;
  total_circulating_money_pence: number;
  department_i_reserve_pence: number;
  department_ii_reserve_pence: number;
  wage_rotation_fund_pence: number;
  idle_hoard_pence: number;
  period: string;
}

export interface DepartmentMoneyReserve {
  id: string;
  department: string;
  reserve_pence: number;
  purpose: string;
  period: string;
}

export interface CirculatingMoneyMass {
  id: string;
  money_stock_pence: number;
  velocity_per_year_basis_points: number;
  effective_circulating_value_pence: number;
  period: string;
}

export interface WageRotationFund {
  id: string;
  fund_pence: number;
  wage_cycle_frequency: number;
  department: string;
  period: string;
}

export interface InterDepartmentSettlement {
  id: string;
  from_department: string;
  to_department: string;
  amount_pence: number;
  purpose: string;
  period: string;
}

export interface ApportionmentBalanceCheck {
  balanced: boolean;
  apportionment_id: string;
}

// Vol. II Ch. 20 — Simple Reproduction

export interface DepartmentalCapital {
  id: string;
  scheme_id: string;
  department: string;
  constant_pence: number;
  variable_pence: number;
  surplus_pence: number;
  total_pence: number;
  turnover_number_bp: number;
}

export interface SimpleReproductionScheme {
  id: string;
  period: string;
  department_i?: DepartmentalCapital;
  department_ii?: DepartmentalCapital;
  is_balanced: boolean;
  tick_count: number;
}

export interface InterDepartmentExchange {
  id: string;
  scheme_id: string;
  from_department: string;
  to_department: string;
  pence: number;
  kind: string;
  description: string;
}

export interface ReproductionTick {
  id: string;
  scheme_id: string;
  tick_number: number;
  period: string;
  is_balanced: boolean;
  worker_pool_size: number;
  wage_bill_pence: number;
  subsistence_basket_pence: number;
  subsistence_covered: boolean;
  status: string;
}

export interface BalanceCheckResult {
  scheme_id: string;
  is_balanced: boolean;
  dept_i_v_plus_s_pence: number;
  dept_ii_constant_pence: number;
  deficit_pence: number;
  dept_i_turnover_bp: number;
  dept_ii_turnover_bp: number;
  dept_i_v_plus_s_annualised_pence: number;
  dept_ii_constant_annualised_pence: number;
  annualised_deficit_pence: number;
  annualised_balanced: boolean;
}

export interface MoneyClosedLoop {
  id: string;
  scheme_id: string;
  period: string;
  is_closed: boolean;
  net_flow_pence: number;
}

// Vol. II Ch. 21 — Extended Reproduction

export interface ExtendedReproductionScheme {
  id: string;
  period: string;
  department_i?: DepartmentalCapital;
  department_ii?: DepartmentalCapital;
  accumulation_rate_i_bps: number;
  accumulation_rate_ii_bps: number;
  is_balanced: boolean;
  tick_count: number;
}

export interface Reinvestment {
  id: string;
  scheme_id: string;
  department: string;
  delta_constant_pence: number;
  delta_variable_pence: number;
  consumed_surplus_pence: number;
  cycle_number: number;
}

export interface MultiPeriodScheme {
  id: string;
  label: string;
  scheme_ids: string[];
}

export interface CompositionShift {
  id: string;
  multi_period_id: string;
  cycle_number: number;
  dept_i_composition_bps: number;
  dept_ii_composition_bps: number;
  social_composition_bps: number;
}

export interface ExtendedMoneyRequirement {
  id: string;
  scheme_id: string;
  base_money_pence: number;
  additional_money_pence: number;
  total_money_pence: number;
}

export interface DepartmentIGrowthLead {
  id: string;
  multi_period_id: string;
  cycle_number: number;
  dept_i_growth_bps: number;
  dept_ii_growth_bps: number;
  lead_confirmed: boolean;
}

// Vol. III Ch. 1 — Cost-Price and Profit. Magnitudes are LabourMinutes; Marx's
// worked examples express them in £ (6s. = one social working-day).

export interface CapitalOutlayV3 {
  constant: number;
  variable: number;
  fixed_wear_and_tear: number;
  fixed_advanced: number;
}

export interface CostPriceResponse {
  id: string;
  outlay: CapitalOutlayV3;
  k: number;
  fixed_component: number;
  circulating_component: number;
  created_at: string;
}

export interface ProfitScenario {
  selling_price: number;
  cost_price: number;
  commodity_value: number;
  amount: number;
  below_value: boolean;
}

export interface ProfitFormResponse {
  cost_price: number;
  surplus_value: number;
  profit: number;
  commodity_value: number;
  mystifies_origin: boolean;
  selling_price_scenarios: ProfitScenario[];
}

// Vol. III Ch. 2 — The Rate of Profit. The same surplus-value measured against
// the total capital C = c + v gives the rate of profit p′; measured against the
// variable capital v alone it gives the rate of surplus-value s′. Both rates are
// basis points (10000 bp = 100%). Mystification is the gap (s′ − p′)/s′ in bp.

export interface RateOfProfitV3 {
  surplus_value: number;
  total_capital: number;
  basis_points: number;
}

export interface ProfitRateAnalysisResponse {
  id: string;
  constant_capital: number;
  variable_capital: number;
  profit_rate: RateOfProfitV3;
  surplus_value_rate: number;
  mystification: number;
  created_at: string;
}

// Vol. III Ch. 3 — Relation of the Rate of Profit to the Rate of Surplus-Value.
// Marx's master equation p′ = s′·(v/C) resolves the rate of profit into two
// factors: the rate of surplus-value s′ and the value-composition v/C. In the
// formula struct, c is the constant component and v the variable component, so
// the total capital C = c + v. s_rate and every *_rate / *_bp field is in basis
// points (10000 bp = 100%).

export type VariationCase =
  | "s_constant_vc_variable"
  | "s_variable_vc_constant"
  | "both_variable";

export interface ProfitRateFormulaV3 {
  c: number;
  v: number;
  s_rate: number;
}

export interface CompositionRatioV3 {
  variable: number;
  total: number;
  basis_points: number;
}

export interface VariationAnalysisResponse {
  id: string;
  case: VariationCase;
  initial: ProfitRateFormulaV3;
  changed: ProfitRateFormulaV3;
  old_profit_rate: number;
  new_profit_rate: number;
  old_composition_bp: number;
  new_composition_bp: number;
  proportional_change: boolean;
  created_at: string;
}

export interface ProfitRateComparisonResponse {
  first: ProfitRateFormulaV3;
  second: ProfitRateFormulaV3;
  first_profit_rate: number;
  second_profit_rate: number;
  first_composition: CompositionRatioV3;
  second_composition: CompositionRatioV3;
  equal_surplus_rate: boolean;
}

// Vol. III Ch. 4 — The Effect of the Turnover on the Rate of Profit. Turnover
// corrects the rate of profit: across a year the variable capital yields its
// surplus once for every turnover, so the annual rate of profit is
// p′ = s′·n·(v/C). c is the total capital C, v the variable capital advanced (the
// v in C), s_rate is s′ and every *_rate / basis_points field is in basis points
// (10000 bp = 100%). n is the whole number of turnovers; annual_wages is the
// year's wage bill vn — what the capitalist's books show in place of v.

export interface AnnualProfitRateV3 {
  surplus_value_rate: number;
  turnovers: number;
  variable_capital: number;
  total_capital: number;
  basis_points: number;
}

export interface TurnoverAnalysisResponse {
  id: string;
  c: number;
  v: number;
  s_rate: number;
  n: number;
  annual_profit_rate: AnnualProfitRateV3;
  single_turnover_profit_rate: number;
  annual_surplus_value_rate: number;
  annual_wages: number;
  created_at: string;
}

// Vol. III Ch. 5 — Economy in the Employment of Constant Capital. A saving in
// constant capital raises the rate of profit s/(c+v) without touching s or v:
// the denominator shrinks while the surplus stays put. The economy takes various
// forms (the kind), and the analysis records the old rate, the new rate, and the
// gain between them, all in basis points (10000 bp = 100%). Magnitudes are
// LabourMinutes, expressed in £ in Marx's examples.

export type EconomyKind =
  | "shared_conditions"
  | "waste_reduction"
  | "workday_extension"
  | "improved_machinery"
  | "raw_material_saving";

export interface ConstantCapitalEconomyV3 {
  kind: EconomyKind;
  constant_capital: number;
  variable_capital: number;
  surplus_value: number;
  saving: number;
}

export interface EconomyAnalysisResponse {
  id: string;
  economy: ConstantCapitalEconomyV3;
  old_profit_rate: number;
  new_profit_rate: number;
  profit_rate_gain: number;
  created_at: string;
}

// §IV — conditions of production used in common: solo_cost is the aggregate the
// participants bear providing the condition each for itself; shared_cost is the
// single shared facility serving them all. The saving per participant is
// (solo_cost − shared_cost) / participant_count.
export interface SharedConditionResponse {
  kind: EconomyKind;
  solo_cost: number;
  shared_cost: number;
  participant_count: number;
}

// §IV — waste reclaimed and returned to production. The constant-capital saving
// is recovered_waste × price_per_unit.
export interface WasteReductionResponse {
  kind: EconomyKind;
  original_waste: number;
  recovered_waste: number;
  price_per_unit: number;
}

// Vol. III Ch. 6 — The Effect of Price Fluctuation on the Rate of Profit. The
// rate of profit p' = s/(c+v) moves with the price of the raw material that makes
// up the circulating part of c, though that price-movement adds no surplus-value.
// Only the circulating raw-material component fluctuates; the fixed component and
// v and s are held fixed. price_factor is the new raw-material price as a
// percentage of its original level (100 = unchanged, 200 = doubled, 50 = halved).
// Rates are basis points (10000 bp = 100%); magnitudes are LabourMinutes (£).

export type PriceFluctuationKind = "rise" | "fall";

export interface ConstantCapitalPriceEffectV3 {
  fixed_capital: number;
  original_material_capital: number;
  price_factor: number;
  variable_capital: number;
  surplus_value: number;
}

export interface PriceFluctuationAnalysisResponse {
  id: string;
  kind: PriceFluctuationKind;
  effect: ConstantCapitalPriceEffectV3;
  old_profit_rate: number;
  new_profit_rate: number;
  created_at: string;
}

// §II — the capital tied up in a stock of raw material. A price fall releases
// capital previously advanced to hold the stock; a price rise ties up more. For a
// given movement, released and tied-up are the same magnitude with opposite sign:
// stock_months × stock_units × (price delta).
export interface CapitalReleaseResponse {
  stock_months: number;
  stock_units: number;
  original_price_per_unit: number;
  new_price_per_unit: number;
}

// Vol. III Ch. 7 — Supplementary Remarks (on the rate of profit). Two threads:
// the organic-composition channel (two capitals of equal s and v but different c
// show different rates of profit) and the refutation of Rodbertus (a nominal
// money-value change, or a proportional change at constant composition, leaves
// the rate untouched; only an organic change moves it). Rates are basis points
// (10000 bp = 100%); magnitudes are LabourMinutes (£).

export interface OrganicCompositionEffectV3 {
  s_a: number;
  v_a: number;
  c_a: number;
  s_b: number;
  v_b: number;
  c_b: number;
}

export interface CompositionEffectAnalysisResponse {
  id: string;
  effect: OrganicCompositionEffectV3;
  profit_rate_a: number;
  profit_rate_b: number;
  rate_difference: number;
  created_at: string;
}

export type MagnitudeChangeKind =
  | "money_value_change"
  | "proportional_change"
  | "organic_change";

export interface CapitalMagnitudeChangeV3 {
  kind: MagnitudeChangeKind;
  original_capital: number;
  original_profit: number;
  factor: number;
}

export interface MagnitudeChangeAnalysisResponse {
  id: string;
  change: CapitalMagnitudeChangeV3;
  old_rate: number;
  new_rate: number;
  rate_unchanged: boolean;
  created_at: string;
}

// The Part I summary names the determinants of the rate of profit established
// across Chs. 1–6 and gathers the rate-of-profit analyses recorded so far.
export interface PartISummaryResponse {
  id: string;
  determinants: string[];
  analyses: ProfitRateAnalysisResponse[];
  created_at: string;
}

// Vol. III Ch. 8 — Different Compositions of Capitals in Different Branches.
// ProductionSphereResponse mirrors avgprofit.ProductionSphere. Rates in basis
// points (10000 = 100%); capital magnitudes in LabourMinutes.
export interface ProductionSphereResponse {
  id: string;
  name: string;
  c: number;
  v: number;
  s_rate: number;
  constant_capital: number;
  surplus_value: number;
  individual_profit_rate: number;
  variable_percent: number;
  labour_power_index: number;
  created_at: string;
}

// OrganicCompositionResponse mirrors avgprofit.OrganicComposition with derived
// percentage fields for rendering.
export interface OrganicCompositionResponse {
  constant: number;
  variable: number;
  constant_percent: number;
  variable_percent: number;
}

// Vol. III Ch. 9 — Formation of a General Rate of Profit.

// GeneralProfitRateResponse mirrors avgprofit.GeneralProfitRate.
export interface GeneralProfitRateResponse {
  id: string;
  spheres: ProductionSphereResponse[];
  rate: number;                    // p̄′ in basis points
  sum_surplus_values: number;
  sum_total_capitals: number;
  average_variable_percent: number;
  created_at: string;
}

// PriceOfProductionResponse mirrors avgprofit.PriceOfProduction.
export interface PriceOfProductionResponse {
  id: string;
  sphere_name: string;
  cost_price: number;
  general_rate: number;
  commodity_value: number;
  average_profit: number;
  price: number;
  deviation: number;               // signed: price − value
  composition: string;             // "higher" | "lower" | "average" | ""
  created_at: string;
}

// ValuePriceDeviationResponse mirrors avgprofit.ValuePriceDeviation.
export interface ValuePriceDeviationResponse {
  value: number;
  price_of_production: number;
  deviation: number;               // signed: price − value
}

// SocialCapitalAggregateResponse mirrors avgprofit.SocialCapitalAggregate.
export interface SocialCapitalAggregateResponse {
  spheres: ProductionSphereResponse[];
  weighted_general_rate: number;
  sum_surplus_values: number;
  sum_total_capitals: number;
  average_variable_percent: number;
  prices_of_production: PriceOfProductionResponse[];
  total_values: number;
  total_prices_of_production: number;
  sum_average_profits: number;
  value_conserved: boolean;
}

// MarketValueResponse mirrors avgprofit.MarketValue.
export interface MarketValueResponse {
  id: string;
  sphere_name: string;
  bulk_condition_value: number;
  best_condition_value: number;
  worst_condition_value: number;
  value: number; // == bulk_condition_value
  created_at: string;
}

// SurplusProfitResponse mirrors avgprofit.SurplusProfit.
export interface SurplusProfitResponse {
  id: string;
  firm_name: string;
  individual_value: number;
  market_value: number;
  output_qty: number;
  general_rate: number;
  amount: number; // signed: positive = extra profit, negative = shortfall
  created_at: string;
}

// CapitalFlowResponse mirrors avgprofit.CapitalFlow.
export interface CapitalFlowResponse {
  sphere_name: string;
  current_sphere_rate: number;
  general_rate: number;
  market_price: number;
  market_value: number;
  direction: "inflow" | "outflow" | "stable";
}

// EqualisationResponse mirrors avgprofit.Equalisation.
export interface EqualisationResponse {
  id: string;
  sphere_name: string;
  initial_rate: number;
  target_rate: number;
  direction: "inflow" | "outflow" | "stable";
  is_converging: boolean;
  flow: CapitalFlowResponse;
  created_at: string;
}

// SphereWageOutcomeResponse mirrors avgprofit.SphereWageOutcome.
// old_price, new_price, price_change are in centi-units (÷100 for display).
export interface SphereWageOutcomeResponse {
  sphere_name: string;
  constant_capital: number;
  variable_capital: number;
  wage_factor: number;
  composition: "higher" | "average" | "lower";
  old_price: number;    // centi-units
  new_price: number;    // centi-units
  price_change: number; // signed centi-units
  direction: "rises" | "falls" | "unchanged";
}

// WageEffectAnalysisResponse mirrors avgprofit.WageEffectAnalysis.
export interface WageEffectAnalysisResponse {
  id: string;
  base_constant: number;
  base_variable: number;
  s_rate: number;
  wage_factor: number;
  kind: "rise" | "fall" | "";
  old_general_rate: number; // basis points
  new_general_rate: number; // basis points
  average_outcome: SphereWageOutcomeResponse;
  outcomes: SphereWageOutcomeResponse[];
  created_at: string;
}

// --- finance-service types (Vol. III Ch. 12: Supplementary Remarks) -----------

// PriceOfProductionChangeResponse mirrors avgprofit.PriceOfProductionChange.
export interface PriceOfProductionChangeResponse {
  id: string;
  sphere_name: string;
  cause: "general_rate_change" | "individual_value_change" | "both" | "";
  rate_changed: boolean;
  value_changed: boolean;
  price_changed: boolean;
  old_price: number; // whole units
  new_price: number; // whole units
  created_at: string;
}

// CompensationGroundResponse mirrors avgprofit.CompensationGround.
export interface CompensationGroundResponse {
  reason: string;
  appears_to_be: string;
  actually_is: string; // always "share in aggregate surplus-value"
  creates_value: boolean; // always false
}

// PartIISummaryResponse mirrors avgprofit.PartIISummary.
export interface PartIISummaryResponse {
  id: string;
  general_rate: number; // basis points
  sphere_count: number;
  prices_of_production_count: number;
  wage_effect_count: number;
  total_values: number;
  total_prices_of_production: number;
  value_conserved: boolean;
  average_composition_check: {
    total_surplus_value: number;
    total_average_profit: number;
    surplus_value_equals_average_profit: boolean;
  };
  created_at: string;
}

// --- finance-service types (Vol. III Ch. 13: The Law as Such) ----------------

// TrajectoryPeriodResponse mirrors tendency.TrajectoryPeriod.
export interface TrajectoryPeriodResponse {
  constant_capital: number; // c; LabourMinutes
  variable_capital: number; // v; LabourMinutes
  surplus_value: number; // s; LabourMinutes
  profit_rate: number; // p' in basis points
}

// CompositionTrajectoryResponse mirrors tendency.CompositionTrajectory.
export interface CompositionTrajectoryResponse {
  id: string;
  label: string;
  surplus_value_rate: number; // s' in basis points
  periods: TrajectoryPeriodResponse[];
  profit_rates: number[]; // derived basis-point rates
  created_at: string;
}

// RateMassContradictionResponse mirrors tendency.RateMassContradiction.
export interface RateMassContradictionResponse {
  id: string;
  old_c: number;
  old_rate: number; // basis points
  new_c: number;
  new_rate: number; // basis points
  old_mass: number;
  new_mass: number;
  mass_change: number; // signed
  created_at: string;
}

// --- finance-service types (Vol. III Ch. 14: Counteracting Influences) --------

// CounteractingForceKind mirrors tendency.CounteractingForceKind.
export type CounteractingForceKind =
  | "intensified_exploitation"
  | "depressed_wages"
  | "cheaper_constant_capital"
  | "relative_overpopulation"
  | "foreign_trade"
  | "stock_capital";

// CounteractingForceResponse mirrors tendency.CounteractingForce.
export interface CounteractingForceResponse {
  id: string;
  kind: CounteractingForceKind;
  excluded_from_general_rate: boolean;
  note: string;
  created_at: string;
}

// CounteractingScenarioResponse mirrors tendency.CounteractingScenario.
export interface CounteractingScenarioResponse {
  id: string;
  label: string;
  forces: CounteractingForceResponse[];
  base_trajectory: CompositionTrajectoryResponse;
  initial_profit_rate: number; // basis points
  final_profit_rate: number; // basis points
  modified_rates: number[]; // basis points per period
  created_at: string;
}

// --- finance-service types (Vol. III Ch. 15: Internal Contradictions of the Law) ---

// ContradictionKind mirrors tendency.ContradictionKind.
export type ContradictionKind =
  | "falling_rate_rising_mass"
  | "overproduction_underconsumption"
  | "overaccumulation_relative_overpopulation"
  | "concentration_centralisation";

// InternalContradictionResponse mirrors tendency.InternalContradiction.
export interface InternalContradictionResponse {
  id: string;
  kind: ContradictionKind;
  is_coexistent: boolean;
  note: string;
  created_at: string;
}

// CrisisResponse mirrors tendency.Crisis.
export interface CrisisResponse {
  id: string;
  constant_capital_writedown: number; // percent integer 1-99
  pre_crisis_profit_rate: number; // basis points
  post_crisis_profit_rate: number; // basis points; always > pre
  created_at: string;
}

// CapitalConcentrationResponse mirrors tendency.CapitalConcentration (value object).
export interface CapitalConcentrationResponse {
  initial_c: number; // LabourMinutes
  surplus_value_rate: number; // basis points
  periods: number;
  final_c: number; // derived
  path: number[]; // per-period values including initial_c as [0]
}

// CapitalCentralisationResponse mirrors tendency.CapitalCentralisation (value object).
export interface CapitalCentralisationResponse {
  absorbed_capitals: number; // count
  absorbed_size: number; // LabourMinutes each
  resulting_capital: number; // derived
}

// PartIIISummaryResponse mirrors tendency.PartIIISummary.
export interface PartIIISummaryResponse {
  trajectory_count: number;
  scenario_count: number;
  crisis_count: number;
  contradiction_count: number;
  latest_trajectory?: CompositionTrajectoryResponse;
  latest_scenario?: CounteractingScenarioResponse;
  latest_crisis?: CrisisResponse;
  latest_contradiction?: InternalContradictionResponse;
}

// --- finance-service types (Vol. III Ch. 16: Commercial Capital) --------------

// CommercialCapitalFunction mirrors merchant.CommercialCapitalFunction.
// 0=realisation, 1=storage, 2=transport, 3=distribution.
export type CommercialCapitalFunction = 0 | 1 | 2 | 3;

// CommercialCapitalResponse mirrors merchant.CommercialCapital.
export interface CommercialCapitalResponse {
  id: string;
  money_advanced: number; // pence; >0
  commodity_description: string;
  function: CommercialCapitalFunction;
  surplus_value_produced: number; // always 0 — commercial capital produces no surplus-value
  created_at: string;
}

// MerchantCircuitResponse mirrors merchant.MerchantCircuit (value object).
export interface MerchantCircuitResponse {
  money_advanced: number; // pence
  commodity_value: number; // pence
  money_returned: number; // pence; must exceed money_advanced
}

// ── Vol. III Ch. 17 — Commercial Profit ─────────────────────────────────────

// CommercialProfitResponse mirrors merchant.CommercialProfit.
export interface CommercialProfitResponse {
  id: string;
  productive_capital: number; // pence; >0
  commercial_capital: number; // pence; >0
  surplus_value: number;      // pence; >=0
  unadjusted_rate_bp: number; // s/productive_capital in bp, round-half-up
  adjusted_rate_bp: number;   // s/(productive+commercial) in bp — always <= unadjusted
  new_value_created: number;  // always 0
  created_at: string;
}

// AdjustedGeneralRateResponse mirrors merchant.AdjustedGeneralRate (value object).
export interface AdjustedGeneralRateResponse {
  productive_capital: number;
  commercial_capital: number;
  surplus_value: number;
  unadjusted_rate_bp: number;
  adjusted_rate_bp: number;
}

// CommercialWorkerResponse mirrors merchant.CommercialWorker (value object).
export interface CommercialWorkerResponse {
  necessary_labour_minutes: number; // LabourMinutes; >0
  workday_minutes: number;          // LabourMinutes; >=necessary
  unpaid_labour_minutes: number;    // = workday - necessary
  new_value_created: number;        // always 0
}

// WholesaleRetailSpreadResponse mirrors merchant.WholesaleRetailSpread (value object).
export interface WholesaleRetailSpreadResponse {
  wholesale_price: number; // pence; >0
  retail_price: number;    // pence; >wholesale
  spread_pence: number;    // = retail - wholesale
}

// MerchantTurnoverResponse mirrors merchant.MerchantTurnover (Vol. III Ch. 18).
export interface MerchantTurnoverResponse {
  id: string;
  money_advanced: number;          // pence; >0
  general_rate_bp: number;         // bp; >0
  turnover_count: number;          // count; >=1
  units_per_turnover: number;      // count; >0
  annual_merchant_profit: number;  // pence; derived = roundHalfUp(money_advanced*general_rate_bp, 10000)
  markup_per_unit: number;         // pence; derived (integer div) — inversely proportional to turnover
  created_at: string;
}

// TurnoverPointResponse mirrors merchant.TurnoverPoint (value object).
export interface TurnoverPointResponse {
  turnover_count: number;  // count
  total_units: number;     // = turnover_count * units_per_turnover
  markup_per_unit: number; // pence; annual_merchant_profit / total_units (integer div)
}

// TurnoverEffectAnalysisResponse mirrors merchant.TurnoverEffectAnalysis (value object, not persisted).
export interface TurnoverEffectAnalysisResponse {
  money_advanced: number;
  general_rate_bp: number;
  annual_merchant_profit: number;
  points: TurnoverPointResponse[];
}

// MerchantPriceOfProductionResponse mirrors merchant.MerchantPriceOfProduction (value object).
export interface MerchantPriceOfProductionResponse {
  cost_price: number;          // pence; >=0
  average_profit: number;      // pence; >=0
  price_of_production: number; // pence; = cost_price + average_profit
}

// MoneyOperationKind mirrors merchant.MoneyOperationKind (Vol. III Ch. 19).
// 1=Receipt, 2=Payment, 3=Exchange, 4=Safekeeping, 5=Bookkeeping.
export type MoneyOperationKind = 1 | 2 | 3 | 4 | 5;

// MoneyDealingCapitalResponse mirrors merchant.MoneyDealingCapital (Vol. III Ch. 19).
// MoneyDealingProfit = roundHalfUp(money_advanced * general_rate_bp, 10000).
// NewValueCreated is always 0 — money-dealing capital appropriates but does not create surplus-value.
export interface MoneyDealingCapitalResponse {
  id: string;
  money_advanced: number;       // pence; > 0
  general_rate_bp: number;      // bp; > 0
  operation_kinds: MoneyOperationKind[]; // each 1–5
  money_dealing_profit: number; // pence; derived
  new_value_created: number;    // always 0
  created_at: string;
}

// MoneyDealingCircuitResponse mirrors merchant.MoneyDealingCircuit (value object, M—M′ circuit).
export interface MoneyDealingCircuitResponse {
  money_advanced: number; // pence; > 0
  money_returned: number; // pence; > money_advanced
}

// DevelopmentStage mirrors merchant.DevelopmentStage (Vol. III Ch. 20).
// 1=pre_capitalist, 2=transition, 3=subordinated.
export type DevelopmentStage = 1 | 2 | 3;

// HistoricalMerchantCapitalResponse mirrors merchant.HistoricalMerchantCapital (Vol. III Ch. 20).
// Invariant 1: pre-capitalist forms (stage=1) cannot employ wage-labour.
// Invariant 2 (inverse-development law): SubordinationIndex rises with Stage.
export interface HistoricalMerchantCapitalResponse {
  id: string;
  name: string;
  stage: DevelopmentStage;
  profit_source: string;
  wage_labour: boolean;
  subordination_index: number; // bp [0, 10000]
  created_at: string;
}

// HistoricalTradingSystemResponse mirrors merchant.HistoricalTradingSystem (Vol. III Ch. 20).
// Named historical exemplar (Venice/Genoa/Dutch) — value object, not persisted.
export interface HistoricalTradingSystemResponse {
  name: string;
  stage: DevelopmentStage;
  profit_source: string;
  wage_labour: boolean;
}

// InterestBearingCapital mirrors credit.InterestBearingCapital (Vol. III Ch. 21).
// Interest-bearing capital advances money at an interest rate for a loan term.
// InterestEarned = roundHalfUp(money_advanced * interest_rate_bp, 10000).
// MoneyReturned = money_advanced + interest_earned.
// new_value_created is always 0 — interest-bearing capital creates no new value.
export interface InterestBearingCapital {
  id: string;
  money_advanced: number;    // pence; > 0
  interest_rate_bp: number;  // bp; > 0
  loan_term_days: number;    // > 0
  interest_earned: number;   // pence; derived
  money_returned: number;    // pence; derived
  new_value_created: number; // always 0
  created_at: string;
}

// MoneyCapitalCircuit mirrors credit.MoneyCapitalCircuit (Vol. III Ch. 21).
// The M—M′ circuit of interest-bearing capital. Invariant: money_returned == money_advanced + interest_earned.
export interface MoneyCapitalCircuit {
  money_advanced: number;  // pence; > 0
  interest_earned: number; // pence; >= 0
  money_returned: number;  // pence; = money_advanced + interest_earned
}

// RateOfInterest mirrors credit.RateOfInterest (Vol. III Ch. 22).
// Observed market rate of interest in basis points. RateBP < AverageProfitRateBP.
// CyclePhase: 1=Inactivity, 2=Revival, 3=Prosperity, 4=Overproduction, 5=Crisis.
export interface RateOfInterest {
  id: string;
  rate_bp: number;                // bp; >= 0
  average_profit_rate_bp: number; // bp; > 0
  cycle_phase: number;            // 1..5
  period: string;
  created_at: string;
}

// InterestRateAnalysis mirrors credit.InterestRateAnalysis (Vol. III Ch. 22).
// Compute-only split of average profit into interest and profit of enterprise.
// profit_of_enterprise_bp == average_profit_rate_bp - interest_rate_bp.
export interface InterestRateAnalysis {
  average_profit_rate_bp: number;  // bp; > 0
  interest_rate_bp: number;        // bp; >= 0
  profit_of_enterprise_bp: number; // bp; = avg - interest
}

// InterestRateBound mirrors credit.InterestRateBound (Vol. III Ch. 22).
// Theoretical bounds: max == average profit rate; min approaches but never reaches 0.
export interface InterestRateBound {
  max_bp: number; // bp; > 0
  min_bp: number; // bp; >= 0
}

// WagesOfSuperintendence mirrors credit.WagesOfSuperintendence (Vol. III Ch. 23).
// The labour component of the functioning capitalist's income in large joint-stock
// companies; their income = wages of skilled labour, not profit of enterprise.
export interface WagesOfSuperintendence {
  labour_minutes: number; // canonical value-magnitude unit; >= 0
}

// ProfitDivision mirrors credit.ProfitDivision (Vol. III Ch. 23).
// Splits total profit (= average profit) into interest and profit of enterprise.
// Invariant: total_profit_bp == interest_bp + profit_of_enterprise_bp.
// ownership_form: 1=OwnerOperator, 2=Separated.
export interface ProfitDivision {
  id: string;
  total_profit_bp: number;             // bp; > 0
  interest_bp: number;                 // bp; >= 0
  profit_of_enterprise_bp: number;     // bp; = total - interest; >= 0
  ownership_form: number;              // 1=OwnerOperator, 2=Separated
  wages_of_superintendence: WagesOfSuperintendence;
  mystification_index: number;         // bp [0,10000]
  created_at: string;
}

// CompoundInterestSchedule mirrors credit.CompoundInterestSchedule (Vol. III Ch. 24).
// FinalValue is derived by integer compound iteration (no float64).
// new_value_created is always 0 — the fetish form.
// fetish_form: 1=SimpleInterest, 2=CompoundInterest, 3=MM.
export interface CompoundInterestSchedule {
  id: string;
  principal: number;           // > 0; base amount
  rate_bp: number;             // > 0; annual interest rate in basis points
  period_years: number;        // > 0; number of compounding periods
  final_value: number;         // DERIVED: integer compound
  new_value_created: number;   // always 0 — fetish form
  fetish_form: number;         // 1=SimpleInterest, 2=CompoundInterest, 3=MM
  created_at: string;
}

// BillOfExchange mirrors credit.BillOfExchange (Vol. III Ch. 25).
// A promissory note with a fixed maturity, circulating as commercial money
// before its due date. face_value and maturity_days must be > 0.
export interface BillOfExchange {
  id: string;
  drawer: string;              // issuer
  drawee: string;              // debtor
  face_value: number;          // > 0; pence
  maturity_days: number;       // > 0
  description: string;
  created_at: string;
}

// FictitiousCapital mirrors credit.FictitiousCapital (Vol. III Ch. 25).
// A capitalised income stream: capitalised_value is DERIVED server-side as
// roundHalfUp(annual_income * 10000, rate_bp). For public debt (kind=3) the
// original capital is consumed, so real_capital_exists is false.
// kind: 1=bill, 2=bank-note, 3=public-debt, 4=share.
export interface FictitiousCapital {
  id: string;
  kind: number;                // 1=bill, 2=bank-note, 3=public-debt, 4=share
  name: string;
  annual_income: number;       // > 0; pence
  rate_bp: number;             // > 0; prevailing interest rate in basis points
  capitalised_value: number;   // DERIVED: roundHalfUp(annual_income*10000, rate_bp)
  real_capital_exists: boolean; // false for public debt (consumed capital)
  description: string;
  created_at: string;
}

// ── Vol. III Ch. 26 — Accumulation of Money-Capital ──────────────────────────

// InterestCycleObservation mirrors credit.InterestCycleObservation (Vol. III Ch. 26).
// Records an interest-rate reading tied to a phase of the industrial cycle.
// phase: 1=Inactivity, 2=Revival, 3=Prosperity, 4=Overproduction, 5=Crisis.
export interface InterestCycleObservation {
  phase: number;             // 1..5
  interest_rate_bp: number;  // bp; >= 0
  period: string;
}

// LoanableCapitalSupply mirrors credit.LoanableCapitalSupply (Vol. III Ch. 26).
// Records the aggregate supply of loanable money-capital at the observed phase.
export interface LoanableCapitalSupply {
  amount: number;    // pence; >= 0
  abundant: boolean; // true when supply exceeds demand
}

// IdleIndustrialCapital mirrors credit.IdleIndustrialCapital (Vol. III Ch. 26).
// Portion of industrial capital temporarily idle as money-capital.
export interface IdleIndustrialCapital {
  amount: number;      // pence; >= 0
  description: string;
}

// AccumulationGap mirrors credit.AccumulationGap (Vol. III Ch. 26).
// Signed difference between accumulation rate of money-capital and productive capital.
export interface AccumulationGap {
  gap_bp: number;      // bp; may be positive or negative
  description: string;
}

// MoneyCapitalAccumulation mirrors credit.MoneyCapitalAccumulation (Vol. III Ch. 26).
// Persisted entity for an observation of money-capital accumulation across
// the industrial cycle. Invariants: phase ∈ {1..5}, interest_rate_bp >= 0,
// idle_capital.amount >= 0, loanable_supply.amount >= 0.
export interface MoneyCapitalAccumulation {
  id: string;
  name: string;
  cycle_observation: InterestCycleObservation;
  loanable_supply: LoanableCapitalSupply;
  idle_capital: IdleIndustrialCapital;
  gap: AccumulationGap;
  description: string;
  created_at: string;
}

// CreateMoneyCapitalAccumulationInput is the flat request body for
// POST /v1/credit/money-capital-accumulation.
export interface CreateMoneyCapitalAccumulationInput {
  name: string;
  phase: number;
  interest_rate_bp: number;
  period: string;
  loanable_amount: number;
  loanable_abundant: boolean;
  idle_amount: number;
  idle_description: string;
  gap_bp: number;
  gap_description: string;
  description: string;
}

// ── Vol. III Ch. 27 — The Role of Credit in Capitalist Production ─────────────

// Manager mirrors credit.Manager (Vol. III Ch. 27).
// wage_minutes > 0 — wages of superintendence, NOT profit of enterprise.
export interface Manager {
  name: string;
  title: string;
  wage_minutes: number; // > 0
}

// StockCompany mirrors credit.StockCompany (Vol. III Ch. 27).
// Invariant: total_capital == fixed_capital + floating_capital.
export interface StockCompany {
  id: string;
  name: string;
  fixed_capital: number;    // pence
  floating_capital: number; // pence
  total_capital: number;    // pence; == fixed + floating
  manager: Manager;
  description: string;
  created_at: string;
}

// CreateStockCompanyInput is the request body for POST /v1/credit/stock-companies.
export interface CreateStockCompanyInput {
  name: string;
  fixed_capital: number;
  floating_capital: number;
  total_capital: number;
  manager_name: string;
  manager_title: string;
  manager_wage_minutes: number;
  description: string;
}

// CooperativeFactory mirrors credit.CooperativeFactory (Vol. III Ch. 27).
// Invariant: worker_owned == true.
export interface CooperativeFactory {
  id: string;
  name: string;
  worker_owned: boolean;    // always true
  worker_count: number;
  capital_advanced: number; // pence
  description: string;
  created_at: string;
}

// CreateCooperativeFactoryInput is the request body for POST /v1/credit/cooperative-factories.
export interface CreateCooperativeFactoryInput {
  name: string;
  worker_owned: boolean;
  worker_count: number;
  capital_advanced: number;
  description: string;
}

// ── Vol. III Ch. 28 — Medium of Circulation and Capital (Tooke and Fullarton) ─

// ReserveFund mirrors credit.ReserveFund (Vol. III Ch. 28).
// amount >= 0 — cash reserve held against convertibility demands.
export interface ReserveFund {
  amount: number;        // pence; >= 0
  description: string;
}

// CurrencyObservation mirrors credit.CurrencyObservation (Vol. III Ch. 28).
// Invariants:
//   total_currency_outstanding > 0
//   coin_function_amount + capital_transfer_amount <= total_currency_outstanding
//   reserve_fund.amount >= 0
export interface CurrencyObservation {
  id: string;
  name: string;
  total_currency_outstanding: number; // pence; > 0
  coin_function_amount: number;       // pence; revenue expenditure (coin function)
  capital_transfer_amount: number;    // pence; productive-capital transfers
  reserve_fund: ReserveFund;
  description: string;
  created_at: string;
}

// CreateCurrencyObservationInput is the flat request body for
// POST /v1/credit/currency-observations.
export interface CreateCurrencyObservationInput {
  name: string;
  total_currency_outstanding: number;
  coin_function_amount: number;
  capital_transfer_amount: number;
  reserve_amount: number;
  reserve_description: string;
  description: string;
}

// ── Vol. III Ch. 30 — Money-Capital and Real Capital, I ──────────────────────

// AccumulationCorrespondence mirrors credit.AccumulationCorrespondence (Vol. III Ch. 30).
// The two growth rates MAY diverge — there is no equality invariant; either may be
// negative (the crisis configuration: money-capital grows while real capital shrinks).
export interface AccumulationCorrespondence {
  money_capital_growth_bp: number; // bp; may be positive or negative
  real_capital_growth_bp: number;  // bp; may be positive or negative
  corresponds: boolean;            // true when the two accumulations align
  explanation: string;
}

// CommercialCreditCircuit mirrors credit.CommercialCreditCircuit (Vol. III Ch. 30).
// A value type — the chain of bills producers draw on one another (spinner →
// weaver → manufacturer → exporter). total_bills_amount > 0.
export interface CommercialCreditCircuit {
  id: string;
  stages: string[];
  total_bills_amount: number; // pence; > 0
  description: string;
}

// RealCapitalAccumulation mirrors credit.RealCapitalAccumulation (Vol. III Ch. 30).
// Invariants:
//   phase ∈ {1,2,3,4,5}
//   credit_limit.reserve_capital >= 0
export interface RealCapitalAccumulation {
  id: string;
  name: string;
  phase: number; // 1..5
  accumulation_correspondence: AccumulationCorrespondence;
  credit_limit: {
    reserve_capital: number; // pence; >= 0
    description: string;
  };
  commercial_credit_contracts: boolean;
  cause_is_money_scarcity: boolean;
  description: string;
  created_at: string;
}

// CreateRealCapitalAccumulationInput is the flat request body for
// POST /v1/credit/real-capital-accumulation.
export interface CreateRealCapitalAccumulationInput {
  name: string;
  phase: number;
  money_capital_growth_bp: number;
  real_capital_growth_bp: number;
  corresponds: boolean;
  correspondence_explanation: string;
  reserve_capital: number;
  credit_limit_description: string;
  commercial_credit_contracts: boolean;
  cause_is_money_scarcity: boolean;
  description: string;
}

// ── Vol. III Ch. 31 — Money-Capital and Real Capital, II ─────────────────────

// LoanCapitalVelocity mirrors credit.LoanCapitalVelocity (Vol. III Ch. 31).
// The same money lent repeatedly performs a larger mass of loan capital:
// loan_capital_created == money_amount × times_lent (the same £20 lent 5× = £100).
export interface LoanCapitalVelocity {
  money_amount: number;          // pence
  times_lent: number;            // >= 1
  loan_capital_created: number;  // = money_amount × times_lent
}

// FloatingCapital mirrors credit.FloatingCapital (Vol. III Ch. 31).
// source: 1 = simple transformation of idle money into loan capital (inverse to
//   real accumulation); 2 = conversion of capital/revenue into money that becomes
//   loan capital (accompanies real accumulation).
// Invariants: amount >= 0; loan_capital_velocity.times_lent >= 1;
//   circulation_reserve_saving.amount >= 0.
export interface FloatingCapital {
  id: string;
  name: string;
  amount: number; // pence; >= 0
  source: number; // 1 or 2
  loan_capital_velocity: LoanCapitalVelocity;
  circulation_reserve_saving: {
    amount: number; // pence; >= 0
    description: string;
  };
  description: string;
  created_at: string;
}

// CreateFloatingCapitalInput is the flat request body for
// POST /v1/credit/floating-capital. The velocity and reserve-saving sub-objects
// are flattened; loan_capital_created is derived server-side.
export interface CreateFloatingCapitalInput {
  name: string;
  amount: number;
  source: number;
  velocity_money_amount: number;
  velocity_times_lent: number;
  reserve_saving_amount: number;
  reserve_saving_description: string;
  description: string;
}

// ── Vol. III Ch. 32 — Money-Capital and Real Capital, III (Conclusion) ────────

// LoanCapitalOverstatement mirrors credit.LoanCapitalOverstatement (Vol. III Ch. 32).
// The gap between apparent (loanable money-capital) and real (productive)
// accumulation, in basis points. Always >= 0.
export interface LoanCapitalOverstatement {
  overstatement_pct: number; // basis points; >= 0
  description: string;
}

// MeansOfPaymentDemand mirrors credit.MeansOfPaymentDemand (Vol. III Ch. 32).
// In a crisis (phase 5) the demand for loan capital is a demand for means of
// payment (convertibility), not productive capital: is_productive_capital_demand
// is false.
export interface MeansOfPaymentDemand {
  id: string;
  phase: number; // 1..5
  amount: number; // pence
  is_productive_capital_demand: boolean;
  description: string;
}

// CapitalRelease mirrors credit.CapitalRelease (Vol. III Ch. 32).
// cause: 1 = input (raw-material) price fall; 2 = merchants' idle money;
//   3 = rentier-class savings; 4 = revenue passing through money form.
// A release reflects real accumulation only when reinvested; otherwise it merely
// overstates it. Invariants: released_amount >= 0; loan_capital_overstatement.
//   overstatement_pct >= 0.
export interface CapitalRelease {
  id: string;
  name: string;
  released_amount: number; // pence; >= 0
  cause: number; // 1..4
  reinvested: boolean;
  loan_capital_overstatement: LoanCapitalOverstatement;
  description: string;
  created_at: string;
}

// CreateCapitalReleaseInput is the flat request body for
// POST /v1/credit/capital-release. The overstatement sub-object is flattened.
export interface CreateCapitalReleaseInput {
  name: string;
  released_amount: number;
  cause: number;
  reinvested: boolean;
  overstatement_pct: number;
  overstatement_description: string;
  description: string;
}

// ── Vol. III Ch. 29 — Component Parts of Bank Capital ────────────────────────

// BankCapitalComponentKind mirrors credit.BankCapitalComponentKind (Vol. III Ch. 29).
//   1=cash, 2=bill of exchange, 3=government bond, 4=stock, 5=mortgage.
export type BankCapitalComponentKind = 1 | 2 | 3 | 4 | 5;

// BankCapitalComponent mirrors credit.BankCapitalComponent (Vol. III Ch. 29).
export interface BankCapitalComponent {
  amount: number; // pence
  kind: BankCapitalComponentKind;
  description: string;
}

// BankCapital mirrors credit.BankCapital (Vol. III Ch. 29).
// Invariant: total_capital == cash_amount + securities_amount (server-computed).
export interface BankCapital {
  id: string;
  name: string;
  cash_amount: number;       // pence; gold + notes
  securities_amount: number; // pence; bills + public securities
  total_capital: number;     // DERIVED: cash + securities
  components: BankCapitalComponent[];
  description: string;
  created_at: string;
}

// CreateBankCapitalInput is the request body for POST /v1/credit/bank-capital.
export interface CreateBankCapitalInput {
  name: string;
  cash_amount: number;
  securities_amount: number;
  components: BankCapitalComponent[];
  description: string;
}

// FictitiousCapitalValuation mirrors credit.FictitiousCapitalValuation (Vol. III Ch. 29).
// Invariants:
//   interest_rate_bp > 0
//   market_value == roundHalfUp(annual_income*10000, interest_rate_bp)
//   market value moves inversely with the rate of interest
export interface FictitiousCapitalValuation {
  id: string;
  name: string;
  nominal_value: number;    // pence; face / par value
  annual_income: number;    // pence; yearly yield
  interest_rate_bp: number; // > 0; prevailing rate in bp
  dividend_rate_bp: number; // optional; for stocks, in bp
  market_value: number;     // DERIVED: roundHalfUp(income*10000, rate_bp)
  description: string;
  created_at: string;
}

// CreateFictitiousCapitalValuationInput is the request body for
// POST /v1/credit/fictitious-capital-valuation.
export interface CreateFictitiousCapitalValuationInput {
  name: string;
  nominal_value: number;
  annual_income: number;
  interest_rate_bp: number;
  dividend_rate_bp: number;
  description: string;
}

// DepositEntry mirrors credit.DepositEntry (Vol. III Ch. 29) — a book-entry
// deposit; the purest form of fictitious bank capital. Value type only.
export interface DepositEntry {
  id: string;
  amount: number; // pence; >= 0
  description: string;
}

// ClearingHouseSettlement mirrors credit.ClearingHouseSettlement (Vol. III Ch. 33).
// The London Clearing House nets mutual claims among banks; only the net_balance
// is settled in actual money. money_used <= total_claims always holds.
export interface ClearingHouseSettlement {
  id: string;
  name: string;
  description: string;
  total_claims: number; // pence; > 0
  money_used: number;   // pence; >= 0, <= total_claims
  net_balance: number;  // pence; = total_claims - money_used
  created_at: string;
}

// CreateClearingHouseSettlementInput is the request body for
// POST /v1/credit/clearing-house-settlements.
export interface CreateClearingHouseSettlementInput {
  name: string;
  description: string;
  total_claims: number;
  money_used: number;
}

// NoteIssueConstraint mirrors credit.NoteIssueConstraint (Vol. III Ch. 34).
// Peel's Bank Act 1844 split the Bank of England into Issue and Banking departments;
// MaxNotes = GoldBacking + UnbackedLimit. Regime encodes the legal period.
export interface NoteIssueConstraint {
  id: string;
  name: string;
  regime: number;          // 1=Pre-Act 1844, 2=Bank Act 1844, 3=Suspended
  max_notes: number;       // pence; derived = gold_backing + unbacked_limit
  gold_backing: number;    // pence; >= 0
  unbacked_limit: number;  // pence; >= 0
  notes_beyond_max: number; // pence; >= 0
  description: string;
  created_at: string;
}

// CreateNoteIssueConstraintInput is the request body for
// POST /v1/credit/note-issue-constraints.
// max_notes is server-derived (gold_backing + unbacked_limit).
export interface CreateNoteIssueConstraintInput {
  name: string;
  regime: number;
  gold_backing: number;
  unbacked_limit: number;
  notes_beyond_max: number;
  description: string;
}

// GoldReserve mirrors credit.GoldReserve (Vol. III Ch. 35).
// The BoE gold reserve is the pivot of the credit system; a falling reserve
// raises the rate of exchange and triggers gold inflows to restore it.
export interface GoldReserve {
  id: string;
  name: string;
  amount: number;      // pence; central-bank gold stock
  description: string;
  created_at: string;
}

export interface CreateGoldReserveInput {
  name: string;
  amount: number;
  description: string;
}

// RateOfExchange mirrors credit.RateOfExchange (Vol. III Ch. 35).
// Tracks deviation from mint parity in basis points (signed; positive = above par).
// Determined by the balance of PAYMENTS, not the balance of trade.
export interface RateOfExchange {
  id: string;
  name: string;
  deviation_bp: number;  // signed basis points from mint parity
  description: string;
  created_at: string;
}

export interface CreateRateOfExchangeInput {
  name: string;
  deviation_bp: number;
  description: string;
}

// UsurersCapital mirrors credit.UsurersCapital (Vol. III Ch. 36).
// Usury's capital is the oldest antediluvian form of interest-bearing capital —
// it predates wage-labour and operates on simple commodity circulation alone.
// Under developed capitalism it is disciplined and subordinated to the
// general rate of profit through the credit system.
export interface UsurersCapital {
  id: string;
  name: string;
  stage: number;              // UsuryDevelopmentStage: 1=Antiquity, 2=Mediaeval, 3=Transition, 4=Subordinated
  wage_labour: boolean;       // false in pre-capitalist stages; true only when subordinated
  borrowers: string[];        // debtor classes (peasants, petty producers, nobles, etc.)
  interest_rate_bp: number;   // interest rate in basis points (e.g. 4800 = 48%)
  subordinated_to: string;    // empty unless stage=4; e.g. "industrial_capital"
  description: string;
  created_at: string;
}

export interface CreateUsurersCapitalInput {
  name: string;
  stage: number;
  wage_labour: boolean;
  borrowers: string[];
  interest_rate_bp: number;
  subordinated_to: string;
  description: string;
}

// LandParcel mirrors rent.LandParcel (Vol. III Ch. 37 — opens Part VI).
// fertility_grade: 1=A (worst, regulating) … 4=D (best).
export interface LandParcel {
  id: string;
  owner_id: string;
  fertility_grade: number;
  area_hectares: number;
  location: string;
  created_at: string;
}
export interface CreateLandParcelInput {
  owner_id: string;
  fertility_grade: number;
  area_hectares: number;
  location: string;
}
// GroundRent mirrors rent.GroundRent. form: 1=Differential, 2=Absolute, 3=Monopoly.
export interface GroundRent {
  id: string;
  parcel_id: string;
  capitalist_id: string;
  land_owner_id: string;
  form: number;
  amount_labour_minutes: number;
  period_years: number;
  created_at: string;
}
export interface CreateGroundRentInput {
  parcel_id: string;
  capitalist_id: string;
  land_owner_id: string;
  form: number;
  amount_labour_minutes: number;
  period_years: number;
}

// Vol. III Ch. 38 — Differential Rent: General Remarks
export interface PriceOfProductionSurplusProfit {
  id: string;
  capital_id: string;
  general_price_of_production_bp: number;
  individual_price_of_production_bp: number;
  surplus_profit_bp: number;
  created_at: string;
}
export interface CreatePoPSurplusProfitInput {
  capital_id: string;
  general_price_of_production_bp: number;
  individual_price_of_production_bp: number;
}
export interface MonopolisedNaturalForce {
  id: string;
  kind: string;
  is_monopolisable: boolean;
  parcel_id: string;
  created_at: string;
}
export interface CreateMonopolisedNaturalForceInput {
  kind: string;
  is_monopolisable: boolean;
  parcel_id: string;
}
export interface DifferentialRentRecord {
  id: string;
  parcel_id: string;
  surplus_profit_bp: number;
  annual_rent_labour_minutes: number;
  created_at: string;
}
export interface CreateDifferentialRentInput {
  parcel_id: string;
  surplus_profit_bp: number;
  annual_rent_labour_minutes: number;
}
export interface CapitalisedRentPrice {
  id: string;
  parcel_id: string;
  annual_rent_labour_minutes: number;
  interest_rate_bp: number;
  capitalised_price_labour_minutes: number;
  created_at: string;
}
export interface CreateCapitalisedRentPriceInput {
  parcel_id: string;
  annual_rent_labour_minutes: number;
  interest_rate_bp: number;
}

// Vol. III Ch. 39 — Differential Rent I
export type SoilGrade = 1 | 2 | 3 | 4;
export interface DifferentialRentIEntry {
  id: string;
  table_id: string;
  grade: SoilGrade;
  capital_advanced_bp: number;
  output_quarters: number;
  price_of_production_bp: number;
  surplus_product_qtrs: number;
  money_rent_bp: number;
  created_at: string;
}
export interface DifferentialRentITable {
  id: string;
  description: string;
  total_rental_bp: number;
  regulating_grade: SoilGrade;
  created_at: string;
}
export interface DifferentialRentITableWithEntries {
  table: DifferentialRentITable;
  entries: DifferentialRentIEntry[];
}
export interface LocationRentFactor {
  id: string;
  parcel_id: string;
  transport_cost_bp: number;
  rent_equivalent_bp: number;
  created_at: string;
}
export interface CreateDifferentialRentIEntryInput {
  grade: SoilGrade;
  capital_advanced_bp: number;
  output_quarters: number;
  price_of_production_bp: number;
}
export interface CreateDifferentialRentITableInput {
  description: string;
  entries: CreateDifferentialRentIEntryInput[];
}
export interface CreateLocationRentFactorInput {
  parcel_id: string;
  transport_cost_bp: number;
  rent_equivalent_bp: number;
}

// Vol. III Ch. 40 — Differential Rent II
export interface SuccessiveInvestment {
  id: string;
  parcel_id: string;
  tranche_number: number;
  capital_bp: number;
  output_quarters: number;
  surplus_profit_bp: number;
  created_at: string;
}
export interface DifferentialRentIITable {
  id: string;
  parcel_id: string;
  total_capital_bp: number;
  total_output_quarters: number;
  total_rent_bp: number;
  rent_per_acre_bp: number;
}
export interface DifferentialRentIITableWithInvestments {
  table: DifferentialRentIITable;
  investments: SuccessiveInvestment[];
}
export interface LeaseTerm {
  id: string;
  parcel_id: string;
  start_year: number;
  duration_years: number;
  fixed_rent_bp: number;
  created_at: string;
}
export interface CreateSuccessiveInvestmentInput {
  parcel_id: string;
  tranche_number: number;
  capital_bp: number;
  output_quarters: number;
  surplus_profit_bp: number;
}
export interface CreateLeaseTermInput {
  parcel_id: string;
  start_year: number;
  duration_years: number;
  fixed_rent_bp: number;
}

// Vol. III Ch. 41 — DR II Constant Price
export type CapitalProductivityCase = 1 | 2 | 3 | 4;
export interface DifferentialRentIIOutcome {
  id: string;
  parcel_id: string;
  case: CapitalProductivityCase;
  additional_capital_bp: number;
  new_rent_per_acre_bp: number;
  total_rental_bp: number;
  rate_of_surplus_profit_bp: number;
  created_at: string;
}
export interface CreateDifferentialRentIIOutcomeInput {
  case: CapitalProductivityCase;
  investments: Array<{
    parcel_id: string;
    tranche_number: number;
    capital_bp: number;
    output_quarters: number;
    surplus_profit_bp: number;
  }>;
}
export interface IntensiveExtensiveComparison {
  id: string;
  intensive_rent_per_acre_bp: number;
  extensive_rent_per_acre_bp: number;
  total_capital_bp: number;
  total_rental_same_bp: number;
  created_at: string;
}
export interface CreateIntensiveExtensiveComparisonInput {
  intensive_rent_per_acre_bp: number;
  extensive_rent_per_acre_bp: number;
  total_capital_bp: number;
  total_rental_same_bp: number;
}

// Vol. III Ch. 42 — DR II Falling Price of Production
export type FallingPriceVariant = 1 | 2 | 3;
export interface SoilExclusionEvent {
  id: string;
  excluded_grade: number;
  new_regulating_grade: number;
  new_price_of_production_bp: number;
  created_at: string;
}
export interface FallingPriceDR2Outcome {
  id: string;
  variant: FallingPriceVariant;
  exclusion_event_id: string;
  total_money_rental_bp: number;
  total_grain_rent_qtrs: number;
  total_output_quarters: number;
  total_capital_bp: number;
  created_at: string;
}
export interface NormalCapitalPerAcre {
  id: string;
  period_label: string;
  capital_per_acre_bp: number;
  created_at: string;
}
export interface CreateSoilExclusionEventInput {
  excluded_grade: number;
  new_regulating_grade: number;
  new_price_of_production_bp: number;
}
export interface CreateFallingPriceDR2OutcomeInput {
  variant: FallingPriceVariant;
  exclusion_event_id: string;
  investments: Array<{
    parcel_id: string;
    tranche_number: number;
    capital_bp: number;
    output_quarters: number;
    surplus_profit_bp: number;
  }>;
}
export interface CreateNormalCapitalPerAcreInput {
  period_label: string;
  capital_per_acre_bp: number;
}

// Vol. III Ch. 43 — DR II Rising Price of Production
export type RisingPriceVariant = 1 | 2 | 3 | 4 | 5;
export interface RisingPriceDR2Outcome {
  id: string;
  variant: RisingPriceVariant;
  total_capital_shillings: number;
  total_rent_shillings: number;
  total_output_bushels: number;
  price_per_bushel_shillings: number;
  created_at: string;
}
export interface CreateRisingPriceDR2OutcomeInput {
  variant: RisingPriceVariant;
  total_capital_shillings: number;
  total_rent_shillings: number;
  total_output_bushels: number;
  price_per_bushel_shillings: number;
}
export interface RentSequence {
  id: string;
  base_rent: number;
  increment: number;
  num_grades: number;
  created_at: string;
}
export interface CreateRentSequenceInput {
  base_rent: number;
  increment: number;
  num_grades: number;
}
export interface CreateRentSequenceResponse extends RentSequence {
  sequence: number[];
}
export interface EngelsTableEntry {
  id: string;
  table_number: number;
  soil_grade_label: string;
  output_bushels: number;
  rent_shillings: number;
  capital_shillings: number;
  created_at: string;
}
export interface CreateEngelsTableEntryInput {
  table_number: number;
  soil_grade_label: string;
  output_bushels: number;
  rent_shillings: number;
  capital_shillings: number;
}

// Vol. III Ch. 44 — Differential Rent on the Worst Cultivated Soil
export type RentEmergenceMechanism = 1 | 2 | 3;
export interface RentOnWorstSoil {
  id: string;
  parcel_id: string;
  mechanism: RentEmergenceMechanism;
  old_price_of_production_bp: number;
  new_price_of_production_bp: number;
  rent_per_acre_bp: number;
  created_at: string;
}
export interface CreateRentOnWorstSoilInput {
  parcel_id: string;
  mechanism: RentEmergenceMechanism;
  old_price_of_production_bp: number;
  new_price_of_production_bp: number;
}
export interface SoilImprovementCapital {
  id: string;
  parcel_id: string;
  capital_invested_bp: number;
  productivity_gain_bp: number;
  amortisation_years: number;
  amortised_years: number;
  created_at: string;
}
export interface CreateSoilImprovementCapitalInput {
  parcel_id: string;
  capital_invested_bp: number;
  productivity_gain_bp: number;
  amortisation_years: number;
  amortised_years: number;
}
export interface RegulatingPriceShift {
  id: string;
  from_grade: number;
  to_grade: number;
  old_price_bp: number;
  new_price_bp: number;
  created_at: string;
}
export interface CreateRegulatingPriceShiftInput {
  from_grade: number;
  to_grade: number;
  old_price_bp: number;
  new_price_bp: number;
}

// Vol. III Ch. 45 — Absolute Ground-Rent
export interface AgriculturalComposition {
  id: string;
  constant_capital_bp: number;
  variable_capital_bp: number;
  composition_ratio_bp: number;
  social_average_ratio_bp: number;
  is_below_average: boolean;
  created_at: string;
}
export interface CreateAgriculturalCompositionInput {
  constant_capital_bp: number;
  variable_capital_bp: number;
  social_average_ratio_bp: number;
}
export interface ValuePriceGap {
  id: string;
  product_id: string;
  value_labour_minutes: number;
  price_of_production_labour_minutes: number;
  gap_labour_minutes: number;
  created_at: string;
}
export interface CreateValuePriceGapInput {
  product_id: string;
  value_labour_minutes: number;
  price_of_production_labour_minutes: number;
}
export interface AbsoluteRent {
  id: string;
  parcel_id: string;
  value_price_gap_id: string;
  surplus_value_bp: number;
  average_profit_bp: number;
  absolute_rent_bp: number;
  created_at: string;
}
export interface CreateAbsoluteRentInput {
  parcel_id: string;
  value_price_gap_id: string;
  surplus_value_bp: number;
  average_profit_bp: number;
}
export interface AbsoluteRentLimit {
  id: string;
  max_rent_bp: number;
  actual_rent_bp: number;
  is_at_limit: boolean;
  created_at: string;
}
export interface CreateAbsoluteRentLimitInput {
  max_rent_bp: number;
  actual_rent_bp: number;
}

// Vol. III Ch. 46 — Building Site Rent, Rent in Mining, Price of Land
export type LandPriceScenarioKind = 1 | 2 | 3 | 4 | 5;
export interface BuildingSiteRent { id: string; parcel_id: string; location_grade: number; annual_rent_labour_minutes: number; lease_term_years: number; created_at: string; }
export interface CreateBuildingSiteRentInput { parcel_id: string; location_grade: number; annual_rent_labour_minutes: number; lease_term_years: number; }
export interface MiningRent { id: string; parcel_id: string; ore_grade: number; annual_rent_labour_minutes: number; surplus_profit_bp: number; created_at: string; }
export interface CreateMiningRentInput { parcel_id: string; ore_grade: number; annual_rent_labour_minutes: number; surplus_profit_bp: number; }
export interface MonopolyPriceRent { id: string; parcel_id: string; product_kind: string; value_labour_minutes: number; monopoly_sale_price_labour_minutes: number; monopoly_rent_labour_minutes: number; created_at: string; }
export interface CreateMonopolyPriceRentInput { parcel_id: string; product_kind: string; value_labour_minutes: number; monopoly_sale_price_labour_minutes: number; }
export interface LandPriceScenario { id: string; kind: LandPriceScenarioKind; parcel_id: string; old_annual_rent_bp: number; new_annual_rent_bp: number; old_interest_rate_bp: number; new_interest_rate_bp: number; old_land_price_bp: number; new_land_price_bp: number; created_at: string; }
export interface CreateLandPriceScenarioInput { kind: LandPriceScenarioKind; parcel_id: string; old_annual_rent_bp: number; new_annual_rent_bp: number; old_interest_rate_bp: number; new_interest_rate_bp: number; old_land_price_bp: number; }

// Vol. III Ch. 47 — Genesis of Capitalist Ground-Rent
export type RentFormStage = 1 | 2 | 3 | 4;
export interface RentHistoricalForm { id: string; stage: RentFormStage; name: string; coercion_kind: string; surplus_form: string; example_region: string; example_period: string; created_at: string; }
export interface CreateRentHistoricalFormInput { stage: RentFormStage; name: string; coercion_kind: string; surplus_form: string; example_region: string; example_period: string; }
export interface HistoricalRentTransition { id: string; from_stage: RentFormStage; to_stage: RentFormStage; preconditions: string; driving_force: string; created_at: string; }
export interface CreateHistoricalRentTransitionInput { from_stage: RentFormStage; to_stage: RentFormStage; preconditions: string; driving_force: string; }
export interface SmallPeasantProduction { id: string; parcel_id: string; total_income_bp: number; rent_component_bp: number; profit_component_bp: number; wages_component_bp: number; is_conflated: boolean; created_at: string; }
export interface CreateSmallPeasantProductionInput { parcel_id: string; rent_component_bp: number; profit_component_bp: number; wages_component_bp: number; is_conflated: boolean; }

// Vol. III Ch. 48 — The Trinity Formula
export type RevenueSourceKind = 1 | 2 | 3;
export interface RevenueStream { id: string; source: RevenueSourceKind; apparent_revenue_bp: number; actual_source_bp: number; is_fetishised: boolean; created_at: string; }
export interface TrinityFormula { id: string; capital_stream_id: string; land_stream_id: string; labour_stream_id: string; total_surplus_value_bp: number; total_apparent_revenue_bp: number; created_at: string; }
export interface RevenueFetishForm { id: string; source: RevenueSourceKind; surface_formula: string; real_relation: string; mystification_kind: string; created_at: string; }
export interface TrinityStreamInput { apparent_revenue_bp: number; actual_source_bp: number; is_fetishised: boolean; }
export interface CreateTrinityFormulaInput { capital_stream: TrinityStreamInput; land_stream: TrinityStreamInput; labour_stream: TrinityStreamInput; }
export interface TrinityFormulaResponse { formula: TrinityFormula; streams: RevenueStream[]; }

// Vol. III Ch. 49 — Concerning the Analysis of the Process of Production
export interface AnnualValueComposition { id: string; constant_capital_labour_minutes: number; variable_capital_labour_minutes: number; surplus_value_labour_minutes: number; total_value_labour_minutes: number; new_value_labour_minutes: number; created_at: string; }
export interface RevenueLimit { id: string; composition_id: string; max_distributable_revenue_lm: number; constant_capital_replacement: number; created_at: string; }
export interface AnnualCompositionResponse { composition: AnnualValueComposition; revenue_limit: RevenueLimit; }
export interface CreateAnnualCompositionInput { constant_capital_labour_minutes: number; variable_capital_labour_minutes: number; surplus_value_labour_minutes: number; }
export interface SurplusValuePartition { id: string; composition_id: string; total_surplus_value_lm: number; profit_enterprise_lm: number; interest_lm: number; ground_rent_lm: number; residual_surplus_lm: number; created_at: string; }
export interface CreateSurplusPartitionInput { composition_id: string; profit_enterprise_lm: number; interest_lm: number; ground_rent_lm: number; residual_surplus_lm: number; }

// Vol. III Ch. 50 — Illusions Created by Competition
export interface CompetitionIllusion { id: string; name: string; surface_appearance: string; real_relation: string; produced_by: string; created_at: string; }
export interface CostPriceFetish { id: string; constant_capital_bp: number; variable_capital_bp: number; cost_price_bp: number; profit_bp: number; actual_value_bp: number; created_at: string; }
export interface CreateCostPriceFetishInput { constant_capital_bp: number; variable_capital_bp: number; profit_bp: number; }
export interface AnnualRevenueAccount { id: string; new_value_created_lm: number; wages_share_lm: number; profit_and_rent_share_lm: number; constant_capital_lm: number; total_social_product_lm: number; created_at: string; }
export interface CreateAnnualRevenueAccountInput { constant_capital_lm: number; wages_share_lm: number; profit_and_rent_share_lm: number; }
export interface ValueAppearanceContrast { id: string; illusion_id: string; surface: string; reality: string; created_at: string; }

// Vol. III Ch. 51 — Distribution Relations and Production Relations
export interface ProductionRelationBasis { id: string; name: string; characteristic_feature: string; is_historically_specific: boolean; created_at: string; }
export interface DistributionRelation { id: string; form: string; production_basis_id: string; apparent_character: string; real_character: string; corresponds_to_surplus_form: string; created_at: string; }
export interface HistoricalTransience { id: string; distribution_form_id: string; mode_of_production: string; precondition: string; succeeding_form: string; created_at: string; }
export interface ProductiveForceContradiction { id: string; productive_force: string; social_form: string; contradiction_nature: string; crisis_signal: string; created_at: string; }
export interface CreateProductiveForceContradictionInput { productive_force: string; social_form: string; contradiction_nature: string; crisis_signal: string; }

// Vol. III Ch. 52 — Classes
export interface SocialClass { id: string; name: string; relation_to_means_of_production: string; primary_revenue_form: string; historical_precondition: string; created_at: string; }
export interface ClassIncomeSource { id: string; class_id: string; revenue_form: string; surplus_value_component_bp: number; created_at: string; }
export interface ClassTendency { id: string; class_id: string; description: string; direction: string; created_at: string; }
export interface ClassAmbiguity { id: string; intermediate_group: string; income_source: string; why_not_a_class: string; created_at: string; }
export interface SocialClassDetail { class: SocialClass; income_sources: ClassIncomeSource[]; tendencies: ClassTendency[]; }
