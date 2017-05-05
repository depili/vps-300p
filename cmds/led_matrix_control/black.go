package main

import (
	"image"
	"image/color"
	"image/draw"
)

var eff_black struct {
	img *image.NRGBA
}

func runBlack(tick chan bool, pgmChan, pstChan chan *image.NRGBA) {
	eff_black.img = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(eff_black.img, eff_black.img.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	for pgm := range tick {
		if pgm {
			pgmChan <- eff_black.img
		} else {
			pstChan <- eff_black.img
		}
	}
}
