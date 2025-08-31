package command

import (
	"io"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func setupImpl(cmd *cobra.Command, _ []string, driver *internal.Driver, _ io.Writer) error {
	return driver.CreateDdbTable(cmd.Context(), table) //nolint:wrapcheck
}

func init() {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup the credential store",
		Args:  cobra.NoArgs,
		RunE:  wrapWithDriver(setupImpl),
	}

	Root.AddCommand(cmd)
}
