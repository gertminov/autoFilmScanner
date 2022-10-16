package customWidgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var Window *fyne.Window

func NewConfirm(title string, content string, callback func(b bool)) {
	fmt.Println("hallo")
	dialog.ShowConfirm(title, content, callback, *Window)
}
