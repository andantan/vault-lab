package misc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
)

type AbiHandler struct{}

func NewAbiHandler() *AbiHandler {
	return &AbiHandler{}
}

// Selector godoc
// @Summary      Compute ABI function selector
// @Description  Returns the 4-byte selector for the given function signature
// @Tags         abi
// @Accept       json
// @Produce      json
// @Param        body  body      SelectorRequest   true  "Function signature"
// @Success      200   {object}  SelectorResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/abi/selector [post]
func (h *AbiHandler) Selector(w http.ResponseWriter, r *http.Request) {
	req := new(SelectorRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	sel := core.ABI.Selector(req.Signature)
	handler.WriteJSON(w, http.StatusOK, NewSelectorResponse(sel))
}

// Encode godoc
// @Summary      ABI-encode a function call
// @Description  Returns the ABI-encoded calldata (4-byte selector + packed arguments) for the given signature and args.
// @Description  Signature must be in canonical form — no spaces, no parameter names (e.g. "transfer(address,uint256)").
// @Description  Supported types: address | bool | string | bytes | bytes1–bytes32 | uint8/16/32/64/128/256 | int8/16/32/64/128/256
// @Description  Tuple, array, and slice types are not yet supported.
// @Tags         abi
// @Accept       json
// @Produce      json
// @Param        body  body      EncodeRequest   true  "Function signature and arguments"
// @Success      200   {object}  EncodeResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/abi/encode [post]
func (h *AbiHandler) Encode(w http.ResponseWriter, r *http.Request) {
	req := new(EncodeRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := core.ABI.EncodeCall(req.Signature, req.Args)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewEncodeResponse(data))
}
