package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/strosel/noerr"

	"github.com/marcusolsson/tui-go"
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

	gbox := tui.NewHBox(inbox, spbox)

	root.Starti = tui.NewEntry()
	root.Endi = tui.NewEntry()
	tbox := tui.NewVBox(
		tui.NewLabel("Start (yy-mm-dd hh:mm):"),
		root.Starti,
		tui.NewLabel("End (yy-mm-dd hh:mm):"),
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
		root.Dateii.SetText(time.Now().Format(timef))
		root.Starti.SetText(time.Now().Format(timef))
		root.Endi.SetText(time.Now().AddDate(0, 1, -1).Format(timef))
	} else {
		root.Starti.SetText(root.Budget.Start.Format(timef))
		root.Endi.SetText(root.Budget.End.Format(timef))
	}

	root.Update()
	return &root
}

func (sv *SetBView) Addi(b *tui.Button) {
	sum, err := strconv.ParseFloat(sv.Sumii.Text(), 64)
	noerr.Panic(err)
	date, err := time.Parse(timef, sv.Dateii.Text())
	noerr.Panic(err)
	sv.Budget.Income[sv.Nameii.Text()] = Income{
		Sum:  int(100 * sum),
		Date: date,
	}
	sv.Update()

	sv.Nameii.SetText("")
	sv.Sumii.SetText("")
	sv.Dateii.SetText(time.Now().Format(timef))
}

func (sv *SetBView) Adds(b *tui.Button) {
	sum, err := strconv.ParseFloat(sv.Sumsi.Text(), 64)
	noerr.Panic(err)
	sv.Budget.Spending[sv.Namesi.Text()] = int(100 * sum)
	sv.Update()

	sv.Namesi.SetText("")
	sv.Sumsi.SetText("")
}

func (sv *SetBView) Update() {
	sv.Seti.Clear()
	sv.Sets.Clear()

	for n, i := range sv.Budget.Income {
		sv.Seti.Append(fmt.Sprintf("%-10v %16v %8.2f", n, i.Date.Format(timef), float64(i.Sum)/100.))
	}

	for n, s := range sv.Budget.Spending {
		sv.Sets.Append(fmt.Sprintf("%-10v %8.2f", n, float64(s)/100.))
	}

}

func (sv *SetBView) Save(b *tui.Button) {
	sv.Budget.Start, err = time.Parse(timef, sv.Starti.Text())
	noerr.Panic(err)
	sv.Budget.End, err = time.Parse(timef, sv.Endi.Text())
	noerr.Panic(err)
	if sv.Budget.ID.IsZero() {
		sv.Budget.ID = primitive.NewObjectID()
		ctx, _ := context.WithTimeout(context.Background(), dTimeout)
		_, err := db.Collection(bDb).InsertOne(ctx, sv.Budget)
		noerr.Panic(err)
	} else {
		//update
	}
	sv.Cancel(b)
}

func (sv *SetBView) Cancel(b *tui.Button) {
	ui.SetWidget(hView)
	ui.SetFocusChain(hView)
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
