package main

import (
	"context"
	"time"

	"github.com/strosel/noerr"

	"github.com/marcusolsson/tui-go"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db    *mongo.Database
	hView *HistoryView
	err   error
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	db, err = Connect(ctx, "finance")
	noerr.Panic(err)

	hView = GetHistoryView()

	ui, err := tui.New(hView)
	noerr.Panic(err)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Up", func() { hView.History.Scroll(0, -1) })
	ui.SetKeybinding("Down", func() { hView.History.Scroll(0, 1) })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}