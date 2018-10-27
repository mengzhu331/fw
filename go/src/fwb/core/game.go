package core

import (
	"er"
	"fwb"
	"fwb/actn"
	"fwb/cards"
	"hlf"
	"sgs"
	"strconv"
	"sutil"

	"github.com/google/uuid"
)

var _gameLog = hlf.MakeLogger("Games")

type timer struct {
	intervalMS int
	elapsedMS  int
	te         timerExe
	id         int
}

type gameConf struct {
	MaxPawn   int
	MinRounds int
	StartGold int
}

type gameImp struct {
	app      fwb.FwApp
	lg       hlf.Logger
	alg      hlf.Logger
	cm       fwb.CardManager
	ap       *actn.ActionParser
	profile  string
	conf     gameConf
	gameuuid uuid.UUID
	timerID  int

	//phase data
	dynamicCmdMap map[int]cmdExe
	timers        []*timer
	pd            interface{}
	phs           phase
	newPhs        phase

	//global game data
	gd        fwb.GameData
	turnOrder []int
}

type cmdExe func(*gameImp, sgs.Command) *er.Err

type timerExe func(*gameImp, sgs.Command) *er.Err

type enterPhase func(*gameImp) *er.Err

var _phaseEnterMap map[phase]enterPhase

var _defaultCmdMap = map[int]cmdExe{
	sgs.CMD_APP_RUN: onRun,
	sgs.CMD_TICK:    onTickDefault,
}

func (me *gameImp) setupPhases() {
	if _phaseEnterMap != nil {
		return
	}
	_phaseEnterMap = map[phase]enterPhase{}
	_phaseEnterMap[_P_GAME_START] = pgsInit
	_phaseEnterMap[_P_ROUNDS_START] = prsInit
	_phaseEnterMap[_P_ROUNDS_TURNS] = prtInit
	_phaseEnterMap[_P_ROUNDS_SETTLEMENT] = pstInit
	_phaseEnterMap[_P_ROUNDS_FINISH] = prfInit
	_phaseEnterMap[_P_GAME_SETTLEMENT] = pgstInit
	_phaseEnterMap[_P_GAME_FINISH] = pgfInit
	me.phs = _P_INIT
	me.newPhs = me.phs
}

func makeGame(app fwb.FwApp, profile string) (*gameImp, *er.Err) {
	game := gameImp{}

	game.app = app

	game.lg = app.GetLogger()

	game.ap = actn.MakeActionParser(&game)

	if (len(app.GetPlayers()) == 3 && profile == fwb.PROFILE_3PVP) ||
		(len(app.GetPlayers()) == 2 && profile == fwb.PROFILE_2PVP) {
		game.profile = profile
	} else {
		er.Throw(fwb.E_PROFILE_MISMATCH, er.EInfo{
			"details": "profile and player number mismatch",
			"player":  strconv.Itoa(len(app.GetPlayers())),
			"profile": profile,
		}).To(game.lg)
	}

	game.cm = cards.MakeCardManager()

	cfile := "./conf/profiles/" + profile + "/game.conf"
	err := sutil.LoadConfFile(cfile, &game.conf)
	if err != nil {
		return nil, er.Throw(fwb.E_MISSING_GAME_SETTINGS, er.EInfo{
			"details": "failed to load settings",
			"file":    cfile,
		}).To(game.lg)
	}

	e := game.cm.LoadCards(profile)

	if e != nil {
		e.To(game.lg)
	}
	game.setupPhases()
	return &game, e
}

func (me *gameImp) GetLogger() hlf.Logger {
	return me.lg
}

func (me *gameImp) GetProfile() string {
	return me.profile
}

func (me *gameImp) GameOver(reasonCode int, details interface{}) *er.Err {

	me.lg.Inf("Game Over: reason %v", reasonCode)

	payload := map[string](interface{}){
		"ReasonCode": reasonCode,
		"Details":    details,
	}

	err := me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_GAME_OVER,
		Who:     fwb.CMD_WHO_APP,
		Payload: payload,
	})

	if err.Importance() >= er.IMPT_UNRECOVERABLE {
		return err
	}

	err = err.Push(me.app.SendToSession(sgs.Command{
		ID:      sgs.CMD_APP_CLOSE,
		Payload: reasonCode,
	}))

	return err
}

