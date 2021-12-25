package cryptonote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dataUint8beTo64 = []*struct {
	i uint64
	b []byte
}{
	{0x0000000000000001, []byte{0x1}},
	{0x0000000000000102, []byte{0x1, 0x2}},
	{0x0000000000010203, []byte{0x1, 0x2, 0x3}},
	{0x0000000001020304, []byte{0x1, 0x2, 0x3, 0x4}},
	{0x0000000102030405, []byte{0x1, 0x2, 0x3, 0x4, 0x5}},
	{0x0000010203040506, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}},
	{0x0001020304050607, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7}},
	{0x0102030405060708, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8}},
}

var encodeBlockPositive = []*struct {
	b []byte
	s string
}{
	{[]byte{0x00}, "11"},
	{[]byte{0x39}, "1z"},
	{[]byte{0xFF}, "5Q"},

	{[]byte{0x00, 0x00}, "111"},
	{[]byte{0x00, 0x39}, "11z"},
	{[]byte{0x01, 0x00}, "15R"},
	{[]byte{0xFF, 0xFF}, "LUv"},

	{[]byte{0x00, 0x00, 0x00}, "11111"},
	{[]byte{0x00, 0x00, 0x39}, "1111z"},
	{[]byte{0x01, 0x00, 0x00}, "11LUw"},
	{[]byte{0xFF, 0xFF, 0xFF}, "2UzHL"},

	{[]byte{0x00, 0x00, 0x00, 0x39}, "11111z"},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF}, "7YXq9G"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x39}, "111111z"},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, "VtB5VXc"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x39}, "11111111z"},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, "3CUsUpv9t"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39}, "111111111z"},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, "Ahg1opVcGW"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39}, "1111111111z"},
	{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, "jpXCZedGfVQ"},

	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "11111111111"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, "11111111112"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08}, "11111111119"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09}, "1111111111A"},
	{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3A}, "11111111121"},
	{[]byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, "1Ahg1opVcGW"},
	{[]byte{0x06, 0x15, 0x60, 0x13, 0x76, 0x28, 0x79, 0xF7}, "22222222222"},
	{[]byte{0x05, 0xE0, 0x22, 0xBA, 0x37, 0x4B, 0x2A, 0x00}, "1z111111111"},
}

