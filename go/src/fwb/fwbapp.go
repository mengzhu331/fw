package main

import (
	"log"
	"sgs"
)

type FW struct {
	conf    *sgs.AppConfig
	game    *fwGame
	players []player
	pnm     map[string]player
}

type IFW interface {
	SendCommand(sgs.Command) error
	GetPlayers() []player
	GetGame() *fwGame
}

func FwCreator() sgs.App {
	return &FW{}
}

func (this *FW) GetPlayers() []player {
	return this.players
}

func (this *FW) GetGame() *fwGame {
	return this.game
}

func (this *FW) GetName() string {
	return "Family War"
}

func (this *FW) GetModes() []int {
	return []int{
		sgs.APP_MODE_2P,
		sgs.APP_MODE_3P,
	}
}

func (this *FW) Init(conf *sgs.AppConfig) error {
	log.Println("Init FW")
	this.conf = conf
	this.conf.TickIntervalMs = 100
	this.conf.CmdMap[sgs.CMD_ANY] = commandDispatcher

	this.game = &fwGame{}

	this.pnm = make(map[string]player)
	for _, c := range this.conf.Clients {
		this.players = append(this.players, &remotePlayer{
			client: c,
			df:     playerDF{},
			fw:     this,
		})
		this.pnm[this.players[len(this.players)-1].getName()] = this.players[len(this.players)-1]
	}

	return this.game.Init(this)
}

func (this *FW) SendCommand(command sgs.Command) error {
	return commandDispatcher(this, command)
}

func commandDispatcher(app sgs.App, command sgs.Command) error {
	this, ok := app.(*FW)
	if !ok {
		return MakeFwErrorByCode(EC_INVALID_APP_THIS)
	}

	target, name := extractCommandParticipant(command.Target)

	switch target {
	case TARGET_GAME:
		return this.game.SendCommand(command)
	case TARGET_PLAYER:
		pl, ok := this.pnm[name]
		if ok {
			return pl.sendCommand(command)
		}
	case sgs.TARGET_SYS_APP:
		return this.process(command)
	}
	return nil
}

func (this *FW) process(command sgs.Command) error {
	var err error
	if command.ID == sgs.CMD_TICK {
		err = this.game.SendCommand(command)
		if err != nil {
			return err
		}
		for _, p := range this.players {
			err = p.sendCommand(command)
		}
	}
	return err
}
