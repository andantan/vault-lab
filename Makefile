PROJECT = vault-lab
DOCKER ?= docker

GETH_CONTAINER = vault-lab-geth
GETH_VOLUME = vault-lab-geth-data

BIN_DIR = bin

GO_DIR = go
GO_BIN = vaultctl

RUST_DIR = rust
RUST_BIN = vaultctl-rs

.PHONY: help \
	up down logs geth-logs geth-attach geth-reset \
	compile hardhat-test test \
	go-build go-run go-test \
	rust-build rust-run rust-test \
	clean

help:
	@echo "$(PROJECT) commands:"
	@echo ""
	@echo "  make up              Start local geth"
	@echo "  make down            Stop local geth"
	@echo "  make logs            Follow docker compose logs"
	@echo "  make geth-logs       Follow geth logs"
	@echo "  make geth-attach     Attach to geth IPC console"
	@echo "  make geth-reset      Remove local geth volume and restart"
	@echo ""
	@echo "  make compile         Compile Solidity contracts with Hardhat"
	@echo "  make hardhat-test    Run Hardhat tests"
	@echo "  make test            Run all available tests"
	@echo ""
	@echo "  make go-build        Build Go client"
	@echo "  make go-run          Run Go client"
	@echo "  make go-test         Run Go tests"
	@echo ""
	@echo "  make rust-build      Build Rust client"
	@echo "  make rust-run        Run Rust client"
	@echo "  make rust-test       Run Rust tests"
	@echo ""
	@echo "  make clean           Remove generated build outputs"

up:
	$(DOCKER) compose up -d --wait

down:
	$(DOCKER) compose down

logs:
	$(DOCKER) compose logs -f

geth-logs:
	$(DOCKER) logs -f $(GETH_CONTAINER)

geth-attach:
	$(DOCKER) exec -it $(GETH_CONTAINER) geth attach /root/.ethereum/geth.ipc

geth-reset:
	$(DOCKER) compose down
	$(DOCKER) volume rm $(GETH_VOLUME) || true
	$(DOCKER) compose up -d --wait

compile:
	npx hardhat compile

hardhat-test:
	npx hardhat test

test: hardhat-test go-test rust-test

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

go-build: $(BIN_DIR)
	cd $(GO_DIR) && go mod tidy && go build -o ../$(BIN_DIR)/$(GO_BIN) ./cmd/vaultctl

go-run: go-build
	./$(BIN_DIR)/$(GO_BIN)

go-test:
	@if [ -d "$(GO_DIR)" ]; then \
		cd $(GO_DIR) && go test ./...; \
	else \
		echo "skip go-test: $(GO_DIR) directory not found"; \
	fi

rust-build:
	@if [ -d "$(RUST_DIR)" ]; then \
		cd $(RUST_DIR) && cargo build; \
	else \
		echo "skip rust-build: $(RUST_DIR) directory not found"; \
	fi

rust-run:
	@if [ -d "$(RUST_DIR)" ]; then \
		cd $(RUST_DIR) && cargo run --bin $(RUST_BIN); \
	else \
		echo "skip rust-run: $(RUST_DIR) directory not found"; \
	fi

rust-test:
	@if [ -d "$(RUST_DIR)" ]; then \
		cd $(RUST_DIR) && cargo test; \
	else \
		echo "skip rust-test: $(RUST_DIR) directory not found"; \
	fi

clean:
	rm -rf $(BIN_DIR)
	rm -rf artifacts cache
	@if [ -d "$(RUST_DIR)" ]; then \
		cd $(RUST_DIR) && cargo clean; \
	fi