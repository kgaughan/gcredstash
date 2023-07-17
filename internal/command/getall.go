package command

import (
	"fmt"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func getAllImpl(cmd *cobra.Command, args []string, driver *internal.Driver) error {
	context, err := internal.ParseContext(args[0:])
	if err != nil {
		return err //nolint:wrapcheck
	}

	creds := map[string]string{}
	items, err := driver.ListSecrets(table)
	if err != nil {
		return err //nolint:wrapcheck
	}
	for name := range items {
		value, err := driver.GetSecret(*name, "", table, context)
		if err != nil {
			continue
		}
		creds[*name] = value
	}

	jsonString, err := internal.JSONMarshal(creds)
	if err != nil {
		return fmt.Errorf("cannot marshal credentials: %w", err)
	}

	cmd.Print(string(jsonString))
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "getall [context ...]",
		Short: "Get all credentials from the store",
		Args:  cobra.MinimumNArgs(0),
		RunE:  wrapWithDriver(getAllImpl),
	}

	Root.AddCommand(cmd)
}
