package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the credential store",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		driver, err := internal.NewDriver()
		if err != nil {
			return err //nolint:wrapcheck
		}

		return driver.CreateDdbTable(table) //nolint:wrapcheck
	},
}

func init() {
	Root.AddCommand(setupCmd)
}
