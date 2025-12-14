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

	clickTTT *TicTacToe
}

func newGame(ctx context.Context, g ebiten.Game, pp *PingPong) *Game {
	return &Game{
		game:     g,
		ctx:      ctx,
		pingPong: pp,
		clickTTT: NewTicTacToe(),
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.clickTTT.Tac()
	}
	return g.game.Update()
}

// Draw implements Ebiten's Draw method.
func (g *Game) Draw(screen *ebiten.Image) {
	g.game.Draw(screen)

	g.SetScreen(ebitenImageToImage(screen))

	g.pingPong.Pong()
	g.pingPong.ClickPong(g.clickTTT)
	g.clickTTT.Toe()

	return
}
