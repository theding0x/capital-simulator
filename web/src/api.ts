import type {
  MoneyCircuit,
  CreateMoneyCircuitInput,
  ProductiveState,
  MoneyCircuitCommodityCapital,
  Realisation,
  CircuitMoment,
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
  CommodityCircuit,
  CreateCommodityCircuitInput,
  SuccessivePartialSale,
  MeansOfProductionSource,
  CommodityCircuitAggregate,
  TurnoverTime,
  CreateTurnoverTimeInput,
  SellingPhase,
  BuyingPhase,
  NaturalProcessSpan,
  LatentProductiveCapital,
  Perishability,
  MarketSeparation,
  ActiveFractionResponse,
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

  // Vol. II Ch. 10 — Theories of Fixed and Circulating Capital: Physiocrats and Adam Smith.
  listEconomistAttributions: (theorist?: string) => {
    const url = theorist
      ? `/v1/circulation/economist-attributions?theorist=${encodeURIComponent(theorist)}`
      : "/v1/circulation/economist-attributions";
    return http<import("./types").EconomistAttribution[]>(url);
  },

  // Vol. II Ch. 12 — The Working Period
  createWorkingPeriod: (input: import("./types").CreateWorkingPeriodInput) =>
    http<import("./types").WorkingPeriod>("/v1/working-periods", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listWorkingPeriods: (params?: {
    industrial_capital_id?: string;
    mode?: string;
    commodity_id?: string;
  }) => {
    const q = new URLSearchParams();
    if (params?.industrial_capital_id)
      q.set("industrial_capital_id", params.industrial_capital_id);
    if (params?.mode) q.set("mode", params.mode);
    if (params?.commodity_id) q.set("commodity_id", params.commodity_id);
    const qs = q.toString();
    return http<import("./types").WorkingPeriod[]>(
      qs ? `/v1/working-periods?${qs}` : "/v1/working-periods"
    );
  },

  getWorkingPeriod: (id: string) =>
    http<import("./types").WorkingPeriod>(`/v1/working-periods/${id}`),

  recordWorkingPeriodInterruption: (
    id: string,
    input: import("./types").CreateWorkingPeriodInterruptionInput
  ) =>
    http<import("./types").WorkingPeriodInterruption>(
      `/v1/working-periods/${id}/interruptions`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  recordWorkingPeriodShortening: (
    id: string,
    input: import("./types").CreateWorkingPeriodShorteningInput
  ) =>
    http<import("./types").WorkingPeriodShortening>(
      `/v1/working-periods/${id}/shortenings`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  recordWorkingPeriodCreditFinancing: (
    id: string,
    input: import("./types").CreateWorkingPeriodCreditFinancingInput
  ) =>
    http<import("./types").WorkingPeriodCreditFinancing>(
      `/v1/working-periods/${id}/credit-financing`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  createNaturalConstraint: (
    input: import("./types").CreateNaturalConstraintInput
  ) =>
    http<import("./types").NaturalWorkingPeriodConstraint>(
      "/v1/natural-constraints",
      { method: "POST", body: JSON.stringify(input) }
    ),

  listNaturalConstraints: () =>
    http<import("./types").NaturalWorkingPeriodConstraint[]>(
      "/v1/natural-constraints"
    ),

  // Vol. II Ch. 13 — The Time of Production
  recordProductionLabourGap: (
    id: string,
    input: import("./types").CreateProductionLabourGapInput
  ) =>
    http<import("./types").ProductionLabourGap>(
      `/v1/working-periods/${id}/production-labour-gap`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  getProductionLabourGap: (id: string) =>
    http<import("./types").ProductionLabourGap>(
      `/v1/working-periods/${id}/production-labour-gap`
    ),

  recordNaturalProcessActivation: (
    id: string,
    input: import("./types").CreateNaturalProcessActivationInput
  ) =>
    http<import("./types").NaturalProcessActivation>(
      `/v1/working-periods/${id}/natural-activations`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  listNaturalProcessActivations: (id: string) =>
    http<import("./types").NaturalProcessActivation[]>(
      `/v1/working-periods/${id}/natural-activations`
    ),

  recordNaturalSubject: (
    id: string,
    input: import("./types").CreateNaturalSubjectInput
  ) =>
    http<import("./types").NaturalSubject>(
      `/v1/working-periods/${id}/natural-subject`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  getNaturalSubject: (id: string) =>
    http<import("./types").NaturalSubject>(
      `/v1/working-periods/${id}/natural-subject`
    ),

  createIndustryBenchmark: (
    input: import("./types").CreateIndustryBenchmarkInput
  ) =>
    http<import("./types").NaturalProcessIndustry>(
      "/v1/industry-benchmarks",
      { method: "POST", body: JSON.stringify(input) }
    ),

  listIndustryBenchmarks: () =>
    http<import("./types").NaturalProcessIndustry[]>(
      "/v1/industry-benchmarks"
    ),
};

// Vol. II Ch. 3 — The Circuit of Commodity-Capital
export const commodityCircuitApi = {
  list: (agentId?: string) => {
    const qs = agentId ? `?agent_id=${encodeURIComponent(agentId)}` : "";
    return http<CommodityCircuit[]>(`/v1/commodity-circuits${qs}`);
  },

  create: (input: CreateCommodityCircuitInput) =>
    http<CommodityCircuit>("/v1/commodity-circuits", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  get: (id: string) => http<CommodityCircuit>(`/v1/commodity-circuits/${id}`),

  aggregate: () => http<CommodityCircuitAggregate>("/v1/commodity-circuits/aggregate"),

  recordPartialSale: (id: string, quantity: number, realisedPence: number) =>
    http<SuccessivePartialSale>(`/v1/commodity-circuits/${id}/partial-sales`, {
      method: "POST",
      body: JSON.stringify({ quantity, realised_pence: realisedPence }),
    }),

  linkMPSource: (id: string, sourceKind: MeansOfProductionSource["source_kind"], sourceCommodityCircuitId?: string) =>
    http<MeansOfProductionSource>(`/v1/commodity-circuits/${id}/mp-source`, {
      method: "POST",
      body: JSON.stringify({
        source_kind: sourceKind,
        ...(sourceCommodityCircuitId ? { source_commodity_circuit_id: sourceCommodityCircuitId } : {}),
      }),
    }),

  close: (id: string, terminal: { constant_pence: number; variable_pence: number; surplus_pence: number; capitalised_pence: number; pounds_total: number }) =>
    http<CommodityCircuit>(`/v1/commodity-circuits/${id}/close`, {
      method: "POST",
      body: JSON.stringify(terminal),
    }),
};

// Vol. II Ch. 1 — The Circuit of Money-Capital
export const moneyCircuitApi = {
  list: (agentId?: string, moment?: CircuitMoment) => {
    const qs = new URLSearchParams();
    if (agentId) qs.set("agent_id", agentId);
    if (moment) qs.set("moment", moment);
    const q = qs.toString();
    return http<MoneyCircuit[]>(`/v1/money-circuits${q ? `?${q}` : ""}`);
  },

  create: (input: CreateMoneyCircuitInput) =>
    http<MoneyCircuit>("/v1/money-circuits", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  get: (id: string) => http<MoneyCircuit>(`/v1/money-circuits/${id}`),

  recordPurchase: (
    id: string,
    legs: {
      labour_amount: number;
      labour_power_hours: number;
      means_amount: number;
      means_capacity_hours: number;
    }
  ) =>
    http<MoneyCircuit>(`/v1/money-circuits/${id}/purchase`, {
      method: "POST",
      body: JSON.stringify(legs),
    }),

  recordProduce: (
    id: string,
    state: Pick<ProductiveState, "constant_pence" | "variable_pence"> & { entered_at?: string }
  ) =>
    http<MoneyCircuit>(`/v1/money-circuits/${id}/produce`, {
      method: "POST",
      body: JSON.stringify(state),
    }),

  recordCommodity: (
    id: string,
    cc: Pick<MoneyCircuitCommodityCapital, "value_original" | "value_surplus"> & {
      commodity_id?: string;
    }
  ) =>
    http<MoneyCircuit>(`/v1/money-circuits/${id}/commodity`, {
      method: "POST",
      body: JSON.stringify(cc),
    }),

  recordRealise: (
    id: string,
    rl: Pick<Realisation, "realised_pence"> & { sold_at?: string }
  ) =>
    http<MoneyCircuit>(`/v1/money-circuits/${id}/realise`, {
      method: "POST",
      body: JSON.stringify(rl),
    }),
};

// Vol. II Ch. 5 — The Time of Circulation
export const turnoverTimeApi = {
  list: () =>
    http<TurnoverTime[]>("/v1/turnover-time"),

  create: (input: CreateTurnoverTimeInput) =>
    http<TurnoverTime>("/v1/turnover-time", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  get: (id: string) =>
    http<TurnoverTime>(`/v1/turnover-time/${id}`),

  addLabourTime: (id: string, nanos: number) =>
    http<TurnoverTime>(`/v1/turnover-time/${id}/labour-time`, {
      method: "POST",
      body: JSON.stringify({ nanos }),
    }),

  addLabourInterruption: (id: string, nanos: number) =>
    http<TurnoverTime>(`/v1/turnover-time/${id}/labour-interruption`, {
      method: "POST",
      body: JSON.stringify({ nanos }),
    }),

  recordLatentMP: (id: string, pence: number, industrial_capital_id: string) =>
    http<LatentProductiveCapital>(`/v1/turnover-time/${id}/latent-mp`, {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, pence, held_at: new Date().toISOString() }),
    }),

  recordNaturalProcess: (id: string, industrial_capital_id: string, process: NaturalProcessSpan["process"], duration_nanos: number) =>
    http<NaturalProcessSpan>(`/v1/turnover-time/${id}/natural-process`, {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, process, duration_nanos, started_at: new Date().toISOString() }),
    }),

  openSellingPhase: (id: string, industrial_capital_id: string, commodity_circuit_id: string) =>
    http<SellingPhase>(`/v1/turnover-time/${id}/selling-phase`, {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, commodity_circuit_id, opened_at: new Date().toISOString() }),
    }),

  closeSellingPhase: (id: string, selling_phase_id: string, outcome: SellingPhase["outcome"]) =>
    http<SellingPhase>(`/v1/turnover-time/${id}/selling-phase`, {
      method: "POST",
      body: JSON.stringify({ selling_phase_id, outcome }),
    }),

  openBuyingPhase: (id: string, industrial_capital_id: string, money_circuit_id: string, market_location: string) =>
    http<BuyingPhase>(`/v1/turnover-time/${id}/buying-phase`, {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, money_circuit_id, market_location, opened_at: new Date().toISOString() }),
    }),

  closeBuyingPhase: (id: string, buying_phase_id: string) =>
    http<BuyingPhase>(`/v1/turnover-time/${id}/buying-phase`, {
      method: "POST",
      body: JSON.stringify({ buying_phase_id }),
    }),

  getActiveFraction: (id: string) =>
    http<ActiveFractionResponse>(`/v1/turnover-time/${id}/active-fraction`),

  setPerishability: (commodity_id: string, window_nanos: number) =>
    http<Perishability>("/v1/perishability", {
      method: "POST",
      body: JSON.stringify({ commodity_id, window_nanos }),
    }),

  setMarketSeparation: (industrial_capital_id: string, selling_market_id: string, buying_market_id: string) =>
    http<MarketSeparation>("/v1/market-separation", {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, selling_market_id, buying_market_id }),
    }),
};

