package actn

import (
	"encoding/json"
	"er"
	"fmt"
	"fwb"
	"sgs"
)

type actnTrain struct {
	skillID  int
	playerID int
}

func actnTrainParser(command sgs.Command) fwb.Action {
	payload, err := json.Marshal(command.Payload)

	if err != nil {
		return nil
	}

	var train struct {
		ActionID int
		Payload  int
	}

	err = json.Unmarshal(payload, &train)

	if err != nil {
		return nil
	}

	return &actnTrain{
		skillID:  train.Payload,
		playerID: command.Source,
	}
}

func (me *actnTrain) String() string {
	return fmt.Sprintf("[Action %v from Player %v, Skill ID %v]", _actionNames[ACTN_TRAIN], me.playerID, me.skillID)
}

func (me *actnTrain) ID() int {
	return ACTN_TRAIN
}

func (me *actnTrain) getCost(gd *fwb.GameData) fwb.PlayerData {
	cost := make(fwb.PlayerData, fwb.PD_MAX)

	i := gd.GetPDIndex(me.playerID)

	p := gd.PData[i]

	targetLv := p[me.skillID] + 1

	if targetLv == 1 {
		cost[fwb.PD_PT_GOLD] = -2
	} else if targetLv == 2 {
		cost[fwb.PD_PT_GOLD] = -5
	} else if targetLv == 3 {
		cost[fwb.PD_PT_GOLD] = -10
	}

	cost[fwb.PD_PAWNS] = -1
	return cost
}

func (me *actnTrain) ValidateAgainst(gd *fwb.GameData) bool {
	if me.skillID > fwb.PD_SK_INTELLIGENCE || me.skillID < fwb.PD_SK_STRENGTH {
		return false
	}

	i := gd.GetPDIndex(me.playerID)

	if i < 0 {
		return false
	}

	p := gd.PData[i]

	cost := me.getCost(gd)

	res := fwb.PDAdd(p, cost)

	targetLv := p[me.skillID] + 1

	return (targetLv < 1 && targetLv > 3) && res.AllAboveZero()
}

func (me *actnTrain) Do(gd *fwb.GameData) *er.Err {
	i := gd.GetPDIndex(me.playerID)
	if i < 0 {
		return er.Throw(fwb.E_DOACTION_INVALID_CLIENTID, er.EInfo{
			"details":  "invalid player ID for do action",
			"playerID": me.playerID,
		})
	}

	cost := me.getCost(gd)
	gain := make(fwb.PlayerData, fwb.PD_MAX)
	gain[me.skillID] = 1

	p := gd.PData[me.playerID]
	p = fwb.PDAdd(cost, p)
	p = fwb.PDAdd(gain, p)

	gd.PData[me.playerID] = p
	return checkCard(gd, ACTN_TRAIN, me.playerID, 1)
}
