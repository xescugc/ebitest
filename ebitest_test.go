package ebitest_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/xescugc/ebitest"
)

func TestGameUI(t *testing.T) {
	face, _ := loadFont(20)
	g := newGameUI()
	et := ebitest.Run(t, g,
		ebitest.WithFace(face),
		ebitest.WithColor(color.White),
		ebitest.WithDumpErrorImages(),
	)
	defer et.Close()

	text1 := "Click Me"
	text2 := "Clicked Me"

	t1s, _ := et.Should(text1)
	et.ShouldNot(text2)

	t1s.Click()
	t.Sleep(time.Second)

	et.ShouldNot(text1)
	et.Should(text2)
}
