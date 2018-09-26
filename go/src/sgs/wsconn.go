package sgs

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	clientId int
	conn     *websocket.Conn
}

func (me *wsConn) Send(cmd Command) error {
	_log.Dbg("WS send command: %v", cmd.HexID(), cmd.Source, cmd.Payload)

	text, e := json.Marshal(cmd)
	if e != nil {
		return e
	}

	return me.conn.WriteMessage(websocket.TextMessage, text)
}

func (me *wsConn) Run(ch chan Command, mch chan Command) {
	_log.Inf("WS listening to client: %v", me.clientId)
	for {

		me.conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
		_, message, err := me.conn.ReadMessage()
		if err != nil {
			_log.Ntf("Failed to read client: %v", me.clientId)
			goto __close
		}

		ch <- Command{
			ID:      CMD_FORWARD_TO_APP,
			Payload: message,
		}

		select {
		case mc := <-mch:
			if mc.ID == _CMD_CLOSE_NET_CLIENT {
				_log.Inf("Close WS client: %v", me.clientId)
				goto __close
			}
		case <-time.After(time.Duration(1) * time.Second):
		}
	}

__close:
	me.conn.Close()
	_log.Inf("WS closed: %v", me.clientId)
}

func (me *wsConn) BindClientID(clientID int) {
	_log.Inf("WS connection is binded to client: %v", clientID)
	me.clientId = clientID
}
