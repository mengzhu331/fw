package main

import (
	"log"
	"sgs"
)

func p_game_rounds_init(this *fwGame, cmd sgs.Command) error {
	log.Println("P_GAME_ROUNDS init")
	return nil
}

func p_game_rounds_tick(this *fwGame, cmd sgs.Command) error {
	log.Println("P_GAME_ROUNDS tick")
	return nil
}
