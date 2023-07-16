package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

func MakeRootCmd(driver *internal.Driver, version string) *cobra.Command {
	common := &CommonFlags{}
	defaultTable := internal.LookupEnvDefault("credential-store", "GCREDSTASH_TABLE", "CREDSTASH_DEFAULT_TABLE")
	cmd := &cobra.Command{
		Use:   "gcredstash",
		Short: "gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB",
	}
	cmd.PersistentFlags().StringVarP(&common.Table, "table", "t", defaultTable, "DynamoDB table to use for credential storage")
	cmd.Version = version
	cmd.AddCommand(
		MakeDeleteCmd(driver, common),
		MakeGetCmd(driver, common),
		MakeGetAllCmd(driver, common),
		MakeListCmd(driver, common),
		MakePutCmd(driver, common),
		MakeSetupCmd(driver, common),
		MakeTemplateCmd(driver, common),
	)
	return cmd
}
