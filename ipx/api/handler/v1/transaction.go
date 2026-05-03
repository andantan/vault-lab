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

// BuildLegacyNativeTransfer godoc
// @Summary      Build an unsigned legacy native transfer transaction
// @Description  Constructs an unsigned EIP-155 legacy transaction for native ETH transfer and returns the RLP encoding and signing hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      BuildLegacyNativeTransferRequest  true  "Transaction fields"
// @Success      200   {object}  BuildLegacyNativeTransferResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/transaction/legacy/build [post]
func (h *TransactionHandler) BuildLegacyNativeTransfer(w http.ResponseWriter, r *http.Request) {
	req := new(BuildLegacyNativeTransferRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := core.Codec.EncodeLegacyUnsigned(req.ToLegacyTx())
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode tx: %s", err))
		return
	}

	signingHash := core.Hasher.Hash(raw)
	handler.WriteJSON(w, http.StatusOK, NewBuildLegacyNativeTransferResponse(raw, signingHash))
}

// SignLegacyNativeTransfer godoc
// @Summary      Sign an unsigned legacy native transfer transaction
// @Description  Decodes an unsigned legacy RLP, signs it with the key for the given address, and returns the signed raw transaction and tx hash
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param        body  body      SignLegacyNativeTransferRequest  true  "Address and unsigned RLP"
// @Success      200   {object}  SignLegacyNativeTransferResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /evm/v1/transaction/legacy/sign [post]
func (h *TransactionHandler) SignLegacyNativeTransfer(w http.ResponseWriter, r *http.Request) {
	req := new(SignLegacyNativeTransferRequest)
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

	tx, err := core.Codec.DecodeLegacyUnsigned(req.UnsignedRaw())
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

	signedRaw, err := core.Codec.EncodeLegacySigned(tx, sig)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode signed tx: %s", err))
		return
	}

	txHash := core.Hasher.Hash(signedRaw)
	handler.WriteJSON(w, http.StatusOK, NewSignLegacyNativeTransferResponse(signedRaw, txHash))
}
