package fwb

import (
	"er"
	"hlf"
	"sgs"
)

//FwApp is the mediator and manager of the FW backend modules.
type FwApp interface {
	GetLogger() hlf.Logger
	GetPlayers() []PlayerAgent
	SendAllPlayers(command sgs.Command) *er.Err
	SendToPlayer(playerID int, command sgs.Command) *er.Err
	SendToMockPlayer(playerID int, command sgs.Command)
	SendToGame(command sgs.Command) *er.Err
	SendToSession(command sgs.Command) *er.Err
	GetSession() sgs.Session
	//GetPlayer retrieve a specific PlayerAgent matching a playerID
	GetPlayer(playerID int) PlayerAgent
}

//PlayerAgent is the proxy of remote player or computer mocked player
type PlayerAgent interface {
	SendCommand(command sgs.Command) *er.Err
	ID() int
	Name() string
}

//Game is the model of the gameplay system
type Game interface {
	SendCommand(command sgs.Command) *er.Err
	GetProfile() string
	GetLogger() hlf.Logger
	GameOver(reasonCode int, details interface{}) *er.Err
}

//Card is the game object describing the currently avaliable actions that is accesible to all game participants
type Card struct {
	ID          int
	MaxSlot     int
	Pawns       []int
	PawnPerTurn int
}

//CardManager is the module in charge of loading cards from configuration files, and it makes card set for each round according
//to a card set making algorithm
type CardManager interface {
	LoadCards(conf string) *er.Err
	MakeCardSet() ([]Card, []Card, []Card)
}

//Action player Action abstract interface
type Action interface {
	String() string
	ID() int
	ValidateAgainst(*GameData) bool
	Do(*GameData) *er.Err
}
