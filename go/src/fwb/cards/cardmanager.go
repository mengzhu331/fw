package cards

import (
	"encoding/json"
	"er"
	"fwb"
	"hlf"
	"io/ioutil"
	"math/rand"
	"strconv"
)

type cardManager struct {
	lg hlf.Logger

	BasicCardSet []int

	OptionalCardSet []int

	shuffledCardSet []int

	OptionalCardsPerRound int

	cards []fwb.Card
}

//MakeCardManager instantiate a cardManage object
func MakeCardManager(game fwb.Game) fwb.CardManager {
	return &cardManager{
		lg: game.GetLogger(),
	}
}

//LoadCards implement CardManager.LoadCards() to load cards for a profile
func (me *cardManager) LoadCards(profile string) *er.Err {
	path := "./profiles/" + profile + "/" + "cards.conf"
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return er.Throw(fwb.E_MISSING_GAME_SETTINGS, er.EInfo{
			"details":  "failed to load card settings",
			"profile":  profile,
			"io error": err.Error(),
		}).To(me.lg)
	}

	err = json.Unmarshal(c, &me.cards)
	if err != nil {
		return er.Throw(fwb.E_INVALID_SETTINGS, er.EInfo{
			"details":      "failed to decode card settings",
			"decode error": err.Error(),
		}).To(me.lg)
	}

	path = "./profiles/" + profile + "/" + "cardset.conf"
	c, err = ioutil.ReadFile(path)
	if err != nil {
		return er.Throw(fwb.E_MISSING_GAME_SETTINGS, er.EInfo{
			"details":  "failed to load cardset settings",
			"profile":  profile,
			"io error": err.Error(),
		}).To(me.lg)
	}

	err = json.Unmarshal(c, me)
	if err != nil {
		return er.Throw(fwb.E_INVALID_SETTINGS, er.EInfo{
			"details":      "failed to decode cardset settings",
			"decode error": err.Error(),
		}).To(me.lg)
	}

	me.lg.Inf("Basic Card Set: %v", cardSetToString(&me.BasicCardSet))
	me.lg.Inf("Optional Card Set: %v", cardSetToString(&me.OptionalCardSet))
	return nil

}

func (me *cardManager) shuffle(cards []int, swaps int) []int {
	shuffled := make([]int, len(cards))
	shuffled = append(shuffled, cards...)
	for i := 0; i < swaps; i++ {
		a := rand.Intn(len(shuffled))
		b := rand.Intn(len(shuffled))
		shuffled[a] = shuffled[b] + shuffled[a]
		shuffled[b] = shuffled[a] - shuffled[b]
		shuffled[a] = shuffled[a] - shuffled[b]
	}
	return shuffled
}

func (me *cardManager) MakeCardSet() ([]fwb.Card, []fwb.Card, []fwb.Card) {
	specialCards := make([]fwb.Card, 0)
	specialCards = append(specialCards, fwb.Card{
		ID:          CARD_VOID,
		MaxSlot:     999,
		Pawns:       make([]int, 0),
		PawnPerTurn: 1,
	})

	if me.shuffledCardSet == nil || len(me.shuffledCardSet) < me.OptionalCardsPerRound {
		me.shuffledCardSet = make([]int, len(me.OptionalCardSet))
		me.shuffledCardSet = me.shuffle(me.shuffledCardSet, len(me.OptionalCardSet)*2)
	}

	basicCards := make([]fwb.Card, len(me.BasicCardSet))

	for i := 0; i < len(me.BasicCardSet); i++ {
		basicCards = append(basicCards, me.cards[me.BasicCardSet[i]])
	}

	shuffledCards := make([]fwb.Card, me.OptionalCardsPerRound)

	for i := 0; i < me.OptionalCardsPerRound; i++ {
		shuffledCards = append(shuffledCards, me.cards[me.shuffledCardSet[i]])
	}

	me.shuffledCardSet = me.shuffledCardSet[me.OptionalCardsPerRound:]

	return specialCards, basicCards, shuffledCards
}

func cardSetToString(me *[]int) string {
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
