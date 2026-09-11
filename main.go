package main

import (
	"bufio"
	"cl/internal"
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	flagLoop := flag.Bool("loop", false, "restart the TUI continously, persistent TUI")
	flagHeight := flag.Int("height", 0, "maximum number of visible options (0 = dynamic terminal height)")
	flag.Parse()

	var items []string

	// read stdio
	scanner := bufio.NewScanner(os.Stdin)
	// append each entry to items variable
	for scanner.Scan() { items = append(items, scanner.Text()) }
	// output error to stderr and exit, if scanner returned error
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// exit if no items were given from stdin
	if len(items) == 0 { return }

	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open terminal: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	p := tea.NewProgram(
		internal.NewModel(items, *flagLoop, *flagHeight),
		tea.WithInput(tty),
	)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// only print after bubble tea exist
	// this means we're back on the normal terminal screen
	final := finalModel.(internal.Model)
	if !final.Quit {
		fmt.Println(final.Result)
	}
}
