package v2

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
)

type BuildNativeLegacyTransactionRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`

	value *big.Int
	from  *types.Address
	to    *types.Address
}

func (r *BuildNativeLegacyTransactionRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}
	if !common.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	r.from = types.NewAddress(common.HexToAddress(r.From))

	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}
	if !common.IsHexAddress(r.To) {
		return errors.New("to: invalid address")
	}
	r.to = types.NewAddress(common.HexToAddress(r.To))

	r.Value = strings.TrimSpace(r.Value)
	if r.Value == "" {
		return errors.New("value is required")
	}
	v, ok := new(big.Int).SetString(r.Value, 10)
	if !ok {
		return errors.New("value: must be a decimal wei amount")
	}
	r.value = v

	return nil
}

func (r *BuildNativeLegacyTransactionRequest) FromAddr() *types.Address { return r.from }
func (r *BuildNativeLegacyTransactionRequest) ToAddr() *types.Address   { return r.to }
func (r *BuildNativeLegacyTransactionRequest) Amount() *big.Int         { return r.value }

type BuildNativeLegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
}

func NewBuildNativeLegacyTransaction(data []byte) *BuildNativeLegacyTransactionResponse {
	return &BuildNativeLegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(data),
	}
}

type BuildNativeEIP1559TransactionRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`

	v *big.Int
	f *types.Address
	t *types.Address
}

func (r *BuildNativeEIP1559TransactionRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}
	if !common.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	r.f = types.NewAddress(common.HexToAddress(r.From))

	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}
	if !common.IsHexAddress(r.To) {
		return errors.New("to: invalid address")
	}
	r.t = types.NewAddress(common.HexToAddress(r.To))

	r.Value = strings.TrimSpace(r.Value)
	if r.Value == "" {
		return errors.New("value is required")
	}
	v, ok := new(big.Int).SetString(r.Value, 10)
	if !ok {
		return errors.New("value: must be a decimal wei amount")
	}
	r.v = v

	return nil
}

func (r *BuildNativeEIP1559TransactionRequest) FromAddr() *types.Address { return r.f }
func (r *BuildNativeEIP1559TransactionRequest) ToAddr() *types.Address   { return r.t }
func (r *BuildNativeEIP1559TransactionRequest) Amount() *big.Int         { return r.v }

type BuildNativeEIP1559TransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
}

func NewBuildNativeEIP1559Transaction(data []byte) *BuildNativeEIP1559TransactionResponse {
	return &BuildNativeEIP1559TransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(data),
	}
}

type BuildERC20LegacyTransactionRequest struct {
	From     string `json:"from"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	To       string `json:"to"       example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Contract string `json:"contract" example:"0x5FbDB2315678afecb367f032d93F642f64180aa3"`
	Amount   string `json:"amount"   example:"1000000000000000000"`

	f *types.Address
	t *types.Address
	c *types.Address
	a *big.Int
}

func (r *BuildERC20LegacyTransactionRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !common.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	r.f = types.NewAddress(common.HexToAddress(r.From))

	r.To = strings.TrimSpace(r.To)
	if !common.IsHexAddress(r.To) {
		return errors.New("to: invalid address")
	}
	r.t = types.NewAddress(common.HexToAddress(r.To))

	r.Contract = strings.TrimSpace(r.Contract)
	if !common.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.c = types.NewAddress(common.HexToAddress(r.Contract))

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	r.a = a

	return nil
}

func (r *BuildERC20LegacyTransactionRequest) FromAddr() *types.Address     { return r.f }
func (r *BuildERC20LegacyTransactionRequest) ToAddr() *types.Address       { return r.t }
func (r *BuildERC20LegacyTransactionRequest) ContractAddr() *types.Address { return r.c }
func (r *BuildERC20LegacyTransactionRequest) ToAmount() *big.Int           { return r.a }

type BuildERC20LegacyTransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
}

func NewBuildERC20LegacyTransaction(data []byte) *BuildERC20LegacyTransactionResponse {
	return &BuildERC20LegacyTransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(data),
	}
}

type BuildERC20EIP1559TransactionRequest struct {
	From     string `json:"from"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	To       string `json:"to"       example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Contract string `json:"contract" example:"0x5FbDB2315678afecb367f032d93F642f64180aa3"`
	Amount   string `json:"amount"   example:"1000000000000000000"`

	f *types.Address
	t *types.Address
	c *types.Address
	a *big.Int
}

func (r *BuildERC20EIP1559TransactionRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if !common.IsHexAddress(r.From) {
		return errors.New("from: invalid address")
	}
	r.f = types.NewAddress(common.HexToAddress(r.From))

	r.To = strings.TrimSpace(r.To)
	if !common.IsHexAddress(r.To) {
		return errors.New("to: invalid address")
	}
	r.t = types.NewAddress(common.HexToAddress(r.To))

	r.Contract = strings.TrimSpace(r.Contract)
	if !common.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.c = types.NewAddress(common.HexToAddress(r.Contract))

	r.Amount = strings.TrimSpace(r.Amount)
	a, ok := new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	r.a = a

	return nil
}

func (r *BuildERC20EIP1559TransactionRequest) FromAddr() *types.Address     { return r.f }
func (r *BuildERC20EIP1559TransactionRequest) ToAddr() *types.Address       { return r.t }
func (r *BuildERC20EIP1559TransactionRequest) ContractAddr() *types.Address { return r.c }
func (r *BuildERC20EIP1559TransactionRequest) ToAmount() *big.Int           { return r.a }

type BuildERC20EIP1559TransactionResponse struct {
	UnsignedRLP string `json:"unsigned_rlp"`
}

func NewBuildERC20EIP1559Transaction(data []byte) *BuildERC20EIP1559TransactionResponse {
	return &BuildERC20EIP1559TransactionResponse{
		UnsignedRLP: "0x" + hex.EncodeToString(data),
	}
}
