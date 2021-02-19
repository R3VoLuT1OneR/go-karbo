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

func (e *Encoder) encode(field metadata) error {
	typ, err := e.getType(field.value.Type(), field.asArray)
	if err != nil {
		return err
	}

	if err := e.writeElementPrefix(typ, field.name); err != nil {
		return err
	}

	if err := e.writeValue(field.value, field.asArray); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) getType(typ reflect.Type, asArray bool) (byte, error) {
	kind := typ.Kind()

	if t, ok := mapSimpleKindToBType[kind]; ok {
		return t, nil
	}

	switch kind {
	case reflect.String:
		return typeBinary, nil
	case reflect.Struct:
		return typeObject, nil
	case reflect.Array, reflect.Slice:
		if asArray {
			itemType, err := e.getType(typ.Elem(), false)
			if err != nil {
				return 0, err
			}

			return flagArray | itemType, nil
		}

		return typeBinary, nil
	}

	return 0, errors.New(fmt.Sprintf("unsuported kind: %s", kind))
}

func (e *Encoder) writeValue(val reflect.Value, asArray bool) error {
	kind := val.Kind()

	if _, ok := mapSimpleKindToBType[kind]; ok {
		if err := binary.Write(e.w, binary.LittleEndian, val.Interface()); err != nil {
			return err
		}

		return nil
	}

	switch kind {
	case reflect.String:
		if err := e.writeVarInt(uint64(val.Len())); err != nil {
			return err
		}

		if _, err := e.w.Write([]byte(val.Interface().(string))); err != nil {
			return err
		}
	case reflect.Struct:
		smd, err := getStructBinaryMetadata(val, true)
		if err != nil {
			return err
		}

		if err := e.writeVarInt(uint64(len(smd.order))); err != nil {
			return err
		}

		for i := 0; i < len(smd.order); i++ {
			if err := e.encode(smd.order[i]); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		if asArray {
			l := val.Len()
			if err := e.writeVarInt(uint64(l)); err != nil {
				return err
			}

			for i := 0; i < l; i++ {
				if err := e.writeValue(val.Index(i), false); err != nil {
					return err
				}
			}

			return nil
		}

		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.LittleEndian, val.Interface()); err != nil {
	    	return err
		}

		if err := e.writeVarInt(uint64(buf.Len())); err != nil {
			return err
		}

		if _, err := e.w.Write(buf.Bytes()); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unsuported kind: %s", kind))
	}

	return nil
}

// writeElementPrefix writes element name and byte of the type right after
func (e *Encoder) writeElementPrefix(t byte, name string) error {
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
	if i <= 63 {
		var v = (uint8(i) << 2) | rawSizeMarkByte
		var b = []byte{v}

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else if i <= 16383 {
		var v = (uint16(i) << 2) | rawSizeMarkWord
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, v)

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else if i <= 1073741823 {
		var v = (uint32(i) << 2) | rawSizeMarkDWord
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, v)

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	} else  {
		if i > 4611686018427387903 {
			return errors.New("failed to pack varInt - too big amount")
		}

		var v = (i << 2) | rawSizeMarkInt64
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)

		if _, err := e.w.Write(b); err != nil {
			return err
		}
	}

	return nil
}
