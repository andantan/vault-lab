package v1

import (
	"errors"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type SignRequest struct {
	Address string `json:"address" example:"0xEbD69375..."`
	Digest  string `json:"digest"  example:"0xa1b2c3d4..."`

	hash *types.Hash
}

func (r *SignRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}

	b, err := util.ParseHex(r.Digest)
	if err != nil {
		return errors.New("digest: " + err.Error())
	}
	r.hash, err = types.NewHashFromBytes(b)
	if err != nil {
		return errors.New("digest: " + err.Error())
	}

	return nil
}

func (r *SignRequest) Hash() *types.Hash {
	return r.hash
}

type SignResponse struct {
	Signature string `json:"signature"`
}

func NewSignResponse(s *types.Signature) *SignResponse {
	return &SignResponse{Signature: "0x" + s.Hex()}
}

type VerifyByPublicKeyRequest struct {
	Hash      string `json:"hash"       example:"0xa1b2c3d4..."`
	PublicKey string `json:"public_key" example:"0x04a1b2c3..."`
	Signature string `json:"signature"  example:"0xa1b2c3d4..."`

	hash   *types.Hash
	pubKey *types.PublicKey
	sig    *types.Signature
}

func (r *VerifyByPublicKeyRequest) ValidateRequest() error {
	b, err := util.ParseHex(r.Hash)
	if err != nil {
		return errors.New("hash: " + err.Error())
	}
	r.hash, err = types.NewHashFromBytes(b)
	if err != nil {
		return errors.New("hash: " + err.Error())
	}

	b, err = util.ParseHex(r.PublicKey)
	if err != nil {
		return errors.New("public_key: " + err.Error())
	}
	r.pubKey, err = types.NewPublicKeyFromBytes(b)
	if err != nil {
		return errors.New("public_key: " + err.Error())
	}

	b, err = util.ParseHex(r.Signature)
	if err != nil {
		return errors.New("signature: " + err.Error())
	}
	r.sig, err = types.NewSignature(b)
	if err != nil {
		return errors.New("signature: " + err.Error())
	}

	return nil
}

func (r *VerifyByPublicKeyRequest) ToHash() *types.Hash {
	return r.hash
}

func (r *VerifyByPublicKeyRequest) ToPublicKey() *types.PublicKey {
	return r.pubKey
}

func (r *VerifyByPublicKeyRequest) ToSignature() *types.Signature {
	return r.sig
}

type VerifyByAddressRequest struct {
	Hash      string `json:"hash"      example:"0xa1b2c3d4..."`
	Address   string `json:"address"   example:"0xAbCd1234..."`
	Signature string `json:"signature" example:"0xa1b2c3d4..."`

	hash *types.Hash
	addr *types.Address
	sig  *types.Signature
}

func (r *VerifyByAddressRequest) ValidateRequest() error {
	b, err := util.ParseHex(r.Hash)
	if err != nil {
		return errors.New("hash: " + err.Error())
	}
	r.hash, err = types.NewHashFromBytes(b)
	if err != nil {
		return errors.New("hash: " + err.Error())
	}

	b, err = util.ParseHex(r.Address)
	if err != nil {
		return errors.New("address: " + err.Error())
	}
	r.addr, err = types.NewAddressFromBytes(b)
	if err != nil {
		return errors.New("address: " + err.Error())
	}

	b, err = util.ParseHex(r.Signature)
	if err != nil {
		return errors.New("signature: " + err.Error())
	}
	r.sig, err = types.NewSignature(b)
	if err != nil {
		return errors.New("signature: " + err.Error())
	}

	return nil
}

func (r *VerifyByAddressRequest) ToHash() *types.Hash {
	return r.hash
}

func (r *VerifyByAddressRequest) ToAddress() *types.Address {
	return r.addr
}

func (r *VerifyByAddressRequest) ToSignature() *types.Signature {
	return r.sig
}

type VerifyByPublicKeyResponse struct {
	Result bool `json:"result"`
}

func NewVerifyByPublicKeyResponse(result bool) *VerifyByPublicKeyResponse {
	return &VerifyByPublicKeyResponse{Result: result}
}

type VerifyByAddressResponse struct {
	Result bool `json:"result"`
}

func NewVerifyByAddressResponse(result bool) *VerifyByAddressResponse {
	return &VerifyByAddressResponse{Result: result}
}
