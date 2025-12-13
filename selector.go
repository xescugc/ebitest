package ebitest

import (
	"image"
	"image/color"

	"github.com/go-vgo/robotgo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Selector struct {
	img  image.Image
	rect image.Rectangle
}

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

func NewFromImage(i image.Image) *Selector {
	return &Selector{
		img: i,
	}
}

func (s *Selector) Click() {
	cx, cy := s.center()
	robotgo.Move(cx, cy)
	robotgo.Click("left", true)
}

func (s *Selector) center() (int, int) {
	return s.rect.Min.X + (s.rect.Dx() / 2), s.rect.Min.Y + (s.rect.Dy() / 2)
}

func (s *Selector) Rec() image.Rectangle {
	return s.rect
}

func (s *Selector) Image() image.Image {
	return s.img
}

func (s *Selector) Type(string) {
}

func (s *Selector) KeyTap(ks ...string) {
}
