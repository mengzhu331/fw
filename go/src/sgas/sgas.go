package main

import (
	"er"
	"hlf"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type client struct {
	id       int
	token    string
	username string
}

var _clients = make(map[string]client)
var _clientID = 0x8000
var _cMtx = sync.RWMutex{}

var _log = hlf.MakeLogger("SGAS")

//Login user login
func Login(username string, password string) (int, string, *er.Err) {
	//TODO: auth client

	_cMtx.Lock()
	defer _cMtx.Unlock()

	_, found := _clients[username]
	if !found {
		_clientID++
		c := client{
			id:       _clientID,
			token:    makeToken(),
			username: username,
		}
		_clients[username] = c
	}

	return _clients[username].id, _clients[username].token, nil
}

//Logout user request to logout
func Logout(username string, clientID int, token string) *er.Err {

	_cMtx.Lock()
	defer _cMtx.Unlock()

	c, found := _clients[username]
	if !found {
		return er.Throw(_E_USER_NOT_LOGIN, er.EInfo{
			"details": "user has not login",
			"user":    username,
		}).To(_log)
	}

	if c.id != clientID || c.token != token {
		return er.Throw(_E_INVALID_CLIENT_INFO, er.EInfo{
			"details": "invalid client info",
			"client":  clientID,
			"token":   token,
		}).To(_log)
	}

	delete(_clients, username)
	return nil
}

//ValidateToken validate a client token
func ValidateToken(clientID int, token string) bool {
	_cMtx.RLock()
	defer _cMtx.RUnlock()

	for _, c := range _clients {
		if c.id == clientID {
			if c.token == token {
				return true
			}
			return false
		}
	}
	return false
}

func makeToken() string {
	token, e := uuid.NewRandom()
	if e != nil {
		_log.Err("Backup token used, failed to generate UUID token: %v", e.Error())
		return strconv.Itoa(int(rand.Uint32()))
	}
	return strings.Replace(token.String(), "-", "", -1)
}
