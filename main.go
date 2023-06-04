package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var tmpFilename string = "tmp.png"

func main() {
	a := app.New()
	w := a.NewWindow("Screenshots")

	_, error := os.Stat(tmpFilename)

	var image *canvas.Image

	if !os.IsNotExist(error) {
		image = canvas.NewImageFromFile(tmpFilename)
		image.FillMode = canvas.ImageFillContain
		image.SetMinSize(fyne.NewSize(200, 250))
		image.Refresh()
	} else {
		image = canvas.NewImageFromResource(nil)
	}

	btn := widget.NewButton("Take Screenshot", func() {
		takeScreenshot(image)
	})

	btn2 := widget.NewButton("Save to Pictures", func() {
		saveToPictures()
	})

	btns := container.NewGridWithColumns(2, btn, btn2)
	w.SetContent(
		container.NewBorder(nil, btns, nil, nil, image),
	)

	w.Resize(fyne.NewSize(380, 220))
	w.ShowAndRun()
}

func takeScreenshot(image *canvas.Image) {
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("tmp.png")
	file, _ := os.Create(fileName)
	defer file.Close()
	png.Encode(file, img)

	image.File = "tmp.png"
	image.Refresh()
}

func saveToPictures() {
	input, err := ioutil.ReadFile("tmp.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now()

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	destinationFile := dirname + "/Pictures/" + now.Format("15:04:05") + ".png"

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}
