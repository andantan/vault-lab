package v1

import (
	"errors"
	"strings"

	"github.com/andantan/evmlab/core/types"
)

type Keccak256LegacyRequest struct {
	Message string `json:"message" example:"hello world!"`
}

func (r *Keccak256LegacyRequest) ValidateRequest() error {
	r.Message = strings.TrimSpace(r.Message)
	if r.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

type Keccak256LegacyResponse struct {
	Digest string `json:"digest"`
}

func NewKeccak256LegacyResponse(h *types.Hash) *Keccak256LegacyResponse {
	return &Keccak256LegacyResponse{
		Digest: h.String(),
	}
}

type Keccak256PersonalRequest struct {
	Message string `json:"message" example:"hello world!"`
}

func (r *Keccak256PersonalRequest) ValidateRequest() error {
	r.Message = strings.TrimSpace(r.Message)
	if r.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

type Keccak256PersonalResponse struct {
	Digest string `json:"digest"`
}

func NewKeccak256PersonalResponse(h *types.Hash) *Keccak256PersonalResponse {
	return &Keccak256PersonalResponse{
		Digest: h.String(),
	}
}
