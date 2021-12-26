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

func TestDerivePublicKey(t *testing.T) {
	file, err := os.Open("./fixtures/derive_public_key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		derivationBytes, _ := hex.DecodeString(line[0])
		outputIndex, _ := strconv.ParseUint(line[1], 0, 64)
		baseBytes, _ := hex.DecodeString(line[2])
		expectedResult, _ := strconv.ParseBool(line[3])

		var derivation KeyDerivation
		copy(derivation[:], derivationBytes)

		var base PublicKey
		copy(base[:], baseBytes)

		actual, actualErr := derivation.derivePublicKey(outputIndex, &base)

		if expectedResult {
			expectedBytes, _ := hex.DecodeString(line[4])
			var expected PublicKey
			copy(expected[:], expectedBytes)

			assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
			assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))
		} else {
			assert.NotNil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		}

		lineNumber++
		times++
	}

	assert.Equal(t, 288, times)
}

func TestDeriveSecretKey(t *testing.T) {
	file, err := os.Open("./fixtures/derive_secret_key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		derivationBytes, _ := hex.DecodeString(line[0])
		outputIndex, _ := strconv.ParseUint(line[1], 0, 64)
		baseBytes, _ := hex.DecodeString(line[2])
		expectedBytes, _ := hex.DecodeString(line[3])

		var derivation KeyDerivation
		copy(derivation[:], derivationBytes)

		var base PublicKey
		copy(base[:], baseBytes)

		var expected SecretKey
		copy(expected[:], expectedBytes[:])

		actual, actualErr := derivation.deriveSecretKey(outputIndex, &base)

		assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
}

func TestUnderiveSecretKey(t *testing.T) {
	file, err := os.Open("./fixtures/underive_public_key.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		derivationBytes, _ := hex.DecodeString(line[0])
		outputIndex, _ := strconv.ParseUint(line[1], 0, 64)
		derivedKeyBytes, _ := hex.DecodeString(line[2])
		expectedResult, _ := strconv.ParseBool(line[3])

		var derivation KeyDerivation
		copy(derivation[:], derivationBytes)

		var derivedKey PublicKey
		copy(derivedKey[:], derivedKeyBytes)

		actual, actualErr := derivation.underivePublicKey(outputIndex, &derivedKey)

		if expectedResult {
			expectedBytes, _ := hex.DecodeString(line[4])
			var expected PublicKey
			copy(expected[:], expectedBytes)

			assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
			assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))
		} else {
			assert.NotNil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		}

		lineNumber++
		times++
	}

	assert.Equal(t, 288, times)
}
