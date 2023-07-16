package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func deleteImpl(cmd *cobra.Command, args []string) error {
	driver, err := internal.NewDriver()
	if err != nil {
		return err //nolint:wrapcheck
	}
	return driver.DeleteSecrets(args[0], version, table) //nolint:wrapcheck
}

func init() {
	cmd := &cobra.Command{
		Use:   "delete credential",
		Short: "Delete a credential from the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err //nolint:wrapcheck
			}
			return internal.CheckVersion(&version) //nolint:wrapcheck
		},
		RunE: deleteImpl,
	}

	flag := cmd.Flags()
	flag.StringVarP(&version, "version", "v", "", "delete a specfic version of the credential")

	Root.AddCommand(cmd)
}