func (me *gameImp) SendCommand(command sgs.Command) *er.Err {
	dce, founddce := me.dynamicCmdMap[command.ID]
	var err *er.Err

	if founddce {
		err = dce(me, command)
		if err.Importance() >= er.IMPT_DEGRADE {
			return err
		}
	}

	exec, found := _defaultCmdMap[command.ID]

	if found {
		err = err.Push(exec(me, command))
		if err.Importance() >= er.IMPT_DEGRADE {
			return err
		}
	}

	if !found && !founddce {
		return er.Throw(fwb.E_CMD_NOT_EXECUTABLE, er.EInfo{
			"details": "gameImp is not supposed to receive the command",
			"command": command.HexID(),
			"phase":   me.phs,
		}).To(me.lg)
	}

	return me.switchPhase()
}

func (me *gameImp) gotoPhase(p phase) *er.Err {
	me.lg.Dbg("Go to phase from %x to %x", me.phs, p)
	me.newPhs = p
	return nil
}

func (me *gameImp) switchPhase() *er.Err {
	if me.phs == me.newPhs {
		return nil
	}

	me.lg.Dbg("Switch phase from %x to %x", me.phs, me.newPhs)
	me.dynamicCmdMap = make(map[int]cmdExe)
	me.timers = make([]*timer, 0)

	exec, found := _phaseEnterMap[me.newPhs]
	if found {
		me.phs = me.newPhs
		return exec(me)
	}

	return er.Throw(fwb.E_INVALID_GAME_PHASE, er.EInfo{
		"details": "cannot switch to invalid game phase",
		"phase":   me.newPhs,
	}).To(me.lg)
}

func onRun(me *gameImp, command sgs.Command) *er.Err {
	me.gotoPhase(_P_GAME_START)
	return nil
}

func onTickDefault(me *gameImp, command sgs.Command) *er.Err {
	if me.newPhs != me.phs {
		return me.switchPhase()
	}

	deltaMS := command.Payload.(int)

	return tickTimer(me, deltaMS)
}

func (me *gameImp) needSwitchPhase() bool {
	return me.phs != me.newPhs
}

func tickTimer(me *gameImp, deltaMS int) *er.Err {
	if len(me.timers) < 1 {
		return nil
	}

	var err *er.Err
	timersLeft := make([]*timer, 0)
	for _, t := range me.timers {
		t.elapsedMS += deltaMS
		if t.elapsedMS >= t.intervalMS {
			err = err.Push(t.te(me, sgs.Command{
				ID:      fwb.CMD_TIMER,
				Who:     fwb.CMD_WHO_APP,
				Payload: t.id,
			}))

			if err.Importance() >= er.IMPT_DEGRADE {
				return err
			}

			if me.needSwitchPhase() {
				break
			}
		} else {
			timersLeft = append(timersLeft, t)
		}
	}

	me.timers = timersLeft
	return nil
}

func (me *gameImp) setDCE(cmdId int, dce cmdExe) {
	me.dynamicCmdMap[cmdId] = dce
}

func (me *gameImp) unsetDCE(cmdId int) {
	delete(me.dynamicCmdMap, cmdId)
}

func (me *gameImp) setTimer(intervalMS int, te timerExe) int {
	t := timer{
		intervalMS: intervalMS,
		te:         te,
		id:         me.timerID,
	}
	me.timerID++

	var i int
	for i = range me.timers {
		if me.timers[i] == nil {
			break
		}
	}
	if i == len(me.timers) {
		me.timers = append(me.timers, &t)
	} else {
		me.timers[i] = &t
	}
	return t.id
}

func (me *gameImp) unsetTimer(id int) {
	var i int
	for i = range me.timers {
		if me.timers[i].id == id {
			break
		}
	}

	if i < len(me.timers) {
		me.timers = append(me.timers[:i], me.timers[i+1:]...)
	}
}
