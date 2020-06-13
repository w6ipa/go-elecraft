package cmd

import (
	"flag"
	"strconv"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/jc-m/go-elecraft/rig"
	"github.com/jc-m/go-elecraft/ui"
	"github.com/mitchellh/cli"
)

type CWTrnCmd struct {
	UI cli.Ui
}

// Help return help for cw trainer command
func (c CWTrnCmd) Help() string {
	helpText := `
Usage: elec cw trainer <port> <speed>
  CW Trainer
`
	return strings.TrimSpace(helpText)
}

func (c CWTrnCmd) Run(args []string) int {

	f := flag.NewFlagSet("trainer", flag.ContinueOnError)

	if err := f.Parse(args); err != nil {
		c.UI.Error("Invalid flag")
		return 1
	}

	if len(f.Args()) < 2 {
		c.UI.Error("Missing arguments")
		return 1
	}

	speed, err := strconv.Atoi(f.Arg(1))
	if err != nil {
		c.UI.Error("invalid baud rate")
		return 1
	}

	k := rig.New(f.Arg(0), speed)

	if err := k.Open(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	cmd := rig.NewTTx()

	time.Sleep(1 * time.Second)

	if _, err := k.SendCommand(cmd, "1"); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	defer g.Close()

	g.SetManagerFunc(ui.CWPracticeLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	done := make(chan struct{})

	go ui.BottomUpdate(g, k.GetDataChan(), done)

	if err := g.MainLoop(); err != nil {
		if gocui.IsQuit(err) {
			k.SendCommand(cmd, "0")
			k.Close()
			close(done)
			return 0
		} else {
			c.UI.Error(err.Error())
			return 1
		}
	}
	return 0
}

// Synopsis return help for init command
func (c CWTrnCmd) Synopsis() string {
	return "CW Trainer"
}
