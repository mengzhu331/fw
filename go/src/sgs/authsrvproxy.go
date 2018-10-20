package sgs

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type authSrvPrx interface {
	setServerURI(string)
	vclient(clientID int, token string) bool
	enableTestClients(bool)
}

type sgasPrx struct {
	serverURI   string
	testEnabled bool
}

func (me *sgasPrx) enableTestClients(enabled bool) {
	me.testEnabled = enabled
}

func (me *sgasPrx) vclient(clientID int, token string) bool {

	//clientIDs started with 0x40000000 are test client IDs
	if me.testEnabled && clientID > 0x40000000 {
		return true
	}
	requestURI := me.serverURI + "/vclient" + "?client=" + strconv.Itoa(clientID) + "&token=" + token
	client := http.Client{}
	request, _ := http.NewRequest("POST", requestURI, nil)
	response, err := client.Do(request)
	if err != nil {
		_log.Ntf("Invalid clientID or token %v %v", clientID, token)
		return false
	}

	registered, _ := ioutil.ReadAll(response.Body)
	return string(registered) != ""
}

func (me *sgasPrx) setServerURI(serverURI string) {
	me.serverURI = serverURI
}

func createSgasPrx() *sgasPrx {
	return &sgasPrx{}
}
