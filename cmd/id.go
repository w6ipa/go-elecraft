package cmd

import (
	"flag"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/w6ipa/go-elecraft/rig"
)

type IDCmd struct {
	UI cli.Ui
}

// Help return help for cw command
func (c IDCmd) Help() string {
	helpText := `
Usage: elec id [options] <port>
  Rig Identification
`
	return strings.TrimSpace(helpText)
}

func (c IDCmd) Run(args []string) int {
	var speed int

	f := flag.NewFlagSet("out", flag.ContinueOnError)
	f.IntVar(&speed, "s", 38400, "baud rate")

	if err := f.Parse(args); err != nil {
		return 1
	}
	if len(f.Args()) < 1 {
		c.UI.Error("Missing arguments")
		return cli.RunResultHelp
	}
	k := rig.New(f.Arg(0), speed)

	if err := k.Open(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	defer k.Close()

	time.Sleep(1 * time.Second)
	cmd := rig.NewOM()

	buff, err := k.SendCommand(cmd, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	rsp := cmd.Parse(buff)
	m, ok := rsp.(map[string]bool)
	if !ok {
		c.UI.Error("unexpected structure")
		return 1
	}
	rig := ""
	if _, ok = m["K3S"]; ok {
		rig = "K3S"
	}
	if _, ok = m["KX2"]; ok {
		rig = "KX2"
	}
	if _, ok = m["KX3"]; ok {
		rig = "KX3"
	}

	if len(rig) == 0 {
		return 1
	}
	c.UI.Output(rig)

	return 0
}

// Synopsis return help for init command
func (c IDCmd) Synopsis() string {
	return "Rig identifier"
}
