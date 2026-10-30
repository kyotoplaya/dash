package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
			{name: "Panel 1"},
			{name: "Panel 2"},
			{name: "Panel 3"},
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
	colWidth := m.Width / len(m.panels)
	colHeight := m.Height / 4 // Set a fixed height for the panels

	views := make([]string, len(m.panels))

	for i, p := range m.panels {
		views[i] = p.View(i == m.focus, colWidth, colHeight)
	}

	panelsStr := lipgloss.JoinHorizontal(lipgloss.Top, views...)
	sizeStr := fmt.Sprintf("Window Size: %dx%d", m.Width, m.Height)
	footerStr := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press 'q' to quit. Press 'tab' to switch focus.")

	vStr := lipgloss.JoinVertical(lipgloss.Left, panelsStr, sizeStr, footerStr)

	v := tea.NewView(vStr)
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
