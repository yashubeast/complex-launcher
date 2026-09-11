package internal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sahilm/fuzzy"
)

type Model struct {
	query          string
	selected       int
	items          []string
	matches        []fuzzy.Match

	Result         string
	Quit           bool

	offset         int
	maxVisible     int // 0 = auto-detect from terminal height
	terminalHeight int

	flagLoop       bool
}

func (m *Model) reset() {
	m.query = ""
	m.selected = 0
	m.offset = 0
	m.filter()
}

func NewModel(items []string, loop bool, maxVisible int) Model {
	m := Model{
		items:    items,
		selected: 0,
		flagLoop: loop,
		maxVisible: maxVisible,
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
	if m.query == "" {
		m.matches = make([]fuzzy.Match, len(m.items))

		for i := range m.items {
			m.matches[i] = fuzzy.Match{
				Str: m.items[i],
				Index: i,
			}
		}
	} else {
		m.matches = fuzzy.Find(m.query, m.items)
	}

	if m.selected >= len(m.matches) {
		m.selected = 0
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
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
				// something matched: output the selected item
				m.Result = m.matches[m.selected].Str
			} else {
				// nothing matched: output whatever the query was
				m.Result = m.query
			}
			if m.flagLoop {
				m.reset()
				return m, nil
			}
			return m, tea.Quit

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

func (m Model) View() tea.View {
	var out string

	out += "> " + m.query + "\n\n"

	start := m.offset
	end := start + m.visibleCount()

	if end > len(m.matches) {
		end = len(m.matches)
	}

	for i := start; i < end; i++ {
		match := m.matches[i]
		prefix := "  "
		if i == m.selected {
			prefix = "> "
		}
		out += prefix + match.Str + "\n"
	}
	v := tea.NewView(out)
	// fullscreen / alternate screen
	v.AltScreen = true
	return v
}
