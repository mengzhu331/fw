package actn

import (
	"encoding/json"
	"er"
	"fmt"
	"fwb"
	"sgs"
)

var _personsRequired = map[int]int{}

type actnBasic struct {
	playerID int
	actionID int
}

func (me *actnBasic) getCost(player fwb.PlayerData) fwb.PlayerData {
	cost := make(fwb.PlayerData, fwb.PD_MAX)

	cost[fwb.PD_PAWNS] = -1

	switch me.actionID {
	case ACTN_FEED_SHEEP:
		cost[fwb.PD_PT_CEREALS] = -2
	case ACTN_EMPLOY:
		cost[fwb.PD_PT_GOLD] = -10
	case ACTN_PARTY:
		cost[fwb.PD_PT_BEER] = -2
		cost[fwb.PD_PT_MEAT] = -2
	case ACTN_TAKE_VACATION:
		cost[fwb.PD_PT_GOLD] = -5
	case ACTN_UPGRADE_HOUSE:
		cost[fwb.PD_PAWNS] = -3
		cost[fwb.PD_PT_GOLD] = -10
	case ACTN_MAKE_WINE:
		cost[fwb.PD_PT_CEREALS] = -8
	}

	return cost
}

func (me *actnBasic) checkRequirement(player fwb.PlayerData) bool {
	cost := me.getCost(player)

	res := fwb.PDAdd(player, cost)

	if !res.AllAboveZero() {
		return false
	}

	fulfill := true
	inte := player[fwb.PD_SK_INTELLIGENCE]
	know := player[fwb.PD_SK_KNOWLEDGE]
	stre := player[fwb.PD_SK_STRENGTH]

	switch me.actionID {
	case ACTN_UPGRADE_HOUSE:
		if player[fwb.PD_HOUSE_LV] == 1 {
			fulfill = inte >= 1 && know >= 1 && stre >= 1
		} else if player[fwb.PD_HOUSE_LV] == 2 {
			fulfill = inte >= 2 && know >= 2 && stre >= 2
		}
	case ACTN_GOLD_MINING:
		fulfill = know >= 1 && stre >= 1
	}

	if !fulfill {
		return false
	}

	return fulfill
}

func (me *actnBasic) makeGain(player fwb.PlayerData) fwb.PlayerData {
	gain := make(fwb.PlayerData, fwb.PD_MAX)

	inte := player[fwb.PD_SK_INTELLIGENCE]
	know := player[fwb.PD_SK_KNOWLEDGE]
	stre := player[fwb.PD_SK_STRENGTH]

	switch me.ID() {
	case ACTN_FARM:
		if stre < 1 {
			gain[fwb.PD_PT_CEREALS] = 2
		} else if stre == 1 {
			gain[fwb.PD_PT_CEREALS] = 4
		} else if stre == 2 {
			gain[fwb.PD_PT_CEREALS] = 6
		} else if stre > 2 {
			gain[fwb.PD_PT_CEREALS] = 10
		}

	case ACTN_TAKE_OFF:
		gain[fwb.PD_PT_HEART] = 1

	case ACTN_PARTTIME_WORK:
		if inte < 1 {
			gain[fwb.PD_PT_GOLD] = 2
		} else if inte == 1 {
			gain[fwb.PD_PT_GOLD] = 4
		} else if inte == 2 {
			gain[fwb.PD_PT_GOLD] = 6
		} else if inte > 2 {
			gain[fwb.PD_PT_GOLD] = 10
		}

	case ACTN_EMPLOY:
		gain[fwb.PD_PAWNS] = 3

	case ACTN_HUNT:
		if stre > 2 {
			gain[fwb.PD_PT_MEAT] = 4
		} else if stre > 0 {
			gain[fwb.PD_PT_MEAT] = 2
		} else {
			gain[fwb.PD_PT_MEAT] = 1
		}

	case ACTN_BEG:
		if inte > 2 {
			gain[fwb.PD_PT_MEAT] = 3
			gain[fwb.PD_PT_GOLD] = 5
		} else {
			gain[fwb.PD_PT_CEREALS] = 2
		}

	case ACTN_FEED_SHEEP:
		if know < 1 {
			gain[fwb.PD_PT_MEAT] = 1
			gain[fwb.PD_PT_WOOL] = 1
		} else if know == 1 {
			gain[fwb.PD_PT_MEAT] = 2
			gain[fwb.PD_PT_WOOL] = 1
		} else if know == 2 {
			gain[fwb.PD_PT_MEAT] = 2
			gain[fwb.PD_PT_WOOL] = 2
		} else if know > 2 {
			gain[fwb.PD_PT_MEAT] = 3
			gain[fwb.PD_PT_WOOL] = 3
		}

	case ACTN_TAKE_VACATION:
		gain[fwb.PD_PT_HEART] = 2

	case ACTN_PARTY:
		gain[fwb.PD_PT_HEART] = 10

	case ACTN_UPGRADE_HOUSE:
		gain[fwb.PD_HOUSE_LV] = 1

	case ACTN_GOLD_MINING:
		if stre > 2 && know > 2 {
			gain[fwb.PD_PT_GOLD] = 20
		} else if stre > 1 && know > 1 {
			gain[fwb.PD_PT_GOLD] = 10
		} else {
			gain[fwb.PD_PT_GOLD] = 5
		}
	case ACTN_MAKE_WINE:
		if know > 1 && inte > 1 {
			gain[fwb.PD_PT_BEER] = 2
		} else {
			gain[fwb.PD_PT_BEER] = 1
		}
	}
	return gain
}

func actnBasicParser(command sgs.Command) fwb.Action {
	payload, err := json.Marshal(command.Payload)

	if err != nil {
		return nil
	}

	var actnCmd struct {
		ActionID int
		Payload  interface{}
	}

	err = json.Unmarshal(payload, &actnCmd)

	if err != nil {
		return nil
	}

	return &actnBasic{
		playerID: command.Who,
		actionID: actnCmd.ActionID,
	}
}

func (me *actnBasic) String() string {
	return fmt.Sprintf("[Action: %v from Player: %v]", ActionNames[me.actionID], me.playerID)
}

func (me *actnBasic) ID() int {
	return me.actionID
}

func (me *actnBasic) Do(gd *fwb.GameData) *er.Err {
	i := gd.GetPDIndex(me.playerID)
	if i < 0 {
		return er.Throw(fwb.E_DOACTION_INVALID_CLIENTID, er.EInfo{
			"details":  "invalid player ID for do action",
			"playerID": me.playerID,
		})
	}

	p := gd.PData[i]

	cost := me.getCost(p)
	gain := me.makeGain(p)

	res := fwb.PDAdd(p, cost)
	res = fwb.PDAdd(res, gain)
	gd.PData[i] = res

	return checkCard(gd, me.actionID, me.playerID, -cost[fwb.PD_PAWNS])
}

func (me *actnBasic) ValidateAgainst(gd *fwb.GameData) bool {
	if !HasCardSlots(gd, me.actionID) {
		return false
	}

	playeri := gd.GetPDIndex(me.playerID)

	if playeri < 0 {
		return false
	}

	return me.checkRequirement(gd.PData[playeri])
}
