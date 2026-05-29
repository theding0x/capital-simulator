# Dockerfile Patterns

> Load when: Multi-stage builds, image optimization, security.

## Node.js Multi-Stage Build

```dockerfile
# Stage 1: Dependencies
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --only=production

# Stage 2: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 3: Production
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

RUN addgroup --system --gid 1001 appgroup
RUN adduser --system --uid 1001 appuser

COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./

USER appuser
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s CMD wget -q --spider http://localhost:3000/health || exit 1
CMD ["node", "dist/main.js"]
```

## Python Multi-Stage Build

```dockerfile
FROM python:3.12-slim AS builder
WORKDIR /app
RUN pip install --no-cache-dir poetry
COPY pyproject.toml poetry.lock ./
RUN poetry export -f requirements.txt -o requirements.txt
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

FROM python:3.12-slim
WORKDIR /app
COPY --from=builder /install /usr/local
COPY src/ src/
RUN useradd -r appuser && chown -R appuser /app
USER appuser
CMD ["python", "-m", "src.main"]
```

## Layer Caching Order

```dockerfile
# 1. Base image (rarely changes)
FROM node:20-alpine

# 2. System dependencies (rarely changes)
RUN apk add --no-cache tini

# 3. App dependencies (changes occasionally)
COPY package.json package-lock.json ./
RUN npm ci

# 4. App source (changes frequently)
COPY . .
RUN npm run build
```