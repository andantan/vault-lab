package types

import (
	"math/big"

	"github.com/andantan/evmlab/internal/util"
)

var (
	weiPerGwei  = new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	weiPerEther = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

func WeiToGwei(wei *big.Int) string {
	return util.FormatScaledInt(wei, 9)
}

func WeiToEther(wei *big.Int) string {
	return util.FormatScaledInt(wei, 18)
}

func GweiToWei(gwei *big.Int) string {
	return util.MultiplyUnit(gwei, weiPerGwei)
}

func GweiToEther(gwei *big.Int) string {
	return util.FormatScaledInt(gwei, 9)
}

func EtherToWei(ether *big.Int) string {
	return util.MultiplyUnit(ether, weiPerEther)
}

func EtherToGwei(ether *big.Int) string {
	return util.MultiplyUnit(ether, weiPerGwei)
}
