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
}

export interface LabourScenarioResult {
  daily_value_minutes: number;
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  labour_power_value_minutes: number;
  rate_of_surplus_value: number;
  law_constant_daily_value: boolean;
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
