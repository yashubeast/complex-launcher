package main

import (
	"bufio"
	"cl/internal"
	"cl/sources"
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

	dmenuMode := !stdinIsTty()
	if dmenuMode {
		readStdio(&items)
	}

	// use default source, if no items were given
	// TODO: separate this with a --dmenu flag or something
	if len(items) == 0 {
		source := sources.PathSource{}

		itemsAsItem, err := source.List()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		items = make([]string, len(itemsAsItem))
		for i, item := range itemsAsItem {
			items[i] = item.Name
		}
	}

	// initiate tty

	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open terminal: %v\n", err)
		os.Exit(1)
	}
	defer tty.Close()

	// fetch config

	prefixes, err := internal.LoadPrefixConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// initiate the model

	p := tea.NewProgram(
		internal.NewModel(items, *flagLoop, *flagHeight, prefixes, dmenuMode),
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
	if dmenuMode && !final.Quit {
		fmt.Println(final.Result)
	}
}

func stdinIsTty() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func readStdio(items *[]string) {
	// read stdio
	scanner := bufio.NewScanner(os.Stdin)

	// append each entry to items variable
	for scanner.Scan() {
		*items = append(*items, scanner.Text())
	}

	// output error to stderr and exit, if scanner returned error
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
