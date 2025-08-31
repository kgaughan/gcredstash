package command

import (
	"io"
	"os"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func wrapWithDriver(fn func(*cobra.Command, []string, *internal.Driver, io.Writer) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		driver, err := internal.NewDriver(cmd.Context())
		if err != nil {
			return err //nolint:wrapcheck
		}
		return fn(cmd, args, driver, os.Stdout) //nolint:wrapcheck
	}
}
