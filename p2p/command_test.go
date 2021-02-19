package p2p

import (
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

var encodedHandshakeReq = []byte{
	0x1, 0x11, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1,
	0x8, 0x9, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0xc, 0x14, 0xa, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0xa, 0x40, 0xd6,
	0x48, 0x2c, 0x89, 0xbc, 0x2d, 0x5b, 0x81, 0xaa, 0x9a,
	0xbd, 0xf1, 0xd7, 0x31, 0x7d, 0xc3, 0x7, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x8, 0x4, 0x7, 0x70,
	0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x5, 0x10, 0x3,
	00, 00, 00, 00, 00, 00, 0xa, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5,
	0x72, 0x14, 0x6, 0x60, 00, 00, 00, 00, 0x7,
	0x6d, 0x79, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x6, 0x5b,
	0x7e, 00, 00, 0xc, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0xc, 0x8,
	0xe, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x6, 0xe7, 0xcc,
	0x8, 00, 0x6, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64,
	0xa, 0x80, 0x8, 0xc, 0xe5, 0xe5, 0x96, 0x77, 0x9,
	0x4b, 0xde, 0xda, 0xae, 0xea, 0xe9, 0xa8, 0x96, 0x4b,
	0x60, 0xfe, 0xf7, 0x26, 0x70, 0x9c, 0xf6, 0x5f, 0x28,
	0x71, 0x38, 0x6e, 0xa2, 0x2c, 0x62, 0x9e,
}

var encodedHandshakeRes = []byte{
	0x1, 0x11, 0x1, 0x1, 0x1, 0x1, 0x2, 0x1, 0x1,
	0xc, 0x9, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0xc, 0x14, 0xa, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0xa, 0x40, 0xd6,
	0x48, 0x2c, 0x89, 0xbc, 0x2d, 0x5b, 0x81, 0xaa, 0x9a,
	0xbd, 0xf1, 0xd7, 0x31, 0x7d, 0xc3, 0x7, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x8, 0x1, 0x7, 0x70,
	0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x5, 0x20, 0x7f,
	0x48, 0x7d, 0xf, 0x87, 0x19, 0x30, 0xa, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5,
	0x98, 0x16, 0x6, 0x60, 00, 00, 00, 00, 0x7,
	0x6d, 0x79, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x6, 0x5b,
	0x7e, 00, 00, 0xc, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0xc, 0x8,
	0xe, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x6, 0xee, 0xcc,
	0x8, 00, 0x6, 0x74, 0x6f, 0x70, 0x5f, 0x69, 0x64,
	0xa, 0x80, 0x55, 0xaa, 0xa4, 0xdc, 0xad, 0xb7, 0xf,
	0x58, 0xb5, 0x48, 0x9d, 0x2f, 0x38, 0xaa, 0xf7, 0xd4,
	0x55, 0x61, 0x1, 0xf6, 0xe2, 0x4, 0x38, 0xb7, 0xed,
	0x8e, 0xa5, 0x60, 0xe8, 0xc9, 0xd1, 0x16, 0xe, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x70, 0x65, 0x65, 0x72,
	0x6c, 0x69, 0x73, 0x74, 0xa, 0x81, 0x5e, 0x5b, 0xf0,
	0x8f, 0x44, 0x5b, 0x7e, 00, 00, 0x10, 0x3, 00,
	00, 00, 00, 00, 00, 0x98, 0x16, 0x6, 0x60,
	00, 00, 00, 00, 0x6d, 0xce, 0x24, 0xb9, 0x5b,
	0x7e, 00, 00, 0x41, 0x2, 00, 00, 00, 00,
	00, 00, 0x98, 0x16, 0x6, 0x60, 00, 00, 00,
	00, 0x88, 0x25, 0xb9, 0xa3, 0x5b, 0x7e, 00, 00,
	0x3, 0x1e, 0xc4, 0xdc, 0xe6, 0x8, 0x6a, 0x5e, 0x97,
	0x16, 0x6, 0x60, 00, 00, 00, 00, 0xd5, 0x88,
	0x59, 0xfc, 0x5b, 0x7e, 00, 00, 0x12, 0x69, 0x6,
	0xf9, 0xa1, 0xa3, 0x48, 0x8e, 0x97, 0x16, 0x6, 0x60,
	00, 00, 00, 00, 0x2d, 0x50, 0x96, 0x21, 0x5b,
	0x7e, 00, 00, 0x80, 0x70, 0x4f, 0x21, 0x26, 0x81,
	0xcf, 0x65, 0x97, 0x16, 0x6, 0x60, 00, 00, 00,
	00, 0x74, 0xcb, 0x5a, 0x2c, 0x5b, 0x7e, 00, 00,
	0x10, 0x3, 00, 00, 00, 00, 00, 00, 0x97,
	0x16, 0x6, 0x60, 00, 00, 00, 00, 0x91, 0xef,
	0x1c, 0x37, 0x5b, 0x7e, 00, 00, 0x92, 0x3, 00,
	00, 00, 00, 00, 00, 0x97, 0x16, 0x6, 0x60,
	00, 00, 00, 00, 0x51, 0x19, 0xe7, 0x52, 0x5b,
	0x7e, 00, 00, 0xee, 0xbf, 0x52, 0x9a, 0x4c, 0xff,
	0x48, 0xf8, 0x94, 0x16, 0x6, 0x60, 00, 00, 00,
	00, 0x5b, 0xf6, 0x4, 0x64, 0x5b, 0x7e, 00, 00,
	0x5c, 0x71, 0xe3, 0x54, 0x79, 0xa0, 0x13, 0xc0, 0x93,
	0x16, 0x6, 0x60, 00, 00, 00, 00, 0x57, 0xf4,
	0x85, 0x9d, 0x5b, 0x7e, 00, 00, 0x86, 0x4f, 0xc9,
	0x69, 0x5, 0x53, 0xe, 0xfa, 0xe1, 0x11, 0x6, 0x60,
	00, 00, 00, 00, 0x5e, 0x9a, 0xd9, 0x6f, 0x5b,
	0x7e, 00, 00, 0xaf, 0x1, 00, 00, 00, 00,
	00, 00, 0xe4, 0xa6, 0x5, 0x60, 00, 00, 00,
	00, 0x56, 0x79, 0x63, 0xe9, 0x5b, 0x7e, 00, 00,
	0x31, 0x4b, 0xc7, 0xf6, 0x8f, 0x71, 0xc0, 0x4c, 0x49,
	0x8a, 0x5, 0x60, 00, 00, 00, 00, 0x4d, 0x7a,
	0xa6, 0x2a, 0x5b, 0x7e, 00, 00, 0xd7, 0xc5, 0x15,
	0xf8, 0x18, 0x6, 0x1b, 0x35, 0x5b, 0x43, 0x5, 0x60,
	00, 00, 00, 00, 0x1f, 0x14, 0x2, 0xf2, 0x5b,
	0x7e, 00, 00, 0x3f, 0x37, 0x68, 0x52, 0xd0, 0x68,
	0xff, 0x83, 0x97, 0xf7, 0x4, 0x60, 00, 00, 00,
	00, 0xb2, 0x5e, 0xf4, 0x4e, 0x5b, 0x7e, 00, 00,
	0xc7, 0x7c, 0x54, 0xad, 0x7f, 0x84, 0xa2, 0xcf, 0x83,
	0xce, 0x4, 0x60, 00, 00, 00, 00, 0x6d, 0xed,
	0xe, 0xc, 0x5b, 0x7e, 00, 00, 0xe8, 0x3, 00,
	00, 00, 00, 00, 00, 0x3e, 0xbd, 0x4, 0x60,
	00, 00, 00, 00, 0xb0, 0x75, 0xab, 0x28, 0x5b,
	0x7e, 00, 00, 0xdb, 0x66, 0xa4, 0xef, 0x46, 0x3b,
	0xda, 0x5f, 0x8c, 0xb9, 0x4, 0x60, 00, 00, 00,
	00, 0x4d, 0x78, 0xb0, 0xc8, 0x5b, 0x7e, 00, 00,
	0x37, 00, 00, 00, 00, 00, 00, 00, 0x7e,
	0x7b, 0x4, 0x60, 00, 00, 00, 00, 0xb2, 0xcc,
	0x2f, 0xb2, 0x5b, 0x7e, 00, 00, 0x3d, 0xde, 0x67,
	0x15, 0xbe, 0x77, 0xe4, 0x4b, 0x42, 0x5d, 0x4, 0x60,
	00, 00, 00, 00, 0xb0, 0x25, 0x32, 0xc7, 0x5b,
	0x7e, 00, 00, 0x37, 0xec, 0xc7, 0xfa, 0x2c, 0x90,
	0xba, 0x22, 0x35, 0x3c, 0x4, 0x60, 00, 00, 00,
	00, 0x6d, 0x57, 0x9b, 0x71, 0x5b, 0x7e, 00, 00,
	0x5b, 0xe8, 0x51, 0x4f, 0x2d, 0x6d, 0x98, 0x76, 0xf1,
	0x34, 0x4, 0x60, 00, 00, 00, 00, 0xb9, 0xcc,
	0x18, 0xab, 0x5b, 0x7e, 00, 00, 0xe8, 0x3, 00,
	00, 00, 00, 00, 00, 0x3b, 0x15, 0x4, 0x60,
	00, 00, 00, 00, 0xb2, 0x9f, 0xe3, 0x2c, 0x5b,
	0x7e, 00, 00, 0x68, 0x2, 00, 00, 00, 00,
	00, 00, 0x75, 0x39, 0x3, 0x60, 00, 00, 00,
	00, 0xd5, 0x6d, 0xe4, 0x11, 0x5b, 0x7e, 00, 00,
	0xe8, 0x68, 0x99, 0x58, 0x3f, 0xc3, 0x4c, 0x72, 0x4e,
	0x2a, 0x3, 0x60, 00, 00, 00, 00, 0x5d, 0x48,
	0x2f, 0xed, 0x5b, 0x7e, 00, 00, 0x52, 0xec, 0xf,
	0x4c, 0x4e, 0xc, 0xf9, 0x67, 0xd2, 0xf, 0x3, 0x60,
	00, 00, 00, 00, 0x80, 0x48, 0xad, 0xe6, 0x5b,
	0x7e, 00, 00, 0x79, 0x9d, 0xdb, 0xd2, 0xc6, 0xf3,
	0xbb, 0x86, 0x48, 0x2c, 0x2, 0x60, 00, 00, 00,
	00, 0x6d, 0xb, 0x82, 0xee, 0x5b, 0x7e, 00, 00,
	0x71, 0x81, 0x41, 0x69, 0x6, 0x49, 0x7, 0xf5, 0xe8,
	0xe4, 0x1, 0x60, 00, 00, 00, 00, 0x53, 0x19,
	0xda, 0xe9, 0x5b, 0x7e, 00, 00, 0xda, 0x27, 0x1b,
	0x12, 0x8c, 0x56, 0x81, 0x6c, 0x1b, 0xcf, 0x1, 0x60,
	00, 00, 00, 00, 0x58, 0xc6, 0x2f, 0x56, 0xdf,
	0x4e, 00, 00, 0xb2, 0xd1, 0x60, 0xf5, 0x48, 0x7,
	0x67, 0x14, 0xb2, 0xc8, 0x1, 0x60, 00, 00, 00,
	00, 0xbc, 0x24, 0xd3, 0x5e, 0x5b, 0x7e, 00, 00,
	0x52, 0xf, 0xaa, 0xd3, 0x70, 0x57, 0x80, 0x83, 0x8b,
	0xc8, 0x1, 0x60, 00, 00, 00, 00, 0x2e, 0x77,
	0x24, 0x21, 0x5b, 0x7e, 00, 00, 0x70, 0x50, 0x32,
	0xc4, 0x16, 0xd2, 0xe7, 0xaf, 0x9b, 0x82, 0x1, 0x60,
	00, 00, 00, 00, 0x53, 0x19, 0xdd, 0xae, 0x5b,
	0x7e, 00, 00, 0xda, 0x27, 0x1b, 0x12, 0x8c, 0x56,
	0x81, 0x6c, 0x24, 0xa, 0x1, 0x60, 00, 00, 00,
	00, 0x5d, 0x49, 0x16, 0xb6, 0x5b, 0x7e, 00, 00,
	0x3c, 0x59, 0x87, 0xf5, 0x5, 0x3e, 0x8f, 0xca, 0x8b,
	0xab, 00, 0x60, 00, 00, 00, 00, 0x25, 0x34,
	0x76, 0x81, 0x5b, 0x7e, 00, 00, 0x3, 0x2b, 0xfc,
	0xae, 0x1d, 0xed, 0x35, 0x26, 0xa6, 0x9e, 00, 0x60,
	00, 00, 00, 00, 0xb2, 0x89, 0x78, 0x9c, 0x5b,
	0x7e, 00, 00, 0x43, 0xf5, 0x75, 0x2c, 0x5a, 0x85,
	0x6b, 0x6e, 0xb8, 0x65, 00, 0x60, 00, 00, 00,
	00, 0xb9, 0x71, 0xd3, 0x3d, 0x5b, 0x7e, 00, 00,
	0xe8, 0x3, 00, 00, 00, 00, 00, 00, 0x25,
	0x5b, 00, 0x60, 00, 00, 00, 00, 0x6d, 0xed,
	0x7, 0x56, 0x5b, 0x7e, 00, 00, 0xe8, 0x3, 00,
	00, 00, 00, 00, 00, 0x4d, 0x51, 00, 0x60,
	00, 00, 00, 00, 0x52, 0x41, 0xa, 0xd7, 0x5b,
	0x7e, 00, 00, 0x5, 0x23, 0xd8, 0xf4, 0xd6, 0xeb,
	0xbb, 0xfc, 0xe7, 0x8, 00, 0x60, 00, 00, 00,
	00, 0x5, 0x3a, 0x45, 0xec, 0x5b, 0x7e, 00, 00,
	0x5f, 0x7d, 0x80, 0xae, 0xff, 0x6d, 0x79, 0x35, 0x90,
	0xc0, 0xff, 0x5f, 00, 00, 00, 00, 0x25, 0x73,
	0x9d, 0xb6, 0x5b, 0x7e, 00, 00, 0xd1, 0x8a, 0x18,
	0xf5, 0xce, 0xf9, 0x74, 0x5c, 0x34, 0x97, 0xff, 0x5f,
	00, 00, 00, 00, 0x55, 0xc2, 0xf1, 0x53, 0x5b,
	0x7e, 00, 00, 0xe8, 0x7d, 0xa1, 0x1, 0x80, 0xaa,
	0x25, 0xfa, 0xce, 0x93, 0xff, 0x5f, 00, 00, 00,
	00, 0x90, 0xca, 0x10, 0xd3, 0x5b, 0x7e, 00, 00,
	0xc4, 0x81, 0xca, 0x11, 0x6, 0x8f, 0xce, 0x7d, 0x20,
	0x8d, 0xff, 0x5f, 00, 00, 00, 00, 0x55, 0x9f,
	00, 0x2c, 0x5b, 0x7e, 00, 00, 0x10, 0x3, 00,
	00, 00, 00, 00, 00, 0xd4, 0x60, 0xff, 0x5f,
	00, 00, 00, 00, 0x5, 0x3a, 0x16, 0x64, 0x5b,
	0x7e, 00, 00, 00, 0x3d, 0x8f, 0x6b, 0x96, 0xef,
	0x27, 0x3b, 0x39, 0x57, 0xff, 0x5f, 00, 00, 00,
	00, 0xb0, 0x25, 0x46, 0xe3, 0x5b, 0x7e, 00, 00,
	0x83, 0xd5, 0x7, 0x97, 0x73, 0xff, 0x81, 0x5d, 0xf6,
	0x9, 0xff, 0x5f, 00, 00, 00, 00, 0x48, 0x5f,
	0x5d, 0x7c, 0x5b, 0x7e, 00, 00, 0x43, 0x41, 0xa7,
	0x16, 0xaa, 0xbf, 0x7e, 0x1e, 0x35, 0xfe, 0xfe, 0x5f,
	00, 00, 00, 00, 0x5b, 0x7b, 0x92, 0x8a, 0x5b,
	0x7e, 00, 00, 0x33, 0xcb, 0xcc, 0xd1, 0xa1, 0xd,
	0x9, 0x89, 0xc4, 0xeb, 0xfe, 0x5f, 00, 00, 00,
	00, 0x59, 0xa5, 0xe0, 0xc3, 0x5b, 0x7e, 00, 00,
	0x77, 0xf0, 0x8a, 0xd1, 0x49, 0x9f, 0x28, 0xd6, 0x71,
	0x9d, 0xfe, 0x5f, 00, 00, 00, 00, 0x56, 0x5,
	0xa6, 0x85, 0x5b, 0x7e, 00, 00, 0x78, 0x6d, 0x3f,
	0x29, 0x87, 0xea, 0x52, 0xbb, 0x62, 0x43, 0xfe, 0x5f,
	00, 00, 00, 00, 0xb5, 0x7a, 0xa2, 0x58, 0x5b,
	0x7e, 00, 00, 0x4a, 0xea, 0xa, 0x71, 0x5f, 0xc7,
	0xb6, 0x6d, 0xb7, 0x42, 0xfe, 0x5f, 00, 00, 00,
	00, 0x5d, 0x4c, 0x91, 0xb3, 0x5b, 0x7e, 00, 00,
	0x1d, 0xae, 0x6b, 0x2, 0xb0, 0x4e, 0x34, 0xe1, 0x98,
	0x1a, 0xfe, 0x5f, 00, 00, 00, 00, 0x5c, 0x34,
	0x98, 0x5, 0x5e, 0xd2, 00, 00, 0x28, 0x1, 0x3a,
	0xf1, 0xcc, 0x47, 0x7b, 0xb3, 0xd1, 0x80, 0xfd, 0x5f,
	00, 00, 00, 00, 0xb0, 0x64, 0xa3, 0xa1, 0x5b,
	0x7e, 00, 00, 0xd5, 0x68, 0xed, 0x91, 0x73, 0xd7,
	0xe5, 0xf4, 0xf5, 0x51, 0xfd, 0x5f, 00, 00, 00,
	00, 0x5c, 0x34, 0x98, 0x5, 0x5b, 0x7e, 00, 00,
	0x95, 0x56, 0xcc, 0xd, 0xe9, 0xd7, 0xd6, 0x9d, 0xcd,
	0x48, 0xfd, 0x5f, 00, 00, 00, 00, 0x4e, 0x38,
	0xfc, 0x23, 0x5b, 0x7e, 00, 00, 0xbf, 0x94, 0xa0,
	0xc2, 0x2f, 0xd4, 0x20, 0x72, 0x3b, 0x32, 0xfd, 0x5f,
	00, 00, 00, 00, 0x5f, 0x1b, 0xf5, 0x8e, 0x5b,
	0x7e, 00, 00, 0x85, 0xb2, 0x4a, 0x69, 0xa9, 0xbd,
	0x59, 0x9b, 0x1f, 0xe4, 0xfc, 0x5f, 00, 00, 00,
	00, 0xb0, 0x25, 0xc5, 0x3d, 0x5b, 0x7e, 00, 00,
	0x90, 0xaf, 0xd4, 0x38, 0xb0, 0xfb, 0x15, 0x5a, 0x68,
	0xbc, 0xfc, 0x5f, 00, 00, 00, 00, 0x6d, 0xfb,
	0x6, 0xf9, 0x5b, 0x7e, 00, 00, 0xbe, 0xd1, 0xe,
	0xd8, 0xfb, 0xff, 0x5c, 0x39, 0x59, 0x71, 0xfc, 0x5f,
	00, 00, 00, 00, 0xb0, 0xcd, 0xc5, 0xad, 0x5b,
	0x7e, 00, 00, 0x4d, 0xce, 0x96, 0x9b, 0xac, 0x62,
	0x29, 0x3a, 0x14, 0x6a, 0xfc, 0x5f, 00, 00, 00,
	00, 0x5b, 0xdd, 0x55, 0x19, 0x5b, 0x7e, 00, 00,
	0x8b, 0xd4, 0x43, 0xdf, 0x54, 0xdc, 0xd, 0xff, 0xc0,
	0x34, 0xfb, 0x5f, 00, 00, 00, 00, 0x4e, 0x89,
	0x1a, 0x71, 0x5b, 0x7e, 00, 00, 0x8d, 0xcb, 0xd9,
	0xb3, 0x81, 0x42, 0x1c, 0x80, 0x73, 0x8, 0xfb, 0x5f,
	00, 00, 00, 00, 0x17, 0xe7, 0x41, 0x38, 0x5b,
	0x7e, 00, 00, 0xac, 0xd4, 0x50, 0x6c, 0x76, 0xf6,
	0x8, 0xd, 0x55, 0xd6, 0xfa, 0x5f, 00, 00, 00,
	00, 0xc2, 0x2c, 0x90, 0xd5, 0x5b, 0x7e, 00, 00,
	0xdd, 0x77, 0xa0, 0x3b, 0x87, 0x28, 0xa2, 0x39, 0x49,
	0x40, 0xfa, 0x5f, 00, 00, 00, 00, 0x46, 0x34,
	0x53, 0xc5, 0x5b, 0x7e, 00, 00, 0x35, 0x1, 00,
	00, 00, 00, 00, 00, 0x75, 0x32, 0xfa, 0x5f,
	00, 00, 00, 00, 0x57, 0xc6, 0x6e, 0xa1, 0x5b,
	0x7e, 00, 00, 0x95, 0xca, 0x83, 0x82, 0xd5, 0x3b,
	0x17, 0x2a, 0xd5, 0xf6, 0xf9, 0x5f, 00, 00, 00,
	00, 0x53, 0x14, 0x85, 0x4c, 0x5b, 0x7e, 00, 00,
	0xda, 0x27, 0x1b, 0x12, 0x8c, 0x56, 0x81, 0x6c, 0xfe,
	0xd3, 0xf9, 0x5f, 00, 00, 00, 00, 0x5a, 0xf9,
	0x98, 0xec, 0x5b, 0x7e, 00, 00, 0xa4, 0x9a, 0x23,
	0x42, 0x1, 0x2f, 0x34, 0x53, 0x58, 0xb5, 0xf9, 0x5f,
	00, 00, 00, 00, 0x59, 0xbf, 0x7d, 0x9d, 0x5b,
	0x7e, 00, 00, 0xbf, 0x81, 0xad, 0xdc, 0x6b, 0xce,
	0xc6, 0xf9, 0x2, 0xb2, 0xf9, 0x5f, 00, 00, 00,
	00, 0x1f, 0x94, 0xf5, 0x39, 0x5b, 0x7e, 00, 00,
	0x6a, 0x8b, 0xd8, 0x59, 0xa3, 0x5e, 0x2, 0x4d, 0xda,
	0xa4, 0xf9, 0x5f, 00, 00, 00, 00, 0x5b, 0xed,
	0x7b, 0xcf, 0x5b, 0x7e, 00, 00, 0xbf, 0x97, 0xd9,
	0x35, 0xb0, 0x12, 0x6a, 0x50, 0xb7, 0x93, 0xf9, 0x5f,
	00, 00, 00, 00, 0x5d, 0x4b, 0xe6, 0x27, 0x5b,
	0x7e, 00, 00, 0xf7, 0xb3, 0x8b, 0xb, 0x7b, 0xc2,
	0xc0, 0x3f, 0xf4, 0xed, 0xf8, 0x5f, 00, 00, 00,
	00, 0xaf, 0xc4, 0xbb, 0xa3, 0x5b, 0x7e, 00, 00,
	0x28, 0x57, 0x90, 0x59, 0x3, 0x39, 0x19, 0x77, 0xd6,
	0xd6, 0xf8, 0x5f, 00, 00, 00, 00, 0x18, 0x14,
	0x11, 0xbc, 0x5b, 0x7e, 00, 00, 0xbc, 0xb, 0xc8,
	0x40, 0x78, 0x9d, 0xff, 0xbb, 0xa1, 0xd1, 0xf8, 0x5f,
	00, 00, 00, 00, 0x1f, 0x28, 0x6e, 0x3c, 0x5b,
	0x7e, 00, 00, 0x31, 0x2, 00, 00, 00, 00,
	00, 00, 0xa4, 0xb9, 0xf8, 0x5f, 00, 00, 00,
	00, 0x63, 0xe9, 0xfb, 0x17, 0x5b, 0x7e, 00, 00,
	0x2b, 00, 00, 00, 00, 00, 00, 00, 0xa4,
	0xb4, 0xf8, 0x5f, 00, 00, 00, 00, 0x5, 0xff,
	0xa1, 0x8d, 0x5b, 0x7e, 00, 00, 0x30, 0x85, 0x7d,
	0xf, 0xe, 0x83, 0x4e, 0x42, 0xab, 0x43, 0xf8, 0x5f,
	00, 00, 00, 00, 0x2e, 0x62, 0xc, 0x25, 0x5b,
	0x7e, 00, 00, 0xce, 0x66, 0xdb, 0xe2, 0xf4, 0x54,
	0xab, 0xd4, 0xa7, 00, 0xf8, 0x5f, 00, 00, 00,
	00, 0x6d, 0x68, 0xb0, 0x58, 0x5b, 0x7e, 00, 00,
	0x15, 0xa0, 0xad, 0x9c, 0x15, 0xb0, 0x82, 0x9f, 0x25,
	0x9e, 0xf7, 0x5f, 00, 00, 00, 00, 0xd4, 0x5c,
	0xf4, 0xe3, 0x5b, 0x7e, 00, 00, 0x61, 0x94, 0xaa,
	0x3c, 0x54, 0x95, 0x49, 0x85, 0x96, 0x65, 0xf7, 0x5f,
	00, 00, 00, 00, 0x5b, 0xd0, 0x41, 0x1a, 0x5b,
	0x7e, 00, 00, 0xb4, 0xb, 0x1, 0xa6, 0xed, 0x58,
	0x1b, 0x85, 0xba, 0x52, 0xf7, 0x5f, 00, 00, 00,
	00, 0x56, 0x79, 0x6c, 0x3d, 0x5b, 0x7e, 00, 00,
	0x31, 0x4b, 0xc7, 0xf6, 0x8f, 0x71, 0xc0, 0x4c, 0xcc,
	0xed, 0xf6, 0x5f, 00, 00, 00, 00, 0xd5, 0xe7,
	0x3b, 0x10, 0x5b, 0x7e, 00, 00, 0x1c, 0x47, 0xeb,
	0x1b, 0x98, 0xb5, 0x10, 0x8e, 0x81, 0xea, 0xf6, 0x5f,
	00, 00, 00, 00, 0x5c, 0x71, 0xb1, 0xaa, 0x5b,
	0x7e, 00, 00, 0xe8, 0x3, 00, 00, 00, 00,
	00, 00, 0x7d, 0x3b, 0xf6, 0x5f, 00, 00, 00,
	00, 0x5b, 0xf4, 0x39, 0x25, 0x5b, 0x7e, 00, 00,
	0x15, 0x88, 0xbe, 0xba, 0xf6, 0xca, 0x91, 0x5a, 0x1f,
	0x8e, 0xf5, 0x5f, 00, 00, 00, 00, 0xb9, 0xd1,
	0x39, 0x5c, 0x5b, 0x7e, 00, 00, 0xb8, 0x73, 0x5d,
	0xc3, 0x82, 00, 0x34, 0x6c, 0x62, 0x7b, 0xf5, 0x5f,
	00, 00, 00, 00, 0xb5, 0x3d, 0x36, 0xda, 0x5b,
	0x7e, 00, 00, 0x51, 0x2a, 0x70, 0x55, 0xaf, 0x9f,
	0x8c, 0x91, 0xc6, 0x35, 0xf5, 0x5f, 00, 00, 00,
	00, 0x5d, 0x48, 0x53, 0xf6, 0x5b, 0x7e, 00, 00,
	0x7, 0x3b, 0xaa, 0x8e, 0x20, 0x61, 0xd5, 0x81, 0xf9,
	0xee, 0xf4, 0x5f, 00, 00, 00, 00, 0xb2, 0x88,
	0xd3, 0xd3, 0x5b, 0x7e, 00, 00, 0x13, 0xa1, 0xf6,
	0xfc, 0x2, 0x9, 0x3f, 0xd4, 0x4e, 0xea, 0xf4, 0x5f,
	00, 00, 00, 00, 0xaf, 0x1c, 0x11, 0x13, 0x5b,
	0x7e, 00, 00, 0x5d, 0x4a, 0x88, 0x23, 0x7b, 0x19,
	0xa2, 0xbd, 0x5c, 0xf9, 0xf3, 0x5f, 00, 00, 00,
	00, 0x5, 0xf8, 0x30, 0xf, 0x5b, 0x7e, 00, 00,
	0x11, 0x68, 0x74, 0x3d, 0xf9, 0xbd, 0x74, 0x46, 0x51,
	0xe1, 0xf3, 0x5f, 00, 00, 00, 00, 0xc9, 0xef,
	0x52, 0xe6, 0x5b, 0x7e, 00, 00, 0x8b, 0x1c, 0x27,
	0x20, 0x9, 0x77, 0xfd, 0x97, 0x5, 0xca, 0xf3, 0x5f,
	00, 00, 00, 00, 0xe, 0xd7, 0x2d, 0xf3, 0x5b,
	0x7e, 00, 00, 0xe8, 0x3, 00, 00, 00, 00,
	00, 00, 0x27, 0x10, 0xf3, 0x5f, 00, 00, 00,
	00, 0xb2, 0x5c, 0x22, 0x7d, 0x5b, 0x7e, 00, 00,
	0xca, 0x85, 0xb4, 0x44, 0xc6, 0x92, 0x17, 0x25, 0x9e,
	00, 0xf3, 0x5f, 00, 00, 00, 00, 0x1f, 0xcf,
	0xa6, 0x20, 0x5b, 0x7e, 00, 00, 0x4c, 0xa4, 0xc6,
	0x39, 0x4a, 0x65, 0xe2, 0x8d, 0x3f, 0xe9, 0xf2, 0x5f,
	00, 00, 00, 00, 0x25, 0x37, 0x59, 0x5e, 0x5b,
	0x7e, 00, 00, 0x3, 0x2b, 0xfc, 0xae, 0x1d, 0xed,
	0x35, 0x26, 0x1a, 0xe0, 0xf2, 0x5f, 00, 00, 00,
	00, 0x5e, 0xf0, 0xb2, 0xcc, 0x5b, 0x7e, 00, 00,
	0x20, 0x37, 0x8c, 0x9b, 0x6e, 0x54, 0x19, 0xbe, 0xfa,
	0xd2, 0xf2, 0x5f, 00, 00, 00, 00, 0x5e, 0xf4,
	0xb4, 0xcf, 0x5b, 0x7e, 00, 00, 0x22, 0xd7, 0x56,
	0x8e, 0xb1, 0xd5, 0x97, 0x45, 0x79, 0x77, 0xf2, 0x5f,
	00, 00, 00, 00, 0xb2, 0xa5, 0x32, 0xd5, 0x5b,
	0x7e, 00, 00, 0x70, 0x19, 0x2, 0xcd, 0x6b, 0xf9,
	0x48, 0x32, 0x89, 0x33, 0xf2, 0x5f, 00, 00, 00,
	00, 0x65, 0xbc, 0x54, 0xe5, 0x5b, 0x7e, 00, 00,
	0x40, 0xba, 0x24, 0xc1, 0x9b, 0xc9, 0x25, 0x97, 0x69,
	0xfa, 0xf1, 0x5f, 00, 00, 00, 00, 0xb2, 0xd2,
	0xda, 0x94, 0x5b, 0x7e, 00, 00, 0xfb, 0xf4, 0xe3,
	0x74, 0x10, 0x6f, 0xb8, 0xd2, 0xba, 0xc9, 0xf1, 0x5f,
	00, 00, 00, 00, 0x5, 0x3a, 0xa8, 0xa7, 0x5b,
	0x7e, 00, 00, 0xc5, 0x9, 0xd, 0xac, 0xad, 0xe9,
	0xad, 0xfc, 0x14, 0x7c, 0xf1, 0x5f, 00, 00, 00,
	00, 0x90, 0xd9, 0x1e, 0x26, 0x5b, 0x7e, 00, 00,
	0x9a, 0xbd, 0x53, 0xbf, 0x3f, 0x4b, 0xb2, 0xf0, 0x1e,
	0x75, 0xf1, 0x5f, 00, 00, 00, 00, 0x5d, 0x48,
	0xf5, 0x94, 0x5b, 0x7e, 00, 00, 0xb0, 0x3d, 0xce,
	0xca, 0x2f, 0x9c, 0xc0, 0x70, 0xa1, 0x18, 0xf1, 0x5f,
	00, 00, 00, 00, 0x25, 0x39, 0x38, 0x85, 0x5b,
	0x7e, 00, 00, 0x6b, 0x7a, 0x2a, 0xdd, 0x27, 0x34,
	0x73, 0x98, 0xa5, 0xf4, 0xf0, 0x5f, 00, 00, 00,
	00, 0xb9, 0x8f, 0x92, 0x56, 0x5b, 0x7e, 00, 00,
	0xff, 0xc9, 0x7f, 0x9a, 0xd7, 0xa4, 0x63, 0xd4, 0x8a,
	0xe6, 0xf0, 0x5f, 00, 00, 00, 00, 0x4d, 0x79,
	0x74, 0xf4, 0x5b, 0x7e, 00, 00, 0x93, 0xda, 0xe2,
	0x36, 0x58, 0xd, 0x12, 0x2c, 0xa2, 0xce, 0xf0, 0x5f,
	00, 00, 00, 00, 0x3, 0x13, 0x1c, 0x5, 0x5b,
	0x7e, 00, 00, 0x56, 0x7a, 0xd1, 0x15, 0x28, 0x43,
	0x2d, 0x95, 0xba, 0xc5, 0xf0, 0x5f, 00, 00, 00,
	00, 0x3, 0x81, 0x3d, 0x4c, 0x5b, 0x7e, 00, 00,
	0x4, 0xaa, 0xdf, 0xda, 0xe0, 0xbe, 0x17, 0xb6, 0xa2,
	0xbe, 0xf0, 0x5f, 00, 00, 00, 00, 0x12, 0xde,
	0xe9, 0x76, 0x5b, 0x7e, 00, 00, 0x98, 0xf1, 0xeb,
	0x72, 0x1e, 0x1c, 0xb8, 0xaf, 0x76, 0xbe, 0xf0, 0x5f,
	00, 00, 00, 00, 0x12, 0xbc, 0x86, 0x4b, 0x5b,
	0x7e, 00, 00, 0xa6, 0x29, 0xea, 0x94, 0x8a, 0xb1,
	0x9a, 0x6f, 0x49, 0xbe, 0xf0, 0x5f, 00, 00, 00,
	00, 0x3, 0x15, 0x6a, 0xb7, 0x5b, 0x7e, 00, 00,
	0x67, 0xc1, 0x59, 0x86, 0x87, 0xce, 0x48, 0x33, 0xa5,
	0xbd, 0xf0, 0x5f, 00, 00, 00, 00, 0x3, 0x16,
	0x63, 0xfe, 0x5b, 0x7e, 00, 00, 0x25, 0xd0, 0xcb,
	0x96, 0xf4, 0x7d, 0x2e, 0x6d, 0x4f, 0xb9, 0xf0, 0x5f,
	00, 00, 00, 00, 0x3, 0xf, 0xd8, 0x2c, 0x5b,
	0x7e, 00, 00, 0x88, 0x9c, 0xff, 0x53, 0x70, 0xa0,
	0x4c, 0xe0, 0xb2, 0xb8, 0xf0, 0x5f, 00, 00, 00,
	00, 0x5, 0x3f, 0x9f, 0xbf, 0x5b, 0x7e, 00, 00,
	0x83, 0xc3, 0x22, 0x9f, 0x4c, 0x7e, 0x2f, 0xa7, 0x8d,
	0xb6, 0xf0, 0x5f, 00, 00, 00, 00, 0x5e, 0x2d,
	0x4a, 0x3d, 0x5b, 0x7e, 00, 00, 0xce, 0x24, 0x49,
	0xc4, 0x33, 0x1, 0xee, 0x64, 0x88, 0x94, 0xf0, 0x5f,
	00, 00, 00, 00, 0x48, 0x5c, 0x9f, 0xe7, 0x5b,
	0x7e, 00, 00, 0x99, 0x25, 0xd1, 0x9e, 0x14, 0x87,
	0xaa, 0xd1, 0x35, 0x71, 0xf0, 0x5f, 00, 00, 00,
	00, 0xb2, 0xcf, 0xc8, 0x6d, 0x5b, 0x7e, 00, 00,
	0x3d, 0xde, 0x67, 0x15, 0xbe, 0x77, 0xe4, 0x4b, 0x72,
	0x54, 0xf0, 0x5f, 00, 00, 00, 00, 0x5d, 0x49,
	0xd0, 0x4c, 0x5b, 0x7e, 00, 00, 0x9, 00, 00,
	00, 00, 00, 00, 00, 0xfd, 0x60, 0xef, 0x5f,
	00, 00, 00, 00, 0x5c, 0x71, 0xb0, 0xd7, 0x5b,
	0x7e, 00, 00, 0xe8, 0x3, 00, 00, 00, 00,
	00, 00, 0xe1, 0x40, 0xef, 0x5f, 00, 00, 00,
	00, 0x55, 0xf1, 0xa9, 0xea, 0x5b, 0x7e, 00, 00,
	0xc4, 0x3d, 0x1f, 0xf4, 0xcd, 0x81, 0x19, 0xed, 0x70,
	0x11, 0xef, 0x5f, 00, 00, 00, 00, 0x44, 0xc5,
	0x26, 0x8, 0x5b, 0x7e, 00, 00, 0x16, 0xdd, 0xc6,
	0xc1, 0x80, 0xe3, 0xee, 0x9e, 0xca, 0x9e, 0xee, 0x5f,
	00, 00, 00, 00, 0x5b, 0xe1, 0x5e, 0xc, 0x5b,
	0x7e, 00, 00, 0xdc, 0xd, 0xba, 0xa9, 0x48, 0x5d,
	0x9e, 0x14, 0x57, 0x24, 0xee, 0x5f, 00, 00, 00,
	00, 0x6d, 0xe5, 0xf, 0xb7, 0x5b, 0x7e, 00, 00,
	0x81, 0x3, 00, 00, 00, 00, 00, 00, 0xe7,
	0xee, 0xed, 0x5f, 00, 00, 00, 00, 0xb9, 0xb0,
	0x6c, 0x41, 0x5b, 0x7e, 00, 00, 0xae, 0xd8, 0x3d,
	0xd6, 0x1f, 0x7e, 0x10, 0xf6, 0x85, 0xc2, 0xed, 0x5f,
	00, 00, 00, 00, 0xb9, 0x67, 0x28, 0x6d, 0x5b,
	0x7e, 00, 00, 0x85, 0x8d, 0x94, 0x22, 0x66, 0x82,
	0xde, 0x5f, 0x77, 0x6e, 0xed, 0x5f, 00, 00, 00,
	00, 0x47, 0x13, 0xd9, 0xba, 0x5b, 0x7e, 00, 00,
	0xe5, 0x3c, 0x75, 0xf, 0x3, 0x91, 0x7b, 0x77, 0x22,
	0x3a, 0xed, 0x5f, 00, 00, 00, 00, 0xb0, 0x7d,
	0x39, 0x26, 0x5b, 0x7e, 00, 00, 0xe8, 0x3, 00,
	00, 00, 00, 00, 00, 0xfd, 0xbb, 0xec, 0x5f,
	00, 00, 00, 00, 0x56, 0x1, 0xa, 0xda, 0x5b,
	0x7e, 00, 00, 0xbd, 0x42, 0x7b, 0x59, 0x7b, 0x85,
	0x4c, 0x33, 0x77, 0x74, 0xec, 0x5f, 00, 00, 00,
	00, 0x2e, 0x62, 0x1e, 0xa2, 0x5b, 0x7e, 00, 00,
	0x5e, 0xea, 0x6d, 0xa5, 0x88, 0x35, 0xbb, 0x5c, 0xca,
	0x22, 0xec, 0x5f, 00, 00, 00, 00, 0xb0, 0x65,
	0xdb, 0xf6, 0x5b, 0x7e, 00, 00, 0x12, 0x2, 00,
	00, 00, 00, 00, 00, 0x9b, 0x1e, 0xec, 0x5f,
	00, 00, 00, 00, 0x2d, 0x76, 0x97, 0xf, 0x5b,
	0x7e, 00, 00, 0x73, 0x10, 0xd2, 0x96, 0xf6, 0x2c,
	0x79, 0xdf, 0xdf, 0x1b, 0xec, 0x5f, 00, 00, 00,
	00, 0x5f, 0xa4, 0x8, 0xcf, 0x5b, 0x7e, 00, 00,
	0x37, 00, 00, 00, 00, 00, 00, 00, 0x79,
	0x5c, 0xeb, 0x5f, 00, 00, 00, 00, 0x2e, 0x95,
	0xb6, 0x97, 0x33, 0x5b, 00, 00, 0xd0, 0xc5, 0xf1,
	0xe0, 0xae, 0x67, 0xe1, 0x40, 0x20, 0xb2, 0xe8, 0x5f,
	00, 00, 00, 00, 0xb9, 0xf7, 0x75, 0xd6, 0x5b,
	0x7e, 00, 00, 0xf, 0xf6, 0xb0, 0x23, 0x21, 0xd8,
	0x47, 0xeb, 0x53, 0xa5, 0xe8, 0x5f, 00, 00, 00,
	00, 0x5b, 0xf4, 0x15, 0x73, 0x5b, 0x7e, 00, 00,
	0x91, 0xb0, 0x87, 0xc1, 0xca, 0x96, 0x68, 0x2, 0x1c,
	0xa1, 0xe8, 0x5f, 00, 00, 00, 00, 0x58, 0x51,
	0xf8, 0x40, 0x5b, 0x7e, 00, 00, 0xc0, 0xf6, 0x6f,
	0x52, 0x31, 0xcf, 0xc0, 0x1c, 0x4a, 0x9e, 0xe8, 0x5f,
	00, 00, 00, 00, 0x4a, 0x49, 0xf8, 0x76, 0x5b,
	0x7e, 00, 00, 0xfa, 0x40, 0xf0, 0xb5, 0xb7, 0x77,
	0x99, 0xa3, 0xfe, 0x2e, 0xe8, 0x5f, 00, 00, 00,
	00, 0x2e, 0xa2, 0x13, 0x22, 0x5b, 0x7e, 00, 00,
	0xe8, 0x3, 00, 00, 00, 00, 00, 00, 0xe9,
	0x85, 0xe7, 0x5f, 00, 00, 00, 00, 0xb2, 0x97,
	0x35, 0x8c, 0x5b, 0x7e, 00, 00, 0x69, 0xf0, 0x37,
	0x1, 0x90, 0x2b, 0xfa, 0xf3, 0xfd, 0x6a, 0xe7, 0x5f,
	00, 00, 00, 00, 0x25, 0x19, 0x7e, 0xb8, 0x5b,
	0x7e, 00, 00, 0x1f, 0x10, 0xef, 0x94, 0xe1, 0xb2,
	0xc6, 0x98, 0x1d, 0xec, 0xe5, 0x5f, 00, 00, 00,
	00, 0xc2, 0x2c, 0x67, 0xe0, 0x5b, 0x7e, 00, 00,
	0x42, 00, 00, 00, 00, 00, 00, 00, 0x96,
	0x4e, 0xe5, 0x5f, 00, 00, 00, 00, 0x2e, 0xaf,
	0xf1, 0x10, 0x5b, 0x7e, 00, 00, 0x2f, 00, 00,
	00, 00, 00, 00, 00, 0xc8, 0x93, 0xe4, 0x5f,
	00, 00, 00, 00, 0xb0, 0x68, 0x34, 0xbd, 0x5b,
	0x7e, 00, 00, 0x24, 00, 00, 00, 00, 00,
	00, 00, 0x87, 0x96, 0xe3, 0x5f, 00, 00, 00,
	00, 0x54, 0x43, 0x2c, 0x15, 0x5b, 0x7e, 00, 00,
	0x33, 0x41, 0xb4, 0x6a, 0xbd, 0x3a, 0xf2, 0xe6, 0x95,
	0x57, 0xe3, 0x5f, 00, 00, 00, 00, 0x57, 0x42,
	0x5d, 0xb7, 0x5b, 0x7e, 00, 00, 0xf3, 0x84, 0xcf,
	0x8a, 0x16, 0x10, 0x28, 0x58, 0x37, 0xf8, 0xe2, 0x5f,
	00, 00, 00, 00, 0xb0, 0x78, 0x2f, 0xa3, 0x5b,
	0x7e, 00, 00, 0xb2, 0xde, 0x9a, 0x63, 0xfb, 0xff,
	0xea, 0x59, 0x8c, 0x21, 0xe2, 0x5f, 00, 00, 00,
	00, 0x56, 0x63, 0x74, 0x9f, 0x5b, 0x7e, 00, 00,
	0xdf, 0xff, 0xbc, 0x40, 0x98, 0x6a, 0xed, 0xb8, 0x1d,
	0xe4, 0xe1, 0x5f, 00, 00, 00, 00, 0xd9, 0xc7,
	0xef, 0x5a, 0x5b, 0x7e, 00, 00, 0x6b, 0x6, 0x89,
	0x4b, 0xfc, 0xd6, 0x5e, 0xae, 0xe7, 0xc1, 0xe0, 0x5f,
	00, 00, 00, 00, 0x6d, 0x56, 0x8a, 0xab, 0x5b,
	0x7e, 00, 00, 0xad, 0xa1, 0x76, 0x57, 0xf0, 0xc2,
	0xb7, 0xbd, 0x82, 0xb9, 0xe0, 0x5f, 00, 00, 00,
	00, 0xb0, 0xcd, 0x55, 0xe1, 0x5b, 0x7e, 00, 00,
	0x4d, 0xce, 0x96, 0x9b, 0xac, 0x62, 0x29, 0x3a, 0x23,
	0x4f, 0xe0, 0x5f, 00, 00, 00, 00, 0xc1, 0x6a,
	0xdc, 0x34, 0x5b, 0x7e, 00, 00, 0x5c, 0x76, 0xd1,
	0xdb, 0xb1, 0xe4, 0x2, 0x88, 0xd2, 0xb5, 0xdf, 0x5f,
	00, 00, 00, 00, 0x4e, 0x89, 0xd, 0x7f, 0x5b,
	0x7e, 00, 00, 0x34, 0x20, 0x19, 0x51, 0x94, 0xbe,
	0x64, 0x26, 0x59, 0x97, 0xdd, 0x5f, 00, 00, 00,
	00, 0xb0, 0x25, 0x57, 0xae, 0x5b, 0x7e, 00, 00,
	0xd0, 0x9a, 0xd5, 0xca, 0x26, 0xf6, 0x90, 0x13, 0xca,
	0x7, 0xdd, 0x5f, 00, 00, 00, 00, 0xb2, 0x96,
	0x7a, 0x99, 0x5b, 0x7e, 00, 00, 0x88, 0xc6, 0x3,
	0x89, 0x61, 0xc2, 0xa0, 0x8a, 0xfa, 0xa7, 0xdb, 0x5f,
	00, 00, 00, 00, 0x56, 0x7c, 0xfa, 0x84, 0x5b,
	0x7e, 00, 00, 0x31, 0x4b, 0xc7, 0xf6, 0x8f, 0x71,
	0xc0, 0x4c, 0x52, 0xc, 0xdb, 0x5f, 00, 00, 00,
	00, 0x4d, 0x79, 0x7b, 0xb, 0x5b, 0x7e, 00, 00,
	0x1d, 00, 00, 00, 00, 00, 00, 00, 0x4,
	0xb0, 0xda, 0x5f, 00, 00, 00, 00, 0xb0, 0x24,
	0x88, 0x28, 0x5b, 0x7e, 00, 00, 0xca, 0x33, 0xe1,
	0x5d, 0xd5, 0x2, 0x85, 0x3, 0xad, 0x5b, 0xda, 0x5f,
	00, 00, 00, 00, 0xb0, 0x25, 0x84, 0xc8, 0x5b,
	0x7e, 00, 00, 0x7c, 0x99, 0x9a, 0xeb, 0x73, 0xad,
	0xce, 0xd2, 0x15, 0x55, 0xda, 0x5f, 00, 00, 00,
	00, 0xb2, 0xd7, 0xbf, 0x2d, 0x5b, 0x7e, 00, 00,
	0xc8, 0x27, 0x74, 0xfa, 0xa6, 0xce, 0xb, 0x14, 0x46,
	0xe, 0xda, 0x5f, 00, 00, 00, 00, 0x2e, 0xaf,
	0x8e, 0xd8, 0x5b, 0x7e, 00, 00, 0xda, 0xc3, 0xed,
	0xa, 0x7, 0xa4, 0x88, 0x7c, 0x1a, 0xd, 0xda, 0x5f,
	00, 00, 00, 00, 0x65, 0xbf, 0x44, 0xb4, 0x5b,
	0x7e, 00, 00, 0xbd, 0x66, 0x92, 0xc0, 0xdf, 0xf0,
	0x17, 0xa0, 0x5b, 0x76, 0xd9, 0x5f, 00, 00, 00,
	00, 0x2e, 0x76, 0x94, 0x7b, 0x5b, 0x7e, 00, 00,
	0x9, 0x78, 0x5, 0x2c, 0x6c, 0x80, 0x4d, 0x2e, 0x9b,
	0xb7, 0xd8, 0x5f, 00, 00, 00, 00, 0x8e, 0xb8,
	0xbe, 0xc2, 0x5b, 0x7e, 00, 00, 0x19, 0x5d, 0x2f,
	0x1, 0x15, 0xa3, 0x5f, 0x2f, 0xa6, 0x33, 0xd8, 0x5f,
	00, 00, 00, 00, 0xc4, 0x32, 0xc6, 0xb0, 0x5b,
	0x7e, 00, 00, 0xa9, 0x1c, 0x50, 0xa5, 0xf1, 0x64,
	0xa3, 0x4b, 0x34, 0x1d, 0xd7, 0x5f, 00, 00, 00,
	00, 0x5, 0x3d, 0x38, 0xa, 0x5b, 0x7e, 00, 00,
	0x70, 0xaf, 0x27, 0x90, 0x86, 0xa9, 0xfc, 0x74, 0x60,
	0x97, 0xd6, 0x5f, 00, 00, 00, 00, 0x6d, 0xe3,
	0x41, 0x5, 0x5b, 0x7e, 00, 00, 0x79, 0xcb, 0x52,
	0x45, 0x43, 0x9b, 0xb9, 0xa1, 0x41, 0x8, 0xd6, 0x5f,
	00, 00, 00, 00, 0xc2, 0x2c, 0x81, 0x28, 0x5b,
	0x7e, 00, 00, 0xe8, 0x3, 00, 00, 00, 00,
	00, 00, 0xea, 0xe1, 0xd5, 0x5f, 00, 00, 00,
	00, 0xb0, 0x25, 0x7f, 0xf3, 0x5b, 0x7e, 00, 00,
	0x96, 0x2, 00, 00, 00, 00, 00, 00, 0xf9,
	0x70, 0xd5, 0x5f, 00, 00, 00, 00, 0x2e, 0x95,
	0xb6, 0x97, 0xc8, 0x2f, 00, 00, 0xb7, 0x37, 0x68,
	0x8, 0x6f, 0x7a, 0x4b, 0x68, 0xd8, 0xd6, 0xd3, 0x5f,
	00, 00, 00, 00, 0x31, 0xb0, 0xe9, 0x6c, 0x5b,
	0x7e, 00, 00, 0x5c, 0x74, 0x2, 0x85, 0xf1, 0xaf,
	0xc, 0x22, 0xb8, 0xbe, 0xd3, 0x5f, 00, 00, 00,
	00, 0x3e, 0xa6, 0x93, 0x26, 0x5b, 0x7e, 00, 00,
	0x45, 0xb8, 0xf2, 0x55, 0x84, 0x9d, 0x13, 0xd1, 0x4f,
	0xaa, 0xd3, 0x5f, 00, 00, 00, 00, 0xb0, 0x72,
	0xd0, 0xcf, 0x5b, 0x7e, 00, 00, 0x18, 0x32, 0x6d,
	0xea, 0xcd, 0xe7, 0xb7, 0x63, 0x9e, 0x9c, 0xd3, 0x5f,
	00, 00, 00, 00, 0x50, 0xd1, 0xe8, 0xef, 0x5b,
	0x7e, 00, 00, 0x89, 0x12, 0xd6, 0xb5, 0x7d, 0x78,
	0x18, 0xfc, 0xaa, 0x3f, 0xd3, 0x5f, 00, 00, 00,
	00, 0x5c, 0x5f, 0x85, 0x97, 0x5b, 0x7e, 00, 00,
	0xf8, 0xa7, 0xd2, 0xe5, 0x9b, 0x62, 0x1b, 0xc6, 0x77,
	0xca, 0xd2, 0x5f, 00, 00, 00, 00, 0x5d, 0x48,
	0x40, 0x88, 0x5b, 0x7e, 00, 00, 0xb9, 0x5, 0x30,
	0x6b, 0xcd, 0x8, 0x57, 0x96, 0xe6, 0x5a, 0xd2, 0x5f,
	00, 00, 00, 00, 0x4d, 0x78, 0xe6, 0xd4, 0x5b,
	0x7e, 00, 00, 0x4d, 0x9c, 0x25, 0x10, 0xea, 0xdd,
	0xfc, 0x37, 0xa3, 0x38, 0xd2, 0x5f, 00, 00, 00,
	00, 0xb0, 0x7d, 0x3c, 0x71, 0x5b, 0x7e, 00, 00,
	0x42, 00, 00, 00, 00, 00, 00, 00, 0xf6,
	0x2e, 0xd2, 0x5f, 00, 00, 00, 00, 0x5b, 0xcf,
	0x68, 0x8d, 0x5b, 0x7e, 00, 00, 0x5b, 00, 00,
	00, 00, 00, 00, 00, 0xbc, 0x28, 0xd2, 0x5f,
	00, 00, 00, 00, 0x5b, 0xf3, 0xcb, 0xc8, 0x5b,
	0x7e, 00, 00, 0x5e, 0xea, 0x6d, 0xa5, 0x88, 0x35,
	0xbb, 0x5c, 0x86, 0xf0, 0xd0, 0x5f, 00, 00, 00,
	00, 0x59, 0xb8, 0x42, 0x4c, 0x5b, 0x7e, 00, 00,
	0x4f, 0x1, 00, 00, 00, 00, 00, 00, 0x89,
	0x36, 0xcf, 0x5f, 00, 00, 00, 00, 0xb2, 0xd2,
	0xca, 0xf3, 0x5b, 0x7e, 00, 00, 0x7b, 0x83, 0x38,
	0xbd, 0xc5, 0x90, 0x50, 0xe2, 0xe, 0x6f, 0xce, 0x5f,
	00, 00, 00, 00, 0x6d, 0x56, 0xb8, 0x8e, 0x5b,
	0x7e, 00, 00, 0x7d, 0xa2, 0xdf, 0x63, 0x71, 0x67,
	0x6, 0xb8, 0xc0, 0x3a, 0xcd, 0x5f, 00, 00, 00,
	00, 0x6d, 0x68, 0xaa, 0xe9, 0x5b, 0x7e, 00, 00,
	0x55, 0x47, 0x6b, 0x4f, 0x8e, 0x6e, 0x5d, 0x69, 0x4e,
	0xe1, 0xcc, 0x5f, 00, 00, 00, 00, 0x4d, 0x78,
	0xae, 0xe2, 0x5b, 0x7e, 00, 00, 0x13, 0x94, 0x30,
	0xd6, 0x17, 0xfc, 0x55, 0x41, 0x59, 0xd2, 0xcc, 0x5f,
	00, 00, 00, 00, 0x4e, 0x1b, 0xa4, 0x8e, 0x5b,
	0x7e, 00, 00, 0xea, 0x21, 0xa2, 0x1e, 0x1c, 0x76,
	0x66, 0x89, 0xa1, 0xf0, 0xcb, 0x5f, 00, 00, 00,
	00, 0x6d, 0x10, 0x84, 0xb3, 0x5b, 0x7e, 00, 00,
	0x4a, 0xdd, 0x55, 0xb6, 0xc9, 0x67, 0xbe, 0x2b, 0x6d,
	0xcc, 0xcb, 0x5f, 00, 00, 00, 00, 0xb0, 0x25,
	0x84, 0xa4, 0x5b, 0x7e, 00, 00, 0xd, 00, 00,
	00, 00, 00, 00, 00, 0x6d, 0x9d, 0xcb, 0x5f,
	00, 00, 00, 00, 0x9f, 0xe0, 0x44, 0xd2, 0x5b,
	0x7e, 00, 00, 0x95, 0xf7, 0xde, 0x1e, 0x29, 0x9e,
	0x39, 0xcc, 0x29, 0x96, 0xcb, 0x5f, 00, 00, 00,
	00, 0x2d, 0x90, 0xc1, 0xe8, 0x5b, 0x7e, 00, 00,
	0x76, 0x3a, 0xb3, 0xf7, 0x74, 0x96, 0xac, 0x43, 0x32,
	0x30, 0xc9, 0x5f, 00, 00, 00, 00, 0x86, 0xf9,
	0x97, 0xd4, 0x5b, 0x7e, 00, 00, 0x5f, 0xed, 0xfa,
	0x26, 0x12, 0x36, 0x44, 0xa8, 0x19, 0x9b, 0xc8, 0x5f,
	00, 00, 00, 00, 0x51, 0xa7, 0x8, 0x64, 0x5b,
	0x7e, 00, 00, 0x97, 0xd3, 0xbc, 0xd9, 0x60, 0xa8,
	0x9, 0x92, 0x18, 0x6f, 0xc8, 0x5f, 00, 00, 00,
	00, 0x8e, 0x2f, 0x69, 0x63, 0x5b, 0x7e, 00, 00,
	0xac, 0x1e, 0x9a, 0x95, 0xc9, 0x50, 0x93, 0x69, 0xa5,
	0x40, 0xc8, 0x5f, 00, 00, 00, 00, 0x25, 0x36,
	0xe6, 0xa1, 0x5b, 0x7e, 00, 00, 0x3, 0x2b, 0xfc,
	0xae, 0x1d, 0xed, 0x35, 0x26, 0xc2, 0xf, 0xc8, 0x5f,
	00, 00, 00, 00, 0x5d, 0xaf, 0xc6, 0xbc, 0x5b,
	0x7e, 00, 00, 0x5a, 00, 00, 00, 00, 00,
	00, 00, 0xb0, 0x6f, 0xc6, 0x5f, 00, 00, 00,
	00, 0xd4, 0xb2, 0x2, 0xdb, 0x5b, 0x7e, 00, 00,
	0xeb, 0x7e, 0x3b, 0xf1, 0x69, 0x5a, 0x5a, 0x81, 0x52,
	0x1f, 0xc6, 0x5f, 00, 00, 00, 00, 0x86, 0xf9,
	0xf6, 0xb, 0x5b, 0x7e, 00, 00, 0x26, 0xa9, 0xe5,
	0xdd, 0xb4, 0x7e, 0x9a, 0xf0, 0x1a, 0xe4, 0xc5, 0x5f,
	00, 00, 00, 00, 0xb0, 0x62, 0x1c, 0xf4, 0x5b,
	0x7e, 00, 00, 0xee, 0xfe, 0x98, 0x9e, 0x2f, 0x7f,
	0xa6, 0x6b, 0xfb, 0xd5, 0xc4, 0x5f, 00, 00, 00,
	00, 0xb0, 0x25, 0x4b, 0xc5, 0x5b, 0x7e, 00, 00,
	0x27, 0x80, 0xcf, 0x5e, 0x48, 0xfe, 0x5d, 0xe8, 0x7a,
	0xe2, 0xc3, 0x5f, 00, 00, 00, 00, 0x5e, 0xf0,
	0xaa, 0x21, 0x5b, 0x7e, 00, 00, 0x20, 0x37, 0x8c,
	0x9b, 0x6e, 0x54, 0x19, 0xbe, 0x8, 0xaa, 0xc3, 0x5f,
	00, 00, 00, 00, 0x4e, 0x89, 0x16, 0x32, 0x5b,
	0x7e, 00, 00, 0x85, 0x8d, 0x94, 0x22, 0x66, 0x82,
	0xde, 0x5f, 0xa3, 0x61, 0xc3, 0x5f, 00, 00, 00,
	00, 0x5a, 0x5a, 0x45, 0x72, 0x5b, 0x7e, 00, 00,
	0x80, 0xae, 0x15, 0x38, 0x20, 0x52, 0x4f, 0xec, 0x2b,
	0xda, 0xc2, 0x5f, 00, 00, 00, 00, 0x2e, 0x76,
	0x86, 0x62, 0x5b, 0x7e, 00, 00, 0x46, 0xd9, 0x44,
	0x1d, 0x41, 0xc3, 0x6c, 0x3f, 0x5b, 0xcb, 0xc2, 0x5f,
	00, 00, 00, 00, 0xb0, 0x26, 0x2e, 0x18, 0x5b,
	0x7e, 00, 00, 0x35, 0xa, 0xb4, 0x8b, 0x66, 0xb5,
	0xed, 0xa5, 0x7a, 0x30, 0xc1, 0x5f, 00, 00, 00,
	00, 0x8d, 0x8a, 0x6e, 0x88, 0x5b, 0x7e, 00, 00,
	0x8, 00, 00, 00, 00, 00, 00, 00, 0xd,
	0x48, 0xbe, 0x5f, 00, 00, 00, 00, 0x18, 0x16,
	0x88, 0xd2, 0x5b, 0x7e, 00, 00, 0x40, 0xec, 0xa6,
	0xfd, 0x85, 0x7c, 0xc, 0x93, 0x59, 0x29, 0xbe, 0x5f,
	00, 00, 00, 00, 0x4d, 0x7a, 0xab, 0x1b, 0x5b,
	0x7e, 00, 00, 0x72, 0xab, 0x5, 0x9c, 0x46, 0x3c,
	0x89, 0x83, 0xc6, 0xfc, 0xbd, 0x5f, 00, 00, 00,
	00, 0x6d, 0x57, 0x64, 0x83, 0x5b, 0x7e, 00, 00,
	0x81, 0x9a, 0x8, 0x1a, 0x14, 0xb3, 0xd5, 0x6, 0x16,
	0x8b, 0xbd, 0x5f, 00, 00, 00, 00, 0x2e, 0xfc,
	0xd4, 0x83, 0x5b, 0x7e, 00, 00, 0xe4, 0x47, 0xf5,
	0xa8, 0xcc, 0x96, 0xc7, 0xa7, 0xd5, 0x59, 0xbd, 0x5f,
	00, 00, 00, 00, 0x5b, 0xe0, 0x55, 0x61, 0x5b,
	0x7e, 00, 00, 0xc3, 0x6, 0x92, 0x73, 0x81, 0xe1,
	0x4a, 0xbd, 0x1d, 0x44, 0xbd, 0x5f, 00, 00, 00,
	00, 0xc2, 0x35, 0xc5, 0x2f, 0x5b, 0x7e, 00, 00,
	0xe1, 0x2d, 0xe4, 0x35, 0xb7, 0x19, 0xc2, 0x8f, 0x32,
	0xec, 0xbc, 0x5f, 00, 00, 00, 00, 0x5f, 0x84,
	0xdc, 0x25, 0x5b, 0x7e, 00, 00, 0x3, 0x2b, 0xfc,
	0xae, 0x1d, 0xed, 0x35, 0x26, 0xd, 0xbb, 0xbc, 0x5f,
	00, 00, 00, 00, 0x5e, 0xf4, 0x1c, 0x31, 0x5b,
	0x7e, 00, 00, 0xb, 00, 00, 00, 00, 00,
	00, 00, 0xe3, 0x8c, 0xbb, 0x5f, 00, 00, 00,
	00, 0x1f, 0x80, 0xe7, 0xbe, 0x5b, 0x7e, 00, 00,
	0x2e, 0x12, 0x78, 0xdf, 0x38, 0xec, 0x3e, 0x56, 0xcf,
	0x60, 0xba, 0x5f, 00, 00, 00, 00, 0xb9, 0xb8,
	0xa8, 0x2b, 0x5b, 0x7e, 00, 00, 0xa, 0x5e, 0x2a,
	0x38, 0x5d, 0xb8, 0xd9, 0x5b, 0xce, 0x52, 0xba, 0x5f,
	00, 00, 00, 00, 0x2e, 0x76, 0x5b, 0xd7, 0x5b,
	0x7e, 00, 00, 0x30, 0x83, 0xbf, 0x8b, 0xd0, 0x94,
	0xfb, 0xbc, 0x5a, 0x43, 0xba, 0x5f, 00, 00, 00,
	00, 0xc2, 0xbb, 0x68, 0x29, 0x5b, 0x7e, 00, 00,
	0xfa, 0xd6, 0x99, 0x49, 0x57, 0xb, 0xd8, 0xd1, 0x99,
	0xff, 0xb9, 0x5f, 00, 00, 00, 00, 0x5, 0xf8,
	0x25, 0xab, 0x5b, 0x7e, 00, 00, 0x8d, 0x93, 0x3e,
	0xfc, 0xf0, 0xeb, 0xd6, 0xc3, 0x2d, 0x59, 0xb9, 0x5f,
	00, 00, 00, 00, 0x6d, 0x57, 0xb7, 0x82, 0x5b,
	0x7e, 00, 00, 0xef, 0x14, 0x69, 0x78, 0x10, 0xa2,
	0xd1, 0xfe, 0xf2, 0x56, 0xb9, 0x5f, 00, 00, 00,
	00, 0x5d, 0x4f, 0xc5, 0x8, 0x5b, 0x7e, 00, 00,
	0x9c, 0xef, 0x68, 0x46, 0x47, 0xc5, 0xbe, 0x75, 0x1c,
	0x2d, 0xb8, 0x5f, 00, 00, 00, 00, 0x4e, 0x1e,
	0xe8, 0xd5, 0x5b, 0x7e, 00, 00, 0x45, 0x49, 0xad,
	0x21, 0xa3, 0xe9, 0x93, 0xc5, 0xe2, 0x11, 0xb8, 0x5f,
	00, 00, 00, 00, 0x2e, 0xb9, 0x69, 0x49, 0x5b,
	0x7e, 00, 00, 0x89, 0x5e, 0xcc, 0x52, 0x76, 0xcf,
	0x5, 0xc2, 0x1a, 0xd4, 0xb7, 0x5f, 00, 00, 00,
	00, 0x51, 0xa2, 0xe6, 0x54, 0x5b, 0x7e, 00, 00,
	0x44, 0x73, 0xa6, 0xbd, 0xd4, 0xcc, 0x5f, 0xad, 0x38,
	0xb8, 0xb7, 0x5f, 00, 00, 00, 00, 0xb2, 0x9e,
	0xc5, 0x32, 0x5b, 0x7e, 00, 00, 0x9e, 0xa, 0x27,
	0x5d, 0x4d, 0xf0, 0x63, 0x1, 0x63, 0xc7, 0xb6, 0x5f,
	00, 00, 00, 00, 0x6f, 0xdd, 0xa1, 0x9d, 0x5b,
	0x7e, 00, 00, 0x16, 0xb4, 0xb4, 0xb9, 0x22, 0x5a,
	0x5f, 0xd1, 0x3c, 0x3e, 0xb6, 0x5f, 00, 00, 00,
	00, 0x55, 0xc4, 0xcc, 0xd4, 0x5b, 0x7e, 00, 00,
	0xa1, 0x5c, 0x96, 0x20, 0x35, 0x84, 0xc7, 0x79, 0x21,
	0xd4, 0xb5, 0x5f, 00, 00, 00, 00, 0xbc, 0xef,
	0x33, 0xc2, 0x5b, 0x7e, 00, 00, 0x23, 0x39, 0xf0,
	0xd4, 0x3e, 0xae, 0x15, 0xcb, 0x7a, 0x9b, 0xb5, 0x5f,
	00, 00, 00, 00, 0x5c, 0xf9, 0x7f, 0x45, 0x5b,
	0x7e, 00, 00, 0x42, 0x7a, 0x8a, 0x4b, 0xed, 0x5f,
	0x26, 00, 0x24, 0x55, 0xb5, 0x5f, 00, 00, 00,
	00, 0xd4, 0x73, 0xe9, 0x4, 0x5b, 0x7e, 00, 00,
	0x73, 0x9a, 0xb5, 0x84, 0x46, 0x78, 0x56, 0x42, 0x8b,
	0xfa, 0xb4, 0x5f, 00, 00, 00, 00, 0xd5, 0xae,
	0x13, 0x4b, 0x5b, 0x7e, 00, 00, 0x23, 0x7b, 0x54,
	0x4d, 0x2a, 0xcd, 0xd2, 0x66, 0x35, 0xf7, 0xb4, 0x5f,
	00, 00, 00, 00, 0x2e, 0x27, 0x45, 0xd6, 0x5b,
	0x7e, 00, 00, 0x3e, 0xec, 0x97, 0x99, 0x48, 0x81,
	0x99, 0x38, 0xa2, 0x16, 0xb4, 0x5f, 00, 00, 00,
	00, 0x5, 0x39, 0x41, 0x8f, 0x5b, 0x7e, 00, 00,
	0x1, 0xd3, 0xa8, 0xb, 0x54, 0x5d, 0x32, 0x2d, 0xa8,
	0x3, 0xb4, 0x5f, 00, 00, 00, 00, 0x2e, 0xad,
	0x83, 0xf2, 0x5b, 0x7e, 00, 00, 0xe8, 0x3, 00,
	00, 00, 00, 00, 00, 0xb7, 0xfa, 0xb3, 0x5f,
	00, 00, 00, 00, 0xb9, 0x91, 0xb4, 0x34, 0x5b,
	0x7e, 00, 00, 0x17, 0x8c, 0xf8, 0x8f, 0x2f, 0xbe,
	0xfa, 0x71, 0x95, 0xb9, 0xb3, 0x5f, 00, 00, 00,
	00, 0x55, 0xc6, 0x8d, 0x5a, 0x5b, 0x7e, 00, 00,
	0x32, 0x2b, 0xd0, 0x56, 0xc4, 0x2b, 0x6f, 0xb2, 0x41,
	0x69, 0xb3, 0x5f, 00, 00, 00, 00, 0x2b, 0xfc,
	0x12, 0x78, 0x5b, 0x7e, 00, 00, 0xc0, 0x6, 0x43,
	0x55, 0xe2, 0xf4, 0x65, 0x19, 0x67, 0x4f, 0xb3, 0x5f,
	00, 00, 00, 00, 0x5, 0xff, 0xb8, 0xdf, 0x5b,
	0x7e, 00, 00, 0x3e, 0xec, 0x97, 0x99, 0x48, 0x81,
	0x99, 0x38, 0x86, 0xb5, 0xb2, 0x5f, 00, 00, 00,
	00, 0x76, 0x46, 0xe9, 0x27, 0x5b, 0x7e, 00, 00,
	0x2f, 0x1f, 0x7a, 0x2b, 0x61, 0xe, 0x2f, 0x7d, 0xd7,
	0xa7, 0xb2, 0x5f, 00, 00, 00, 00, 0xb9, 0xce,
	0x24, 0x2d, 0x5b, 0x7e, 00, 00, 0xa, 0x5e, 0x2a,
	0x38, 0x5d, 0xb8, 0xd9, 0x5b, 0x1d, 0x64, 0xb1, 0x5f,
	00, 00, 00, 00, 0x5f, 0x43, 0x52, 0x20, 0x5b,
	0x7e, 00, 00, 0x88, 0xd5, 0x25, 0x42, 0x41, 0x20,
	0x1b, 0x64, 0x64, 0x6a, 0xb0, 0x5f, 00, 00, 00,
	00, 0x5c, 0x71, 0x52, 0x6f, 0x5b, 0x7e, 00, 00,
	0x80, 0x20, 0xe0, 0xb8, 0x53, 0xad, 0x32, 0x2b, 0x2,
	0x50, 0xb0, 0x5f, 00, 00, 00, 00, 0x86, 0xf9,
	0xb0, 0x8b, 0x5b, 0x7e, 00, 00, 0x3c, 0x52, 0xe3,
	0xb6, 0x99, 0x7b, 0xe3, 0x37, 0x58, 0x12, 0xb0, 0x5f,
	00, 00, 00, 00, 0xb2, 0x36, 0xc0, 0x30, 0x5b,
	0x7e, 00, 00, 0x84, 0x6b, 0x4b, 0xeb, 0x2, 0x51,
	0x6a, 0x53, 0x6d, 0x34, 0xaf, 0x5f, 00, 00, 00,
	00, 0xc3, 0x1a, 0x5d, 0xd3, 0x5b, 0x7e, 00, 00,
	0xd5, 0x2, 00, 00, 00, 00, 00, 00, 0x51,
	0x67, 0xae, 0x5f, 00, 00, 00, 00, 0xc3, 0x36,
	0x2a, 0xc7, 0x5b, 0x7e, 00, 00, 0xac, 0x99, 0x1,
	0x15, 0xce, 0x66, 0x8d, 0xeb, 0xcc, 0x2b, 0xae, 0x5f,
	00, 00, 00, 00, 0x5b, 0xf6, 0x4, 0x90, 0x5b,
	0x7e, 00, 00, 0x5c, 0x71, 0xe3, 0x54, 0x79, 0xa0,
	0x13, 0xc0, 0x69, 0x1c, 0xab, 0x5f, 00, 00, 00,
	00, 0xb0, 0x69, 0x10, 0xf4, 0x5b, 0x7e, 00, 00,
	0x34, 0xf8, 0x80, 0xc, 0xb3, 0xb6, 0xaf, 0xfc, 0x70,
	0xad, 0xa9, 0x5f, 00, 00, 00, 00, 0xc9, 0x5d,
	0x3a, 0xa9, 0x5b, 0x7e, 00, 00, 0xa3, 0x3e, 0xe1,
	0xcb, 0x13, 0x7d, 0x54, 0x3a, 0x7d, 0x48, 0xa8, 0x5f,
	00, 00, 00, 00, 0x18, 0xf1, 0xf4, 0x42, 0x5b,
	0x7e, 00, 00, 0x13, 0x2f, 0xea, 0xbc, 0x93, 0x6b,
	0x7a, 0xe2, 0xfc, 0x2f, 0xa7, 0x5f, 00, 00, 00,
	00, 0xb9, 0x5, 0x69, 0xee, 0x5b, 0x7e, 00, 00,
	0x74, 0x60, 0x62, 0x4e, 0xc9, 0x43, 0xde, 0x52, 0xdd,
	0x54, 0xa6, 0x5f, 00, 00, 00, 00, 0x6d, 0x7a,
	0x1e, 0x84, 0x5b, 0x7e, 00, 00, 0x9c, 0x15, 0x8b,
	0x70, 0x9c, 0x51, 0xa8, 00, 0x37, 0x6d, 0xa4, 0x5f,
	00, 00, 00, 00, 0x6d, 0xfb, 0xda, 0x32, 0x5b,
	0x7e, 00, 00, 0x5e, 0x78, 0xf2, 0xa2, 0xc2, 0x2,
	0xeb, 0xfb, 0xb7, 0x7, 0xa4, 0x5f, 00, 00, 00,
	00, 0x40, 0xe2, 0xac, 0xa, 0x5b, 0x7e, 00, 00,
	0x1f, 0x8a, 0x79, 0xdd, 0x92, 0x82, 0xea, 0x52, 0xe7,
	0xca, 0xa2, 0x5f, 00, 00, 00, 00,
}

func TestDecodeHandshakeResponse(t *testing.T) {
	var decoded HandshakeResponse

	err := binary.Unmarshal(encodedHandshakeRes, &decoded)
	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, config.MainNet().NetworkID, decoded.NodeData.NetworkID)
	assert.Equal(t, 252, len(decoded.Peers))
	assert.Equal(t, uint64(0x6e6b855a2c75f543), decoded.Peers[34].ID)
	assert.Equal(t, uint32(0x56928fb9), decoded.Peers[104].Address.IP)
	assert.Equal(t, 32347, int(decoded.Peers[104].Address.Port))
}

