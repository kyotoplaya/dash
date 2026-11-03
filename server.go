package main

import (
	"fmt"
	"net"
	"time"
)

type serverConfig struct {
	Ip      string `json:"ip"`
	Port    string `json:"port"`
	Timeout int    `json:"timeout"`
}

type serverModule struct {
	cfg serverConfig
}

func NewServerModule() Module {
	return &serverModule{
		cfg: serverConfig{Ip: "10.10.3.28", Port: "22", Timeout: 5},
	}
}

func (m *serverModule) Config() any {
	return &m.cfg
}

func (*serverModule) Id() string   { return "server" }
func (*serverModule) Name() string { return "Home Server" }
func (m *serverModule) Collect() (string, error) {
	host := m.cfg.Ip
	port := m.cfg.Port
	timeout := time.Duration(m.cfg.Timeout) * time.Second

	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	duration := time.Since(start)

	if err != nil {
		return fmt.Sprintf("Server is offline: %s", err), err
	}
	defer conn.Close()

	return fmt.Sprintf("Server is online\nPing: %s", duration), nil
}
func (*serverModule) Interval() time.Duration {
	return time.Minute
}
