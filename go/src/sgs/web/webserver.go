package web

import (
	"encoding/json"
	"log"
	"net/http"
	"sgs/aop"
	"sgs/server"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var _conf *Conf

//StartUp start up the sgs web server
func StartUp(conf *Conf) error {
	_conf = conf
	server.Init(server.ServerParam{
		CPS: _conf.cps(),
	})
	router := mux.NewRouter()
	router.HandleFunc(conf.epmv(EP_LOGIN), loginRest).Methods("POST")
	router.HandleFunc(conf.epmv(EP_JOINSESSION), joinSessionRest).Methods("POST")
	router.HandleFunc(conf.epmv(EP_WEBSOCKET), connectWS)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.port()), router))
	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func loginRest(w http.ResponseWriter, r *http.Request) {
	aop.Ib()
	defer aop.Ie()

	enableCors(&w)
	queries := r.URL.Query()
	un, _ := queries["username"]
	pw, _ := queries["password"]

	if len(un) != 1 || len(pw) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username or password is missing"))
		return
	}

	cid, err := server.Login(un[0], pw[0])

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
	aop.Ib()
	defer aop.Ie()

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

	err = server.JoinSession(cidi)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode("Success")
}

func connectWS(w http.ResponseWriter, r *http.Request) {
	aop.Ib()
	defer aop.Ie()

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

	server.BindNetClient(cidi, &wsClient{
		conn: conn,
	})

	w.Write([]byte("Success"))
}

func makeWSConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  _conf.wsReadBuff(),
		WriteBufferSize: _conf.wsWriteBuff(),
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	return upgrader.Upgrade(w, r, nil)
}
