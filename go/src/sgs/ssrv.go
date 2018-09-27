package sgs

import (
	"er"
	"strconv"
	"sync"
)

var _param SSrvParam

var _slog = _log.Child("Sessions")

type clientMap map[int]int

var _clients = make(clientMap)

var _cMutex = sync.RWMutex{}

var _csMutex sync.Mutex

var _currentSession *session

var _sessions = make(map[int]*session)

var _sMutex sync.Mutex

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

func initSSrv(param SSrvParam) error {
	_log.Inf("Starting SSVR...")
	if !validate(&param) {
		return er.Throw(_E_INVALID_SERVER_PARAM, er.EInfo{
			"details": "invalid server parameters used for SSVR",
			"param":   param,
		}).To(_log)
	}
	_param = param
	_log.Inf("SSVR started")
	return nil
}

func joinSessionQueue(username string, clientID int, conn NetConn) *er.Err {
	_log.Inf("Client %v requested to join session queue", clientID)

	_csMutex.Lock()
	defer _csMutex.Unlock()

	if _currentSession == nil {
		_sessionID++
		_currentSession = makeSession(_sessionID)

		_log.Inf("Created new session %v", _sessionID)
	}

	c, found := _currentSession.clients[clientID]

	if found {
		_log.Ntf("already added to session %v, duplicated request", c)
		return nil
	}

	_cMutex.Lock()
	var s int
	s, found = _clients[clientID]
	_cMutex.Unlock()

	if found {

		_sMutex.Lock()
		sess, founds := _sessions[s]
		if founds {
			founds = !sess.closed
		}
		_sMutex.Unlock()

		if founds {
			return er.Throw(_E_CLIENT_ALREADY_JOIN_SESSION, er.EInfo{
				"details": "client duplicated join session request",
				"client":  strconv.Itoa(clientID),
			})
		}
	}

	_currentSession.clients[clientID] = netClient{
		id:       clientID,
		username: username,
		s:        nil,
		conn:     conn,
		mch:      make(chan Command),
	}

	if len(_currentSession.clients) == _param.DefaultClients {

		_log.Inf("session %v has sufficient user joined, is to be started", _currentSession.id)
		_sMutex.Lock()
		_sessions[_currentSession.id] = _currentSession
		_sMutex.Unlock()
		if !startSession(_currentSession) {
			closeSession(_currentSession)
		}
		_currentSession = nil

	}
	return nil
}

func quitSessionQueue(clientID int) *er.Err {
	_log.Inf("Client %v requested to quit session queue", clientID)
	_csMutex.Lock()
	defer _csMutex.Unlock()

	c, found := _currentSession.clients[clientID]

	if !found {
		return er.Throw(_E_QUIT_SESSION_QUEUE_INVALID, er.EInfo{
			"details": "quit request from client not in session queue",
			"client":  clientID,
		}).To(_log)
	}

	c.close()

	delete(_currentSession.clients, clientID)
	return nil
}

func reconnectClient(clientID int, conn NetConn) *er.Err {
	_log.Inf("Client %v reconnect")
	_csMutex.Lock()
	c, found := _currentSession.clients[clientID]
	if found {
		c.conn = conn
		_currentSession.clients[clientID] = c
		_csMutex.Unlock()
		return nil
	}
	_csMutex.Unlock()

	_cMutex.Lock()
	defer _cMutex.Unlock()

	_, found = _clients[clientID]
	if !found {
		return er.Throw(_E_CLIENT_FAILE_RECONNECT, er.EInfo{
			"details": "reconnect client not in session",
			"client":  clientID,
		}).To(_log)
	}

	c.s.mch <- Command{
		ID:      _CMD_CLIENT_RECONNECT,
		Payload: conn,
	}
	return nil
}

func startSession(s *session) bool {
	err := s.run()
	fail := (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE
	return !fail
}

func closeSession(s *session) {
	_log.Inf("Closing session %v...", s.id)

	_cMutex.Lock()
	for _, c := range s.clients {
		c.close()
		delete(_clients, c.id)
	}
	_cMutex.Unlock()

	_log.Inf("Session %v closed", s.id)
}

func validate(p *SSrvParam) bool {
	return p.ABF != nil && p.BaseTickMs > 0 && p.DefaultClients > 0 && p.MinimalClients > 0
}
