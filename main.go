package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func slicePanels(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

/*MAIN MODEL*/
type model struct {
	Width  int
	Height int

	modules         []Module
	panels          []panel
	maxPanelsPerRow int
	focus           int

	time time.Time
}

func loadModules(modules []Module) error {

	err := isDirExists("modules")
	if err != nil {
		return err
	}

	for _, m := range modules {
		path := filepath.Join("modules", m.Id()+".json")
		fileData, err := os.ReadFile(path)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				cfg := m.Config()

				data, err := json.MarshalIndent(cfg, "", " ")
				if err != nil {
					return err
				}

				err = os.WriteFile(path, data, 0644)
				if err != nil {
					return err
				}
				continue
			}
			return err

		}

		err = json.Unmarshal(fileData, m.Config())
		if err != nil {
			return err
		}

	}

	return nil

}

func initialModel(modules []Module) model {

	panels := []panel{}
	for _, module := range modules {
		p := panel{name: module.Name(), body: "Loading..."}
		panels = append(panels, p)
	}

	return model{
		modules:         modules,
		panels:          panels,
		maxPanelsPerRow: 2,
		focus:           0,
	}
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	for _, module := range m.modules {
		cmds = append(cmds, poll(module, 0)) // Запускаем первый тик сразу
	}

	return tea.Batch(
		cmds...,
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

	case moduleMsg:
		for i, p := range m.panels {
			if p.name == msg.name {
				if msg.err != nil {
					m.panels[i].body = msg.err.Error()
				} else {
					m.panels[i].body = msg.text
				}
				break
			}
		}

		for i, module := range m.modules {
			if module.Name() == msg.name {
				cmd = poll(m.modules[i], m.modules[i].Interval())
			}
		}

	}

	return m, cmd
}

func (m model) View() tea.View {
	const panelOverhead = 2

	footerStr := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press 'q' to quit. Press 'tab' to switch focus.")

	colWidth := m.Width / m.maxPanelsPerRow

	views := make([]string, len(m.panels))

	rowsCount := (len(m.panels) + m.maxPanelsPerRow - 1) / m.maxPanelsPerRow

	footerHeight := lipgloss.Height(footerStr)
	colHeight := (m.Height-footerHeight)/rowsCount - panelOverhead

	rem := m.Width % m.maxPanelsPerRow

	for i, p := range m.panels {
		w := colWidth
		if i%m.maxPanelsPerRow < rem {
			w++ // раздать остаток (свободные пиксели) первым панелям, чтобы они были равномерно распределены
		}

		views[i] = p.View(i == m.focus, w, colHeight)

	}

	sliced := slicePanels(views, m.maxPanelsPerRow)
	panels := []string{}

	for _, slice := range sliced {
		panels = append(panels, lipgloss.JoinHorizontal(lipgloss.Top, slice...))
	}

	panelsStr := lipgloss.JoinVertical(lipgloss.Top, panels...)

	vStr := lipgloss.JoinVertical(lipgloss.Left, panelsStr, footerStr)

	v := tea.NewView(vStr)
	v.AltScreen = true
	return v
}

func main() {
	modules := []Module{
		NewClockModule(),
		NewWeatherModule(),
		NewDembelModule(),
		NewRandomModule(),
		NewServerModule(),
		NewHardwareModule(),
	}

	if err := loadModules(modules); err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(initialModel(modules))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
