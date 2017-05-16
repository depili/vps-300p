package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/matrix"
	"github.com/depili/vps-300p/mixer"
	"github.com/depili/vps-300p/patterns"
	"github.com/hsluv/hsluv-go"
	"image"
	"image/color"
	"image/draw"
)

var images struct {
	PGM  *image.NRGBA
	PST  *image.NRGBA
	Out  *image.NRGBA
	Key1 *image.RGBA
}

func imageWorker(pgmChan, pstChan chan *image.NRGBA, key1Chan chan [][]bool) {
	m := matrix.Init(Options.Matrix, matrix_height, matrix_width)
	fmt.Printf("Connected to the matrix: %s\n", Options.Matrix)
	defer m.Close()

	black := color.RGBA{0, 0, 0, 255}
	transparent := color.RGBA{0, 0, 0, 0}

	images.PGM = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.PGM, images.PGM.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	images.PST = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.PST, images.PST.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	images.Out = image.NewNRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.Out, images.Out.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	images.Key1 = image.NewRGBA(image.Rect(0, 0, matrix_width, matrix_height))
	draw.Draw(images.Key1, images.Key1.Bounds(), &image.Uniform{transparent}, image.ZP, draw.Src)

	for {
		images.PGM = <-pgmChan
		images.PST = <-pstChan
		key1mask := <-key1Chan

		state := <-stateChan

		var alpha image.Image

		// Determinate the alpha for the transition
		if state.Type == mixer.TransMix {
			// Mix transition
			alpha = image.NewUniform(color.RGBA{0, 0, 0, byte(state.Value / 4)})
			// alpha_inv := image.NewUniform(color.RGBA{0, 0, 0, 255 - byte(state.Value/4)})
		} else {
			// Wipe transition

			v := state.Value
			r := state.PatternRev

			switch state.Pattern {
			case 0:
				alpha = patterns.NewPattern000(matrix_width, matrix_height, v, r)
			case 1:
				alpha = patterns.NewPattern001(matrix_width, matrix_height, v, r)
			case 2:
				alpha = patterns.NewPattern002(matrix_width, matrix_height, v, r)
			case 3:
				alpha = patterns.NewPattern003(matrix_width, matrix_height, v, r)
			case 4:
				alpha = patterns.NewPattern004(matrix_width, matrix_height, v, r)
			case 5:
				alpha = patterns.NewPattern005(matrix_width, matrix_height, v, r)
			case 6:
				alpha = patterns.NewPattern006(matrix_width, matrix_height, v, r)
			case 41:
				alpha = patterns.NewPattern041(matrix_width, matrix_height, v, r)
			default:
				alpha = patterns.NewPattern000(matrix_width, matrix_height, v, r)
			}
		}

		if state.Layers&mixer.LayerBKGD != 0 && state.Value != 0 {
			// Background transition in progress
			draw.Draw(images.Out, images.Out.Bounds(), images.PGM, image.ZP, draw.Src)
			draw.DrawMask(images.Out, images.Out.Bounds(), images.PST,
				image.ZP, alpha, image.ZP, draw.Over)
		} else {
			draw.Draw(images.Out, images.Out.Bounds(), images.PGM, image.ZP, draw.Src)
		}

		// Keyer 1
		hue := float64(state.Key1State.MattHue) / 1023.0 * 360.0
		sat := float64(state.Key1State.MattSat) / 1023.0 * 100.0
		lum := float64(state.Key1State.MattLum) / 1023.0 * 100.0
		matt_r, matt_g, matt_b := hsluv.HsluvToRGB(hue, sat, lum)
		matt := color.RGBA{byte(matt_r * 255.0), byte(matt_g * 255.0), byte(matt_b * 255.0), 255}
		key1Insert := &image.Uniform{matt}

		fill := color.RGBA{0, 0, 0, 255}
		if state.Layers&mixer.LayerKey1 != 0 && state.Value != 0 {
			// Key1 transition
			if state.Key1 {
				fill = color.RGBA{0, 0, 0, byte(255 - state.Value/4)}
			} else {
				fill = color.RGBA{0, 0, 0, byte(state.Value / 4)}
			}
			for r := range key1mask {
				for c := range key1mask[0] {
					if key1mask[r][c] {
						images.Key1.Set(c, r, fill)
					} else {
						images.Key1.Set(c, r, transparent)
					}
				}
			}
		} else if state.Key1 {
			for r := range key1mask {
				for c := range key1mask[0] {
					if key1mask[r][c] {
						images.Key1.Set(c, r, fill)
					} else {
						images.Key1.Set(c, r, transparent)
					}
				}
			}
		} else {
			draw.Draw(images.Key1, images.Key1.Bounds(), &image.Uniform{transparent},
				image.ZP, draw.Src)
		}

		draw.DrawMask(images.Out, images.Out.Bounds(), key1Insert,
			image.ZP, images.Key1, image.ZP, draw.Over)

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
