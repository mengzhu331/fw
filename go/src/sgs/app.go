package sgs

import (
	"time"
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
}
