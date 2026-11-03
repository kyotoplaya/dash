package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type randomConfig struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type randomModule struct {
	cfg randomConfig
}

func NewRandomModule() Module {
	return &randomModule{
		cfg: randomConfig{Min: 0, Max: 100},
	}
}

func (m *randomModule) Config() any {
	return &m.cfg
}

func (*randomModule) Id() string   { return "random" }
func (*randomModule) Name() string { return "Random" }

func (m *randomModule) Collect() (string, error) {
	randomNumber := rand.IntN(m.cfg.Max - m.cfg.Min + 1)
	return fmt.Sprintf("Random Number (%d : %d): %d", m.cfg.Min, m.cfg.Max, randomNumber), nil
}
func (*randomModule) Interval() time.Duration {
	return 5 * time.Second
}
