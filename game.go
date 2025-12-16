package ebitest

import (
	"context"
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	mxScreen sync.RWMutex
	screen   image.Image

	game ebiten.Game
	ctx  context.Context

	pingPong *PingPong

	keyTapKeys []ebiten.Key

	clickTTT  *TicTacToe
	keyTapTTT *TicTacToe
}

func newGame(ctx context.Context, g ebiten.Game, pp *PingPong) *Game {
	return &Game{
		game:      g,
		ctx:       ctx,
		pingPong:  pp,
		clickTTT:  NewTicTacToe(),
		keyTapTTT: NewTicTacToe(),
	}
}

func (g *Game) GetScreen() image.Image {
	g.mxScreen.Lock()
	defer g.mxScreen.Unlock()

	return g.screen
}

func (g *Game) SetScreen(s image.Image) {
	g.mxScreen.Lock()
	defer g.mxScreen.Unlock()

	g.screen = s
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.game.Layout(outsideWidth, outsideHeight)
}

// Update implements Game.
func (g *Game) Update() error {
	select {
	case <-g.ctx.Done():
		return ebiten.Termination
	default:
	}
	// Check for click
	if g.clickTTT.HasTic() {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.clickTTT.Tac()
		}
	}

	// Check for Pressed
	if g.keyTapTTT.HasTic() {
		pressed := true
		for _, k := range g.keyTapKeys {
			pressed = pressed && ebiten.IsKeyPressed(k)
		}
		if pressed {
			g.keyTapTTT.Tac()
		}
	}

	return g.game.Update()
}

// Draw implements Ebiten's Draw method.
func (g *Game) Draw(screen *ebiten.Image) {
	g.game.Draw(screen)

	g.SetScreen(ebitenImageToImage(screen))

	g.pingPong.Pong()
	g.pingPong.ClickPong(g.clickTTT)
	g.pingPong.KeyTapPong(g.keyTapTTT, g)

	g.clickTTT.Toe()
	g.keyTapTTT.Toe()
}
