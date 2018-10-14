package actn

import (
	"encoding/json"
	"er"
	"fwb"
	"fwb/cards"
	"hlf"
	"sgs"
)

const (

	//ACTN_SKIP skip one turn
	ACTN_SKIP = iota

	//ACTN_FARM farm action
	ACTN_FARM

	//ACTN_FEED_SHEEP feed sheep action
	ACTN_FEED_SHEEP

	//ACTN_TAKE_OFF take off action
	ACTN_TAKE_OFF

	//ACTN_PARTTIME_WORK parttime work action
	ACTN_PARTTIME_WORK

	//ACTN_UPGRADE_HOUSE upgrade house action
	ACTN_UPGRADE_HOUSE

	//ACTN_TRAIN train skills action
	ACTN_TRAIN

	//ACTN_EMPLOY employ action
	ACTN_EMPLOY

	//ACTN_TRADE trade action
	ACTN_TRADE

	//ACTN_HUNT hunt action
	ACTN_HUNT

	//ACTN_STEAL steal action
	ACTN_STEAL

	//ACTN_WEAVE weave wool action
	ACTN_WEAVE

	//ACTN_MAKE_WINE make wine action
	ACTN_MAKE_WINE

	//ACTN_PARTY party action
	ACTN_PARTY

	//ACTN_BEG beg action
	ACTN_BEG

	//ACTN_GOLD_MINING gold mining action
	ACTN_GOLD_MINING

	//ACTN_TAKE_VACATION take vacation action
	ACTN_TAKE_VACATION

	//ACTN_MAX max action index
	ACTN_MAX
)

var _actionNames = []string{
	"Skip",
	"Farm",
	"FeedSheep",
	"TakeOff",
	"PartTimeWork",
	"UpgradeHouse",
	"Train",
	"Employ",
	"Trade",
	"Hunt",
	"Steal",
	"Weave",
	"MakeWine",
	"Party",
	"Beg",
	"GoldMining",
	"TakeVacation",
}

type actnParser func(sgs.Command) fwb.Action

var _actnParser = map[int]actnParser{
	ACTN_SKIP:          actnSkipParser,
	ACTN_FARM:          actnBasicParser,
	ACTN_FEED_SHEEP:    actnBasicParser,
	ACTN_TAKE_OFF:      actnBasicParser,
	ACTN_PARTTIME_WORK: actnBasicParser,
	ACTN_UPGRADE_HOUSE: actnBasicParser,
	ACTN_TRAIN:         actnTrainParser,

	ACTN_EMPLOY:      actnBasicParser,
	ACTN_TRADE:       actnTradeParser,
	ACTN_HUNT:        actnBasicParser,
	ACTN_STEAL:       actnStealParser,
	ACTN_WEAVE:       actnWeaveParser,
	ACTN_MAKE_WINE:   actnBasicParser,
	ACTN_PARTY:       actnBasicParser,
	ACTN_BEG:         actnBasicParser,
	ACTN_GOLD_MINING: actnBasicParser,

	ACTN_TAKE_VACATION: actnBasicParser,
}

var _actnPersons = map[int]int{
	ACTN_SKIP:          1,
	ACTN_FARM:          1,
	ACTN_FEED_SHEEP:    1,
	ACTN_TAKE_OFF:      1,
	ACTN_PARTTIME_WORK: 1,
	ACTN_UPGRADE_HOUSE: 3,
	ACTN_TRAIN:         1,

	ACTN_EMPLOY:      1,
	ACTN_TRADE:       1,
	ACTN_HUNT:        1,
	ACTN_STEAL:       1,
	ACTN_WEAVE:       1,
	ACTN_MAKE_WINE:   1,
	ACTN_PARTY:       1,
	ACTN_BEG:         1,
	ACTN_GOLD_MINING: 1,

	ACTN_TAKE_VACATION: 1,
}

