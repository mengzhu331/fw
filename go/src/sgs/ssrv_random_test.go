package sgs

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomPlayers(t *testing.T) {
	srv, _ := makeSSrv(SSrvParam{
		Profile:        "test",
		DefaultClients: 2,
		MinimalClients: 2,
		OptimalWS:      30,
		BaseTickMs:     10,
		ABF:            buildMockApp,
	})

	players := [50]*randomPlayer{}

	for i := 0; i < len(players); i++ {
		name := ""
		for j := 0; j < 4; j++ {
			name += string('a' + rune(rand.Intn(26)))
		}
		minRqIntr := 100 + rand.Intn(300)
		maxRqIntr := 500 + rand.Intn(1000)
		player := makeRandomPlayer(name, i+0x3000, minRqIntr, maxRqIntr, 10000, srv)
		_log.Inf("player %v spawned, %v %v", player.getUsername(), player.minRqIntr, player.maxRqIntr)
		go player.run(t)
		players[i] = player
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	time.Sleep(time.Duration(10000) * time.Millisecond)

	ise := 0
	for i := 0; i < len(players); i++ {
		if players[i].inSession {
			ise++
		}
	}
	if len(srv.sessions) != ise/2 {
		t.Errorf("not all players in session, miss %v", ise-len(srv.sessions)*2)
	}
}
