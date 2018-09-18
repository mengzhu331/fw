package server

//NetClient client interface independent of protocol
type NetClient interface {
	Send(msg string) error
	BindChan(ch chan string)
}

type clientMap map[int]NetClient

var _clients clientMap

func (me *clientMap) set(clientID int, client NetClient) {
	(*me)[clientID] = client
}

func (me *clientMap) get(clientID int) (NetClient, bool) {
	return (*me).get(clientID)
}
