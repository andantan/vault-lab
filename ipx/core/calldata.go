package core

import (
	"math/big"

	"github.com/andantan/evmlab/core/types"
)

// BalanceOfCalldata builds calldata for balanceOf(address).
func BalanceOfCalldata(account *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.BalanceOfSelector)
	copy(data[16:36], account.Bytes())
	return data
}

// ApproveCalldata builds calldata for approve(address,uint256).
func ApproveCalldata(spender *types.Address, amount *big.Int) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.ApproveSelector)
	copy(data[16:36], spender.Bytes())
	amount.FillBytes(data[36:68])
	return data
}

// TransferCalldata builds calldata for transfer(address,uint256).
func TransferCalldata(to *types.Address, amount *big.Int) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.TransferSelector)
	copy(data[16:36], to.Bytes())
	amount.FillBytes(data[36:68])
	return data
}

// TransferFromCalldata builds calldata for transferFrom(address,address,uint256).
func TransferFromCalldata(from, to *types.Address, amount *big.Int) []byte {
	data := make([]byte, 100)
	copy(data[:4], types.TransferFromSelector)
	copy(data[16:36], from.Bytes())
	copy(data[48:68], to.Bytes())
	amount.FillBytes(data[68:100])
	return data
}

// EIP712DomainCalldata builds calldata for eip712Domain().
func EIP712DomainCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.EIP712DomainSelector)
	return data
}

// NameCalldata builds calldata for name().
func NameCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.NameSelector)
	return data
}

// VersionCalldata builds calldata for version().
func VersionCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.VersionSelector)
	return data
}

// SymbolCalldata builds calldata for symbol().
func SymbolCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.SymbolSelector)
	return data
}

// DecimalsCalldata builds calldata for decimals().
func DecimalsCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.DecimalsSelector)
	return data
}

// TotalSupplyCalldata builds calldata for totalSupply().
func TotalSupplyCalldata() []byte {
	data := make([]byte, 4)
	copy(data, types.TotalSupplySelector)
	return data
}

// NoncesCalldata builds calldata for nonces(address).
func NoncesCalldata(owner *types.Address) []byte {
	data := make([]byte, 36)
	copy(data[:4], types.NoncesSelector)
	copy(data[16:36], owner.Bytes())
	return data
}

// AllowanceCalldata builds calldata for allowance(address,address).
func AllowanceCalldata(owner, spender *types.Address) []byte {
	data := make([]byte, 68)
	copy(data[:4], types.AllowanceSelector)
	copy(data[16:36], owner.Bytes())
	copy(data[48:68], spender.Bytes())
	return data
}
