package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"github.com/fogleman/imview"
	"github.com/nfnt/resize"
)

type Size struct {
	width  int
	height int
}
type MyImage struct {
	value *image.RGBA
}

func Width(i image.Image) int {
	return i.Bounds().Max.X - i.Bounds().Min.X
}

func Height(i image.Image) int {
	return i.Bounds().Max.Y - i.Bounds().Min.Y
}

//function to resize image and put it to a new resulting image. takes in a decoded image, the start point of which the image should be pasted on the new bg, the width and the height of the image
func (bgImg *MyImage) drawRaw(img image.Image, sp image.Point, width uint, height uint) {

	// height = 0
	// resize to width using Lanczos resampling and preserve aspect ratio by setting either width or height to 0
	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
	w := int(Width(resizedImg))
	h := int(Height(resizedImg))

	draw.Draw(bgImg, image.Rectangle{sp, image.Point{sp.X + w, sp.Y + h}}, resizedImg, image.ZP, draw.Src)
}

func (i *MyImage) Set(x, y int, c color.Color) {
	i.value.Set(x, y, c)
}

func (i *MyImage) ColorModel() color.Model {
	return i.value.ColorModel()
}

func (i *MyImage) Bounds() image.Rectangle {
	return i.value.Bounds()
}

func (i *MyImage) At(x, y int) color.Color {
	return i.value.At(x, y)
}

func main() {
	//open file
	file, err := os.Open("Bobo2.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//decode image into image.Image
	img, fortmaName, err := image.Decode(file)
	if err != nil {
		log.Println(fortmaName)
		log.Fatal(err)
	}

	out := MyImage{
		value: image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{400, 400}}),
	}
	out.drawRaw(img, image.Point{100, 100}, 180, 150)

	//imview just to give the image
	imview.Show(out.value)

}
