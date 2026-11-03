package main

import "time"

type clockConfig struct{}

type clockModule struct {
	cfg clockConfig
}

func NewClockModule() Module {
	return &clockModule{
		cfg: clockConfig{},
	}
}

func (m *clockModule) Config() any {
	return &m.cfg
}

func (*clockModule) Id() string   { return "clock" }
func (*clockModule) Name() string { return "Clock" }
func (m *clockModule) Collect() (string, error) {
	return time.Now().Format("15:04:05"), nil
}

func (*clockModule) Interval() time.Duration {
	now := time.Now()

	// Находим точное время начала СЛЕДУЮЩЕЙ системной секунды
	nextSecond := now.Truncate(time.Second).Add(time.Second)

	delay := time.Until(nextSecond)

	return delay
}
