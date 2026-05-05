package core

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type abiCodec struct{}

var ABI = new(abiCodec)

// Selector returns the 4-byte ABI selector for the given canonical function signature.
func (c *abiCodec) Selector(signature string) []byte {
	return Hasher.Hash([]byte(signature)).Bytes()[:4]
}

// EncodeCall ABI-encodes a function call: 4-byte selector + packed arguments.
// Signature must be in canonical form (no spaces, no parameter names) e.g. "transfer(address,uint256)".
// TODO: tuple, array, slice types are not yet supported in signature parsing or argument conversion.
func (c *abiCodec) EncodeCall(signature string, args []string) ([]byte, error) {
	idx := strings.Index(signature, "(")
	if idx < 0 || !strings.HasSuffix(signature, ")") {
		return nil, fmt.Errorf("invalid signature: expected name(type,...)")
	}

	name := strings.TrimSpace(signature[:idx])
	if name == "" {
		return nil, fmt.Errorf("invalid signature: missing function name")
	}

	var typeStrs []string
	if inner := signature[idx+1 : len(signature)-1]; inner != "" {
		typeStrs = strings.Split(inner, ",")
		for i := range typeStrs {
			typeStrs[i] = strings.TrimSpace(typeStrs[i])
		}
	}

	if len(typeStrs) != len(args) {
		return nil, fmt.Errorf("signature has %d param(s) but got %d arg(s)", len(typeStrs), len(args))
	}

	abiArgs := make(abi.Arguments, len(typeStrs))
	for i, ts := range typeStrs {
		t, err := abi.NewType(ts, "", nil)
		if err != nil {
			return nil, fmt.Errorf("invalid type %q: %s", ts, err)
		}
		abiArgs[i] = abi.Argument{Type: t}
	}

	goArgs := make([]interface{}, len(args))
	for i, arg := range args {
		v, err := convertArg(abiArgs[i].Type, arg)
		if err != nil {
			return nil, fmt.Errorf("arg[%d] (%s): %s", i, typeStrs[i], err)
		}
		goArgs[i] = v
	}

	packed, err := abiArgs.Pack(goArgs...)
	if err != nil {
		return nil, fmt.Errorf("pack: %s", err)
	}

	canonical := name + "(" + strings.Join(typeStrs, ",") + ")"
	sel := c.Selector(canonical)
	return append(sel, packed...), nil
}

func convertArg(t abi.Type, arg string) (any, error) {
	arg = strings.TrimSpace(arg)
	switch t.T {
	case abi.AddressTy:
		if !common.IsHexAddress(arg) {
			return nil, fmt.Errorf("invalid address: %s", arg)
		}
		return common.HexToAddress(arg), nil

	case abi.UintTy:
		n := new(big.Int)
		if _, ok := n.SetString(arg, 0); !ok {
			return nil, fmt.Errorf("invalid uint: %s", arg)
		}

		switch t.Size {
		case 8:
			return uint8(n.Uint64()), nil
		case 16:
			return uint16(n.Uint64()), nil
		case 32:
			return uint32(n.Uint64()), nil
		case 64:
			return n.Uint64(), nil
		default:
			return n, nil
		}

	case abi.IntTy:
		n := new(big.Int)
		if _, ok := n.SetString(arg, 0); !ok {
			return nil, fmt.Errorf("invalid int: %s", arg)
		}

		switch t.Size {
		case 8:
			return int8(n.Int64()), nil
		case 16:
			return int16(n.Int64()), nil
		case 32:
			return int32(n.Int64()), nil
		case 64:
			return n.Int64(), nil
		default:
			return n, nil
		}

	case abi.BoolTy:
		switch strings.ToLower(arg) {
		case "true", "1":
			return true, nil
		case "false", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("invalid bool: %s", arg)
		}

	case abi.StringTy:
		return arg, nil

	case abi.BytesTy:
		b, err := hexutil.Decode(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid bytes: %s", err)
		}
		return b, nil

	case abi.FixedBytesTy:
		b, err := hexutil.Decode(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid bytes%d: %s", t.Size, err)
		}
		arrType := reflect.ArrayOf(t.Size, reflect.TypeOf(byte(0)))
		arr := reflect.New(arrType).Elem()
		for i := 0; i < len(b) && i < t.Size; i++ {
			arr.Index(i).Set(reflect.ValueOf(b[i]))
		}
		return arr.Interface(), nil

	default:
		return nil, fmt.Errorf("unsupported type: %v", t)
	}
}
