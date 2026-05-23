import type {
  ExtendedReproductionInput,
  ExtendedReproductionResult,
  RepaymentPeriodInput,
  RepaymentPeriodResult,
  SimpleReproductionInput,
  SimpleReproductionResult,
  SplitSurplusInput,
  SplitSurplusResult,
  Agent,
  CapitalCircuit,
  CapitalCompositionResult,
  Circuit,
  CircuitProof,
  Commodity,
  ComputeCircuitInput,
  ComputePriceInput,
  CreateAgentCircuitInput,
  CreateAgentInput,
  CreateCircuitInput,
  CreateCommodityInput,
  CreateExchangeInput,
  CreateHoardInput,
  CreateLabourCapitalistInput,
  CreateLabourPowerOfferingInput,
  CreateLabourPowerPurchaseInput,
  CreateLabourWorkerInput,
  CreateOfferInput,
  CreateOwnerInput,
  CreatePaymentObligationInput,
  CreateProductionAccountInput,
  CreateWorldMoneyTransferInput,
  DecomposeCapitalInput,
  DecomposeCapitalResult,
  Exchange,
  ExchangeRatio,
  ExchangeSimulation,
  Hoard,
  LabourCapitalist,
  LabourPowerOffering,
  LabourPowerPurchase,
  LabourWorker,
  MoneyCommodity,
  Offer,
  Owner,
  PaymentObligation,
  Price,
  ProductionAccountResult,
  RateOfSurplusValueInput,
  RateOfSurplusValueResult,
  SimulateExchangeInput,
  SocialRelations,
  UniversalEquivalent,
  UpdateAgentInput,
  UpdateCommodityInput,
  ValueResponse,
  WorldMoneyTransfer,
  RunLabourProcessInput,
  RunLabourProcessResult,
  CreateWorkingDayInput,
  ValidateWorkingDayInput,
  ValidateWorkingDayResponse,
  WorkingDayResponse,
  CreateRelayScheduleInput,
  RelaySchedule,
  ComputeMassInput,
  SurplusLimitsResponse,
  SurplusValueSnapshot,
  ProductionWorkingDay,
  ShortenWorkingDayResponse,
  ProductionRateResult,
  ExtraSurplusValueResult,
  RecordWorkingDayInput,
  ShortenWorkingDayInput,
  ExtraSurplusValueInput,
  Cooperation,
  CreateCooperationInput,
  CollectiveWorkingDayResult,
  AverageSocialLabourResult,
  MinimumCapitalInput,
  MinimumCapitalResult,
  Manufacture,
  CreateManufactureInput,
  ProportionalGroupSizeResult,
  ManufactureMinimumCapitalResult,
  ManufactureForm,
  Machine,
  CreateMachineInput,
  CreateFactoryInput,
  Factory,
  MachineWearResponse,
  FactoryTickResult,
  RelativeSurplusInput,
  RelativeSurplusFromSourceInput,
  RelativeSurplusResult,
  FactoryTick,
  AbsoluteSurplusValueInput,
  AbsoluteSurplusValueResult,
  LabourScenarioInput,
  LabourScenarioResult,
  RelativeSurplusValueInput,
  RelativeSurplusValueResult,
  SurplusValueRateResult,
  RatesOfSurplusValueInput,
  RatesOfSurplusValueResult,
  CreateWageFormInput,
  WageForm,
  ComputeHourlyPriceInput,
  HourlyPriceOfLabour,
  CreateWorkingSessionInput,
  WorkingSession,
  ComputePiecePriceInput,
  ComputePiecePriceResult,
  CreatePieceWageInput,
  PieceWage,
  CreateSubContractInput,
  SubContract,
  NationalIntensity,
  DayWage,
  StandardisedWage,
  WageComparison,
  RegisterIntensityInput,
  RegisterDayWageInput,
  OrganicCompositionInput,
  OrganicCompositionResult,
  LabourDemandInput,
  LabourDemandResult,
  ReserveArmyInput,
  ReserveArmyResult,
  GeneralLawScenarioInput,
  GeneralLawScenarioResult,
  HistoricalStage,
  HistoricalStageInput,
  SeedScenarioInput,
  SeedScenarioResult,
  EnclosureEvent,
  EnclosureEventInput,
  WageStatute,
  WageStatuteInput,
  VagrancyLaw,
  VagrancyLawInput,
  StatutoryWage,
  LabourDisciplineRegime,
  FarmTenure,
  FarmTenureInput,
  RealRentResult,
  DomesticIndustryInput,
  DomesticIndustry,
  MarketFormationInput,
  MarketFormation,
  HomeMarketSize,
  CapitalOriginInput,
  CapitalOrigin,
  ColonialTransferInput,
  ColonialTransfer,
  NationalDebtInput,
  NationalDebt,
  IndustrialCapitalGenesis,
  AccumulationTrajectory,
  ColonialLabourMarket,
  ComputeIndependenceInput,
  CreateColonialLabourMarketInput,
  NegationOfNegationResponse,
  RegulateColonialMarketInput,
  RunCentralisationInput,
  WageWorkerIndependence,
  ProductiveCircuit,
  CreateProductiveCircuitInput,
  LatentMoneyCapital,
  RevenueCircuit,
  CapitalisationStep,
  ReserveFund,
  ReserveDraw,
} from "./types";

