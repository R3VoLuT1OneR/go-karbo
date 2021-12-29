package crypto

import (
	"bufio"
	"encoding/hex"
	"fmt"
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
		h.FromBytes(testBytes)

		assert.Equal(t, expected, h[:])
		times++
	}

	assert.Equal(t, 320, times)
}

func TestHashToScalar(t *testing.T) {
	file, err := os.Open("./fixtures/hash_to_scalar.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		dataBytes, _ := hex.DecodeString(line[0])
		expectedBytes, _ := hex.DecodeString(line[1])

		var expected EllipticCurveScalar
		copy(expected[:], expectedBytes[:32])

		actualHash := HashFromBytes(dataBytes)
		actual := actualHash.toScalar()

		assert.Equal(t, expected, actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
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
			copy(h[:], testBytes[(i*32):(i*32)+32])
			list = append(list, h)
		}

		mh := list.MerkleRootHash()

		assert.Equal(t, expected, mh[:])
		times++
	}

	assert.Equal(t, 16, times)
}

func TestHash_ToEc(t *testing.T) {
	file, err := os.Open("./fixtures/hash_to_ec.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		pkBytes, _ := hex.DecodeString(line[0])
		expectedBytes, _ := hex.DecodeString(line[1])

		var expected EllipticCurvePoint
		copy(expected[:], expectedBytes[:32])

		hash := HashFromBytes(pkBytes)
		actual, actualErr := hash.toPoint()

		assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
}
