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
