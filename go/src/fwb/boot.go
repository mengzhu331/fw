package main

import (
	"log"
	"sgs"
	"time"
)

const (
	WAIT_SECOND time.Duration = 10
)

type Exec func(i int)

type D struct{}

func (this *D) d(i int) {
	log.Println("D.d", i)
}

func main() {
	log.Println("FW Backend 1.0 started")
	sgs.HookUpAppCreator(fwGameCreator)
	cli1, _ := sgs.Login("regn", "regn")
	cli2, _ := sgs.Login("yaya", "regn")
	sgs.JoinSession(cli1)
	sgs.JoinSession(cli2)
	<-time.After(WAIT_SECOND * time.Second)
	log.Println("FW Backend 1.0 ended")
}
