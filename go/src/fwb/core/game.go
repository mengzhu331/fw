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
)

type timer struct {
	intervalMS int
	elapsedMS  int
	te         timerExe
}

type gameConf struct {
	MaxPawn   int
	MinRounds int
}

type gameImp struct {
	app     fwb.FwApp
	lg      hlf.Logger
	cm      fwb.CardManager
	ap      *actn.ActionParser
	profile string
	conf    gameConf

	//phase data
	dynamicCmdMap map[int]cmdExe
	timers        []*timer
	pd            interface{}
	phs           phase

	//global game data
	gd        fwb.GameData
	turnOrder []int
}

type cmdExe func(*gameImp, sgs.Command) *er.Err

type timerExe func(*gameImp, sgs.Command) (bool, *er.Err)

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

	game.cm = cards.MakeCardManager(&game)

	cfile := "./" + profile + "/game.conf"
	err := sutil.LoadConfFile(cfile, &game.conf)
	if err != nil {
		return nil, er.Throw(fwb.E_MISSING_GAME_SETTINGS, er.EInfo{
			"details": "failed to load settings",
			"file":    cfile,
		})
	}

	e := game.cm.LoadCards(profile)

	return &game, e
}

func (me *gameImp) GetLogger() hlf.Logger {
	return me.lg
}

func (me *gameImp) GetProfile() string {
	return me.profile
}

func (me *gameImp) GameOver(reasonCode int, details interface{}) *er.Err {

	payload := map[string](interface{}){
		"ReasonCode": reasonCode,
		"Details":    details,
	}

	err := me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_GAME_OVER,
		Source:  fwb.CMD_SOURCE_APP,
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

	return err
}

func (me *gameImp) gotoPhase(p phase) *er.Err {
	me.dynamicCmdMap = make(map[int]cmdExe)
	me.timers = make([]*timer, 0)
	me.phs = p
	exec, found := _phaseEnterMap[p]
	if found {
		return exec(me)
	}
	return nil
}

func onRun(me *gameImp, command sgs.Command) *er.Err {
	me.gotoPhase(_P_GAME_START)
	return nil
}

func onTickDefault(me *gameImp, command sgs.Command) *er.Err {
	var err *er.Err

	deltaMS := command.Payload.(int)

	for id, t := range me.timers {
		if t != nil {
			t.elapsedMS += deltaMS
			if t.elapsedMS >= t.intervalMS {
				c, e := t.te(me, sgs.Command{
					ID:      fwb.CMD_TIMER,
					Source:  fwb.CMD_SOURCE_APP,
					Payload: id,
				})
				if !c {
					me.timers[id] = nil
				}
				err = err.Push(e)
				if err.Importance() >= er.IMPT_DEGRADE {
					return err
				}
			}
		}
	}

	return err
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
	}

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
	return i
}

func (me *gameImp) unsetTimer(id int) {
	me.timers[id] = nil
}
