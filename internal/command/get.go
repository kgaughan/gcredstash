package command

import (
	"fmt"
	"io"
	"strings"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
)

var (
	noNL  bool
	noErr bool
)

func getImpl(cmd *cobra.Command, args []string, driver *internal.Driver, out io.Writer) error {
	ctx := cmd.Context()

	encCtx, err := internal.ParseContext(args[1:])
	if err != nil {
		return err //nolint:wrapcheck
	}

	credential := args[0]
	if strings.Contains(credential, "*") {
		items, err := driver.ListSecrets(ctx, table)
		if err != nil {
			return err //nolint:wrapcheck
		}
		creds := map[string]string{}
		for name := range items {
			if !glob.Glob(credential, name) {
				continue
			}
			value, err := driver.GetSecret(ctx, name, version, table, encCtx)
			if err != nil {
				continue
			}
			creds[name] = value
		}
		result, err := internal.JSONMarshal(creds)
		if err != nil {
			return fmt.Errorf("can't marshal credential: %w", err)
		}
		fmt.Fprint(out, string(result))
	} else {
		value, err := driver.GetSecret(ctx, credential, version, table, encCtx)
		if err != nil {
			if noErr {
				return nil
			}
			return err //nolint:wrapcheck
		}
		if noNL {
			fmt.Fprint(out, value)
		} else {
			fmt.Fprintln(out, value)
		}
	}
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "get credential [context ...]",
		Short: "Get a credential from the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err //nolint:wrapcheck
			}
			return internal.CheckVersion(&version) //nolint:wrapcheck
		},
		RunE: wrapWithDriver(getImpl),
	}

	flag := cmd.Flags()
	flag.BoolVarP(&noNL, "noline", "n", false, "don't append newline to returned value")
	flag.BoolVarP(&noErr, "noerr", "s", false, "don't exit with an error if the secret is not found")
	flag.StringVarP(&version, "version", "v", "", "get a specific version of the credential")

	Root.AddCommand(cmd)
}
