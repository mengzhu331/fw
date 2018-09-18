package sgs

import (
	"time"
)

const (
	APP_MODE_1P int = 0x1
	APP_MODE_2P int = 0x2
	APP_MODE_3P int = 0x3
)

type CommandExecutor func(App, Command) error

type AppConfig struct {
	TickIntervalMs time.Duration
	CmdMap         map[uint]CommandExecutor
	Clients        []*Client
	SendCommand    func(Command)
	S              *Session
}

type App interface {
	Init(*AppConfig) error
	GetName() string
	GetModes() []int
}
