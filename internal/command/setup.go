package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func MakeSetupCmd(driver *internal.Driver, common *CommonFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Setup the credential store",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return driver.CreateDdbTable(common.Table) //nolint:wrapcheck
		},
	}
}
