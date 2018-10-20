package sutil

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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

//Shuffle randomly order a list of elements
func Shuffle(list ...interface{}) []interface{} {
	for i := 0; i < len(list)-1; i++ {
		n := rand.Intn(len(list)-i) + i
		m := list[i]
		list[i] = list[n]
		list[n] = m
	}
	return list
}

//ShuffleInt randomly order a list of int
func ShuffleInt(list ...int) []int {
	interfaceList := make([]interface{}, len(list))
	for i, v := range list {
		interfaceList[i] = v
	}
	interfaceList = Shuffle(interfaceList...)
	for i, v := range interfaceList {
		list[i], _ = v.(int)
	}
	return list
}
