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

		actual, actualErr := GenerateKeyImage(&publicKey, &privateKey)

		assert.Nil(t, actualErr, fmt.Sprintf("failed at line: %d", lineNumber))
		assert.Equal(t, expected, *actual, fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 256, times)
}
