package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "gcredstash",
	Short: "gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB",
}

var table string

func init() {
	defaultTable := internal.LookupEnvDefault("credential-store", "GCREDSTASH_TABLE", "CREDSTASH_DEFAULT_TABLE")
	Root.PersistentFlags().StringVarP(&table, "table", "t", defaultTable, "DynamoDB table to use for credential storage")
	Root.Version = internal.Version
}
