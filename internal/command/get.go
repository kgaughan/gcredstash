package command

import (
	"fmt"
	"strings"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
)

var (
	noNL  bool
	noErr bool
)

var getCmd = &cobra.Command{
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

		driver, err := internal.NewDriver()
		if err != nil {
			return err //nolint:wrapcheck
		}

		credential := args[0]
		if strings.Contains(credential, "*") {
			items, err := driver.ListSecrets(table)
			if err != nil {
				return err //nolint:wrapcheck
			}
			creds := map[string]string{}
			for name := range items {
				if !glob.Glob(credential, *name) {
					continue
				}
				value, err := driver.GetSecret(*name, version, table, context)
				if err != nil {
					continue
				}
				creds[*name] = value
			}
			result, err := internal.JSONMarshal(creds)
			if err != nil {
				return fmt.Errorf("cannot marshal credential: %w", err)
			}
			cmd.Println(string(result))
		} else {
			// TOOD: wildcard support
			value, err := driver.GetSecret(credential, version, table, context)
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
		}
		return nil
	},
}

func init() {
	flag := getCmd.Flags()
	flag.BoolVarP(&noNL, "noline", "n", false, "don't append newline to returned value")
	flag.BoolVarP(&noErr, "noerr", "s", false, "don't exit with an error if the secret is not found")
	flag.StringVarP(&version, "version", "v", "", "get a specific version of the credential")
	Root.AddCommand(getCmd)
}
