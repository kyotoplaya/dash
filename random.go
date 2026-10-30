package main

import (
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
)

type random int

func tickRandom() tea.Cmd {
	rnd := rand.IntN(100)

	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return random(rnd)
	})
}
