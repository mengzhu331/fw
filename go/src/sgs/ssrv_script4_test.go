package sgs

import (
	"testing"
	"time"
)

func TestScript4(t *testing.T) {
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
	fls := false
	p1.s = script{
		{1000, scriptedJSQ, &tr},
		{6000, scriptedRC, &tr},
	}
	p2.s = script{
		{2000, scriptedRC, &fls},
		{3000, scriptedJSQ, &tr},
	}

	lg := _log.Child("script4")

	go rl.run()
	go p1.run(t)
	go p2.run(t)

	<-time.After(time.Duration(10) * time.Second)

	rl.mch <- "quit"

	<-time.After(time.Duration(50) * time.Millisecond)
	expected := commandLog{
		commandLE{3000, _CMD_INIT_APP, 22},
		commandLE{3000, _CMD_INIT_APP, 33},
		commandLE{3000, CMD_APP_RUN, 22},
		commandLE{3000, CMD_APP_RUN, 33},
		commandLE{6000, _CMD_CLIENT_RECONNECT, 33},
		commandLE{6000, _CMD_CLIENT_RECONNECT, 22},
	}
	if !rl.cl.conformTo(expected) {
		lg.Ntf(rl.cl.String())
		lg.Ntf(expected.String())
		t.Errorf("Command Log does not conform to expectation")
	}

}
