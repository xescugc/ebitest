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

// HasTic checks if Tic has content
func (ttt *TicTacToe) HasTic() bool {
	return len(ttt.tic) != 0
}

// Tac second
func (ttt *TicTacToe) Tac() {
	if ttt.HasTic() {
		<-ttt.tic
		ttt.tac <- struct{}{}
	}
}

// HasTac checks if Tac has content
func (ttt *TicTacToe) HasTac() bool {
	return len(ttt.tac) != 0
}

// Toe third
func (ttt *TicTacToe) Toe() {
	if ttt.HasTac() {
		<-ttt.tac
		ttt.toe <- struct{}{}
	}
}

// HasToe checks if Toe has content
func (ttt *TicTacToe) HasToe() bool {
	return len(ttt.toe) != 0
}
