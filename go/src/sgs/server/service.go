package server

import "errors"

//ServerParam parameters for the server
type ServerParam struct {
	CPS int
}

var _param ServerParam

//Init set server param
func Init(param ServerParam) {
	_param = param
}

//Login log into system with user credential
func Login(username string, password string) (int, error) {
	if c, found := _clients.find(username); found {
		return c.id, nil
	}
	c := _clients.add(netClient{
		username: username,
	})
	return c, nil
}

//JoinSession Join a game session
func JoinSession(clientID int) error {
	return nil
}

//BindNetConn Bind a NetConn to the client ID
func BindNetConn(clientID int, net NetConn) error {
	client, ok := _clients.get(clientID)
	if !ok {
		return errors.New("client not found")
	}
	client.conn = net
	_clients.set(clientID, client)
	return nil
}
