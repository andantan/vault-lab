package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
)

// SignatureLength 64 bytes ECDSA signature + 1 byte recovery id
const (
	SignatureLength    = 65
	SignatureHexLength = 130
)

type Signature struct {
	bytes []byte
	r     *big.Int
	s     *big.Int
	v     byte
	hex   string
}

func NewSignature(b []byte) (*Signature, error) {
	if len(b) != SignatureLength {
		return nil, fmt.Errorf("signature must be %d but got: %d", SignatureLength, len(b))
	}

	cp := make([]byte, SignatureLength)
	copy(cp, b)

	return &Signature{
		bytes: cp,
		r:     new(big.Int).SetBytes(cp[:32]),
		s:     new(big.Int).SetBytes(cp[32:64]),
		v:     cp[64],
	}, nil
}

func (s *Signature) IsNil() bool {
	if s == nil || s.bytes == nil {
		return true
	}

	if s.r == nil || s.s == nil {
		return true
	}

	return false
}

func (s *Signature) Bytes() []byte {
	return s.bytes[:]
}

func (s *Signature) Equal(o *Signature) bool {
	return bytes.Equal(s.bytes, o.bytes)
}

func (s *Signature) Hex() string {
	if s.hex == "" {
		s.hex = hex.EncodeToString(s.bytes)
	}

	return s.hex
}

func (s *Signature) R() *big.Int {
	if s.r == nil {
		s.r = new(big.Int).SetBytes(s.bytes[:32])
	}

	return s.r
}

func (s *Signature) S() *big.Int {
	if s.s == nil {
		s.s = new(big.Int).SetBytes(s.bytes[32:64])
	}

	return s.s
}

func (s *Signature) V() byte {
	return s.v
}

func (s *Signature) LegacyV() byte {
	return s.V() + 27
}
