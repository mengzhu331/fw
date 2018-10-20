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
	_log.Dbg("WS send command: 0x%v, 0x%x, %v", cmd.HexID(), cmd.Who, cmd.Payload)

	text, e := json.Marshal(cmd)
	if e != nil {
		return e
	}

	return me.conn.WriteMessage(websocket.TextMessage, text)
}

func (me *wsConn) Run(ch chan Command, mch chan Command) {
	_log.Inf("WS listening to client: %v", me.clientId)

	readch := make(chan Command)

	go func() {
		for {
			_, message, err := me.conn.ReadMessage()
			if err != nil {
				_log.Ntf("Failed to read client: %v, %v", me.clientId, err.Error())
				mch <- Command{
					ID:      _CMD_CLOSE_NET_CLIENT,
					Who:     _CMD_WHO_WSCONN,
					Payload: err.Error(),
				}
			}

			readch <- Command{
				ID:      CMD_FORWARD_TO_APP,
				Payload: message,
			}
		}
	}()

	for {
		select {
		case mc := <-mch:
			if mc.ID == _CMD_CLOSE_NET_CLIENT {
				_log.Inf("Close WS client: %v, %v, %v", me.clientId, mc.Who, mc.Payload.(string))
				goto __close
			}
		case command := <-readch:
			ch <- command
		case <-time.After(time.Duration(1) * time.Second):
		}
	}

__close:
	me.conn.Close()
	_log.Inf("WS closed: %v", me.clientId)
}
