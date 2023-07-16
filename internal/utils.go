package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidContext = errors.New("invalid context")
	ErrBadVersion     = errors.New("malformed version")
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

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
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

func LookupEnvDefault(defaultVal string, envVars ...string) string {
	for _, envVar := range envVars {
		if val, ok := os.LookupEnv(envVar); ok {
			return val
		}
	}
	return defaultVal
}

func CheckVersion(version *string) error {
	if *version != "" {
		parsed, err := strconv.Atoi(*version)
		if err != nil {
			return fmt.Errorf("%w: %q", ErrBadVersion, *version)
		}
		*version = fmt.Sprintf("%019d", parsed)
	}
	return nil
}

func ParseContext(strs []string) (map[string]string, error) {
	context := map[string]string{}

	for _, ctx := range strs {
		kv := strings.SplitN(ctx, "=", 2)

		if len(kv) < 2 || kv[0] == "" || kv[1] == "" {
			return nil, fmt.Errorf("%w: %q", ErrInvalidContext, ctx)
		}

		context[kv[0]] = kv[1]
	}

	return context, nil
}
