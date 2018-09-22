package web

import (
	"encoding/json"
	"er"
	"hlf"
	"net/http"
	"sgs/ssvr"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var _param WebSrvParam

var _log = hlf.MakeLogger("SGS Web Server")

const (
	_EWeb = 0x1000
)

const (
	_E_START_WEB_SERVER_FAIL = er.IMPT_UNRECOVERABLE | er.ET_SERVICE | _EWeb | 0x1
)

//endpointMap name of endpoint
type endpointMap map[string]string

func (me *endpointMap) toString() string {
	es := "("
	first := true
	for _, v := range *me {
		if !first {
			es += ", "
		}
		es += v
		first = false
	}
	es += ")"

	return es
}

//WebSrvParam parameter for rest/ws web server
type WebSrvParam struct {
	Port        int
	EPM         endpointMap
	WSReadBuff  int
	WSWriteBuff int
}

//StartUp start up the sgs web server
func StartUp(param WebSrvParam) error {
	_log.Inf("Starting SGS Web Server...")
	router := mux.NewRouter()
	_param = param

	_log.Ntf("SGS Web Server Params: Port %v, Endpoints %v", _param.Port, _param.EPM.toString())

	router.HandleFunc(param.EPM["login"], loginRest).Methods("POST")
	router.HandleFunc(param.EPM["join_session"], joinSessionRest).Methods("POST")
	router.HandleFunc(param.EPM["ws"], connectWS)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func loginRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("LoginRest() enter")
	defer _log.Trc("LoginRest() leave")

	enableCors(&w)
	queries := r.URL.Query()
	un, _ := queries["username"]
	pw, _ := queries["password"]

	_log.Inf("Login request: Username %v", un)

	if len(un) != 1 || len(pw) != 1 {
		info := "Username or password is missing"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	cid, err := ssvr.Login(un[0], pw[0])

	if err != nil {
		info := "Incorrect username or password"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cid)

	_log.Inf("Login successful")
	return
}

func joinSessionRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("JoinSessionRest() enter")
	defer _log.Trc("JoinSessionRest() leave")

	enableCors(&w)
	queries := r.URL.Query()
	cid, _ := queries["clientid"]

	_log.Inf("Join Session request: Client ID %v", cid)

	if len(cid) != 1 {
		info := "Invalid client ID"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	cidi, err := strconv.Atoi(cid[0])

	if err != nil {
		info := "Invalid client ID"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	err = ssvr.JoinSession(cidi)

	if err != nil {
		info := err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(info))
		return
	}

	json.NewEncoder(w).Encode("Success")

	_log.Inf("Join Session successful")
}

func connectWS(w http.ResponseWriter, r *http.Request) {
	_log.Trc("connectWS() enter")
	defer _log.Trc("connectWS() leave")

	enableCors(&w)
	queries := r.URL.Query()
	cid, _ := queries["clientid"]

	_log.Inf("Connect WS request: Client ID %v", cid)

	if len(cid) != 1 {
		info := "Invalid client ID"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	cidi, err := strconv.Atoi(cid[0])

	if err != nil {
		info := "Invalid client ID"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	conn, wserr := makeWSConnection(w, r)
	if wserr != nil {
		_log.Err(wserr.Error())
		return
	}

	ssvr.BindNetConn(cidi, &wsConn{
		conn: conn,
	})

	w.Write([]byte("Success"))
	_log.Inf("Connect WS successful")
}

func makeWSConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  _param.WSReadBuff,
		WriteBufferSize: _param.WSWriteBuff,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	return upgrader.Upgrade(w, r, nil)
}
