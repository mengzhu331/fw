package ssvr

import (
	"encoding/json"
	"er"
	"hlf"
	"strconv"
	"sync"
	"time"
)

var _sessionID int = 0x2000

var _csMutex sync.Mutex

var _currentSession *session

var _sessions = make(map[int]*session)

var _sMutex sync.Mutex

//Session session interface for users
type Session interface {
	CmdChan() chan []byte
	ForwardToClient(command Command) *er.Err
	GetLogger() hlf.Logger
}

type session struct {
	id      int
	cch     chan []byte
	clients clientMap
	app     App
	running bool

	baseTickMs int
	lg         hlf.Logger
}

type commandMap map[int]func(*session, Command) *er.Err

var _cm = commandMap{
	CMD_FORWARD_TO_APP:    execForwardToApp,
	CMD_FORWARD_TO_CLIENT: execForwardToClient,
}

func (me *session) CmdChan() chan []byte {
	return me.cch
}

func (me *session) ForwardToClient(command Command) *er.Err {
	return execForwardToClient(me, command)
}

func (me *session) GetLogger() hlf.Logger {
	return me.lg
}

func (me *session) run() *er.Err {

	me.lg.Inf("Starting session %v", me.id)

	me.cch = make(chan []byte)

	clients := make([]int, 0)
	for _, c := range me.clients {
		if c.conn != nil {
			go c.conn.Run(me.cch)
		}
		clients = append(clients, c.ID())
	}

	me.app = _param.ABF()

	err := me.app.Init(me, clients)
	if (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE {
		me.lg.Err("Session failed to start due to application failed to init")
		return err
	}

	me.running = true

	me.lg.Inf("Session started: %v", me.id)

	me.app.SendCommand(Command{
		ID: CMD_APP_RUN,
	})

	go sessionRoutine(me)

	return nil

}

func sessionRoutine(me *session) {
	t := time.Now()
	for me.running {
		select {
		case <-time.After(time.Duration(me.baseTickMs) * time.Millisecond):
			tt := time.Now()
			dms := int(tt.Sub(t) / time.Millisecond)

			me.app.SendCommand(Command{
				ID:      CMD_TICK,
				Payload: dms,
			})
			t = tt

		case cmdBytes := <-me.cch:

			me.lg.Dbg("receive command: %v", cmdBytes)

			err := me.exec(cmdBytes)
			if (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE {
				me.running = false
			}
		}
	}
}

func (me *session) exec(cmdBytes []byte) *er.Err {
	command := Command{}

	err := json.Unmarshal(cmdBytes, &command)

	if err != nil {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "illegal command for session",
			"command": cmdBytes,
		}).To(me.lg)
	}

	execlet, found := _cm[command.ID]
	if !found {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "invalid command ID",
			"command": CmdHexID(command),
		}).To(me.lg)
	}

	return execlet(me, command)
}

func execForwardToApp(s *session, command Command) *er.Err {
	s.lg.Dbg("Forward command from client: Payload %v", command.Payload)

	cmdBytes, isBytes := command.Payload.([]byte)
	if !isBytes {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "payload from client is not bytes type",
			"payload": command.Payload,
		}).To(s.lg)
	}

	var clientCommand Command
	err := json.Unmarshal(cmdBytes, &clientCommand)

	if err != nil {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "payload from client is not ssvr.Command json",
			"payload": command.Payload,
		}).To(s.lg)
	}

	return s.app.SendCommand(clientCommand)
}

//ForwardToClient app use the interface to forward command
func execForwardToClient(s *session, command Command) *er.Err {
	s.lg.Dbg("Forward command to client: Payload %v", command.Payload)

	plBytes, isBytes := command.Payload.([]byte)
	if !isBytes {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "payload to client is not bytes type",
			"payload": command.Payload,
		}).To(s.lg)
	}

	var pl PlForwardToClient
	err := json.Unmarshal(plBytes, &pl)

	if err != nil {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "invalid payload to client",
			"payload": command.Payload,
		}).To(s.lg)
	}

	client, found := s.clients[pl.ClientID]

	if !found {
		return er.Throw(_E_SESSION_INVALID_CLIENT, er.EInfo{
			"details": "invalid client to forward",
			"client":  pl.ClientID,
		}).To(s.lg)
	}

	if client.conn == nil {
		return er.Throw(_E_NO_DUAL_CONNECTION_SUPPORT, er.EInfo{
			"details": "no dual connection",
		}).To(s.lg)
	}

	err = client.conn.Send(pl.Payload)

	if err != nil {
		return er.Throw(_E_CLIENT_CONNECTION_FAIL, er.EInfo{
			"details": "failed to interact with client",
			"client":  pl.ClientID,
		}).To(s.lg)
	}

	return nil
}

func makeSession(sessionID int) *session {
	_slog.Dbg("Make session: ID %v baseTickMs %v", sessionID, _param.BaseTickMs)

	return &session{
		id:      sessionID,
		cch:     make(chan []byte),
		clients: make(clientMap),

		baseTickMs: _param.BaseTickMs,
		lg:         _slog.Child("Session " + strconv.Itoa(_sessionID)),
	}
}
