package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

/*MAIN MODEL*/
type model struct {
	Width  int
	Height int

	panels []panel
	focus  int
}

func initialModel() model {
	return model{
		panels: []panel{
			{name: "panel1", index: 0},
			{name: "panel2", index: 1},
			{name: "panel3", index: 2},
		},
		focus: 0,
	}
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

		case "tab":
			focus := (m.focus + 1) % len(m.panels)
			m.focus = focus
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}

func (m model) View() tea.View {

	panelsStr := m.getPanels(m.panels)

	v := tea.NewView(panelsStr)
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
