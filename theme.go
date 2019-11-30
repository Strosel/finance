package main

import (
	"fmt"

	"github.com/marcusolsson/tui-go"
)

func GetTheme() *tui.Theme {
	t := tui.NewTheme()

	t.SetStyle("normal", tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite})
	t.SetStyle("label.warning", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})
	t.SetStyle("label.good", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorGreen})
	t.SetStyle("button.focused", tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack})
	t.SetStyle("list.item.selected", tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack})

	return t
}

func SummaryFormat(plan, fact int, warn bool) *tui.Box {
	slbl := tui.NewLabel(fmt.Sprintf("%9.2f", float64(fact)/100.))
	if warn {
		slbl.SetStyleName("warning")
	} else {
		slbl.SetStyleName("good")
	}

	if plan >= 0 {
		return tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%9.2f", float64(plan)/100.)),
			slbl,
		)
	}
	return tui.NewHBox(
		tui.NewLabel(fmt.Sprintf("%9v", "")),
		slbl,
	)
}

func SavingsFormat(plan, fact int) *tui.Box {
	think := tui.NewLabel(fmt.Sprintf("%9.2f", float64(plan)/100.))
	if plan < 0 {
		think.SetStyleName("warning")
	}

	real := tui.NewLabel(fmt.Sprintf("%9.2f", float64(fact)/100.))
	if fact < plan {
		real.SetStyleName("warning")
	} else {
		real.SetStyleName("good")
	}

	return tui.NewHBox(
		think,
		real,
	)
}
