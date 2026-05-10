package core

import (
	"math/big"

	"github.com/andantan/evmlab/core/types"
)

const (
	selectorSize = 4
	wordSize     = 32
	addressSize  = 20
)

// BalanceOfCalldata builds calldata for balanceOf(address).
func BalanceOfCalldata(account *types.Address) []byte {
	data := make([]byte, selectorSize+wordSize)
	copy(data[:selectorSize], types.BalanceOfSelector)
	copy(data[selectorSize+wordSize-addressSize:selectorSize+wordSize], account.Bytes())
	return data
}

// ApproveCalldata builds calldata for approve(address,uint256).
func ApproveCalldata(spender *types.Address, amount *big.Int) []byte {
	data := make([]byte, selectorSize+wordSize+wordSize)
	copy(data[:selectorSize], types.ApproveSelector)
	copy(data[selectorSize+wordSize-addressSize:selectorSize+wordSize], spender.Bytes())
	amount.FillBytes(data[selectorSize+wordSize : selectorSize+wordSize+wordSize])
	return data
}

// TransferCalldata builds calldata for transfer(address,uint256).
func TransferCalldata(to *types.Address, amount *big.Int) []byte {
	data := make([]byte, selectorSize+wordSize+wordSize)
	copy(data[:selectorSize], types.TransferSelector)
	copy(data[selectorSize+wordSize-addressSize:selectorSize+wordSize], to.Bytes())
	amount.FillBytes(data[selectorSize+wordSize : selectorSize+wordSize+wordSize])
	return data
}

// TransferFromCalldata builds calldata for transferFrom(address,address,uint256).
func TransferFromCalldata(from, to *types.Address, amount *big.Int) []byte {
	data := make([]byte, selectorSize+wordSize+wordSize+wordSize)
	copy(data[:selectorSize], types.TransferFromSelector)
	copy(data[selectorSize+wordSize-addressSize:selectorSize+wordSize], from.Bytes())
	copy(data[selectorSize+wordSize+wordSize-addressSize:selectorSize+wordSize+wordSize], to.Bytes())
	amount.FillBytes(data[selectorSize+wordSize+wordSize : selectorSize+wordSize+wordSize+wordSize])
	return data
}

// AllowanceCalldata builds calldata for allowance(address,address).
func AllowanceCalldata(owner, spender *types.Address) []byte {
	data := make([]byte, selectorSize+wordSize+wordSize)
	copy(data[:selectorSize], types.AllowanceSelector)
	copy(data[selectorSize+wordSize-addressSize:selectorSize+wordSize], owner.Bytes())
	copy(data[selectorSize+wordSize+wordSize-addressSize:selectorSize+wordSize+wordSize], spender.Bytes())
	return data
}
