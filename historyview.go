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
		if e.GetType() == "R/"+expand || e.GetType() == expand {
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

func (hv *HistoryView) updateSummary(events []Event, budget Budget) {
	hv.Summary.Clear()

	spent := map[string]int{}
	for _, e := range events {
		switch v := e.(type) {
		case Transaction:
			cat := e.GetCategory()
			if _, ok := budget.Budget[cat]; ok {
				spent[cat] += e.GetSum()
			} else {
				spent["other"] += e.GetSum()
			}
		case Receipt:
			for _, p := range v.Products {
				cat := p.GetCategory()
				if _, ok := budget.Budget[cat]; ok {
					spent[cat] += p.GetSum()
				} else {
					spent["other"] += p.GetSum()
				}
			}
		}
	}

	//TODO print after set order of budget
	//TODO conditional color REMEMBER savings
	for c, s := range budget.Budget {
		hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:\n%-5.2v%5.2v\n", c, float64(s)/100., float64(spent[c])/100.)))
	}
	bt, st := 0, 0
	for k, v := range budget.Budget {
		bt += v
		st += spent[k]
	}
	hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:\n%-5.2v%5.2v\n", "Total", float64(bt)/100., float64(st)/100.)))
}

func (hv *HistoryView) Update(expand string) {
	//events := RandEventStub(40)
	events := GetEvents()
	budget := GetBudget()
	hv.updateHistory(events, expand)
	hv.updateSummary(events, budget)
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
