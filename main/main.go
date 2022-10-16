package main

import (
	"fyne.io/fyne/v2"
	"fyneTest/backend/config"
	"fyneTest/backend/keyboard"
	"fyneTest/backend/serial"
	"fyneTest/backend/watcher"
	"fyneTest/ui"
)

func main() {

	appConfig := config.InitConfig()
	ports := serial.GetAvailablePorts()
	var window fyne.Window
	if len(ports) < 1 {
		window = ui.GetNoPortsError()
	} else {

		serial.InitSerial(ports[0], &appConfig)
		slateReady := make(chan bool, 1)
		fileWatcher := watcher.InitFileWatcher(&appConfig)

		window = ui.GetUI(&appConfig, ports, fileWatcher, slateReady)

		go keyboard.InitKeyboard(&appConfig.Images)

		go serial.InitSlate(appConfig.Calibrate, slateReady)
	}
	//go watcher.Watch(fileWatcher, &appConfig)
	window.ShowAndRun()
}
