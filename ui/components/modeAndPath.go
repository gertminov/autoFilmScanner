package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyneTest/backend/config"
	"fyneTest/backend/serial"
	"github.com/fsnotify/fsnotify"
	"github.com/ncruces/zenity"
	"os"
)

var fsWatcher *fsnotify.Watcher

func GetModeAndPath(config *config.Config, fileWatcher *fsnotify.Watcher) *fyne.Container {
	fsWatcher = fileWatcher
	spacer := layout.NewSpacer()
	spacer.Resize(fyne.Size{Width: 20})
	return container.New(
		layout.NewFormLayout(),
		filmDiaSelector(config),
		//widget.NewSeparator(),
		container.New(
			layout.NewPaddedLayout(),

			container.NewVBox(
				widget.NewLabel("Path to watch"),
				fileSelector(config),
				layout.NewSpacer()),
		))
}

func filmDiaSelector(appConfig *config.Config) *fyne.Container {
	options := []string{"Dia", "Film"}
	radioGroup := widget.NewRadioGroup(options, func(s string) {
		if s == "Dia" {
			appConfig.Dia = true
		} else {
			appConfig.Dia = false
		}
		serial.SwitchFrameType(appConfig)
	})

	if appConfig.Dia {
		radioGroup.Selected = "Dia"
	} else {
		radioGroup.Selected = "Film"
	}

	return container.NewHBox(radioGroup, widget.NewSeparator())
}

func fileSelector(appConfig *config.Config) *fyne.Container {
	defaultPath := ""
	if len(appConfig.DefaultPath) > 0 {
		defaultPath = appConfig.DefaultPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			defaultPath = ""
			fmt.Println(err)
		} else {
			defaultPath = homeDir
		}
	}

	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		fsWatcher.Remove(appConfig.DefaultPath)
		appConfig.DefaultPath = s // write current path to config
		fsWatcher.Add(appConfig.DefaultPath)
	}

	selectBtn := widget.NewButton("select Folder", func() {
		file := openFileSelect(defaultPath)
		entry.SetText(file)
	})

	entry.SetText(defaultPath)

	return container.New(
		layout.NewFormLayout(),
		selectBtn,
		entry,
	)
}

func openFileSelect(defaultPath string) string {
	file, err := zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.Directory())
	if err != nil {
		fmt.Println(err)
	}
	return file
}
