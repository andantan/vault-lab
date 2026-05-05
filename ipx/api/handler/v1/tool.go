package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
)

type ToolHandler struct{}

func NewToolHandler() *ToolHandler {
	return &ToolHandler{}
}

// ChecksumEIP55 godoc
// @Summary      Convert address to EIP-55 checksum format
// @Description  Returns the EIP-55 mixed-case checksum encoding for the given address
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        body  body      ChecksumEIP55Request  true  "Address"
// @Success      200   {object}  ChecksumEIP55Response
// @Failure      400   {object}  map[string]string
// @Router       /evm/tool/address/checksum/eip55 [post]
func (h *ToolHandler) ChecksumEIP55(w http.ResponseWriter, r *http.Request) {
	req := new(ChecksumEIP55Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewChecksumEIP55Response(req.ToAddress()))
}

// DeriveKey godoc
// @Summary      Derive key set from private key
// @Description  Returns the public key and address derived from the given private key
// @Tags         tool
// @Accept       json
// @Produce      json
// @Param        body  body      DeriveKeyRequest   true  "Private key"
// @Success      200   {object}  DeriveKeyResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/tool/crypto/derive [post]
func (h *ToolHandler) DeriveKey(w http.ResponseWriter, r *http.Request) {
	req := new(DeriveKeyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	key, err := core.DeriveKeyFromPrivHex(req.PrivateKey)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid private_key: %s", err))
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewDeriveKeyResponse(key))
}