var encodePositive = []*struct {
	s string
	b []byte
}{
	{"", []byte{}},
	{"11", []byte{0x00}},
	{"5Q", []byte{0xFF}},
	{"111", []byte{0x00, 0x00}},
	{"LUv", []byte{0xFF, 0xFF}},
	{"11111", []byte{0x00, 0x00, 0x00}},
	{"2UzHL", []byte{0xFF, 0xFF, 0xFF}},
	{"111111", []byte{0x00, 0x00, 0x00, 0x00}},
	{"7YXq9G", []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	{"1111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
	{"VtB5VXc", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"3CUsUpv9t", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"1111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"Ahg1opVcGW", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"11111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQ", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"1111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQ5Q", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"11111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQLUv", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"1111111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQ2UzHL", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"11111111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQ7YXq9G", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"111111111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQVtB5VXc", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"22222222222VtB5VXc", []byte{0x06, 0x15, 0x60, 0x13, 0x76, 0x28, 0x79, 0xF7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"11111111111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQ3CUsUpv9t", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"111111111111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQAhg1opVcGW", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"1111111111111111111111", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"jpXCZedGfVQjpXCZedGfVQ", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var decodeBlockNegative = []string{
	"1",
	"z",
	"5R",
	"zz",
	"LUw",
	"zzz",
	"1111",
	"zzzz",
	"2UzHM",
	"zzzzz",
	"7YXq9H",
	"zzzzzz",
	"VtB5VXd",
	"zzzzzzz",
	"11111111",
	"zzzzzzzz",
	"3CUsUpv9u",
	"zzzzzzzzz",
	"Ahg1opVcGX",
	"zzzzzzzzzz",
	"jpXCZedGfVR",
	"zzzzzzzzzzz",
	"01111111111",
	"11111111110",
	"11111011111",
	"I1111111111",
	"O1111111111",
	"l1111111111",
	"_1111111111",
}

var decodeNegative = []string{
	"1",
	"z",
	"1111",
	"zzzz",
	"11111111",
	"zzzzzzzz",
	"123456789AB1",
	"123456789ABz",
	"123456789AB1111",
	"123456789ABzzzz",
	"123456789AB11111111",
	"123456789ABzzzzzzzz",
	"5R",
	"zz",
	"LUw",
	"zzz",
	"2UzHM",
	"zzzzz",
	"7YXq9H",
	"zzzzzz",
	"VtB5VXd",
	"zzzzzzz",
	"3CUsUpv9u",
	"zzzzzzzzz",
	"Ahg1opVcGX",
	"zzzzzzzzzz",
	"jpXCZedGfVR",
	"zzzzzzzzzzz",
	"123456789AB5R",
	"123456789ABzz",
	"123456789ABLUw",
	"123456789ABzzz",
	"123456789AB2UzHM",
	"123456789ABzzzzz",
	"123456789AB7YXq9H",
	"123456789ABzzzzzz",
	"123456789ABVtB5VXd",
	"123456789ABzzzzzzz",
	"123456789AB3CUsUpv9u",
	"123456789ABzzzzzzzzz",
	"123456789ABAhg1opVcGX",
	"123456789ABzzzzzzzzzz",
	"123456789ABjpXCZedGfVR",
	"123456789ABzzzzzzzzzzz",
	"zzzzzzzzzzz11",
	"10",
	"11I",
	"11O11",
	"11l111",
	"11_11111111",
	"1101111111111",
	"11I11111111111111",
	"11O1111111111111111111",
	"1111111111110",
	"111111111111l1111",
	"111111111111_111111111",
}

var encodeDecodeAddr = []*struct {
	addr string
	tag  uint64
	data []byte
}{
	{"21D35quxec71111111111111111111111111111111111111111111111111111111111111111111111111111116Q5tCH", 6, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"2Aui6ejTFscjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQVqegMoV", 6, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"1119XrkPuSmLzdHXgVgrZKjepg5hZAxffLzdHXgVgrZKjepg5hZAxffLzdHXgVgrZKjepg5hZAxffLzdHXgVgrZKVphZRvn", 0, []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}},
	{"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111115TXfiA", 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"PuT7GAdgbA83qvSEivPLYo11111111111111111111111111111111111111111111111111111111111111111111111111111169tWrH", 0x1122334455667788, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"PuT7GAdgbA841d7FXjswpJjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQjpXCZedGfVQVq4LL1v", 0x1122334455667788, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	{"PuT7GAdgbA819VwdWVDP", 0x1122334455667788, []byte{0x11}},
	{"PuT7GAdgbA81efAfdCjPg", 0x1122334455667788, []byte{0x22, 0x22}},
	{"PuT7GAdgbA83sryEt3YC8Q", 0x1122334455667788, []byte{0x33, 0x33, 0x33}},
	{"PuT7GAdgbA83tWUuc54PFP3b", 0x1122334455667788, []byte{0x44, 0x44, 0x44, 0x44}},
	{"PuT7GAdgbA83u9zaKrtRKZ1J6", 0x1122334455667788, []byte{0x55, 0x55, 0x55, 0x55, 0x55}},
	{"PuT7GAdgbA83uoWF3eanGG1aRoG", 0x1122334455667788, []byte{0x66, 0x66, 0x66, 0x66, 0x66, 0x66}},
	{"PuT7GAdgbA83vT1umSHMYJ4oNVdu", 0x1122334455667788, []byte{0x77, 0x77, 0x77, 0x77, 0x77, 0x77, 0x77}},
	{"PuT7GAdgbA83w6XaVDyvpoGQBEWbB", 0x1122334455667788, []byte{0x88, 0x88, 0x88, 0x88, 0x88, 0x88, 0x88, 0x88}},
	{"PuT7GAdgbA83wk3FD1gW7J2KVGofA1r", 0x1122334455667788, []byte{0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99}},
	{"15p2yAV", 0, []byte{}},
	{"FNQ3D6A", 0x7F, []byte{}},
	{"26k9QWweu", 0x80, []byte{}},
	{"3BzAD7n3y", 0xFF, []byte{}},
	{"11efCaY6UjG7JrxuB", 0, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}},
	{"21rhHRT48LN4PriP9", 6, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}},
}

var decodeAddrNegative = []*struct {
	addr string
	err  string
}{
	{"zuT7GAdgbA819VwdWVDP", "Block overflow"},
	{"0uT7GAdgbA819VwdWVDP", "invalid character '0' found"},
	{"IuT7GAdgbA819VwdWVDP", "invalid character 'I' found"},
	{"OuT7GAdgbA819VwdWVDP", "invalid character 'O' found"},
	{"luT7GAdgbA819VwdWVDP", "invalid character 'l' found"},
	{string(rune(0x00)) + "uT7GAdgbA819VwdWVDP", "invalid character '\x00' found"},
	{"PuT7GAdgbA819VwdWVD", "invalid block size"},
	{"11efCaY6UjG7JrxuC", "invalid checksum"},
	// {"jerj2e4mESo", "handles_non_correct_tag"}, // "jerj2e4mESo" == "\xFF\x00\xFF\xFF\x5A\xD9\xF1\x1C"
	{"1", "invalid block size"},
	{"1111", "invalid block size"},
	{"11", "Decoded size is too short 1"},
	{"111", "Decoded size is too short 2"},
	{"11111", "Decoded size is too short 3"},
	{"111111", "Decoded size is too short 4"},
	{"999999", "Block overflow"},
	{"ZZZZZZ", "Block overflow"},
}

func TestUint8beTo64(t *testing.T) {
	for _, td := range dataUint8beTo64 {
		assert.Equal(t, td.i, uint8beTo64(td.b))
	}
}

func TestUint64ToUint8be(t *testing.T) {
	for _, td := range dataUint8beTo64 {
		assert.Equal(t, td.b, uint64toUint8be(td.i))
	}
}

func TestDecodeBlockPositive(t *testing.T) {
	for _, td := range encodeBlockPositive {
		res, err := decodeBlock(td.s)
		assert.Nil(t, err)
		assert.Equal(t, td.b, res)
	}
}

func TestDecodeBlockNegative(t *testing.T) {
	for _, td := range decodeBlockNegative {
		res, err := decodeBlock(td)
		assert.NotNilf(t, err, "No error found in test '%s'", td)
		assert.Nil(t, res, "Response is not nil for test '%s'", td)
	}
}

func TestDecodeNegative(t *testing.T) {
	for _, td := range decodeNegative {
		res, err := decode(td)
		assert.NotNilf(t, err, "No error found in test '%s'", td)
		assert.Nil(t, res, "Response is not nil for test '%s'", td)
	}
}

func TestPositiveEncodeBlocks(t *testing.T) {
	for _, td := range encodeBlockPositive {
		assert.Equal(t, td.s, encodeBlock(td.b))
	}
}

func TestPositiveEncode(t *testing.T) {
	for _, td := range encodePositive {
		assert.Equal(t, td.s, encode(td.b))
	}
}

func TestEncodeDecodeAddr(t *testing.T) {
	for _, td := range encodeDecodeAddr {
		rencode := EncodeAddr(td.tag, td.data)

		assert.Equal(t, td.addr, rencode)

		tag, data, err := DecodeAddr(rencode)

		assert.Nil(t, err)
		assert.Equal(t, td.tag, tag)
		assert.Equal(t, td.data, data)
	}
}

func TestDecodeAddrNegative(t *testing.T) {
	for _, td := range decodeAddrNegative {
		tag, data, err := DecodeAddr(td.addr)

		assert.Equal(t, uint64(0), tag)
		assert.Nil(t, data)

		if err != nil {
			assert.Equal(t, td.err, err.Error())
		} else {
			assert.NotNil(t, err)
		}
	}
}
