package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"

	"io"
	"os"
)

func wrapWithDriver(fn func(*cobra.Command, []string, *internal.Driver, io.Writer) error) func(*cobra.Command, []string) error {
	driver, err := internal.NewDriver()
	if err != nil {
		return func(_ *cobra.Command, _ []string) error {
			return err //nolint:wrapcheck
		}
	}
	return func(cmd *cobra.Command, args []string) error {
		return fn(cmd, args, driver, os.Stdout) //nolint:wrapcheck
	}
}
