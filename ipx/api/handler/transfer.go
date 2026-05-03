package handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

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
// @Tags         eth
// @Accept       json
// @Produce      json
// @Param        body  body      transferRequest   true  "Transfer request"
// @Success      200   {object}  transferResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Failure      504   {object}  map[string]string
// @Router       /eth/transfers [post]
func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.From == "" || req.To == "" || req.Value == "" {
		writeError(w, http.StatusBadRequest, "from, to, value are required")
		return
	}

	value, ok := new(big.Int).SetString(req.Value, 10)
	if !ok {
		writeError(w, http.StatusBadRequest, "value must be a decimal wei amount")
		return
	}

	key, err := h.cfg.KeyByAddress(req.From)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to derive key")
		return
	}

	ctx := r.Context()

	chainIDHex, err := h.client.ChainID(ctx)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to get chain id")
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse chain id")
		return
	}

	fromAddr := evmKey.Address.String()

	nonceHex, err := h.client.GetTransactionCount(ctx, fromAddr, "pending")
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to get nonce")
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse nonce")
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(ctx)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to get tip cap")
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse tip cap")
		return
	}

	block, err := h.client.BlockByNumber(ctx, "latest")
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to get latest block")
		return
	}
	baseFee, err := util.HexToBigInt(block["baseFeePerGas"].(string))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to parse base fee")
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

	unsigned, err := core.Codec.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encode tx")
		return
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to sign tx")
		return
	}

	rawTxBytes, err := core.Codec.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encode signed tx")
		return
	}

	txHashHex, err := h.client.SendRawTransaction(ctx, "0x"+hex.EncodeToString(rawTxBytes))
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to send tx: "+err.Error())
		return
	}

	_, err = h.client.WaitForReceipt(context.Background(), txHashHex, 30*time.Second)
	if err != nil {
		writeError(w, http.StatusGatewayTimeout, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, transferResponse{TxHash: txHashHex})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
