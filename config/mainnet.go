package config

// MainNetParams provides mainnet network params
func MainNetParams() *Params {

	return &Params{
		PublicAddressBase58Prefix: 111,         // addresses start with "K"
		TxProofBase58Prefix:       0x369488,    // (0x369488), starts with "Proof..."
		ReserveProofBase58Prefix:  0xa74ad1d14, // (0xa74ad1d14), starts with "RsrvPrf..."
		KeysSignatureBase58Prefix: 0xa7f2119,   // (0xa7f2119), starts with "SigV1..."
	}
}
