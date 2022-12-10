package ui

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyneTest/backend/config"
	"fyneTest/backend/watcher"
	"fyneTest/ui/components"
	"fyneTest/ui/components/customWidgets"
	"fyneTest/ui/customTheme"
)

var window fyne.Window

func GetUI(appConfig *config.Config, serialPorts []string, fileWatcher *watcher.FileWatcher, slateReady chan bool) fyne.Window {

	a := app.New()

	a.Lifecycle().SetOnStopped(func() {
		config.WriteConfig(appConfig) //write config to disk when app is closed
	})
	a.Settings().SetTheme(customTheme.MyTheme{})
	window = a.NewWindow("ScanWatcher")
	window.Resize(fyne.Size{
		Width:  300,
		Height: 300,
	})

	customWidgets.Window = &window

	fileWatcher.FinishedDialog = ShowPopUp
	if len(serialPorts) > 1 {
		dialog.ShowInformation("Multiple Ports", "multiple serial ports connected is not tested but should work.\n please select the correct port", window)
	}
	mainWindow := container.NewVBox(
		components.GetModeAndPath(appConfig, fileWatcher.Watcher),
		widget.NewSeparator(),
		components.GetImagesSection(appConfig, fileWatcher),
		widget.NewSeparator(),
		components.GetButtons(appConfig, serialPorts, fileWatcher, &window, slateReady),
	)

	window.SetContent(mainWindow)
	return window
}

func GetNoPortsError() fyne.Window {
	a := app.New()
	errWindow := a.NewWindow("No Serial Ports Found")
	errWindow.Resize(fyne.Size{
		Width:  800,
		Height: 300,
	})
	err := errors.New("No Serial Ports found. please connect your arduino and restart the software")
	errWindow.SetContent(
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	)
	dialog.ShowError(err, errWindow)
	return errWindow
}

func ShowPopUp(title string, content string) {
	dialog.ShowInformation(title, content, window)
	//dialog.ShowInformation(title, content, window)
}
