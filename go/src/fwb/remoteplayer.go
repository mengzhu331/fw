package main

import "sgs"

type remotePlayer struct {
	client *sgs.Client
	df     playerDF
	fw     *FW
}

func (this *remotePlayer) getName() string {
	return this.client.Username
}

func (this *remotePlayer) getDF() *playerDF {
	return &this.df
}

func (this *remotePlayer) SendCommand(cmd sgs.Command) error {
	if cmd.Target != makeCommandParticipantUri(TARGET_PLAYER, this.getName()) {
		return MakeFwErrorByCode(EC_ILLEGAL_COMMAND_TARGET)
	}

	this.fw.SendCommand(sgs.Command{
		ID:     CMD_GAME_START_ACK,
		Source: makeCommandParticipantUri(TARGET_PLAYER, this.getName()),
		Target: makeCommandParticipantUri(TARGET_GAME, ""),
	})
	return nil
}
