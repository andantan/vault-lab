package v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/util"
)

type ChainIDResponse struct {
	ChainID    uint64 `json:"chain_id"`
	ChainIDHex string `json:"chain_id_hex"`
}

func NewChainIDResponse(chainID uint64, hex string) *ChainIDResponse {
	return &ChainIDResponse{
		ChainID:    chainID,
		ChainIDHex: hex,
	}
}

type GasPriceResponse struct {
	GasPrice    string `json:"gas_price"`
	GasPriceHex string `json:"gas_price_hex"`
	Wei         string `json:"wei"`
	Gwei        string `json:"gwei"`
	Ether       string `json:"ether"`
}

func NewGasPriceResponse(gasPrice *big.Int, hex string) *GasPriceResponse {
	return &GasPriceResponse{
		GasPrice:    gasPrice.String(),
		GasPriceHex: hex,
		Wei:         gasPrice.String(),
		Gwei:        types.WeiToGwei(gasPrice),
		Ether:       types.WeiToEther(gasPrice),
	}
}

type BlockNumberResponse struct {
	BlockNumber    uint64 `json:"block_number"`
	BlockNumberHex string `json:"block_number_hex"`
}

func NewBlockNumberResponse(blockNumber uint64, hex string) *BlockNumberResponse {
	return &BlockNumberResponse{
		BlockNumber:    blockNumber,
		BlockNumberHex: hex,
	}
}

type NonceRequest struct {
	Address string `json:"address" example:"0xEbD69375..."`
	Block   string `json:"block"   example:"pending"`
}

func (r *NonceRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	if r.Address == "" {
		return errors.New("address is required")
	}
	if r.Block == "" {
		r.Block = "pending"
	}
	return nil
}

type NonceResponse struct {
	Nonce    uint64 `json:"nonce"`
	NonceHex string `json:"nonce_hex"`
}

func NewNonceResponse(nonce uint64, hex string) *NonceResponse {
	return &NonceResponse{
		Nonce:    nonce,
		NonceHex: hex,
	}
}

type BalanceRequest struct {
	Address string `json:"address" example:"0xEbD69375..."`
	Block   string `json:"block"   example:"latest"`
}

func (r *BalanceRequest) ValidateRequest() error {
	r.Address = strings.TrimSpace(r.Address)
	r.Block = strings.TrimSpace(r.Block)
	if r.Address == "" {
		return errors.New("address is required")
	}
	if _, err := util.ParseHex(r.Address); err != nil {
		return errors.New("address: " + err.Error())
	}
	return nil
}

type BalanceResponse struct {
	Balance    string `json:"balance"`
	BalanceHex string `json:"balance_hex"`
	Wei        string `json:"wei"`
	Gwei       string `json:"gwei"`
	Ether      string `json:"ether"`
}

func NewBalanceResponse(wei *big.Int, hex string) *BalanceResponse {
	return &BalanceResponse{
		Balance:    wei.String(),
		BalanceHex: hex,
		Wei:        wei.String(),
		Gwei:       types.WeiToGwei(wei),
		Ether:      types.WeiToEther(wei),
	}
}

type TransactionRequest struct {
	TxHash string `json:"tx_hash" example:"0xabc123..."`
}

func (r *TransactionRequest) ValidateRequest() error {
	r.TxHash = strings.TrimSpace(r.TxHash)
	if r.TxHash == "" {
		return errors.New("tx_hash is required")
	}
	bare := strings.TrimPrefix(r.TxHash, "0x")
	if len(bare) != types.HashHexLength {
		return fmt.Errorf("tx_hash: must be %d hex chars (got %d)", types.HashHexLength, len(bare))
	}
	return nil
}

type TransactionReceiptRequest struct {
	TxHash string `json:"tx_hash" example:"0xabc123..."`
}

func (r *TransactionReceiptRequest) ValidateRequest() error {
	r.TxHash = strings.TrimSpace(r.TxHash)
	if r.TxHash == "" {
		return errors.New("tx_hash is required")
	}
	bare := strings.TrimPrefix(r.TxHash, "0x")
	if len(bare) != types.HashHexLength {
		return fmt.Errorf("tx_hash: must be %d hex chars (got %d)", types.HashHexLength, len(bare))
	}
	return nil
}

type EstimateGasRequest struct {
	From  string `json:"from"  example:"0xEbD69375..."`
	To    string `json:"to"    example:"0x8336c196..."`
	Value string `json:"value" example:"1000000000000000000"`
	Data  string `json:"data"  example:"0x"`
	Block string `json:"block" example:"latest"`

	p map[string]string
}

func (r *EstimateGasRequest) ValidateRequest() error {
	r.From = strings.TrimSpace(r.From)
	if r.From == "" {
		return errors.New("from is required")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	r.p = make(map[string]string, 4)
	r.p["from"] = r.From

	if r.To != "" {
		r.p["to"] = r.To
	}
	if r.Value != "" {
		r.p["value"] = r.Value
	}
	if r.Data != "" {
		r.p["data"] = r.Data
	}
	return nil
}

func (r *EstimateGasRequest) Params() map[string]string {
	return r.p
}

type EstimateGasResponse struct {
	GasLimit    uint64 `json:"gas_limit"`
	GasLimitHex string `json:"gas_limit_hex"`
	Wei         string `json:"wei"`
	Gwei        string `json:"gwei"`
	Ether       string `json:"ether"`
}

func NewEstimateGasResponse(gasLimit uint64, hex string) *EstimateGasResponse {
	wei := new(big.Int).SetUint64(gasLimit)
	return &EstimateGasResponse{
		GasLimit:    gasLimit,
		GasLimitHex: hex,
		Wei:         wei.String(),
		Gwei:        types.WeiToGwei(wei),
		Ether:       types.WeiToEther(wei),
	}
}

type CallRequest struct {
	From  string `json:"from"  example:"0xEbD69375..."`
	To    string `json:"to"    example:"0x8336c196..."`
	Data  string `json:"data"  example:"0x70a08231..."`
	Block string `json:"block" example:"latest"`

	p map[string]string
}

func (r *CallRequest) ValidateRequest() error {
	r.To = strings.TrimSpace(r.To)
	if r.To == "" {
		return errors.New("to is required")
	}
	r.Data = strings.TrimSpace(r.Data)
	if r.Data == "" {
		return errors.New("data is required")
	}
	if r.Block == "" {
		r.Block = "latest"
	}
	r.p = make(map[string]string, 3)
	r.p["to"] = r.To
	r.p["data"] = r.Data
	if r.From != "" {
		r.p["from"] = r.From
	}
	return nil
}

func (r *CallRequest) Params() map[string]string {
	return r.p
}
