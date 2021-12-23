package cryptonote

import (
	"bufio"
	"encoding/hex"
	"github.com/r3volut1oner/go-karbo/crypto"
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

		var h crypto.Hash
		h.FromBytes(testBytes)

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

		list := crypto.HashList{}
		listLen := len(testBytes) / 32
		for i := 0; i < listLen; i++ {
			var h crypto.Hash
			copy(h[:], testBytes[(i*32):(i*32)+32])
			list = append(list, h)
		}

		mh := list.MerkleRootHash()

		assert.Equal(t, expected, mh[:])
		times++
	}

	assert.Equal(t, 16, times)
}
