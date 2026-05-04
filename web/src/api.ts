import type {
  Agent,
  CapitalCircuit,
  Circuit,
  Commodity,
  ComputePriceInput,
  CreateAgentCircuitInput,
  CreateAgentInput,
  CreateCircuitInput,
  CreateCommodityInput,
  CreateExchangeInput,
  CreateHoardInput,
  CreateOfferInput,
  CreateOwnerInput,
  CreatePaymentObligationInput,
  CreateWorldMoneyTransferInput,
  Exchange,
  ExchangeRatio,
  Hoard,
  MoneyCommodity,
  Offer,
  Owner,
  PaymentObligation,
  Price,
  SocialRelations,
  UniversalEquivalent,
  UpdateAgentInput,
  UpdateCommodityInput,
  ValueResponse,
  WorldMoneyTransfer,
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
};
