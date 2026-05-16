package types

var (
	BalanceOfSelector    = []byte{0x70, 0xa0, 0x82, 0x31} // balanceOf(address)
	ApproveSelector      = []byte{0x09, 0x5e, 0xa7, 0xb3} // approve(address,uint256)
	TransferSelector     = []byte{0xa9, 0x05, 0x9c, 0xbb} // transfer(address,uint256)
	AllowanceSelector    = []byte{0xdd, 0x62, 0xed, 0x3e} // allowance(address,address)
	TransferFromSelector = []byte{0x23, 0xb8, 0x72, 0xdd} // transferFrom(address,address,uint256)
	EIP712DomainSelector = []byte{0x84, 0xb0, 0x19, 0x6e} // eip712Domain()
	NoncesSelector       = []byte{0x7e, 0xce, 0xbe, 0x00} // nonces(address)
	NameSelector         = []byte{0x06, 0xfd, 0xde, 0x03} // name()
	VersionSelector      = []byte{0x54, 0xfd, 0x4d, 0x50} // version()
	SymbolSelector       = []byte{0x95, 0xd8, 0x9b, 0x41} // symbol()
	DecimalsSelector     = []byte{0x31, 0x3c, 0xe5, 0x67} // decimals()
	TotalSupplySelector  = []byte{0x18, 0x16, 0x0d, 0xdd} // totalSupply()
)

// Parameter represents a single ABI function parameter with its name and type.
type Parameter struct {
	Name string
	Type string
}

// Function represents a parsed ABI function signature.
// Types and Names are both empty when the function takes no parameters (e.g. "totalSupply()").
// Names is empty when the signature contained no parameter names (e.g. "approve(address,uint256)").
// Names is populated only when all parameters include a name (e.g. "approve(address spender,uint256 amount)").
type Function struct {
	Signature string
	Hash      *Hash
	Name      string
	Types     []string
	Names     []string
}

func NewFunction(signature string, hash *Hash, name string) *Function {
	return &Function{
		Signature: signature,
		Hash:      hash,
		Name:      name,
		Types:     []string{},
		Names:     []string{},
	}
}

// Selector returns the 4-byte ABI selector: the first 4 bytes of the function's keccak256 hash.
func (f *Function) Selector() []byte {
	return f.Hash.Bytes()[:4]
}

// NamedArgs maps parameter names to the corresponding argument values.
func (f *Function) NamedArgs(args []any) map[string]any {
	m := make(map[string]any, len(f.Names))
	for i, name := range f.Names {
		m[name] = args[i]
	}
	return m
}

// Parameters zips Types and Names into a []Parameter slice.
// Names may be empty (len 0) when the signature contained no parameter names.
func (f *Function) Parameters() []Parameter {
	params := make([]Parameter, len(f.Types))
	for i := range f.Types {
		p := Parameter{
			Type: f.Types[i],
		}

		if i < len(f.Names) {
			p.Name = f.Names[i]
		}
		params[i] = p
	}
	return params
}
