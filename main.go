package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

/*MAIN MODEL*/
type model struct {
	Width  int
	Height int

	panels          []panel
	maxPanelsPerRow int
	focus           int

	time time.Time
}

func initialModel() model {
	return model{
		panels: []panel{
			{name: "Clock", body: "--:--:--"},
			{name: "Random", body: "--"},
			{name: "Weather", body: "--"},
			{name: "Hardware", body: "--"},
		},
		maxPanelsPerRow: 3,
		focus:           0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchClock(),
		fetchRandom(),
		fetchWeather(),
		fetchHardware(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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

	case clock:
		m.time = time.Time(msg)
		cmd = tickClock() // Запускаем последующие тики каждую секунду
		m.panels[0].body = m.time.Format("15:04:05")

	case random:
		m.panels[1].body = fmt.Sprintf("%d", msg)
		cmd = tickRandom() // Запускаем последующие тики каждые 3 секунды

	case weather:
		m.panels[2].body = string(msg)

	case hardware:
		m.panels[3].body = string(msg)
		cmd = tickHardware() // Запускаем последующие тики каждые 5 секунд
	}

	return m, cmd
}

func (m model) View() tea.View {
	const panelOverhead = 2

	footerStr := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press 'q' to quit. Press 'tab' to switch focus.")

	colWidth := m.Width / m.maxPanelsPerRow

	views := make([]string, len(m.panels))

	rowsCount := (len(m.panels) + m.maxPanelsPerRow - 1) / m.maxPanelsPerRow

	footerHeight := strings.Count(footerStr, "\n")
	colHeight := (m.Height-footerHeight)/rowsCount - panelOverhead

	rem := m.Width % m.maxPanelsPerRow

	for i, p := range m.panels {
		w := colWidth
		if i%m.maxPanelsPerRow < rem {
			w++ // раздать остаток (свободные пиксели) первым панелям, чтобы они были равномерно распределены
		}

		views[i] = p.View(i == m.focus, w, colHeight)

	}

	var panels []string
	step := 0
	for i := m.maxPanelsPerRow; i <= len(views); i += m.maxPanelsPerRow {
		panels = append(panels, lipgloss.JoinHorizontal(lipgloss.Top, views[step:i]...))
		step += m.maxPanelsPerRow
	}

	panels = append(panels, lipgloss.JoinHorizontal(lipgloss.Top, views[step:]...))
	panelsStr := lipgloss.JoinVertical(lipgloss.Top, panels...)

	vStr := lipgloss.JoinVertical(lipgloss.Left, panelsStr, footerStr)

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
