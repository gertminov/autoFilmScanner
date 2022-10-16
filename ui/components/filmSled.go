package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyneTest/backend/config"
	"fyneTest/backend/watcher"
	"fyneTest/ui/components/filmSled"
)

func GetImagesSection(appConfig *config.Config, fileWatcher *watcher.FileWatcher) *fyne.Container {

	sled := filmSled.NewFilmSled(fileWatcher)
	//entry := customWidgets.NewNumericalEntry()
	slider := widget.NewSlider(1, 6)
	slider.OnChanged = func(val float64) {
		valInt := int(val)
		fmt.Println(val)
		sled.UpdateActiveFrames(valInt)
		appConfig.ImagesPerStrip = valInt
	}
	slider.SetValue(float64(appConfig.ImagesPerStrip))
	//entry.SetNumber(appConfig.ImagesPerStrip)

	imagesPerStrip := container.NewVBox(
		widget.NewLabel("Images per Strip"),
		slider,
	)

	box := container.NewVBox(
		imagesPerStrip,
		//slider,
		sled.Sled,
	)
	return container.New(
		layout.NewPaddedLayout(),
		box)
}