// Vol. II Ch. 6 — Costs of Circulation
export const circulationCostsApi = {
  list: (industrial_capital_id?: string, kind?: string, nature?: string) => {
    const params = new URLSearchParams();
    if (industrial_capital_id) params.set("industrial_capital_id", industrial_capital_id);
    if (kind) params.set("kind", kind);
    if (nature) params.set("nature", nature);
    const qs = params.toString();
    return http<{ items: import("./types").CirculationCost[] }>(`/v1/circulation-costs${qs ? `?${qs}` : ""}`);
  },

  get: (id: string) =>
    http<import("./types").CirculationCost>(`/v1/circulation-costs/${id}`),

  create: (industrial_capital_id: string, kind: import("./types").CirculationCostKind, pence: number) =>
    http<import("./types").CirculationCost>("/v1/circulation-costs", {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, kind, pence }),
    }),

  aggregate: (industrial_capital_id: string, period?: string) => {
    const params = new URLSearchParams({ industrial_capital_id });
    if (period) params.set("period", period);
    return http<import("./types").AggregateCirculationCostsResult>(`/v1/circulation-costs/aggregate?${params}`);
  },

  systemFauxFrais: (period?: string) => {
    const qs = period ? `?period=${encodeURIComponent(period)}` : "";
    return http<import("./types").SystemFauxFraisResult>(`/v1/circulation-costs/system-faux-frais${qs}`);
  },

  createAgent: (
    industrial_capital_id: string,
    role: string,
    wage_rate_pence: number,
    labour_minutes_necessary: number,
    labour_minutes_surplus: number,
  ) =>
    http<import("./types").CirculationAgent>("/v1/circulation-agents", {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, role, wage_rate_pence, labour_minutes_necessary, labour_minutes_surplus }),
    }),

  createMoneyAsFauxFrais: (pence: number, annual_replacement_pence: number, industrial_capital_id?: string) =>
    http<import("./types").MoneyAsFauxFrais>("/v1/money-as-faux-frais", {
      method: "POST",
      body: JSON.stringify({ pence, annual_replacement_pence, industrial_capital_id }),
    }),

  createCommoditySupply: (industrial_capital_id: string, pence: number, is_voluntary: boolean, commodity_id?: string) =>
    http<import("./types").CommoditySupply>("/v1/commodity-supplies", {
      method: "POST",
      body: JSON.stringify({ industrial_capital_id, pence, is_voluntary, commodity_id }),
    }),

  addStorageCost: (
    supply_id: string,
    building_pence: number,
    labour_pence: number,
    preservation_labour_pence: number,
  ) =>
    http<import("./types").StorageCost>(`/v1/commodity-supplies/${supply_id}/storage-cost`, {
      method: "POST",
      body: JSON.stringify({ building_pence, labour_pence, preservation_labour_pence }),
    }),

  setTransportTariff: (
    commodity_id: string,
    base_pence_per_ton_mile: number,
    fragility_multiplier_basis_points: number,
    breakage_risk_multiplier_basis_points: number,
  ) =>
    http<import("./types").TransportTariff>("/v1/transport-tariffs", {
      method: "POST",
      body: JSON.stringify({ commodity_id, base_pence_per_ton_mile, fragility_multiplier_basis_points, breakage_risk_multiplier_basis_points }),
    }),

  createTransportLeg: (leg: {
    commodity_id: string;
    quantity: number;
    origin_id: string;
    destination_id: string;
    distance_meters: number;
    weight_grams: number;
    labour_cost_pence: number;
    means_of_transport_pence: number;
    surplus_pence: number;
  }) =>
    http<import("./types").TransportLeg>("/v1/transport-legs", {
      method: "POST",
      body: JSON.stringify(leg),
    }),
};

// Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers
export const turnoversApi = {
  create: (t: { industrial_capital_id?: string; lens: string; turnover_time_minutes: number }) =>
    http<import("./types").Turnover>("/v1/turnovers", {
      method: "POST",
      body: JSON.stringify(t),
    }),

  get: (id: string) => http<import("./types").Turnover>(`/v1/turnovers/${id}`),

  list: (params?: { industrial_capital_id?: string; lens?: string }) => {
    const qs = new URLSearchParams();
    if (params?.industrial_capital_id) qs.set("industrial_capital_id", params.industrial_capital_id);
    if (params?.lens) qs.set("lens", params.lens);
    const q = qs.toString();
    return http<{ items: import("./types").Turnover[] }>(`/v1/turnovers${q ? "?" + q : ""}`);
  },

  recordCycle: (
    id: string,
    cycle: {
      started_at: string;
      ended_at: string;
      advance_pence: number;
      returned_pence: number;
      production_minutes: number;
      circulation_minutes: number;
    },
  ) =>
    http<import("./types").TurnoverCycle>(`/v1/turnovers/${id}/cycles`, {
      method: "POST",
      body: JSON.stringify(cycle),
    }),

  recomputeNumber: (id: string) =>
    http<import("./types").TurnoverNumber>(`/v1/turnovers/${id}/recompute-number`, {
      method: "POST",
    }),

  getNumber: (id: string) => http<import("./types").TurnoverNumber>(`/v1/turnovers/${id}/number`),
};


