package config

import "github.com/google/uuid"

// MainNetParams provides mainnet network params
func MainNetParams() *Params {
	nid, _ := uuid.Parse("d6482c89-bc2d-5b81-aa9a-bdf1d7317dc3")

	return &Params{
		NetworkID: nid,
		PublicAddressBase58Prefix: 111,         // addresses start with "K"
		TxProofBase58Prefix:       0x369488,    // (0x369488), starts with "Proof..."
		ReserveProofBase58Prefix:  0xa74ad1d14, // (0xa74ad1d14), starts with "RsrvPrf..."
		KeysSignatureBase58Prefix: 0xa7f2119,   // (0xa7f2119), starts with "SigV1..."

		P2PMinimumVersion: P2PVersion1,
		P2PCurrentVersion: P2PVersion4,

		SeedNodes: []string{
			"seed1.karbowanec.com:32347",
			"seed2.karbowanec.com:32347",
			"seed.karbo.cloud:32347",
			"seed.karbo.org:32347",
			"seed.karbo.io:32347",
			"185.86.78.40:32347",
			"108.61.198.115:32347",
			"45.32.232.11:32347",
			"46.149.182.151:32347",
			"144.91.94.65:32347",
		},
	}
}
