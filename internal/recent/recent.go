package recent

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/indimeco/jyggalag/internal/state"
)

type model struct {
	recents  []string
	cursor   int
	selected string
}

func initialModel() model {
	recents, err := state.ReadRecent()
	if err != nil {
		log.Println(fmt.Sprintf("Warning: failed to read recents: %v", err))
	}
	return model{
		recents: recents,
		cursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

var cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("82")).Blink(true)
var selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("82"))
var footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

func (m model) View() string {
	s := "Select recent note\n\n"

	for i, recent := range m.recents {

		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
			s += fmt.Sprintf("%s %s\n", cursor, selectedStyle.Render(recent))
		} else {
			s += fmt.Sprintf("%s %s\n", cursor, recent)
		}
	}

	s += footerStyle.Render("\nPress esc to quit.\n")

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
