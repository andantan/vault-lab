package core

import (
	"errors"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type signer struct{}

var Signer = new(signer)

// Sign signs the given 32-byte Keccak-256 digest with the provided private key
// and returns a recoverable EVM signature in [R || S || V] format.
//
// The input hash must already be computed by the caller.
// This method does not hash the original message or payload.
//
// Internally, it uses secp256k1 ECDSA signing and returns a 65-byte signature
// containing:
//   - R: 32 bytes
//   - S: 32 bytes
//   - V: 1 byte recovery id
func (s *signer) Sign(hash *types.Hash, priv types.PrivateKey) (*types.Signature, error) {
	raw, err := crypto.Sign(hash.Bytes(), priv.Key)
	if err != nil {
		return nil, err
	}

	sig, err := types.NewSignature(raw)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// EcrecoverPubKey recovers the uncompressed secp256k1 public key from a 32-byte
// Keccak-256 hash and a recoverable EVM signature in [R || S || V] format.
//
// This mirrors the behavior of the Ethereum ecrecover precompile (EIP-2).
// The recovered public key corresponds to the private key that produced the
// signature for the given hash.
//
// The input signature must be 65 bytes in [R || S || V] format where:
//   - R: 32 bytes
//   - S: 32 bytes
//   - V: 1 byte recovery id (0 or 1)
//
// Returns the recovered public key in uncompressed form (65 bytes, 0x04 prefix).
func (s *signer) EcrecoverPubKey(hash *types.Hash, sig *types.Signature) (*types.PublicKey, error) {
	if hash.IsNil() || sig.IsNil() {
		return nil, errors.New("hash or signature is nil")
	}

	r, err := secp256k1.RecoverPubkey(hash.Bytes(), sig.Bytes())
	if err != nil {
		return nil, errors.New(err.Error())
	}

	pk, err := crypto.UnmarshalPubkey(r)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return types.NewPublicKey(pk), nil
}

// Verify checks whether the given recoverable signature was produced for the
// provided hash by the expected public key.
//
// It recovers the public key from the hash and signature, then compares the
// recovered key bytes with the given public key bytes.
//
// This verification requires the signature to be in recoverable EVM format
// ([R || S || V], 65 bytes), and the public key bytes to use the same encoding
// format as the recovered public key.
func (s *signer) Verify(hash *types.Hash, pubKey *types.PublicKey, sig *types.Signature) error {
	if hash.IsNil() || pubKey.IsNil() || sig.IsNil() {
		return errors.New("given hash or public key or signature is nil")
	}

	r, err := s.EcrecoverPubKey(hash, sig)
	if err != nil {
		return err
	}

	if !pubKey.Equal(r) {
		return errors.New("public key mismatch")
	}

	return nil
}

// VerifyByAddress checks whether the given recoverable signature was produced
// for the provided hash by the private key corresponding to the expected address.
//
// It recovers the public key from the hash and signature, derives the Ethereum
// address from the recovered key, and compares it against the provided address.
//
// This is the address-based equivalent of Verify, and matches the semantics of
// the Solidity ecrecover pattern commonly used in smart contract authentication.
//
// This verification requires the signature to be in recoverable EVM format
// ([R || S || V], 65 bytes).
func (s *signer) VerifyByAddress(hash *types.Hash, address *types.Address, sig *types.Signature) error {
	if hash.IsNil() || address.IsNil() || sig.IsNil() {
		return errors.New("given hash or address or signature is nil")
	}

	r, err := s.EcrecoverPubKey(hash, sig)
	if err != nil {
		return err
	}

	if !address.Equal(r.Address()) {
		return errors.New("address mismatch")
	}

	return nil
}
