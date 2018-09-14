package sgs

import (
	"log"
	"time"
)

type Session struct {
	ID           int
	maxClient    int
	clients      []*Client
	appConfig    AppConfig
	app          App
	lastTickTime time.Time
	cmdChOut     chan Command
}

func (this *Session) run() {
	log.Println("[Session] ID ", this.ID, " started")
	this.appConfig.Clients = this.clients
	this.appConfig.S = this
	if err := this.app.Init(&this.appConfig); err != nil {
		log.Panicln("[Session] Application init failed: " + err.Error())
	}

	for {
		select {
		case <-time.After(time.Duration(this.appConfig.TickIntervalMs) * time.Millisecond):
			this.doTick()
		}
	}
}

func (this *Session) doTick() {
	tt := time.Now()
	dt := time.Duration(0)
	if !this.lastTickTime.IsZero() {
		dt = tt.Sub(this.lastTickTime)
	}

	tick := Command{
		ID:     CMD_TICK,
		Param:  int(dt / time.Millisecond),
		Source: "SYS://framework",
	}

	this.sendCommandToApplication(tick)
	this.lastTickTime = tt
}

func (this *Session) sendCommandToApplication(command Command) {
	var err error = nil
	if this.appConfig.CmdMap[command.ID] != nil {
		exec := this.appConfig.CmdMap[command.ID]
		err = exec(this.app, command)
	}

	if err == nil && this.appConfig.CmdMap[CMD_ANY] != nil {
		exec := this.appConfig.CmdMap[CMD_ANY]
		err = exec(this.app, command)
	}

	this.handleError(err)
}

func (this *Session) Exec(command Command) error {
	if command.Target == "APP://"+this.app.GetName() {
		this.sendCommandToApplication(command)
	}
	return nil
}

func (this *Session) handleError(err error) {
}
