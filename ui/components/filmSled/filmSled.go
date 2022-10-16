package filmSled

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyneTest/backend/serial"
	"fyneTest/backend/watcher"
	"fyneTest/ui/components/customWidgets"
	"fyneTest/ui/customTheme"
	"image/color"
	"strconv"
)

type FilmSled struct {
	FrameButtons   []*fyne.Container
	Sled           *fyne.Container
	FrameBtnColors []*canvas.Rectangle
}

var filmSled FilmSled

func NewFilmSled(fileWatcher *watcher.FileWatcher) FilmSled {
	grid := container.NewGridWithRows(1)
	var frameBtns []*fyne.Container
	var frameBtnColors []*canvas.Rectangle
	for i := 0; i < 6; i++ {
		btn, btnColor := createFilmSledBtns(i + 1)
		frameBtns = append(frameBtns, btn)
		frameBtnColors = append(frameBtnColors, btnColor)
		grid.Add(btn)
	}
	filmSled = FilmSled{
		FrameButtons:   frameBtns,
		Sled:           grid,
		FrameBtnColors: frameBtnColors,
	}

	fileWatcher.SetCurFrame = setCurFrame

	return filmSled
}

func (sled *FilmSled) UpdateActiveFrames(amt int) {

	for i, button := range sled.FrameBtnColors {
		if i < amt {
			button.FillColor = customTheme.MyTheme{}.Color("primary", 1)
			//button.FillColor = color.RGBA{
			//	R: 70,
			//	G: 40,
			//	B: 40,
			//	A: 255,
			//}
		} else {
			button.FillColor = color.RGBA{
				R: 40,
				G: 40,
				B: 40,
				A: 255,
			}
		}
		button.Refresh()
	}
}

func createFilmSledBtns(i int) (*fyne.Container, *canvas.Rectangle) {
	btn := widget.NewButton(strconv.Itoa(i), func() {
		customWidgets.NewConfirm("Move to Frame", "Are you sure you want to move the frame "+strconv.Itoa(i)+"?", func(b bool) {
			if b {
				serial.MoveToFrame(i)
			}
		})
		//filmSled.UpdateActiveFrames(i)
	})
	recColor := canvas.NewRectangle(
		color.RGBA{
			R: 40,
			G: 40,
			B: 40,
			A: 255,
		})

	return container.New(
		layout.NewMaxLayout(),
		recColor,
		btn), recColor
}

func setCurFrame(frame int) {
	curBtn := filmSled.FrameBtnColors[frame]
	curBtn.FillColor = color.RGBA{
		R: 244,
		G: 67,
		B: 54,
		A: 255,
	}
	curBtn.Refresh()
}
