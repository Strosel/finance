package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

	lbls map[string]*tui.Label

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
		lbls:        map[string]*tui.Label{},
	}

	var boxd *tui.Box
	if p == nil {
		root.Datei = tui.NewEntry()
		root.Datei.SetFocused(true)
		root.Datei.SetSizePolicy(tui.Expanding, tui.Minimum)
		root.lbls["date"] = tui.NewLabel(fmt.Sprintf("%25v", "Date (yy-mm-dd hh:mm): "))
		boxd = tui.NewHBox(root.lbls["date"], root.Datei)
		boxd.SetSizePolicy(tui.Expanding, tui.Maximum)
	}

	root.Namei = tui.NewEntry()
	root.Namei.SetSizePolicy(tui.Expanding, tui.Minimum)
	nbox := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Name: ")), root.Namei)
	nbox.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Sumi = tui.NewEntry()
	root.Sumi.SetSizePolicy(tui.Expanding, tui.Minimum)
	root.lbls["sum"] = tui.NewLabel(fmt.Sprintf("%25v", "Sum: "))
	rbox := tui.NewHBox(root.lbls["sum"], root.Sumi)
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
			root.Datei.SetText(time.Now().Format(timef))
		} else {
			root.Datei.SetText(root.Transaction.Datetime.Format(timef))
		}
	}

	root.Box.SetBorder(true)

	if t.ID.IsZero() {
		root.Box.SetTitle("Add")
	} else {
		root.Box.SetTitle("Update")
		root.Namei.SetText(root.Transaction.Name)
		root.Sumi.SetText(strings.TrimSpace(root.Transaction.GetSumS()))
		root.Cati.SetText(root.Transaction.Category)
		root.Notei.SetText(root.Transaction.Note)
	}

	return &root
}

func (av *AddTView) Save(b *tui.Button) {
	for _, l := range av.lbls {
		l.SetStyleName("normal")
	}

	av.Transaction.Name = av.Namei.Text()
	av.Transaction.Category = av.Cati.Text()
	av.Transaction.Note = av.Notei.Text()

	sum, err := strconv.Atoi(flre.ReplaceAllString(av.Sumi.Text(), ""))
	if err != nil {
		av.lbls["sum"].SetStyleName("warning")
		return
	}
	av.Transaction.Sum = sum

	if av.Parent == nil {
		av.Transaction.Datetime, err = time.Parse(timef, av.Datei.Text())
		if err != nil {
			av.lbls["date"].SetStyleName("warning")
			return
		}
		if av.Transaction.ID.IsZero() {
			av.Transaction.ID = primitive.NewObjectID()
			ctx, _ := context.WithTimeout(context.Background(), dTimeout)
			_, err := db.Collection(tDb).InsertOne(ctx, av.Transaction)
			if err != nil {
				ui.SetWidget(NewErrorView(err))
				ui.SetFocusChain(nil)
			}
		} else {
			ctx, _ := context.WithTimeout(context.Background(), dTimeout)
			_, err := db.Collection(tDb).ReplaceOne(ctx,
				bson.M{
					"_id": av.Transaction.ID,
				}, av.Transaction)
			if err != nil {
				ui.SetWidget(NewErrorView(err))
				ui.SetFocusChain(nil)
			}
		}
	} else {
		if av.Transaction.ID.IsZero() {
			av.Transaction.ID = primitive.ObjectID{1}
			av.Parent.Receipt.Products = append(av.Parent.Receipt.Products, av.Transaction)
		}
	}
	av.Cancel(b)
}

func (av *AddTView) Cancel(b *tui.Button) {
	if av.Parent == nil {
		ui.SetWidget(hView)
		ui.SetFocusChain(hView)
		hView.Update("")
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
