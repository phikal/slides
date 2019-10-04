package main

import (
	"fmt"
	"image"
	"os"
)

var img image.Image

type Image struct{}

func (r *Image) Set(val string) {
	file, err := os.Open(val)
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't open image file: ", err)
		return
	}
	defer file.Close()

	img, _, err = image.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while decoding image: ", err)
	}
}
func (r *Image) Reset() { img = nil }
func (r *Image) Push()  {}

func printImage() {
	var scale, xoff, yoff float64
	bounds := img.Bounds()
	iheight := bounds.Max.Y - bounds.Min.Y
	iwidth := bounds.Max.X - bounds.Min.X
	if fill != (iwidth/width > iheight/height) {
		scale = float64(width-2*padding) / float64(iwidth)
		yoff += (float64(iheight)*scale - float64(height-2*padding)) / 2
	} else {
		scale = float64(height-2*padding) / float64(iheight)
		xoff += (float64(iwidth)*scale - float64(width-2*padding)) / 2
	}
	fmt.Printf("%d %d 8 [%g 0 0 %g %g %g] {<",
		iwidth, iheight, 1/scale, 1/scale,
		(-float64(padding)+xoff)/scale,
		(-float64(padding)+yoff)/scale)
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Printf("%02x%02x%02x", r/0x101, g/0x101, b/0x0101)
		}
	}
	fmt.Println(">} false 3 colorimage showpage")
}
