package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func setupImpl(cmd *cobra.Command, args []string, driver *internal.Driver) error {
	return driver.CreateDdbTable(table) //nolint:wrapcheck
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
