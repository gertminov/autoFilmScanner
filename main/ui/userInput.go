package ui

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/sqweek/dialog"
	"github.com/tcnksm/go-input"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"scanwatcher/main/serial"
)

var ui *input.UI

func AskForPathAndPort(withUI bool) (string, string) {
	ui = &input.UI{Writer: os.Stdout, Reader: os.Stdin}

	OpenConfig()
	path := getPathToWatch(withUI)
	port := getArduinoPort()

	return path, port
}

func getArduinoPort() string {
	ports := serial.GetAvailablePorts()

	if len(ports) == 1 {
		return ports[0]
	}

	portQuery := "mit welchem Port soll kommuniziert werden?"
	port, err := ui.Select(portQuery, ports, &input.Options{
		Default: ports[0],
	})
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func getPathToWatch(withUI bool) string {
	var err error
	var path string

	if withUI {
		path, err = dialog.Directory().Title("Ordner der Überwacht werden soll").Browse()
	} else {
		prompt := promptui.Prompt{
			Label: "Pfad zu dem Ordner der Überwacht werden soll",
		}
		path, err = prompt.Run()
	}
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		fmt.Println("Filepath does not exits")
		path = getPathToWatch(false)
	}
	return path
}

func Alert(message string, title string) {
	dialog.Message(message).Title(title).Info()
}

func OpenConfig() {
	prompt := promptui.Select{
		Label: "Konfigurationsdatei öffnen?",
		Items: []string{"Nein", "Ja"},
	}
	_, answer, err := prompt.Run()
	//answer, err := ui.Select("Konfigurationsdatei öffnen?", []string{"Ja", "Nein"}, &input.Options{Default: "Nein"})
	if err != nil {
		log.Printf("input: %v", answer)
		log.Fatal(err)
	}

	if answer == "Ja" {
		abstPath, _ := filepath.Abs("./config.yaml")
		switch runtime.GOOS {
		case "windows":
			err := exec.Command("C:\\Windows\\system32\\notepad.exe", abstPath).Run()
			if err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Println("Du findest die Konfigurationsdatei unter: ")
			fmt.Println(filepath.Abs("./config.yaml"))
		}

	}
}
