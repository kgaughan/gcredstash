package command

import (
	"fmt"
	"os"
	"strings"

	gcredstash "github.com/kgaughan/gcredstash/internal"
)

type DeleteCommand struct {
	Meta
}

func (c *DeleteCommand) parseArgs(args []string) (string, string, error) {
	newArgs, version, err := gcredstash.ParseVersion(args)
	if err != nil {
		//nolint:wrapcheck
		return "", "", err
	}

	if len(newArgs) < 1 {
		return "", "", ErrTooFewArgs
	}

	if len(newArgs) > 1 {
		return "", "", ErrTooManyArgs
	}

	credential := args[0]

	return credential, version, nil
}

func (c *DeleteCommand) RunImpl(args []string) error {
	credential, version, err := c.parseArgs(args)
	if err != nil {
		return err
	}

	//nolint:wrapcheck
	return c.Driver.DeleteSecrets(credential, version, c.Meta.Table)
}

func (c *DeleteCommand) Run(args []string) int {
	if err := c.RunImpl(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	return 0
}

func (c *DeleteCommand) Synopsis() string {
	return "Delete a credential from the store"
}

func (c *DeleteCommand) Help() string {
	helpText := `
usage: gcredstash delete [-v VERSION] credential
`
	return strings.TrimSpace(helpText)
}
