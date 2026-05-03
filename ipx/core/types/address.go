package types

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

const (
	AddressLength    = 20
	AddressHexLength = 40
)

type Address struct {
	Addr common.Address

	bytes []byte
	hex   string
}

func NewAddress(a common.Address) *Address {
	return &Address{
		Addr: a,
	}
}

func NewAddressFromBytes(b []byte) (*Address, error) {
	if len(b) != AddressLength {
		return nil, fmt.Errorf("address must be %d bytes but got: %d", AddressLength, len(b))
	}

	return &Address{
		Addr: common.BytesToAddress(b),
	}, nil
}

func (a *Address) IsNil() bool {
	if a == nil {
		return true
	}

	return false
}

func (a *Address) Bytes() []byte {
	if a.bytes == nil {
		a.bytes = a.Addr.Bytes()
	}

	return a.bytes
}

func (a *Address) String() string {
	if a.hex == "" {
		a.hex = a.Addr.Hex()
	}

	return a.hex
}

func (a *Address) Equal(o *Address) bool {
	return bytes.Equal(a.Bytes(), o.Bytes())
}
