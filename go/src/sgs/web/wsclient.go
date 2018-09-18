package web

import (
	"github.com/gorilla/websocket"
)

type wsClient struct {
	conn *websocket.Conn
	ch   chan string
}

func (me *wsClient) Send(msg string) error {
	return me.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (me *wsClient) BindChan(ch chan string) {
	me.ch = ch
}
