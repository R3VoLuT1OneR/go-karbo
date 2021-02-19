package binary

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const tagBinary = "binary"

const tagOptionAsArray = "array"
const tagOptionOmitEmpty = "omitempty"

type metadata struct {
	order int
	value reflect.Value
	name string

	asArray bool
	omitEmpty bool
}

type structMetadata struct {
	fields map[string]metadata
	order []metadata
}

// getStructBinaryMetadata reads interface fields and returns map and fields order
// Order of the fields is very important for the encoding.
func getStructBinaryMetadata(val reflect.Value, omitEmpty bool) (*structMetadata, error) {
	smd := structMetadata{
		fields: map[string]metadata{},
		order: []metadata{},
	}

	typ := val.Type()

	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// TODO: Implement metadata caching
		if tagString, ok := field.Tag.Lookup(tagBinary); ok {
			tagValues := strings.Split(tagString, ",")

			if len(tagValues) == 0 {
				return nil, errors.New("missing field name")
			}

			md := metadata{
				name: tagValues[0],
				value: val.Field(i),
				order: len(smd.order),
			}

			for _, tagValue := range tagValues[1:] {
				switch tagValue {
				case tagOptionAsArray:
					md.asArray = true
				case tagOptionOmitEmpty:
					md.omitEmpty = true
				}
			}

			if _, ok := smd.fields[md.name]; ok {
				return nil, errors.New(fmt.Sprintf("duplicate key '%s' found", md.name))
			}

			smd.fields[md.name] = md

			if !omitEmpty || !(md.omitEmpty && md.value.IsZero()) {
				smd.order = append(smd.order, md)
			}
		}
	}

	return &smd, nil
}
