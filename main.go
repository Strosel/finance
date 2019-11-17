package main

import "github.com/marcusolsson/tui-go"

var (
	cView = 0
	views = []tui.Widget{
		GetHistoryView(),
	}
)

func main() {
	ui, err := tui.New(views[cView])
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
