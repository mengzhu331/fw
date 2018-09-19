package ssvr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
				Param: map[string]interface{}{
					"DeltaMs": dms,
				},
			})
		case cmdStr := <-me.cch:
			err := me.exec(cmdStr)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

}

func (me *session) exec(cmdStr string) error {
	command := Command{}

	err := json.Unmarshal([]byte(cmdStr), &command)

	if err != nil {
		return err
	}

	if command.InCate(CMD_C_CLIENT) {
		return me.app.SendCommand(command)
	} else if command.InCate(CMD_C_APP) {
		ciu, found := command.Param["ClientID"]

		if !found {
			return errors.New("Internal error, invalid command: " + cmdStr)
		}

		cis, valid := ciu.(string)
		if !valid {
			return errors.New("Internal error, invalid client id: " + cmdStr)
		}

		clti, e := strconv.Atoi(cis)

		if e != nil {
			return errors.New("Internal error, invalid client id: " + cmdStr)
		}

		c, found := me.clients[clti]

		if !found {
			return errors.New("Internal error, invalid client: " + cmdStr)
		}
		return c.conn.Send(cmdStr)
	}

	return nil
}
