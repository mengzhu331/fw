package actn

import (
	"encoding/json"
	"er"
	"fmt"
	"fwb"
	"log"
	"sgs"
)

type actnTrade struct {
	playerID       int
	tradeDirection int
	amount         []int
}

//PtPrices get the price for each property
func PtPrices() fwb.PlayerData {
	p := make(fwb.PlayerData, fwb.PD_MAX)
	p[fwb.PD_PT_CEREALS] = 2
	p[fwb.PD_PT_MEAT] = 5
	p[fwb.PD_PT_SWEATER] = 10
	p[fwb.PD_PT_WOOL] = 5
	p[fwb.PD_PT_WINE] = 20
	return p
}

func actnTradeParser(command sgs.Command) fwb.Action {
	payload, err := json.Marshal(command.Payload)

	if err != nil {
		return nil
	}

	var trade struct {
		ActionID int
		Payload  struct {
			Direction int
			Amount    []int
		}
	}

	err = json.Unmarshal(payload, &trade)

	if err != nil {
		return nil
	}

	return &actnTrade{
		playerID:       command.Who,
		tradeDirection: trade.Payload.Direction,
		amount:         trade.Payload.Amount,
	}
}

func (me *actnTrade) String() string {
	var direction string
	if me.tradeDirection > 0 {
		direction = "BUY"
	} else if me.tradeDirection < 0 {
		direction = "SELL"
	} else {
		direction = "UNDEFINED"
	}
	return fmt.Sprintf("[Action %v from Player %v, Direction %v, Amount %v]", ActionNames[ACTN_TRADE], me.playerID, direction, me.amount)
}

func (me *actnTrade) ID() int {
	return ACTN_TRADE
}

func (me *actnTrade) getCost() fwb.PlayerData {
	cost := make(fwb.PlayerData, fwb.PD_MAX)
	p := PtPrices()

	if me.tradeDirection < 0 {
		for i, v := range me.amount {
			cost[i] = v * -1
		}
		cost[fwb.PD_PT_GOLD] -= 2
	} else {
		for i, v := range me.amount {
			cost[fwb.PD_PT_GOLD] -= p[i] * v
		}
		cost[fwb.PD_PT_GOLD] -= 2
	}
	cost[fwb.PD_PAWNS] = -1
	return cost
}

func (me *actnTrade) makeGain() fwb.PlayerData {

	gain := make(fwb.PlayerData, fwb.PD_MAX)

	p := PtPrices()

	if me.tradeDirection < 0 {
		for i, v := range me.amount {
			gain[fwb.PD_PT_GOLD] += p[i] * v
		}
	} else {
		gain = me.amount
	}
	return gain
}

func (me *actnTrade) ValidateAgainst(gd *fwb.GameData) bool {
	playeri := gd.GetPDIndex(me.playerID)
	if playeri < 0 {
		return false
	}

	pd := gd.PData[playeri]
	cost := me.getCost()
	log.Printf("trade cost %v", cost)
	log.Printf("trade est %v", fwb.PDAdd(cost, pd))
	return HasCardSlots(gd, ACTN_TRADE) && fwb.PDAdd(cost, pd).AllAboveZero()
}

func (me *actnTrade) Do(gd *fwb.GameData) *er.Err {
	i := gd.GetPDIndex(me.playerID)
	if i < 0 {
		return er.Throw(fwb.E_DOACTION_INVALID_CLIENTID, er.EInfo{
			"details":  "invalid player ID for do action",
			"playerID": me.playerID,
		})
	}

	cost := me.getCost()
	gain := me.makeGain()
	p := gd.PData[i]
	p = fwb.PDAdd(cost, p)
	p = fwb.PDAdd(gain, p)

	gd.PData[i] = p
	return checkCard(gd, ACTN_TRADE, me.playerID, 1)
}