func TestDecodeHandShakeRequest(t *testing.T) {

	var decoded HandshakeRequest

	err := binary.Unmarshal(encodedHandshakeReq, &decoded)

	mainnet := config.MainNet()

	assert.Nil(t, err)

	assert.Equal(t, mainnet.NetworkID, decoded.NodeData.NetworkID)
	assert.Equal(t, mainnet.P2PCurrentVersion, decoded.NodeData.Version)
	assert.Equal(t, uint64(0x310), decoded.NodeData.PeerID)
	assert.Equal(t, uint64(0x60061472), decoded.NodeData.LocalTime)
	assert.Equal(t, uint32(0x7e5b), decoded.NodeData.MyPort)

	assert.Equal(t, uint32(0x8cce7), decoded.PayloadData.CurrentHeight)
	assert.Equal(t,
		cryptonote.Hash{
			0x8, 0xc, 0xe5, 0xe5, 0x96, 0x77, 0x9, 0x4b,
			0xde, 0xda, 0xae, 0xea, 0xe9, 0xa8, 0x96, 0x4b,
			0x60, 0xfe, 0xf7, 0x26, 0x70, 0x9c, 0xf6, 0x5f,
			0x28, 0x71, 0x38, 0x6e, 0xa2, 0x2c, 0x62, 0x9e,
		},
		decoded.PayloadData.TopBlockHash,
	)

	encoded, err := binary.Marshal(decoded)

	assert.Nil(t, err)
	assert.Equal(t, encodedHandshakeReq, encoded)
}

