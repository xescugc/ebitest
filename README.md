# Ebitest

[![Go Reference](https://pkg.go.dev/badge/github.com/xescugc/ebitest.svg)](https://pkg.go.dev/github.com/xescugc/ebitest)

Ebitest is a lib to test Ebiten UI through inputs and asserting on what should be on the screen.

## Requirements

There are a few dependencies/requirements in order to run Ebitest at it's fullest.

To run the test **headless**, meaning without having the game open on a screen, you need to install [Xvfb](https://www.x.org/archive/X11R7.7/doc/man/man1/Xvfb.1.xhtml), but
that is only available on Linux (and X server) for others I did not investigate yet. With this then you can just do `xvfb-run go test ./...` (check the [Makefile](./Makefile)).

To fake all the events (click, scroll, mouse move) I use [robotgo](https://github.com/go-vgo/robotgo) so I would recommend checking their README for specific
requirements depending on the OS you have.

## Installation

With Go module support (Go 1.11+), just import:

```golang
import "github.com/xescugc/ebitest"
```

Otherwise, to install the ebiteset package, run the command:

```
go get github.com/xescugc/ebitest
```

## Usage

Simple API that has:
* `Should(s)` and `ShouldNot(s)`: Not stop execution when fails
* `Must(s)` and `MustNot(s)`: Stop execution if assertion fails

When asserting the `s` can be many things:
* `string`: To search that string on the screen (the Color and Face have to be provided on the initialization of Ebitest.Run)
* `image.Image`: Search that specific image on the screen
* `*ebiten.Image`: Searches that specific image
* `*ebitest.Selector`: Searches for the selector internal image

When using a positive assertion (`Should` or `Must`) they return also the `*ebitest.Selector` so then you can interact with it
like doing a `.Click()`.

Initialize Ebitest with `ebitest.Run(t, g)` with `t *testing.Test` and `g ebiten.Game`. A few extra options are available like:
* `WithFace|Color`: To set the default values when the using the assertions with a text value.
* `WithDumpErrorImages`: Which will generate an image when a test fail with the failed assertion on the folder `_ebitest_dump/`

If you need some extra interactions that are not implemented (yet) you can directly use [robotgo](https://github.com/go-vgo/robotgo),
but those may fail as they are not synchronized internally so I would recommend opening an issue and I'll add it.

## Example

This is a simple test in which there is a Game with just a button that when clicked switches the test in it from `Click Me` to `Clicked Me`

```golang
package ebitest_test

import (
	"image/color"
	"testing"

	"github.com/xescugc/ebitest"
)

func TestGameUI(t *testing.T) {
	face, _ := loadFont(20)
	g := newGameUI()
	et := ebitest.Run(g,
		ebitest.WithFace(face),
		ebitest.WithColor(color.White),
		ebitest.WithDumpErrorImages(),
	)
	defer et.Close()

	text1 := "Click Me"
	text2 := "Clicked Me"

	t1s, _ := et.Should(t, text1)
	et.ShouldNot(t, text2)

	t1s.Click()

	et.Should(t, text1)
	et.Should(t, text2)
}
```

The output of this test (that fails) is the following:

```
--- FAIL: TestGameUI (7.52s)
    ebitest_test.go:28: 
                Error Trace:    /home/xescugc/repos/ebitest/ebitest.go:113
                                                        /home/xescugc/repos/ebitest/ebitest_test.go:28
                Error:          selector not found
                                image at: _ebitest_dump/019b1537-1c60-7041-ad54-0297ea4b0eef.png
                Test:           TestGameUI
FAIL
FAIL    github.com/xescugc/ebitest      7.578s
FAIL
make: *** [Makefile:7: test] Error 1
```

And if you open the `_ebitest_dump/019b1537-1c60-7041-ad54-0297ea4b0eef.png` (on the current path) you see

<p align="center">
    <img src="docs/error_image.png" width=50% height=50%>
</p>

## Run it on a CI

If the CI has low resources (like GitHub Actions) it'll most likely fail (check `Known issues#2`) but you
can check what I install for it to run on the [`go.yml`](.github/workflows/go.yml)

## Known issues and Limitations

1/ You cannot have more than 1 test case

Basically you cannot run more than one test as even calling `Ebitest.Close()` there are some resources missing and you may get an error like

> panic: ebiten: NewImage cannot be called after RunGame finishes [recovered, repanicked]

2/ Some false positive/negative

Due to the nature of this test (the game is running on a goroutine) there may be the case in which an input is not registered by the game
so an expectation may randomly fail.

I kind of fixed it (100 consecutive test pass) using a custom [PingPong](./ping_pong.go) and [TicTacToe](./tic_tac_toe.go) that basically forces a context switch and synchronizes Input+Game.Update+Game.Draw but it still fails in low resource like GitHub [Actions](https://github.com/xescugc/ebitest/actions) for example.

## Plans

* Add more helpers for assertions (like animations)
* Add more inputs (potentially just port all the [robotgo](https://github.com/go-vgo/robotgo) lib) synchronized
* Others
