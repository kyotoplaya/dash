package main

import (
	"fmt"
	"time"
)

type dembelConfig struct {
	DembelDate  string `json:"dembelDate"`
	SoldierName string `json:"soldierName"`
}

type dembelModule struct {
	cfg dembelConfig
}

func NewDembelModule() Module {
	return &dembelModule{
		cfg: dembelConfig{DembelDate: "23.11.2026", SoldierName: "Ваня"},
	}
}

func (m *dembelModule) Config() any {
	return &m.cfg
}

func (*dembelModule) Id() string   { return "dembel" }
func (*dembelModule) Name() string { return "Dembel Timer" }

func (m dembelModule) Collect() (string, error) {
	demobilizationDate := m.cfg.DembelDate
	demobilizationTime := "00:00:00"
	demobilizationDateTime := demobilizationDate + " " + demobilizationTime

	layout := "02.01.2006 15:04:05"
	demobilization, err := time.ParseInLocation(layout, demobilizationDateTime, time.Local)
	if err != nil {
		return "", err
	}

	remaining := time.Until(demobilization)
	if remaining < 0 {
		remaining = 0
	}

	totalSeconds := int64(remaining / time.Second)
	days := totalSeconds / 86400
	hours := totalSeconds % 86400 / 3600
	minutes := totalSeconds % 3600 / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%s will be demobilized in %02d:%02d:%02d:%02d", m.cfg.SoldierName, days, hours, minutes, seconds), nil
}

func (dembelModule) Interval() time.Duration {
	return time.Second
}
