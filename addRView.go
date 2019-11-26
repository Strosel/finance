package main

import (
	"context"
	"fmt"
	"time"

	"github.com/marcusolsson/tui-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddRView struct {
	*tui.Box
	Datei   *tui.Entry
	Storei  *tui.Entry
	Saveb   *tui.Button
	Cancelb *tui.Button
	Addb    *tui.Button
	Prodb   *ScrollList

	lbls map[string]*tui.Label

	Receipt *Receipt
}

func NewAddRView(r *Receipt) *AddRView {
	if r == nil {
		r = &Receipt{}
	}
	root := AddRView{
		Receipt: r,
		lbls:    map[string]*tui.Label{},
	}

	root.Datei = tui.NewEntry()
	root.Datei.SetFocused(true)
	root.Datei.SetSizePolicy(tui.Expanding, tui.Minimum)
	root.lbls["date"] = tui.NewLabel(fmt.Sprintf("%25v", "Date (yy-mm-dd hh:mm): "))
	boxd := tui.NewHBox(root.lbls["date"], root.Datei)
	boxd.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Storei = tui.NewEntry()
	root.Storei.SetSizePolicy(tui.Expanding, tui.Minimum)
	boxs := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Store: ")), root.Storei)
	boxs.SetSizePolicy(tui.Expanding, tui.Maximum)

	root.Cancelb = tui.NewButton("[Cancel]")
	root.Cancelb.OnActivated(root.Cancel)
	root.Saveb = tui.NewButton("[Save]")
	root.Saveb.OnActivated(root.Save)

	bbox := tui.NewHBox(tui.NewSpacer(), root.Cancelb, tui.NewLabel(" "), root.Saveb, tui.NewLabel(" "))

	root.Addb = tui.NewButton("[Add]")
	root.Addb.OnActivated(root.Add)
	root.Prodb = NewScrollList()
	root.Prodb.SetOnDelete(root.Delete)
	root.Prodb.SetOnConfirm(root.Edit)
	boxp := tui.NewVBox(root.Addb, root.Prodb)
	boxp.SetBorder(true)
	boxp.SetTitle("Products")

	root.Box = tui.NewVBox(boxd, boxs, boxp, bbox)
	root.Box.SetBorder(true)

	if r.ID.IsZero() {
		root.Box.SetTitle("Add")
		root.Datei.SetText(time.Now().Format(timef))
	} else {
		root.Box.SetTitle("Update")
		root.Datei.SetText(root.Receipt.Datetime.Format(timef))
		root.Storei.SetText(root.Receipt.Store)
	}

	root.Update()
	return &root
}

func (av *AddRView) Save(b *tui.Button) {
	for _, l := range av.lbls {
		l.SetStyleName("normal")
	}

	av.Receipt.Datetime, err = time.Parse(timef, av.Datei.Text())
	if err != nil {
		av.lbls["date"].SetStyleName("warning")
		return
	}
	av.Receipt.Store = av.Storei.Text()
	//! handle errors
	if av.Receipt.ID.IsZero() {
		av.Receipt.ID = primitive.NewObjectID()
		ctx, _ := context.WithTimeout(context.Background(), dTimeout)
		_, err := db.Collection(rDb).InsertOne(ctx, av.Receipt)
		if err != nil {
			ui.SetWidget(NewErrorView(err))
			ui.SetFocusChain(nil)
		}
	} else {
		ctx, _ := context.WithTimeout(context.Background(), dTimeout)
		_, err := db.Collection(rDb).ReplaceOne(ctx,
			bson.M{
				"_id": av.Receipt.ID,
			}, av.Receipt)
		if err != nil {
			ui.SetWidget(NewErrorView(err))
			ui.SetFocusChain(nil)
		}
	}
	av.Cancel(b)
}

func (av *AddRView) Cancel(b *tui.Button) {
	ui.SetWidget(hView)
	ui.SetFocusChain(hView)
	hView.Update("")
}

func (av *AddRView) Add(b *tui.Button) {
	aView := NewAddTView(av, nil)
	ui.SetWidget(aView)
	ui.SetFocusChain(aView)
}

func (av *AddRView) Update() {
	av.Prodb.Clear()
	for _, p := range av.Receipt.Products {
		av.Prodb.Append(
			fmt.Sprintf(
				"%10v %8v %10v",
				p.GetName(),
				p.GetSumS(),
				p.GetCategory(),
			),
		)
	}
}

func (av *AddRView) Delete(item string) {
	for i, p := range av.Receipt.Products {
		lbl := fmt.Sprintf(
			"%10v %8v %10v",
			p.GetName(),
			p.GetSumS(),
			p.GetCategory(),
		)

		if item == lbl {
			av.Receipt.Products = append(av.Receipt.Products[:i], av.Receipt.Products[i+1:]...)
			break
		}
	}
	av.Update()
}

func (av *AddRView) Edit(item string) {
	for i, p := range av.Receipt.Products {
		lbl := fmt.Sprintf(
			"%10v %8v %10v",
			p.GetName(),
			p.GetSumS(),
			p.GetCategory(),
		)

		if item == lbl {
			aView := NewAddTView(av, av.Receipt.Products[i])
			ui.SetWidget(aView)
			ui.SetFocusChain(aView)
			break
		}
	}
}

func (av *AddRView) FocusNext(w tui.Widget) tui.Widget {
	switch w {
	case av.Datei:
		return av.Storei
	case av.Storei:
		return av.Addb
	case av.Addb:
		return av.Prodb
	case av.Prodb:
		return av.Cancelb
	case av.Cancelb:
		return av.Saveb
	case av.Saveb:
		return av.Datei
	default:
		return nil
	}
}

func (av *AddRView) FocusPrev(w tui.Widget) tui.Widget {
	switch w {
	case av.Datei:
		return av.Saveb
	case av.Storei:
		return av.Datei
	case av.Addb:
		return av.Storei
	case av.Prodb:
		return av.Addb
	case av.Cancelb:
		return av.Prodb
	case av.Saveb:
		return av.Cancelb
	default:
		return nil
	}
}

func (av *AddRView) FocusDefault() tui.Widget {
	return av.Datei
}
