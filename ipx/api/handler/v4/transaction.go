package v4

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
)

type TransactionHandler struct {
	cfg    *config.Config
	client *rpc.Client
}

func NewTransactionHandler(cfg *config.Config, client *rpc.Client) *TransactionHandler {
	return &TransactionHandler{cfg: cfg, client: client}
}

// BuildNativeLegacyTransaction godoc
// @Summary      Build, sign, and send legacy transfer tx
// @Description  Fetches chain state, builds and signs an EIP-155 legacy native transfer, broadcasts it, and returns unsigned_rlp, signed_rlp, tx_hash, r, s, v
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeLegacyTransactionRequest  true  "Transfer request"
// @Success      200   {object}  BuildNativeLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v4/transaction/native/legacy [post]
func (h *TransactionHandler) BuildNativeLegacyTransaction(w http.ResponseWriter, r *http.Request) {
	var req BuildNativeLegacyTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
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

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 21000,
		To:       &req.ToAddr().Addr,
		Value:    req.Amount(),
		Data:     nil,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}

	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	if _, err := h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeLegacyTransactionResponse(unsigned, signed, txHash, sig))
}

// BuildNativeEIP1559Transaction godoc
// @Summary      Build, sign, and send EIP-1559 transfer tx
// @Description  Fetches chain state, builds and signs an EIP-1559 native transfer, broadcasts it, and returns unsigned_rlp, signed_rlp, tx_hash, r, s, v
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildNativeEIP1559TransactionRequest  true  "Transfer request"
// @Success      200   {object}  BuildNativeEIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v4/transaction/native/eip1559 [post]
func (h *TransactionHandler) BuildNativeEIP1559Transaction(w http.ResponseWriter, r *http.Request) {
	var req BuildNativeEIP1559TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
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

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  21000,
		To:        &req.ToAddr().Addr,
		Value:     req.Amount(),
		Data:      nil,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}

	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	if _, err := h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildNativeEIP1559TransactionResponse(unsigned, signed, txHash, sig))
}

// BuildERC20LegacyTransaction godoc
// @Summary      Build, sign, and send ERC-20 legacy transfer tx
// @Description  Builds transfer(address,uint256) calldata, estimates gas, signs with the configured key, broadcasts, and returns unsigned_rlp, signed_rlp, tx_hash, r, s, v
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildERC20LegacyTransactionRequest  true  "ERC-20 transfer request"
// @Success      200   {object}  BuildERC20LegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v4/transaction/erc20/legacy [post]
func (h *TransactionHandler) BuildERC20LegacyTransaction(w http.ResponseWriter, r *http.Request) {
	var req BuildERC20LegacyTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
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

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	gasPriceHex, err := h.client.GasPrice(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get gas price: %s", err))
		return
	}
	gasPrice, err := util.HexToBigInt(gasPriceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas price: %s", err))
		return
	}

	calldata := core.TransferCalldata(req.ToAddr(), req.ToAmount())

	gasEstHex, err := h.client.EstimateGas(r.Context(), map[string]any{
		"from":  req.From,
		"to":    req.ContractAddr().String(),
		"value": "0x0",
		"data":  "0x" + hex.EncodeToString(calldata),
	}, "latest")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to estimate gas: %s", err))
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.LegacyTx{
		ChainID:  chainID,
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasEst * 12 / 10,
		To:       &req.ContractAddr().Addr,
		Value:    big.NewInt(0),
		Data:     calldata,
	}

	unsigned, err := core.RLP.EncodeLegacyUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}

	signed, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	if _, err := h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildERC20LegacyTransactionResponse(unsigned, signed, txHash, sig))
}

// BuildERC20EIP1559Transaction godoc
// @Summary      Build, sign, and send ERC-20 EIP-1559 transfer tx
// @Description  Builds transfer(address,uint256) calldata, estimates gas, signs with the configured key, broadcasts, and returns unsigned_rlp, signed_rlp, tx_hash, r, s, v
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildERC20EIP1559TransactionRequest  true  "ERC-20 transfer request"
// @Success      200   {object}  BuildERC20EIP1559TransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      502   {object}  map[string]string
// @Router       /evm/v4/transaction/erc20/eip1559 [post]
func (h *TransactionHandler) BuildERC20EIP1559Transaction(w http.ResponseWriter, r *http.Request) {
	var req BuildERC20EIP1559TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
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

	chainIDHex, err := h.client.ChainID(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get chain id: %s", err))
		return
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse chain id: %s", err))
		return
	}

	nonceHex, err := h.client.GetTransactionCount(r.Context(), req.From, "pending")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get nonce: %s", err))
		return
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse nonce: %s", err))
		return
	}

	tipCapHex, err := h.client.MaxPriorityFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get tip cap: %s", err))
		return
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse tip cap: %s", err))
		return
	}

	baseFeeHex, err := h.client.BaseFeePerGas(r.Context())
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to get base fee: %s", err))
		return
	}
	baseFee, err := util.HexToBigInt(baseFeeHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse base fee: %s", err))
		return
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	calldata := core.TransferCalldata(req.ToAddr(), req.ToAmount())

	gasEstHex, err := h.client.EstimateGas(r.Context(), map[string]any{
		"from":  req.From,
		"to":    req.ContractAddr().String(),
		"value": "0x0",
		"data":  "0x" + hex.EncodeToString(calldata),
	}, "latest")
	if err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to estimate gas: %s", err))
		return
	}
	gasEst, err := util.HexToUint64(gasEstHex)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse gas estimate: %s", err))
		return
	}

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasEst * 12 / 10,
		To:        &req.ContractAddr().Addr,
		Value:     big.NewInt(0),
		Data:      calldata,
	}

	unsigned, err := core.RLP.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode unsigned tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign tx: %s", err))
		return
	}

	signed, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	if _, err := h.client.SendRawTransaction(r.Context(), "0x"+hex.EncodeToString(signed)); err != nil {
		handler.WriteError(w, http.StatusBadGateway, fmt.Sprintf("failed to send transaction: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewBuildERC20EIP1559TransactionResponse(unsigned, signed, txHash, sig))
}