// Vol. II Ch. 8 — Fixed Capital and Circulating Capital
export const compositionApi = {
  createComponent: (c: Omit<import("./types").CapitalComponent, "id" | "entered_at">) =>
    http<import("./types").CapitalComponent>("/v1/capital-components", {
      method: "POST",
      body: JSON.stringify(c),
    }),

  getComponent: (id: string) =>
    http<import("./types").CapitalComponent>(`/v1/capital-components/${id}`),

  listComponents: (params?: { industrial_capital_id?: string; kind?: string; role?: string }) => {
    const q = new URLSearchParams();
    if (params?.industrial_capital_id) q.set("industrial_capital_id", params.industrial_capital_id);
    if (params?.kind) q.set("kind", params.kind);
    if (params?.role) q.set("role", params.role);
    const qs = q.toString();
    return http<import("./types").CapitalComponent[]>(`/v1/capital-components${qs ? "?" + qs : ""}`);
  },

  registerFixedItem: (item: Omit<import("./types").FixedCapitalItem, "id" | "purchased_at">) =>
    http<import("./types").FixedCapitalItem>("/v1/fixed-items", {
      method: "POST",
      body: JSON.stringify(item),
    }),

  getFixedItem: (id: string) =>
    http<import("./types").FixedCapitalItem>(`/v1/fixed-items/${id}`),

  recordSubcomponent: (itemId: string, sub: Omit<import("./types").FixedCapitalSubcomponent, "id" | "fixed_capital_item_id">) =>
    http<import("./types").FixedCapitalSubcomponent>(`/v1/fixed-items/${itemId}/subcomponents`, {
      method: "POST",
      body: JSON.stringify(sub),
    }),

  recordWear: (itemId: string, wear: Omit<import("./types").WearAndTear, "id" | "fixed_capital_item_id" | "calculated_at">) =>
    http<import("./types").WearAndTear>(`/v1/fixed-items/${itemId}/wear`, {
      method: "POST",
      body: JSON.stringify(wear),
    }),

  getSinkingFund: (itemId: string) =>
    http<import("./types").SinkingFundForItem>(`/v1/fixed-items/${itemId}/sinking-fund`),

  recordRepair: (itemId: string, repair: Omit<import("./types").Repair, "id" | "fixed_capital_item_id">) =>
    http<import("./types").Repair>(`/v1/fixed-items/${itemId}/repairs`, {
      method: "POST",
      body: JSON.stringify(repair),
    }),

  recordReinvestment: (itemId: string, r: Omit<import("./types").SinkingFundReinvestment, "id" | "fixed_capital_item_id" | "occurred_at">) =>
    http<import("./types").SinkingFundReinvestment>(`/v1/fixed-items/${itemId}/reinvestments`, {
      method: "POST",
      body: JSON.stringify(r),
    }),

  recordCirculatingCycle: (c: Omit<import("./types").CirculatingCycle, "id">) =>
    http<import("./types").CirculatingCycle>("/v1/circulating-cycles", {
      method: "POST",
      body: JSON.stringify(c),
    }),
};

