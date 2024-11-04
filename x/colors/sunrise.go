package colors

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai"
)

var (
	p                  = 3.14
	i          float64 = 0
	SunriseCmd         = &bonzai.Cmd{
		Name:    `sunrise`,
		Usage:   `sunrise`,
		Short:   `Showcase all colors of terminal`,
		Comp:    comp.Cmds,
		Cmds:    []*bonzai.Cmd{},
		MaxArgs: 1,
		Call: func(x *bonzai.Cmd, args ...string) error {
			var ms int64 = 10
			if len(args) > 0 {
				ms, _ = strconv.ParseInt(args[0], 10, 64)
			}
			Sunrise(time.Duration(ms) * time.Millisecond)
			return nil
		},
	}
)

func reset() {
	exec.Command("reset").Run()       // Reset terminal
	fmt.Print("\033[H\033[2J\033[0m") // Clear terminal
	fmt.Print("\033[?25h")            // Show cursor
	fmt.Print(
		"\033[?1049l",
	) // Restore original screen on exit
	fmt.Print("\033[?25h") // Show cursor on exit
	os.Exit(0)
}

// Sunrise prints a wall of all colors in the terminal. The speed of the
// animation can be adjusted with the speed parameter.
func Sunrise(speed time.Duration) {
	// Set up signal handling to catch SIGINT and SIGTERM
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel
		reset()
	}()

	// Hide cursor and disable echo
	fmt.Print("\033[?25l")   // Hide cursor
	fmt.Print("\033[?1049h") // Switch to alternate screen buffer

	for {
		i += 0.04
		r := int(128 + 127*math.Sin(i))
		g := int(128 + 127*math.Sin(i+p*(1.0/3)))
		b := int(128 + 127*math.Sin(i+p*(2.0/3)))

		fmt.Printf("\033[48;2;%d;%d;%dm\n", r, g, b)
		time.Sleep(speed) // Sleep for 10ms
	}
}
