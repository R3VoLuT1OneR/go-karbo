package binary

import (
	"errors"
	"fmt"
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
	flagArray = 0x80

	annotationBinary = "binary"
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

type mapBinaryKeyValueField map[string]reflect.Value

// mapStructFields reads interface fields and returns map and fields order
// Order of the fields is very important for the encoding.
func mapStructFields(v interface{}) (mapBinaryKeyValueField, []string, error) {
	var mp = mapBinaryKeyValueField{}
	var order []string

	val := reflect.ValueOf(v)
	typ := val.Type()

	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}


	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if name, ok := field.Tag.Lookup(annotationBinary); ok {
			if _, ok := mp[name]; ok {
				return nil, nil, errors.New(fmt.Sprintf("duplicate key '%s' found", name))
			}

			mp[name] = val.Field(i)
			order = append(order, name)
		}
	}

	return mp, order, nil
}
