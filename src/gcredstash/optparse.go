package gcredstash

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrArgumentRequired = errors.New("option requires an argument")
	ErrInvalidContext   = errors.New("invalid context")
	ErrBadVersion       = errors.New("malformed version")
)

func ParseOptionWithValue(args []string, key string) ([]string, string, error) {
	newArgs := []string{}
	val := ""
	nextOpt := false

	for _, arg := range args {
		switch {
		case nextOpt:
			if strings.HasPrefix(arg, "-") {
				return nil, "", fmt.Errorf("%w: %s", ErrArgumentRequired, key)
			}
			val = arg
			nextOpt = false

		case arg == key:
			nextOpt = true

		default:
			newArgs = append(newArgs, arg)
		}
	}

	if nextOpt {
		return nil, "", fmt.Errorf("%w: %s", ErrArgumentRequired, key)
	}

	return newArgs, val, nil
}

func ParseVersion(args []string) ([]string, string, error) {
	newArgs, version, err := ParseOptionWithValue(args, "-v")
	if err != nil {
		//nolint:wrapcheck
		return nil, "", err
	}

	if version != "" {
		ver, err := strconv.Atoi(version)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %q", ErrBadVersion, version)
		}

		version = fmt.Sprintf("%019d", ver)
	}

	return newArgs, version, nil
}

func ParseContext(strs []string) (map[string]string, error) {
	context := map[string]string{}

	for _, ctx := range strs {
		kv := strings.SplitN(ctx, "=", 2)

		if len(kv) < 2 || kv[0] == "" || kv[1] == "" {
			return nil, fmt.Errorf("%w: %s", ErrInvalidContext, ctx)
		}

		context[kv[0]] = kv[1]
	}

	return context, nil
}

func HasOption(args []string, opt string) ([]string, bool) {
	newArgs := []string{}
	hasOpt := false

	for _, arg := range args {
		if arg == opt {
			hasOpt = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs, hasOpt
}
