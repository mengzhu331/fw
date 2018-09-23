package main

import (
	"encoding/json"
	"er"
	"io/ioutil"
)

const (
	_CARD_FARM = iota + 1
	_CARD_FEED_SHEEP
	_CARD_TAKE_VACATION
	_CARD_PARTTIME_WORK
	_CARD_EMPLOY
	_CARD_TRADE

	_CARD_HUNT
	_CARD_FISH
	_CARD_MAJONG
	_CARD_HEAVY_WORK

	_CARD_STEAL
	_CARD_ARSON

	_CARD_UPGRADE_HOUSE
	_CARD_TRAIN_BRAWN
	_CARD_TRAIN_KNOWLEDGE
	_CARD_TRAIN_STEALTH
	_CARD_TRAIN_KEENNESS

	_CARD_GOLD_MINING
	_CARD_TRAP

	_CARD_MAX
)

type card struct {
	ID             int
	MaxSlot        int
	InUseSlot      int
	WorkerRequired int
	PersonRequired int
	Prob           int
	ActionID       int
}

func loadCards(gm *gameImp, cards *[]card, profile string) *er.Err {
	path := "./profiles/" + profile + "/" + "cards.conf"
	c, err := ioutil.ReadFile(path)
	if err != nil {
		er.Throw(_E_MISSING_GAME_SETTINGS, er.EInfo{
			"details":  "failed to load card settings",
			"profile":  profile,
			"io error": err.Error(),
		}).To(gm.lg)
	}

	err = json.Unmarshal(c, cards)
	if err != nil {
		er.Throw(_E_INVALID_SETTINGS, er.EInfo{
			"details":      "failed to decode card settings",
			"decode error": err.Error(),
		}).To(gm.lg)
	}

	return nil
}
