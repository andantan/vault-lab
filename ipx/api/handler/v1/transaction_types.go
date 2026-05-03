package v1

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
	"github.com/ethereum/go-ethereum/common"
)

type SignLegacyNativeTransferRequest struct {
	Address     string `json:"address"      example:"0xEbD69375..."`
	UnsignedRLP string `json:"unsigned_rlp" example:"0xec8085..."`

	unsignedRaw []byte
}

func (r *SignLegacyNativeTransferRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}

	b, err := util.ParseHex(r.UnsignedRLP)
	if err != nil {
		return errors.New("unsigned_rlp: " + err.Error())
	}
	if len(b) == 0 {
		return errors.New("unsigned_rlp: must not be empty")
	}
	r.unsignedRaw = b

	return nil
}

func (r *SignLegacyNativeTransferRequest) UnsignedRaw() []byte {
	return r.unsignedRaw
}

type SignLegacyNativeTransferResponse struct {
	RawTransaction string `json:"raw_transaction"`
	TxHash         string `json:"tx_hash"`
}

func NewSignLegacyNativeTransferResponse(raw []byte, hash *types.Hash) *SignLegacyNativeTransferResponse {
	return &SignLegacyNativeTransferResponse{
		RawTransaction: "0x" + hex.EncodeToString(raw),
		TxHash:         hash.String(),
	}
}

type BuildLegacyNativeTransferRequest struct {
	ChainID  string `json:"chain_id"  example:"20001209"`
	Nonce    uint64 `json:"nonce"     example:"0"`
	GasPrice string `json:"gas_price" example:"20000000000"`
	GasLimit uint64 `json:"gas_limit" example:"21000"`
	To       string `json:"to"        example:"0x8336c196ABb9E7092C879C28D352b39d3f2f3D7A"`
	Value    string `json:"value"     example:"1000000000000000000"`
	Data     string `json:"data"      example:"0x"`

	chainID  *big.Int
	gasPrice *big.Int
	to       common.Address
	value    *big.Int
	data     []byte
}

func (r *BuildLegacyNativeTransferRequest) ValidateRequest() error {
	var ok bool

	r.chainID, ok = new(big.Int).SetString(strings.TrimSpace(r.ChainID), 10)
	if !ok || r.chainID.Sign() <= 0 {
		return errors.New("chain_id: must be a positive decimal integer")
	}

	r.gasPrice, ok = new(big.Int).SetString(strings.TrimSpace(r.GasPrice), 10)
	if !ok || r.gasPrice.Sign() <= 0 {
		return errors.New("gas_price: must be a positive decimal integer")
	}

	if r.GasLimit == 0 {
		return errors.New("gas_limit: must be greater than zero")
	}

	toBytes, err := util.ParseHex(r.To)
	if err != nil {
		return errors.New("to: " + err.Error())
	}
	r.to = common.BytesToAddress(toBytes)

	r.value, ok = new(big.Int).SetString(strings.TrimSpace(r.Value), 10)
	if !ok || r.value.Sign() < 0 {
		return errors.New("value: must be a non-negative decimal integer")
	}

	if d := strings.TrimSpace(r.Data); d == "" || d == "0x" {
		return nil
	}

	if r.data, err = util.ParseHex(r.Data); err != nil {
		return errors.New("data: " + err.Error())
	}

	return nil
}

func (r *BuildLegacyNativeTransferRequest) ToLegacyTx() *types.LegacyTx {
	return &types.LegacyTx{
		ChainID:  r.chainID,
		Nonce:    r.Nonce,
		GasPrice: r.gasPrice,
		GasLimit: r.GasLimit,
		To:       r.to,
		Value:    r.value,
		Data:     r.data,
	}
}

type BuildLegacyNativeTransferResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SigningHash string `json:"signing_hash"`
}

func NewBuildLegacyNativeTransferResponse(raw []byte, hash *types.Hash) *BuildLegacyNativeTransferResponse {
	return &BuildLegacyNativeTransferResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(raw),
		SigningHash: hash.String(),
	}
}
