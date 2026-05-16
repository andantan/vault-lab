package misc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andantan/evmlab/api/handler"
	"github.com/andantan/evmlab/core"
)

type ABIHandler struct{}

func NewAbiHandler() *ABIHandler {
	return &ABIHandler{}
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
func (h *ABIHandler) Selector(w http.ResponseWriter, r *http.Request) {
	req := new(SelectorRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	fn, err := core.ABI.ParseFunctionSignature(req.Signature)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewSelectorResponse(fn.Selector()))
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
func (h *ABIHandler) Encode(w http.ResponseWriter, r *http.Request) {
	req := new(EncodeRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	if err := req.ValidateRequest(); err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	fn, err := core.ABI.ParseFunctionSignature(req.Signature)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := core.ABI.EncodeCall(fn, req.Args)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewEncodeResponse(data))
}

// DecodeResult godoc
// @Summary      Decode ABI-encoded return data
// @Description  Decodes the raw hex output of an eth_call into human-readable values given a list of ABI types
// @Tags         abi
// @Accept       json
// @Produce      json
// @Param        body  body      DecodeResultRequest   true  "ABI types and hex-encoded return data"
// @Success      200   {object}  DecodeResultResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/abi/decode/result [post]
func (h *ABIHandler) DecodeResult(w http.ResponseWriter, r *http.Request) {
	req := new(DecodeResultRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	data, err := req.ValidateRequest()
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	values, err := core.ABI.DecodeResult(req.Types, data)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewDecodeResultResponse(values))
}

// DecodeCall godoc
// @Summary      Decode ABI-encoded function calldata
// @Description  Decodes calldata (4-byte selector + args) into a name→value map using the given function signature. Parameter names are used as keys; falls back to "arg0", "arg1", … when the signature has no names.
// @Tags         abi
// @Accept       json
// @Produce      json
// @Param        body  body      DecodeCallRequest   true  "Function signature and hex-encoded calldata"
// @Success      200   {object}  DecodeCallResponse
// @Failure      400   {object}  map[string]string
// @Router       /evm/abi/decode/call [post]
func (h *ABIHandler) DecodeCall(w http.ResponseWriter, r *http.Request) {
	req := new(DecodeCallRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handler.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", err))
		return
	}
	data, err := req.ValidateRequest()
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	fn, err := core.ABI.ParseFunctionSignature(req.Signature)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	values, err := core.ABI.DecodeCall(fn, data)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.WriteJSON(w, http.StatusOK, NewDecodeCallResponse(data, values))
}

// EIP712DomainCalldata godoc
// @Summary      Build eip712Domain calldata
// @Description  Returns ABI-encoded calldata for eip712Domain()
// @Tags         abi
// @Produce      json
// @Success      200   {object}  EIP712DomainCalldataResponse
// @Router       /evm/abi/encode/eip712-domain [post]
func (h *ABIHandler) EIP712DomainCalldata(w http.ResponseWriter, r *http.Request) {
	handler.WriteJSON(w, http.StatusOK, NewEIP712DomainCalldataResponse(core.EIP712DomainCalldata()))
}
