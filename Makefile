# Capital Simulator - top-level Makefile

GO          ?= go
GOFLAGS     ?=
BIN_DIR     ?= bin
SERVICES    := api-gateway commodity-service agent-service market-service simulation-engine
DOCKER_REPO ?= ghcr.io/theding0x/capital-simulator
TAG         ?= dev

.PHONY: all build test vet tidy clean $(addprefix run-,$(SERVICES)) docker docker-% k8s-apply k8s-delete web-install web-dev web-build

all: vet test build

build:
	@mkdir -p $(BIN_DIR)
	@for svc in $(SERVICES); do \
		echo "==> building $$svc"; \
		$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$$svc ./services/$$svc/cmd/$$svc || exit 1; \
	done

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR)

# run-<service>: run a single service from source
$(addprefix run-,$(SERVICES)):
	$(GO) run ./services/$(@:run-%=%)/cmd/$(@:run-%=%)

# docker: build all service images
docker: $(addprefix docker-,$(SERVICES))

docker-%:
	docker build -f services/$*/Dockerfile -t $(DOCKER_REPO)/$*:$(TAG) .

k8s-apply:
	kubectl apply -k deploy/k8s

k8s-delete:
	kubectl delete -k deploy/k8s

# Web (Vite + React)
web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build
