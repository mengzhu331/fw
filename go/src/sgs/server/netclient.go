package server

//NetConn network interface independent of protocol
type NetConn interface {
	Send(msg string) error
	BindChan(ch chan string)
}

type netClient struct {
	id       int
	username string
	conn     NetConn
}

type clientMap map[int]netClient

var _clients clientMap
var _clientID int = 0x8000

func (me *clientMap) add(client netClient) int {
	_clientID++
	client.id = _clientID
	(*me)[_clientID] = client
	return client.id
}

func (me *clientMap) set(clientID int, client netClient) {
	(*me)[clientID] = client
}

func (me *clientMap) get(clientID int) (netClient, bool) {
	return (*me).get(clientID)
}

func (me *clientMap) find(username string) (netClient, bool) {
	for _, c := range *me {
		if c.username == username {
			return c, true
		}
	}
	return netClient{}, false
}
