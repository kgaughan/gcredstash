package command

import (
	"fmt"
	"io"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

var (
	key         string
	autoVersion bool
)

func putImpl(_ *cobra.Command, args []string, driver *internal.Driver, out io.Writer) error {
	context, err := internal.ParseContext(args[2:])
	if err != nil {
		return err //nolint:wrapcheck
	}
	credential := args[0]
	value := args[1]
	if value == "-" {
		value = internal.ReadStdin()
	}

	if autoVersion {
		latestVersion, err := driver.GetHighestVersion(credential, table)
		if err != nil {
			return fmt.Errorf("can't fetch highest version: %w", err)
		}
		latestVersion++
		version = internal.VersionNumToStr(latestVersion)
	} else if version == "" {
		version = internal.VersionNumToStr(1)
	}

	if err := driver.PutSecret(credential, value, version, key, table, context); err != nil {
		return fmt.Errorf("can't store secret: %w", err)
	}

	fmt.Fprintf(out, "%v has been stored\n", credential)
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "put credential value [context ...]",
		Short: "Put a credential into the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
				return err
			}
			return internal.CheckVersion(&version) //nolint:wrapcheck
		},
		RunE: wrapWithDriver(putImpl),
	}

	defaultKMSKey := internal.LookupEnvDefault("alias/credstash", "GCREDSTASH_KMS_KEY", "CREDSTASH_KMS_KEY")

	flag := cmd.Flags()
	flag.StringVarP(&key, "key", "k", defaultKMSKey, "the KMS key-id of the master key to use")
	flag.BoolVarP(&autoVersion, "autoversion", "a", false, "automatically increment the version of the credential to be stored; causes the -v flag to be ignored")
	flag.StringVarP(&version, "version", "v", "", "put a specific version of the credential")

	Root.AddCommand(cmd)
}
