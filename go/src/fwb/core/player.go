package core

import (
	"encoding/json"
	"er"
	"fwb"
	"hlf"
	"sgs"
)

type playerImp struct {
	app  fwb.FwApp
	id   int
	name string
	lg   hlf.Logger
}

func (me *playerImp) ID() int {
	return me.id
}

func (me *playerImp) Name() string {
	return me.name
}

func (me *playerImp) SendCommand(command sgs.Command) *er.Err {
	if command.InCategory(sgs.CMD_C_APP_TO_CLIENT) {
		pld, e := json.Marshal(command.Payload)
		if e != nil {
			return er.Throw(fwb.E_CMD_PAYLOAD_NOT_ENCODABLE, er.EInfo{
				"details": "payload with command sent to client is not able to be encoded to json",
				"payload": command.Payload,
			}).To(me.lg)
		}

		return me.app.GetSession().ForwardToClient(me.id, sgs.Command{
			ID:      sgs.CMD_FORWARD_TO_CLIENT,
			Payload: pld,
		})
	}

	if command.ID == sgs.CMD_FORWARD_TO_APP {
		actualCmd := sgs.Command{}
		cmdBytes, ok := command.Payload.([]byte)

		if !ok {
			return er.Throw(fwb.E_INVALID_CMD, er.EInfo{
				"details": "payload to forward is not bytes",
				"payload": command.Payload,
			}).To(me.lg)
		}

		e := json.Unmarshal(cmdBytes, &actualCmd)

		if e != nil {
			return er.Throw(fwb.E_CMD_PAYLOAD_NOT_DECODABLE, er.EInfo{
				"details": "payload with command sent to app is not able to be decoded to sgs.Command",
				"payload": command.Payload,
			}).To(me.lg)
		}

		return me.app.SendToGame(actualCmd)
	}

	return er.Throw(fwb.E_CMD_NOT_EXECUTABLE, er.EInfo{
		"details": "playerImp is not supposed to receive the command",
		"command": command.HexID(),
	})
}

//makePlayer Init a new player
func makePlayer(app fwb.FwApp, id int, name string) fwb.PlayerAgent {
	player := playerImp{}
	player.app = app
	player.id = id
	player.name = name
	player.lg = app.GetLogger()
	return &player
}
