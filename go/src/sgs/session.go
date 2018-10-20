package sgs

import (
	"encoding/json"
	"er"
	"hlf"
	"strconv"
	"time"
)

//Session session interface for users
type Session interface {
	CmdChan() chan Command
	ForwardToClient(cid int, command Command) *er.Err
	GetLogger() hlf.Logger
	GetClientName(cid int) string
}

type session struct {
	id  int
	cch chan Command
	mch chan Command

	clients map[int]*netClient
	app     App
	closed  bool
	closing bool
	appRun  bool

	baseTickMs int
	lg         hlf.Logger
}

type commandMap map[int]func(*session, Command) *er.Err

var _cm = commandMap{
	CMD_FORWARD_TO_APP:    execForwardToApp,
	CMD_FORWARD_TO_CLIENT: execForwardToClient,
	_CMD_CLIENT_RECONNECT: execClientReconnect,
	CMD_APP_CLOSE:         execAppClose,
}

func (me *session) CmdChan() chan Command {
	return me.cch
}

func (me *session) GetClientName(cid int) string {
	client, found := me.clients[cid]
	if !found {
		er.Throw(_E_REQUEST_WITH_INVALID_CLIENT_ID, er.EInfo{
			"details":  "invalid client ID when invoking GetClientName()",
			"clientid": cid,
		}).To(me.lg)

		return "unknown"
	}

	return client.username
}

func (me *session) ForwardToClient(cid int, command Command) *er.Err {
	client, found := me.clients[cid]

	if !found {
		return er.Throw(_E_SESSION_INVALID_CLIENT, er.EInfo{
			"details": "invalid client to forward",
			"client":  cid,
		}).To(me.lg)
	}

	e := client.send(command)

	if e != nil {
		return er.Throw(_E_CLIENT_CONNECTION_FAIL, er.EInfo{
			"details": "failed to interact with client",
			"client":  cid,
		}).To(me.lg)
	}

	return nil
}

func (me *session) GetLogger() hlf.Logger {
	return me.lg
}

func (me *session) run(abf AppBuildFunc, profile string) *er.Err {

	me.lg.Inf("Starting session %v", me.id)

	me.cch = make(chan Command, 100)

	clients := make([]int, 0)
	for _, c := range me.clients {
		c.run(me.cch)
		clients = append(clients, c.id)
	}

	me.app = abf()

	err := me.app.Init(me, clients, profile)
	if (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE {
		me.lg.Err("Session failed to start due to application failed to init")
		return err
	}

	go sessionRoutine(me)
	me.lg.Inf("Session started: %v", me.id)

	return nil

}

func (me *session) runApp() *er.Err {
	if me.appRun {
		return nil
	}

	me.appRun = true

	return me.app.SendCommand(Command{
		ID: CMD_APP_RUN,
	})

}

func sessionRoutine(me *session) {
	t := time.Now()
	for !me.closing {
		select {
		case mc := <-me.mch:
			me.lg.Dbg("receive management: %v %v %v", mc.HexID, mc.Who, mc.Payload)

			err := me.exec(mc)
			if (err.Code() & er.E_IMPORTANCE) >= er.IMPT_DEGRADE {
				goto __close
			}

		case <-time.After(time.Duration(me.baseTickMs) * time.Millisecond):
			err := me.runApp()
			if (err.Code() & er.E_IMPORTANCE) >= er.IMPT_DEGRADE {
				goto __close
			}

			tt := time.Now()
			dms := int(tt.Sub(t) / time.Millisecond)

			me.app.SendCommand(Command{
				ID:      CMD_TICK,
				Payload: dms,
			})
			t = tt

		case cmd := <-me.cch:

			me.lg.Dbg("receive command: %v %v %v", cmd.HexID, cmd.Who, cmd.Payload)

			err := me.exec(cmd)
			if (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE {
				goto __close
			}
		}
	}
__close:
	me.lg.Inf("Closing session")
	for _, c := range me.clients {
		c.close()
	}
	me.closed = true
	me.lg.Inf("Session closed")
}

func (me *session) exec(command Command) *er.Err {
	execlet, found := _cm[command.ID]
	if !found {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "invalid command ID",
			"command": command.HexID(),
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

	clientCommand := Command{}

	err := json.Unmarshal(cmdBytes, &clientCommand)

	if err != nil {
		return er.Throw(_E_SESSION_INVALID_COMMAND, er.EInfo{
			"details": "payload from client is not Command json",
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

	return s.ForwardToClient(pl.ClientID, pl.Cmd)
}

func execClientReconnect(s *session, command Command) *er.Err {
	s.lg.Dbg("Client reconnect: client %v", command.Who)

	return s.app.SendCommand(command)
}

func execAppClose(s *session, command Command) *er.Err {
	s.lg.Inf("App is to be closed by itself")
	s.closing = true
	return nil
}

func makeSession(sessionID int, ssrv *sessionServer) *session {
	ssrv.lg.Dbg("Make session: ID %v baseTickMs %v", sessionID, ssrv.param.BaseTickMs)

	return &session{
		id:      sessionID,
		cch:     make(chan Command),
		mch:     make(chan Command),
		clients: make(map[int]*netClient),

		baseTickMs: ssrv.param.BaseTickMs,
		lg:         ssrv.lg.Child("Session " + strconv.Itoa(sessionID)),
	}
}
