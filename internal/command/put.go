package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/kgaughan/gcredstash/internal"
)

type PutCommand struct {
	Meta
}

func (c *PutCommand) parseArgs(args []string) (string, string, string, map[string]string, bool, error) {
	argsWithoutA, autoVersion := internal.HasOption(args, "-a")
	newArgs, version, err := internal.ParseVersion(argsWithoutA)
	if err != nil {
		//nolint:wrapcheck
		return "", "", "", nil, false, err
	}

	if len(newArgs) < 2 {
		return "", "", "", nil, false, ErrTooFewArgs
	}

	credential := newArgs[0]
	value := newArgs[1]
	context, err := internal.ParseContext(newArgs[2:])

	//nolint:wrapcheck
	return credential, value, version, context, autoVersion, err
}

func (c *PutCommand) RunImpl(args []string) error {
	credential, value, version, context, autoVersion, err := c.parseArgs(args)
	if err != nil {
		return fmt.Errorf("can't parse arguments: %w", err)
	}

	if value == "-" {
		value = internal.ReadStdin()
	}

	if autoVersion {
		latestVersion, err := c.Driver.GetHighestVersion(credential, c.Table)
		if err != nil {
			return fmt.Errorf("can't fetch highest version: %w", err)
		}

		latestVersion++
		version = internal.VersionNumToStr(latestVersion)
	} else if version == "" {
		version = internal.VersionNumToStr(1)
	}

	if err := c.Driver.PutSecret(credential, value, version, c.KmsKey, c.Table, context); err != nil {
		return fmt.Errorf("can't store secret: %w", err)
	}

	fmt.Printf("%s has been stored\n", credential)

	return nil
}

func (c *PutCommand) Run(args []string) int {
	if err := c.RunImpl(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	return 0
}

func (c *PutCommand) Synopsis() string {
	return "Put a credential into the store"
}

func (c *PutCommand) Help() string {
	helpText := `
usage: gcredstash put [-k KEY] [-v VERSION] [-a] credential value [context [context ...]]
`
	return strings.TrimSpace(helpText)
}
