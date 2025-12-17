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
	baseDumpFoler    = "_ebitest_dump/"
	findAllSelectors = true
)

var (
	emptyRec image.Rectangle
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
			p := dumpErrorImages(sc, sel)
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
		p := dumpErrorImages(sc, sel)
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
			p := dumpErrorImages(sc, sel)
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
		p := dumpErrorImages(sc, sel)
		msg += "\nimage at: " + p
	}
	require.Fail(t, msg)
}

// GetAll returns all the repeated instances of s or none if nothing is found
func (e *Ebitest) GetAll(s interface{}) []*Selector {
	sc := e.game.GetScreen()
	sels, _ := e.findSelectors(sc, s, findAllSelectors)

	return sels
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
	case *ebiten.Image:
		return NewFromImage(ebitenImageToImage(v))
	case image.Image:
		return NewFromImage(v)
	case *Selector:
		return NewFromImage(v.Image())
	default:
		panic(fmt.Sprintf("Invalid Selector of type %T, the supported ones are: 'string', 'image.Image', '*ebiten.Image' and '*ebitest.Selector'", s))
	}
}

// findSelector returns a Selector from ss if found. `all` will basically mean it'll return all of them
func (e *Ebitest) findSelector(sc image.Image, ss interface{}) (*Selector, bool) {
	sels, sel := e.findSelectors(sc, ss, !findAllSelectors)
	if len(sels) == 0 {
		return sel, false
	}
	return sels[0], true
}

// findSelector returns a Selector from ss if found. `all` will basically mean it'll return all of them
func (e *Ebitest) findSelectors(sc image.Image, ss interface{}, all bool) ([]*Selector, *Selector) {
	selectors := make([]*Selector, 0)
	bsel := e.getSelector(ss)

	sx := sc.Bounds().Dx()
	sy := sc.Bounds().Dy()
	for x := range sx {
		for y := range sy {
			if hasImageAt(sc, bsel.Image(), x, y) {
				sel := NewFromImage(bsel.Image())
				selx := sel.Image().Bounds().Dx()
				sely := sel.Image().Bounds().Dy()

				sel.rect = image.Rect(x, y, x+selx, y+sely)
				sel.PingPong = e.PingPong
				selectors = append(selectors, sel)
				if !all {
					return selectors, bsel
				}
			}
		}
	}

	return selectors, bsel
}

// hasImageAt checks if the image sub is image i at the ix, iy
func hasImageAt(i, sub image.Image, ix, iy int) bool {
	sx, sy := sub.Bounds().Dx(), sub.Bounds().Dy()
	for x := range sx {
		for y := range sy {
			sc := sub.At(x, y)
			ic := i.At(ix+x, iy+y)

			nsc := sc.(color.NRGBA)
			// If the source it's transparent we ignore it
			// we want to only compare colors so we consider
			// it as good
			if nsc.A != 255 {
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

// dumpErrorImages dumps a composition of the 2 images into 1 so it displays
// what was checked
func dumpErrorImages(s image.Image, sel *Selector) string {
	i := sel.Image()
	sb := s.Bounds()
	ib := i.Bounds()
	x := sb.Dx() + ib.Dx()
	y := sb.Dy()
	img := image.NewRGBA(image.Rect(0, 0, x, y))

	draw.Draw(img, sb, s, image.Point{}, draw.Over)
	draw.Draw(img, image.Rect(sb.Dx(), 0, x, ib.Dy()), i, image.Point{}, draw.Over)

	if sel.Rec() != emptyRec {
		drawRectangle(img, sel.Rec(), 2)
	}

	u, _ := uuid.NewV7()

	ip := filepath.Join(baseDumpFoler, u.String()+".png")
	writeImage(ip, img)

	wd, _ := os.Getwd()
	return filepath.Join(wd, ip)
}

// drawRectangle will draw in the image(img) the rectangel(rec) with thiknes
func drawRectangle(img *image.RGBA, rec image.Rectangle, thickness int) {
	col := color.RGBA{255, 0, 0, 255}

	for t := 0; t < thickness; t++ {
		// draw horizontal lines
		for x := rec.Min.X; x <= rec.Max.X; x++ {
			img.Set(x, rec.Min.Y+t, col)
			img.Set(x, rec.Max.Y-t, col)
		}
		// draw vertical lines
		for y := rec.Min.Y; y <= rec.Max.Y; y++ {
			img.Set(rec.Min.X+t, y, col)
			img.Set(rec.Max.X-t, y, col)
		}
	}
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
