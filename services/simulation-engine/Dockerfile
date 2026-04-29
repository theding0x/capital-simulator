# syntax=docker/dockerfile:1.7
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY pkg ./pkg
COPY services/simulation-engine ./services/simulation-engine
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/app ./services/simulation-engine/cmd/simulation-engine

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8084
ENTRYPOINT ["/app"]
