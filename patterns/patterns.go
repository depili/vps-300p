package patterns

import (
	"image"
	"image/color"
)

func mask(pattern bool) color.Color {
	if pattern {
		return color.Alpha{255}
	} else {
		return color.Alpha{0}
	}
}

func reverse(x, y int, rect image.Rectangle, rev bool) (int, int) {
	if rev {
		x = rect.Dx() - x
		y = rect.Dy() - y
	}
	return x, y
}

// Pattern 000: horizontal wipe

type pattern000 struct {
	Rect  image.Rectangle
	State int
	Rev   bool
}

func (p *pattern000) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *pattern000) Bounds() image.Rectangle {
	return p.Rect
}

func (p *pattern000) At(x, y int) color.Color {
	x, y = reverse(x, y, p.Rect, p.Rev)
	return mask(x < (p.Rect.Dx() * p.State / 1023))
}

func NewPattern000(w, h, value int, r bool) *pattern000 {
	return &pattern000{image.Rect(0, 0, w, h), value, r}
}

// Pattern 001: vertical wipe

type pattern001 struct {
	Rect  image.Rectangle
	State int
	Rev   bool
}

func (p *pattern001) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *pattern001) Bounds() image.Rectangle {
	return p.Rect
}

func (p *pattern001) At(x, y int) color.Color {
	x, y = reverse(x, y, p.Rect, p.Rev)
	return mask(y < (p.Rect.Dy() * p.State / 1023))
}

func NewPattern001(w, h, value int, r bool) *pattern001 {
	return &pattern001{image.Rect(0, 0, w, h), value, r}
}

// PatternCorner wipe from a corner (2-5)

type patternCorner struct {
	Rect  image.Rectangle
	State int
	Rev   bool
	flipX bool
	flipY bool
}

func (p *patternCorner) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *patternCorner) Bounds() image.Rectangle {
	return p.Rect
}

func (p *patternCorner) At(x, y int) color.Color {
	x, y = reverse(x, y, p.Rect, p.Rev)
	if p.flipX {
		x = p.Rect.Dx() - x
	}
	if p.flipY {
		y = p.Rect.Dy() - y
	}

	m := y < (p.Rect.Dy()*p.State/1023) &&
		x < (p.Rect.Dx()*p.State/1023)
	return mask(m)
}

func NewPattern002(w, h, value int, r bool) *patternCorner {
	return &patternCorner{image.Rect(0, 0, w, h), value, r, false, false}
}

func NewPattern003(w, h, value int, r bool) *patternCorner {
	return &patternCorner{image.Rect(0, 0, w, h), value, r, true, false}
}

func NewPattern004(w, h, value int, r bool) *patternCorner {
	return &patternCorner{image.Rect(0, 0, w, h), value, r, true, true}
}

func NewPattern005(w, h, value int, r bool) *patternCorner {
	return &patternCorner{image.Rect(0, 0, w, h), value, r, false, false}
}

// Pattern 006: all corners wipe

type pattern006 struct {
	Rect  image.Rectangle
	State int
	Rev   bool
}

func (p *pattern006) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *pattern006) Bounds() image.Rectangle {
	return p.Rect
}

func (p *pattern006) At(x, y int) color.Color {
	s := p.State
	if p.Rev {
		s = 1023 - s
	}

	m := y < (p.Rect.Dy()*s/1023/2) &&
		x < (p.Rect.Dx()*s/1023/2)
	m = m || p.Rect.Dy()-y < (p.Rect.Dy()*s/1023/2) &&
		x < (p.Rect.Dx()*s/1023/2)
	m = m || y < (p.Rect.Dy()*s/1023/2) &&
		p.Rect.Dx()-x < (p.Rect.Dx()*s/1023/2)
	x, y = reverse(x, y, p.Rect, true)
	m = m || y < (p.Rect.Dy()*s/1023/2) &&
		x < (p.Rect.Dx()*s/1023/2)
	if p.Rev {
		return mask(!m)
	}
	return mask(m)
}

func NewPattern006(w, h, value int, r bool) *pattern006 {
	return &pattern006{image.Rect(0, 0, w, h), value, r}
}

// Pattern 041: towards center wipe

type pattern041 struct {
	Rect  image.Rectangle
	State int
	Rev   bool
}

func (p *pattern041) ColorModel() color.Model {
	return color.AlphaModel
}

func (p *pattern041) Bounds() image.Rectangle {
	return p.Rect
}

func (p *pattern041) At(x, y int) color.Color {
	s := p.State
	if !p.Rev {
		s = 1023 - s
	}
	m := y > (p.Rect.Dy()*s/1023/2) &&
		x > (p.Rect.Dx()*s/1023/2)
	m = m && p.Rect.Dy()-y > (p.Rect.Dy()*s/1023/2) &&
		x > (p.Rect.Dx()*s/1023/2)
	m = m && y > (p.Rect.Dy()*s/1023/2) &&
		p.Rect.Dx()-x > (p.Rect.Dx()*s/1023/2)
	x, y = reverse(x, y, p.Rect, true)
	m = m && y > (p.Rect.Dy()*s/1023/2) &&
		x > (p.Rect.Dx()*s/1023/2)
	if p.Rev {
		return mask(!m)
	}
	return mask(m)
}

func NewPattern041(w, h, value int, r bool) *pattern041 {
	return &pattern041{image.Rect(0, 0, w, h), value, r}
}
