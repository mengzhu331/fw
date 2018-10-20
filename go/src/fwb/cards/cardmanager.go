package cards

import (
	"er"
	"fwb"
	"strconv"
	"sutil"
)

type cardManager struct {
	BasicCardSet []int

	OptionalCardSet []int

	SpecialCardSet []int

	shuffledCardSet []int

	OptionalCardsPerRound int

	cards map[int]fwb.Card
}

//MakeCardManager instantiate a cardManage object
func MakeCardManager() fwb.CardManager {
	return &cardManager{}
}

//LoadCards implement CardManager.LoadCards() to load cards for a profile
func (me *cardManager) LoadCards(profile string) *er.Err {
	path := "./conf/profiles/" + profile + "/" + "cards.conf"
	cardList := make([]fwb.Card, 0)
	e := sutil.LoadConfFile(path, &cardList)
	if e != nil {
		return er.Throw(fwb.E_MISSING_GAME_SETTINGS, er.EInfo{
			"details":  "failed to load cards",
			"path":     path,
			"io error": e.Error(),
		})
	}

	me.cards = make(map[int]fwb.Card)
	for _, c := range cardList {
		me.cards[c.ID] = c
	}

	path = "./conf/profiles/" + profile + "/" + "cardset.conf"
	e = sutil.LoadConfFile(path, &me)
	if e != nil {
		return er.Throw(fwb.E_MISSING_GAME_SETTINGS, er.EInfo{
			"details":  "failed to load card configuration",
			"path":     path,
			"io error": e.Error(),
		})
	}
	return nil

}

func (me *cardManager) shuffle() {
	me.shuffledCardSet = make([]int, len(me.OptionalCardSet))
	copy(me.shuffledCardSet, me.OptionalCardSet)
	me.shuffledCardSet = sutil.ShuffleInt(me.shuffledCardSet...)
}

func (me *cardManager) MakeCardSet() ([]fwb.Card, []fwb.Card, []fwb.Card) {
	specialCards := make([]fwb.Card, 0, len(me.SpecialCardSet))
	for i := 0; i < len(me.SpecialCardSet); i++ {
		specialCards = append(specialCards, me.cards[me.SpecialCardSet[i]])
	}

	if me.shuffledCardSet == nil || len(me.shuffledCardSet) < me.OptionalCardsPerRound {
		me.shuffle()
	}

	basicCards := make([]fwb.Card, 0, len(me.BasicCardSet))

	for i := 0; i < len(me.BasicCardSet); i++ {
		basicCards = append(basicCards, me.cards[me.BasicCardSet[i]])
	}

	shuffledCards := make([]fwb.Card, 0, me.OptionalCardsPerRound)

	for i := 0; i < me.OptionalCardsPerRound; i++ {
		shuffledCards = append(shuffledCards, me.cards[me.shuffledCardSet[i]])
	}

	me.shuffledCardSet = me.shuffledCardSet[me.OptionalCardsPerRound:]

	return specialCards, basicCards, shuffledCards
}

//CardSetToString format a card set as string
func CardSetToString(me *[]int) string {
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
