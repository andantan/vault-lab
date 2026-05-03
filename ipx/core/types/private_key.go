package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	PrivateKeyLength    = 32
	PrivateKeyHexLength = 64
)

type PrivateKey struct {
	Key *ecdsa.PrivateKey

	bytes []byte
	hex   string
}

func NewPrivateKey(k *ecdsa.PrivateKey) *PrivateKey {
	return &PrivateKey{
		Key: k,
	}
}

func (k *PrivateKey) IsNil() bool {
	if k == nil || k.Key == nil {
		return true
	}

	return false
}

func (k *PrivateKey) Bytes() []byte {
	if k.bytes == nil {
		k.bytes = crypto.FromECDSA(k.Key)
	}

	return k.bytes
}

func (k *PrivateKey) Hex() string {
	if k.hex == "" {
		k.hex = hex.EncodeToString(k.Bytes())
	}

	return k.hex
}

func (k *PrivateKey) Equal(o *PrivateKey) bool {
	return bytes.Equal(k.Bytes(), o.Bytes())
}

func (k *PrivateKey) PublicKey() *PublicKey {
	return NewPublicKey(k.Key.Public().(*ecdsa.PublicKey))
}
