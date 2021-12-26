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

var testingData = []*struct {
	spendKey       string
	viewKey        string
	publicSpendKey string
	publicViewKey  string
}{
	{
		"6390482f5b3a1fe7fef34577b2cd0d14f12c075578e21ecaf48d1fbc300cf80b",
		"2d92d42406f972c51bce29af0d7ece15284c3decc8d15afa9d72ac76e0d07508",
		"711e2156025f8b8d66aeb2908e21a08f971a3b3b722de0e0876b68bcf0c71b74",
		"ba8e26760a9262408f4cf67cf0b5f4c3e69a8a07367b77149dac04834b300f29",
	},
}

func TestGenerateViewFromSpend(t *testing.T) {
	for _, td := range testingData {
		spendKeyBytes, _ := hex.DecodeString(td.spendKey)
		viewKeyBytes, _ := hex.DecodeString(td.viewKey)

		var spendKey SecretKey
		copy(spendKey[:], spendKeyBytes[:])

		var viewKey SecretKey
		copy(viewKey[:], viewKeyBytes[:])

		assert.Equal(t, viewKey, ViewFromSpend(&spendKey))
	}
}

func TestGetPublicKey(t *testing.T) {
	for _, td := range testingData {
		spendKeyBytes, _ := hex.DecodeString(td.spendKey)
		viewKeyBytes, _ := hex.DecodeString(td.viewKey)

		var spendKey SecretKey
		copy(spendKey[:], spendKeyBytes[:])

		var viewKey SecretKey
		copy(viewKey[:], viewKeyBytes[:])

		publicSpendKey, _ := PublicFromSecret(&spendKey)
		publicViewKey, _ := PublicFromSecret(&viewKey)

		assert.Equal(t, td.publicSpendKey, hex.EncodeToString(publicSpendKey[:]))
		assert.Equal(t, td.publicViewKey, hex.EncodeToString(publicViewKey[:]))
	}
}

func TestKey_Check(t *testing.T) {
	file, err := os.Open("./fixtures/key_check.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		publicKeyBytes, _ := hex.DecodeString(line[0])
		expected, _ := strconv.ParseBool(line[1])

		var publicKey PublicKey
		copy(publicKey[:], publicKeyBytes[:32])

		assert.Equal(t, expected, publicKey.Check(), fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 372, times)
}

func TestPublicFromPrivate(t *testing.T) {
	file, err := os.Open("./fixtures/public_from_private.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		secretKeyBytes, _ := hex.DecodeString(line[0])
		expectedSuccess, _ := strconv.ParseBool(line[1])

		var secretKey SecretKey
		copy(secretKey[:], secretKeyBytes)

		actualPublicKey, err := PublicFromSecret(&secretKey)

		assert.Equal(t, expectedSuccess, err == nil, fmt.Sprintf("failed at line: %d", lineNumber))

		if len(line) > 2 {
			expectedBytes, _ := hex.DecodeString(line[2])

			var expectedPublicKey PublicKey
			copy(expectedPublicKey[:], expectedBytes[:32])

			assert.Equal(t, expectedPublicKey, *actualPublicKey, fmt.Sprintf("failed at line: %d", lineNumber))
		}

		lineNumber++
		times++
	}

	assert.Equal(t, 272, times)
}
