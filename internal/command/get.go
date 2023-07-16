package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func MakeGetCmd(driver *internal.Driver, common *CommonFlags) *cobra.Command {
	var version string
	var noNL bool
	var noErr bool

	cmd := &cobra.Command{
		Use:   "get credential [context ...]",
		Short: "Get a credential from the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err //nolint:wrapcheck
			}
			return internal.CheckVersion(&version) //nolint:wrapcheck
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := internal.ParseContext(args[1:])
			if err != nil {
				return err //nolint:wrapcheck
			}
			// TOOD: wildcard support
			value, err := driver.GetSecret(args[0], version, common.Table, context)
			if err != nil {
				if noErr {
					return nil
				}
				return err //nolint:wrapcheck
			}
			if noNL {
				cmd.Print(value)
			} else {
				cmd.Println(value)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&noNL, "noline", "n", false, "don't append newline to returned value")
	cmd.Flags().BoolVarP(&noErr, "noerr", "s", false, "don't exit with an error if the secret is not found")
	cmd.Flags().StringVarP(&version, "version", "v", "", "get a specific version of the credential")

	return cmd
}
