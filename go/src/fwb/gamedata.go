package fwb

import (
	"math"
)

const (
	//PD_CLIENT_ID index in PlayerData for player client ID
	PD_CLIENT_ID = iota

	//PD_PT_HEART index in PlayerData for property Heart
	PD_PT_HEART

	//PD_PT_GOLD index in PlayerData for property Gold
	PD_PT_GOLD

	//PD_PT_CEREALS index in PlayerData for property Cereals
	PD_PT_CEREALS

	//PD_PT_MEAT index in PlayerData for property Meat
	PD_PT_MEAT

	//PD_PT_WOOL index in PlayerData for property Woo
	PD_PT_WOOL

	//PD_PT_SWEATER index in PlayerData for property Sweater
	PD_PT_SWEATER

	//PD_PT_BEER index in PlayerData for property Wine
	PD_PT_BEER

	//PD_SK_STRENGTH index in PlayerData for skill Strength
	PD_SK_STRENGTH

	//PD_SK_KNOWLEDGE index in PlayerData for skill Knowledge
	PD_SK_KNOWLEDGE

	//PD_SK_INTELLIGENCE index in PlayerData for skill Intelligence
	PD_SK_INTELLIGENCE

	//PD_HOUSE_LV index in PlayerData for the level of the player house
	PD_HOUSE_LV

	//PD_MAX_PAWNS index in PlayerData for the max pawns
	PD_MAX_PAWNS

	//PD_PAWNS index in PlayerData for the pawns left
	PD_PAWNS

	//PD_MAX max index in PlayerData
	PD_MAX
)

//PlayerData is the model of runtime player data
type PlayerData []int

//GameData is the data object for exchanging runtime game information between systems and modules
type GameData struct {
	Round     int
	MaxPawn   int
	MinRounds int
	Cards     []Card
	PData     []PlayerData
}

func (me *PlayerData) init(clientID int, maxPawn int) {
	*me = make(PlayerData, PD_MAX)
	(*me)[PD_CLIENT_ID] = clientID
	(*me)[PD_PT_HEART] = 0
	(*me)[PD_PT_GOLD] = 30
	(*me)[PD_PT_CEREALS] = 0
	(*me)[PD_PT_MEAT] = 0
	(*me)[PD_PT_WOOL] = 0
	(*me)[PD_PT_SWEATER] = 0
	(*me)[PD_PT_BEER] = 0

	(*me)[PD_SK_STRENGTH] = 0
	(*me)[PD_SK_KNOWLEDGE] = 0
	(*me)[PD_SK_INTELLIGENCE] = 0

	(*me)[PD_HOUSE_LV] = 0

	(*me)[PD_MAX_PAWNS] = maxPawn
	(*me)[PD_PAWNS] = maxPawn
}

//Init create players data and set all values of game data to default
func (me *GameData) Init(players []PlayerAgent, maxPawn int, minRounds int) {
	me.Round = 0
	me.MaxPawn = maxPawn
	me.MinRounds = minRounds
	me.Cards = nil
	me.PData = make([]PlayerData, len(players))
	for i := range me.PData {
		me.PData[i].init(players[i].ID(), maxPawn)
	}
}

//PDAdd sum up fields of twp player data object
func PDAdd(left PlayerData, right PlayerData) PlayerData {
	sum := make(PlayerData, int(math.Max(float64(len(left)), float64(len(right)))))
	for i := range sum {
		if i < len(left) {
			sum[i] += left[i]
		}

		if i < len(right) {
			sum[i] += right[i]
		}
	}
	return sum
}

//AllAboveZero all fields of player data object valued above zero
func (me PlayerData) AllAboveZero() bool {
	for k, v := range me {
		if k != PD_CLIENT_ID && v < 0 {
			return false
		}
	}
	return true
}

//GetPDIndex get the player ID corresponded player data index in the game data object
func (me *GameData) GetPDIndex(playerID int) int {
	for i, v := range me.PData {
		if v[PD_CLIENT_ID] == playerID {
			return i
		}
	}
	return -1
}
