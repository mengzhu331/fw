package ssvr

import (
	"errors"
	"hlf"
)

//SSrvParam parameters for the server
type SSrvParam struct {
	CPS        int
	BaseTickMs int
	ABF        AppBuildFunc
}

var _param SSrvParam

var _log = hlf.CreateLogger("SSVR")

//Init set server param
func Init(param SSrvParam) error {
	_log.Inf("Starting SSVR...")
	_param = param
	_log.Inf("SSVR started")
	return nil
}

//Login log into system with user credential
func Login(username string, password string) (int, error) {
	_cMutex.Lock()
	defer _cMutex.Unlock()

	for id, c := range _clients {
		if c.username == username {
			return id, nil
		}
	}

	_clientID++
	_clients[_clientID] = netClient{
		username: username,
	}

	return _clientID, nil
}

//JoinSession Join a game session
func JoinSession(clientID int) error {
	_csMutex.Lock()
	defer _csMutex.Unlock()

	if _currentSession == nil {
		_sessionID++
		_currentSession = &session{
			ID:      _sessionID,
			cch:     make(chan string),
			clients: make(clientMap),
		}
	}

	_, found := _currentSession.clients[clientID]

	if found {
		return errors.New("Already added to session")
	}

	_cMutex.Lock()
	c, cfound := _clients[clientID]
	_cMutex.Unlock()

	if !cfound {
		return errors.New("Client not found")
	}

	_currentSession.clients[clientID] = c

	if len(_currentSession.clients) == _param.CPS {

		_sMutex.Lock()
		_sessions[_currentSession.ID] = _currentSession
		_sMutex.Unlock()

		go _currentSession.run(_param.BaseTickMs)
		_currentSession = nil
	} else if len(_currentSession.clients) > _param.CPS {
		return errors.New("client number exceeds max allowed for session")
	}
	return nil
}

//BindNetConn Bind a NetConn to the client ID
func BindNetConn(clientID int, net NetConn) error {
	_cMutex.Lock()
	defer _cMutex.Unlock()

	client, ok := _clients[clientID]
	if !ok {
		return errors.New("client not found")
	}
	client.conn = net
	_clients[clientID] = client
	return nil
}
