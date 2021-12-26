package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
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

		assert.Equal(t, expected, sig.Check(&hash, &publicKey), fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 512, times)
}

func TestGenerateRingSignatureThanCheck(t *testing.T) {
	prefixHashBytes, _ := hex.DecodeString("c163cfc7e5c5d7b136155d26b79d91d380c7d7bd7bc9c608ba5ad9f9bb12dc4b")
	imageBytes, _ := hex.DecodeString("c713c67710f09280417aedd46d8c8e65957e7f363c2c1a9929733227e5997836")

	pub1Bytes, _ := hex.DecodeString("a1065e0d19926521b8af3859eacf90abe26336b7e90ad649aaa860884cb16127")
	pub2Bytes, _ := hex.DecodeString("6d63e898708022da25c4cefab9d3da5940e5e9831ef336befc3c8b8712df63f1")

	secretKeyBytes, _ := hex.DecodeString("2126bd4fbb4a1e9f82c176d0ab594ff9cd6cdf13bf9c2416800cd67f3a249405")
	secretIndex := uint64(1)

	var prefixHash Hash
	copy(prefixHash[:], prefixHashBytes[:32])

	var image KeyImage
	copy(image[:], imageBytes[:32])

	var pub1 PublicKey
	copy(pub1[:], pub1Bytes[:32])

	var pub2 PublicKey
	copy(pub2[:], pub2Bytes[:32])

	var sec SecretKey
	copy(sec[:], secretKeyBytes[:32])

	pubs := []PublicKey{pub1, pub2}
	sigs, err := GenerateRingSignature(&prefixHash, &image, &pubs, &sec, secretIndex)

	assert.Nil(t, err)
	assert.NotNil(t, sigs)

	assert.True(t, CheckRingSignature(&prefixHash, &image, &pubs, &sigs, true))
}

func TestGenerateRingSignature(t *testing.T) {
	file, err := os.Open("./fixtures/generate_ring_signature.txt")
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

		prefixHashBytes, _ := hex.DecodeString(line[0])
		imageBytes, _ := hex.DecodeString(line[1])
		pubsCount, _ := strconv.ParseUint(line[2], 10, 64)

		var prefixHash Hash
		copy(prefixHash[:], prefixHashBytes)

		var image KeyImage
		copy(image[:], imageBytes[:])

		pubs := make([]PublicKey, pubsCount)
		for i := 0; i < int(pubsCount); i++ {
			pubBytes, _ := hex.DecodeString(line[i+3])
			copy(pubs[i][:], pubBytes[:32])
		}

		secretKeyBytes, _ := hex.DecodeString(line[3+pubsCount])

		var sec SecretKey
		copy(sec[:], secretKeyBytes)

		secIndex, _ := strconv.ParseUint(line[3+int(pubsCount)+1], 10, 64)

		expectedSigsBytes, _ := hex.DecodeString(line[3+int(pubsCount)+2])

		expectedSigs := make([]Signature, pubsCount)
		_ = binary.Read(bytes.NewReader(expectedSigsBytes), binary.LittleEndian, &expectedSigs)

		actualSigs, actualErr := GenerateRingSignature(&prefixHash, &image, &pubs, &sec, secIndex)

		assert.NotNil(t, actualSigs, fmt.Sprintf("failed at line: %d", lineNumber))
		assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))

		assert.Equal(t, expectedSigs, actualSigs, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 254, times)
}

func TestCheckRingSignature(t *testing.T) {
	file, err := os.Open("./fixtures/check_ring_signature.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	times := 0
	lineNumber := 1
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		prefixHashBytes, _ := hex.DecodeString(line[0])
		imageBytes, _ := hex.DecodeString(line[1])
		pubsCount, _ := strconv.ParseUint(line[2], 10, 64)

		var prefixHash Hash
		copy(prefixHash[:], prefixHashBytes)

		var image KeyImage
		copy(image[:], imageBytes[:])

		pubs := make([]PublicKey, pubsCount)
		for i := 0; i < int(pubsCount); i++ {
			pubBytes, _ := hex.DecodeString(line[i+3])
			copy(pubs[i][:], pubBytes[:32])
		}

		sigsBytes, _ := hex.DecodeString(line[pubsCount+3])

		sigs := make([]Signature, pubsCount)
		_ = binary.Read(bytes.NewReader(sigsBytes), binary.LittleEndian, &sigs)

		expected, _ := strconv.ParseBool(line[pubsCount+4])
		actual := CheckRingSignature(&prefixHash, &image, &pubs, &sigs, true)

		assert.Equal(t, expected, actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 1024, times)
}
