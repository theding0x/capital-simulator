import type {
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
};
