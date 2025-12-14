package ebitest

import (
	"github.com/go-vgo/robotgo"
)

// PingPong used to switch contexts between goroutines
type PingPong struct {
	ping chan struct{}
	pong chan struct{}

	clickPing chan Ball
	clickPong chan struct{}
}

type Ball struct {
	X, Y int
}

// NewPingPong initialines a new PingPong
func NewPingPong() *PingPong {
	return &PingPong{
		ping: make(chan struct{}, 1),
		pong: make(chan struct{}, 1),

		clickPing: make(chan Ball, 1),
		clickPong: make(chan struct{}, 1),
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

// Ping sends a ping and waits for a Pong
func (pp *PingPong) ClickPing(b Ball) {
	pp.clickPing <- b
	<-pp.clickPong
}

// Pong will read on ping if any present and then
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
