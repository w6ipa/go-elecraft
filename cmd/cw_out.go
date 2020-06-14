package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/cli"
	"github.com/w6ipa/go-elecraft/rig"
)

type CWOutCmd struct {
	UI cli.Ui
}

// Help return help for cw trainer command
func (c CWOutCmd) Help() string {
	helpText := `
Usage: elec cw out <port> <speed>
  CW redirect to stdout
`
	return strings.TrimSpace(helpText)
}

func (c CWOutCmd) Run(args []string) int {

	f := flag.NewFlagSet("out", flag.ContinueOnError)

	if err := f.Parse(args); err != nil {
		c.UI.Error("Invalid flag")
		return cli.RunResultHelp
	}

	if len(f.Args()) < 2 {
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
	defer end(k)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
Loop:
	for {
		select {
		case data, ok := <-k.GetDataChan():
			if !ok {
				break Loop
			}
			fmt.Fprintf(os.Stdout, string(data))
		case <-sigChan:
			os.Exit(0)
		}
	}
	return 0
}

// Synopsis return help for init command
func (c CWOutCmd) Synopsis() string {
	return "CW to stdout"
}

func end(k *rig.Connection) {
	cmd := rig.NewTTx()
	k.SendCommand(cmd, "0")
	k.Close()
	close(k.GetDataChan())
}
