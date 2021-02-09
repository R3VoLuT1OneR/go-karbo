package binary

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

type testElement struct {
	Enabled bool `binary:"boolean_key"`

	UInteger8 uint8 `binary:"uint8"`
	UInteger16 uint16 `binary:"uint16"`
	UInteger32 uint32 `binary:"uint32"`
	UInteger64 uint64 `binary:"uint64"`

	Integer8 int8 `binary:"int8"`
	Integer16 int16 `binary:"int16"`
	Integer32 int32 `binary:"int32"`
	Integer64 int64 `binary:"int64"`

	Float64 float64 `binary:"float64"`

	StringData string `binary:"string_data"`

	Child testChildElement `binary:"child"`
}

type testChildElement struct {
	SomeData string `binary:"some_data"`
	Enabled bool `binary:"boolean_key"`
}

func TestSimpleTestElement(t *testing.T)  {
	testData1 := testElement{
		StringData: "hello",
		UInteger8: math.MaxUint8 - 1,
		UInteger16: math.MaxUint8 - 3,
		UInteger32: math.MaxUint8 - 4,
		UInteger64: math.MaxUint8 - 9,

		Integer8: math.MaxInt8 - 122,
		Integer16: math.MaxInt16 - 123,
		Integer32: math.MaxInt32 - 124,
		Integer64: math.MaxInt64 - 126,

		Float64: math.MaxFloat64 - 239.39383,

		Enabled: true,

		Child: testChildElement{
			SomeData: "string_child_data",
			Enabled: false,
		},
	}

	encoded, err := Marshal(testData1)
	assert.Nil(t, err)

	var decoded testElement
	err = Unmarshal(encoded, &decoded)
	assert.Nil(t, err)
	assert.Equal(t, testData1, decoded)
}

