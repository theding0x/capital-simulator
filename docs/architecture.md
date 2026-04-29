# Architecture

Capital Simulator is a microservices simulation of an economy as described in Karl Marx's *Capital, Volume I*. The architecture is intentionally modular: each Marxist economic category (the commodity, the agent, the market, the production process, the circulation of capital) gets its own service so that chapter-by-chapter additions stay localized.

## Topology

```
                ┌──────────────┐
        browser │     web      │  React + Vite + TS, served by nginx
                └──────┬───────┘
                       │  /api/*
                       ▼
                ┌──────────────┐
                │ api-gateway  │  external HTTP entrypoint, fans out
                └──┬─┬─┬─┬─────┘
                   │ │ │ │
       ┌───────────┘ │ │ └────────────┐
       ▼             ▼ ▼              ▼
┌────────────┐ ┌──────────┐ ┌──────────────────┐
│ commodity- │ │  agent-  │ │ market-service   │
│  service   │ │ service  │ │                  │
└─────┬──────┘ └────┬─────┘ └────────┬─────────┘
      │             │                │
      └─────────────┼────────────────┘
                    │
            ┌───────▼────────────┐
            │ simulation-engine  │  drives ticks
            └───────┬────────────┘
                    │
   ┌────────────────┴───────────────┐
   ▼                                ▼
┌────────┐                       ┌───────┐
│ MongoDB │  durable state        │ Redis │  hot caches & tick state
└────────┘                       └───────┘
```

## Services

| Service              | Port  | Marxist role                                                            | Persistence       |
|----------------------|-------|-------------------------------------------------------------------------|-------------------|
| `api-gateway`        | 8080  | External entrypoint; fans out to domain services.                       | —                 |
| `commodity-service`  | 8081  | Use-value, exchange-value, value (Ch. 1).                               | MongoDB           |
| `agent-service`      | 8082  | Workers, capitalists, and other class-bearers (Ch. 4+).                 | MongoDB           |
| `market-service`     | 8083  | Exchange and circulation; C-M-C and M-C-M' (Ch. 2-3).                   | MongoDB + Redis   |
| `simulation-engine`  | 8084  | Time-step orchestrator; advances the economy one period at a time.      | MongoDB + Redis   |

All Go services share a single root `go.mod` (`github.com/theding0x/capital-simulator`). Cross-cutting concerns live under `pkg/`:

- `pkg/log` — structured logging via `log/slog`.
- `pkg/httpx` — HTTP server scaffolding with `/healthz`, `/readyz`, and graceful shutdown.
- `pkg/mongo` — MongoDB connection config (driver wired in by a later chapter).
- `pkg/redis` — Redis connection config (driver wired in by a later chapter).

## Data flow per simulation tick (target shape)

1. `simulation-engine` advances the clock and tells `agent-service` to act.
2. Agents form intentions to produce / sell / buy and notify `commodity-service` and `market-service`.
3. `market-service` matches trades and updates prices.
4. State changes are persisted to MongoDB; hot tick state is cached in Redis.
5. `api-gateway` exposes a read view of the world to the React UI.

This is the *target*; the initial scaffold ships health endpoints only.

## Roadmap (chapter-driven)

Each chapter of *Capital* turns into a feature branch and PR. Approximate mapping:

| Chapter | Concepts                                                 | Primary services              |
|---------|----------------------------------------------------------|-------------------------------|
| Ch. 1   | Commodity, use-value, value, exchange-value, fetishism   | commodity-service             |
| Ch. 2-3 | Exchange, money, hoarding, means of payment              | market-service, commodity     |
| Ch. 4   | Money → capital; the general formula M-C-M'              | agent-service, market         |
| Ch. 5-7 | Labour-process, valorization, surplus-value              | agent-service, simulation-eng |
| Ch. 8-9 | Constant/variable capital, rate of surplus-value         | commodity, simulation-eng     |
| Ch. 10  | The working day                                          | agent-service, simulation-eng |
| Ch. 11+ | Cooperation, machinery, wages, accumulation              | all                           |
