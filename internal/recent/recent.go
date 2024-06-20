package recent

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	recents  []string
	cursor   int
	selected string
}

func initialModel() model {
	return model{
		recents: []string{"First", "Second", "Third"},
		cursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	s := "Select recent note\n\n"

	for i, recent := range m.recents {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, recent)
	}

	s += "\nPress esc to quit.\n"

	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.recents)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = m.recents[m.cursor]
			return m, tea.Quit
		}

	}
	return m, nil
}

func SelectRecent() (string, error) {
	p := tea.NewProgram(initialModel())
	tm, err := p.Run()
	if err != nil {
		return "", err
	}
	m := tm.(model)
	return m.selected, nil
}
