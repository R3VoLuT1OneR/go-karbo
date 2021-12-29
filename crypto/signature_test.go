package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testReader struct {
	seed []byte
}

func (r *testReader) Read(b []byte) (int, error) {
	copy(b, r.seed)
	return len(r.seed), nil
}

func TestSignatureGenerateAndThenCheck(t *testing.T) {
	hashKey, _ := GenerateKey()
	hash := HashFromBytes(hashKey[:])

	secretKey, _ := GenerateKey()
	publicKey, _ := PublicFromSecret(&secretKey)

	sig, _ := hash.Sign(&secretKey)

	assert.True(t, sig.Check(&hash, publicKey))
}

func TestGenerateSignature(t *testing.T) {
	file, err := os.Open("./fixtures/signature_generate.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Replace random
	rand.Reader = &testReader{[]byte{1, 2, 3}}
	saveReader := rand.Reader
	defer func(reader io.Reader) {
		rand.Reader = reader
	}(saveReader)

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		assert.Equal(t, 3, len(line))

		hashBytes, _ := hex.DecodeString(line[0])
		secretKeyBytes, _ := hex.DecodeString(line[1])
		sigBytes, _ := hex.DecodeString(line[2])

		var hash Hash
		copy(hash[:], hashBytes)

		var secretKey SecretKey
		copy(secretKey[:], secretKeyBytes)

		var expectedSignature Signature
		_ = expectedSignature.Deserialize(bytes.NewReader(sigBytes))

		sig, _ := hash.Sign(&secretKey)

		assert.Equal(t, expectedSignature, *sig, fmt.Sprintf("failed on line: %d", lineNumber))
		lineNumber++
		times++
	}

	assert.Equal(t, 253, times)
}

func TestSignature_Check(t *testing.T) {
	file, err := os.Open("./fixtures/signature_check.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		assert.Equal(t, 4, len(line))

		hashBytes, _ := hex.DecodeString(line[0])
		publicKeyBytes, _ := hex.DecodeString(line[1])
		sigBytes, _ := hex.DecodeString(line[2])
		expected, _ := strconv.ParseBool(line[3])

		var hash Hash
		copy(hash[:], hashBytes)

		var publicKey PublicKey
		copy(publicKey[:], publicKeyBytes[:])

		var sig Signature
		_ = sig.Deserialize(bytes.NewReader(sigBytes))

		if lineNumber == 4 {
			assert.Equal(t, expected, sig.Check(&hash, &publicKey), fmt.Sprintf("failed at line: %d", lineNumber))
		}

		lineNumber++
		times++
	}

	assert.Equal(t, 512, times)
}
