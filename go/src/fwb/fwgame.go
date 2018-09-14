package main

import (
	"log"
	"sgs"
)

const (
	CMD_ID_SEND_RESULT uint = 0x0001
)

type fwGame struct {
	fw   IFW
	cp   uint64
	pcem map[uint]phaseMap
	lsm  map[uint64]interface{}
}

func (this *fwGame) Init(fw IFW) error {
	log.Println("Init FW Game")

	this.fw = fw

	this.lsm = make(map[uint64]interface{})

	this.mapCmdExecutor(sgs.CMD_TICK, P_GAME_ROUNDS, p_game_rounds_tick)
	this.mapCmdExecutor(sgs.CMD_TICK, P_ROUNDS_TURNS, p_rounds_turns_tick)
	this.mapCmdExecutor(CMD_INIT, P_GAME_START, p_game_start_init)
	this.mapCmdExecutor(CMD_GAME_START_ACK, P_GAME_START, p_game_start_game_start_ack)
	this.mapCmdExecutor(CMD_INIT, P_GAME_ROUNDS, p_game_rounds_init)

	this.gotoPhase(P_GAME_START)
	return nil
}

func (this *fwGame) mapCmdExecutor(cmd uint, phase uint64, exe fwGameExecutor) {
	if this.pcem == nil {
		this.pcem = make(map[uint]phaseMap)
	}

	if this.pcem[cmd] == nil {
		this.pcem[cmd] = make(phaseMap)
	}

	this.pcem[cmd][phase] = exe
}

func (this *fwGame) SendCommand(cmd sgs.Command) error {
	if this.pcem[cmd.ID] != nil {
		return phaseDispatcherSimple(this, this.pcem[cmd.ID], this.cp, cmd)
	}
	return nil
}

func (this *fwGame) gotoPhase(phase uint64) {
	this.cp = phase
	this.SendCommand(sgs.Command{
		ID:     CMD_INIT,
		Target: makeCommandParticipantUri(TARGET_GAME, ""),
		Source: makeCommandParticipantUri(TARGET_GAME, ""),
	})
}
