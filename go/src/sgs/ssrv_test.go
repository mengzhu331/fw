package sgs

import (
	"er"
	"math"
	"os"
	"testing"
	"time"
)

type mockApp struct {
	s   Session
	cid []int
}

const (
	_CMD_INIT_APP = iota + 100
)

func (me *mockApp) Init(s Session, clientIDs []int, profile string) *er.Err {
	me.cid = clientIDs
	me.s = s
	for _, cid := range me.cid {
		s.ForwardToClient(cid, Command{
			ID: _CMD_INIT_APP,
		})
	}
	return nil
}

func (me *mockApp) SendCommand(command Command) *er.Err {
	if command.ID == CMD_TICK {
		return nil
	}

	for _, cid := range me.cid {
		command.Source = cid
		me.s.ForwardToClient(cid, command)
	}
	return nil
}

func buildMockApp() App {
	return &mockApp{}
}

type mockConn struct {
	clientID int
	ch       chan Command
}

func (me *mockConn) Send(cmd Command) error {
	cmd.Source = me.clientID
	me.ch <- cmd
	return nil
}

func (me *mockConn) Run(ch chan Command, mch chan Command) {
}

func (me *mockConn) BindClientID(clientID int) {}

type scriptedRequest struct {
	t        int64
	f        func(*scriptedPlayer) bool
	expected *bool
}

type script []scriptedRequest

type scriptedPlayer struct {
	username string
	clientID int
	conn     *mockConn
	s        script
}

func makePlayer(username string, clientID int) *scriptedPlayer {
	return &scriptedPlayer{
		username: username,
		clientID: clientID,
		conn:     &mockConn{clientID: clientID},
	}
}

func (me *scriptedPlayer) run(t *testing.T) {
	crtrq := 0
	initTime := time.Now()

	for crtrq < len(me.s) {
		select {
		case <-time.After(time.Duration(10) * time.Millisecond):
			if time.Now().Sub(initTime) > time.Duration(me.s[crtrq].t)*time.Millisecond {
				s := me.s[crtrq].f(me)
				if me.s[crtrq].expected != nil && *(me.s[crtrq].expected) != s {
					t.Errorf("request return incorrect")
				}
				crtrq++
			}
		}
	}
}

func scriptedJSQ(p *scriptedPlayer) bool {
	return joinSessionQueue(p.username, p.clientID, p.conn) == nil
}

func scriptedQSQ(p *scriptedPlayer) bool {
	return quitSessionQueue(p.clientID) == nil
}

func scriptedRC(p *scriptedPlayer) bool {
	return reconnectClient(p.clientID, p.conn) == nil
}

type commandLE struct {
	t      int64
	id     int
	client int
}
type commandLog []commandLE

func (me commandLog) conformTo(cl commandLog) bool {
	if len(me) != len(cl) {
		return false
	}
	for i := 0; i < len(me); i++ {
		found := false
		for j := 0; j < len(cl); j++ {
			if cl[j].id == me[i].id && cl[j].client == me[i].client && math.Abs(float64(cl[j].t-me[i].t)) < 100 {
				cl = append(cl[:j], cl[j+1:]...)
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type resLogger struct {
	initT time.Time
	mch   chan string
	ch    chan Command
	cl    commandLog
}

func (me *resLogger) run() {
	for {
		select {
		case mc := <-me.mch:
			if mc == "quit" {
				goto __quit
			}
		case c := <-me.ch:
			t := time.Now().Sub(me.initT) / time.Millisecond
			me.cl = append(me.cl, commandLE{
				client: c.Source,
				t:      int64(t),
				id:     c.ID,
			})
		}

	}
__quit:
}

func (me *resLogger) getMockConn(clientID int) *mockConn {
	return &mockConn{
		clientID: clientID,
		ch:       me.ch,
	}
}

func makeresLogger() *resLogger {
	return &resLogger{
		initT: time.Now(),
		mch:   make(chan string),
		ch:    make(chan Command),
		cl:    make(commandLog, 0),
	}
}

func TestMain(m *testing.M) {
	initSSrv(SSrvParam{
		Profile:        "test",
		DefaultClients: 2,
		MinimalClients: 2,
		OptimalWS:      30,
		BaseTickMs:     10,
		ABF:            buildMockApp,
	})

	os.Exit(m.Run())
}

func TestScript1(t *testing.T) {
	p1 := makePlayer("regn", 22)
	p2 := makePlayer("yaya", 33)
	tr := true
	p1.s = script{
		{1000, scriptedJSQ, &tr},
	}
	p2.s = script{
		{2000, scriptedJSQ, &tr},
	}

	rl := makeresLogger()
	p1.conn = rl.getMockConn(p1.clientID)
	p2.conn = rl.getMockConn(p2.clientID)
	go rl.run()
	go p1.run(t)
	go p2.run(t)

	<-time.After(time.Duration(3) * time.Second)

	rl.mch <- "quit"

	<-time.After(time.Duration(50) * time.Millisecond)
	if !rl.cl.conformTo(commandLog{
		commandLE{2000, _CMD_INIT_APP, 22},
		commandLE{2000, _CMD_INIT_APP, 33},
		commandLE{2000, CMD_APP_RUN, 22},
		commandLE{2000, CMD_APP_RUN, 33},
	}) {
		t.Errorf("Command Log does not conform to expectation")
	}
}
