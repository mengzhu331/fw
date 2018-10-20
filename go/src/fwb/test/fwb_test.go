package test

import (
	"fwb/core"
	"log"
	"runtime"
	"sgs"
	"strings"
	"testing"
	"time"
)

func TestPlay2PvP(t *testing.T) {
	loadConf()
	mockClient1 := mockClient{
		name:     _1P_NAME,
		clientID: _1P_ID,
		t:        t,
	}
	log.Print("client 1 ", mockClient1.name)

	mockClient2 := mockClient{
		name:     _2P_NAME,
		clientID: _2P_ID,
		t:        t,
	}
	log.Print("client 2 ", mockClient2.name)

	gameOverMap := make(map[int]bool)

	go bootServer(t)
	<-time.After(time.Duration(3) * time.Second)

	go func() {
		<-time.After(time.Duration(600) * time.Second)
		log.Print("test timeout")
		t.Fail()
	}()

	for i := 0; i < 3; i++ {
		go mockClient1.connect()
		go mockClient2.connect()

		select {
		case c := <-_ch:
			gameOverMap[c] = true
			c = <-_ch
			gameOverMap[c] = true
			if len(gameOverMap) != 2 {
				log.Print("incorrect client state")
				t.Fail()
			}
		}
	}

}

func bootServer(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	c := strings.Split(file, "/")
	path := file[:len(file)-len(c[len(c)-1])-1]
	_tl.Inf("Configuration Path: " + path)
	e := sgs.Run(core.FwAppBuildFunc, path)
	if e != nil {
		t.Fail()
	}
}
