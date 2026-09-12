package internal

import (
	"cl/sources"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type Model struct {
	query          string
	selected       int
	items          []sources.Item
	matches        []fuzzy.Match

	// TODO: assign default items according to user's selected default source
	defaultItems   []sources.Item
	prefixSource   *sources.PrefixSource

	Result         string
	Quit           bool

	offset         int
	maxVisible     int // 0 = auto-detect from terminal height
	terminalHeight int
	terminalWidth  int

	flagLoop       bool
	dmenuMode      bool
}

func (m *Model) reset() {
	m.query = ""
	m.selected = 0
	m.offset = 0
	m.filter()
}

func NewModel(items []string, loop bool, maxVisible int, prefixes []sources.Prefix, dmenuMode bool) Model {
	idkItems := sources.GetItems(items)
	m := Model{
		items:    idkItems,
		selected: 0,
		flagLoop: loop,
		dmenuMode: dmenuMode,
		maxVisible: maxVisible,
		defaultItems: idkItems,
		prefixSource: &sources.PrefixSource{ Prefixes: prefixes },
	}
	m.filter()
	return m
}

func (m *Model) visibleCount() int {
	if m.maxVisible > 0 {
		return m.maxVisible
	}
	// query line + blank line
	n := m.terminalHeight - 2
	if n < 1 {
		return 1
	}
	return n
}

func (m *Model) keepSelectedVisible() {
	visible := m.visibleCount()
	if visible <= 0 { return }

	if m.selected < m.offset {
		m.offset = m.selected
	}

	if m.selected >= m.offset + visible {
		m.offset = m.selected - visible + 1
	}

	maxOffset := max(len(m.matches) - visible, 0)

	if m.offset > maxOffset {
		m.offset = maxOffset
	}

	if m.offset < 0 {
		m.offset = 0
	}
}

// filter runs the fuzzy matcher against the curreny query
func (m *Model) filter() {
	// if prefixSource is enabled check for prefix
	if m.prefixSource != nil {
		m.prefixSource.Query = m.query
		// fetch prefix items, halt on error
		prefixItems, err := m.prefixSource.List()
		if err != nil {
			m.matches = nil
			return
		}

		// if we have prefixItems then show those
		if prefixItems != nil {
			m.items = prefixItems
			// show all prefix items, no fuzzy matching
			m.matches = make([]fuzzy.Match, len(m.items))
			for i := range m.items {
				m.matches[i] = fuzzy.Match{
					Str: m.items[i].Name,
					Index: i,
				}
			}
			m.resetSelection()
			return
		}
		// otherwise regular fuzzy filtering
		m.items = m.defaultItems
	}

	// initial item source
	if m.query == "" {
		m.matches = make([]fuzzy.Match, len(m.items))

		for i := range m.items {
			m.matches[i] = fuzzy.Match{
				Str: m.items[i].Name,
				Index: i,
			}
		}
	} else {
		names := sources.GetNames(m.items)
		m.matches = fuzzy.Find(m.query, names)
	}
	m.resetSelection()
}
func (m *Model) resetSelection() {
	if m.selected >= len(m.matches) {
		m.selected = 0
	}
	m.offset = 0
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.keepSelectedVisible()

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			if m.flagLoop {
				m.reset()
				return m, nil
			}
			return m, tea.Quit

		case "up", "ctrl+k":
			if m.selected > 0 {
				m.selected--
				m.keepSelectedVisible()
			}

		case "down", "ctrl+j":
			if m.selected < len(m.matches) - 1 {
				m.selected++
				m.keepSelectedVisible()
			}

		case "backspace":
			runes := []rune(m.query)
			if len(runes) > 0 {
				m.query = string(runes[:len(runes)-1])
				m.filter()
			}

		case "enter":
			if len(m.matches) > 0 {
				selected := m.items[m.matches[m.selected].Index]

				// prefix item
				if selected.PrefixExecuteType != "" {
					err := executePrefix(selected)
					// TODO: handle errors properly
					if err != nil { m.Result = err.Error() }

					return m.handleFlagLoop()
				}
				// something matched: output the selected item
				// TODO: handle apps
				m.Result = m.matches[m.selected].Str
				return m.handleFlagLoop()
			}
			// no matches: user entered random text
			if !m.dmenuMode && m.query != "" {
				if err := executeCopyAndNotify("%s", m.query); err != nil {
					m.Result = err.Error()
				}
			}
			// dmenu mode still returns the raw query
			m.Result = m.query
			return m.handleFlagLoop()

		default:
			s := msg.String()

			if s == "space" {
				m.query += " "
				m.filter()
			} else if len(s) == 1 {
				m.query += s
				m.filter()
			}
		}
	}
	return m, nil
}

func (m Model) handleFlagLoop() (tea.Model, tea.Cmd) {
	if m.flagLoop {
		m.reset()
		return m, nil
	}
	return m, tea.Quit
}

func (m Model) View() tea.View {
	var out string

	out += m.center(m.query) + "\n\n"

	start := m.offset
	end := min(start + m.visibleCount(), len(m.matches))

	selectedStyle := lipgloss.NewStyle().
	  Width(m.terminalWidth).
		Background(lipgloss.Color("255")).
		Foreground(lipgloss.Color("0")).
		Align(lipgloss.Center)

	normalStyle := lipgloss.NewStyle().
	  Width(m.terminalWidth).
		Align(lipgloss.Center)

	for i := start; i < end; i++ {
		match := m.matches[i]
		style := normalStyle
		if i == m.selected {
			style = selectedStyle
		}
		out += style.Render(match.Str) + "\n"
	}
	v := tea.NewView(out)
	// fullscreen / alternate screen
	v.AltScreen = true
	return v
}

func (m Model) center(s string) string {
	if m.terminalWidth <= 0 {
		return s
	}

	return lipgloss.NewStyle().
    Width(m.terminalWidth).
		Align(lipgloss.Center).
		Render(s)
}
