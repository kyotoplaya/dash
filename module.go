package main

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type Module interface {
	Name() string
	Collect() (string, error)
	Interval() time.Duration
}

// moduleMsg — единственное сообщение для всех модулей.
type moduleMsg struct {
	name string
	text string
	err  error
}

func poll(m Module, delay time.Duration) tea.Cmd {
	return tea.Tick(delay, func(time.Time) tea.Msg {
		msg, err := m.Collect()
		return moduleMsg{
			name: m.Name(),
			text: msg,
			err:  err,
		}
	})
}
