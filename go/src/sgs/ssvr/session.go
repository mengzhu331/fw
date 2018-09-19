package ssvr

import (
	"sync"
	"time"
)

type session struct {
	ID      int
	cch     chan string
	clients clientMap
	app     App
	running bool
}

var _sessionID int = 0x2000

var _currentSession *session

var _sessions = make(map[int]*session)

var _csMutex sync.Mutex

var _sMutex sync.Mutex

func (me *session) run(baseTickMs int) {
	me.cch = make(chan string)
	clients := make([]NetClient, 0)
	for _, c := range me.clients {
		go c.conn.Run(me.cch)
		clients = append(clients, &c)
	}

	me.app = _param.ABF()
	me.app.Init(me.cch, clients)

	me.running = true

	t := time.Now()

	for me.running {
		select {
		case <-time.After(time.Duration(baseTickMs) * time.Millisecond):
			tt := time.Now()
			dms := int(tt.Sub(t) / time.Millisecond)

			me.app.SendCommand(Command{
				ID: CMD_TICK,
				Param: TickParam{
					DeltaMs: dms,
				},
			})
		}
	}

}
