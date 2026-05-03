package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/internal/config"
)

type SignHandler struct {
	cfg *config.Config
}

func NewSignHandler(cfg *config.Config) *SignHandler {
	return &SignHandler{cfg: cfg}
}

// Sign godoc
// @Summary      Sign a pre-computed hash
// @Description  Signs a hex-encoded 32-byte digest with the given private key using raw secp256k1 ECDSA. Hashing is the caller's responsibility — use the hash endpoints to produce the digest before calling this.
// @Tags         sign
// @Accept       json
// @Produce      json
// @Param        body  body      SignRequest   true  "Address and digest"
// @Success      200   {object}  SignResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/sign [post]
func (h *SignHandler) Sign(w http.ResponseWriter, r *http.Request) {
	req := new(SignRequest)
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

	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to derive key: %s", err))
		return
	}

	sig, err := core.Signer.Sign(req.Hash(), *evmKey.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to sign: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewSignResponse(sig))
}

// VerifyByPublicKey godoc
// @Summary      Verify a signature against a public key
// @Description  Recovers the signer's public key from the signature via ecrecover and compares it against the provided public key
// @Tags         sign
// @Accept       json
// @Produce      json
// @Param        body  body      VerifyByPublicKeyRequest   true  "Hash, public key, and signature"
// @Success      200   {object}  VerifyByPublicKeyResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/sign/verify/by-public-key [post]
func (h *SignHandler) VerifyByPublicKey(w http.ResponseWriter, r *http.Request) {
	req := new(VerifyByPublicKeyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := core.Signer.Verify(req.ToHash(), req.ToPublicKey(), req.ToSignature()) == nil
	handler.WriteJSON(w, http.StatusOK, NewVerifyByPublicKeyResponse(result))
}

// VerifyByAddress godoc
// @Summary      Verify a signature against an address
// @Description  Recovers the signer's address from the signature via ecrecover and compares it against the provided address
// @Tags         sign
// @Accept       json
// @Produce      json
// @Param        body  body      VerifyByAddressRequest   true  "Hash, address, and signature"
// @Success      200   {object}  VerifyByAddressResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/sign/verify/by-address [post]
func (h *SignHandler) VerifyByAddress(w http.ResponseWriter, r *http.Request) {
	req := new(VerifyByAddressRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := core.Signer.VerifyByAddress(req.ToHash(), req.ToAddress(), req.ToSignature()) == nil
	handler.WriteJSON(w, http.StatusOK, NewVerifyByAddressResponse(result))
}
