package v2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	coretypes "github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
	"github.com/ethereum/go-ethereum/common"
)

type TransferHandler struct {
	cfg    *config.Config
	client *rpc.Client
}

func NewTransferHandler(cfg *config.Config, client *rpc.Client) *TransferHandler {
	return &TransferHandler{cfg: cfg, client: client}
}

// Transfer godoc
// @Summary      Send ETH (EIP-1559)
// @Description  Signs and broadcasts a dynamic-fee ETH transfer transaction
// @Tags         transfer
// @Accept       json
// @Produce      json
// @Param        body  body      transferRequest   true  "Transfer request"
// @Success      200   {object}  transferResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Failure      504   {object}  map[string]string
// @Router       /evm/v2/transfers/native/eip1559 [post]
func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if req.From == "" || req.To == "" || req.Value == "" {
		handler.WriteError(w, http.StatusBadRequest, "from, to, value are required")
		return
	}

	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		handler.WriteError(w, http.StatusBadRequest, "value must be a decimal wei amount")
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("from: %s", err))
		return
	}

	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	ctx := r.Context()

	chainIDHex, err := h.client.ChainID(ctx)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	fromAddr := evmKey.Address.String()

	nonceHex, err := h.client.GetTransactionCount(ctx, fromAddr, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(ctx)
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	block, err := h.client.BlockByNumber(ctx, "latest")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get latest block: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(block["baseFeePerGas"].(string))
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	tx := &coretypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  21000,
		To:        new(common.HexToAddress(req.To)),
		Value:     value,
		Data:      nil,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}

	rawTxBytes, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	txHashHex, err := h.client.SendRawTransaction(ctx, "0x"+hex.EncodeToString(rawTxBytes))
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send tx: %s", err))
		return
	}

	_, err = h.client.WaitForReceipt(context.Background(), txHashHex, 30*time.Second)
	if err != nil {
		handler.WriteError(w, http.StatusGatewayTimeout, fmt.Sprintf("failed to wait for receipt: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, transferResponse{TxHash: txHashHex})
}
