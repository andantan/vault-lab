package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type RPCHandler struct {
	client *rpc.Client
}

func NewRPCHandler(client *rpc.Client) *RPCHandler {
	return &RPCHandler{client: client}
}

// ChainID godoc
// @Summary      Get chain ID
// @Description  Returns the chain ID of the connected network as a decimal integer
// @Tags         rpc
// @Produce      json
// @Success      200  {object}  ChainIDResponse
// @Failure      500  {object}  map[string]string
// @Router       /evm/rpc/chain-id [post]
func (h *RPCHandler) ChainID(w http.ResponseWriter, r *http.Request) {
	hex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}

	chainID, err := util.HexToUint64(hex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewChainIDResponse(chainID, hex))
}

// GasPrice godoc
// @Summary      Get gas price
// @Description  Returns the current gas price in decimal and hex
// @Tags         rpc
// @Produce      json
// @Success      200  {object}  GasPriceResponse
// @Failure      500  {object}  map[string]string
// @Router       /evm/rpc/gas-price [post]
func (h *RPCHandler) GasPrice(w http.ResponseWriter, r *http.Request) {
	hex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}

	gasPrice, err := util.HexToBigInt(hex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewGasPriceResponse(gasPrice, hex))
}

// BlockNumber godoc
// @Summary      Get block number
// @Description  Returns the latest block number in decimal and hex
// @Tags         rpc
// @Produce      json
// @Success      200  {object}  BlockNumberResponse
// @Failure      500  {object}  map[string]string
// @Router       /evm/rpc/block-number [post]
func (h *RPCHandler) BlockNumber(w http.ResponseWriter, r *http.Request) {
	hex, err := h.client.BlockNumber(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get block number: %s", err))
		return
	}

	blockNumber, err := util.HexToUint64(hex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse block number: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBlockNumberResponse(blockNumber, hex))
}

// Nonce godoc
// @Summary      Get nonce for address
// @Description  Returns the pending transaction count (nonce) for the given address as a decimal integer
// @Tags         rpc
// @Accept       json
// @Produce      json
// @Param        body  body      NonceRequest  true  "Address"
// @Success      200   {object}  NonceResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/rpc/nonce [post]
func (h *RPCHandler) Nonce(w http.ResponseWriter, r *http.Request) {
	req := new(NonceRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hex, err := h.client.GetTransactionCount(r.Context(), req.Address, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}

	nonce, err := util.HexToUint64(hex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewNonceResponse(nonce, hex))
}

// Balance godoc
// @Summary      Get balance for address
// @Description  Returns the balance for the given address at the specified block in decimal and hex
// @Tags         rpc
// @Accept       json
// @Produce      json
// @Param        body  body      BalanceRequest  true  "Address and block"
// @Success      200   {object}  BalanceResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/rpc/balance [post]
func (h *RPCHandler) Balance(w http.ResponseWriter, r *http.Request) {
	req := new(BalanceRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	balanceHex, err := h.client.GetBalance(r.Context(), req.Address, req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get balance: %s", err))
		return
	}

	wei, err := util.HexToBigInt(balanceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse balance: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBalanceResponse(wei, balanceHex))
}

// Transaction godoc
// @Summary      Get transaction by hash
// @Description  Returns transaction details for the given tx hash
// @Tags         rpc
// @Accept       json
// @Produce      json
// @Param        body  body      TransactionRequest  true  "Transaction hash"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/rpc/transaction [post]
func (h *RPCHandler) Transaction(w http.ResponseWriter, r *http.Request) {
	req := new(TransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.client.GetTransactionByHash(r.Context(), req.TxHash)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get transaction: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, result)
}

// TransactionReceipt godoc
// @Summary      Get transaction receipt by hash
// @Description  Returns the transaction receipt for the given tx hash
// @Tags         rpc
// @Accept       json
// @Produce      json
// @Param        body  body      TransactionRequest  true  "Transaction hash"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/rpc/transaction/receipt [post]
func (h *RPCHandler) TransactionReceipt(w http.ResponseWriter, r *http.Request) {
	req := new(TransactionReceiptRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.client.TransactionReceipt(r.Context(), req.TxHash)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get receipt: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, result)
}

// EstimateGas godoc
// @Summary      Estimate gas for a transaction
// @Description  Returns the estimated gas limit for the given transaction parameters
// @Tags         rpc
// @Accept       json
// @Produce      json
// @Param        body  body      EstimateGasRequest   true  "Transaction parameters"
// @Success      200   {object}  EstimateGasResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/rpc/estimate-gas [post]
func (h *RPCHandler) EstimateGas(w http.ResponseWriter, r *http.Request) {
	req := new(EstimateGasRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hex, err := h.client.EstimateGas(r.Context(), req.Params(), req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to estimate gas: %s", err))
		return
	}

	gasLimit, err := util.HexToUint64(hex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas limit: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewEstimateGasResponse(gasLimit, hex))
}

// Call godoc
// @Summary      Execute a read-only contract call
// @Description  Executes eth_call and returns the raw result
// @Tags         rpc
// @Accept       json
// @Produce      json
// @Param        body  body      CallRequest  true  "Call parameters"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/rpc/call [post]
func (h *RPCHandler) Call(w http.ResponseWriter, r *http.Request) {
	req := new(CallRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.client.CallContract(r.Context(), req.Params(), req.Block)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to call contract: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, map[string]string{"result": result})
}
