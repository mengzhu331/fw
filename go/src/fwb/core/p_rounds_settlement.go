package core

import (
	"encoding/json"
	"er"
	"fwb"
	"sgs"
)

type playerSData struct {
	Cereals int
	Meat    int
	Sweater int
}

type pstData struct {
	ps map[int]*playerSData
}

func pstInit(me *gameImp) *er.Err {
	me.pd = &pstData{
		ps: make(map[int]*playerSData),
	}

	me.setDCE(fwb.CMD_COMMIT_ROUND_SETTLEMENT, pstOnCommitRoundSettlement)
	me.setTimer(30000, pstOnTimeOut)

	return me.app.SendAllPlayers(sgs.Command{
		ID:     fwb.CMD_ROUND_SETTLEMENT,
		Source: fwb.CMD_SOURCE_APP,
	})
}

func parsePlayerSData(me *gameImp, command sgs.Command) (*playerSData, *er.Err) {
	bytePsd, err := json.Marshal(command.Payload)
	if err != nil {
		return nil, er.Throw(fwb.E_CMD_PAYLOAD_NOT_DECODABLE, er.EInfo{
			"details":   "failed to decode command payload",
			"commandID": command.ID,
			"payload":   command.Payload,
		}).To(me.lg)
	}

	psd := playerSData{}
	err = json.Unmarshal(bytePsd, &psd)
	if err != nil {
		return nil, er.Throw(fwb.E_CMD_PAYLOAD_NOT_DECODABLE, er.EInfo{
			"details":   "failed to decode command payload",
			"commandID": command.ID,
			"payload":   command.Payload,
		}).To(me.lg)
	}

	return &psd, nil
}

func pstOnCommitRoundSettlement(me *gameImp, command sgs.Command) *er.Err {
	pid := command.Source
	if me.gd.GetPDIndex(pid) < 0 {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details": "invalid player ID when commit round settlement",
			"ID":      pid,
		})
	}

	pd := me.pd.(*pstData)

	psd, err := parsePlayerSData(me, command)
	if err != nil {
		return err
	}

	if !validateSettlement(me, pid, *psd) {
		err = me.app.SendToPlayer(pid, sgs.Command{
			ID:     fwb.CMD_ROUND_SETTLEMENT_INVALID,
			Source: fwb.CMD_SOURCE_APP,
		})
		return err
	}

	pd.ps[pid] = psd
	pn := len(me.app.GetPlayers())

	if len(pd.ps) == pn {
		return applyPS(me)
	}

	return err
}

func applyPS(me *gameImp) *er.Err {
	pd := me.pd.(*pstData)

	for k, v := range pd.ps {
		hearts := v.Cereals*2 + v.Meat*3 + v.Sweater*2
		delta := make(fwb.PlayerData, fwb.PD_MAX)
		delta[fwb.PD_PT_HEART] = hearts
		delta[fwb.PD_PT_CEREALS] = -v.Cereals
		delta[fwb.PD_PT_MEAT] = -v.Meat
		delta[fwb.PD_PT_SWEATER] = -v.Sweater

		me.gd.PData[k] = fwb.PDAdd(me.gd.PData[k], delta)
	}

	err := me.app.SendAllPlayers(sgs.Command{
		ID:      fwb.CMD_ROUND_SETTLEMENT_UPDATE,
		Source:  fwb.CMD_SOURCE_APP,
		Payload: me.gd,
	})

	if err.Importance() >= er.IMPT_DEGRADE {
		return err
	}

	return me.gotoPhase(_P_ROUNDS_FINISH)
}

func pstOnTimeOut(me *gameImp, command sgs.Command) (bool, *er.Err) {
	pd := me.pd.(*pstData)

	for k := range pd.ps {
		_, found := pd.ps[k]
		if !found {
			me.app.SendToMockPlayer(k, sgs.Command{
				ID:     fwb.CMD_ROUND_SETTLEMENT,
				Source: fwb.CMD_SOURCE_APP,
			})
		}
	}
	return false, nil
}

func validateSettlement(me *gameImp, playerID int, psd playerSData) bool {
	cerealsHas := me.gd.PData[playerID][fwb.PD_PT_CEREALS]
	meatHas := me.gd.PData[playerID][fwb.PD_PT_MEAT]
	sweaterHas := me.gd.PData[playerID][fwb.PD_PT_SWEATER]
	maxPawns := me.gd.PData[playerID][fwb.PD_MAX_PAWNS]

	return psd.Cereals >= 0 && psd.Cereals <= cerealsHas &&
		psd.Meat >= 0 && psd.Meat <= meatHas &&
		psd.Sweater >= 0 && psd.Sweater <= sweaterHas &&
		psd.Cereals+psd.Meat <= maxPawns &&
		psd.Sweater <= maxPawns
}
