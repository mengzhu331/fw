package sgs

import "errors"

const (
	BASE_CLIENT_ID int = 0x00010003
)

var _nextClientID = BASE_CLIENT_ID

var _clients = make(map[int]*Client)

type AppCreator func() App

var _appCreator AppCreator

func HookUpAppCreator(appCreator AppCreator) {
	_appCreator = appCreator
}

func Login(username string, password string) (int, int) {
	for _, v := range _clients {
		if v.Username == username {
			return v.ID, 1
		}
	}

	_clients[_nextClientID] = &Client{
		Username: username,
		ID:       _nextClientID,
	}
	_nextClientID++
	return _clients[_nextClientID-1].ID, 0
}

var _sessionID = 0x5000
var _sessions = make([]Session, 0)
var _csi = 0
var _cmdCh chan Command

func JoinSession(clientID int) error {
	if _clients[clientID] == nil {
		return errors.New("Invalid client ID")
	}

	if _csi >= len(_sessions) {
		s, err := initNewSession()
		if err != nil {
			return err
		}
		_sessions = append(_sessions, s)
	}

	_sessions[_csi].clients = append(_sessions[_csi].clients, _clients[clientID])
	_clients[clientID].Owner = _sessions[_csi]
	if len(_sessions[_csi].clients) == _sessions[_csi].maxClient {
		go _sessions[_csi].run()
		_csi++
	}
	return nil
}

func initNewSession() (Session, error) {
	if _appCreator == nil {
		return Session{}, errors.New("Application not found")
	}

	_sessionID++

	return Session{
		ID:        _sessionID,
		maxClient: 2,
		clients:   make([]*Client, 0),
		appConfig: AppConfig{
			TickIntervalMs: 100,
			CmdMap:         make(map[uint]CommandExecutor),
		},
		app:      _appCreator(),
		cmdChOut: _cmdCh,
	}, nil
}
