package misc

import (
	"encoding/hex"
	"errors"
	"strings"
)

type SelectorRequest struct {
	Signature string `json:"signature" example:"transfer(address,uint256)"`
}

func (r *SelectorRequest) ValidateRequest() error {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

type SelectorResponse struct {
	Selector string `json:"selector"`
}

func NewSelectorResponse(selector []byte) *SelectorResponse {
	return &SelectorResponse{
		Selector: "0x" + hex.EncodeToString(selector),
	}
}

type EncodeRequest struct {
	Signature string   `json:"signature" example:"transfer(address,uint256)"`
	Args      []string `json:"args"      example:"[\"0xDa70aA79...\",\"1000000000000000000\"]"`
}

func (r *EncodeRequest) ValidateRequest() error {
	r.Signature = strings.TrimSpace(r.Signature)
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

type EncodeResponse struct {
	Data string `json:"data"`
}

func NewEncodeResponse(data []byte) *EncodeResponse {
	return &EncodeResponse{
		Data: "0x" + hex.EncodeToString(data),
	}
}
