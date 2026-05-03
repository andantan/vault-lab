package core

import (
	"fmt"
	"math/big"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type codec struct{}

var Codec = new(codec)

// legacyTxRLP is the RLP-serializable form for legacy transactions.
// Field order must match: [nonce, gasPrice, gasLimit, to, value, data, v, r, s].
type legacyTxRLP struct {
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	To       common.Address
	Value    *big.Int
	Data     []byte
	V        *big.Int
	R        *big.Int
	S        *big.Int
}

// EncodeLegacyUnsigned returns the RLP encoding for EIP-155 signing:
// [nonce, gasPrice, gasLimit, to, value, data, chainId, 0, 0]
func (c *codec) EncodeLegacyUnsigned(tx *types.LegacyTx) ([]byte, error) {
	enc, err := rlp.EncodeToBytes(&legacyTxRLP{
		Nonce:    tx.Nonce,
		GasPrice: tx.GasPrice,
		GasLimit: tx.GasLimit,
		To:       tx.To,
		Value:    tx.Value,
		Data:     tx.Data,
		V:        tx.ChainID,
		R:        big.NewInt(0),
		S:        big.NewInt(0),
	})
	if err != nil {
		return nil, err
	}

	return enc, nil
}

// EncodeLegacySigned returns the signed RLP with EIP-155 v = chainId * 2 + 35 + recoveryId.
func (c *codec) EncodeLegacySigned(tx *types.LegacyTx, sig *types.Signature) ([]byte, error) {
	v := new(big.Int).Mul(tx.ChainID, big.NewInt(2))
	v.Add(v, big.NewInt(35))
	v.Add(v, new(big.Int).SetUint64(uint64(sig.V())))

	raw, err := rlp.EncodeToBytes(&legacyTxRLP{
		Nonce:    tx.Nonce,
		GasPrice: tx.GasPrice,
		GasLimit: tx.GasLimit,
		To:       tx.To,
		Value:    tx.Value,
		Data:     tx.Data,
		V:        v,
		R:        sig.R(),
		S:        sig.S(),
	})
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// DecodeLegacyUnsigned decodes an unsigned EIP-155 RLP back into a LegacyTx.
// The V field of the unsigned encoding carries chainId, which is extracted here.
func (c *codec) DecodeLegacyUnsigned(raw []byte) (*types.LegacyTx, error) {
	var dec legacyTxRLP
	if err := rlp.DecodeBytes(raw, &dec); err != nil {
		return nil, err
	}

	return &types.LegacyTx{
		ChainID:  dec.V,
		Nonce:    dec.Nonce,
		GasPrice: dec.GasPrice,
		GasLimit: dec.GasLimit,
		To:       dec.To,
		Value:    dec.Value,
		Data:     dec.Data,
	}, nil
}

type accessTuple struct {
	Address     common.Address
	StorageKeys []common.Hash
}

// dynamicFeeTxRLPUnsigned is the RLP payload for EIP-1559 signing hash input.
// Field order: [chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gasLimit, to, value, data, accessList]
type dynamicFeeTxRLPUnsigned struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int
	GasFeeCap  *big.Int
	GasLimit   uint64
	To         *common.Address `rlp:"nil"`
	Value      *big.Int
	Data       []byte
	AccessList []accessTuple
}

// dynamicFeeTxRLPSigned is the RLP payload for a signed EIP-1559 transaction.
// Field order: [chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gasLimit, to, value, data, accessList, v, r, s]
type dynamicFeeTxRLPSigned struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int
	GasFeeCap  *big.Int
	GasLimit   uint64
	To         *common.Address `rlp:"nil"`
	Value      *big.Int
	Data       []byte
	AccessList []accessTuple
	V          *big.Int
	R          *big.Int
	S          *big.Int
}

// EncodeDynamicFeeUnsigned returns the EIP-1559 signing preimage: 0x02 || rlp([chainId, ...]).
// Pass the result to Hasher.Hash to get the digest for signing.
func (c *codec) EncodeDynamicFeeUnsigned(tx *types.DynamicFeeTx) ([]byte, error) {
	rlpBytes, err := rlp.EncodeToBytes(&dynamicFeeTxRLPUnsigned{
		ChainID:    tx.ChainID,
		Nonce:      tx.Nonce,
		GasTipCap:  tx.GasTipCap,
		GasFeeCap:  tx.GasFeeCap,
		GasLimit:   tx.GasLimit,
		To:         tx.To,
		Value:      tx.Value,
		Data:       tx.Data,
		AccessList: []accessTuple{},
	})
	if err != nil {
		return nil, err
	}

	return append([]byte{0x02}, rlpBytes...), nil
}

// EncodeDynamicFeeSigned returns the EIP-1559 wire format: 0x02 || rlp([chainId, ..., v, r, s]).
// v is the raw recovery id (0 or 1) — no EIP-155 adjustment.
func (c *codec) EncodeDynamicFeeSigned(tx *types.DynamicFeeTx, sig *types.Signature) ([]byte, error) {
	rlpBytes, err := rlp.EncodeToBytes(&dynamicFeeTxRLPSigned{
		ChainID:    tx.ChainID,
		Nonce:      tx.Nonce,
		GasTipCap:  tx.GasTipCap,
		GasFeeCap:  tx.GasFeeCap,
		GasLimit:   tx.GasLimit,
		To:         tx.To,
		Value:      tx.Value,
		Data:       tx.Data,
		AccessList: []accessTuple{},
		V:          new(big.Int).SetUint64(uint64(sig.V())),
		R:          sig.R(),
		S:          sig.S(),
	})
	if err != nil {
		return nil, err
	}

	return append([]byte{0x02}, rlpBytes...), nil
}

// DecodeDynamicFeeUnsigned decodes the EIP-1559 signing preimage (0x02 || rlp([chainId, ...])) back into a DynamicFeeTx.
func (c *codec) DecodeDynamicFeeUnsigned(raw []byte) (*types.DynamicFeeTx, error) {
	if len(raw) < 1 || raw[0] != 0x02 {
		return nil, fmt.Errorf("not an EIP-1559 transaction: missing 0x02 prefix")
	}

	var dec dynamicFeeTxRLPUnsigned
	if err := rlp.DecodeBytes(raw[1:], &dec); err != nil {
		return nil, err
	}

	return &types.DynamicFeeTx{
		ChainID:   dec.ChainID,
		Nonce:     dec.Nonce,
		GasTipCap: dec.GasTipCap,
		GasFeeCap: dec.GasFeeCap,
		GasLimit:  dec.GasLimit,
		To:        dec.To,
		Value:     dec.Value,
		Data:      dec.Data,
	}, nil
}
