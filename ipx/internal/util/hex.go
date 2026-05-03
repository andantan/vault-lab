package util

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func HexToUint64(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
}

func HexToBigInt(s string) (*big.Int, error) {
	n := new(big.Int)
	if _, ok := n.SetString(strings.TrimPrefix(s, "0x"), 16); !ok {
		return nil, fmt.Errorf("invalid hex integer: %s", s)
	}
	return n, nil
}
