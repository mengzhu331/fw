package web

import (
	"encoding/json"
	"log"
	"net/http"
	"sgs/ssvr"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	//EP_MIN minimal endpoint index
	EP_MIN = iota

	//EP_LOGIN login endpoint
	EP_LOGIN

	//EP_WEBSOCKET open websocket endpoint
	EP_WEBSOCKET

	//EP_JOINSESSION join session endpoint
	EP_JOINSESSION

	//EP_MAX maximal endpoint index
	EP_MAX
)

//EndPointMap name of endpoint
type EndPointMap map[int]string

//WebSrvParam parameter for rest/ws web server
type WebSrvParam struct {
	Port        int
	EPM         EndPointMap
	WSReadBuff  int
	WSWriteBuff int
}

var _param WebSrvParam

//StartUp start up the sgs web server
func StartUp(param WebSrvParam) error {
	router := mux.NewRouter()
	_param = param
	router.HandleFunc(param.EPM[EP_LOGIN], loginRest).Methods("POST")
	router.HandleFunc(param.EPM[EP_JOINSESSION], joinSessionRest).Methods("POST")
	router.HandleFunc(param.EPM[EP_WEBSOCKET], connectWS)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(param.Port), router))
	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func loginRest(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	queries := r.URL.Query()
	un, _ := queries["username"]
	pw, _ := queries["password"]

	if len(un) != 1 || len(pw) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username or password is missing"))
		return
	}

	cid, err := ssvr.Login(un[0], pw[0])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incorrect username or password"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cid)
	return
}

func joinSessionRest(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	queries := r.URL.Query()
	cid, _ := queries["clientid"]

	if len(cid) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid client ID"))
		return
	}

	cidi, err := strconv.Atoi(cid[0])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid client ID"))
		return
	}

	err = ssvr.JoinSession(cidi)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode("Success")
}

func connectWS(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	queries := r.URL.Query()
	cid, _ := queries["clientid"]

	if len(cid) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid client ID"))
		return
	}

	cidi, err := strconv.Atoi(cid[0])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid client ID"))
		return
	}

	conn, wserr := makeWSConnection(w, r)
	if wserr != nil {
		return
	}

	ssvr.BindNetConn(cidi, &wsConn{
		conn: conn,
	})

	w.Write([]byte("Success"))
}

func makeWSConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  _param.WSReadBuff,
		WriteBufferSize: _param.WSWriteBuff,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	return upgrader.Upgrade(w, r, nil)
}
