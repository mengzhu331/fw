package sgs

import (
	"er"
	"hlf"
	"strconv"
	"sync"
)

var _serverID = 0x6000
var _serverIDMutex = sync.Mutex{}

type clientMap map[int]int

type sessionServer struct {
	id    int
	param *SSrvParam

	lg hlf.Logger

	clients clientMap

	clientsMutex sync.RWMutex

	currentSession *session

	currentSessionMutex sync.Mutex

	sessions map[int]*session

	sessionsMutex sync.RWMutex

	sessionID int
}

//SSrvParam parameters for the server
type SSrvParam struct {
	Profile        string
	DefaultClients int
	MinimalClients int
	OptimalWS      int
	BaseTickMs     int
	ABF            AppBuildFunc
	WSReadBuff     int
	WSWriteBuff    int
}

func makeSSrv(param SSrvParam) (*sessionServer, error) {
	_log.Inf("Init SSVR...")

	if !validate(&param) {
		return nil, er.Throw(_E_INVALID_SERVER_PARAM, er.EInfo{
			"details": "invalid server parameters used for SSVR",
			"param":   param,
		}).To(_log)
	}

	server := &sessionServer{}

	_serverIDMutex.Lock()
	_serverID++
	server.id = _serverID
	_serverIDMutex.Unlock()

	server.param = &param

	server.sessionID = 0x2000

	server.lg = hlf.MakeLogger("SSRV" + strconv.Itoa(server.id))

	server.clients = make(clientMap)

	server.clientsMutex = sync.RWMutex{}

	server.currentSessionMutex = sync.Mutex{}

	server.currentSession = nil

	server.sessions = make(map[int]*session)

	server.sessionsMutex = sync.RWMutex{}

	_log.Inf("Init SSVR successful")
	return server, nil
}

func (me *sessionServer) joinSessionQueue(username string, clientID int, conn NetConn) *er.Err {
	me.lg.Inf("Client %v requested to join session queue", clientID)

	me.currentSessionMutex.Lock()
	defer me.currentSessionMutex.Unlock()

	if me.currentSession == nil {
		me.sessionID++
		me.currentSession = makeSession(me.sessionID, me)

		me.lg.Inf("Created new session %v", me.sessionID)
	}

	c, found := me.currentSession.clients[clientID]

	if found {
		me.lg.Ntf("already added to session %v, duplicated request", c)
		return nil
	}

	me.clientsMutex.RLock()
	var s int
	s, found = me.clients[clientID]
	me.clientsMutex.RUnlock()

	if found {
		me.sessionsMutex.Lock()
		sess, founds := me.sessions[s]
		if founds {
			founds = !sess.closed
		}
		me.sessionsMutex.Unlock()

		if founds {
			return er.Throw(_E_CLIENT_ALREADY_JOIN_SESSION, er.EInfo{
				"details": "client duplicated join session request",
				"client":  strconv.Itoa(clientID),
			})
		}
	}

	me.currentSession.clients[clientID] = &netClient{
		id:       clientID,
		username: username,
		s:        nil,
		conn:     conn,
		mch:      make(chan Command),
	}

	if len(me.currentSession.clients) == me.param.DefaultClients {

		me.lg.Inf("session %v has sufficient user joined, is to be started", me.currentSession.id)
		me.sessionsMutex.Lock()
		me.sessions[me.currentSession.id] = me.currentSession
		me.sessionsMutex.Unlock()

		me.clientsMutex.Lock()
		for _, c := range me.currentSession.clients {
			me.clients[c.id] = me.currentSession.id
		}
		me.clientsMutex.Unlock()

		if !me.startSession(me.currentSession) {
			me.closeSession(me.currentSession)
		}
		me.currentSession = nil
	}
	return nil
}

func (me *sessionServer) quitSessionQueue(clientID int) *er.Err {
	me.lg.Inf("Client %v requested to quit session queue", clientID)

	me.currentSessionMutex.Lock()
	defer me.currentSessionMutex.Unlock()

	if me.currentSession == nil || me.currentSession.clients[clientID] == nil {
		return er.Throw(_E_QUIT_SESSION_QUEUE_INVALID, er.EInfo{
			"details": "quit request from client not in session queue",
			"client":  clientID,
		}).To(me.lg)
	}

	delete(me.currentSession.clients, clientID)
	return nil
}

func (me *sessionServer) reconnectClient(clientID int, conn NetConn) *er.Err {
	me.lg.Inf("Client %v reconnect", clientID)

	me.currentSessionMutex.Lock()

	if me.currentSession != nil && me.currentSession.clients[clientID] != nil {
		c := me.currentSession.clients[clientID]
		c.conn = conn
		me.currentSession.clients[clientID] = c
		me.currentSessionMutex.Unlock()
		return nil
	}
	me.currentSessionMutex.Unlock()

	me.clientsMutex.Lock()

	sid, found := me.clients[clientID]
	me.clientsMutex.Unlock()

	if !found {
		return er.Throw(_E_CLIENT_FAILE_RECONNECT, er.EInfo{
			"details": "reconnect client not in session",
			"client":  clientID,
		}).To(me.lg)
	}

	me.sessionsMutex.Lock()
	s := me.sessions[sid]
	me.sessionsMutex.Unlock()

	if s == nil {
		return er.Throw(_E_INVALID_SESSION_ID, er.EInfo{
			"details": "invalid session id",
		}).To(me.lg)
	}

	s.mch <- Command{
		ID:      _CMD_CLIENT_RECONNECT,
		Source:  clientID,
		Payload: conn,
	}
	return nil
}

func (me *sessionServer) startSession(s *session) bool {
	err := s.run(me.param.ABF, me.param.Profile)
	fail := (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE
	return !fail
}

func (me *sessionServer) closeSession(s *session) {
	me.lg.Inf("Closing session %v...", s.id)

	me.clientsMutex.Lock()
	for _, c := range s.clients {
		c.close()
		delete(me.clients, c.id)
	}
	me.clientsMutex.Unlock()

	me.lg.Inf("Session %v closed", s.id)
}

func validate(p *SSrvParam) bool {
	return p.ABF != nil && p.BaseTickMs > 0 && p.DefaultClients > 0 && p.MinimalClients > 0
}
