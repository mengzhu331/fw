package main

import (
	"sgs"
	"strconv"
)

type playerDF struct{}

type player interface {
	getName() string
	getURI() string
	getDF() *playerDF
	exec(sgs.Command) error
}

type remotePlayer struct {
	client *sgs.Client
	df     playerDF
	game   *fwGame
}

func (this *remotePlayer) getName() string {
	return this.client.Username
}

func (this *remotePlayer) getDF() *playerDF {
	return &this.df
}

func (this *remotePlayer) exec(cmd sgs.Command) error {
	if cmd.Target != this.getURI() {
		return MakeFwErrorByCode(EC_ILLEGAL_COMMAND_TARGET)
	}

	this.game.sendCommand(sgs.Command{
		ID:     CMD_GAME_START_ACK,
		Source: this.getURI(),
		Target: this.game.getURI(),
	})
	//	this.game.sendCommand(cmd)
	return nil
}

func (this *remotePlayer) getURI() string {
	return "PLY://" + strconv.Itoa(int(this.client.ID))
}
