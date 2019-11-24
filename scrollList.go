package main

import (
	"github.com/marcusolsson/tui-go"
)

type ScrollList struct {
	*tui.Box
	scroll *tui.ScrollArea
	store  *tui.List
}

func NewScrollList() *ScrollList {
	sl := ScrollList{}

	sl.store = tui.NewList()
	sl.store.SetSizePolicy(tui.Preferred, tui.Expanding)
	sl.scroll = tui.NewScrollArea(sl.store)
	sl.scroll.SetSizePolicy(tui.Preferred, tui.Expanding)
	sl.Box = tui.NewVBox(sl.scroll)

	return &sl
}

func (sl *ScrollList) Clear() {
	sl.store.RemoveItems()
}

func (sl *ScrollList) Append(items ...string) {
	sl.store.AddItems(items...)
}

func (sl *ScrollList) Scroll(dx, dy int) {
	sl.scroll.Scroll(dx, dy)
}

func (sl *ScrollList) OnKeyEvent(ev tui.KeyEvent) {
	i := sl.store.Selected()
	l := sl.store.Length()
	if i != -1 {
		switch ev.Name() {
		case "Up":
			if i > 0 {
				sl.store.Select(i - 1)
				if i > 3 {
					sl.Scroll(0, -1)
				}
			}
		case "Down":
			if i < l-1 {
				sl.store.Select(i + 1)
				if i > 2 {
					sl.Scroll(0, 1)
				}
			}
		case "Left":
			sl.store.Select(0)
		case "Right":
			sl.store.Select(l - 1)
		}
	}
}

func (sl *ScrollList) SetFocused(b bool) {
	if b {
		sl.store.SetSelected(0)
	} else {
		sl.store.SetSelected(-1)
	}
}
