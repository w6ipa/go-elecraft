package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/w6ipa/go-elecraft/cmd"
)

func main() {
	c := cli.NewCLI("elec", "0.0.0")

	ui := &cli.BasicUi{Writer: os.Stdout, Reader: os.Stdin}

	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"cw": func() (cli.Command, error) {
			return &cmd.CWCmd{
				UI: ui,
			}, nil
		},
		"cw trainer": func() (cli.Command, error) {
			return &cmd.CWTrnCmd{
				UI: ui,
			}, nil
		},

		"cw out": func() (cli.Command, error) {
			return &cmd.CWOutCmd{
				UI: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
