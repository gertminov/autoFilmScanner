package main

import (
	"fmt"
	config2 "scanwatcher/main/config"
	"scanwatcher/main/serial"
	"scanwatcher/main/ui"
	"scanwatcher/main/vuescan"
)

func main() {
	config := config2.InitConfig()
	path, portName := ui.AskForPathAndPort(config.WithUI)
	fmt.Println(path)
	watcher := initFileWatcher(config, path)
	serial.InitSerial(portName, config)
	vuescan.InitKeyboard(config.Images.Dia)
	if config.Calibrate {
		serial.CalibrateSlate()
	}
	serial.MoveToStartPosition()
	ui.Alert("Please insert the film holder. AFTER that click ok to move to first image", "first image")
	ui.Alert("please click once into the Vuescan window", "focus vuescan")
	serial.MoveToFirstFrame()
	vuescan.Scan()

	fmt.Println("everything set up, you can start scanning now")
	watch(watcher, path)
}
