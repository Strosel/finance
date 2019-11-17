package main

import (
	"github.com/marcusolsson/tui-go"
)

type HistoryView struct {
	*tui.Box
	Summary *tui.Box
	History *tui.Box
	Input   *tui.Entry
}

func GetHistoryView() *HistoryView {
	root := HistoryView{}

	root.Summary = tui.NewVBox(tui.NewLabel("          ")) //minimal size padding label
	root.Summary.SetBorder(true)
	root.Summary.SetTitle("Summary")
	root.Summary.SetSizePolicy(tui.Maximum, tui.Maximum)

	root.History = tui.NewHBox(tui.NewLabel("Hist"))
	root.History.SetBorder(true)
	root.History.SetTitle("History")
	root.History.SetSizePolicy(tui.Expanding, tui.Expanding)

	root.Input = tui.NewEntry()
	root.Input.SetFocused(true)
	root.Input.SetSizePolicy(tui.Expanding, tui.Expanding)
	root.Input.OnSubmit(root.Command)

	input := tui.NewHBox(root.Input)
	input.SetBorder(true)
	input.SetSizePolicy(tui.Minimum, tui.Maximum)

	mbox := tui.NewVBox(root.History, input)

	root.Box = tui.NewHBox(root.Summary, mbox)

	return &root
}

func (hv *HistoryView) Command(e *tui.Entry) {
	cmd := e.Text()
	if len(cmd) > 0 {
		e.SetText("")
	}
}
