PROJECT = vault-lab
DOCKER ?= docker

GETH_CONTAINER = vault-lab-geth
GETH_VOLUME = vault-lab-geth-data
BLOCKSCOUT_DB_VOLUME = vault-lab-blockscout-db-data

BIN_DIR = bin
BUILD_DIR = build

GO_DIR = go
GO_BIN = contract_deployer

EXPLORER_URL ?= http://localhost:3000

.PHONY: help \
	up down logs reset \
	geth-logs geth-attach geth-reset \
	explorer \
	compile deploy test \
	go-build go-run go-test \
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
	@echo "  make deploy          Compile and deploy to local geth (CONTRACT=<path>)"
	@echo "  make test            Run all tests"
	@echo ""
	@echo "  make go-build        Build Go binary"
	@echo "  make go-run          Run Go binary"
	@echo "  make go-test         Run Go tests"
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

compile:
	@[ -n "$(CONTRACT)" ] || (echo "error: CONTRACT is required (e.g. make compile CONTRACT=contracts/vault/MultiAccountVault.sol)" && exit 1)
	mkdir -p $(BUILD_DIR)
	npx solcjs --abi --bin --base-path . -o $(BUILD_DIR) $(CONTRACT)

deploy: go-build compile
	@[ -n "$(DEPLOYER)" ] || (echo "error: DEPLOYER is required (e.g. make deploy CONTRACT=... DEPLOYER=key0)" && exit 1)
	./$(BIN_DIR)/$(GO_BIN) --contract $(CONTRACT) --deployer $(DEPLOYER)

test: go-test

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

go-build: $(BIN_DIR)
	cd $(GO_DIR) && go mod tidy && go build -o ../$(BIN_DIR)/$(GO_BIN) ./cmd/contract_deployer

go-run: go-build
	./$(BIN_DIR)/$(GO_BIN)

go-test:
	cd $(GO_DIR) && go test ./...

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	rm -rf artifacts cache
