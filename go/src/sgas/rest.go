package main

import (
	"encoding/json"
	"er"
	"net/http"
	"strconv"
	"sutil"

	"github.com/gorilla/mux"
)

type conf struct {
	Port int
}

func main() {
	_log.Inf("Init SGAS")
	router := mux.NewRouter()

	c := conf{}

	if e := loadConf(&c); e != nil && (e.Code()&er.E_IMPORTANCE) >= er.IMPT_UNRECOVERABLE {
		_log.Inf("Failed to start SGAS")
		return
	}

	epLogin := "/login"
	epLogout := "/logout"
	epVclient := "/vclient"
	router.HandleFunc(epLogin, loginRest).Methods("POST")
	router.HandleFunc(epLogout, logoutRest).Methods("POST")
	router.HandleFunc(epVclient, vclientRest).Methods("POST")

	_log.Ntf("SGAS start listening Port %v", c.Port)
	_log.Ntf("SGAS endpoints: %v %v %v", epLogin, epLogout, epVclient)

	err := http.ListenAndServe(":"+strconv.Itoa(c.Port), router)

	if err != nil {
		er.Throw(_E_START_WEB_SERVER, er.EInfo{
			"details":   "failed to start web",
			"fail info": err}).To(_log)
	}

}

func loginRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("LoginRest() enter")
	defer _log.Trc("LoginRest() leave")

	sutil.EnableCors(&w)

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

	cid, token, err := Login(un[0], pw[0])

	if err != nil {
		info := "Incorrect username or password"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(strconv.Itoa(cid) + " " + token)

	_log.Inf("Login successful, cid %v token %v", cid, token)
	return
}

func logoutRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("LogoutRest() enter")
	defer _log.Trc("LogoutRest() leave")

	sutil.EnableCors(&w)

	queries := r.URL.Query()
	un, _ := queries["username"]
	cids, _ := queries["client"]
	token, _ := queries["token"]

	_log.Inf("Login request: Username %v, client %v, token %v", un, cids, token)

	if len(un) != 1 || len(cids) != 1 || len(token) != 1 {
		info := "Username, client ID or token is missing"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	cid, err := strconv.Atoi(cids[0])

	if err != nil {
		info := "client ID is not a valid integer"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	e := Logout(un[0], cid, token[0])

	if e != nil {
		info := "Incorrect client information"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))

	_log.Inf("Logout successful")
	return
}

func vclientRest(w http.ResponseWriter, r *http.Request) {
	_log.Trc("vclientRest() enter")
	defer _log.Trc("vclientRest() leave")

	sutil.EnableCors(&w)

	queries := r.URL.Query()
	cids, _ := queries["client"]
	token, _ := queries["token"]

	_log.Inf("Validate client request: client %v, token %v", cids, token)

	if len(cids) != 1 || len(token) != 1 {
		info := "Client ID or token is missing"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	cid, err := strconv.Atoi(cids[0])

	if err != nil {
		info := "Client ID is not a valid integer"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(info))
		_log.Ntf(info)
		return
	}

	username, valid := ValidateToken(cid, token[0])

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(username))

	_log.Inf("Validate client ID finish, validity %v", valid)
	return
}

func loadConf(c *conf) *er.Err {
	cfgf := "sgas.conf"
	_log.Inf("Load settings from %v", cfgf)

	if e := sutil.LoadConfFile(cfgf, c); e != nil {
		return er.Throw(_E_LOAD_SETTINGS_FAIL, er.EInfo{
			"details": "fail to load settings",
			"info":    e,
		}).To(_log)
	}
	return nil
}
