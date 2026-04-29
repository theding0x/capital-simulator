import type {
  Commodity,
  CreateCommodityInput,
  ExchangeRatio,
  SocialRelations,
  UpdateCommodityInput,
  ValueResponse,
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
};
