package core

import (
	"bytes"
	"crypto/ecdsa"
	"errors"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EVMSecp256k1Key struct {
	PrivateKey *types.PrivateKey
	PublicKey  *types.PublicKey
	Address    *types.Address
}

// GenerateKey generates a new ECDSA private key for the specified curve.
func GenerateKey() (*EVMSecp256k1Key, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	pub := priv.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)

	return &EVMSecp256k1Key{
		PrivateKey: types.NewPrivateKey(priv),
		PublicKey:  types.NewPublicKey(pub),
		Address:    types.NewAddress(addr),
	}, nil
}

// DeriveKeyFromHex reconstructs an EVMSecp256k1Key from hex-encoded strings and verifies
// that the stored public key and address are consistent with the private key.
func DeriveKeyFromHex(privHex, pubHex, addrHex string) (*EVMSecp256k1Key, error) {
	priv, err := crypto.HexToECDSA(privHex)
	if err != nil {
		return nil, err
	}

	pub := priv.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)

	if types.NewPublicKey(pub).Hex() != pubHex {
		return nil, errors.New("stored public key does not match private key")
	}

	if addr != common.HexToAddress(addrHex) {
		return nil, errors.New("stored address does not match private key")
	}

	return &EVMSecp256k1Key{
		PrivateKey: types.NewPrivateKey(priv),
		PublicKey:  types.NewPublicKey(pub),
		Address:    types.NewAddress(addr),
	}, nil
}

func (k *EVMSecp256k1Key) VerifyKeyPair() error {
	if k == nil || k.PrivateKey.Key == nil || k.PublicKey.Key == nil {
		return errors.New("invalid key pair")
	}

	dpk, ok := k.PrivateKey.Key.Public().(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to derive public key")
	}

	got := crypto.FromECDSAPub(dpk)
	want := k.PublicKey.Bytes()
	if !bytes.Equal(got, want) {
		return errors.New("private key and public key do not match")
	}

	da := crypto.PubkeyToAddress(*dpk)
	if da != k.Address.Addr {
		return errors.New("public key and address do not match")
	}

	return nil
}
