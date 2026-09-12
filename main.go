package main

import (
	"cl/internal"
	"flag"
	"fmt"
	"os"
)

func main() {
	flags := internal.Flags{}

	flag.BoolVar(&flags.Loop, "loop", false, "restart the TUI continously, persistent TUI")
	flag.IntVar(&flags.Height, "height", 0, "maximum number of visible options (0 = dynamic terminal height)")
	flag.Parse()

	err := internal.Run(flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
