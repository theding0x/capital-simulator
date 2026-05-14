# commodity-service

The simulation's model of the **commodity** — the cell-form of bourgeois wealth from which Marx begins *Capital, Vol. I*.

This service implements **Capital Vol. I, Ch. 1 — The Commodity** in full:

- **§1.** The two factors: use-value and value, with value-magnitude as socially necessary labour-time (`SNLT`). The inverse-proportionality law between productivity and value is encoded in `ProductivityChange.Apply`.
- **§2.** The dual character of labour: every commodity records its `ConcreteLabour` (weaving, tailoring, smelting, …); value is computed by reducing concrete labour to homogeneous abstract human labour via `AsAbstractLabour`.
- **§3.** The four value-forms — Simple, Expanded, General, and Money — derived from a population of commodities through `SimpleFormOf`, `ExpandedFormOf`, `GeneralFormOf`, and `MoneyFormOf`. The money-form treats the universal equivalent as a piece of social information rather than a separate computation.
- **§4.** The fetishism of commodities, surfaced explicitly via `SocialRelationsOf` and the `/social-relations` endpoint: every exchange-value response can be unwrapped into its underlying labour relations.

## Persistence

MongoDB. The `commodities` collection has a unique case-insensitive index on `name`. Connection settings come from `MONGO_URI` and `MONGO_DATABASE` (see `pkg/mongo`).

For local iteration without a live MongoDB instance, set:

- `MONGO_DISABLED=true` to skip the dial entirely and use the in-memory store; or
- `FALLBACK_MEMORY=true` to fall back to memory on a failed dial.

## API

| Method | Path                                      | Description                                                                |
|--------|-------------------------------------------|----------------------------------------------------------------------------|
| GET    | `/healthz`                                | Liveness                                                                   |
| GET    | `/readyz`                                 | Readiness                                                                  |
| POST   | `/v1/commodities`                         | Register a commodity                                                       |
| GET    | `/v1/commodities`                         | List all commodities                                                       |
| GET    | `/v1/commodities/{id}`                    | Fetch one                                                                  |
| PATCH  | `/v1/commodities/{id}`                    | Partial update                                                             |
| DELETE | `/v1/commodities/{id}`                    | Delete                                                                     |
| POST   | `/v1/commodities/{id}/value`              | Compute value (labour-minutes) of a quantity                               |
| GET    | `/v1/commodities/{id}/value-form`         | Render a value-form: `?kind=simple\|expanded\|general\|money[&quote_id=][&money_id=]` |
| GET    | `/v1/commodities/{id}/social-relations`   | Un-mystified view: the labour relations behind exchange (§4)               |
| POST   | `/v1/exchange-ratio`                      | Pairwise ratio: `{ base_id, quote_id, base_qty }`                          |

The api-gateway forwards `/v1/commodities/*` and `/v1/exchange-ratio` to this service.

## Run locally

```bash
make run-commodity-service
# or via the full stack:
docker compose up commodity-service
```

## Marx's classic example as a quick smoke test

```bash
# 30 minutes of weaving per yard of linen
curl -X POST localhost:8081/v1/commodities -d '{
  "name": "linen",
  "use_value": {"description":"linen for clothing","unit":"yards"},
  "concrete_labour": {"kind":"weaving"},
  "snlt_per_unit": 30
}'

# 10 hours of tailoring per coat
curl -X POST localhost:8081/v1/commodities -d '{
  "name": "coat",
  "use_value": {"description":"a coat to wear","unit":"coats"},
  "concrete_labour": {"kind":"tailoring"},
  "snlt_per_unit": 600
}'

# 20 yards of linen = 1 coat
curl -X POST localhost:8081/v1/exchange-ratio -d '{
  "base_id": "<linen-id>", "quote_id": "<coat-id>", "base_qty": 20
}'
```
