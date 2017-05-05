package main

import (
	"github.com/hsluv/hsluv-go"
	"image"
	"image/color"
	"math/rand"
)

var eff_flame struct {
	palette []color.RGBA
	buffer  [][]byte
}

func runFlame(tick chan bool, pgmChan, pstChan chan *image.NRGBA) {
	initFlame()
	for pgm := range tick {
		flameSeed()
		img := flameFill()
		if pgm {
			pgmChan <- img
		} else {
			pstChan <- img
		}
	}
}

func initFlame() {
	eff_flame.palette = make([]color.RGBA, 256)
	eff_flame.buffer = make([][]byte, matrix_height)
	for r, _ := range eff_flame.buffer {
		eff_flame.buffer[r] = make([]byte, matrix_width)
	}

	for i, _ := range eff_flame.palette {
		l := float64(i) / 256.0 * 100.0
		if l > 50 {
			l = 50
		}
		r, g, b := hsluv.HsluvToRGB(float64(i)/2.0, 100, l)
		color := color.RGBA{byte(r * 255.0), byte(g * 255.0), byte(b * 255.0), 255}
		eff_flame.palette[i] = color
	}
}

func flameSeed() {
	// Seed the bottom row
	for c, _ := range eff_flame.buffer[matrix_height-1] {
		eff_flame.buffer[matrix_height-1][c] = byte(rand.Float32() * 255.0)
	}
}

func flameFill() *image.NRGBA {
	for r, row := range eff_flame.buffer[0 : matrix_height-1] {
		value := 0
		for c, _ := range row {
			if y := (r + 1); y < matrix_height {
				value = int(eff_flame.buffer[y][c])
				if x := (c - 1); x >= 0 {
					value += int(eff_flame.buffer[y][x])
				}
				if x := c + 1; x < matrix_width {
					value += int(eff_flame.buffer[y][x])
				}
			}
			value *= 40
			value /= 129
			eff_flame.buffer[r][c] = byte(value)
		}
	}

	image := image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	b := image.Bounds()
	for r := b.Min.Y; r < b.Max.Y; r++ {
		for c := b.Min.X; c < b.Max.X; c++ {
			image.Set(c, r, eff_flame.palette[eff_flame.buffer[r][c]])
		}
	}
	return image
}
