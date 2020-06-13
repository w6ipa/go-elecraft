package main

import (
	"log"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/jc-m/go-elecraft/rig"
	"github.com/jc-m/go-elecraft/ui"
)

func main() {

	k := rig.New("/dev/tty.usbserial-A600UJU4", 38400)

	if err := k.Open(); err != nil {
		log.Fatal(err)
	}

	cmd := rig.NewTTx()

	time.Sleep(1 * time.Second)

	_, err := k.SendCommand(cmd, "1")
	if err != nil {
		log.Fatal(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.CWPracticeLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		log.Panicln(err)
	}

	done := make(chan struct{})

	go ui.BottomUpdate(g, k.GetDataChan(), done)

	if err := g.MainLoop(); err != nil {
		if gocui.IsQuit(err) {
			k.SendCommand(cmd, "0")
			k.Close()
			close(done)
		} else {
			log.Panicln(err)
		}
	}
}
