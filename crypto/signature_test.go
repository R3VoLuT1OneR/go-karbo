package crypto

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGenerateSignature(t *testing.T) {
	t.SkipNow()
	file, err := os.Open("./fixtures/signature_generate.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		assert.Equal(t, 4, len(line))

		hashBytes, _ := hex.DecodeString(line[0])
		pkBytes, _ := hex.DecodeString(line[1])
		skBytes, _ := hex.DecodeString(line[2])
		sigBytes, _ := hex.DecodeString(line[3])

		hash := HashFromBytes(hashBytes)
		pk, _ := KeyFromBytes(pkBytes)
		sk, _ := KeyFromBytes(skBytes)

		var expectedSignature Signature
		_ = expectedSignature.Deserialize(bytes.NewReader(sigBytes))

		sig, err := GenerateSignature(hash, pk, sk)

		assert.Nil(t, err)
		assert.Equal(t, expectedSignature, *sig)
	}

	assert.Equal(t, 256, times)
}

func TestSignature_Check(t *testing.T) {
	file, err := os.Open("./fixtures/signature_check.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		assert.Equal(t, 4, len(line))

		hashBytes, _ := hex.DecodeString(line[0])
		pkBytes, _ := hex.DecodeString(line[1])
		sigBytes, _ := hex.DecodeString(line[2])
		expected, _ := strconv.ParseBool(line[3])

		hash := HashFromBytes(hashBytes)
		pk, _ := KeyFromBytes(pkBytes)

		var sig Signature
		_ = sig.Deserialize(bytes.NewReader(sigBytes))

		assert.Equal(t, expected, sig.Check(hash, pk))
	}

	assert.Equal(t, 512, times)
}
