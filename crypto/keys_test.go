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
		keyBytes, _ := hex.DecodeString(line[0])
		expected, _ := strconv.ParseBool(line[1])
		pk, _ := KeyFromBytes(keyBytes)

		assert.Equal(t, expected, pk.Check(), fmt.Sprintf("failed at line: %d", lineNumber))

		lineNumber++
		times++
	}

	assert.Equal(t, 372, times)
}
