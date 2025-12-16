package ebitest

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseDumpFoler = "_ebitest_dump/"
)

type Ebitest struct {
	game     *Game
	PingPong *PingPong

	ctxCancelFn context.CancelFunc
	endGameChan chan struct{}

	options options
}

type options struct {
	face            text.Face
	color           color.Color
	dumpErrorImages bool
}

type optionsFn func(*options)

// WithFace set's the default Face for when checking for texts
func WithFace(f text.Face) optionsFn {
	return func(o *options) {
		o.face = f
	}
}

// WithColor set's the default Color for when checking for texts
func WithColor(c color.Color) optionsFn {
	return func(o *options) {
		o.color = c
	}
}

// WithDumpErrorImages enables the option to output a custom image when a test fails
// that has the screen and the image that was tried to match in order to debug it
func WithDumpErrorImages() optionsFn {
	return func(o *options) {
		o.dumpErrorImages = true
	}
}

func Run(game ebiten.Game, opts ...optionsFn) *Ebitest {
	ctx, cfn := context.WithCancel(context.TODO())
	pingPong := NewPingPong()
	g := newGame(ctx, game, pingPong)
	endGameChan := make(chan struct{})
	go func() {
		ebiten.RunGame(g)
		endGameChan <- struct{}{}
	}()

	et := &Ebitest{
		game:        g,
		ctxCancelFn: cfn,
		PingPong:    pingPong,
		endGameChan: endGameChan,
	}

	op := options{}

	for _, ofn := range opts {
		ofn(&op)
	}
	et.options = op

	et.PingPong.Ping()

	if et.options.dumpErrorImages {
		os.RemoveAll(baseDumpFoler)
		os.MkdirAll(baseDumpFoler, 0777)
	}

	return et
}

// Close stops the underlying game
func (e *Ebitest) Close() {
	e.ctxCancelFn()
	<-e.endGameChan
	close(e.endGameChan)
}

// Should checks if selector(s) is present in the game and returns it
// s can be a: 'string', 'image.Image', '*ebiten.Image' and '*ebitest.Selector'
func (e *Ebitest) Should(t *testing.T, s interface{}) (*Selector, bool) {
	t.Helper()
	e.PingPong.Ping()
	sc := e.game.GetScreen()

	sel, ok := e.findSelector(sc, s)
	if !ok {
		msg := "selector not found"
		if e.options.dumpErrorImages {
			p := dumpErrorImages(sc, sel.Image())
			msg += "\nimage at: " + p
		}
		assert.Fail(t, msg)
		return nil, false
	}

	return sel, true
}

// ShouldNot checks if selector(s) is not present in the game
// s can be a: 'string', 'image.Image', '*ebiten.Image' and '*ebitest.Selector'
func (e *Ebitest) ShouldNot(t *testing.T, s interface{}) bool {
	t.Helper()
	e.PingPong.Ping()
	sc := e.game.GetScreen()

	sel, ok := e.findSelector(sc, s)
	if !ok {
		return true
	}

	msg := "selector found"
	if e.options.dumpErrorImages {
		p := dumpErrorImages(sc, sel.Image())
		msg += "\nimage at: " + p
	}
	assert.Fail(t, msg)
	return false
}

// Must checks if selector(s) is present in the game and returns it.
// If it's not present it'll fail the test
// s can be a: 'string', 'image.Image', '*ebiten.Image' and '*ebitest.Selector'
func (e *Ebitest) Must(t *testing.T, s interface{}) *Selector {
	t.Helper()
	e.PingPong.Ping()
	sc := e.game.GetScreen()

	sel, ok := e.findSelector(sc, s)
	if !ok {
		msg := "selector not found"
		if e.options.dumpErrorImages {
			p := dumpErrorImages(sc, sel.Image())
			msg += "\nimage at: " + p
		}
		require.Fail(t, msg)
		return nil
	}

	return sel
}

