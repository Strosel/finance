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
	sb.store.SetSizePolicy(tui.Preferred, tui.Expanding)
	sb.scroll = tui.NewScrollArea(sb.store)
	sb.scroll.SetSizePolicy(tui.Preferred, tui.Expanding)
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

func (sb *ScrollBox) OnKeyEvent(ev tui.KeyEvent) {
	if sb.IsFocused() {
		switch ev.Name() {
		case "Up":
			sb.scroll.Scroll(0, -1)
		case "Down":
			sb.scroll.Scroll(0, 1)
		case "Left":
			sb.scroll.ScrollToTop()
		case "Right":
			sb.scroll.ScrollToBottom()
		}
	}
}

func (sb *ScrollBox) SetFocused(b bool) {
	sb.scroll.SetFocused(b)
}
