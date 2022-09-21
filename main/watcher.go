package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"scanwatcher/main/config"
	"scanwatcher/main/serial"
	"scanwatcher/main/ui"
	"scanwatcher/main/vuescan"
	"strconv"
)

var imagesScanned int
var imagesPerStrip int
var oldFiles []string

func initFileWatcher(config config.Config, path string) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range dir {
		name := entry.Name()
		oldFiles = append(oldFiles, name)
	}
	imagesScanned = 0
	imagesPerStrip = config.Images.ImagesPerStrip
	return watcher
}

// test comment
func watch(watcher *fsnotify.Watcher, path string) {
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
					vuescan.Scan()
				}
			case err, ok := <-watcher.Errors:

				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := watcher.Add(path)
	if err != nil {

		log.Fatal(err)
	}
	<-done
}

func countImagesScanned() bool {
	imagesScanned++
	if imagesScanned >= imagesPerStrip-1 {
		serial.MoveToStartPosition()
		message := "Alle Bilder in dem Strip (" + strconv.Itoa(imagesPerStrip) + ") gescanned, lege neue Bilder ein, warte bis das erste bild in position ist und drücke scannen in Vuescan"
		ui.Alert(message, "finished scanning")
		serial.MoveToFirstFrame()
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
