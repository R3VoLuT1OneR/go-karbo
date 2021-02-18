package cryptonote

import (
	"bufio"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestHashFromBytes(t *testing.T) {
	file, err := os.Open("./fixtures/hash32.txt")
	check(err)
	defer file.Close()

	times := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cases := strings.Fields(scanner.Text())

		expected, err := hex.DecodeString(cases[0])
		check(err)

		testBytes, err := hex.DecodeString(cases[1])
		check(err)

		var h Hash
		h.FromBytes(&testBytes)

		assert.Equal(t, expected, h[:])
		times++
	}

	assert.Equal(t, 320, times)
}

func TestMerkleRootTree(t *testing.T) {
	file, err := os.Open("./fixtures/merkleTreeHash.txt")
	check(err)
	defer file.Close()

	times := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cases := strings.Fields(scanner.Text())

		expected, err := hex.DecodeString(cases[0])
		check(err)

		testBytes, err := hex.DecodeString(cases[1])
		check(err)

		list := HashList{}
		listLen := len(testBytes) / 32
		for i := 0; i < listLen; i++ {
			var h Hash
			copy(h[:], testBytes[(i*32):(i*32) + 32])
			list = append(list, h)
		}

		mh, err := list.merkleRootHash()
		check(err)

		assert.Equal(t, expected, mh[:])
		times++
	}

	assert.Equal(t, 16, times)
}
