package web

import (
	"github.com/gorilla/websocket"
)

type wsConn struct {
	conn *websocket.Conn
}

func (me *wsConn) Send(msg string) error {
	return me.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (me *wsConn) Run(ch chan string) {
	for {
		_, message, err := me.conn.ReadMessage()
		if err != nil {
			break
		}
		ch <- string(message)
	}
}
