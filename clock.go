package main

import "time"

type clockModule struct{}

func (clockModule) Name() string { return "Clock" }
func (clockModule) Collect() (string, error) {
	return time.Now().Format("15:04:05"), nil
}
func (clockModule) Interval() time.Duration {
	now := time.Now()

	// Находим точное время начала СЛЕДУЮЩЕЙ системной секунды
	nextSecond := now.Truncate(time.Second).Add(time.Second)

	delay := time.Until(nextSecond)

	return time.Second - delay // Возвращаем оставшееся время до следующей секунды
}
