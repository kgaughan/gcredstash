package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/kgaughan/gcredstash/src/gcredstash"
)

type DeleteCommand struct {
	Meta
}

func (c *DeleteCommand) parseArgs(args []string) (string, string, error) {
	newArgs, version, err := gcredstash.ParseVersion(args)
	if err != nil {
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

	err = c.Driver.DeleteSecrets(credential, version, c.Meta.Table)

	if err != nil {
		return err
	}

	return nil
}

func (c *DeleteCommand) Run(args []string) int {
	err := c.RunImpl(args)
	if err != nil {
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
