package config

import "github.com/google/uuid"

// MainNet provides mainnet network params
func MainNet() *Network {
	nid, _ := uuid.Parse("d6482c89-bc2d-5b81-aa9a-bdf1d7317dc3")

	return &Network{
		NetworkID: nid,
		Name: "karbowanec",
		Ticker: "KRB",

		PublicAddressBase58Prefix: 111,         // addresses start with "K"
		TxProofBase58Prefix:       0x369488,    // (0x369488), starts with "Proof..."
		ReserveProofBase58Prefix:  0xa74ad1d14, // (0xa74ad1d14), starts with "RsrvPrf..."
		KeysSignatureBase58Prefix: 0xa7f2119,   // (0xa7f2119), starts with "SigV1..."

		GenesisCoinbaseTxHex: "010a01ff0001fac484c69cd608029b2e4c0281c0b02e7c53291a94d1d0cbff8883f8024f5142ee494ffbbd0880712101f904925cc23f86f9f3565188862275dc556a9bdfb6aec22c5aca7f0177c45ba8",
		P2PMinimumVersion: P2PVersion4,
		P2PCurrentVersion: P2PVersion4,

		SeedNodes: []string{
			"localhost:32347",
			//"node.karbo.network:32347",
			//"seed1.karbowanec.com:32347",
			//"seed2.karbowanec.com:32347",
			//"seed.karbo.cloud:32347",
			//"seed.karbo.org:32347",
			//"seed.karbo.io:32347",
			//"185.86.78.40:32347",
			//"108.61.198.115:32347",
			//"45.32.232.11:32347",
			//"46.149.182.151:32347",
			//"144.91.94.65:32347",
		},
	}
}
