package sutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//EnableCors enable web CORS
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//LoadConfFile load settings from storage
func LoadConfFile(f string, conf interface{}) error {
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(c, conf)
	return err
}
