package main

import (
	"log"
	"sgs"
	"time"
)

const (
	CMD_ID_SEND_RESULT uint = 0x0001
)

const (
	MAX_RESPONSE_GAME_START time.Duration = time.Duration(32) * time.Second
)

type fwGame struct {
	fw   IFW
	cp   uint64
	pcem map[uint]phaseMap
	dcem map[uint]fwGameExecutor
	lsm  map[uint64]interface{}

	sacs []actionCard
	dacs []actionCard

	df gameDF

	resTimeOut int
	timerSet   time.Time

	tm map[uint64]int
}

type gameDF struct {
	cr int
}

func (this *fwGame) setTm() {
	this.tm[P_GAME_START] = 32
}

func (this *fwGame) Init(fw IFW) error {
	log.Println("Init FW Game")

	this.sacs = []actionCard{
		acFarm,
		acFeedSheep,
		acDoParttimeJob,
		acTakeVacation,
		acTrade,
		acEmploy,
	}

	this.setTm()

	this.fw = fw
	this.df.cr = 0

	this.lsm = make(map[uint64]interface{})

	this.mapCmdExecutor(sgs.CMD_TICK, P_ROUNDS_TURNS, p_rounds_turns_tick)
	this.mapCmdExecutor(CMD_ENTER, P_GAME_START, p_game_start_enter)
	this.mapCmdExecutor(CMD_GAME_START_ACK, P_GAME_START, p_game_start_game_start_ack)
	this.mapCmdExecutor(CMD_CLAIM_ACTION_CARD, P_ROUNDS_TURNS, p_rounds_turns_claim_action_card)

	this.dcem[sgs.CMD_TICK] = defaultTick

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
	var err error
	if this.pcem[cmd.ID] != nil {
		err = phaseDispatcherSimple(this, this.pcem[cmd.ID], this.cp, cmd)
	}

	if err != nil {
		return err
	}

	f := this.dcem[cmd.ID]
	if f != nil {
		return f(this, cmd)
	}
	return nil
}

func (this *fwGame) gotoPhase(phase uint64) {
	this.cp = phase
	this.SendCommand(sgs.Command{
		ID:     CMD_ENTER,
		Target: makeCommandParticipantUri(TARGET_GAME, ""),
		Source: makeCommandParticipantUri(TARGET_GAME, ""),
	})
	this.timerSet = time.Now()
	this.resTimeOut = this.tm[phase]
}

func defaultTick(this *fwGame, cmd sgs.Command) error {
	if this.timerSet.IsZero() {
		c := time.Now()
		if c.Sub(this.timerSet) > time.Duration(this.resTimeOut)*time.Second {
			return this.fw.SendCommand(sgs.Command{
				ID:     CMD_GAME_END_PLAYER_NO_RESPONSE,
				Source: TARGET_GAME,
				Target: sgs.TARGET_SYS_APP,
			})
		}
	}
	return nil
}
