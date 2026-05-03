package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
)

type HashHandler struct{}

func NewHashHandler() *HashHandler {
	return &HashHandler{}
}

// Keccak256Legacy godoc
// @Summary      Compute raw Keccak256 hash
// @Description  Computes the Keccak256 hash of the given message with no prefix applied (no EIP standard)
// @Tags         hash
// @Accept       json
// @Produce      json
// @Param        body  body      Keccak256LegacyRequest   true  "Message to hash"
// @Success      200   {object}  Keccak256LegacyResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/hash/keccak256/legacy [post]
func (h *HashHandler) Keccak256Legacy(w http.ResponseWriter, r *http.Request) {
	req := new(Keccak256LegacyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash := core.Hasher.Hash([]byte(req.Message))
	handler.WriteJSON(w, http.StatusOK, NewKeccak256LegacyResponse(hash))
}

// Keccak256Personal godoc
// @Summary      Compute Keccak256 hash with EIP-191 prefix
// @Description  Prepends the EIP-191 personal sign prefix ("\x19Ethereum Signed Message:\n" + length) to the message and returns the Keccak256 hash — matches the digest produced by eth_sign / personal_sign
// @Tags         hash
// @Accept       json
// @Produce      json
// @Param        body  body      Keccak256PersonalRequest   true  "Message to hash"
// @Success      200   {object}  Keccak256PersonalResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/v1/hash/keccak256/personal [post]
func (h *HashHandler) Keccak256Personal(w http.ResponseWriter, r *http.Request) {
	req := new(Keccak256PersonalRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash := core.Hasher.PersonalHash([]byte(req.Message))
	handler.WriteJSON(w, http.StatusOK, NewKeccak256PersonalResponse(hash))
}
