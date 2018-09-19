package ssvr

import (
	"sync"
)

//NetConn network interface independent of protocol
type NetConn interface {
	Send(msg string) error
	Run(ch chan string)
}

//NetClient external interface for netClient
type NetClient interface {
	ID() int
	Username() string
}

type netClient struct {
	id       int
	username string
	conn     NetConn
}

func (me *netClient) ID() int {
	return me.id
}

func (me *netClient) Username() string {
	return me.username
}

type clientMap map[int]netClient

var _clients = make(clientMap)
var _clientID int = 0x8000

var _cMutex = sync.Mutex{}
