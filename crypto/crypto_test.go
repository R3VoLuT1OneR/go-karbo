package crypto

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

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

		actual := HashToScalar(dataBytes)

		assert.Equal(t, expected, actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
}

func TestGenerateKeyImage(t *testing.T) {
	file, err := os.Open("./fixtures/generate_key_image.txt")
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
		skBytes, _ := hex.DecodeString(line[1])
		expectedBytes, _ := hex.DecodeString(line[2])

		var publicKey PublicKey
		copy(publicKey[:], pkBytes[:32])

		var privateKey SecretKey
		copy(privateKey[:], skBytes[:32])

		var expected KeyImage
		copy(expected[:], expectedBytes[:32])

		actual, actualErr := GenerateKeyImage(publicKey, privateKey)

		assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
}

func TestGenerateKeyDerivation(t *testing.T) {
	file, err := os.Open("./fixtures/generate_key_derivation.txt")
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
		skBytes, _ := hex.DecodeString(line[1])
		expectedResult, _ := strconv.ParseBool(line[2])

		var publicKey PublicKey
		var privateKey SecretKey

		copy(publicKey[:], pkBytes[:32])
		copy(privateKey[:], skBytes[:32])

		actual, actualErr := GenerateKeyDerivation(publicKey, privateKey)

		if expectedResult {
			expectedBytes, _ := hex.DecodeString(line[3])
			var expected KeyDerivation
			copy(expected[:], expectedBytes[:32])

			assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
			assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))
		} else {
			assert.NotNil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		}

		lineNumber++
		times++
	}

	assert.Equal(t, 272, times)
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
		actual, actualErr := hash.ToPoint()

		assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
}
