package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/jc-m/go-elecraft/ui"
	"github.com/jc-m/go-elecraft/utils"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func bottomUpdate(g *gocui.Gui, done chan struct{}) {
	charset := "ABCDEGHIJKLMNOPQRSTUVWXYZ1234567890 ,."
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("bottom")
				if err != nil {
					return err
				}
				fmt.Fprintf(v, "%s", string(charset[seededRand.Intn(len(charset))]))
				return nil
			})
		}
	}
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(ui.CWPracticeLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.MoveCursor(1, 0, false)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			x, y := v.Cursor()
			w, err := v.Word(x, y)
			if err != nil {
				return nil
			}
			b, _ := g.View("bottom")
			fmt.Fprint(b, w)
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			g.CurrentView().Title = ""
			if g.CurrentView().Name() == "top" {
				g.SetCurrentView("bottom")
			} else {
				g.SetCurrentView("top")
			}
			g.CurrentView().Title = "Active"
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	done := make(chan struct{})

	content, err := ioutil.ReadFile("pg5200.txt")
	if err != nil {
		log.Fatal(err)
	}

	loadText(g, utils.FilterCW(content))
	//go bottomUpdate(g, done)

	if err := g.MainLoop(); err != nil {
		if gocui.IsQuit(err) {
			close(done)
		} else {
			log.Panicln(err)
		}
	}
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
