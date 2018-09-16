package main

import (
	"log"
	"sgs"
)

func p_rounds_turns_tick(this *fwGame, cmd sgs.Command) error {
	log.Println("P_ROUNDS_TURNS tick")
	return nil
}

func p_rounds_turns_claim_action_card(this *fwGame, cmd sgs.Command) error {
	protocol, name := extractCommandParticipant(cmd.Source)

	pl := this.fw.GetPlayer(name)

	if protocol != TARGET_PLAYER || pl == nil {
		return MakeFwErrorByCode(EC_UNKNOWN_PLAYER_NAME)
	}

	acType := cmd.Param.(int)

	for _, ac := range this.sacs {
		if ac.acType == acType {
			return this.claimActionCard(pl, ac)
		}
	}
	return nil
}

func (this *fwGame) claimActionCard(p player, ac actionCard) error {
	err := ac.claimedBy(p)
	return err
}
