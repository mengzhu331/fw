package server

import (
	"sync"
)

//NetConn network interface independent of protocol
type NetConn interface {
	Send(msg string) error
	BindChan(ch chan string)
}

type netClient struct {
	id       int
	username string
	conn     NetConn
}

type clientMap map[int]netClient

var _clients clientMap
var _clientID int = 0x8000

var _cMutex = sync.Mutex{}
