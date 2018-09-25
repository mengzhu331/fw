package authsvr

import (
	"er"
	"hlf"

	"github.com/google/uuid"
)

const (
	_REQ_LOGIN = iota + 1
	_REQ_LOGOUT
	_REQ_VALIDATE_TOKEN
	_REQ_SHUT_DOWN
)

const (
	_RES_SUCCESSFUL = iota
	_RES_FAIL
	_RES_ILLEGAL_REQUEST
)

var _em = map[int]func(*authReq) *authRes{
	_REQ_LOGIN:          onLogin,
	_REQ_LOGOUT:         onLogout,
	_REQ_VALIDATE_TOKEN: onValidateToken,
	_REQ_SHUT_DOWN:      onShutDown,
}

var _log = hlf.MakeLogger("AuthSrv")
var _reqch = make(chan *authReq)
var _clients = make(map[string]client)
var _clientId = 0x8000
var _running = false

type client struct {
	id       int
	token    string
	username string
}

type authReq struct {
	id    int
	param interface{}
	resch chan *authRes
}

type loginParam struct {
	username string
	password string
}

type logoutParam struct {
	username string
	clientId int
	token    string
}

type vtParam struct {
	clientId int
	token    string
}

type shutdownParam struct {
	info string
}

type authRes struct {
	statusCode int
	param      interface{}
}

type loginResParam struct {
	clientId int
	token    string
}

func srvRoutine(errch chan *er.Err, reqch chan *authReq) {
	_running = true

	for _running {
		select {
		case req := <-reqch:
			r := onReq(req)
			if req.resch != nil {
				req.resch <- r
			}
		}
	}
}

func onReq(rq *authReq) *authRes {
	if rq == nil || _em[rq.id] == nil {
		er.Throw(_E_INVALID_AUTH_REQUEST, er.EInfo{
			"details": "invalid auth request",
			"request": rq,
		}).To(_log)
		return &authRes{
			statusCode: _RES_ILLEGAL_REQUEST,
		}
	}

	return _em[rq.id](rq)
}

func onLogin(rq *authReq) *authRes {
	p, ok := rq.param.(loginParam)
	if !ok {
		er.Throw(_E_INVALID_AUTH_REQUEST, er.EInfo{
			"details": "invalid login request",
			"param":   rq.param,
		}).To(_log)
		return &authRes{
			statusCode: _RES_ILLEGAL_REQUEST,
		}
	}

	//TODO: auth client

	_, found := _clients[p.username]
	if !found {
		_clientId++
		token, _ := uuid.NewRandom()
		c := client{
			id:       _clientId,
			token:    token.String(),
			username: p.username,
		}
		_clients[p.username] = c
	}

	return &authRes{
		statusCode: _RES_SUCCESSFUL,
		param: loginResParam{
			clientId: _clients[p.username].id,
			token:    _clients[p.username].token,
		},
	}
}

func onLogout(rq *authReq) *authRes {
	p, ok := rq.param.(logoutParam)
	if !ok {
		er.Throw(_E_INVALID_AUTH_REQUEST, er.EInfo{
			"details": "invalid logout request",
			"param":   rq.param,
		}).To(_log)
		return &authRes{
			statusCode: _RES_ILLEGAL_REQUEST,
		}
	}

	_, found := _clients[p.username]

	if !found {
		er.Throw(_E_USER_NOT_LOGIN, er.EInfo{
			"details": "user has not login",
			"user":    p.username,
		}).To(_log)
		return &authRes{
			statusCode: _RES_FAIL,
		}
	}

	delete(_clients, p.username)

	return &authRes{
		statusCode: _RES_SUCCESSFUL,
	}
}

func onValidateToken(rq *authReq) *authRes {
	p, ok := rq.param.(vtParam)
	if !ok {
		er.Throw(_E_INVALID_AUTH_REQUEST, er.EInfo{
			"details": "invalid validate token request",
			"param":   rq.param,
		}).To(_log)
		return &authRes{
			statusCode: _RES_ILLEGAL_REQUEST,
		}
	}

	for _, c := range _clients {
		if c.id == p.clientId {
			if c.token == p.token {
				return &authRes{
					statusCode: _RES_SUCCESSFUL,
				}
			}
			break
		}
	}

	return &authRes{
		statusCode: _RES_FAIL,
	}
}

func onShutDown(rq *authReq) *authRes {
	p, ok := rq.param.(shutdownParam)
	if !ok {
		er.Throw(_E_INVALID_AUTH_REQUEST, er.EInfo{
			"details": "invalid shutdown request",
			"param":   rq.param,
		}).To(_log)
		return &authRes{
			statusCode: _RES_ILLEGAL_REQUEST,
		}
	}
	_log.Ntf("AuthSrv shut down: " + p.info)
	_running = false

	return &authRes{
		statusCode: _RES_SUCCESSFUL,
	}
}
