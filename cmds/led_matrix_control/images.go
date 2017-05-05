package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/matrix"
	"image"
	"image/color"
	"image/draw"
)

var images struct {
	PGM *image.NRGBA
	PST *image.NRGBA
	Out *image.NRGBA
}

func imageWorker(pgmChan, pstChan chan *image.NRGBA) {
	m := matrix.Init(Options.Matrix, matrix_height, matrix_width)
	fmt.Printf("Connected to the matrix: %s\n", Options.Matrix)
	defer m.Close()

	black := color.RGBA{0, 0, 0, 255}
	images.PGM = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.PGM, images.PGM.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	images.PST = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.PST, images.PST.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	images.Out = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.Out, images.Out.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	for {
		select {
		case img := <-pgmChan:
			images.PGM = img
		case img := <-pstChan:
			images.PST = img
		}
		alpha := image.NewUniform(color.RGBA{0, 0, 0, eff_alpha})
		draw.Draw(images.Out, images.Out.Bounds(), images.PST, image.ZP, draw.Src)
		draw.DrawMask(images.Out, images.Out.Bounds(), images.PGM, image.ZP, alpha, image.ZP, draw.Over)

		b := images.Out.Bounds()
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				var c [3]byte
				p := images.Out.NRGBAAt(x, y)
				c[0] = p.R
				c[1] = p.G
				c[2] = p.B
				m.SetPixel(y, x, c)
			}
		}
		m.Send()
	}
}
