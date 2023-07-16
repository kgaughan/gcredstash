package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var errBadVersion = errors.New("malformed version")

func checkVersion(version *string) error {
	if *version != "" {
		parsed, err := strconv.Atoi(*version)
		if err != nil {
			return fmt.Errorf("%w: %q", errBadVersion, *version)
		}
		*version = fmt.Sprintf("%019d", parsed)
	}
	return nil
}

func lookupEnvDefault(defaultVal string, envVars ...string) string {
	for _, envVar := range envVars {
		if val, ok := os.LookupEnv(envVar); ok {
			return val
		}
	}
	return defaultVal
}

func main() {
	var secretVersion string
	var key string
	var noNL bool
	var noErr bool
	var autoVersion bool
	var inplace bool
	var table string

	defaultTable := lookupEnvDefault("credential-store", "GCREDSTASH_TABLE", "CREDSTASH_DEFAULT_TABLE")
	defaultKMSKey := lookupEnvDefault("alias/credstash", "GCREDSTASH_KMS_KEY", "CREDSTASH_KMS_KEY")

	var deleteCmd = &cobra.Command{
		Use:   "delete credential",
		Short: "Delete a credential from the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			if err := checkVersion(&secretVersion); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// TODO
			fmt.Printf("%+v, %q\n", args, secretVersion)
		},
	}
	deleteCmd.Flags().StringVarP(&secretVersion, "version", "v", "", "delete a specfic version of the credential")

	var getCmd = &cobra.Command{
		Use:   "get credential [context ...]",
		Short: "Get a credential from the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}
			if err := checkVersion(&secretVersion); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	getCmd.Flags().BoolVarP(&noNL, "noline", "n", false, "don't append newline to returned value")
	getCmd.Flags().BoolVarP(&noErr, "noerr", "s", false, "don't exit with an error if the secret is not found")
	getCmd.Flags().StringVarP(&secretVersion, "version", "v", "", "get a specific version of the credential")

	var getAllCmd = &cobra.Command{
		Use:   "getall [context ...]",
		Short: "Get all credentials from the store",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List credentials and their version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var putCmd = &cobra.Command{
		Use:   "put credential value [context ...]",
		Short: "Put a credential into the store",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
				return err
			}
			if err := checkVersion(&secretVersion); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	putCmd.Flags().StringVarP(&key, "key", "k", defaultKMSKey, "the KMS key-id of the master key to use")
	putCmd.Flags().BoolVarP(&autoVersion, "autoversion", "a", false, "automatically increment the version of the credential to be stored; causes the -v flag to be ignored")
	putCmd.Flags().StringVarP(&secretVersion, "version", "v", "", "put a specific version of the credential")

	var setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup the credential store",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var templateCmd = &cobra.Command{
		Use:   "template",
		Short: "Parse a tempalte file with credentials",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	templateCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "overwrite the template file")

	var rootCmd = &cobra.Command{
		Use:   "gcredstash",
		Short: "gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB",
	}
	rootCmd.PersistentFlags().StringVarP(&table, "table", "t", defaultTable, "DynamoDB table to use for credential storage")

	rootCmd.AddCommand(
		deleteCmd,
		getCmd,
		getAllCmd,
		listCmd,
		putCmd,
		setupCmd,
		templateCmd,
	)
	rootCmd.Execute()
}