// Vol. II Ch. 9 — The Aggregate Turnover of Advanced Capital
export const aggregateTurnoverApi = {
  create: (body: {
    industrial_capital_id: string;
    contributions: Array<{
      capital_component_id?: string;
      kind: string;
      advanced_pence: number;
      turnover_number_basis_points: number;
      difference_kind?: string;
    }>;
  }) =>
    http<import("./types").AggregateTurnover>("/v1/aggregate-turnovers", {
      method: "POST",
      body: JSON.stringify(body),
    }),

  list: (industrial_capital_id?: string) => {
    const qs = industrial_capital_id
      ? `?industrial_capital_id=${encodeURIComponent(industrial_capital_id)}`
      : "";
    return http<import("./types").AggregateTurnover[]>(`/v1/aggregate-turnovers${qs}`);
  },

  get: (id: string) =>
    http<import("./types").AggregateTurnover>(`/v1/aggregate-turnovers/${id}`),

  recompute: (id: string) =>
    http<import("./types").AggregateTurnover>(`/v1/aggregate-turnovers/${id}/recompute`, {
      method: "POST",
      body: "{}",
    }),

  recordLifetimeCycle: (id: string, duration_minutes: number) =>
    http<import("./types").LifetimeCycle>(`/v1/aggregate-turnovers/${id}/lifetime-cycle`, {
      method: "POST",
      body: JSON.stringify({ duration_minutes }),
    }),

  recordCrisisPhase: (id: string, phase: string, notes?: string) =>
    http<import("./types").CrisisCyclePhaseRecord>(
      `/v1/aggregate-turnovers/${id}/crisis-phase`,
      {
        method: "POST",
        body: JSON.stringify({ phase, notes }),
      }
    ),
};

