package cryptonote

// Address represents karbo address
type Address struct {
	base58 string
}

// NewAddress build cryptonote address
// TODO: Implement proper address parsing
func NewAddress(address string) (Address, error) {
	return Address{
		base58: address,
	}, nil
}

// Base58 provides string representation of the address
func (address *Address) Base58() string {
	return address.base58
}

func (address *Address) String() string {
	return address.base58
}
