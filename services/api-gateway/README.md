# api-gateway

External HTTP entrypoint for the capital-simulator. The React UI talks to this service, and this service fans out to the domain services (commodity, agent, market, simulation-engine).

- Default port: `8080`
- Endpoints (current scaffold):
  - `GET /healthz` — liveness
  - `GET /readyz` — readiness
  - `GET /v1/info` — service metadata and downstream list

Run locally:

```bash
make run-api-gateway
```
