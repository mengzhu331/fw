package actn

import (
	"encoding/json"
	"er"
	"fmt"
	"fwb"
	"sgs"
)

type actnSteal struct {
	targetPlayer int
	playerID     int
}

func actnStealParser(command sgs.Command) fwb.Action {
	payload, err := json.Marshal(command.Payload)

	if err != nil {
		return nil
	}

	var steal struct {
		ActionID int
		Payload  int
	}

	err = json.Unmarshal(payload, &steal)

	if err != nil {
		return nil
	}

	return &actnSteal{
		playerID:     command.Who,
		targetPlayer: steal.Payload,
	}
}

func (me *actnSteal) String() string {
	return fmt.Sprintf("[Action %v from Player %v, Target Player %v]", ActionNames[ACTN_STEAL], me.playerID, me.targetPlayer)
}

func (me *actnSteal) ID() int {
	return ACTN_STEAL
}

func (me *actnSteal) ValidateAgainst(gd *fwb.GameData) bool {
	playeri := gd.GetPDIndex(me.playerID)
	if playeri < 0 {
		return false
	}

	playerInte := gd.PData[playeri][fwb.PD_SK_INTELLIGENCE]
	targetKnow := gd.PData[me.targetPlayer][fwb.PD_SK_KNOWLEDGE]
	playerPawns := gd.PData[playeri][fwb.PD_PAWNS]
	return HasCardSlots(gd, ACTN_STEAL) && playerInte > targetKnow && playerPawns > 0
}

func (me *actnSteal) Do(gd *fwb.GameData) *er.Err {
	i := gd.GetPDIndex(me.playerID)
	if i < 0 {
		return er.Throw(fwb.E_DOACTION_INVALID_CLIENTID, er.EInfo{
			"details":  "invalid player ID for do action",
			"playerID": me.playerID,
		})
	}

	playerInte := gd.PData[i][fwb.PD_SK_INTELLIGENCE]
	targetKnow := gd.PData[me.targetPlayer][fwb.PD_SK_KNOWLEDGE]
	diff := playerInte - targetKnow
	amount := 0

	if diff > 2 {
		amount = 30
	} else if diff > 1 {
		amount = 15
	} else if diff > 0 {
		amount = 10
	}

	gd.PData[i][fwb.PD_PT_GOLD] += amount
	gd.PData[i][fwb.PD_PAWNS] -= 1
	gd.PData[me.targetPlayer][fwb.PD_PT_GOLD] -= amount
	return checkCard(gd, ACTN_STEAL, me.playerID, 1)
}
