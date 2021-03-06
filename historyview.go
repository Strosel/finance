package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/marcusolsson/tui-go"
	"github.com/strosel/finance/finance"
	"go.mongodb.org/mongo-driver/bson"
)

type HistoryView struct {
	*tui.Box
	Summary *ScrollBox
	History *ScrollList
	Input   *tui.Entry

	time time.Time
}

func GetHistoryView() *HistoryView {
	root := HistoryView{
		time: time.Now(),
	}

	root.Summary = NewScrollBox()
	root.Summary.SetBorder(true)
	root.Summary.SetTitle("Summary")
	root.Summary.SetSizePolicy(tui.Maximum, tui.Maximum)
	root.Summary.Box.Append(tui.NewLabel(fmt.Sprintf("%19v", "")))

	root.History = NewScrollList()
	root.Update("")
	root.History.SetBorder(true)
	root.History.SetTitle("History")
	root.History.SetSizePolicy(tui.Expanding, tui.Expanding)
	root.History.SetOnDelete(root.Delete)
	root.History.SetOnConfirm(root.Edit)

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

func (hv *HistoryView) updateHistory(events []finance.Event, expand string) {
	hv.History.Clear()

	for _, e := range events {
		hv.History.Append(
			fmt.Sprintf(
				"[%v] %20v %8v %10v %v",
				e.GetTime().Format(timef),
				e.GetName(),
				e.GetSumS(),
				e.GetCategory(),
				e.GetType(),
			),
		)
		if e.GetType() == "R/"+expand || e.GetType() == expand {
			r := e.(finance.Receipt)
			for _, p := range r.Products {
				hv.History.Append(
					fmt.Sprintf(
						"%16v %20v %8v %10v %v",
						"",
						p.GetName(),
						p.GetSumS(),
						p.GetCategory(),
						p.GetType(),
					),
				)
			}
		}
	}
}

func (hv *HistoryView) updateSummary(events []finance.Event, budget finance.Budget) {
	hv.Summary.Clear()

	spent := map[string]int{}
	for _, e := range events {
		switch v := e.(type) {
		case finance.Transaction:
			cat := e.GetCategory()
			if _, ok := budget.Spending[cat]; ok {
				spent[cat] += e.GetSum()
			} else {
				spent["other"] += e.GetSum()
			}
		case finance.Receipt:
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
		hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:", strings.Title(c))))
		hv.Summary.Append(SummaryFormat(-1, s.Sum, !s.Received()))
	}

	hv.Summary.Append(tui.NewLabel(" "))
	hv.Summary.Append(tui.NewLabel("Spending"))
	for c, s := range budget.Spending {
		hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:", strings.Title(c))))

		if finance.IsSavings(c) {
			hv.Summary.Append(SavingsFormat(s, spent[c]))
		} else {
			hv.Summary.Append(SummaryFormat(s, spent[c], spent[c] > s))
		}
	}
	bt, st := 0, 0
	for k, v := range budget.Spending {
		bt += v
		st += spent[k]
	}

	hv.Summary.Append(tui.NewLabel(" "))
	hv.Summary.Append(tui.NewLabel("Total:"))
	hv.Summary.Append(SummaryFormat(bt, st, st > bt))

	inc, cinc := 0, 0
	for _, v := range budget.Income {
		if v.Received() {
			cinc += v.Sum
		}
		inc += v.Sum
	}
	hv.Summary.Append(tui.NewLabel(" "))
	hv.Summary.Append(tui.NewLabel(fmt.Sprintf("%v:\n%9.2f %9.2f", "Balance", float64(inc-bt)/100., float64(cinc-st)/100.)))
}

func (hv *HistoryView) Update(expand string) {
	budget, err := finance.GetBudget(db.Collection(bDb), dTimeout, hv.time)
	ResolveError(err)
	events, err := finance.GetEvents(db.Collection(tDb), db.Collection(rDb), dTimeout, budget.Start, budget.End)
	ResolveError(err)
	hv.updateHistory(events, expand)
	hv.updateSummary(events, budget)
}

func (hv *HistoryView) Delete(item string) {
	ids := idre.FindAllStringSubmatch(item, -1)[0][1]
	id, err := primitive.ObjectIDFromHex(ids)
	ResolveError(err)
	ctx, _ := context.WithTimeout(context.Background(), dTimeout)
	if strings.Contains(item, "R") {
		_, err = db.Collection(rDb).DeleteOne(ctx, bson.M{
			"_id": id,
		})
		ResolveError(err)
	} else {
		_, err = db.Collection(tDb).DeleteOne(ctx, bson.M{
			"_id": id,
		})
		ResolveError(err)
	}

	hv.Update("")
}

func (hv *HistoryView) Edit(item string) {
	ids := idre.FindAllStringSubmatch(item, -1)[0][1]
	id, err := primitive.ObjectIDFromHex(ids)
	ResolveError(err)
	ctx, _ := context.WithTimeout(context.Background(), dTimeout)
	if strings.Contains(item, "R") {
		r := &finance.Receipt{}
		res := db.Collection(rDb).FindOne(ctx, bson.M{
			"_id": id,
		})
		err = res.Decode(r)
		ResolveError(err)

		aView := NewAddRView(r)
		ui.SetWidget(aView)
		ui.SetFocusChain(aView)
	} else {
		t := &finance.Transaction{}
		res := db.Collection(tDb).FindOne(ctx, bson.M{
			"_id": id,
		})
		err = res.Decode(t)
		ResolveError(err)

		aView := NewAddTView(nil, t)
		ui.SetWidget(aView)
		ui.SetFocusChain(aView)
	}
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
	case "update", "u":
		hv.Update("")
		fallthrough
	case "expand", "exp":
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
		budget, err := finance.GetBudget(db.Collection(bDb), dTimeout, hv.time)
		ResolveError(err)
		sView := NewSetBView(&budget)
		ui.SetWidget(sView)
		ui.SetFocusChain(sView)
	case "time", "view":
		if len(cmd) > 1 {
			if t, err := time.Parse(timefs, cmd[1]); err == nil {
				hv.time = t
			} else if strings.ToLower(cmd[1]) == "now" {
				hv.time = time.Now()
			}
			hv.Update("")
		}
	}

	e.SetText("")
}
