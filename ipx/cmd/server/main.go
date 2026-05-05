// @title           evmlab API
// @version         1.0
// @description     Ethereum transaction API
// @host            localhost:33152
// @BasePath        /
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andantan/evmlab/api/handler/v1"
	"github.com/andantan/evmlab/api/handler/v2"
	_ "github.com/andantan/evmlab/docs"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := util.FindProjectRoot()
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(root, "config.yaml"))
	if err != nil {
		return err
	}

	client := rpc.NewClient(cfg.RPCURL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/evm/rpc", func(r chi.Router) {
		rpcHandler := v1.NewRPCHandler(client)
		r.Post("/chain-id", rpcHandler.ChainID)
		r.Post("/gas-price", rpcHandler.GasPrice)
		r.Post("/block-number", rpcHandler.BlockNumber)
		r.Post("/nonce", rpcHandler.Nonce)
		r.Post("/balance", rpcHandler.Balance)
		r.Post("/transaction", rpcHandler.Transaction)
		r.Post("/transaction/receipt", rpcHandler.TransactionReceipt)
		r.Post("/estimate-gas", rpcHandler.EstimateGas)
		r.Post("/call", rpcHandler.Call)

		tool := v1.NewToolHandler()
		r.Post("/tool/address/checksum/eip55", tool.ChecksumEIP55)
		r.Post("/tool/crypto/derive", tool.DeriveKey)
	})

	r.Route("/evm/tool", func(r chi.Router) {
		tool := v1.NewToolHandler()
		r.Post("/address/checksum/eip55", tool.ChecksumEIP55)
		r.Post("/crypto/derive", tool.DeriveKey)
	})

	r.Route("/evm/v1", func(r chi.Router) {
		hash := v1.NewHashHandler()
		r.Post("/hash/keccak256/legacy", hash.Keccak256Legacy)
		r.Post("/hash/keccak256/personal", hash.Keccak256Personal)

		sign := v1.NewSignHandler(cfg)
		r.Post("/sign", sign.Sign)
		r.Post("/sign/verify/by-public-key", sign.VerifyByPublicKey)
		r.Post("/sign/verify/by-address", sign.VerifyByAddress)

		tx := v1.NewTransactionHandler(cfg)
		r.Post("/transaction/legacy/build", tx.BuildLegacyNativeTransfer)
		r.Post("/transaction/legacy/sign", tx.SignLegacyNativeTransfer)
	})

	r.Route("/evm/v2", func(r chi.Router) {
		transfer := v2.NewTransferHandler(cfg, client)

		r.Post("/transfers/native/eip1559", transfer.Transfer)
	})

	fmt.Println("Listening on", cfg.ServerAddr)
	return http.ListenAndServe(cfg.ServerAddr, r)
}
