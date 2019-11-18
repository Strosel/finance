package main

import (
	"fmt"
	"strings"

	"github.com/marcusolsson/tui-go"
)

type HistoryView struct {
	*tui.Box
	Summary *ScrollBox
	History *ScrollBox
	Input   *tui.Entry
}

func GetHistoryView() *HistoryView {
	root := HistoryView{}

	root.Summary = NewScrollBox()
	root.Summary.SetBorder(true)
	root.Summary.SetTitle("Summary")
	root.Summary.SetSizePolicy(tui.Maximum, tui.Maximum)

	root.History = NewScrollBox()
	root.Update("")
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

func (hv *HistoryView) updateHistory(events []Event, expand string) {
	hv.History.Clear()

	for _, e := range events {
		hv.History.Append(tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("[%v]", e.GetTime().Format("06-01-02 15:04"))),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("%10v", e.GetName()))),
			tui.NewLabel(fmt.Sprintf("%5v", e.GetSum())),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("%10v", e.GetCategory()))),
			tui.NewLabel(e.GetType()),
		))
		if e.GetType() == "R/"+expand {
			r := e.(Receipt)
			for _, p := range r.Products {
				hv.History.Append(tui.NewHBox(
					tui.NewLabel("                "),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("%10v", p.GetName()))),
					tui.NewLabel(fmt.Sprintf("%5v", p.GetSum())),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("%10v", p.GetCategory()))),
					tui.NewLabel(p.GetType()),
				))
			}
		}
	}
}

func (hv *HistoryView) updateSummary(events []Event) {
	hv.Summary.Clear()

	budget := map[string]int{}
	spent := map[string]int{}
	for _, e := range events {
		switch v := e.(type) {
		case Transaction:
			spent[e.GetCategory()] += e.GetSum()
			budget[e.GetCategory()] = 2000
		case Receipt:
			for _, p := range v.Products {
				spent[p.GetCategory()] += e.GetSum()
				budget[p.GetCategory()] = 2000
			}
		}
	}

	//TODO print after set order of budget
	//TODO conditional color
	for c, s := range spent {
		hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:\n%-6v%6v\n", c, budget[c], s)))
	}
}

func (hv *HistoryView) Update(expand string) {
	//events := RandEventStub(40)
	events := GetEvents()
	fmt.Println(len(events))
	hv.updateHistory(events, expand)
	hv.updateSummary(events)
}

func (hv *HistoryView) Command(e *tui.Entry) {
	cmd := strings.Split(e.Text(), " ")
	switch strings.ToLower(cmd[0]) {
	case "update":
		hv.Update("")
		fallthrough
	case "expand":
		if len(cmd) > 1 {
			hv.Update(cmd[1])
		}
		fallthrough
	case "top":
		hv.History.ScrollToTop()
	}
	e.SetText("")
}
