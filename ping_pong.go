package ebitest

// PingPong used to switch contexts between goroutines
type PingPong struct {
	ping chan struct{}
	pong chan struct{}
}

// NewPingPong initialines a new PingPong
func NewPingPong() *PingPong {
	return &PingPong{
		ping: make(chan struct{}, 1),
		pong: make(chan struct{}, 1),
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
