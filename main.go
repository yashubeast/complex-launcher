package main

import (
	"bufio"
	"cl/internal"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
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
		fmt.Fprintf(os.Stderr, "cannot open terminal: ", err)
		os.Exit(1)
	}
	defer tty.Close()

	p := tea.NewProgram(
		internal.NewModel(items),
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
