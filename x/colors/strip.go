package colors

import (
	"fmt"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/comp"
)

var StripCmd = &bonzai.Cmd{
	Name:  `strip`,
	Usage: `strip`,
	Short: `Print 24 bit colors in strips`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Call: func(x *bonzai.Cmd, args ...string) error {
		Strip()
		return nil
	},
}

func setBackgroundColor(r, g, b int) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

func resetOutput() {
	fmt.Print("\x1b[0m\n")
}

func rainbowColor(h int) (int, int, int) {
	var r, g, b int
	h = h % 256

	switch {
	case h < 43:
		r, g, b = 255, h*255/43, 0
	case h < 86:
		r, g, b = 255-(h-43)*255/43, 255, 0
	case h < 129:
		r, g, b = 0, 255, (h-86)*255/43
	case h < 172:
		r, g, b = 0, 255-(h-129)*255/43, 255
	case h < 215:
		r, g, b = (h-172)*255/43, 0, 255
	default:
		r, g, b = 255, 0, 255-(h-215)*255/43
	}
	return r, g, b
}

func Strip() {
	for i := 0; i <= 127; i++ {
		setBackgroundColor(i, 0, 0)
		fmt.Print(" ")
	}
	resetOutput()
	for i := 255; i >= 128; i-- {
		setBackgroundColor(i, 0, 0)
		fmt.Print(" ")
	}
	resetOutput()

	for i := 0; i <= 127; i++ {
		setBackgroundColor(0, i, 0)
		fmt.Print(" ")
	}
	resetOutput()
	for i := 255; i >= 128; i-- {
		setBackgroundColor(0, i, 0)
		fmt.Print(" ")
	}
	resetOutput()

	for i := 0; i <= 127; i++ {
		setBackgroundColor(0, 0, i)
		fmt.Print(" ")
	}
	resetOutput()
	for i := 255; i >= 128; i-- {
		setBackgroundColor(0, 0, i)
		fmt.Print(" ")
	}
	resetOutput()

	for i := 0; i <= 127; i++ {
		r, g, b := rainbowColor(i)
		setBackgroundColor(r, g, b)
		fmt.Print(" ")
	}
	resetOutput()
	for i := 255; i >= 128; i-- {
		r, g, b := rainbowColor(i)
		setBackgroundColor(r, g, b)
		fmt.Print(" ")
	}
	resetOutput()
}
