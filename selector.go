package ebitest

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Selector represents the thing you are searching for
type Selector struct {
	img  image.Image
	rect image.Rectangle

	PingPong *PingPong
}

// NewFromText crates a new Selector from a txt
func NewFromText(txt string, f text.Face, c color.Color) *Selector {
	top := &text.DrawOptions{}
	top.ColorScale.ScaleWithColor(c)

	x, y := text.Measure(txt, f, 0)
	rec := image.Rect(0, 0, int(x), int(y))

	imgp := image.NewPaletted(rec, color.Palette{color.Transparent})
	img := ebiten.NewImageFromImage(imgp)

	text.Draw(img, txt, f, top)

	return NewFromImage(ebitenImageToImage(img))
}

// NewFromImage creates a new Selector from an image
func NewFromImage(i image.Image) *Selector {
	return &Selector{
		img: i,
	}
}

// Click will click on the center of the Selectore
func (s *Selector) Click() {
	cx, cy := s.center()
	s.PingPong.ClickPing(Ball{X: cx, Y: cy})
}

// center returns the center of the selector
func (s *Selector) center() (int, int) {
	return s.rect.Min.X + (s.rect.Dx() / 2), s.rect.Min.Y + (s.rect.Dy() / 2)
}

// Rec returns the image.Rectangle of the image
func (s *Selector) Rec() image.Rectangle {
	return s.rect
}

// Image returns the underlying image
func (s *Selector) Image() image.Image {
	return s.img
}
