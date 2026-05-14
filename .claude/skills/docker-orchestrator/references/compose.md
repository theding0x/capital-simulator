# Docker Compose Patterns

> Load when: Multi-service setups, networking, development environments.

## Production-Ready Compose

```yaml
services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
      target: runner
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgresql://user:pass@db:5432/mydb
      - REDIS_URL=redis://redis:6379
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: "1.0"

  db:
    image: postgres:16-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: mydb
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d mydb"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    volumes:
      - redisdata:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s

volumes:
  pgdata:
  redisdata:
```

## Development Compose (Override)

```yaml
# docker-compose.override.yml (auto-merged in dev)
services:
  api:
    build:
      target: builder  # Use build stage with dev deps
    volumes:
      - .:/app         # Hot reload
      - /app/node_modules  # Exclude node_modules
    command: npm run dev
    environment:
      - DEBUG=true
```