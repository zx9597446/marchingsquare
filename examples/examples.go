package main

import (
	"github.com/zx9597446/marchingsquare"
	"image"
	"image/color"
	"image/png"
	"os"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func debugDraw(result []marchingsquare.Point, oldfile, newfile string, what color.Color) {
	old, err := os.Open(oldfile)
	defer old.Close()
	panicIfErr(err)
	oldimg, _, err := image.Decode(old)
	panicIfErr(err)
	b := oldimg.Bounds()
	newimg := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			newimg.Set(x, y, oldimg.At(x, y))
		}
	}
	for _, pt := range result {
		newimg.Set(pt.X, pt.Y, what)
	}
	file, err := os.Create(newfile)
	defer file.Close()
	panicIfErr(err)
	png.Encode(file, newimg)
}

func main() {
	ret := marchingsquare.ProcessWithFile("terrain.png", marchingsquare.TransparentTest)
	debugDraw(ret, "terrain.png", "new.png", color.RGBA{255, 0, 0, 255})
}
