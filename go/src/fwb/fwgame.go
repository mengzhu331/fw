package main

import (
	"log"
	"sgs"
)

const (
	CMD_ID_SEND_RESULT uint = 0x0001
)

type fwGame struct {
	appconfig *sgs.AppConfig
	cp        uint64
	pcem      map[uint]phaseMap
	lsm       map[uint64]interface{}
	players   []player
}

func fwGameCreator() sgs.App {
	return &fwGame{}
}

func (this *fwGame) GetName() string {
	return "Family War"
}

func (this *fwGame) Init(appConfig *sgs.AppConfig) error {
	log.Println("Init FW Game")

	this.lsm = make(map[uint64]interface{})

	this.appconfig = appConfig
	this.appconfig.CmdMap[sgs.CMD_ANY] = phasedCmdExecutor

	this.mapCmdExecutor(sgs.CMD_TICK, P_GAME_ROUNDS, p_game_rounds_tick)
	this.mapCmdExecutor(sgs.CMD_TICK, P_ROUNDS_TURNS, p_rounds_turns_tick)
	this.mapCmdExecutor(CMD_INIT, P_GAME_START, p_game_start_init)
	this.mapCmdExecutor(CMD_GAME_START_ACK, P_GAME_START, p_game_start_game_start_ack)
	this.mapCmdExecutor(CMD_INIT, P_GAME_ROUNDS, p_game_rounds_init)

	for _, c := range this.appconfig.Clients {
		this.players = append(this.players, &remotePlayer{
			client: c,
			df:     playerDF{},
			game:   this,
		})
	}

	this.gotoPhase(P_GAME_START)
	return nil
}

func (this *fwGame) getURI() string {
	return "APP://" + this.GetName()
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

func phasedCmdExecutor(game sgs.App, cmd sgs.Command) error {
	this, ok := game.(*fwGame)
	if !ok {
		return MakeFwErrorByCode(EC_INVALID_GAME_THIS)
	}

	if this.pcem[cmd.ID] != nil {
		return phaseDispatcherSimple(this, this.pcem[cmd.ID], this.cp, cmd)
	}
	return nil
}

func (this *fwGame) sendCommand(cmd sgs.Command) {
	this.appconfig.S.Exec(cmd)
}

func (this *fwGame) gotoPhase(phase uint64) {
	this.cp = phase
	this.sendCommand(sgs.Command{
		ID:     CMD_INIT,
		Target: this.getURI(),
		Source: this.getURI(),
	})
}
