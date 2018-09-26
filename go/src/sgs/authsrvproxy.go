package sgs

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type authSrvPrx interface {
	setServerURI(string)
	vclient(clientID int, token string) string
}

type sgasPrx struct {
	serverURI string
}

func (me *sgasPrx) vclient(clientID int, token string) string {
	requestURI := me.serverURI + "/vclient" + "?client=" + strconv.Itoa(clientID) + "&token=" + token
	client := http.Client{}
	request, _ := http.NewRequest("POST", requestURI, nil)
	response, err := client.Do(request)
	if err != nil {
		_log.Ntf("Invalid clientID or token %v %v", clientID, token)
		return ""
	}

	username, _ := ioutil.ReadAll(response.Body)
	return string(username)
}

func (me *sgasPrx) setServerURI(serverURI string) {
	me.serverURI = serverURI
}

func createSgasPrx() *sgasPrx {
	return &sgasPrx{}
}