var _actnCards = map[int]int{
	ACTN_SKIP:          cards.CARD_VOID,
	ACTN_FARM:          cards.CARD_FARM,
	ACTN_FEED_SHEEP:    cards.CARD_FEED_SHEEP,
	ACTN_TAKE_OFF:      cards.CARD_TAKE_OFF,
	ACTN_PARTTIME_WORK: cards.CARD_PARTTIME_WORK,
	ACTN_UPGRADE_HOUSE: cards.CARD_UPGRADE_HOUSE,
	ACTN_TRAIN:         cards.CARD_TRAIN,

	ACTN_EMPLOY:      cards.CARD_EMPLOY,
	ACTN_TRADE:       cards.CARD_TRADE,
	ACTN_HUNT:        cards.CARD_HUNT,
	ACTN_STEAL:       cards.CARD_STEAL,
	ACTN_WEAVE:       cards.CARD_WEAVE,
	ACTN_MAKE_WINE:   cards.CARD_MAKE_WINE,
	ACTN_PARTY:       cards.CARD_PARTY,
	ACTN_BEG:         cards.CARD_BEG,
	ACTN_GOLD_MINING: cards.CARD_GOLD_MINING,

	ACTN_TAKE_VACATION: cards.CARD_TAKE_OFF,
}

//ActionParser object for parsing action command
type ActionParser struct {
	lg hlf.Logger
}

//MakeActionParser init ActionParser object
func MakeActionParser(game fwb.Game) *ActionParser {
	return &ActionParser{
		lg: game.GetLogger(),
	}
}

//Parse main service of ActionParser
func (me *ActionParser) Parse(command sgs.Command) (fwb.Action, *er.Err) {
	if command.ID != fwb.CMD_ACTION {
		return nil, er.Throw(fwb.E_INVALID_CMD, er.EInfo{
			"details":   "Command ID is invalid as an Action Command",
			"commandID": command.ID,
		}).To(me.lg)
	}

	payload, err := json.Marshal(command.Payload)

	if err != nil {
		return nil, er.Throw(fwb.E_INVALID_CMD, er.EInfo{
			"details": "Payload is invalid",
			"payload": command.Payload,
		}).To(me.lg)
	}

	var actnCmd struct {
		ActionID int
		Payload  interface{}
	}

	err = json.Unmarshal(payload, &actnCmd)

	if err != nil {
		return nil, er.Throw(fwb.E_INVALID_CMD, er.EInfo{
			"details": "Payload is invalid as an Action Command",
			"payload": command.Payload,
		}).To(me.lg)
	}

	parser, ok := _actnParser[actnCmd.ActionID]
	if !ok {
		return nil, er.Throw(fwb.E_ACTION_PARSER_NOT_FOUND, er.EInfo{
			"details":  "Parser not found in the Action Parser map for the Action Command",
			"actionID": actnCmd.ActionID,
		}).To(me.lg)
	}

	action := parser(command)

	if action == nil {
		return nil, er.Throw(fwb.E_INVALID_CMD, er.EInfo{
			"details":   "Failed to parse the Action Command",
			"commandID": command.ID,
			"payload":   command.Payload,
		}).To(me.lg)
	}

	return action, nil
}

func getTargetCard(gd *fwb.GameData, actionID int) (*fwb.Card, *er.Err) {
	minSlot, ok := _actnPersons[actionID]
	if !ok {
		return nil, er.Throw(fwb.E_INVALID_ACTION, er.EInfo{
			"details": "cannot find action information",
			"ID":      actionID,
		})
	}

	cardID, cardOk := _actnCards[actionID]
	if !cardOk {
		return nil, er.Throw(fwb.E_INVALID_ACTION, er.EInfo{
			"details": "cannot find action information",
			"ID":      actionID,
		})
	}

	card := FindCard(gd, cardID, minSlot)

	return card, nil
}

//FindCard find the card with the specified ID and empty slots no less than minSlot
func FindCard(me *fwb.GameData, cardID int, minSlot int) *fwb.Card {
	for i := range me.Cards {
		if me.Cards[i].ID == cardID && me.Cards[i].MaxSlot-len(me.Cards[i].Pawns) >= minSlot {
			return &me.Cards[i]
		}
	}
	return nil
}

func hasCardSlots(gd *fwb.GameData, actionID int) bool {
	card, _ := getTargetCard(gd, actionID)
	if card == nil {
		return false
	}
	return card.MaxSlot-len(card.Pawns) >= card.PawnPerTurn
}

func checkCard(gd *fwb.GameData, actionID int, playerID int, number int) *er.Err {
	card, err := getTargetCard(gd, actionID)
	if card != nil {
		for i := 0; i < number; i++ {
			card.Pawns = append(card.Pawns, playerID)
		}
	}
	return err
}
