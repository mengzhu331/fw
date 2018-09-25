package ssvr

import (
	"er"
	"hlf"
	"strconv"
)

var _param SSrvParam

var _log = hlf.MakeLogger("SSVR")

var _slog = _log.Child("Sessions")

//SSrvParam parameters for the server
type SSrvParam struct {
	Profile        string
	DefaultClients int
	MinimalClients int
	OptimalWS      int
	BaseTickMs     int
	ABF            AppBuildFunc
}

//Init set server param
func Init(param SSrvParam) error {
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

//Login log into system with user credential
func Login(username string, password string) (int, error) {

	_log.Inf("User login: %v", username)

	_cMutex.Lock()
	defer _cMutex.Unlock()

	for id, c := range _clients {
		if c.username == username {
			_log.Ntf("User already login: %v, client %v", username, id)
			return id, nil
		}
	}

	_clientID++
	_clients[_clientID] = netClient{
		username: username,
	}

	_log.Inf("User %v login successful, assigned %v", username, _clientID)

	return _clientID, nil
}

//JoinSession Join a game session
func JoinSession(clientID int) *er.Err {
	_log.Inf("Client %v requested to join session", clientID)

	_csMutex.Lock()
	defer _csMutex.Unlock()

	if _currentSession == nil {
		_sessionID++
		_currentSession = makeSession(_sessionID)

		_log.Inf("Created new session %v", _sessionID)
	}

	c, found := _currentSession.clients[clientID]

	if found {
		_log.Ntf("already added to session %v, duplicated request", c.id)
		return nil
	}

	_cMutex.Lock()
	c, found = _clients[clientID]
	if c.s != nil {
		cid := c.id
		_cMutex.Unlock()
		return er.Throw(_E_CLIENT_ALREADY_JOIN_SESSION, er.EInfo{
			"details": "client duplicated join session request",
			"client":  strconv.Itoa(cid),
		})
	}
	c.s = _currentSession
	_cMutex.Unlock()

	if !found {
		return er.Throw(_E_JOIN_SESSION_INVALID_CLIENT, er.EInfo{
			"details": "joing session request with illegal client",
			"client":  clientID,
		}).To(_log)
	}

	_currentSession.clients[clientID] = c

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

func startSession(s *session) bool {

	err := s.run()
	fail := (err.Code() & er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE
	return !fail
}

func closeSession(s *session) {
	_log.Inf("Closing session %v...", s.id)

	_cMutex.Lock()
	for _, c := range s.clients {
		c.s = nil
	}
	_cMutex.Unlock()

	_sMutex.Lock()
	delete(_sessions, s.id)
	_sMutex.Unlock()

	_log.Inf("Session %v closed", s.id)
}

//BindNetConn Bind a NetConn to the client ID
func BindNetConn(clientID int, net NetConn) error {
	_cMutex.Lock()
	defer _cMutex.Unlock()

	client, ok := _clients[clientID]

	if !ok {
		return er.Throw(_E_BIND_CONN_INVALID_CLIENT, er.EInfo{
			"details": "bind connection with illegal client",
			"client":  clientID,
		}).To(_log)
	}

	client.conn = net
	net.BindClientID(clientID)
	_clients[clientID] = client
	return nil
}

func validate(p *SSrvParam) bool {
	return p.ABF != nil && p.BaseTickMs > 0 && p.DefaultClients > 0 && p.MinimalClients > 0
}
