package main

import (
	"github.com/depili/go-rgb-led-matrix/bdf"
)

var eff_text struct {
	bitmap [][]bool
	row    int
	col    int
	length int
}

func initText(text, font string, x, y, length int) {
	eff_text.row = y
	eff_text.col = x
	eff_text.length = length
	f, err := bdf.Parse(font)
	if err != nil {
		panic(err)
	}
	eff_text.bitmap = f.TextBitmap(text)
}

func runText(tick chan bool, key1Chan, key2Chan chan [][]bool) {
	initText("Huomenta helsinki.hacklab.fi", "10x20.bdf", 0, 20, matrix_width)
	offset := 0
	for key := range tick {
		img := scrollText(offset)
		offset++
		if key {
			key1Chan <- img
		} else {
			key2Chan <- img
		}
	}
}

func scrollText(offset int) [][]bool {
	img := make([][]bool, matrix_height)
	for r := range img {
		img[r] = make([]bool, matrix_width)
	}
	scroll := getScrollBitmask(offset)

	for r := range scroll {
		for c := range scroll[r] {
			if scroll[r][c] {
				img[r+eff_text.row][c+eff_text.col] = true
			} else {
				img[r+eff_text.row][c+eff_text.col] = false
			}
		}
	}
	return img
}

func getScrollBitmask(offset int) [][]bool {
	s := len(eff_text.bitmap[0])
	if s < eff_text.length {
		return eff_text.bitmap
	}
	// Seamlessly loop the bitmap around
	offset = offset % s
	size := s - offset
	if size >= eff_text.length {
		return offsetBitmask(offset, eff_text.length)
	} else {
		ret := make([][]bool, len(eff_text.bitmap))
		for r, row := range offsetBitmask(offset, eff_text.length) {
			ret[r] = row
		}
		for r, row := range offsetBitmask(0, eff_text.length-size) {
			ret[r] = append(ret[r], row...)
		}
		return ret
	}
}

func offsetBitmask(offset, length int) [][]bool {
	ret := make([][]bool, len(eff_text.bitmap))
	for r, row := range eff_text.bitmap {
		end := offset + length
		if end > len(row) {
			end = len(row)
		}
		ret[r] = make([]bool, end-offset)
		for c, b := range row[offset:end] {
			ret[r][c] = b
		}
	}
	return ret
}
