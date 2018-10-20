package sgs

import (
	"er"
	"fmt"
	"hlf"
	"math"
	"math/rand"
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

var _tl = hlf.MakeLogger("TestSSRV")

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
		command.Who = cid
		me.s.ForwardToClient(cid, command)
	}
	return nil
}

func buildMockApp() App {
	return &mockApp{}
}

type mockConn struct {
	clientID int
	player   mockPlayer
}

func (me *mockConn) Send(cmd Command) error {
	cmd.Who = me.clientID
	me.player.sendCommand(cmd)
	return nil
}

func (me *mockConn) Run(ch chan Command, mch chan Command) {
	for {
		select {
		case mc := <-mch:
			me.Send(mc)
			if mc.ID == _CMD_CLOSE_NET_CLIENT {
				goto __quit
			}
		}
	}
__quit:
}

func (me *mockConn) BindClientID(clientID int) {}

type requestSender func(mockPlayer) bool

type scriptedRequest struct {
	t        int64
	f        requestSender
	expected *bool
}

type script []scriptedRequest

type mockPlayer interface {
	sendCommand(command Command)
	run(t *testing.T)
	getConn() NetConn
	getClientID() int
	getUsername() string
	reconn()
	getSessionServer() *sessionServer
}

type scriptedPlayer struct {
	username string
	clientID int
	conn     *mockConn
	s        script
	rl       *resLogger
	srv      *sessionServer
	lg       hlf.Logger
}

func makePlayer(username string, clientID int, rl *resLogger, srv *sessionServer) *scriptedPlayer {
	player := scriptedPlayer{
		username: username,
		clientID: clientID,
		conn:     &mockConn{clientID: clientID},
		rl:       rl,
		srv:      srv,
	}
	player.conn.player = &player
	player.lg = srv.lg.Child(username)
	return &player
}

func (me *scriptedPlayer) run(t *testing.T) {
	crtrq := 0
	initTime := time.Now()

	for crtrq < len(me.s) {
		select {
		case <-time.After(time.Duration(500) * time.Millisecond):
			if time.Now().Sub(initTime) > time.Duration(me.s[crtrq].t)*time.Millisecond {
				_tl.Dbg("client %v, rq %v", me.clientID, crtrq)
				s := me.s[crtrq].f(me)
				if me.s[crtrq].expected != nil && *(me.s[crtrq].expected) != s {
					t.Errorf("request return incorrect")
				}
				crtrq++
			}
		}
	}
}

func (me *scriptedPlayer) sendCommand(command Command) {
	me.lg.Inf("Command received: 0x%x, %v", command.ID, command.Who)
	me.rl.ch <- command
}

func (me *scriptedPlayer) getConn() NetConn {
	return me.conn
}

func (me *scriptedPlayer) getClientID() int {
	return me.clientID
}

func (me *scriptedPlayer) getSessionServer() *sessionServer {
	return me.srv
}

func (me *scriptedPlayer) getUsername() string {
	return me.username
}

func (me *scriptedPlayer) reconn() {
	me.conn = &mockConn{
		clientID: me.clientID,
		player:   me,
	}
}

func scriptedJSQ(p mockPlayer) bool {
	return p.getSessionServer().joinSessionQueue(p.getUsername(), p.getClientID(), p.getConn()) == nil
}

func scriptedQSQ(p mockPlayer) bool {
	return p.getSessionServer().quitSessionQueue(p.getClientID()) == nil
}

func scriptedRC(p mockPlayer) bool {
	p.reconn()
	return p.getSessionServer().reconnectClient(p.getClientID(), p.getConn()) == nil
}

type commandLE struct {
	t      int64
	id     int
	client int
}
type commandLog []commandLE

func (me commandLog) conformTo(cl commandLog) bool {
	if len(me) != len(cl) {
		_tl.Ntf("command log length incorrect: %v", me.String())

		return false
	}
	for i := 0; i < len(me); i++ {
		found := false
		for j := 0; j < len(cl); j++ {
			if cl[j].id == me[i].id && cl[j].client == me[i].client && math.Abs(float64(cl[j].t-me[i].t)) < 1200 {
				cl = append(cl[:j], cl[j+1:]...)
				found = true
				break
			}
		}
		if !found {
			_tl.Ntf("command log not found: %x %v %v", me[i].id, me[i].client, me[i].t)
			return false
		}
	}
	return true
}

func (me commandLog) String() string {
	s := fmt.Sprintf("[\n")
	for _, cl := range me {
		s += fmt.Sprintf("Command 0x%x, Client %v, Time %v\n", cl.id, cl.client, cl.t)
	}
	s += fmt.Sprintf("]")
	return s
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
				client: c.Who,
				t:      int64(t),
				id:     c.ID,
			})
		}

	}
__quit:
}

func makeresLogger() *resLogger {
	return &resLogger{
		initT: time.Now(),
		mch:   make(chan string, 100),
		ch:    make(chan Command, 100),
		cl:    make(commandLog, 0),
	}
}

type randomPlayer struct {
	username   string
	clientID   int
	minRqIntr  int
	maxRqIntr  int
	conn       *mockConn
	sleepAfter int
	lg         hlf.Logger
	inSession  bool
	srv        *sessionServer
}

func (me *randomPlayer) run(t *testing.T) {
	rqc := []requestSender{
		scriptedJSQ, scriptedQSQ, scriptedRC,
	}
	rqn := []string{
		"Join Session Queue", "Quit Session Queue", "Reconnect",
	}

	intr := me.minRqIntr + rand.Intn(me.maxRqIntr-me.minRqIntr)

	for {
		select {
		case <-time.After(time.Duration(me.sleepAfter) * time.Millisecond):
			goto __quit
		case <-time.After(time.Duration(intr) * time.Millisecond):
			i := rand.Intn(len(rqc))
			rq := rqc[i]
			rs := rq(me)
			me.lg.Inf("Player %v send request %v, response %v", me.username, rqn[i], rs)
			intr = me.minRqIntr + rand.Intn(me.maxRqIntr-me.minRqIntr)
			if i == 0 && rs {
				me.inSession = true
			} else if i == 1 && rs {
				me.inSession = false
			}
		}

	}
__quit:
}

func (me *randomPlayer) sendCommand(command Command) {
	me.lg.Inf("Command receiverd 0x%x, %v", command.ID, command.Who)
}

func (me *randomPlayer) getClientID() int {
	return me.clientID
}

func (me *randomPlayer) getConn() NetConn {
	return me.conn
}

func (me *randomPlayer) getSessionServer() *sessionServer {
	return me.srv
}

func (me *randomPlayer) getUsername() string {
	return me.username
}

func (me *randomPlayer) reconn() {
	me.conn = &mockConn{
		clientID: me.clientID,
		player:   me,
	}
}

var _playerLog = hlf.MakeLogger("Players")

func makeRandomPlayer(username string, clientID int, minRqIntr int, maxRqIntr int, sleepAfter int, srv *sessionServer) *randomPlayer {
	player := randomPlayer{
		username: username,
		clientID: clientID,
		conn: &mockConn{
			clientID: clientID,
		},
		minRqIntr:  minRqIntr,
		maxRqIntr:  maxRqIntr,
		sleepAfter: sleepAfter,
		lg:         _playerLog.Child(username),
		inSession:  false,
		srv:        srv,
	}
	player.conn.player = &player
	return &player
}
