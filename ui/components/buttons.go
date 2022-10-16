package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyneTest/backend/config"
	"fyneTest/backend/serial"
	"fyneTest/backend/watcher"
)

var StartBtn *widget.Button

func GetButtons(appConfig *config.Config, ports []string, fileWatcher *watcher.FileWatcher, window *fyne.Window, slateReady chan bool) *fyne.Container {
	resetAndCal := container.NewHBox(
		widget.NewButton("Calibrate", func() {
			serial.CalibrateSlate()
		}),
		widget.NewButton("Reset", func() {
			serial.MoveToStartPosition()
		}))

	StartBtn = widget.NewButton("Calibrating...", func() {
		dialog.ShowInformation("Start Scanning", "please start the first scan in VueScan.\n all following frames will be scanned automatically", *window)
		go serial.MoveToFirstFrame()
		go watcher.Watch(fileWatcher.Watcher, appConfig)
		StartBtn.SetText("Running")
	})

	fileWatcher.StartBtn = StartBtn

	portselector := widget.NewSelect(ports, func(s string) {
		go serial.InitSerial(s, appConfig)
	})
	if len(ports) > 0 {
		portselector.Selected = ports[0]
	}
	portSelect := container.NewHBox(
		widget.NewLabel("Port:"),
		portselector,
	)

	StartBtn.Importance = widget.HighImportance
	StartBtn.Disable()
	go func() {
		for {
			if <-slateReady {
				StartBtn.SetText("Start scanning")
				StartBtn.Enable()
				return
			}
		}
	}()

	return container.NewGridWithRows(1,
		resetAndCal,
		portSelect,
		layout.NewSpacer(),
		StartBtn,
	)
}
