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

type SignLegacyTransactionRequest struct {
	Address     string `json:"address"      example:"0xEbD69375..."`
	UnsignedRLP string `json:"unsigned_rlp" example:"0xec8085..."`

	unsignedRaw []byte
}

func (r *SignLegacyTransactionRequest) ValidateRequest() error {
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

func (r *SignLegacyTransactionRequest) UnsignedRaw() []byte {
	return r.unsignedRaw
}

type SignLegacyTransactionResponse struct {
	RawTransaction string `json:"raw_transaction"`
	TxHash         string `json:"tx_hash"`
}

func NewSignLegacyNativeTransferResponse(raw []byte, hash *types.Hash) *SignLegacyTransactionResponse {
	return &SignLegacyTransactionResponse{
		RawTransaction: "0x" + hex.EncodeToString(raw),
		TxHash:         hash.String(),
	}
}

type SignDynamicFeeTransactionRequest struct {
	Address     string `json:"address"      example:"0xEbD69375..."`
	UnsignedRLP string `json:"unsigned_rlp" example:"0x02f8..."`

	unsignedRaw []byte
}

func (r *SignDynamicFeeTransactionRequest) ValidateRequest() error {
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

func (r *SignDynamicFeeTransactionRequest) UnsignedRaw() []byte {
	return r.unsignedRaw
}

type SignDynamicFeeTransactionResponse struct {
	RawTransaction string `json:"raw_transaction"`
	TxHash         string `json:"tx_hash"`
}

func NewSignDynamicFeeTransactionResponse(raw []byte, hash *types.Hash) *SignDynamicFeeTransactionResponse {
	return &SignDynamicFeeTransactionResponse{
		RawTransaction: "0x" + hex.EncodeToString(raw),
		TxHash:         hash.String(),
	}
}

type BuildLegacyTransactionRequest struct {
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

func (r *BuildLegacyTransactionRequest) ValidateRequest() error {
	var ok bool

	r.chainID, ok = new(big.Int).SetString(strings.TrimSpace(r.ChainID), 10)
	if !ok {
		return errors.New("chain_id: must be a decimal integer")
	}
	if r.chainID.Sign() <= 0 {
		return errors.New("chain_id: must be positive")
	}

	r.gasPrice, ok = new(big.Int).SetString(strings.TrimSpace(r.GasPrice), 10)
	if !ok {
		return errors.New("gas_price: must be a decimal integer")
	}
	if r.gasPrice.Sign() <= 0 {
		return errors.New("gas_price: must be positive")
	}

	if r.GasLimit == 0 {
		return errors.New("gas_limit: must be greater than zero")
	}

	toBytes, err := util.ParseHex(r.To)
	if err != nil {
		return errors.New("to: " + err.Error())
	}
	r.to = common.BytesToAddress(toBytes)

	if v := strings.TrimSpace(r.Value); v == "" {
		r.value = new(big.Int)
	} else {
		r.value, ok = new(big.Int).SetString(v, 10)
		if !ok {
			return errors.New("value: must be a decimal integer")
		}
		if r.value.Sign() < 0 {
			return errors.New("value: must be non-negative")
		}
	}

	if d := strings.TrimSpace(r.Data); d == "" || d == "0x" {
		return nil
	}

	if r.data, err = util.ParseHex(r.Data); err != nil {
		return errors.New("data: " + err.Error())
	}

	return nil
}

func (r *BuildLegacyTransactionRequest) ToLegacyTx() *types.LegacyTx {
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

type BuildLegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SigningHash string `json:"signing_hash"`
}

func NewBuildLegacyNativeTransferResponse(raw []byte, hash *types.Hash) *BuildLegacyTransactionResponse {
	return &BuildLegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(raw),
		SigningHash: hash.String(),
	}
}

type BuildDynamicFeeTransactionRequest struct {
	ChainID              string `json:"chain_id"                 example:"20001209"`
	Nonce                uint64 `json:"nonce"                    example:"0"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas" example:"1500000000"`
	MaxFeePerGas         string `json:"max_fee_per_gas"          example:"3000000000"`
	GasLimit             uint64 `json:"gas_limit"                example:"21000"`
	To                   string `json:"to"                       example:"0x8336c196ABb9E7092C879C28D352b39d3f2f3D7A"`
	Value                string `json:"value"                    example:"1000000000000000000"`
	Data                 string `json:"data"                     example:"0x"`

	chainID   *big.Int
	gasTipCap *big.Int
	gasFeeCap *big.Int
	to        common.Address
	value     *big.Int
	data      []byte
}

func (r *BuildDynamicFeeTransactionRequest) ValidateRequest() error {
	var ok bool

	r.chainID, ok = new(big.Int).SetString(strings.TrimSpace(r.ChainID), 10)
	if !ok {
		return errors.New("chain_id: must be a decimal integer")
	}
	if r.chainID.Sign() <= 0 {
		return errors.New("chain_id: must be positive")
	}

	r.gasTipCap, ok = new(big.Int).SetString(strings.TrimSpace(r.MaxPriorityFeePerGas), 10)
	if !ok {
		return errors.New("max_priority_fee_per_gas: must be a decimal integer")
	}
	if r.gasTipCap.Sign() <= 0 {
		return errors.New("max_priority_fee_per_gas: must be positive")
	}

	r.gasFeeCap, ok = new(big.Int).SetString(strings.TrimSpace(r.MaxFeePerGas), 10)
	if !ok {
		return errors.New("max_fee_per_gas: must be a decimal integer")
	}
	if r.gasFeeCap.Sign() <= 0 {
		return errors.New("max_fee_per_gas: must be positive")
	}

	if r.GasLimit == 0 {
		return errors.New("gas_limit: must be greater than zero")
	}

	toBytes, err := util.ParseHex(r.To)
	if err != nil {
		return errors.New("to: " + err.Error())
	}
	r.to = common.BytesToAddress(toBytes)

	if v := strings.TrimSpace(r.Value); v == "" {
		r.value = new(big.Int)
	} else {
		r.value, ok = new(big.Int).SetString(v, 10)
		if !ok {
			return errors.New("value: must be a decimal integer")
		}
		if r.value.Sign() < 0 {
			return errors.New("value: must be non-negative")
		}
	}

	if d := strings.TrimSpace(r.Data); d == "" || d == "0x" {
		return nil
	}

	if r.data, err = util.ParseHex(r.Data); err != nil {
		return errors.New("data: " + err.Error())
	}

	return nil
}

func (r *BuildDynamicFeeTransactionRequest) ToDynamicFeeTx() *types.DynamicFeeTx {
	return &types.DynamicFeeTx{
		ChainID:   r.chainID,
		Nonce:     r.Nonce,
		GasTipCap: r.gasTipCap,
		GasFeeCap: r.gasFeeCap,
		GasLimit:  r.GasLimit,
		To:        &r.to,
		Value:     r.value,
		Data:      r.data,
	}
}

type BuildDynamicFeeTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
	SigningHash string `json:"signing_hash"`
}

func NewBuildDynamicFeeTransactionResponse(raw []byte, hash *types.Hash) *BuildDynamicFeeTransactionResponse {
	return &BuildDynamicFeeTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(raw),
		SigningHash: hash.String(),
	}
}
