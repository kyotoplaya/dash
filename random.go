package main

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type randomModule struct{}

func (randomModule) Name() string { return "Random" }
func (randomModule) Collect() (string, error) {
	randomNumber := rand.IntN(100)
	return "Random Number: " + strconv.Itoa(randomNumber), nil
}
func (randomModule) Interval() time.Duration {
	return 5 * time.Second
}
