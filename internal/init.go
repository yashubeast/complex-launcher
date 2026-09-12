package internal

import (
	"bufio"
	"cl/sources"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type Flags struct {
	Loop   bool
	Height int
}

func Run(config Config, flags Flags) error {
	items, dmenu, err := loadItems()
	if err != nil { return err }

	tty, err := os.Open("/dev/tty")
	if err != nil {
		return fmt.Errorf("open terminal: %w", err)
	}
	defer tty.Close()

	program := tea.NewProgram(
		NewModel(
			items,
			flags.Loop,
			flags.Height,
			config,
			dmenu,
		),
		tea.WithInput(tty),
	)

	finalModel, err := program.Run()
	if err != nil { return err }

	if dmenu {
		final := finalModel.(Model)
		if !final.Quit {
			fmt.Println(final.Result)
		}
	}
	return nil
}

func loadItems() ([]sources.Item, bool, error) {
	// load from dmenu if there is stdin
	dmenu := !stdinIsTty()
	if dmenu {
		items, err := readStdin()
		if err != nil { return nil, true, err }

		if len(items) > 0 { return items, true, nil }
	}

	// otherwise load default items
	items, err := defaultItems()
	return items, dmenu, err
}

func defaultItems() ([]sources.Item, error) {
	items, err := (sources.PathSource{}).List()
	if err != nil {
		return nil, fmt.Errorf("load Path items: %w", err)
	}

	return items, nil
}

func stdinIsTty() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func readStdin() ([]sources.Item, error) {
	var items []sources.Item

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		items = append(items, sources.Item{
			Name: scanner.Text(),
			Cmd: scanner.Text(),
		})
	}
	err := scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("read stdin: %w", err)
	}

	return items, nil
}
