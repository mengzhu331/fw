package main

import (
	"sgs"
	"testing"
)

type TestFW struct {
	game    *fwGame
	players []player
}

func (this *TestFW) GetGame() *fwGame {
	return this.game
}

func (this *TestFW) GetPlayers() []player {
	return this.players
}

func (this *TestFW) SendCommand(command sgs.Command) error {
	return nil
}

func TestSmokeFwGame(t *testing.T) {
	gametested := fwGame{}
	fw := &TestFW{
		game: &gametested,
	}

	gametested.Init(fw)
}
