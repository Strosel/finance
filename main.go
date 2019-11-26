package main

import (
	"context"
	"time"

	"github.com/strosel/noerr"

	"github.com/marcusolsson/tui-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	savestr  = "spar|spara|sparande|save|saving|savings"
	rDb      = "testR"
	tDb      = "test"
	bDb      = "testB"
	dTimeout = time.Minute
	timef    = "06-01-02 15:04"
)

var (
	db    *mongo.Database
	ui    tui.UI
	hView *HistoryView
	err   error
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), dTimeout)
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
