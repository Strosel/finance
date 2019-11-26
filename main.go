package main

import (
	"context"
	"regexp"
	"time"

	"github.com/strosel/noerr"

	"github.com/marcusolsson/tui-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	savestr  = "spar|spara|sparande|save|saving|savings"
	rDb      = "receipts"     //"testR"
	tDb      = "transactions" //"test"
	bDb      = "budgets"      //"testB"
	dTimeout = time.Minute
	timef    = "06-01-02 15:04"
	timefs   = "06-01-02"
)

var (
	db    *mongo.Database
	ui    tui.UI
	hView *HistoryView
	err   error

	idre = regexp.MustCompile(`[RT]/(\w+)`)
	flre = regexp.MustCompile(`[\.,]`)
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), dTimeout)
	db, err = Connect(ctx, "finance")
	noerr.Fatal(err)

	hView = GetHistoryView()

	ui, err = tui.New(hView)
	noerr.Fatal(err)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetFocusChain(hView)
	ui.SetTheme(GetTheme())

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
