package main

import (
	"github.com/marcusolsson/tui-go"
)

type ScrollBox struct {
	*tui.Box
	scroll *tui.ScrollArea
	store  *tui.Box
}

func NewScrollBox() *ScrollBox {
	sb := ScrollBox{}

	sb.store = tui.NewVBox()
	sb.scroll = tui.NewScrollArea(sb.store)
	sb.Box = tui.NewVBox(sb.scroll)

	return &sb
}

func (sb *ScrollBox) Clear() {
	for sb.store.Length() != 0 {
		sb.store.Remove(0)
	}
}

func (sb *ScrollBox) Append(w tui.Widget) {
	sb.store.Append(w)
}

func (sb *ScrollBox) Scroll(dx, dy int) {
	sb.scroll.Scroll(dx, dy)
}

func (sb *ScrollBox) ScrollToTop() {
	sb.scroll.ScrollToTop()
}