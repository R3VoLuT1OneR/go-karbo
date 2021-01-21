package binary

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

// Decode binary data into provided interface
//
// v should be pointer to needed value
func (d *Decoder) Decode(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("interface must be pointer")
	}

	size, err := d.readVarInt()
	if err == io.EOF {
		return nil
	}

	mp, order, err := mapStructFields(v)
	if err != nil {
		return err
	}

	for i := uint64(0); i < size; i++ {
		name, err := d.readName()
		if err != nil {
			return err
		}

		if _, ok := mp[name]; !ok {
			return errors.New(fmt.Sprintf("field '%s' not found in %T", name, v))
		}

		if order[i] != name {
			return errors.New(fmt.Sprintf("field '%s' placed in wrong order", name))
		}

		field := mp[name]

		if err := d.readValue(field); err != nil {
			return errors.New(fmt.Sprintf("Error on '%s' field decode: %s", name, err))
		}
	}

	return nil
}

func (d *Decoder) readName() (string, error) {
	var sizeByte [1]byte

	if _, err := d.r.Read(sizeByte[:]); err != nil {
		return "", err
	}

	str := make([]byte, sizeByte[0])

	if _, err := d.r.Read(str); err != nil {
		return "", err
	}

	return string(str), nil
}

func (d *Decoder) readValue(value reflect.Value) error {
	typeByte := make([]byte, 1)
	if _, err := d.r.Read(typeByte[:]); err != nil {
		return nil
	}

	// Some simple kinds can be read with binary package.
	// It is gonna read exact amount of needed bytes and put them as value into value reflection.
	if value.Kind() == mapBTypeToSimpleKind[typeByte[0]] {
		if err := binary.Read(d.r, binary.LittleEndian, value.Addr().Interface()); err != nil {
			return err
		}

		return nil
	}

	switch typeByte[0] {
	// In binary types like String, Slice, Array data can be encoded.
	// We depend on receiving interface to define how exactly we need to read the data.
	case typeBinary:
		size, err := d.readVarInt()
		if err != nil {
			return err
		}

		b := make([]byte, size)
		if _, err := d.r.Read(b); err != nil {
			return err
		}

		switch value.Kind() {
		// For slice we don't know exact size of the data encoded.
		// We read item by item with N size, where N defined by receiving slice element type.
		case reflect.Slice:
			newSlice := reflect.New(value.Type()).Elem()
			itemType := value.Type().Elem()
			reader := bytes.NewReader(b)

			for {
				item := reflect.New(itemType)
				err := binary.Read(reader, binary.LittleEndian, item.Interface())
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				newSlice = reflect.Append(newSlice, item.Elem())
			}

			value.Set(newSlice)
		// For binary data only [N]byte array supported.
		case reflect.Array:
			newArray := reflect.New(value.Type()).Elem()
			reflect.Copy(newArray, reflect.ValueOf(b))

			value.Set(newArray)
		case reflect.String:
			value.SetString(string(b))
		default:
			return errors.New(fmt.Sprintf("not supported kind '%s' for binary received", value.Kind()))
		}
	case typeObject:
		v := reflect.New(value.Type())
		if err := d.Decode(v.Interface()); err != nil {
			return err
		}

		value.Set(v.Elem())
	default:
		return errors.New(fmt.Sprintf("unknown value type %v", typeByte[0]))
	}

	return nil
}

func (d *Decoder) readVarInt() (uint64, error) {
	var sizeBytes [8]byte
	if _, err := d.r.Read(sizeBytes[0:1]); err != nil {
		return 0, err
	}

	var bytesLeft = 0
	switch sizeBytes[0] & rawSizeMarkMask {
	case rawSizeMarkByte:
		bytesLeft = 0
		break
	case byte(rawSizeMarkWord):
		bytesLeft = 1
		break
	case byte(rawSizeMarkDWord):
		bytesLeft = 3
		break
	case byte(rawSizeMarkInt64):
		bytesLeft = 7
		break
	}

	allBytes := make([]byte, 8)
	allBytes[0] = sizeBytes[0]

	if _, err := d.r.Read(allBytes[1:bytesLeft+1]); err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(allBytes) >> 2, nil
}
