package ebitest

import (
	"github.com/go-vgo/robotgo"
	"github.com/hajimehoshi/ebiten/v2"
)

// PingPong used to switch contexts between goroutines
type PingPong struct {
	ping chan struct{}
	pong chan struct{}

	clickPing chan Ball
	clickPong chan struct{}

	keyTapPing chan Ball
	keyTapPong chan struct{}
}

type Ball struct {
	X, Y   int
	KeyTap BallKeyTap
}

type BallKeyTap struct {
	Keys []ebiten.Key
}

// NewPingPong initialines a new PingPong
func NewPingPong() *PingPong {
	return &PingPong{
		ping: make(chan struct{}, 1),
		pong: make(chan struct{}, 1),

		clickPing: make(chan Ball, 1),
		clickPong: make(chan struct{}, 1),

		keyTapPing: make(chan Ball, 1),
		keyTapPong: make(chan struct{}, 1),
	}
}

// Ping sends a ping and waits for a Pong
func (pp *PingPong) Ping() {
	pp.ping <- struct{}{}
	<-pp.pong
}

// Pong will read on ping if any present and then
// will send to pong
func (pp *PingPong) Pong() {
	if len(pp.ping) != 0 {
		<-pp.ping
		pp.pong <- struct{}{}
	}
}

// ClickPing sends a ping and waits for a Pong
func (pp *PingPong) ClickPing(b Ball) {
	pp.clickPing <- b
	<-pp.clickPong
}

// ClickPong will read on ping if any present and then
// will send to pong
func (pp *PingPong) ClickPong(ttt *TicTacToe) {
	if len(pp.clickPing) != 0 {
		b := <-pp.clickPing
		robotgo.Move(b.X, b.Y)
		robotgo.Click("left")
		ttt.Tic()
		go func() {
			<-ttt.toe
			pp.clickPong <- struct{}{}
		}()
	}
}

// Ping sends a ping and waits for a Pong
func (pp *PingPong) KeyTapPing(b Ball) {
	pp.keyTapPing <- b
	<-pp.keyTapPong
}

// Pong will read on ping if any present and then
// will send to pong
func (pp *PingPong) KeyTapPong(ttt *TicTacToe, g *Game) {
	if len(pp.keyTapPing) != 0 {
		b := <-pp.keyTapPing
		keys := b.KeyTap.Keys
		key := ebitenToRobotgoKeys[keys[0]]
		args := make([]interface{}, 0, 0)
		for _, k := range keys[1:] {
			args = append(args, ebitenToRobotgoKeys[k])
		}
		g.keyTapKeys = keys
		robotgo.KeyTap(key, args...)
		ttt.Tic()
		go func() {
			<-ttt.toe
			g.keyTapKeys = make([]ebiten.Key, 0, 0)
			pp.keyTapPong <- struct{}{}
		}()
	}
}
