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

	"github.com/andantan/evmlab/api/handler/misc"
	"github.com/andantan/evmlab/api/handler/v1"
	"github.com/andantan/evmlab/api/handler/v2"
	"github.com/andantan/evmlab/api/handler/v3"
	"github.com/andantan/evmlab/api/handler/v4"
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
		rpcHandler := misc.NewRPCHandler(client)
		r.Post("/", rpcHandler.Raw)
		r.Post("/chain-id", rpcHandler.ChainID)
		r.Post("/block-number", rpcHandler.BlockNumber)
		r.Post("/nonce", rpcHandler.Nonce)
		r.Post("/balance", rpcHandler.Balance)
		r.Post("/code", rpcHandler.Code)
		r.Post("/transaction", rpcHandler.Transaction)
		r.Post("/transaction/receipt", rpcHandler.TransactionReceipt)
		r.Post("/transaction/send", rpcHandler.SendTransaction)
		r.Post("/fee/base", rpcHandler.BaseFeePerGas)
		r.Post("/fee/priority", rpcHandler.MaxPriorityFeePerGas)
		r.Post("/fee/max", rpcHandler.MaxFeePerGas)
		r.Post("/gas/price", rpcHandler.GasPrice)
		r.Post("/gas/estimate", rpcHandler.EstimateGas)
		r.Post("/call", rpcHandler.Call)
	})

	r.Route("/evm/abi", func(r chi.Router) {
		abi := misc.NewAbiHandler()
		r.Post("/selector", abi.Selector)
		r.Post("/decode/result", abi.DecodeResult)
		r.Post("/decode/call", abi.DecodeCall)
		r.Post("/encode", abi.Encode)
		r.Post("/encode/balance-of", abi.BalanceOfCalldata)
		r.Post("/encode/approve", abi.ApproveCalldata)
		r.Post("/encode/transfer", abi.TransferCalldata)
		r.Post("/encode/allowance", abi.AllowanceCalldata)
		r.Post("/encode/transfer-from", abi.TransferFromCalldata)
	})

	r.Route("/evm/tool", func(r chi.Router) {
		tool := misc.NewToolHandler()
		r.Post("/address/eip55", tool.EIP55)
		r.Post("/crypto/derive", tool.DeriveKey)
		r.Post("/unit/convert", tool.ConvertUnit)
	})

	r.Route("/evm/hash", func(r chi.Router) {
		hash := misc.NewHashHandler()
		r.Post("/keccak256/legacy", hash.Keccak256Legacy)
		r.Post("/keccak256/eip191", hash.Keccak256EIP191)
		r.Post("/keccak256/eip712", hash.Keccak256EIP712)
	})

	r.Route("/evm/sign", func(r chi.Router) {
		sign := misc.NewSignHandler(cfg)
		r.Post("/", sign.Sign)
		r.Post("/ecrecover", sign.Ecrecover)
		r.Post("/verify/by-public-key", sign.VerifyByPublicKey)
		r.Post("/verify/by-address", sign.VerifyByAddress)
		r.Post("/transaction/legacy", sign.SignLegacyTransaction)
		r.Post("/transaction/eip1559", sign.SignEIP1559Transaction)
	})

	r.Route("/evm/v1", func(r chi.Router) {
		tx := v1.NewTransactionHandler(cfg)
		r.Post("/transaction/legacy/build", tx.BuildLegacyTransaction)
		r.Post("/transaction/eip1559/build", tx.BuildEIP1559Transaction)
	})

	r.Route("/evm/v2", func(r chi.Router) {
		transfer := v2.NewTransactionHandler(client)
		r.Post("/transaction/native/legacy", transfer.BuildNativeLegacyTransaction)
		r.Post("/transaction/native/eip1559", transfer.BuildNativeEIP1559Transaction)
		r.Post("/transaction/erc20/legacy", transfer.BuildERC20LegacyTransaction)
		r.Post("/transaction/erc20/eip1559", transfer.BuildERC20EIP1559Transaction)
	})

	r.Route("/evm/v3", func(r chi.Router) {
		tx := v3.NewTransactionHandler(cfg, client)
		r.Post("/transaction/native/legacy", tx.BuildNativeLegacyTransaction)
		r.Post("/transaction/native/eip1559", tx.BuildNativeEIP1559Transaction)
		r.Post("/transaction/erc20/legacy", tx.BuildERC20LegacyTransaction)
		r.Post("/transaction/erc20/eip1559", tx.BuildERC20EIP1559Transaction)
	})

	r.Route("/evm/v4", func(r chi.Router) {
		tx := v4.NewTransactionHandler(cfg, client)
		r.Post("/transaction/native/legacy", tx.BuildNativeLegacyTransaction)
		r.Post("/transaction/native/eip1559", tx.BuildNativeEIP1559Transaction)
		r.Post("/transaction/erc20/legacy", tx.BuildERC20LegacyTransaction)
		r.Post("/transaction/erc20/eip1559", tx.BuildERC20EIP1559Transaction)
	})

	fmt.Println("Listening on", cfg.ServerAddr)
	return http.ListenAndServe(cfg.ServerAddr, r)
}
