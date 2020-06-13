package cmd

import (
	"strings"

	"github.com/mitchellh/cli"
)

type CWCmd struct {
	UI cli.Ui
}

// Help return help for cw command
func (c CWCmd) Help() string {
	helpText := `
Usage: elec cw <subcommand> [options] [args]
  CW operations
`
	return strings.TrimSpace(helpText)
}

func (c CWCmd) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis return help for init command
func (c CWCmd) Synopsis() string {
	return "CW Operations"
}
