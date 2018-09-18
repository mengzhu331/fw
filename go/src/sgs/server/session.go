package server

import (
	"sync"
)

type session struct {
	ID      int
	cch     chan string
	clients clientMap
	running bool
}

var _sessionID int = 0x2000

var _currentSession *session

var _sessions = make(map[int]*session)

var _csMutex sync.Mutex

var _sMutex sync.Mutex

func (me *session) run() {}
