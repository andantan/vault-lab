package eth

import (
	"fmt"
	"strings"
)

func DecodeAddress(wordHex string) (string, error) {
	wordHex = strings.TrimPrefix(wordHex, "0x")
	if len(wordHex) != 64 {
		return "", fmt.Errorf("invalid abi word length: %d", len(wordHex))
	}
	addr := wordHex[24:]
	return "0x" + addr, nil
}
