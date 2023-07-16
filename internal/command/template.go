package command

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kgaughan/gcredstash/internal"
	gcredstash "github.com/kgaughan/gcredstash/internal"
	"github.com/kgaughan/gcredstash/internal/templating"
	"github.com/spf13/cobra"
)

var inplace bool

func templateImpl(cmd *cobra.Command, args []string) error {
	tmplFile := args[0]

	var content string
	if tmplFile == "-" {
		content = gcredstash.ReadStdin()
	} else {
		var err error
		content, err = gcredstash.ReadFile(tmplFile)
		if err != nil {
			return fmt.Errorf("cannot read %q: %w", tmplFile, err)
		}
	}

	driver, err := internal.NewDriver()
	if err != nil {
		return err //nolint:wrapcheck
	}

	tmpl, err := templating.MakeTemplate(driver, table).Parse(content)
	if err != nil {
		return fmt.Errorf("cannot parse %q template: %w", tmplFile, err)
	}

	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, nil); err != nil {
		return fmt.Errorf("cannot execute %q template: %w", tmplFile, err)
	}

	if inplace {
		if err := os.WriteFile(tmplFile, buf.Bytes(), 0o644); err != nil { //nolint:gosec
			return fmt.Errorf("cannot write to %q: %w", tmplFile, err)
		}
	}

	cmd.Print(buf.String())
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Parse a template file with credentials",
		Args:  cobra.ExactArgs(1),
		RunE:  templateImpl,
	}

	flag := cmd.Flags()
	flag.BoolVarP(&inplace, "inplace", "i", false, "overwrite the template file")

	Root.AddCommand(cmd)
}
