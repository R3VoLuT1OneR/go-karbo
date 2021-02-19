package binary

import (
	"reflect"
)

const (
	// #define MAX_STRING_LEN_POSSIBLE       2000000000 //do not let string be so big
	maxStringLen = 2000000000 //do not let string be so big

	storageSignatureA uint32 = 0x01011101
	storageSignatureB uint32 = 0x01020101
	storageFormatVer byte = 1
	headSize int = 9

	rawSizeMarkMask byte = 0x03

	rawSizeMarkByte uint8 = 0
	rawSizeMarkWord uint16 = 1
	rawSizeMarkDWord uint32 = 2
	rawSizeMarkInt64 uint64 = 3

	typeInt64   byte = 1
	typeInt32   byte = 2
	typeInt16   byte = 3
	typeInt8    byte = 4
	typeUInt64  byte = 5
	typeUInt32  byte = 6
	typeUInt16  byte = 7
	typeUInt8   byte = 8
	typeFloat64 byte = 9 // C++ double
	typeBinary  byte = 10
	typeBool    byte = 11
	typeObject  byte = 12
	typeArray   byte = 13

	//const uint8_t BIN_KV_SERIALIZE_FLAG_ARRAY = 0x80;
	flagArray byte = 0x80
)

var mapBTypeToSimpleKind = map[byte]reflect.Kind{
	typeUInt8:   reflect.Uint8,
	typeUInt16:  reflect.Uint16,
	typeUInt32:  reflect.Uint32,
	typeUInt64:  reflect.Uint64,
	typeInt8:    reflect.Int8,
	typeInt16:   reflect.Int16,
	typeInt32:   reflect.Int32,
	typeInt64:   reflect.Int64,
	typeFloat64: reflect.Float64,
	typeBool:    reflect.Bool,
}

var mapSimpleKindToBType = map[reflect.Kind]byte{
	reflect.Uint8:   typeUInt8,
	reflect.Uint16:  typeUInt16,
	reflect.Uint32:  typeUInt32,
	reflect.Uint64:  typeUInt64,
	reflect.Int8:    typeInt8,
	reflect.Int16:   typeInt16,
	reflect.Int32:   typeInt32,
	reflect.Int64:   typeInt64,
	reflect.Float64: typeFloat64,
	reflect.Bool:    typeBool,
}