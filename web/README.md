# web

Vite + React + TypeScript dashboard for the capital-simulator economy.

```bash
cd web
npm install
npm run dev   # http://localhost:5173 - proxies /api -> http://localhost:8080 (api-gateway)
npm run build
```

Production build is served by nginx (see `Dockerfile` and `nginx.conf`), with `/api/` proxied to the `api-gateway` service.