const BASE = "/api";

async function http<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (res.status === 204) {
    return undefined as T;
  }
  const text = await res.text();
  if (!res.ok) {
    let msg = res.statusText;
    try {
      const body = JSON.parse(text) as { error?: string };
      if (body.error) msg = body.error;
    } catch {
      // ignore
    }
    throw new Error(`${res.status}: ${msg}`);
  }
  return text ? (JSON.parse(text) as T) : (undefined as T);
}

export const api = {
  listCommodities: () =>
    http<{ items: Commodity[] }>("/v1/commodities").then((r) => r.items),

  getCommodity: (id: string) => http<Commodity>(`/v1/commodities/${id}`),

  createCommodity: (input: CreateCommodityInput) =>
    http<Commodity>("/v1/commodities", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  updateCommodity: (id: string, input: UpdateCommodityInput) =>
    http<Commodity>(`/v1/commodities/${id}`, {
      method: "PATCH",
      body: JSON.stringify(input),
    }),

  deleteCommodity: (id: string) =>
    http<void>(`/v1/commodities/${id}`, { method: "DELETE" }),

  computeValue: (id: string, quantity: number) =>
    http<ValueResponse>(`/v1/commodities/${id}/value`, {
      method: "POST",
      body: JSON.stringify({ quantity }),
    }),

  exchangeRatio: (baseId: string, quoteId: string, baseQty: number) =>
    http<ExchangeRatio>("/v1/exchange-ratio", {
      method: "POST",
      body: JSON.stringify({
        base_id: baseId,
        quote_id: quoteId,
        base_qty: baseQty,
      }),
    }),

  socialRelations: (id: string) =>
    http<SocialRelations>(`/v1/commodities/${id}/social-relations`),

  // --- market-service (Ch. 2) ---

  listOwners: () =>
    http<{ items: Owner[] }>("/v1/owners").then((r) => r.items),

  createOwner: (input: CreateOwnerInput) =>
    http<Owner>("/v1/owners", { method: "POST", body: JSON.stringify(input) }),

  listOffers: () =>
    http<{ items: Offer[] }>("/v1/offers").then((r) => r.items),

  createOffer: (input: CreateOfferInput) =>
    http<Offer>("/v1/offers", { method: "POST", body: JSON.stringify(input) }),

  deleteOffer: (id: string) =>
    http<void>(`/v1/offers/${id}`, { method: "DELETE" }),

  listExchanges: () =>
    http<{ items: Exchange[] }>("/v1/exchanges").then((r) => r.items),

  createExchange: (input: CreateExchangeInput) =>
    http<Exchange>("/v1/exchanges", { method: "POST", body: JSON.stringify(input) }),

  getUniversalEquivalent: () =>
    http<UniversalEquivalent>("/v1/universal-equivalent"),

  setUniversalEquivalent: (commodityId: string) =>
    http<UniversalEquivalent>("/v1/universal-equivalent", {
      method: "POST",
      body: JSON.stringify({ commodity_id: commodityId }),
    }),

  getMoneyCommodity: () =>
    http<MoneyCommodity>("/v1/money-commodity"),

  setMoneyCommodity: (commodityId: string) =>
    http<MoneyCommodity>("/v1/money-commodity", {
      method: "POST",
      body: JSON.stringify({ commodity_id: commodityId }),
    }),

  listPrices: () =>
    http<{ items: Price[] }>("/v1/prices").then((r) => r.items),

  computePrice: (input: ComputePriceInput) =>
    http<Price>("/v1/prices", { method: "POST", body: JSON.stringify(input) }),

  // --- market-service (Ch. 3: Money) ---

  getMoneyRequired: (sumOfPrices: number, velocity: number) =>
    http<{ money_required: number }>(
      `/v1/circulation/money-required?sum_of_prices=${sumOfPrices}&velocity=${velocity}`
    ),

  listCircuits: () =>
    http<{ items: Circuit[] }>("/v1/circuits").then((r) => r.items),

  createCircuit: (input: CreateCircuitInput) =>
    http<Circuit>("/v1/circuits", { method: "POST", body: JSON.stringify(input) }),

  listHoards: () =>
    http<{ items: Hoard[] }>("/v1/hoards").then((r) => r.items),

  createHoard: (input: CreateHoardInput) =>
    http<Hoard>("/v1/hoards", { method: "POST", body: JSON.stringify(input) }),

  listPaymentObligations: () =>
    http<{ items: PaymentObligation[] }>("/v1/payment-obligations").then((r) => r.items),

  createPaymentObligation: (input: CreatePaymentObligationInput) =>
    http<PaymentObligation>("/v1/payment-obligations", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  settlePaymentObligation: (id: string) =>
    http<PaymentObligation>(`/v1/payment-obligations/${id}/settle`, { method: "POST" }),

  listWorldMoneyTransfers: () =>
    http<{ items: WorldMoneyTransfer[] }>("/v1/world-money-transfers").then((r) => r.items),

  createWorldMoneyTransfer: (input: CreateWorldMoneyTransferInput) =>
    http<WorldMoneyTransfer>("/v1/world-money-transfers", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- agent-service (Ch. 4) ---

  listAgents: (classFilter?: string) =>
    http<{ items: Agent[] }>(
      classFilter ? `/v1/agents?class=${encodeURIComponent(classFilter)}` : "/v1/agents"
    ).then((r) => r.items),

  createAgent: (input: CreateAgentInput) =>
    http<Agent>("/v1/agents", { method: "POST", body: JSON.stringify(input) }),

  getAgent: (id: string) => http<Agent>(`/v1/agents/${id}`),

  updateAgent: (id: string, input: UpdateAgentInput) =>
    http<Agent>(`/v1/agents/${id}`, { method: "PATCH", body: JSON.stringify(input) }),

  deleteAgent: (id: string) =>
    http<void>(`/v1/agents/${id}`, { method: "DELETE" }),

  createAgentCircuit: (agentId: string, input: CreateAgentCircuitInput) =>
    http<CapitalCircuit>(`/v1/agents/${agentId}/circuits`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listAgentCircuits: (agentId: string) =>
    http<{ items: CapitalCircuit[] }>(`/v1/agents/${agentId}/circuits`).then((r) => r.items),

  reinvestAgent: (agentId: string, commodityId: string, mReturned: number) =>
    http<CapitalCircuit>(`/v1/agents/${agentId}/reinvest`, {
      method: "POST",
      body: JSON.stringify({ commodity_id: commodityId, m_returned: mReturned }),
    }),

  hoardAgent: (agentId: string) =>
    http<Agent>(`/v1/agents/${agentId}/hoard`, { method: "POST", body: JSON.stringify({}) }),

  // --- agent-service (Ch. 5) ---

  computeCircuit: (input: ComputeCircuitInput) =>
    http<CircuitProof>("/v1/circuit-probes", { method: "POST", body: JSON.stringify(input) }),

  simulateExchange: (input: SimulateExchangeInput) =>
    http<ExchangeSimulation>("/v1/exchange-simulations", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- agent-service (Ch. 6: Labour-Power) ---

  listLabourWorkers: () =>
    http<{ items: LabourWorker[] }>("/v1/workers").then((r) => r.items),

  createLabourWorker: (input: CreateLabourWorkerInput) =>
    http<LabourWorker>("/v1/workers", { method: "POST", body: JSON.stringify(input) }),

  getLabourWorker: (id: string) => http<LabourWorker>(`/v1/workers/${id}`),

  listLabourCapitalists: () =>
    http<{ items: LabourCapitalist[] }>("/v1/capitalists").then((r) => r.items),

  createLabourCapitalist: (input: CreateLabourCapitalistInput) =>
    http<LabourCapitalist>("/v1/capitalists", { method: "POST", body: JSON.stringify(input) }),

  getLabourCapitalist: (id: string) => http<LabourCapitalist>(`/v1/capitalists/${id}`),

  listOfferings: () =>
    http<{ items: LabourPowerOffering[] }>("/v1/labour-power/offerings").then((r) => r.items),

  createOffering: (input: CreateLabourPowerOfferingInput) =>
    http<LabourPowerOffering>("/v1/labour-power/offerings", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listLabourPurchases: () =>
    http<{ items: LabourPowerPurchase[] }>("/v1/labour-power/purchases").then((r) => r.items),

  createLabourPurchase: (input: CreateLabourPowerPurchaseInput) =>
    http<LabourPowerPurchase>("/v1/labour-power/purchases", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getLabourPurchase: (id: string) =>
    http<LabourPowerPurchase>(`/v1/labour-power/purchases/${id}`),

  // --- agent-service (Ch. 7: Labour-Process) ---

  runLabourProcess: (input: RunLabourProcessInput) =>
    http<RunLabourProcessResult>("/v1/labour-processes", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getLabourProcess: (id: string) =>
    http<RunLabourProcessResult>(`/v1/labour-processes/${id}`),

  // --- commodity-service (Ch. 8: Constant & Variable Capital) ---

  decomposeCapital: (input: DecomposeCapitalInput) =>
    http<DecomposeCapitalResult>("/v1/capital/decompose", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getCapitalComposition: (constantValue: number, variableValue: number) =>
    http<CapitalCompositionResult>(
      `/v1/capital/composition?constant_value=${constantValue}&variable_value=${variableValue}`
    ),

  // --- commodity-service (Ch. 9: Rate of Surplus-Value) ---

  createProductionAccount: (input: CreateProductionAccountInput) =>
    http<ProductionAccountResult>("/v1/production-accounts", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listProductionAccounts: () =>
    http<{ items: ProductionAccountResult[] }>("/v1/production-accounts").then(
      (r) => r.items
    ),

  getProductionAccount: (id: string) =>
    http<ProductionAccountResult>(`/v1/production-accounts/${id}`),

  computeRateOfSurplusValue: (input: RateOfSurplusValueInput) =>
    http<RateOfSurplusValueResult>("/v1/rate-of-surplus-value", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- agent-service (Ch. 10: The Working-Day) ---

  createWorkingDay: (input: CreateWorkingDayInput) =>
    http<WorkingDayResponse>("/v1/working-days", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getWorkingDay: (id: string) =>
    http<WorkingDayResponse>(`/v1/working-days/${id}`),

  validateWorkingDay: (input: ValidateWorkingDayInput) =>
    http<ValidateWorkingDayResponse>("/v1/working-days/validate", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createRelaySchedule: (input: CreateRelayScheduleInput) =>
    http<RelaySchedule>("/v1/relay-schedules", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getRelaySchedule: (id: string) =>
    http<RelaySchedule>(`/v1/relay-schedules/${id}`),

  // --- simulation-engine (Ch. 11: Rate and Mass of Surplus-Value) ---

  computeSurplusMass: (input: ComputeMassInput) =>
    http<SurplusValueSnapshot>("/v1/surplus/mass", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getSurplusLimits: (labourPowerValue?: number, workerCount?: number) => {
    const params = new URLSearchParams();
    if (labourPowerValue !== undefined) params.set("labour_power_value", String(labourPowerValue));
    if (workerCount !== undefined) params.set("worker_count", String(workerCount));
    const qs = params.toString();
    return http<SurplusLimitsResponse>(`/v1/surplus/limits${qs ? `?${qs}` : ""}`);
  },

  // --- simulation-engine (Ch. 12: Relative Surplus-Value) ---

  recordWorkingDay: (input: RecordWorkingDayInput) =>
    http<ProductionWorkingDay>("/v1/production/working-day", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  shortenWorkingDay: (input: ShortenWorkingDayInput) =>
    http<ShortenWorkingDayResponse>("/v1/production/working-day/shorten", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getProductionRate: (necessary: number, surplus: number) =>
    http<ProductionRateResult>(
      `/v1/production/rate-of-surplus-value?necessary=${necessary}&surplus=${surplus}`
    ),

  computeExtraSurplusValue: (input: ExtraSurplusValueInput) =>
    http<ExtraSurplusValueResult>("/v1/production/extra-surplus-value", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- agent-service (Ch. 13: Co-operation) ---

  listCooperations: (capitalistId?: string) =>
    http<{ items: Cooperation[] }>(
      capitalistId
        ? `/v1/cooperations?capitalist_id=${encodeURIComponent(capitalistId)}`
        : "/v1/cooperations"
    ).then((r) => r.items),

  getCooperation: (id: string) => http<Cooperation>(`/v1/cooperations/${id}`),

  createCooperation: (input: CreateCooperationInput) =>
    http<Cooperation>("/v1/cooperations", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeCollectiveWorkingDay: (id: string) =>
    http<CollectiveWorkingDayResult>(`/v1/cooperations/${id}/collective-working-day`, {
      method: "POST",
      body: JSON.stringify({}),
    }),

  computeAverageSocialLabour: (id: string) =>
    http<AverageSocialLabourResult>(`/v1/cooperations/${id}/average-social-labour`, {
      method: "POST",
      body: JSON.stringify({}),
    }),

  computeMinimumCapital: (input: MinimumCapitalInput) =>
    http<MinimumCapitalResult>("/v1/cooperations/minimum-capital", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- agent-service (Ch. 14: Division of Labour and Manufacture) ---

  listManufactures: (filters?: { capitalistId?: string; form?: ManufactureForm }) => {
    const q = new URLSearchParams();
    if (filters?.capitalistId) q.set("capitalist_id", filters.capitalistId);
    if (filters?.form) q.set("form", filters.form);
    const qs = q.toString();
    return http<{ items: Manufacture[] }>(
      qs ? `/v1/manufactures?${qs}` : "/v1/manufactures"
    ).then((r) => r.items);
  },

  getManufacture: (id: string) => http<Manufacture>(`/v1/manufactures/${id}`),

  createManufacture: (input: CreateManufactureInput) =>
    http<Manufacture>("/v1/manufactures", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  proportionalGroupSize: (id: string, targetOutputRate: number) =>
    http<ProportionalGroupSizeResult>(
      `/v1/manufactures/${id}/proportional-group-size`,
      {
        method: "POST",
        body: JSON.stringify({ target_output_rate: targetOutputRate }),
      }
    ),

  scaleManufacture: (id: string, multiplier: number) =>
    http<Manufacture>(`/v1/manufactures/${id}/scale`, {
      method: "POST",
      body: JSON.stringify({ multiplier }),
    }),

  getManufactureMinimumCapital: (id: string, rawMaterialCostFactor: number) =>
    http<ManufactureMinimumCapitalResult>(
      `/v1/manufactures/${id}/minimum-capital?raw_material_cost_factor=${encodeURIComponent(
        String(rawMaterialCostFactor)
      )}`
    ),

  // --- simulation-engine (Ch. 15: Machinery and Modern Industry) ---

  listMachines: () =>
    http<{ items: Machine[] }>("/v1/machines").then((r) => r.items),

  getMachine: (id: string) => http<Machine>(`/v1/machines/${id}`),

  createMachine: (input: CreateMachineInput) =>
    http<Machine>("/v1/machines", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getMachineWear: (id: string) =>
    http<MachineWearResponse>(`/v1/machines/${id}/wear`),

  listFactories: () =>
    http<{ items: Factory[] }>("/v1/factories").then((r) => r.items),

  getFactory: (id: string) => http<Factory>(`/v1/factories/${id}`),

  createFactory: (input: CreateFactoryInput) =>
    http<Factory>("/v1/factories", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  tickFactory: (id: string) =>
    http<FactoryTickResult>(`/v1/factories/${id}/tick`, {
      method: "POST",
      body: JSON.stringify({}),
    }),

  listFactoryTicks: (id: string, limit?: number) => {
    const qs = limit ? `?limit=${limit}` : "";
    return http<{ items: FactoryTick[] }>(`/v1/factories/${id}/ticks${qs}`).then((r) => r.items);
  },

  // --- simulation-engine (Part IV bridge: Ch.13/14/15 → Ch.12) ---

  relativeSurplus: (input: RelativeSurplusInput) =>
    http<RelativeSurplusResult>("/v1/production/relative-surplus", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  relativeSurplusFromSource: (input: RelativeSurplusFromSourceInput) =>
    http<RelativeSurplusResult>("/v1/production/relative-surplus-from-productivity", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- simulation-engine (Ch. 16: Absolute and Relative Surplus-Value) ---

  computeAbsoluteSurplusValue: (input: AbsoluteSurplusValueInput) =>
    http<AbsoluteSurplusValueResult>("/v1/surplus-value/absolute", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeRelativeSurplusValue: (input: RelativeSurplusValueInput) =>
    http<RelativeSurplusValueResult>("/v1/surplus-value/relative", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getSurplusValueRate: (params: {
    surplus_labour: number;
    necessary_labour: number;
    total_capital: number;
    surplus_value?: number;
  }) => {
    const qs = new URLSearchParams();
    qs.set("surplus_labour", String(params.surplus_labour));
    qs.set("necessary_labour", String(params.necessary_labour));
    qs.set("total_capital", String(params.total_capital));
    if (params.surplus_value !== undefined) qs.set("surplus_value", String(params.surplus_value));
    return http<SurplusValueRateResult>(`/v1/surplus-value/rate?${qs.toString()}`);
  },

  // --- agent-service (Ch. 17: Changes of Magnitude in the Price of Labour-Power and in Surplus-Value) ---

  computeLabourScenario: (input: LabourScenarioInput) =>
    http<LabourScenarioResult>("/v1/labour-scenarios", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- simulation-engine (Ch. 18: Various Formula for the Rate of Surplus-Value) ---

  computeRatesOfSurplusValue: (input: RatesOfSurplusValueInput) =>
    http<RatesOfSurplusValueResult>("/v1/surplus-value/rates", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- agent-service (Ch. 19: The Transformation of the Value of Labour-Power into Wages) ---

  createWageForm: (input: CreateWageFormInput) =>
    http<WageForm>("/v1/wage-forms", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getWageForm: (agentID: string) =>
    http<WageForm>(`/v1/wage-forms/${agentID}`),

  listWageForms: () =>
    http<WageForm[]>("/v1/wage-forms"),

  // --- agent-service (Ch. 20: Time-Wages) ---

  computeHourlyPrice: (input: ComputeHourlyPriceInput) =>
    http<HourlyPriceOfLabour>("/v1/time-wages/hourly-price", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createWorkingSession: (input: CreateWorkingSessionInput) =>
    http<WorkingSession>("/v1/time-wages/sessions", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getWorkingSession: (id: string) =>
    http<WorkingSession>(`/v1/time-wages/sessions/${id}`),

  listWorkingSessions: (agentID: string) =>
    http<WorkingSession[]>(`/v1/agents/${agentID}/time-wages/sessions`),

  // --- agent-service (Ch. 21: Piece-Wages) ---

  computePiecePrice: (input: ComputePiecePriceInput) =>
    http<ComputePiecePriceResult>("/v1/piece-price", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createPieceWage: (agentID: string, input: CreatePieceWageInput) =>
    http<PieceWage>(`/v1/agents/${agentID}/piece-wages`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getPieceWage: (agentID: string) =>
    http<PieceWage>(`/v1/agents/${agentID}/piece-wages`),

  createSubContract: (input: CreateSubContractInput) =>
    http<SubContract>("/v1/sub-contracts", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getSubContract: (id: string) =>
    http<SubContract>(`/v1/sub-contracts/${id}`),

  // --- agent-service (Ch. 22: National Differences of Wages) ---

  registerIntensity: (input: RegisterIntensityInput) =>
    http<NationalIntensity>("/v1/intensities", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listIntensities: () =>
    http<{ items: NationalIntensity[] }>("/v1/intensities"),

  registerDayWage: (input: RegisterDayWageInput) =>
    http<DayWage>("/v1/wages", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getStandardisedWage: (country: string, referenceDayMinutes = 600) =>
    http<StandardisedWage>(`/v1/wages/${country}/standardised?reference_day_minutes=${referenceDayMinutes}`),

  getWageComparison: (referenceDayMinutes = 600) =>
    http<WageComparison>(`/v1/comparisons?reference_day_minutes=${referenceDayMinutes}`),

  // --- simulation-engine (Ch. 23: Simple Reproduction) ---

  runSimpleReproduction: (input: SimpleReproductionInput) =>
    http<SimpleReproductionResult>("/v1/reproductions/simple", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeRepaymentPeriod: (input: RepaymentPeriodInput) =>
    http<RepaymentPeriodResult>("/v1/reproductions/repayment-period", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- simulation-engine (Ch. 24: The Transformation of Surplus-Value into Capital) ---

  runExtendedReproduction: (input: ExtendedReproductionInput) =>
    http<ExtendedReproductionResult>("/v1/reproductions/extended", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  splitSurplus: (input: SplitSurplusInput) =>
    http<SplitSurplusResult>("/v1/reproductions/split-surplus", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- simulation-engine (Ch. 25: The General Law of Capitalist Accumulation) ---

  computeOrganicComposition: (input: OrganicCompositionInput) =>
    http<OrganicCompositionResult>("/v1/accumulation/organic-composition", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeLabourDemand: (input: LabourDemandInput) =>
    http<LabourDemandResult>("/v1/accumulation/labour-demand", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeReserveArmy: (input: ReserveArmyInput) =>
    http<ReserveArmyResult>("/v1/accumulation/reserve-army", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createGeneralLawScenario: (input: GeneralLawScenarioInput) =>
    http<GeneralLawScenarioResult>("/v1/accumulation/scenarios", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getGeneralLawScenario: (id: string) =>
    http<GeneralLawScenarioResult>(`/v1/accumulation/scenarios/${id}`),

  // --- simulation-engine (Ch. 26: The Secret of Primitive Accumulation) ---

  createHistoricalStage: (input: HistoricalStageInput) =>
    http<HistoricalStage>("/v1/historical-stages", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listHistoricalStages: () =>
    http<HistoricalStage[]>("/v1/historical-stages"),

  seedScenarioFromStage: (id: string, input: SeedScenarioInput) =>
    http<SeedScenarioResult>(`/v1/historical-stages/${id}/seed-scenario`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // --- simulation-engine (Ch. 27: Expropriation of the Agricultural Population) ---

  createEnclosureEvent: (input: EnclosureEventInput) =>
    http<EnclosureEvent>("/v1/enclosure-events", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listEnclosureEvents: () =>
    http<EnclosureEvent[]>("/v1/enclosure-events"),

  // --- simulation-engine (Ch. 28: Bloody Legislation Against the Expropriated) ---

  createWageStatute: (stageId: string, input: WageStatuteInput) =>
    http<WageStatute>(`/v1/historical-stages/${stageId}/wage-statutes`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createVagrancyLaw: (stageId: string, input: VagrancyLawInput) =>
    http<VagrancyLaw>(`/v1/historical-stages/${stageId}/vagrancy-laws`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getLabourDiscipline: (stageId: string) =>
    http<LabourDisciplineRegime>(`/v1/historical-stages/${stageId}/labour-discipline`),

  compareStatutoryWage: (acted: number, market: number) =>
    http<StatutoryWage>("/v1/statutory-wages/compare", {
      method: "POST",
      body: JSON.stringify({ acted_wage_pence: acted, market_wage_pence: market }),
    }),

  // --- simulation-engine (Ch. 29: Genesis of the Capitalist Farmer) ---

  createFarmTenure: (stageId: string, input: FarmTenureInput) =>
    http<FarmTenure>(`/v1/historical-stages/${stageId}/farm-tenures`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listFarmTenures: (stageId: string) =>
    http<FarmTenure[]>(`/v1/historical-stages/${stageId}/farm-tenures`),

  computeRealRent: (nominalRentPence: number, depreciationFactor: number) =>
    http<RealRentResult>("/v1/farm-tenures/real-rent", {
      method: "POST",
      body: JSON.stringify({ nominal_rent_pence: nominalRentPence, depreciation_factor: depreciationFactor }),
    }),

  // --- simulation-engine (Ch. 30: Reaction of the Agricultural Revolution on Industry) ---

  createDomesticIndustry: (stageId: string, input: DomesticIndustryInput) =>
    http<DomesticIndustry>(`/v1/historical-stages/${stageId}/domestic-industries`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeMarketFormation: (input: MarketFormationInput) =>
    http<MarketFormation>("/v1/market-formation", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getHomeMarket: (stageId: string) =>
    http<HomeMarketSize>(`/v1/historical-stages/${stageId}/home-market`),

  // --- simulation-engine (Ch. 31: Genesis of the Industrial Capitalist) ---

  createCapitalOrigin: (stageId: string, input: CapitalOriginInput) =>
    http<CapitalOrigin>(`/v1/historical-stages/${stageId}/capital-origins`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createColonialTransfer: (stageId: string, input: ColonialTransferInput) =>
    http<ColonialTransfer>(`/v1/historical-stages/${stageId}/colonial-transfers`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createNationalDebt: (stageId: string, input: NationalDebtInput) =>
    http<NationalDebt>(`/v1/historical-stages/${stageId}/national-debts`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getIndustrialCapitalGenesis: (stageId: string) =>
    http<IndustrialCapitalGenesis>(`/v1/historical-stages/${stageId}/genesis`),

  // --- simulation-engine (Ch. 32: Historical Tendency of Capitalist Accumulation) ---

  getNegationOfNegation: () =>
    http<NegationOfNegationResponse>("/v1/accumulation/negation-of-negation"),

  runCentralisation: (input: RunCentralisationInput) =>
    http<AccumulationTrajectory>("/v1/accumulation/centralisation", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listAccumulationTrajectories: () =>
    http<AccumulationTrajectory[]>("/v1/accumulation/trajectories"),

  getAccumulationTrajectory: (id: string) =>
    http<AccumulationTrajectory>(`/v1/accumulation/trajectories/${id}`),

  // --- simulation-engine (Ch. 33: The Modern Theory of Colonisation) ---

  createColonialMarket: (input: CreateColonialLabourMarketInput) =>
    http<ColonialLabourMarket>("/v1/colonial-markets", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listColonialMarkets: () =>
    http<ColonialLabourMarket[]>("/v1/colonial-markets"),

  regulateColonialMarket: (id: string, input: RegulateColonialMarketInput) =>
    http<ColonialLabourMarket>(`/v1/colonial-markets/${id}/regulate`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  computeIndependence: (id: string, input: ComputeIndependenceInput) =>
    http<WageWorkerIndependence>(`/v1/colonial-markets/${id}/independence`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  // Vol. II Ch. 2 — The Circuit of Productive Capital
  listProductiveCircuits: () =>
    http<ProductiveCircuit[]>("/v1/productive-circuits"),

  createProductiveCircuit: (input: CreateProductiveCircuitInput) =>
    http<ProductiveCircuit>("/v1/productive-circuits", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getProductiveCircuit: (id: string) =>
    http<ProductiveCircuit>(`/v1/productive-circuits/${id}`),

  recordRevenue: (id: string, amount: number) =>
    http<RevenueCircuit>(`/v1/productive-circuits/${id}/revenue`, {
      method: "POST",
      body: JSON.stringify({ amount }),
    }),

  accumulateLatentCapital: (id: string, amount: number) =>
    http<LatentMoneyCapital>(`/v1/productive-circuits/${id}/accumulate`, {
      method: "POST",
      body: JSON.stringify({ amount }),
    }),

  capitaliseCircuit: (id: string, amount: number, delta_constant_pence: number, delta_variable_pence: number) =>
    http<CapitalisationStep>(`/v1/productive-circuits/${id}/capitalise`, {
      method: "POST",
      body: JSON.stringify({ amount, delta_constant_pence, delta_variable_pence }),
    }),

  depositReserve: (id: string, amount: number) =>
    http<ReserveFund>(`/v1/productive-circuits/${id}/reserve/deposit`, {
      method: "POST",
      body: JSON.stringify({ amount }),
    }),

  drawReserve: (id: string, amount: number, reason: string) =>
    http<ReserveDraw>(`/v1/productive-circuits/${id}/reserve/draw`, {
      method: "POST",
      body: JSON.stringify({ amount, reason }),
    }),
};
