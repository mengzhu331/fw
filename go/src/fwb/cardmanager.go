package main

import (
	"encoding/json"
	"er"
	"io/ioutil"
	"strconv"
)

type cardManager struct {
	gm *gameImp

	BasicCardSet []int

	OptionalCardSet map[int]int

	cards []card
}

func (me *cardManager) loadCards() *er.Err {
	profile := me.gm.profile
	e := loadCards(me.gm, &me.cards)
	if e != nil {
		return e
	}

	path := "./profiles/" + profile + "/" + "cardset.conf"
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return er.Throw(_E_MISSING_GAME_SETTINGS, er.EInfo{
			"details":  "failed to load cardset settings",
			"profile":  profile,
			"io error": err.Error(),
		}).To(me.gm.lg)
	}

	err = json.Unmarshal(c, me)
	if err != nil {
		return er.Throw(_E_INVALID_SETTINGS, er.EInfo{
			"details":      "failed to decode cardset settings",
			"decode error": err.Error(),
		}).To(me.gm.lg)
	}

	me.gm.lg.Inf("Basic Card Set: %v", basicCardSetToString(&me.BasicCardSet))
	me.gm.lg.Inf("Optional Card Set: %v", optionalCardSetToString(&me.OptionalCardSet))
	return nil

}

func (me *cardManager) makeCardSet() []card {
	return nil
}

func basicCardSetToString(me *[]int) string {
	s := "["
	for i := 0; i < len(*me); i++ {
		s += strconv.Itoa((*me)[i])

		if i != len(*me)-1 {
			s += ", "
		}
	}
	s += "]"
	return s
}

func optionalCardSetToString(me *map[int]int) string {
	s := "["
	l := len(*me)
	for k, v := range *me {
		s += strconv.Itoa(k) + ":" + strconv.Itoa(v)
		if l != 1 {
			s += ", "
		}
		l--
	}
	s += "]"
	return s
}
