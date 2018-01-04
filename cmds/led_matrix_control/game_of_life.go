package main

import (
	"github.com/depili/vps-300p/life"
	"image"
	"image/color"
	"image/draw"
)

// Game of Life

var eff_gol struct {
	life *life.Life
	step int
}

func runGol(tick chan bool, pgmChan, pstChan chan *image.NRGBA) {
	eff_gol.life = life.NewLife(matrix_width, matrix_height)
	for pgm := range tick {
		eff_gol.life.Step()
		eff_gol.step++
		if eff_gol.step > 500 {
			eff_gol.life = life.NewLife(matrix_width, matrix_height)
			eff_gol.step = 0
		}
		mask := eff_gol.life.Image()
		img := image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
		draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
		draw.DrawMask(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 255, 255}},
			image.ZP, mask, image.ZP, draw.Over)
		if pgm {
			pgmChan <- img
		} else {
			pstChan <- img
		}
	}
}