// Vol. II Ch. 14 — The Time of Circulation (refined)
export const circulationTimeV2Api = {
  listDistanceLagRelations: () =>
    http<{ items: import("./types").DistanceLagRelation[] }>("/v1/distance-lag"),
  createDistanceLagRelation: (input: import("./types").CreateDistanceLagRelationInput) =>
    http<import("./types").DistanceLagRelation>("/v1/distance-lag", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  listSpeedImprovements: (market_id?: string) => {
    const qs = market_id ? `?market_id=${encodeURIComponent(market_id)}` : "";
    return http<{ items: import("./types").CirculationSpeedImprovement[] }>(`/v1/circulation-speed-improvements${qs}`);
  },
  createSpeedImprovement: (input: import("./types").CreateCirculationSpeedImprovementInput) =>
    http<import("./types").CirculationSpeedImprovement>("/v1/circulation-speed-improvements", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  upgradeSellingPhase: (turnoverID: string, phaseID: string, input: import("./types").CreateSellingPhaseDetailedInput) =>
    http<import("./types").SellingPhaseDetailed>(
      `/v1/turnover-time/${turnoverID}/selling-phase/${phaseID}/detail`,
      { method: "POST", body: JSON.stringify(input) }
    ),
  getAnnualSurplusPenalty: (turnoverID: string, period: string) =>
    http<import("./types").AnnualSurplusPenalty>(
      `/v1/turnover-time/${turnoverID}/annual-surplus-penalty?period=${encodeURIComponent(period)}`
    ),
};

// ── Vol. II Ch. 15 — The Effects of a Change of Prices ──────────────────────

export const priceRevolutionApi = {
  create: (input: import("./types").CreatePriceRevolutionInput) =>
    http<import("./types").PriceRevolutionRecord>("/v1/price-revolutions", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  get: (id: string) =>
    http<import("./types").PriceRevolutionRecord>(`/v1/price-revolutions/${id}`),

  recordPriceCase: (
    id: string,
    input: { fall_case?: string; rise_case?: string }
  ) =>
    http<import("./types").PriceCaseRecord>(
      `/v1/price-revolutions/${id}/price-case`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  recordInventoryRevaluation: (
    input: import("./types").RecordInventoryRevaluationInput
  ) =>
    http<import("./types").InventoryRevaluation>("/v1/inventory-revaluations", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  recordSpeculativeHold: (
    input: import("./types").RecordSpeculativeHoldInput
  ) =>
    http<import("./types").SpeculativeHold>("/v1/speculative-holds", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getCompound: (industrialCapitalId: string, period: string) =>
    http<import("./types").CompoundCapitalAdjustment>(
      `/v1/price-revolutions/compound?industrial_capital_id=${encodeURIComponent(industrialCapitalId)}&period=${encodeURIComponent(period)}`
    ),
};

// Vol. II Ch. 16 — The Turnover of Variable Capital
export const annualSurplusRateApi = {
  createAdvance: (input: import("./types").CreateVariableCapitalAdvanceInput) =>
    http<import("./types").VariableCapitalAdvance>("/v1/variable-capital-advances", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getAdvance: (id: string) =>
    http<import("./types").VariableCapitalAdvance>(`/v1/variable-capital-advances/${id}`),

  listAdvances: (industrialCapitalId?: string) =>
    http<{ items: import("./types").VariableCapitalAdvance[] }>(
      `/v1/variable-capital-advances${industrialCapitalId ? `?industrial_capital_id=${encodeURIComponent(industrialCapitalId)}` : ""}`
    ),

  recordPerCircuitRate: (
    id: string,
    input: { surplus_pence: number; variable_pence: number }
  ) =>
    http<import("./types").PerCircuitSurplusRate>(
      `/v1/variable-capital-advances/${id}/per-circuit-rate`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  computeAnnualRate: (id: string, period: string) =>
    http<import("./types").AnnualSurplusRate>(
      `/v1/variable-capital-advances/${id}/compute-annual-rate?period=${encodeURIComponent(period)}`,
      { method: "POST", body: "{}" }
    ),

  listAnnualRates: (industrialCapitalId?: string, period?: string) => {
    const params = new URLSearchParams();
    if (industrialCapitalId) params.set("industrial_capital_id", industrialCapitalId);
    if (period) params.set("period", period);
    const qs = params.toString();
    return http<{ items: import("./types").AnnualSurplusRate[] }>(
      `/v1/annual-surplus-rates${qs ? `?${qs}` : ""}`
    );
  },

  createContrast: (input: {
    capital_a_id: string;
    capital_b_id: string;
    period: string;
  }) =>
    http<import("./types").AnnualSurplusRateContrast>(
      "/v1/annual-surplus-rate-contrasts",
      { method: "POST", body: JSON.stringify(input) }
    ),

  getReproduction: (id: string, yearMinutes?: number) => {
    const qs = yearMinutes ? `?year_minutes=${yearMinutes}` : "";
    return http<import("./types").VariableCapitalReproduction>(
      `/v1/variable-capital-reproductions/${id}${qs}`
    );
  },
};

// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction
export const moneyCapitalApi = {
  listApportionments: () =>
    http<import("./types").MoneySupplyApportionment[]>("/v1/reproduction/apportionments"),

  listInterDepartmentSettlements: () =>
    http<import("./types").InterDepartmentSettlement[]>("/v1/reproduction/inter-department-settlements"),

  createApportionment: (input: {
    total_circulating_money_pence: number;
    department_i_reserve_pence: number;
    department_ii_reserve_pence: number;
    wage_rotation_fund_pence: number;
    idle_hoard_pence: number;
    period: string;
  }) =>
    http<import("./types").MoneySupplyApportionment>("/v1/reproduction/apportionments", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  checkApportionmentBalance: (id: string) =>
    http<import("./types").ApportionmentBalanceCheck>(
      `/v1/reproduction/apportionments/${id}/balance-check`
    ),

  createDepartmentReserve: (input: {
    department: string;
    reserve_pence: number;
    purpose: string;
    period: string;
  }) =>
    http<import("./types").DepartmentMoneyReserve>("/v1/reproduction/department-reserves", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createCirculatingMoneyMass: (input: {
    money_stock_pence: number;
    velocity_per_year_basis_points: number;
    period: string;
  }) =>
    http<import("./types").CirculatingMoneyMass>("/v1/reproduction/circulating-money-mass", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createWageRotationFund: (input: {
    fund_pence: number;
    wage_cycle_frequency: number;
    department: string;
    period: string;
  }) =>
    http<import("./types").WageRotationFund>("/v1/reproduction/wage-rotation-funds", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createInterDepartmentSettlement: (input: {
    from_department: string;
    to_department: string;
    amount_pence: number;
    purpose: string;
    period: string;
  }) =>
    http<import("./types").InterDepartmentSettlement>("/v1/reproduction/inter-department-settlements", {
      method: "POST",
      body: JSON.stringify(input),
    }),
};

// Vol. II Ch. 17 — The Circulation of Surplus-Value
export const surplusCirculationApi = {
  create: (input: { period: string; total_surplus_pence: number }) =>
    http<import("./types").SurplusCirculation>("/v1/reproduction/surplus-circulations", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  get: (id: string) =>
    http<import("./types").SurplusCirculation>(`/v1/reproduction/surplus-circulations/${id}`),

  list: () =>
    http<{ items: import("./types").SurplusCirculation[] }>("/v1/reproduction/surplus-circulations"),

  recordSource: (id: string, input: { source: string; pence: number }) =>
    http<import("./types").RealisationSourceEntry>(
      `/v1/reproduction/surplus-circulations/${id}/realisation-source`,
      { method: "POST", body: JSON.stringify(input) }
    ),

  socialAggregate: (period: string) =>
    http<import("./types").SocialCapitalAggregate>(
      `/v1/reproduction/social-aggregate?period=${encodeURIComponent(period)}`
    ),

  realisationPuzzles: () =>
    http<{ items: import("./types").RealisationPuzzle[] }>("/v1/reproduction/realisation-puzzles"),
};

// Vol. II Ch. 20 — Simple Reproduction
export const simpleReproductionApi = {
  createScheme: (input: { period: string }) =>
    http<import("./types").SimpleReproductionScheme>("/v1/reproduction/simple/schemes", {
      method: "POST", body: JSON.stringify(input),
    }),

  getScheme: (id: string) =>
    http<import("./types").SimpleReproductionScheme>(`/v1/reproduction/simple/schemes/${id}`),

  listSchemes: (period?: string) =>
    http<import("./types").SimpleReproductionScheme[]>(
      `/v1/reproduction/simple/schemes${period ? `?period=${encodeURIComponent(period)}` : ""}`
    ),

  addDepartment: (id: string, input: { department: string; constant_pence: number; variable_pence: number; surplus_pence: number; total_pence: number }) =>
    http<import("./types").DepartmentalCapital>(`/v1/reproduction/simple/schemes/${id}/departments`, {
      method: "POST", body: JSON.stringify(input),
    }),

  recordExchange: (id: string, input: { from_department: string; to_department: string; pence: number; kind: string; description: string }) =>
    http<import("./types").InterDepartmentExchange>(`/v1/reproduction/simple/schemes/${id}/exchanges`, {
      method: "POST", body: JSON.stringify(input),
    }),

  advanceTick: (id: string) =>
    http<import("./types").ReproductionTick>(`/v1/reproduction/simple/schemes/${id}/tick`, {
      method: "POST", body: "{}",
    }),

  balanceCheck: (id: string) =>
    http<import("./types").BalanceCheckResult>(`/v1/reproduction/simple/schemes/${id}/balance-check`),

  recordMoneyLoop: (id: string, input: { period: string; is_closed: boolean; net_flow_pence: number }) =>
    http<import("./types").MoneyClosedLoop>(`/v1/reproduction/simple/schemes/${id}/money-loop`, {
      method: "POST", body: JSON.stringify(input),
    }),
};

// Vol. II Ch. 21 — Extended Reproduction
export const extendedReproductionApi = {
  listSchemes: (period?: string) =>
    http<import("./types").ExtendedReproductionScheme[]>(
      `/v1/reproduction/extended/schemes${period ? `?period=${encodeURIComponent(period)}` : ""}`
    ),

  createScheme: (input: { period: string; accumulation_rate_i_bps: number; accumulation_rate_ii_bps: number }) =>
    http<import("./types").ExtendedReproductionScheme>("/v1/reproduction/extended/schemes", {
      method: "POST", body: JSON.stringify(input),
    }),

  getScheme: (id: string) =>
    http<import("./types").ExtendedReproductionScheme>(`/v1/reproduction/extended/schemes/${id}`),

  addDepartment: (id: string, input: { department: string; constant_pence: number; variable_pence: number; surplus_pence: number; total_pence: number }) =>
    http<import("./types").DepartmentalCapital>(`/v1/reproduction/extended/schemes/${id}/departments`, {
      method: "POST", body: JSON.stringify(input),
    }),

  createReinvestment: (id: string, input: { department: string; delta_constant_pence: number; delta_variable_pence: number; consumed_surplus_pence: number; cycle_number: number }) =>
    http<import("./types").Reinvestment>(`/v1/reproduction/extended/schemes/${id}/reinvestments`, {
      method: "POST", body: JSON.stringify(input),
    }),

  tickScheme: (id: string) =>
    http<import("./types").ExtendedReproductionScheme>(`/v1/reproduction/extended/schemes/${id}/tick`, {
      method: "POST", body: "{}",
    }),

  getMoneyRequirement: (id: string, baseMoneyPence?: number) =>
    http<import("./types").ExtendedMoneyRequirement>(
      `/v1/reproduction/extended/schemes/${id}/money-requirement${baseMoneyPence !== undefined ? `?base_money_pence=${baseMoneyPence}` : ""}`
    ),

  createMultiPeriodScheme: (input: { label: string; scheme_ids: string[] }) =>
    http<import("./types").MultiPeriodScheme>("/v1/reproduction/extended/multi-period", {
      method: "POST", body: JSON.stringify(input),
    }),

  getMultiPeriodScheme: (id: string) =>
    http<import("./types").MultiPeriodScheme>(`/v1/reproduction/extended/multi-period/${id}`),

  listCompositionShifts: (id: string) =>
    http<import("./types").CompositionShift[]>(`/v1/reproduction/extended/multi-period/${id}/composition-shift`),

  listGrowthLeads: (id: string) =>
    http<import("./types").DepartmentIGrowthLead[]>(`/v1/reproduction/extended/multi-period/${id}/growth-lead`),
};

// Vol. III Ch. 1 — Cost-Price and Profit
export const financeApi = {
  listCostPrices: () =>
    http<{ items: import("./types").CostPriceResponse[] }>("/v1/profit/cost-price").then((r) => r.items),

  getCostPrice: (id: string) =>
    http<import("./types").CostPriceResponse>(`/v1/profit/cost-price/${id}`),

  computeCostPrice: (input: {
    constant: number;
    variable: number;
    fixed_wear_and_tear?: number;
    fixed_advanced?: number;
  }) =>
    http<import("./types").CostPriceResponse>("/v1/profit/cost-price", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getProfitForm: (input: { cost_price: number; surplus_value: number }) =>
    http<import("./types").ProfitFormResponse>("/v1/profit/profit-form", {
      method: "POST",
      body: JSON.stringify(input),
    }),
};
