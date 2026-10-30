package main

import "charm.land/lipgloss/v2"

var defaultStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
var focusedStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(1, 2)

type panel struct {
	name string
	body string
}

func (p panel) View(focused bool, width, height int) string {
	panelStr := lipgloss.NewStyle().Render(lipgloss.JoinVertical(lipgloss.Top, p.name, p.body))

	if focused {
		return focusedStyle.Width(width).Height(height).Render(panelStr)
	}
	return defaultStyle.Width(width).Height(height).Render(panelStr)
}
