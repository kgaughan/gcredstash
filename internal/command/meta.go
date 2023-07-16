package command

import (
	"github.com/kgaughan/gcredstash/internal"
	"github.com/mitchellh/cli"
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	UI      cli.Ui
	Table   string
	KmsKey  string
	Version string
	Driver  *internal.Driver
}
