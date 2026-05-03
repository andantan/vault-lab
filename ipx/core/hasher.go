package core

import (
	"strconv"

	"github.com/andantan/evmlab/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EIP191Prefix is the personal sign prefix defined by EIP-191.
const EIP191Prefix = "\x19Ethereum Signed Message:\n"

type hasher struct{}

var Hasher = new(hasher)

// Hash computes the raw Keccak256 hash of the input data and returns it as a wrapped Hash.
// No prefix or transformation is applied — the digest is the direct Keccak256 output of m.
func (h *hasher) Hash(m []byte) *types.Hash {
	digest := crypto.Keccak256(m)
	hash := common.BytesToHash(digest)
	return types.NewHash(hash)
}

// PersonalHash applies the EIP-191 personal sign prefix and returns the Keccak256 hash.
// The prefix "\x19Ethereum Signed Message:\n" + len(m) is prepended before hashing,
// matching the digest produced by eth_sign / personal_sign in wallets such as MetaMask.
func (h *hasher) PersonalHash(m []byte) *types.Hash {
	msg := make([]byte, 0, len(EIP191Prefix)+len(strconv.Itoa(len(m)))+len(m))
	msg = append(msg, EIP191Prefix...)
	msg = strconv.AppendInt(msg, int64(len(m)), 10)
	msg = append(msg, m...)

	return h.Hash(msg)
}
