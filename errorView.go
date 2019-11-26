package main

import "github.com/marcusolsson/tui-go"

type errorView struct {
	*tui.Box
}

func NewErrorView(e error) *errorView {
	title := tui.NewLabel("Ooops something went wrong!")
	body := tui.NewLabel(e.Error())
	body.SetStyleName("warning")
	button := tui.NewButton("[ok, take me home]")
	button.OnActivated(func(b *tui.Button) {
		ui.SetWidget(hView)
		ui.SetFocusChain(hView)
		hView.Update("")
	})
	button.SetFocused(true)
	return &errorView{tui.NewVBox(title, body, button)}
}
