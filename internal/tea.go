package internal

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sahilm/fuzzy"
)

type Model struct {
	query    string
	selected int
	items    []string
	matches  []fuzzy.Match
	Result   string
	Quit     bool
	flagLoop bool
}

func (m *Model) reset() {
	m.query = ""
	m.selected = 0
	m.filter()
}

func NewModel(items []string, loop bool) Model {
	m := Model{
		items:    items,
		selected: 0,
		flagLoop: loop,
	}
	m.filter()
	return m
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
			}

		case "down", "ctrl+j":
			if m.selected < len(m.matches) - 1 {
				m.selected++
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

	for i, match := range m.matches {
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
