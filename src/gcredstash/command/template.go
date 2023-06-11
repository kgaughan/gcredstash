package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/kgaughan/gcredstash/src/gcredstash"
	"github.com/mattn/go-shellwords"
)

var ErrCannotCast = errors.New("cannot cast to string")

type TemplateCommand struct {
	Meta
}

func (c *TemplateCommand) parseArgs(args []string) (string, bool, error) {
	newArgs, inPlace := gcredstash.HasOption(args, "-i")

	if len(newArgs) < 1 {
		return "", false, ErrTooFewArgs
	}

	if len(newArgs) > 1 {
		return "", false, ErrTooManyArgs
	}

	tmplFile := newArgs[0]

	return tmplFile, inPlace, nil
}

func (c *TemplateCommand) readTemplate(filename string) (string, error) {
	var content string

	if filename == "-" {
		content = gcredstash.ReadStdin()
	} else {
		var err error
		content, err = gcredstash.ReadFile(filename)

		if err != nil {
			return "", fmt.Errorf("can't read template: %w", err)
		}
	}

	return content, nil
}

func (c *TemplateCommand) getCredential(credential string, context map[string]string) (string, error) {
	value, err := c.Driver.GetSecret(credential, "", c.Table, context)
	if err != nil {
		//nolint:wrapcheck
		return "", err
	}

	return value, nil
}

func (c *TemplateCommand) executeTemplate(name, content string) (string, error) {
	tmpl := template.New(name)

	tmpl = tmpl.Funcs(template.FuncMap{
		"get": func(args ...interface{}) (string, error) {
			if len(args) < 1 {
				return "", ErrTooFewArgs
			}

			newArgs := []string{}

			for _, arg := range args {
				str, ok := arg.(string)

				if !ok {
					return "", fmt.Errorf("%w: %v", ErrCannotCast, arg)
				}

				newArgs = append(newArgs, str)
			}

			credential := newArgs[0]
			context, err := gcredstash.ParseContext(newArgs[1:])
			if err != nil {
				return "", fmt.Errorf("could not parse context: %w", err)
			}

			value, err := c.getCredential(credential, context)
			if err != nil {
				return "", fmt.Errorf("could not fetch credentials: %w", err)
			}

			return value, nil
		},
		"env": func(args ...interface{}) (string, error) {
			if len(args) < 1 {
				return "", ErrTooFewArgs
			}

			if len(args) > 1 {
				return "", ErrTooManyArgs
			}

			key, ok := args[0].(string)

			if !ok {
				return "", fmt.Errorf("%w: %v", ErrCannotCast, args[0])
			}

			return os.Getenv(key), nil
		},
		"sh": func(args ...interface{}) (string, error) {
			if len(args) < 1 {
				return "", ErrTooFewArgs
			}

			if len(args) > 1 {
				return "", ErrTooManyArgs
			}

			line, ok := args[0].(string)

			if !ok {
				return "", fmt.Errorf("%w: %v", ErrCannotCast, args[0])
			}

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

	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return "", fmt.Errorf("can't parse %q template: %w", name, err)
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)

	return buf.String(), err
}

func (c *TemplateCommand) RunImpl(args []string) (string, error) {
	tmplFile, inPlace, err := c.parseArgs(args)
	if err != nil {
		return "", fmt.Errorf("can't parse arguments: %w", err)
	}

	tmplContent, err := c.readTemplate(tmplFile)
	if err != nil {
		//nolint:wrapcheck
		return "", err
	}

	out, err := c.executeTemplate(tmplFile, tmplContent)
	if err != nil {
		return "", fmt.Errorf("can't execute template: %w", err)
	}

	if inPlace {
		//nolint:gosec
		if err := os.WriteFile(tmplFile, []byte(out), 0o644); err != nil {
			return "", fmt.Errorf("could not write output: %w", err)
		}
	}

	return out, nil
}

func (c *TemplateCommand) Run(args []string) int {
	out, err := c.RunImpl(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return 1
	}

	fmt.Print(out)

	return 0
}

func (c *TemplateCommand) Synopsis() string {
	return "Parse a template file with credentials"
}

func (c *TemplateCommand) Help() string {
	helpText := `
usage: gcredstash template [-i] template_file
`
	return strings.TrimSpace(helpText)
}
