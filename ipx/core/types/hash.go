package types

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

const (
	HashLength    = 32
	HashHexLength = 64
)

type Hash struct {
	Hash common.Hash

	bytes []byte
	hex   string
}

func NewHash(h common.Hash) *Hash {
	return &Hash{
		Hash: h,
	}
}

func NewHashFromBytes(b []byte) (*Hash, error) {
	if len(b) != HashLength {
		return nil, fmt.Errorf("hash must be %d but got: %d", HashLength, len(b))
	}

	return &Hash{
		Hash: common.BytesToHash(b),
	}, nil
}

func (h *Hash) IsNil() bool {
	if h == nil {
		return true
	}

	return false
}

func (h *Hash) Bytes() []byte {
	if h.bytes == nil {
		h.bytes = h.Hash.Bytes()
	}

	return h.bytes
}

func (h *Hash) String() string {
	if h.hex == "" {
		h.hex = h.Hash.Hex()
	}

	return h.hex
}

func (h *Hash) Equal(o *Hash) bool {
	return bytes.Equal(h.Bytes(), o.Bytes())
}
