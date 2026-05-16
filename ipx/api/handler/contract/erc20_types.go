package contract

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type ERC20MetadataRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Block    string `json:"block"    example:"latest"`
}

func (r *ERC20MetadataRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

type ERC20MetadataResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply string `json:"total_supply"`
}

func NewERC20MetadataResponse(n, s, ts string, d uint8) *ERC20MetadataResponse {
	return &ERC20MetadataResponse{
		Name:        n,
		Symbol:      s,
		TotalSupply: ts,
		Decimals:    d,
	}
}

type ERC20BalanceRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Account  string `json:"account"  example:"0xAbcD1234..."`
	Block    string `json:"block"    example:"latest"`

	account *types.Address
}

func (r *ERC20BalanceRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Account = strings.TrimSpace(r.Account)
	var err error
	if r.account, err = types.NewAddressFromHex(r.Account); err != nil {
		return errors.New("account: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *ERC20BalanceRequest) ToAccount() *types.Address { return r.account }

type ERC20BalanceResponse struct {
	Balance string `json:"balance"`
}

func NewERC20BalanceResponse(s string) *ERC20BalanceResponse {
	return &ERC20BalanceResponse{
		Balance: s,
	}
}

type ERC20AllowanceRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Owner    string `json:"owner"    example:"0xAbcD1234..."`
	Spender  string `json:"spender"  example:"0xAbcD1234..."`
	Block    string `json:"block"    example:"latest"`

	owner   *types.Address
	spender *types.Address
}

func (r *ERC20AllowanceRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Owner = strings.TrimSpace(r.Owner)
	var err error
	if r.owner, err = types.NewAddressFromHex(r.Owner); err != nil {
		return errors.New("owner: invalid address")
	}
	r.Spender = strings.TrimSpace(r.Spender)
	if r.spender, err = types.NewAddressFromHex(r.Spender); err != nil {
		return errors.New("spender: invalid address")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *ERC20AllowanceRequest) ToOwner() *types.Address   { return r.owner }
func (r *ERC20AllowanceRequest) ToSpender() *types.Address { return r.spender }

type ERC20AllowanceResponse struct {
	Allowance string `json:"allowance"`
}

func NewERC20AllowanceResponse(s string) *ERC20AllowanceResponse {
	return &ERC20AllowanceResponse{
		Allowance: s,
	}
}

type ERC20ApproveRequest struct {
	Contract string `json:"contract" example:"0x8336c196..."`
	Spender  string `json:"spender"  example:"0xAbcD1234..."`
	Value    string `json:"value"    example:"1000000000000000000"`
	Block    string `json:"block"    example:"latest"`

	spender *types.Address
	value   *big.Int
}

func (r *ERC20ApproveRequest) ValidateRequest() error {
	r.Contract = strings.TrimSpace(r.Contract)
	if !util.IsHexAddress(r.Contract) {
		return errors.New("contract: invalid address")
	}
	r.Spender = strings.TrimSpace(r.Spender)
	var err error
	if r.spender, err = types.NewAddressFromHex(r.Spender); err != nil {
		return errors.New("spender: invalid address")
	}
	r.value = new(big.Int)
	if _, ok := r.value.SetString(strings.TrimSpace(r.Value), 10); !ok {
		return errors.New("value: invalid uint256")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	return nil
}

func (r *ERC20ApproveRequest) ToSpender() *types.Address { return r.spender }
func (r *ERC20ApproveRequest) ToValue() *big.Int         { return r.value }

type ERC20ApproveResponse struct {
	Approved bool `json:"approved"`
}

func NewERC20ApproveResponse(a bool) *ERC20ApproveResponse {
	return &ERC20ApproveResponse{
		Approved: a,
	}
}

type BalanceOfCalldataRequest struct {
	Account string `json:"account" example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`

	a *types.Address
}

func (r *BalanceOfCalldataRequest) ValidateRequest() error {
	r.Account = strings.TrimSpace(r.Account)
	a, err := types.NewAddressFromHex(r.Account)
	if err != nil {
		return errors.New("account: invalid address")
	}
	r.a = a
	return nil
}

func (r *BalanceOfCalldataRequest) ToAccount() *types.Address { return r.a }

type BalanceOfCalldataResponse struct {
	Data string `json:"data"`
}

func NewBalanceOfCalldataResponse(b []byte) *BalanceOfCalldataResponse {
	return &BalanceOfCalldataResponse{Data: "0x" + hex.EncodeToString(b)}
}

type ApproveCalldataRequest struct {
	Spender string `json:"spender" example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Amount  string `json:"amount"  example:"1000000000000000000"`

	s *types.Address
	a *big.Int
}

func (r *ApproveCalldataRequest) ValidateRequest() error {
	r.Spender = strings.TrimSpace(r.Spender)
	s, err := types.NewAddressFromHex(r.Spender)
	if err != nil {
		return errors.New("spender: invalid address")
	}
	r.s = s

	r.Amount = strings.TrimSpace(r.Amount)
	var ok bool
	r.a, ok = new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	return nil
}

func (r *ApproveCalldataRequest) ToSpender() *types.Address { return r.s }
func (r *ApproveCalldataRequest) ToAmount() *big.Int        { return r.a }

type ApproveCalldataResponse struct {
	Data string `json:"data"`
}

func NewApproveCalldataResponse(b []byte) *ApproveCalldataResponse {
	return &ApproveCalldataResponse{Data: "0x" + hex.EncodeToString(b)}
}

type TransferCalldataRequest struct {
	To     string `json:"to"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Amount string `json:"amount" example:"1000000000000000000"`

	t *types.Address
	a *big.Int
}

func (r *TransferCalldataRequest) ValidateRequest() error {
	r.To = strings.TrimSpace(r.To)
	t, err := types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}
	r.t = t

	r.Amount = strings.TrimSpace(r.Amount)
	var ok bool
	r.a, ok = new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	return nil
}

func (r *TransferCalldataRequest) ToAddress() *types.Address { return r.t }
func (r *TransferCalldataRequest) ToAmount() *big.Int        { return r.a }

type TransferCalldataResponse struct {
	Data string `json:"data"`
}

func NewTransferCalldataResponse(b []byte) *TransferCalldataResponse {
	return &TransferCalldataResponse{Data: "0x" + hex.EncodeToString(b)}
}

type AllowanceCalldataRequest struct {
	Owner   string `json:"owner"   example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Spender string `json:"spender" example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`

	o *types.Address
	s *types.Address
}

func (r *AllowanceCalldataRequest) ValidateRequest() error {
	r.Owner = strings.TrimSpace(r.Owner)
	o, err := types.NewAddressFromHex(r.Owner)
	if err != nil {
		return errors.New("owner: invalid address")
	}
	r.o = o

	r.Spender = strings.TrimSpace(r.Spender)
	s, err := types.NewAddressFromHex(r.Spender)
	if err != nil {
		return errors.New("spender: invalid address")
	}
	r.s = s
	return nil
}

func (r *AllowanceCalldataRequest) ToOwner() *types.Address   { return r.o }
func (r *AllowanceCalldataRequest) ToSpender() *types.Address { return r.s }

type AllowanceCalldataResponse struct {
	Data string `json:"data"`
}

func NewAllowanceCalldataResponse(b []byte) *AllowanceCalldataResponse {
	return &AllowanceCalldataResponse{Data: "0x" + hex.EncodeToString(b)}
}

type TransferFromCalldataRequest struct {
	From   string `json:"from"   example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	To     string `json:"to"     example:"0xDa70aA79f1a329719b9cf9d334b0a82b1d5269f3"`
	Amount string `json:"amount" example:"1000000000000000000"`

	f *types.Address
	t *types.Address
	a *big.Int
}

func (r *TransferFromCalldataRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	f, err := types.NewAddressFromHex(r.From)
	if err != nil {
		return errors.New("from: invalid address")
	}
	r.f = f

	r.To = strings.TrimSpace(r.To)
	t, err := types.NewAddressFromHex(r.To)
	if err != nil {
		return errors.New("to: invalid address")
	}
	r.t = t

	r.Amount = strings.TrimSpace(r.Amount)
	var ok bool
	r.a, ok = new(big.Int).SetString(r.Amount, 10)
	if !ok {
		return errors.New("amount: invalid integer")
	}
	return nil
}

func (r *TransferFromCalldataRequest) ToFrom() *types.Address { return r.f }
func (r *TransferFromCalldataRequest) ToTo() *types.Address   { return r.t }
func (r *TransferFromCalldataRequest) ToAmount() *big.Int     { return r.a }

type TransferFromCalldataResponse struct {
	Data string `json:"data"`
}

func NewTransferFromCalldataResponse(b []byte) *TransferFromCalldataResponse {
	return &TransferFromCalldataResponse{Data: "0x" + hex.EncodeToString(b)}
}
