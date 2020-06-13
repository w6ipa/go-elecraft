package ui

import (
	"fmt"
	"math"

	"github.com/awesome-gocui/gocui"
)

func CWPracticeLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	maxTopY := 0 + math.Round(float64(maxY/2))
	maxBottomY := maxTopY + math.Round(float64(maxY/2))

	if v, err := g.SetView("top", 1, 1, maxX-1, int(maxTopY), 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Wrap = true
		if _, err := g.SetCurrentView("top"); err != nil {
			return err
		}
		g.CurrentView().Title = "Active"
	}
	if v, err := g.SetView("bottom", 1, int(maxTopY+1), maxX-1, int(maxBottomY-1), 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func BottomUpdate(g *gocui.Gui, c chan []byte, done chan struct{}) {
Loop:
	for {
		select {
		case <-done:
			return
		case data, ok := <-c:
			if !ok {
				break Loop
			}
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("bottom")
				if err != nil {
					return err
				}
				fmt.Fprintf(v, "%s", data)
				return nil
			})
		}
	}
	return
}
