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
	toggle := flag.Bool("toggle", false, "toggle the terminal window")
	flag.Parse()

	// load config
	config, configErr := internal.LoadConfig()
	if configErr != nil {
		fmt.Fprintln(os.Stderr, configErr)
		os.Exit(1)
	}

	if *toggle {
		err := internal.Toggle(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	err := internal.Run(config, flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
