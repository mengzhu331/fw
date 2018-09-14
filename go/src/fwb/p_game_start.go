package main

import (
	"log"
	"sgs"
)

func p_game_start_init(this *fwGame, cmd sgs.Command) error {
	log.Println("P_GAME_START init")
	var err error
	this.lsm[P_GAME_START] = 0
	for _, p := range this.fw.GetPlayers() {
		err = p.sendCommand(sgs.Command{
			ID:     CMD_GAME_START,
			Source: makeCommandParticipantUri(TARGET_GAME, ""),
			Target: makeCommandParticipantUri(TARGET_PLAYER, p.getName()),
		})
		if err != nil {
			break
		}
	}
	return err
}

func p_game_start_game_start_ack(this *fwGame, cmd sgs.Command) error {
	log.Println("Player ACK: ", cmd.Source)

	for i, v := range this.fw.GetPlayers() {
		if cmd.Source == makeCommandParticipantUri(TARGET_PLAYER, v.getName()) {
			playerAck, _ := this.lsm[P_GAME_START].(int)
			playerAck |= (1 << uint(i))
			this.lsm[P_GAME_START] = playerAck
			if playerAck == (1<<uint(len(this.fw.GetPlayers())) - 1) {
				this.gotoPhase(P_GAME_ROUNDS)
				break
			}
		}
	}
	return nil
}