// MustNot checks if selector(s) is not present in the game.
// If it's present it'll fail the test
// s can be a: 'string', 'image.Image', '*ebiten.Image' and '*ebitest.Selector'
func (e *Ebitest) MustNot(t *testing.T, s interface{}) {
	t.Helper()
	e.PingPong.Ping()
	sc := e.game.GetScreen()

	sel, ok := e.findSelector(sc, s)
	if !ok {
		return
	}

	msg := "selector found"
	if e.options.dumpErrorImages {
		p := dumpErrorImages(sc, sel.Image())
		msg += "\nimage at: " + p
	}
	require.Fail(t, msg)
	return
}

// KeyTap taps all the keys at once
func (e *Ebitest) KeyTap(keys ...ebiten.Key) {
	if len(keys) == 0 {
		return
	}
	e.PingPong.KeyTapPing(Ball{KeyTap: BallKeyTap{Keys: keys}})
}

// getSelector converts s to the right Selector initialization
func (e *Ebitest) getSelector(s interface{}) *Selector {
	switch v := s.(type) {
	case string:
		return NewFromText(v, e.options.face, e.options.color)
	case image.Image:
		return NewFromImage(v)
	case *ebiten.Image:
		return NewFromImage(ebitenImageToImage(v))
	case *Selector:
		return NewFromImage(v.Image())
	default:
		panic(fmt.Sprintf("Invalid Selector of type %T, the supported ones are: 'string', 'image.Image', '*ebiten.Image' and '*ebitest.Selector'", s))
	}
}

// findSelector returns a Selector from ss if found
func (e *Ebitest) findSelector(sc image.Image, ss interface{}) (*Selector, bool) {
	sel := e.getSelector(ss)

	sx := sc.Bounds().Dx()
	sy := sc.Bounds().Dy()
	for x := range sx {
		for y := range sy {
			if hasImageAt(sc, sel.Image(), x, y) {
				selx := sel.Image().Bounds().Dx()
				sely := sel.Image().Bounds().Dy()

				sel.rect = image.Rect(x, y, x+selx, y+sely)
				sel.PingPong = e.PingPong
				return sel, true
			}
		}
	}

	return sel, false
}

// dumpErrorImages dumps a composition of the 2 images into 1 so it displays
// what was checked
func dumpErrorImages(s, i image.Image) string {
	sb := s.Bounds()
	ib := i.Bounds()
	x := sb.Dx() + ib.Dx()
	y := sb.Dy()
	img := image.NewRGBA(image.Rect(0, 0, x, y))

	draw.Draw(img, sb, s, image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(sb.Dx(), 0, x, ib.Dy()), i, image.Point{}, draw.Over)

	u, _ := uuid.NewV7()

	ip := filepath.Join(baseDumpFoler, u.String()+".png")
	writeImage(ip, img)

	wd, _ := os.Getwd()
	return filepath.Join(wd, ip)
}

// writeImage writes on the path the image i
func writeImage(path string, i image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, i); err != nil {
		log.Fatal(err)
	}
}

// hasImageAt checks if the image sub is image i at the ix, iy
func hasImageAt(i, sub image.Image, ix, iy int) bool {
	sx, sy := sub.Bounds().Dx(), sub.Bounds().Dy()
	for x := range sx {
		for y := range sy {
			ic := toNRGBA(i.At(ix+x, iy+y))
			sc := toNRGBA(sub.At(x, y))
			sr, sg, sb, sa := sc.RGBA()

			// If the source it's transparent we ignore it
			// we want to only compare colors so we consider
			// it as good
			if sa == 0 || (sr == 0 && sg == 0 && sb == 0) {
				continue
			}

			if !equalColors(sc, ic) {
				return false
			}
		}
	}
	return true
}

// equalColors checks if c1 and c2 have the same RGB
func equalColors(c1, c2 color.Color) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2
}

// toNRGBA convers a pre-multiplied alpha color to a non pre-multiplied alpha one
func toNRGBA(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return color.NRGBA{0, 0, 0, 0}
	}

	// Since color.Color is alpha pre-multiplied, we need to divide the
	// RGB values by alpha again in order to get back the original RGB.
	r *= 0xffff
	r /= a
	g *= 0xffff
	g /= a
	b *= 0xffff
	b /= a

	return color.NRGBA{uint8(r / 65535), uint8(g / 65535), uint8(b / 65535), 255}
}
