package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/config"
)

type TransactionHandler struct {
	cfg *config.Config
}

func NewTransactionHandler(cfg *config.Config) *TransactionHandler {
	return &TransactionHandler{cfg: cfg}
}

// BuildLegacyTransaction godoc
// @Summary      Build an unsigned legacy native transfer transaction
// @Description  Constructs an unsigned EIP-155 legacy transaction for native ETH transfer and returns the RLP encoding and signing hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildLegacyTransactionRequest  true  "Transaction fields"
// @Success      200   {object}  BuildLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/transaction/legacy/build [post]
func (h *TransactionHandler) BuildLegacyTransaction(w http.ResponseWriter, r *http.Request) {
	req := new(BuildLegacyTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := core.RLP.EncodeLegacyUnsigned(req.ToLegacyTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(raw)
	handler.WriteJSON(w, http.StatusOK, NewBuildLegacyNativeTransferResponse(raw, signingHash))
}

// BuildDynamicFeeTransaction godoc
// @Summary      Build an unsigned dynamic fee native transfer transaction
// @Description  Constructs an unsigned EIP-1559 transaction for native ETH transfer and returns the encoded signing payload and signing hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildDynamicFeeTransactionRequest  true  "Transaction fields"
// @Success      200   {object}  BuildDynamicFeeTransactionResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/transaction/dynamic-fee/build [post]
func (h *TransactionHandler) BuildDynamicFeeTransaction(w http.ResponseWriter, r *http.Request) {
	req := new(BuildDynamicFeeTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := core.RLP.EncodeDynamicFeeUnsigned(req.ToDynamicFeeTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(raw)
	handler.WriteJSON(w, http.StatusOK, NewBuildDynamicFeeTransactionResponse(raw, signingHash))
}

// SignLegacyTransaction godoc
// @Summary      Sign an unsigned legacy native transfer transaction
// @Description  Decodes an unsigned legacy RLP, signs it with the key for the given address, and returns the signed raw transaction and tx hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      SignLegacyTransactionRequest  true  "Address and unsigned RLP"
// @Success      200   {object}  SignLegacyTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/v1/transaction/legacy/sign [post]
func (h *TransactionHandler) SignLegacyTransaction(w http.ResponseWriter, r *http.Request) {
	req := new(SignLegacyTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.Address)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("address: %s", err))
		return
	}

	priv, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	tx, err := core.RLP.DecodeLegacyUnsigned(req.UnsignedRaw())
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode unsigned_rlp: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(req.UnsignedRaw())

	sig, err := core.Signer.Sign(signingHash, *priv.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign: %s", err))
		return
	}

	signedRaw, err := core.RLP.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(signedRaw)
	handler.WriteJSON(w, http.StatusOK, NewSignLegacyNativeTransferResponse(signedRaw, txHash))
}

// SignDynamicFeeTransaction godoc
// @Summary      Sign an unsigned dynamic fee native transfer transaction
// @Description  Decodes an unsigned EIP-1559 payload, signs it with the key for the given address, and returns the signed raw transaction and tx hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      SignDynamicFeeTransactionRequest  true  "Address and unsigned RLP"
// @Success      200   {object}  SignDynamicFeeTransactionResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/v1/transaction/dynamic-fee/sign [post]
func (h *TransactionHandler) SignDynamicFeeTransaction(w http.ResponseWriter, r *http.Request) {
	req := new(SignDynamicFeeTransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := h.cfg.KeyByAddress(req.Address)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("address: %s", err))
		return
	}

	priv, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	tx, err := core.RLP.DecodeDynamicFeeUnsigned(req.UnsignedRaw())
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode unsigned_rlp: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(req.UnsignedRaw())

	sig, err := core.Signer.Sign(signingHash, *priv.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign: %s", err))
		return
	}

	signedRaw, err := core.RLP.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(signedRaw)
	handler.WriteJSON(w, http.StatusOK, NewSignDynamicFeeTransactionResponse(signedRaw, txHash))
}
