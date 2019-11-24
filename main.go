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
	ui    tui.UI
	hView *HistoryView
	err   error
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	db, err = Connect(ctx, "finance")
	noerr.Panic(err)

	hView = GetHistoryView()

	ui, err = tui.New(hView)
	noerr.Panic(err)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetFocusChain(hView)
	ui.SetTheme(GetTheme())

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
