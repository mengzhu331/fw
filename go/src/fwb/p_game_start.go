package main

import (
	"log"
	"sgs"
)

type p_game_start_state struct {
	ack int
}

func p_game_start_enter(this *fwGame, cmd sgs.Command) error {
	log.Println("P_GAME_START init")
	var err error

	this.df.cr = 1

	this.lsm[P_GAME_START] = &p_game_start_state{
		ack: 0,
	}

	for _, p := range this.fw.GetPlayers() {
		err = p.SendCommand(sgs.Command{
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
			s, _ := this.lsm[P_GAME_START].(*p_game_start_state)
			s.ack |= (1 << uint(i))
			if s.ack == (1<<uint(len(this.fw.GetPlayers())) - 1) {
				this.gotoPhase(P_ROUNDS_START)
				break
			}
		}
	}
	return nil
}
