package main

import (
	"encoding/json"
	"er"
	"hlf"
	"sgs"
)

type player interface {
	init(app fwApp, id int)
	sendCommand(command sgs.Command) *er.Err
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

func (me *playerImp) sendCommand(command sgs.Command) *er.Err {
	if command.InCategory(sgs.CMD_C_APP_TO_CLIENT) {
		pld, e := json.Marshal(command.Payload)
		if e != nil {
			return er.Throw(_E_CMD_PAYLOAD_NOT_ENCODABLE, er.EInfo{
				"details": "payload with command sent to client is not able to be encoded to json",
				"payload": command.Payload,
			}).To(me.lg)
		}

		return me.app.getSession().ForwardToClient(sgs.Command{
			ID:      sgs.CMD_FORWARD_TO_CLIENT,
			Payload: pld,
		})
	}

	if command.ID == sgs.CMD_FORWARD_TO_APP {
		actualCmd := sgs.Command{}
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
				"details": "payload with command sent to app is not able to be decoded to sgs.Command",
				"payload": command.Payload,
			}).To(me.lg)
		}

		return me.app.getGame().sendCommand(actualCmd)
	}

	return er.Throw(_E_CMD_NOT_EXEC, er.EInfo{
		"details": "playerImp is not supposed to receive the command",
		"command": command.HexID(),
	})
}

func makePlayer() player {
	return &playerImp{}
}
