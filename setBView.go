package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/marcusolsson/tui-go"
	"github.com/strosel/finance/finance"
)

type SetBView struct {
	*tui.Box
	//General
	Starti  *tui.Entry
	Endi    *tui.Entry
	Saveb   *tui.Button
	Cancelb *tui.Button
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

	lblsi map[string]*tui.Label
	lblss map[string]*tui.Label
	lblsg map[string]*tui.Label

	Budget *finance.Budget
}

func NewSetBView(b *finance.Budget) *SetBView {
	if b == nil {
		b = &finance.Budget{
			Spending: map[string]int{},
			Income:   map[string]finance.Income{},
		}
	}
	root := SetBView{
		Budget: b,
		lblsi: map[string]*tui.Label{
			"sum":  tui.NewLabel("Sum:"),
			"date": tui.NewLabel("Date (yy-mm-dd):"),
		},
		lblss: map[string]*tui.Label{
			"sum": tui.NewLabel("Sum:"),
		},
		lblsg: map[string]*tui.Label{
			"start": tui.NewLabel("Start (yy-mm-dd):"),
			"end":   tui.NewLabel("End (yy-mm-dd):"),
		},
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
	root.Seti.SetOnConfirm(root.Editi)
	root.Seti.SetOnDelete(root.Deli)

	inbox := tui.NewVBox(
		tui.NewLabel("Name:"),
		root.Nameii,
		root.lblsi["sum"],
		root.Sumii,
		root.lblsi["date"],
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
	root.Sets.SetOnConfirm(root.Edits)
	root.Sets.SetOnDelete(root.Dels)
	spbox := tui.NewVBox(
		tui.NewLabel("Name:"),
		root.Namesi,
		root.lblss["sum"],
		root.Sumsi,
		tui.NewSpacer(),
		root.Addsb,
		root.Sets,
	)
	spbox.SetBorder(true)
	spbox.SetTitle("Spending")

	gbox := tui.NewHBox(inbox, spbox)

	root.Starti = tui.NewEntry()
	root.Endi = tui.NewEntry()
	tbox := tui.NewVBox(
		root.lblsg["start"],
		root.Starti,
		root.lblsg["end"],
		root.Endi,
	)
	tbox.SetBorder(true)

	root.Saveb = tui.NewButton("[save]")
	root.Saveb.OnActivated(root.Save)
	root.Cancelb = tui.NewButton("[cancel]")
	root.Cancelb.OnActivated(root.Cancel)
	bbox := tui.NewHBox(tui.NewSpacer(), root.Cancelb, tui.NewLabel(" "), root.Saveb, tui.NewLabel(" "))

	root.Box = tui.NewVBox(gbox, tbox, bbox)

	if b.ID.IsZero() {
		root.Dateii.SetText(time.Now().Format(timefs))
		root.Starti.SetText(time.Now().Format(timefs))
		root.Endi.SetText(time.Now().AddDate(0, 1, -1).Format(timefs))
	} else {
		root.Dateii.SetText(time.Now().Format(timefs))
		root.Starti.SetText(root.Budget.Start.Format(timefs))
		root.Endi.SetText(root.Budget.End.Format(timefs))
	}

	root.Update()
	return &root
}

func (sv *SetBView) Addi(b *tui.Button) {
	for _, l := range sv.lblsi {
		l.SetStyleName("normal")
	}

	sum, err := strconv.Atoi(flre.ReplaceAllString(sv.Sumii.Text(), ""))
	if err != nil {
		sv.lblsi["sum"].SetStyleName("warning")
		return
	}
	date, err := time.Parse(timefs, sv.Dateii.Text())
	if err != nil {
		sv.lblsi["date"].SetStyleName("warning")
		return
	}
	sv.Budget.Income[sv.Nameii.Text()] = finance.Income{
		Sum:  sum,
		Date: date,
	}
	sv.Update()

	sv.Nameii.SetText("")
	sv.Sumii.SetText("")
	sv.Dateii.SetText(time.Now().Format(timefs))
}

func (sv *SetBView) Adds(b *tui.Button) {
	for _, l := range sv.lblss {
		l.SetStyleName("normal")
	}

	sum, err := strconv.Atoi(flre.ReplaceAllString(sv.Sumsi.Text(), ""))
	if err != nil {
		sv.lblss["sum"].SetStyleName("warning")
		return
	}
	sv.Budget.Spending[sv.Namesi.Text()] = sum
	sv.Update()

	sv.Namesi.SetText("")
	sv.Sumsi.SetText("")
}

func (sv *SetBView) Editi(item string) {
	for n, i := range sv.Budget.Income {
		lbl := fmt.Sprintf("%-10v %16v %8.2f", n, i.Date.Format(timefs), float64(i.Sum)/100.)

		if item == lbl {
			sv.Nameii.SetText(n)
			sv.Dateii.SetText(i.Date.Format(timefs))
			sv.Sumii.SetText(fmt.Sprintf("%.2f", float64(i.Sum)/100.))
			delete(sv.Budget.Income, n)
			break
		}
	}
}

func (sv *SetBView) Edits(item string) {
	for n, s := range sv.Budget.Spending {
		lbl := fmt.Sprintf("%-10v %8.2f", n, float64(s)/100.)

		if item == lbl {
			sv.Namesi.SetText(n)
			sv.Sumsi.SetText(fmt.Sprintf("%.2f", float64(s)/100.))
			delete(sv.Budget.Spending, n)
			break
		}
	}
}

func (sv *SetBView) Deli(item string) {
	for n, i := range sv.Budget.Income {
		lbl := fmt.Sprintf("%-10v %16v %8.2f", n, i.Date.Format(timefs), float64(i.Sum)/100.)

		if item == lbl {
			delete(sv.Budget.Income, n)
			break
		}
	}
	sv.Update()
}

func (sv *SetBView) Dels(item string) {
	for n, s := range sv.Budget.Spending {
		lbl := fmt.Sprintf("%-10v %8.2f", n, float64(s)/100.)

		if item == lbl {
			delete(sv.Budget.Spending, n)
			break
		}
	}
	sv.Update()
}

func (sv *SetBView) Update() {
	sv.Seti.Clear()
	sv.Sets.Clear()

	for n, i := range sv.Budget.Income {
		sv.Seti.Append(fmt.Sprintf("%-10v %16v %8.2f", n, i.Date.Format(timefs), float64(i.Sum)/100.))
	}

	for n, s := range sv.Budget.Spending {
		sv.Sets.Append(fmt.Sprintf("%-10v %8.2f", n, float64(s)/100.))
	}

}

func (sv *SetBView) Save(b *tui.Button) {
	for _, l := range sv.lblsg {
		l.SetStyleName("normal")
	}

	sv.Budget.Start, err = time.Parse(timefs, sv.Starti.Text())
	if err != nil {
		sv.lblsg["start"].SetStyleName("warning")
		return
	}
	sv.Budget.End, err = time.Parse(timefs, sv.Endi.Text())
	if err != nil {
		sv.lblsg["end"].SetStyleName("warning")
		return
	}
	if sv.Budget.ID.IsZero() {
		sv.Budget.ID = primitive.NewObjectID()
		ctx, _ := context.WithTimeout(context.Background(), dTimeout)
		_, err := db.Collection(bDb).InsertOne(ctx, sv.Budget)
		ResolveError(err)
	} else {
		ctx, _ := context.WithTimeout(context.Background(), dTimeout)
		_, err := db.Collection(bDb).ReplaceOne(ctx,
			bson.M{
				"_id": sv.Budget.ID,
			}, sv.Budget)
		ResolveError(err)
	}
	sv.Cancel(b)
}

func (sv *SetBView) Cancel(b *tui.Button) {
	ui.SetWidget(hView)
	ui.SetFocusChain(hView)
	hView.Update("")
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
		return sv.Starti
	//General Block
	case sv.Starti:
		return sv.Endi
	case sv.Endi:
		return sv.Cancelb
	case sv.Cancelb:
		return sv.Saveb
	case sv.Saveb:
		return sv.Nameii
	default:
		return nil
	}
}

func (sv *SetBView) FocusPrev(w tui.Widget) tui.Widget {
	switch w {
	//Income Block
	case sv.Nameii:
		return sv.Saveb
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
	//General Block
	case sv.Starti:
		return sv.Sets
	case sv.Endi:
		return sv.Starti
	case sv.Cancelb:
		return sv.Endi
	case sv.Saveb:
		return sv.Cancelb
	default:
		return nil
	}
}

func (sv *SetBView) FocusDefault() tui.Widget {
	return sv.Nameii
}
