package v1

import (
	"errors"
	"strings"

	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type ChecksumEIP55Request struct {
	Address string `json:"address" example:"0xEbD69375..."`

	addr *types.Address
}

func (r *ChecksumEIP55Request) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}

	b, err := util.ParseHex(r.Address)
	if err != nil {
		return errors.New("address: " + err.Error())
	}

	r.addr, err = types.NewAddressFromBytes(b)
	if err != nil {
		return errors.New("address: " + err.Error())
	}

	return nil
}

func (r *ChecksumEIP55Request) ToAddress() *types.Address {
	return r.addr
}

type ChecksumEIP55Response struct {
	Address string `json:"address"`
}

func NewChecksumEIP55Response(addr *types.Address) *ChecksumEIP55Response {
	return &ChecksumEIP55Response{
		Address: addr.Checksum(),
	}
}

type DeriveKeyRequest struct {
	PrivateKey string `json:"private_key" example:"0xea66255f..."`
}

func (r *DeriveKeyRequest) ValidateRequest() error {
	r.PrivateKey = strings.TrimSpace(r.PrivateKey)
	if r.PrivateKey == "" {
		return errors.New("private_key is required")
	}
	return nil
}

type DeriveKeyResponse struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func NewDeriveKeyResponse(key *core.EVMSecp256k1Key) *DeriveKeyResponse {
	return &DeriveKeyResponse{
		Address:    key.Address.Checksum(),
		PublicKey:  key.PublicKey.Hex(),
		PrivateKey: key.PrivateKey.Hex(),
	}
}
