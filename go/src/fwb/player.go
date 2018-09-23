package main

import (
	"encoding/json"
	"er"
	"hlf"
	"sgs/ssvr"
)

type player interface {
	init(app fwApp, id int)
	sendCommand(command ssvr.Command) *er.Err
}

type playerImp struct {
	app fwApp
	id  int
	lg  hlf.Logger
}

func (me *playerImp) init(app fwApp, id int) {
	me.app = app
	me.id = id

	me.lg = app.getLogger()
}

func (me *playerImp) sendCommand(command ssvr.Command) *er.Err {
	if ssvr.CmdInCategory(command, ssvr.CMD_C_APP_TO_CLIENT) {
		pld, e := json.Marshal(command.Payload)
		if e != nil {
			return er.Throw(_E_CMD_PAYLOAD_NOT_ENCODABLE, er.EInfo{
				"details": "payload with command sent to client is not able to be encoded to json",
				"payload": command.Payload,
			}).To(me.lg)
		}

		return me.app.getSession().ForwardToClient(ssvr.Command{
			ID:      ssvr.CMD_FORWARD_TO_CLIENT,
			Payload: pld,
		})
	}

	if command.ID == ssvr.CMD_FORWARD_TO_APP {
		actualCmd := ssvr.Command{}
		cmdBytes, ok := command.Payload.([]byte)

		if !ok {
			return er.Throw(_E_INVALID_CMD, er.EInfo{
				"details": "payload to forward is not bytes",
				"payload": command.Payload,
			}).To(me.lg)
		}

		e := json.Unmarshal(cmdBytes, &actualCmd)

		if e != nil {
			return er.Throw(_E_CMD_PAYLOAD_NOT_DECODABLE, er.EInfo{
				"details": "payload with command sent to app is not able to be decoded to ssvr.Command",
				"payload": command.Payload,
			}).To(me.lg)
		}

		return me.app.getGame().sendCommand(actualCmd)
	}

	return er.Throw(_E_CMD_NOT_EXEC, er.EInfo{
		"details": "playerImp is not supposed to receive the command",
		"command": ssvr.CmdHexID(command),
	})
}

func makePlayer() player {
	return &playerImp{}
}
