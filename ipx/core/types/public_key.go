package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	PublicKeyLength    = 65
	PublicKeyHexLength = 130
)

type PublicKey struct {
	Key *ecdsa.PublicKey

	bytes []byte
	hex   string
}

func NewPublicKey(k *ecdsa.PublicKey) *PublicKey {
	return &PublicKey{
		Key: k,
	}
}

func NewPublicKeyFromBytes(b []byte) (*PublicKey, error) {
	if len(b) != PublicKeyLength {
		return nil, fmt.Errorf("public key must be %d bytes but got: %d", PublicKeyLength, len(b))
	}

	k, err := crypto.UnmarshalPubkey(b)
	if err != nil {
		return nil, err
	}

	cp := make([]byte, PublicKeyLength)
	copy(cp, b)

	return &PublicKey{
		Key:   k,
		bytes: cp,
	}, nil
}

func (k *PublicKey) IsNil() bool {
	if k == nil || k.Key == nil {
		return true
	}

	return false
}

func (k *PublicKey) Bytes() []byte {
	if k.bytes == nil {
		k.bytes = crypto.FromECDSAPub(k.Key)
	}

	return k.bytes
}

func (k *PublicKey) Hex() string {
	if k.hex == "" {
		k.hex = hex.EncodeToString(k.Bytes())
	}

	return k.hex
}

func (k *PublicKey) Equal(o *PublicKey) bool {
	return bytes.Equal(k.Bytes(), o.Bytes())
}

func (k *PublicKey) Address() *Address {
	return NewAddress(crypto.PubkeyToAddress(*k.Key))
}
