PROJECT = evmlab
DOCKER ?= docker

GETH_CONTAINER = evmlab-geth
GETH_VOLUME = evmlab-geth-data
BLOCKSCOUT_DB_VOLUME = evmlab-blockscout-db-data

BIN_DIR = bin
BUILD_DIR = build

GO_DIR = ipx
GO_DEPLOYER_BIN = contract_deployer
GO_SERVER_BIN   = server

EXPLORER_URL ?= http://localhost:3000

.PHONY: help \
	up down logs reset \
	geth-logs geth-attach geth-reset \
	explorer \
	compile standard-json deploy test \
	go-build-deployer go-build-server go-build \
	server go-test swag \
	clean

help:
	@echo "$(PROJECT) commands:"
	@echo ""
	@echo "  make up              Start local geth and Blockscout"
	@echo "  make down            Stop local stack"
	@echo "  make logs            Follow docker compose logs"
	@echo "  make reset           Remove local geth and Blockscout data, then restart"
	@echo ""
	@echo "  make geth-logs       Follow geth logs"
	@echo "  make geth-attach     Attach to geth IPC console"
	@echo "  make geth-reset      Remove local geth volume and restart"
	@echo ""
	@echo "  make explorer        Print local Blockscout URL"
	@echo ""
	@echo "  make compile         Compile Solidity contract (CONTRACT=<path>)"
	@echo "  make standard-json   Generate standard JSON input (CONTRACT=<path>)"
	@echo "  make deploy          Compile and deploy to local geth (CONTRACT=<path>)"
	@echo "  make test            Run all tests"
	@echo ""
	@echo "  make go-build          Build all Go binaries"
	@echo "  make go-build-deployer Build contract_deployer binary"
	@echo "  make go-build-server   Build server binary"
	@echo "  make server            Build and run API server"
	@echo "  make go-test           Run Go tests"
	@echo "  make swag            Regenerate Swagger docs"
	@echo ""
	@echo "  make clean           Remove generated build outputs"

up:
	$(DOCKER) compose up -d --wait

down:
	$(DOCKER) compose down

logs:
	$(DOCKER) compose logs -f

reset:
	$(DOCKER) compose down -v
	$(DOCKER) compose up -d --wait

geth-logs:
	$(DOCKER) logs -f $(GETH_CONTAINER)

geth-attach:
	$(DOCKER) exec -it $(GETH_CONTAINER) geth attach /root/.ethereum/geth.ipc

geth-reset:
	$(DOCKER) compose down
	$(DOCKER) volume rm $(GETH_VOLUME) || true
	$(DOCKER) compose up -d --wait

explorer:
	@echo "$(EXPLORER_URL)"

CONTRACT_DIR    = $(shell dirname $(CONTRACT))
CONTRACT_SUBDIR = $(shell dirname $(CONTRACT) | sed 's|^contracts/||')
CONTRACT_NAME   = $(shell basename $(CONTRACT) .sol)

compile:
	@[ -n "$(CONTRACT)" ] || (echo "error: CONTRACT is required (e.g. make compile CONTRACT=contracts/vault/MultiAccountVault.sol)" && exit 1)
	mkdir -p $(BUILD_DIR)/$(CONTRACT_SUBDIR)
	npx solcjs --abi --bin --base-path . -o $(BUILD_DIR)/$(CONTRACT_SUBDIR) $(CONTRACT)

standard-json:
	@[ -n "$(CONTRACT)" ] || (echo "error: CONTRACT is required" && exit 1)
	@[ -f "$(CONTRACT_DIR)/build-standard-json.sh" ] || (echo "error: $(CONTRACT_DIR)/build-standard-json.sh not found" && exit 1)
	mkdir -p $(BUILD_DIR)/$(CONTRACT_SUBDIR)
	bash $(CONTRACT_DIR)/build-standard-json.sh $(CONTRACT_NAME)

deploy: go-build-deployer compile standard-json
	@[ -n "$(DEPLOYER)" ] || (echo "error: DEPLOYER is required (e.g. make deploy CONTRACT=... DEPLOYER=key0)" && exit 1)
	./$(BIN_DIR)/$(GO_DEPLOYER_BIN) --contract $(CONTRACT) --deployer $(DEPLOYER)

test: go-test

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

go-build-deployer: $(BIN_DIR)
	cd $(GO_DIR) && go build -o ../$(BIN_DIR)/$(GO_DEPLOYER_BIN) ./cmd/contract_deployer

go-build-server: $(BIN_DIR)
	cd $(GO_DIR) && go build -o ../$(BIN_DIR)/$(GO_SERVER_BIN) ./cmd/server

go-build: go-build-deployer go-build-server

server: go-build-server swag
	./$(BIN_DIR)/$(GO_SERVER_BIN)

go-test:
	cd $(GO_DIR) && go test ./...

swag:
	cd $(GO_DIR) && swag init -g cmd/server/main.go -o docs

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	rm -rf artifacts cache
