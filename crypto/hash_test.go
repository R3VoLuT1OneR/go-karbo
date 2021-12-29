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
