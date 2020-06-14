package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/mitchellh/cli"
	"github.com/w6ipa/go-elecraft/rig"
	"github.com/w6ipa/go-elecraft/ui"
	"github.com/w6ipa/go-elecraft/utils"
)

type CWTrnCmd struct {
	UI cli.Ui
}

// Help return help for cw trainer command
func (c CWTrnCmd) Help() string {
	helpText := `
Usage: elec cw trainer <port> <speed> <filename>
  CW Trainer
`
	return strings.TrimSpace(helpText)
}

func (c CWTrnCmd) Run(args []string) int {

	f := flag.NewFlagSet("trainer", flag.ContinueOnError)

	if err := f.Parse(args); err != nil {
		c.UI.Error("Invalid flag")
		return cli.RunResultHelp
	}

	if len(f.Args()) < 3 {
		c.UI.Error("Missing arguments")
		return cli.RunResultHelp
	}

	speed, err := strconv.Atoi(f.Arg(1))
	if err != nil {
		c.UI.Error("invalid baud rate")
		return cli.RunResultHelp
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

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			ui.ScrollView(v, -1)
			return nil
		}); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			ui.ScrollView(v, 1)
			return nil
		}); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, y := v.Size()
			ui.ScrollView(v, y-1)
			return nil
		}); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.MoveCursor(1, 0, false)
			return nil
		}); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.MoveCursor(-1, 0, false)
			return nil
		}); err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	content, err := ioutil.ReadFile(f.Arg(2))
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	loadText(g, utils.FilterCW(content))

	done := make(chan struct{})

	go ui.CWUpdate(g, k.GetDataChan(), done)

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

func loadText(g *gocui.Gui, content []byte) {
	g.Update(
		func(g *gocui.Gui) error {
			v, err := g.View("top")
			if err != nil {
				return err
			}
			fmt.Fprint(v, string(content))
			g.Cursor = true
			return nil
		})
}
