package main

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type clock time.Time

func collectClock() clock {
	return clock(time.Now())
}

func fetchClock() tea.Cmd {
	return func() tea.Msg {
		return collectClock()
	}
}

func tickClock() tea.Cmd {
	now := time.Now()

	// Находим точное время начала СЛЕДУЮЩЕЙ системной секунды
	nextSecond := now.Truncate(time.Second).Add(time.Second)

	// time.Until посчитает, сколько миллисекунд осталось до этого момента
	delay := time.Until(nextSecond)

	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return clock(t)
	})
}
