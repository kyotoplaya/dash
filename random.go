package main

import (
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
)

type random int

func collectRandom() random {
	return random(rand.IntN(100))
}

func fetchRandom() tea.Cmd {
	return func() tea.Msg {
		return collectRandom()
	}
}

func tickRandom() tea.Cmd {

	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return collectRandom()
	})
}
