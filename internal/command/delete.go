package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func MakeDeleteCmd(driver *internal.Driver, common *CommonFlags) *cobra.Command {
	var version string

	cmd := &cobra.Command{
		Use:   "delete credential",
		Short: "Delete a credential from the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
	}
			if err := internal.CheckVersion(&version); err != nil {
		return err
	}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return driver.DeleteSecrets(args[0], version, common.Table)
		},
}
	cmd.Flags().StringVarP(&version, "version", "v", "", "delete a specfic version of the credential")

	return cmd
}
