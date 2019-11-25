package main

import (
	"fmt"

	"github.com/marcusolsson/tui-go"
)

type SetBView struct {
	*tui.Box
	//General
	Starti *tui.Entry
	Endi   *tui.Entry
	Saveb  *tui.Button
	//Income
	Nameii *tui.Entry
	Sumii  *tui.Entry
	Dateii *tui.Entry
	Addib  *tui.Button
	Seti   *ScrollList
	//Spending
	Namesi *tui.Entry
	Sumsi  *tui.Entry
	Addsb  *tui.Button
	Sets   *ScrollList

	Budget *Budget
}

func NewSetBView(b *Budget) *SetBView {
	if b == nil {
		b = &Budget{
			Spending: map[string]int{},
			Income:   map[string]Income{},
		}
	}
	root := SetBView{
		Budget: b,
	}

	root.Nameii = tui.NewEntry()
	root.Nameii.SetSizePolicy(tui.Expanding, tui.Minimum)
	root.Sumii = tui.NewEntry()
	root.Sumii.SetSizePolicy(tui.Expanding, tui.Minimum)
	root.Dateii = tui.NewEntry()
	root.Dateii.SetSizePolicy(tui.Expanding, tui.Minimum)
	root.Addib = tui.NewButton("[Add]")
	root.Addib.OnActivated(root.Addi)
	root.Seti = NewScrollList()
	root.Seti.SetBorder(true)

	inbox := tui.NewVBox(
		tui.NewLabel("Name:"),
		root.Nameii,
		tui.NewLabel("Sum:"),
		root.Sumii,
		tui.NewLabel("Date (yy-mm-dd hh:mm):"),
		root.Dateii,
		tui.NewSpacer(),
		root.Addib,
		root.Seti,
	)
	inbox.SetBorder(true)
	inbox.SetTitle("income")

	root.Namesi = tui.NewEntry()
	root.Namesi.SetSizePolicy(tui.Minimum, tui.Maximum)
	root.Sumsi = tui.NewEntry()
	root.Sumsi.SetSizePolicy(tui.Minimum, tui.Maximum)
	root.Addsb = tui.NewButton("[Add]")
	root.Addsb.OnActivated(root.Adds)
	root.Sets = NewScrollList()
	root.Sets.SetBorder(true)

	spbox := tui.NewVBox(
		tui.NewLabel("Name:"),
		root.Namesi,
		tui.NewLabel("Sum:"),
		root.Sumsi,
		tui.NewSpacer(),
		root.Addsb,
		root.Sets,
	)
	spbox.SetBorder(true)
	spbox.SetTitle("Spending")

	root.Box = tui.NewHBox(inbox, spbox)

	root.Update()
	return &root
}

func (sv *SetBView) Addi(b *tui.Button) {
	sv.Budget.Income[sv.Nameii.Text()] = Income{}
	sv.Update()
}

func (sv *SetBView) Adds(b *tui.Button) {
	sv.Budget.Spending[sv.Namesi.Text()] = 0
	sv.Update()
}

func (sv *SetBView) Update() {
	sv.Seti.Clear()
	sv.Sets.Clear()

	for n, i := range sv.Budget.Income {
		sv.Seti.Append(fmt.Sprintf("%-10v %16v %8.3f", n, i.Date.Format("06-01-02 15:04"), float64(i.Sum)/100.))
	}

	for n, s := range sv.Budget.Spending {
		sv.Sets.Append(fmt.Sprintf("%-20v %8.3f", n, float64(s)/100.))
	}

}

func (sv *SetBView) FocusNext(w tui.Widget) tui.Widget {
	switch w {
	//Income Block
	case sv.Nameii:
		return sv.Sumii
	case sv.Sumii:
		return sv.Dateii
	case sv.Dateii:
		return sv.Addib
	case sv.Addib:
		return sv.Seti
	case sv.Seti:
		return sv.Namesi
	//Spending Block
	case sv.Namesi:
		return sv.Sumsi
	case sv.Sumsi:
		return sv.Addsb
	case sv.Addsb:
		return sv.Sets
	case sv.Sets:
		return sv.Nameii
	default:
		return nil
	}
}

func (sv *SetBView) FocusPrev(w tui.Widget) tui.Widget {
	switch w {
	//Income Block
	case sv.Nameii:
		return sv.Sets
	case sv.Sumii:
		return sv.Nameii
	case sv.Dateii:
		return sv.Sumii
	case sv.Addib:
		return sv.Dateii
	case sv.Seti:
		return sv.Addib
	//Spending Block
	case sv.Namesi:
		return sv.Seti
	case sv.Sumsi:
		return sv.Namesi
	case sv.Addsb:
		return sv.Sumsi
	case sv.Sets:
		return sv.Addsb
	default:
		return nil
	}
}

func (sv *SetBView) FocusDefault() tui.Widget {
	return sv.Nameii
}
