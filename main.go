package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puggeraponanna/mycal/calendar"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Some error: ", err)
		os.Exit(1)
	}
}

type model struct {
	cal *calendar.Calendar
}

func initialModel() model {
	return model{
		cal: calendar.New(),
	}
}

func (m model) View() string {
	return m.cal.Render()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cal.SelectR.Row = (m.cal.SelectR.Row + m.cal.NRows - 1) % m.cal.NRows
		case "down", "j":
			m.cal.SelectR.Row = (m.cal.SelectR.Row + 1) % m.cal.NRows
		case "left", "h":
			m.cal.SelectR.Col = (m.cal.SelectR.Col + 6) % 7
		case "right", "l":
			m.cal.SelectR.Col = (m.cal.SelectR.Col + 1) % 7
        case "enter", " ":
            fmt.Println("Selected date is: ", m.cal.Date())
            return m, tea.Quit
        }
	}
    return m, nil
}
