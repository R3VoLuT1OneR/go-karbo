package binary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v interface{}, name string) error {
	if v == nil {
		return nil
	}

	typ := reflect.TypeOf(v)
	kind := typ.Kind()

	switch kind {
	case reflect.Bool:
		if err := e.writePrefix(typeBool, name); err != nil {
			return err
		}

		var val byte = 0
		if v.(bool) {
			val = 1
		}

		if _, err := e.w.Write([]byte{val}); err != nil {
			return err
		}
	case reflect.Uint8:
		if err := e.writePrefix(typeUInt8, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Uint16:
		if err := e.writePrefix(typeUInt16, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Uint32:
		if err := e.writePrefix(typeUInt32, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Uint64:
		if err := e.writePrefix(typeUInt64, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Int8:
		if err := e.writePrefix(typeInt8, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Int16:
		if err := e.writePrefix(typeInt16, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Int32:
		if err := e.writePrefix(typeInt32, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Int64:
		if err := e.writePrefix(typeInt64, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.Float64:
		if err := e.writePrefix(typeFloat64, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	case reflect.String:
		if err := e.writePrefix(typeBinary, name); err != nil {
			return err
		}

		if err := e.writeVarInt(uint64(len(v.(string)))); err != nil {
			return err
		}

		if _, err := e.w.Write([]byte(v.(string))); err != nil {
			return err
		}
	case reflect.Struct:
		if err := e.writePrefix(typeObject, name); err != nil {
			return err
		}

		fieldsMap, order, err := mapStructFields(v)
		if err != nil {
			return err
		}

		if err := e.writeVarInt(uint64(len(order))); err != nil {
			return err
		}

		for i := 0; i < len(order); i++ {
			name := order[i]
			fieldValue := fieldsMap[name]

			if err := e.Encode(fieldValue.Interface(), name); err != nil {
				return err
			}
		}
	case reflect.Array:
		if err := e.writePrefix(typeBinary, name); err != nil {
			return err
		}

		var arrayBytesBuf bytes.Buffer
		if err := binary.Write(&arrayBytesBuf, binary.LittleEndian, v); err != nil {
			return err
		}
		arrayBytes := arrayBytesBuf.Bytes()

		if err := e.writeVarInt(uint64(len(arrayBytes))); err != nil {
			return err
		}

		if _, err := e.w.Write(arrayBytes); err != nil {
			return err
		}
	case reflect.Slice:
		if err := e.writePrefix(typeBinary, name); err != nil {
			return err
		}

		if err := binary.Write(e.w, binary.LittleEndian, v); err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unsuported kind: %s", kind))
		//panic("unsupported type")
		//return errors.New(fmt.Sprintf("unsuported type: %T", v))
	}

	return nil
}

// writePrefix writes element name and byte of the type right after
func (e *Encoder) writePrefix(t byte, name string) error {
	if name == "" {
		return nil
	}

	if err := e.writeElementName(name); err != nil {
		return err
	}

	if _, err := e.w.Write([]byte{t}); err != nil {
		return err
	}

	return nil
}

// writeElementName writes byte of string length and then string as bytes
func (e *Encoder) writeElementName(name string) error {
	nameLen := len(name)
	if nameLen > math.MaxUint8 {
		return errors.New("element name is too long")
	}

	if _, err := e.w.Write([]byte{byte(nameLen)}); err != nil {
		return err
	}

	if _, err := e.w.Write([]byte(name)); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) writeVarInt(i uint64) (err error) {
	if i <= math.MaxUint8 {
		var v = (uint8(i) << 2) | rawSizeMarkByte
		var b = []byte{v}

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else if i <= math.MaxUint16 {
		var v = (uint16(i) << 2) | rawSizeMarkWord
		var b []byte
		binary.LittleEndian.PutUint16(b, v)

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else if i <= math.MaxUint32 {
		var v = (uint32(i) << 2) | rawSizeMarkDWord
		var b []byte
		binary.LittleEndian.PutUint32(b, v)

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else if i <= math.MaxUint64 {
		var v = (i << 2) | rawSizeMarkInt64
		var b []byte
		binary.LittleEndian.PutUint64(b, v)

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else {
		return errors.New("failed to pack varInt - too big amount")
	}

	return nil
}
