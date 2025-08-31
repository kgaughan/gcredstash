package command

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/kgaughan/gcredstash/internal"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var inplace bool

func templateImpl(cmd *cobra.Command, args []string, driver *internal.Driver) error {
	ctx := cmd.Context()
	out := cmd.OutOrStdout()

	tmplFile := args[0]

	var content string
	if tmplFile == "-" {
		content = internal.ReadStdin()
	} else {
		var err error
		content, err = internal.ReadFile(tmplFile)
		if err != nil {
			return fmt.Errorf("can't read %q: %w", tmplFile, err)
		}
	}

	tmpl, err := makeTemplate(ctx, driver, table).Parse(content)
	if err != nil {
		return fmt.Errorf("can't parse %q template: %w", tmplFile, err)
	}

	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, nil); err != nil {
		return fmt.Errorf("can't execute %q template: %w", tmplFile, err)
	}

	if inplace {
		if err := os.WriteFile(tmplFile, buf.Bytes(), 0o644); err != nil { //nolint:gosec
			return fmt.Errorf("can't write to %q: %w", tmplFile, err)
		}
	}

	fmt.Fprint(out, buf.String())
	return nil
}

func makeTemplate(ctx context.Context, driver *internal.Driver, table string) *template.Template {
	return template.New("template").Funcs(template.FuncMap{
		"get": func(credential string, encCtxLines ...string) (string, error) {
			encCtx, err := internal.ParseContext(encCtxLines)
			if err != nil {
				return "", fmt.Errorf("could not parse context: %w", err)
			}

			value, err := driver.GetSecret(ctx, credential, "", table, encCtx)
			if err != nil {
				return "", fmt.Errorf("could not fetch credentials: %w", err)
			}

			return value, nil
		},
		"env": func(key string) (string, error) {
			return os.Getenv(key), nil
		},
		"sh": func(line string) (string, error) {
			cmd, err := shellwords.Parse(line)
			if err != nil {
				return "", fmt.Errorf("could not parse command: %w", err)
			}

			var out []byte

			switch len(cmd) {
			case 0:
				out = []byte{}
			case 1:
				//nolint:gosec
				out, err = exec.Command(cmd[0]).Output()
			default:
				//nolint:gosec
				out, err = exec.Command(cmd[0], cmd[1:]...).Output()
			}

			if err != nil {
				return "", fmt.Errorf("could not run command: %w", err)
			}

			str := string(out)

			return strings.TrimRight(str, "\n"), nil
		},
	})
}

func init() {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Parse a template file with credentials",
		Args:  cobra.ExactArgs(1),
		RunE:  wrapWithDriver(templateImpl),
	}

	flag := cmd.Flags()
	flag.BoolVarP(&inplace, "inplace", "i", false, "overwrite the template file")

	Root.AddCommand(cmd)
}
