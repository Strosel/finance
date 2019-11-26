package main

import (
	"github.com/marcusolsson/tui-go"
)

type ScrollList struct {
	*tui.Box
	scroll *tui.ScrollArea
	store  *tui.List

	sel int

	delete  func(string)
	confirm func(string)
}

func NewScrollList() *ScrollList {
	sl := ScrollList{}

	sl.store = tui.NewList()
	sl.store.SetSizePolicy(tui.Expanding, tui.Expanding)
	sl.scroll = tui.NewScrollArea(sl.store)
	sl.scroll.SetSizePolicy(tui.Expanding, tui.Expanding)
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
		case "Backspace", "Backspace2", "Delete":
			if sl.delete != nil {
				sl.delete(sl.store.SelectedItem())
				if i-1 == -1 {
					sl.store.Select(0)
				} else {
					sl.store.Select(i - 1)
				}
			}
		case "Enter":
			if sl.confirm != nil {
				sl.confirm(sl.store.SelectedItem())
			}
		}
	}
}

func (sl *ScrollList) SetFocused(b bool) {
	if b {
		sl.store.SetSelected(sl.sel)
	} else {
		if s := sl.store.Selected(); s != -1 {
			sl.sel = s
		} else {
			sl.sel = 0
		}
		sl.store.SetSelected(-1)
	}
}

func (sl *ScrollList) SetOnDelete(f func(string)) {
	sl.delete = f
}

func (sl *ScrollList) SetOnConfirm(f func(string)) {
	sl.confirm = f
}
