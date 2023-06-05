package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/kbinani/screenshot"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var tmpFilename string

func main() {
	a := app.New()
	w := a.NewWindow("Screenshots")

	var image *canvas.Image

	image = canvas.NewImageFromFile(tmpFilename)
	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(200, 250))
	image.Refresh()

	btn := widget.NewButton("Take Screenshot", func() {
		takeScreenshot(image, w)
	})

	btn2 := widget.NewButton("Save to Pictures", func() {
		saveToPictures(w)
	})

	btns := container.NewGridWithColumns(2, btn, btn2)
	w.SetContent(
		container.NewBorder(nil, btns, nil, nil, image),
	)

	w.Resize(fyne.NewSize(610, 420))
	w.ShowAndRun()
}

func takeScreenshot(image *canvas.Image, w fyne.Window) {
	playSound()

	w.Hide()
	time.Sleep(time.Millisecond * 300)
	w.Show()
	os.Remove(tmpFilename)

	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("screenshot-")

	file, err := os.CreateTemp("", fileName)

	defer file.Close()

	png.Encode(file, img)

	tmpFilename = file.Name()

	image.File = tmpFilename
	image.Refresh()

}

func saveToPictures(w fyne.Window) {
	input, err := ioutil.ReadFile(tmpFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	now := time.Now()

	d := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
		if write == nil {
			return
		}
		write.Write(input)
		defer write.Close()
		return
	}, w)

	d.SetFileName(now.Format("Screenshot at 15:04:05") + ".png")

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	fileURI := storage.NewFileURI(dirname + "/Pictures")
	fileLister, _ := storage.ListerForURI(fileURI)
	d.SetLocation(fileLister)

	d.Show()
}

func playSound() {
	streamer, format, err := mp3.Decode(resourceCameraShutterClickMp3)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
