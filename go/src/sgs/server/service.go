package server

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
	return 0, nil
}

//JoinSession Join a game session
func JoinSession(clientID int) error {
	return nil
}

//BindNetClient Bind a NetClient to the client ID
func BindNetClient(clientID int, client NetClient) error {
	_clients.set(clientID, client)
	return nil
}
