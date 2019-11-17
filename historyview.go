package main

import (
	"fmt"
	"strings"

	"github.com/marcusolsson/tui-go"
)

type HistoryView struct {
	*tui.Box
	Summary *tui.Box
	History *tui.ScrollArea
	Input   *tui.Entry
}

func GetHistoryView() *HistoryView {
	root := HistoryView{}

	root.Summary = tui.NewVBox(tui.NewLabel("          ")) //minimal size padding label
	root.Summary.SetBorder(true)
	root.Summary.SetTitle("Summary")
	root.Summary.SetSizePolicy(tui.Maximum, tui.Maximum)

	root.LoadHistory()
	hbox := tui.NewVBox(root.History)
	hbox.SetBorder(true)
	hbox.SetTitle("History")
	hbox.SetSizePolicy(tui.Expanding, tui.Expanding)

	root.Input = tui.NewEntry()
	root.Input.SetFocused(true)
	root.Input.SetSizePolicy(tui.Expanding, tui.Expanding)
	root.Input.OnSubmit(root.Command)

	input := tui.NewHBox(root.Input)
	input.SetBorder(true)
	input.SetSizePolicy(tui.Minimum, tui.Maximum)

	mbox := tui.NewVBox(hbox, input)

	root.Box = tui.NewHBox(root.Summary, mbox)

	return &root
}

func (hv *HistoryView) LoadHistory() {
	events := RandEventStub(40)
	history := tui.NewVBox()

	for _, e := range events {
		history.Append(tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("[%v]", e.GetTime().Format("06-01-02 15:04"))),
			tui.NewPadder(5, 0, tui.NewLabel(e.GetName())),
			tui.NewLabel(fmt.Sprintf("%v", e.GetSum())),
			tui.NewPadder(5, 0, tui.NewLabel(e.GetType())),
		))
	}

	hv.History = tui.NewScrollArea(history)
}

func (hv *HistoryView) Command(e *tui.Entry) {
	cmd := strings.ToLower(e.Text())
	if cmd == "top" {
		hv.History.ScrollToTop()
	}
	e.SetText("")
}
