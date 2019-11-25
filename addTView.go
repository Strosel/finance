package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/strosel/noerr"

	"github.com/marcusolsson/tui-go"
)

type AddTView struct {
	*tui.Box
	Datei   *tui.Entry
	Namei   *tui.Entry
	Sumi    *tui.Entry
	Cati    *tui.Entry
	Notei   *tui.Entry
	Saveb   *tui.Button
	Cancelb *tui.Button

	Parent      *AddRView
	Transaction *Transaction
}

func NewAddTView(p *AddRView, t *Transaction) *AddTView {
	if t == nil {
		t = &Transaction{}
	}
	root := AddTView{
		Parent:      p,
		Transaction: t,
	}

	var boxd *tui.Box
	if p == nil {
		root.Datei = tui.NewEntry()
		root.Datei.SetFocused(true)
		root.Datei.SetSizePolicy(tui.Expanding, tui.Minimum)
		boxd = tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Date (yy-mm-dd hh:mm): ")), root.Datei)
		boxd.SetSizePolicy(tui.Expanding, tui.Maximum)
	}

	root.Namei = tui.NewEntry()
	root.Namei.SetSizePolicy(tui.Expanding, tui.Minimum)
	nbox := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Name: ")), root.Namei)
	nbox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Sumi = tui.NewEntry()
	root.Sumi.SetSizePolicy(tui.Expanding, tui.Minimum)
	rbox := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Sum: ")), root.Sumi)
	rbox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Cati = tui.NewEntry()
	root.Cati.SetSizePolicy(tui.Expanding, tui.Minimum)
	cbox := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Category: ")), root.Cati)
	cbox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Notei = tui.NewEntry()
	root.Notei.SetSizePolicy(tui.Expanding, tui.Minimum)
	nobox := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Note: ")), root.Notei)
	nobox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Cancelb = tui.NewButton("[Cancel]")
	root.Cancelb.OnActivated(root.Cancel)
	root.Saveb = tui.NewButton("[Save]")
	root.Saveb.OnActivated(root.Save)

	bbox := tui.NewHBox(tui.NewSpacer(), root.Cancelb, tui.NewLabel(" "), root.Saveb, tui.NewLabel(" "))

	root.Box = tui.NewVBox(nbox, rbox, cbox, nobox, bbox)

	if p == nil {
		root.Box.Insert(0, boxd)
		if t.ID.IsZero() {
			root.Datei.SetText(time.Now().Format("06-01-02 15:04"))
		} else {
			root.Datei.SetText(root.Transaction.Datetime.Format("06-01-02 15:04"))
		}
	}

	root.Box.SetBorder(true)

	if t.ID.IsZero() {
		root.Box.SetTitle("Add")
	} else {
		root.Box.SetTitle("Update")
		root.Namei.SetText(root.Transaction.Name)
		root.Sumi.SetText(fmt.Sprintf("%8.2f", float64(root.Transaction.Sum)/100))
		root.Cati.SetText(root.Transaction.Category)
		root.Notei.SetText(root.Transaction.Note)
	}

	return &root
}

func (av *AddTView) Save(b *tui.Button) {
	av.Transaction.Name = av.Namei.Text()
	av.Transaction.Category = av.Cati.Text()
	av.Transaction.Note = av.Notei.Text()

	sum, err := strconv.ParseFloat(av.Sumi.Text(), 64)
	noerr.Panic(err)
	av.Transaction.Sum = int(sum * 100)

	//! handle errors
	if av.Parent == nil {
		av.Transaction.Datetime, err = time.Parse("06-01-02 15:04", av.Datei.Text())
		noerr.Panic(err)
		//save to db
	} else {
		av.Parent.Receipt.Products = append(av.Parent.Receipt.Products, *av.Transaction)
		av.Cancel(b)
	}
}

func (av *AddTView) Cancel(b *tui.Button) {
	if av.Parent == nil {
		ui.SetWidget(hView)
		ui.SetFocusChain(hView)
	} else {
		ui.SetWidget(av.Parent)
		ui.SetFocusChain(av.Parent)
		av.Parent.Update()
	}
}

func (av *AddTView) FocusNext(w tui.Widget) tui.Widget {
	switch w {
	case av.Datei:
		return av.Namei
	case av.Namei:
		return av.Sumi
	case av.Sumi:
		return av.Cati
	case av.Cati:
		return av.Notei
	case av.Notei:
		return av.Cancelb
	case av.Cancelb:
		return av.Saveb
	case av.Saveb:
		if av.Parent == nil {
			return av.Datei
		}
		return av.Namei
	default:
		return nil
	}
}

func (av *AddTView) FocusPrev(w tui.Widget) tui.Widget {
	switch w {
	case av.Datei:
		return av.Saveb
	case av.Namei:
		if av.Parent == nil {
			return av.Datei
		}
		return av.Saveb
	case av.Sumi:
		return av.Namei
	case av.Cati:
		return av.Sumi
	case av.Notei:
		return av.Cati
	case av.Cancelb:
		return av.Notei
	case av.Saveb:
		return av.Cancelb
	default:
		return nil
	}
}

func (av *AddTView) FocusDefault() tui.Widget {
	if av.Parent == nil {
		return av.Datei
	}
	return av.Namei
}
