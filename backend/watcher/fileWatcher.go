package watcher

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"fyneTest/backend/config"
	"fyneTest/backend/keyboard"
	"fyneTest/backend/serial"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strconv"
	//"strconv"
)

var imagesScanned int
var imagesPerStrip int
var oldFiles []string
var appConfig *config.Config

//var FinishedDialog func(title string, content string, finished chan bool)

type FileWatcher struct {
	Watcher        *fsnotify.Watcher
	SetCurFrame    func(frame int)
	FinishedDialog func(title string, content string)
	StartBtn       *widget.Button
}

var fileWatcher FileWatcher

func InitFileWatcher(config *config.Config) *FileWatcher {
	appConfig = config
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//TODO
	var watchPath = config.DefaultPath
	if config.DefaultPath == "" {
		watchPath, err = os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
	}
	dir, err := os.ReadDir(watchPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range dir {
		name := entry.Name()
		oldFiles = append(oldFiles, name)
	}
	imagesScanned = 0
	imagesPerStrip = config.Images.ImagesPerStrip
	fileWatcher = FileWatcher{
		Watcher: watcher,
	}
	return &fileWatcher
}

// test comment
func Watch(watcher *fsnotify.Watcher, appConfig *config.Config) {
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event) // TODO remove for production
				if event.Op&fsnotify.Write == fsnotify.Write {
					if eventFromOldFile(event.Name) {
						continue
					}
					fmt.Println("created file:", event.Name)
					newStrip := countImagesScanned()
					if newStrip {
						for len(watcher.Events) > 0 {
							<-watcher.Events
						}
						continue
					}
					serial.SendTurn()
					serial.WaitForMotor()
					keyboard.Scan()
				}
			case err, ok := <-watcher.Errors:

				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := watcher.Add(appConfig.DefaultPath)
	if err != nil {

		log.Fatal(err)
	}
	<-done
}

func countImagesScanned() bool {
	imagesScanned++
	//fileWatcher.SetCurFrame(imagesScanned)
	if imagesScanned >= appConfig.ImagesPerStrip {
		serial.MoveToStartPosition()
		//message := "Alle Bilder in dem Strip (" + strconv.Itoa(imagesPerStrip) + ") gescanned, lege neue Bilder ein, warte bis das erste bild in position ist und drücke scannen in Vuescan"
		//ui.Alert(message, "finished scanning")
		finished := make(chan bool)
		fileWatcher.FinishedDialog("finished", "Alle Bilder in dem Strip ("+strconv.Itoa(imagesPerStrip)+") gescanned, lege neue Bilder ein und klick auf continue. Warte bis das erste bild in position ist und drücke scannen in Vuescan")
		fileWatcher.StartBtn.SetText("Continue")
		fileWatcher.StartBtn.OnTapped = func() {
			finished <- true
		}
		<-finished
		serial.MoveToFirstFrame()
		fileWatcher.StartBtn.SetText("Running")
		imagesScanned = 0
		return true
	}
	return false
}

func eventFromOldFile(path string) bool {
	filename := filepath.Base(path)
	for _, entry := range oldFiles {
		if entry == filename {
			return true
		}
	}
	return false
}
