package sgs

import (
	"er"
	"net/http"
	"strconv"
	"sutil"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	_EWeb = 0x1000
)

const (
	_E_START_WEB_SERVER_FAIL = er.IMPT_UNRECOVERABLE | er.ET_SERVICE | _EWeb | 0x1
)

var _auths authSrvPrx

//WebSrvParam parameter for rest/ws web server
type WebSrvParam struct {
	Port        int
	WSReadBuff  int
	WSWriteBuff int
}

func webStartUp(param WebSrvParam, auths authSrvPrx) error {
	_log.Inf("Starting SGS Web Server...")
	_auths = auths

	router := mux.NewRouter()

	_log.Ntf("SGS Web Server Params: Port %v", param.Port)

	router.HandleFunc("/join_session", joinSessionRest)
	router.HandleFunc("/quit_session", quitSessionRest).Methods("POST")
	router.HandleFunc("/reconnect", reconnectClientRest)

	_ch := make(chan string)

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(param.Port), router)
		if err != nil {
			_ch <- err.Error()
		}
	}()

	select {
	case es := <-_ch:
		{
			e := er.Throw(_E_START_WEB_SERVER_FAIL, er.EInfo{"Fail info": es})
			e.To(_log)
			return e
		}
	case <-time.After(time.Duration(1) * time.Second):
		{
			_log.Inf("SGS Web Server started")
		}
	}
	return nil
}

func joinSessionRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("joinSessionRest() enter")
	defer _log.Trc("joinSessionRest() leave")

	var info string
	var conn *websocket.Conn
	var err error
	var icid int

	sutil.EnableCors(&w)

	queries := r.URL.Query()
	cid, _ := queries["clientid"]
	token, _ := queries["token"]
	username, _ := queries["username"]

	_log.Inf("Join Session WS request: client ID %v, username %v", cid, username)

	if len(cid) != 1 || len(token) != 1 || len(username) != 1 {
		goto __invalidparameter
	}

	icid, err = strconv.Atoi(cid[0])

	if err != nil {
		goto __invalidparameter
	}

	if username[0] != _auths.vclient(icid, token[0]) {
		goto __failtovalidateclient
	}

	conn, err = makeWSConnection(w, r)
	if err != nil {
		_log.Err(err.Error())
		goto __failtoconnectws
	}

	err = _srv.joinSessionQueue(username[0], icid, &wsConn{
		clientId: icid,
		conn:     conn})

	if err != nil {
		_log.Err("Failed to join session")
		goto __failtojoinsession
	}

	w.Write([]byte("Success"))
	_log.Inf("Connect WS successful")
	return

__invalidparameter:
	info = "Invalid parameter or mandatory parameter missing"
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtoconnectws:
	info = "Fail to create WS connection"
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtovalidateclient:
	info = "Fail to validate client"
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtojoinsession:
	info = err.Error()
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(info))
	_log.Ntf(info)
	return
}

func quitSessionRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("quitSessionRest() enter")
	defer _log.Trc("quitSessionRest() leave")

	var info string
	var err error
	var icid int

	sutil.EnableCors(&w)

	queries := r.URL.Query()
	cid, _ := queries["clientid"]
	token, _ := queries["token"]

	_log.Inf("Quit Session WS request: client ID %v", cid)

	if len(cid) != 1 || len(token) != 1 {
		goto __invalidparameter
	}

	icid, err = strconv.Atoi(cid[0])

	if err != nil {
		goto __invalidparameter
	}

	if _auths.vclient(icid, token[0]) == "" {
		goto __failtovalidateclient
	}

	err = _srv.quitSessionQueue(icid)

	if err != nil {
		_log.Err("Failed to quit session")
		goto __failtoquitsession
	}

	w.Write([]byte("Success"))
	_log.Inf("Quit session successful")
	return

__invalidparameter:
	info = "Invalid parameter or mandatory parameter missing"
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtoquitsession:
	info = err.Error()
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtovalidateclient:
	info = err.Error()
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

}

func reconnectClientRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("reconnectClientRest() enter")
	defer _log.Trc("reconnectClientRest() leave")

	var info string
	var conn *websocket.Conn
	var err error
	var icid int

	sutil.EnableCors(&w)

	queries := r.URL.Query()
	cid, _ := queries["clientid"]
	token, _ := queries["token"]

	_log.Inf("Reconnect session WS request: client ID %v", cid)

	if len(cid) != 1 || len(token) != 1 {
		goto __invalidparameter
	}

	icid, err = strconv.Atoi(cid[0])

	if err != nil {
		goto __invalidparameter
	}

	if "" == _auths.vclient(icid, token[0]) {
		goto __failtovalidateclient
	}

	conn, err = makeWSConnection(w, r)
	if err != nil {
		_log.Err(err.Error())
		goto __failtoconnectws
	}

	err = _srv.reconnectClient(icid, &wsConn{
		conn: conn,
	})

	if err != nil {
		_log.Err("Failed to reconnect client")
		goto __failtoreconnectclient
	}

	w.Write([]byte("Success"))
	_log.Inf("Connect WS successful")
	return

__invalidparameter:
	info = "Invalid parameter or mandatory parameter missing"
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtoconnectws:
	info = "Fail to create WS connection"
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtoreconnectclient:
	info = err.Error()
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

__failtovalidateclient:
	info = "Failed to validate client token"
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(info))
	_log.Ntf(info)
	return

}

func makeWSConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  _srv.param.WSReadBuff,
		WriteBufferSize: _srv.param.WSWriteBuff,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	return upgrader.Upgrade(w, r, nil)
}
