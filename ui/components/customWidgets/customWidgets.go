package customWidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type NumericalEntry struct {
	widget.Entry
}

func NewNumericalEntry() *NumericalEntry {
	//e := widget.Entry{Wrapping: fyne.TextTruncate, OnChanged: func(s string) { fmt.Println("hallo") }}
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *NumericalEntry) SetNumber(number int) {
	e.Entry.SetText(strconv.Itoa(number))
}

func (e *NumericalEntry) TypedRune(r rune) {

	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *NumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}
