package main

import (
	"github.com/hsluv/hsluv-go"
	"image"
	"image/color"
	"math"
)

var plasma struct {
	palette []color.NRGBA
}

func runPlasma(tick chan bool, pgmChan, pstChan chan *image.NRGBA) {
	step := 0
	for pgm := range tick {
		img := generatePlasma(step)
		step++
		if pgm {
			pgmChan <- img
		} else {
			pstChan <- img
		}
	}
}

func generatePlasma(step int) *image.NRGBA {
	if len(plasma.palette) != 360 {
		plasma.palette = make([]color.NRGBA, 360)
		for i, _ := range plasma.palette {
			var c color.NRGBA
			r, g, b := hsluv.HsluvToRGB(float64(i), 100, 50)
			c.R = byte(r * 255.0)
			c.G = byte(g * 255.0)
			c.B = byte(b * 255.0)
			c.A = 255
			plasma.palette[i] = c
		}
	}

	image := image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	var value float64
	b := image.Bounds()
	for r := b.Min.Y; r < b.Max.Y; r++ {
		for c := b.Min.X; c < b.Max.X; c++ {
			value = math.Sin(dist(c+step, r, 32.0, 128.8) / 8.0)
			value += math.Sin(dist(c, r, 16.0, 64.0) / 8.0)
			value += math.Sin(dist(c, r+step/7.0, 0.0, 0.0) / 7.0)
			value += math.Sin(dist(c, r, 192.0, 100.0) / 16.0)
			image.Set(c, r, plasma.palette[(int((4+value)*45)+step)%360])
		}
	}
	return image
}

func dist(y, x int, c, d float64) float64 {
	a := float64(y)
	b := float64(x)
	return math.Sqrt((a-c)*(a-c) + (b-d)*(b-d))
}
