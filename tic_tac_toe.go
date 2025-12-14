package ebitest

// TicTacToe is a synchronizer of 3 channesl
type TicTacToe struct {
	tic chan struct{}
	tac chan struct{}
	toe chan struct{}
}

// NewTicTacToe returns a new TicTacToe
func NewTicTacToe() *TicTacToe {
	return &TicTacToe{
		tic: make(chan struct{}, 1),
		tac: make(chan struct{}, 1),
		toe: make(chan struct{}, 1),
	}
}

// Tic first
func (ttt *TicTacToe) Tic() {
	ttt.tic <- struct{}{}
}

// Tac second
func (ttt *TicTacToe) Tac() {
	if len(ttt.tic) != 0 {
		<-ttt.tic
		ttt.tac <- struct{}{}
	}
}

// Toe third
func (ttt *TicTacToe) Toe() {
	if len(ttt.tac) != 0 {
		<-ttt.tac
		ttt.toe <- struct{}{}
	}
}
