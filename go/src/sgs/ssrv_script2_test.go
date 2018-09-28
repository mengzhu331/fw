package sgs

import (
	"testing"
	"time"
)

func TestScript2(t *testing.T) {
	srv, _ := makeSSrv(SSrvParam{
		Profile:        "test",
		DefaultClients: 2,
		MinimalClients: 2,
		OptimalWS:      30,
		BaseTickMs:     10,
		ABF:            buildMockApp,
	})

	rl := makeresLogger()
	p1 := makePlayer("regn", 22, rl, srv)
	p2 := makePlayer("yaya", 33, rl, srv)
	tr := true
	p1.s = script{
		{1000, scriptedJSQ, &tr},
		{2000, scriptedQSQ, &tr},
		{4000, scriptedJSQ, &tr},
	}
	p2.s = script{
		{3000, scriptedJSQ, &tr},
	}

	go rl.run()
	go p1.run(t)
	go p2.run(t)

	<-time.After(time.Duration(5) * time.Second)

	rl.mch <- "quit"

	<-time.After(time.Duration(50) * time.Millisecond)
	if !rl.cl.conformTo(commandLog{
		commandLE{4000, _CMD_INIT_APP, 22},
		commandLE{4000, _CMD_INIT_APP, 33},
		commandLE{4000, CMD_APP_RUN, 22},
		commandLE{4000, CMD_APP_RUN, 33},
	}) {
		t.Errorf("Command Log does not conform to expectation")
	}

}
