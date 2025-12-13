package ebitest_test

import (
	"image/color"
	"testing"

	"github.com/go-vgo/robotgo"
	"github.com/xescugc/ebitest"
	"github.com/xescugc/ebitest/testdata"
	"github.com/zeebo/assert"
)

func TestGameButton(t *testing.T) {
	face, _ := testdata.LoadFont(20)
	g := testdata.NewGame()
	et := ebitest.Run(t, g,
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
	text2 := "Clicked Me"

	t1s, _ := et.Should(text1)
	et.ShouldNot(text2)

	t1s.Click()

	et.ShouldNot(text1)
	et.Should(text2)
}
