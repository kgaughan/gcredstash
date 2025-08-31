package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func wrapWithDriver(fn func(*cobra.Command, []string, *internal.Driver) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		driver, err := internal.NewDriver(cmd.Context())
		if err != nil {
			return err //nolint:wrapcheck
		}
		return fn(cmd, args, driver) //nolint:wrapcheck
	}
}
