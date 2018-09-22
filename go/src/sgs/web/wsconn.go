package web

import (
	"encoding/json"
	"sgs/ssvr"

	"github.com/gorilla/websocket"
)

type wsConn struct {
	clientId int
	conn     *websocket.Conn
}

func (me *wsConn) Send(msg []byte) error {
	_log.Dbg("WS send message: %v", msg)
	return me.conn.WriteMessage(websocket.TextMessage, msg)
}

func (me *wsConn) Run(ch chan []byte) {
	_log.Inf("WS listening to client: %v", me.clientId)
	for {
		_, message, err := me.conn.ReadMessage()
		if err != nil {
			_log.Ntf("Client WS disconnected: %v", me.clientId)
			break
		}

		message, _ = json.Marshal(ssvr.Command{
			ID:      ssvr.CMD_FORWARD_CLIENT,
			Payload: message,
		})

		ch <- message

	}
	me.conn.Close()
	_log.Inf("WS closed: %v", me.clientId)
}

func (me *wsConn) BindClientID(clientID int) {
	_log.Inf("WS connection is binded to client: %v", clientID)
	me.clientId = clientID
}
