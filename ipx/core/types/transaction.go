package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// LegacyTx represents an unsigned EIP-155 legacy transaction.
type LegacyTx struct {
	ChainID  *big.Int
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	To       common.Address
	Value    *big.Int
	Data     []byte
}

// DynamicFeeTx represents an unsigned EIP-1559 transaction.
type DynamicFeeTx struct {
	ChainID   *big.Int
	Nonce     uint64
	GasTipCap *big.Int
	GasFeeCap *big.Int
	GasLimit  uint64
	To        *common.Address
	Value     *big.Int
	Data      []byte
}
