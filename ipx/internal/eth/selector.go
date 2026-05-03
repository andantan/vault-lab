package eth

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func Selector(signature string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(signature))

	sum := hash.Sum(nil)
	return "0x" + hex.EncodeToString(sum[:4])
}
