package ebitest_test

import (
	"image/color"
	"testing"

	"github.com/go-vgo/robotgo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
	"github.com/xescugc/ebitest"
	"github.com/xescugc/ebitest/testdata"
)

func TestGameButton(t *testing.T) {
	face, _ := testdata.LoadFont(20)
	g := testdata.NewGame()
	et := ebitest.Run(g,
		ebitest.WithFace(face),
		ebitest.WithColor(color.White),
		ebitest.WithDumpErrorImages(),
	)
	defer et.Close()

	robotgo.Move(0, 0)
	robotgo.Click("left", true)

	assert.True(t, g.Clicked)

	et.PingPong.Ping()

	text1 := "Click Me"
	text1_2 := "Click Me 2"
	text2 := "Clicked Me"

	t1s, _ := et.Should(t, text1)
	t1_2s, _ := et.Should(t, text1_2)
	et.ShouldNot(t, text2)

	t1s.Click()

	et.Should(t, text1_2)
	et.Should(t, text2)

	t1_2s.Click()

	et.ShouldNot(t, text1)
	et.ShouldNot(t, text1_2)
	assert.Len(t, et.GetAll(text2), 2)

	//et.KeyTap(ebiten.KeyShift, ebiten.KeyI)
	et.KeyTap(ebiten.KeyI, ebiten.KeyShift)
	assert.True(t, g.ClickedShiftI)
}
