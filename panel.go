package main

import "charm.land/lipgloss/v2"

var defaultStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
var focusedStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(1, 2)

type panel struct {
	name  string
	index int
}

func (m model) getPanels(panels []panel) string {

	panelViews := make([]string, len(m.panels))
	for i, p := range m.panels {
		style := defaultStyle
		if i == m.focus {
			style = focusedStyle
		}
		panelViews[i] = style.Render(p.name)
	}

	panelsStr := lipgloss.JoinHorizontal(lipgloss.Top, panelViews...)

	return panelsStr
}
