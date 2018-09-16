package main

import "sgs"

type remotePlayer struct {
	client *sgs.Client
	df     playerDF
	fw     *FW
	si     []settlementItem
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
	return nil
}

func (this *remotePlayer) setDF(df *playerDF) {
	this.df = *df
}

func (this *remotePlayer) getSI() []settlementItem {
	return this.si
}

func (this *remotePlayer) removeSI(i int) {
	this.si = append(this.si[:i], this.si[i+1:]...)
}
