package actn

import (
	"er"
	"fmt"
	"fwb"
	"sgs"
)

type actnSkip int

func actnSkipParser(command sgs.Command) fwb.Action {
	cid := command.Who
	return (*actnSkip)(&cid)
}

func (me *actnSkip) String() string {
	return fmt.Sprintf("Action %v from Player 0x%x", ActionNames[ACTN_SKIP], *me)
}

func (me *actnSkip) ID() int {
	return ACTN_SKIP
}

func (me *actnSkip) ValidateAgainst(gd *fwb.GameData) bool {
	return true
}

func (me *actnSkip) Do(gd *fwb.GameData) *er.Err {
	px := gd.GetPDIndex(int(*me))
	if px < 0 {
		return er.Throw(fwb.E_INVALID_ACTION, er.EInfo{
			"details": "skip action with invalid player ID",
			"player":  *me,
		})
	}
	gd.PData[px][fwb.PD_PAWNS] -= 1
	return nil
}
