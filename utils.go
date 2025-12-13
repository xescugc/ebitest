package ebitest

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func ebitenImageToImage(ei *ebiten.Image) image.Image {
	b := ei.Bounds()
	img := image.NewGray(image.Rect(0, 0, b.Dx(), b.Dy()))

	ix, iy := ei.Size()
	for x := range ix {
		for y := range iy {
			sc := ei.At(x, y)
			img.Set(x, y, sc)
		}
	}

	return img
}
