package util

import (
	"math/big"
	"strings"
)

// MultiplyUnit multiplies an integer unit value by its scale factor and returns
// the exact base-10 integer representation.
func MultiplyUnit(value *big.Int, multiplier *big.Int) string {
	if value == nil {
		return "0"
	}

	return new(big.Int).Mul(value, multiplier).String()
}

// FormatScaledInt formats an integer as a fixed-point decimal string by placing
// the decimal point `scale` digits from the right.
func FormatScaledInt(value *big.Int, scale int) string {
	if value == nil {
		return "0." + strings.Repeat("0", scale)
	}

	sign := ""
	abs := new(big.Int).Set(value)
	if abs.Sign() < 0 {
		sign = "-"
		abs.Neg(abs)
	}

	digits := abs.String()
	if scale == 0 {
		return sign + digits
	}

	if len(digits) <= scale {
		return sign + "0." + strings.Repeat("0", scale-len(digits)) + digits
	}

	split := len(digits) - scale
	return sign + digits[:split] + "." + digits[split:]
}
