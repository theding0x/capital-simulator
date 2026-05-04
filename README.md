# Capital Simulator

A simulation of an economy modeled chapter by chapter on Karl Marx's *Capital, Volume I*. Each chapter is implemented as its own branch + pull request, accompanied by an HTML chapter summary that lives in `chapters/volume-1/`.

## Stack

- **Architecture:** Microservices (Go), single Go module monorepo
- **UI:** React + TypeScript (Vite)
- **Database:** MongoDB
- **Caching:** Redis
- **Platform:** Docker / Kubernetes
- **Repository host:** GitHub

## Repository layout

```
capital-simulator/
├── chapters/
│   └── volume-1/       HTML summaries, one per Capital Vol. I chapter
├── docs/               Architecture and design notes
├── deploy/
│   ├── docker/         Service Dockerfiles (also live next to each service)
│   └── k8s/            Kubernetes manifests
├── pkg/                Shared Go packages (log, httpx, mongo, redis)
├── services/
│   ├── api-gateway/        External entrypoint, fans out to domain services
│   ├── commodity-service/  Use-value, exchange-value, value (Ch. 1)
│   ├── agent-service/      Workers, capitalists, and other economic agents
│   ├── market-service/     Exchange, prices, circulation
│   └── simulation-engine/  Time-step orchestrator that drives the world
└── web/                Vite + React + TypeScript dashboard
```

## Workflow

We progress one chapter of *Capital* at a time:

1. Cut a branch named `chapter-NN-short-slug` off `main`.
2. Implement the economic concepts introduced in that chapter across the relevant services.
3. Drop an HTML summary of the chapter into `chapters/volume-1/` (e.g. `chapters/volume-1/01-the-commodity.html`).
4. Open a pull request whose description summarizes the chapter and the simulation changes it produced.
5. Merge into `main` once reviewed.

## Local development

Prerequisites: Go 1.22+, Node 20+, Docker, kubectl (for k8s deploys).

```bash
# Run everything locally (mongo + redis + all Go services)
docker compose up --build

# Build all Go binaries
make build

# Run the API gateway only
make run-api-gateway

# Run the React UI
cd web && npm install && npm run dev
```

## Deploying to Kubernetes

```bash
kubectl apply -k deploy/k8s
```

See `docs/architecture.md` for the high-level service topology.
