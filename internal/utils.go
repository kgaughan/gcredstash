package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	VersionFormat = "%019d"
)

func Atoi(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return num
}

func VersionNumToStr(version int) string {
	return fmt.Sprintf(VersionFormat, version)
}

func ReadStdin() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	return strings.TrimRight(string(input), "\n")
}

func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		//nolint:wrapcheck
		return "", err
	}

	return string(content), nil
}

func MapToJSON(m map[string]string) string {
	jsonString, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}

	jsonString = bytes.ReplaceAll(jsonString, []byte("\\u003c"), []byte("<"))
	jsonString = bytes.ReplaceAll(jsonString, []byte("\\u003e"), []byte(">"))
	jsonString = bytes.ReplaceAll(jsonString, []byte("\\u0026"), []byte("&"))

	return string(jsonString)
}

func MaxKeyLen(items map[*string]*string) int {
	max := 0

	for key := range items {
		keyLen := len(*key)

		if keyLen > max {
			max = keyLen
		}
	}

	return max
}
