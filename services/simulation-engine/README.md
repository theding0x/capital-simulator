# simulation-engine

Time-step orchestrator that drives the simulated economy forward. On each tick it instructs the domain services (commodity, agent, market) to produce, exchange, and accumulate, then advances simulation time.

- Default port: `8084`
- Persistence: MongoDB for run history; Redis for live tick state

Run locally:

```bash
make run-simulation-engine
```
