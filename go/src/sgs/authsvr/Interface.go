package authsvr

import (
	"er"
	"time"
)

const (
	_TIMEOUT = 5
)

//StartUp run AuthSrv
func StartUp() *er.Err {
	_log.Inf("Starting AuthServer...")

	errch := make(chan *er.Err)

	go srvRoutine(errch, _reqch)

	select {
	case e := <-errch:
		if e != nil {
			_log.Err("Failed to start AuthServer")
			return e
		}
	case <-time.After(time.Duration(10) * time.Second):
		return er.Throw(_E_AUTH_SRV_TIMEOUT, er.EInfo{
			"details": "starting server time out",
		}).To(_log)
	}

	_log.Inf("AuthServer started")
	return nil
}

//Login user login
func Login(username string, password string) (int, string, *er.Err) {
	resch := make(chan *authRes, 100)
	_reqch <- &authReq{
		id: _REQ_LOGIN,
		param: loginParam{
			username: username,
			password: password,
		},
		resch: resch,
	}
	select {
	case rs := <-resch:
		if rs.statusCode == _RES_SUCCESSFUL {
			p, ok := rs.param.(loginResParam)
			if !ok {
				return 0, "", er.Throw(_E_INVALID_AUTH_RESULT, er.EInfo{
					"details": "invalid auth result",
					"param":   rs.param,
				}).To(_log)
			}
			return p.clientId, p.token, nil
		}
	case <-time.After(time.Duration(_TIMEOUT) * time.Second):
	}

	return 0, "", er.Throw(_E_AUTH_SRV_TIMEOUT, er.EInfo{
		"details": "login timeout",
	})
}

//Logout user logout
func Logout(username string, clientId int, token string) *er.Err {
	resch := make(chan *authRes, 100)
	_reqch <- &authReq{
		id: _REQ_LOGOUT,
		param: logoutParam{
			username: username,
			clientId: clientId,
			token:    token,
		},
		resch: resch,
	}

	select {
	case rs := <-resch:
		if rs.statusCode == _RES_SUCCESSFUL {
			return nil
		} else {
			return er.Throw(_E_LOGOUT_FAIL, er.EInfo{
				"details":  "logout failed, incorrect parameter",
				"username": username,
				"clientId": clientId,
				"token":    token,
			})
		}
	case <-time.After(time.Duration(_TIMEOUT) * time.Second):
	}

	return er.Throw(_E_AUTH_SRV_TIMEOUT, er.EInfo{
		"details": "login timeout",
	})
}

//ValidateToken validate client token
func ValidateToken(clientID int, token string) (bool, *er.Err) {
	resch := make(chan *authRes, 100)
	_reqch <- &authReq{
		id: _REQ_VALIDATE_TOKEN,
		param: vtParam{
			clientId: clientID,
			token:    token,
		},
		resch: resch,
	}

	select {
	case rs := <-resch:
		if rs.statusCode == _RES_SUCCESSFUL {
			return false, nil
		} else {
			return false, er.Throw(_E_LOGOUT_FAIL, er.EInfo{
				"details": "logout failed, incorrect parameter",
			})
		}
	case <-time.After(time.Duration(_TIMEOUT) * time.Second):
	}

	return false, er.Throw(_E_AUTH_SRV_TIMEOUT, er.EInfo{
		"details": "login timeout",
	})

}
