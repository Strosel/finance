package main

import "github.com/marcusolsson/tui-go"

var (
	hView = GetHistoryView()
)

func main() {
	ui, err := tui.New(hView)
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Up", func() { hView.History.Scroll(0, -1) })
	ui.SetKeybinding("Down", func() { hView.History.Scroll(0, 1) })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
