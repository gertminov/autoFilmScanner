package vuescan

import (
	"github.com/micmonay/keybd_event"
	"log"
)

var kb keybd_event.KeyBonding
var dia bool

func InitKeyboard(isDia bool) {
	var err error
	dia = isDia

	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatal(err)
	}
}

func Scan() {
	if dia {
		preview()
	}
	scan()
}

func scan() {
	kb.SetKeys(keybd_event.VK_N)
	addShiftAndLaunchShortCut()
}

func addShiftAndLaunchShortCut() {
	kb.HasCTRL(true)

	err := kb.Launching()
	if err != nil {
		log.Fatal(err)
	}
}

func preview() {
	kb.SetKeys(keybd_event.VK_I)
	addShiftAndLaunchShortCut()
}
