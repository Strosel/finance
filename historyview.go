package main

import (
	"fmt"
	"strings"

	"github.com/marcusolsson/tui-go"
)

const savestr = "spar|spara|sparande|save|saving|savings"

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
	root.Summary.Box.Append(tui.NewLabel(fmt.Sprintf("%18v", "")))

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
			if _, ok := budget.Spending[cat]; ok {
				spent[cat] += e.GetSum()
			} else {
				spent["other"] += e.GetSum()
			}
		case Receipt:
			for _, p := range v.Products {
				cat := p.GetCategory()
				if _, ok := budget.Spending[cat]; ok {
					spent[cat] += p.GetSum()
				} else {
					spent["other"] += p.GetSum()
				}
			}
		}
	}

	hv.Summary.Append(tui.NewLabel("Income"))
	for c, s := range budget.Income {
		hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:", c)))
		hv.Summary.Append(SummaryFormat(-1, s.Sum, !s.Received))
	}

	hv.Summary.Append(tui.NewLabel(" "))
	hv.Summary.Append(tui.NewLabel("Spending"))
	for c, s := range budget.Spending {
		hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:", c)))

		if strings.Contains(savestr, strings.ToLower(c)) {
			hv.Summary.Append(SummaryFormat(s, spent[c], spent[c] < s))
		} else {
			hv.Summary.Append(SummaryFormat(s, spent[c], spent[c] > s))
		}
	}
	bt, st := 0, 0
	for k, v := range budget.Spending {
		bt += v
		st += spent[k]
	}

	hv.Summary.Append(tui.NewLabel(" ")) //?
	hv.Summary.Append(tui.NewLabel("Total:"))
	hv.Summary.Append(SummaryFormat(bt, st, st > bt))

	inc := 0
	for _, v := range budget.Income {
		if v.Received {
			inc += v.Sum
		}
	}
	hv.Summary.Append(tui.NewLabel(" ")) //?
	hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:\n%8.2f %8.2f", "Balance", float64(inc-bt)/100., float64(inc-st)/100.)))
}

func (hv *HistoryView) Update(expand string) {
	//events := RandEventStub(40)
	events := GetEvents()
	budget := GetBudget()
	hv.updateHistory(events, expand)
	hv.updateSummary(events, budget)
}

func (hv *HistoryView) FocusNext(w tui.Widget) tui.Widget {
	switch w {
	case hv.Input:
		return hv.Summary
	case hv.Summary:
		return hv.History
	case hv.History:
		return hv.Input
	default:
		return nil
	}
}

func (hv *HistoryView) FocusPrev(w tui.Widget) tui.Widget {
	switch w {
	case hv.Input:
		return hv.History
	case hv.Summary:
		return hv.Input
	case hv.History:
		return hv.Summary
	default:
		return nil
	}
}

func (hv *HistoryView) FocusDefault() tui.Widget {
	return hv.Input
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
	case "add":
		if len(cmd) > 1 {
			if strings.ToLower(cmd[1])[0] == 'r' {
				aView := NewAddRView(nil)
				ui.SetWidget(aView)
				ui.SetFocusChain(aView)
			} else if strings.ToLower(cmd[1])[0] == 't' {
				aView := NewAddTView(nil, nil)
				ui.SetWidget(aView)
				ui.SetFocusChain(aView)
			}
		}
	case "set":
		sView := NewSetBView(nil)
		ui.SetWidget(sView)
		ui.SetFocusChain(sView)
		//todo check for existing in given timeframe, use that
	}

	e.SetText("")
}
