package serial

import (
	"fmt"
	"fyneTest/backend/config"
	"go.bug.st/serial"
	"log"
	"strconv"
	"strings"
	"time"
)

var steps uint32
var reverse bool
var reverseSteps uint32
var selectedPort serial.Port
var timeout time.Duration
var startPosition uint32
var firstFrame uint32

func InitSerial(portString string, config *config.Config) {
	initPort(portString, config)

	setSteps(config)

	initReverseSteps(config)
	setStartPosition(config)

	timeout = time.Duration(config.Arduino.TimeOut)

	time.Sleep(2 * time.Second)
}

func SwitchFrameType(appConfig *config.Config) {
	setStartPosition(appConfig)
	setSteps(appConfig)
}

func setStartPosition(config *config.Config) {
	startPosition = config.Arduino.StartPosition

	if config.Dia {
		firstFrame = config.Arduino.FirstDiaPos
	} else {
		firstFrame = config.Arduino.FirstImagePos
	}
}

func initReverseSteps(config *config.Config) {
	reverse = config.Arduino.GoBack
	reverseSteps = config.Arduino.GoBackSteps
}

func initPort(portString string, config *config.Config) {
	mode := &serial.Mode{
		BaudRate: config.Arduino.BoutRate,
	}
	var err error
	selectedPort, err = serial.Open(portString, mode)
	if err != nil {
		log.Fatal(err)
	}
}

func setSteps(config *config.Config) {
	if config.Images.Dia {
		steps = config.Arduino.StepsPerDia
	} else {
		steps = config.Arduino.StepsPerPhoto
	}
}

func GetAvailablePorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Println("No Serieal ports")
	}
	return ports
}

func SendTurn() {
	fmt.Println("steps: ", steps)
	message := buildMessage(steps)
	readAndWriteToSerial(message)
	if reverse {
		readAndWriteToSerial(buildMessage(-reverseSteps))
	}
}

func MoveToStartPosition() {
	start := "to" + strconv.FormatUint(uint64(startPosition), 10)
	readAndWriteToSerial(start)
	WaitForMotor()
}

func MoveToFirstFrame() {
	first := "to" + strconv.FormatUint(uint64(firstFrame), 10)
	readAndWriteToSerial(first)
	WaitForMotor()
}

func MoveToFrame(frame int) {

	frameNr := uint32(frame - 1)
	fmt.Println(frameNr)
	pos := firstFrame + (steps * frameNr)
	message := "to" + strconv.FormatUint(uint64(pos), 10)
	readAndWriteToSerial(message)
	WaitForMotor()
}

func readAndWriteToSerial(message string) {
	writeToSerial(message)
	message = readSerial()
	fmt.Println(message)
}

func writeToSerial(message string) {
	//message = validateMessage(message)
	_, err := selectedPort.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Sent steps")
}

func validateMessage(message string) string {
	if message[len(message)-1:] == "\n" {
		return message
	}
	return message + "\n"
}

func buildMessage(steps uint32) string {
	stepsString := strconv.FormatUint(uint64(steps), 10)
	return stepsString
}

func readSerial() string {
	buff := make([]byte, 300)
	var message string
	for {
		n, err := selectedPort.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		message = message + string(buff[:n])

		//Break after message is finisched
		if strings.Contains(string(buff[:n]), "\n") {
			//fmt.Printf("Message ist: %v", message)
			break
		}
	}
	return message
}

func InitSlate(cal bool, ready chan bool) {
	if cal {
		CalibrateSlate()
	}
	MoveToStartPosition()
	ready <- true
}

func CalibrateSlate() {
	fmt.Println("calibrating transport mechanism. please wait")

	readAndWriteToSerial("c")
	WaitForMotor()

}

func WaitForMotor() {
	channel := make(chan bool, 1)
	go waitForFinishMessage(channel)

	handleTimeOut(channel)
}

func handleTimeOut(channel chan bool) {
	select {
	case <-channel:
		return
	case <-time.After(time.Second * timeout):
		log.Fatalf("Something went wrong: Motor did not finish in %v seconds.\n"+
			"you can change the timeout duration in the config.yaml file", timeout.String())
		return
	}
}

func waitForFinishMessage(channel chan bool) {
	//fmt.Println("Waiting for Motor")
	var message string
	for {
		if strings.Contains(message, "finished") {
			channel <- true
			return
		}
		message = readSerial()
		time.Sleep(30 * time.Millisecond)
	}
}
