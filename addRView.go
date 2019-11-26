package main

import (
	"context"
	"fmt"
	"time"

	"github.com/marcusolsson/tui-go"
	"github.com/strosel/noerr"
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

	Receipt *Receipt
}

func NewAddRView(r *Receipt) *AddRView {
	if r == nil {
		r = &Receipt{}
	}
	root := AddRView{
		Receipt: r,
	}

	root.Datei = tui.NewEntry()
	root.Datei.SetFocused(true)
	root.Datei.SetSizePolicy(tui.Expanding, tui.Minimum)
	boxd := tui.NewHBox(tui.NewLabel(fmt.Sprintf("%25v", "Date (yy-mm-dd hh:mm): ")), root.Datei)
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
	boxp := tui.NewVBox(root.Addb, root.Prodb)
	boxp.SetBorder(true)
	boxp.SetTitle("Products")

	root.Box = tui.NewVBox(boxd, boxs, boxp, bbox)
	root.Box.SetBorder(true)

	if r.ID.IsZero() {
		root.Box.SetTitle("Add")
		root.Datei.SetText(time.Now().Format("06-01-02 15:04"))
	} else {
		root.Box.SetTitle("Update")
		root.Datei.SetText(root.Receipt.Datetime.Format("06-01-02 15:04"))
		root.Storei.SetText(root.Receipt.Store)
	}

	root.Update()
	return &root
}

func (av *AddRView) Save(b *tui.Button) {
	av.Receipt.Datetime, err = time.Parse("06-01-02 15:04", av.Datei.Text())
	noerr.Panic(err)
	av.Receipt.Store = av.Storei.Text()
	//! handle errors
	if av.Receipt.ID.IsZero() {
		av.Receipt.ID = primitive.NewObjectID()
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		_, err := db.Collection("testR").InsertOne(ctx, av.Receipt)
		noerr.Panic(err)
	} else {
		//update
	}
	av.Cancel(b)
}

func (av *AddRView) Cancel(b *tui.Button) {
	ui.SetWidget(hView)
	ui.SetFocusChain(hView)
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
				"%10v %8.2f %10v",
				p.GetName(),
				float64(p.GetSum())/100,
				p.GetCategory(),
			),
		)
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
