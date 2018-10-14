package actn

import (
	"er"
	"fmt"
	"fwb"
	"sgs"
)

type actnSkip int

func actnSkipParser(command sgs.Command) fwb.Action {
	cid := command.Source
	return (*actnSkip)(&cid)
}

func (me *actnSkip) String() string {
	return fmt.Sprintf("[Action %v from Player %v]", _actionNames[ACTN_SKIP], *me)
}

func (me *actnSkip) ID() int {
	return ACTN_SKIP
}

func (me *actnSkip) ValidateAgainst(gd *fwb.GameData) bool {
	return true
}

func (me *actnSkip) Do(gd *fwb.GameData) *er.Err {
	gd.PData[*me][fwb.PD_PAWNS] -= 1
	return nil
}
