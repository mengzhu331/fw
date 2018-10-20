package test

import (
	"encoding/json"
	"fmt"
	"fwb"
	"fwb/actn"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sgs"
	"strconv"
	"sutil"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const (
	_1P_ID   = 0x80000001
	_1P_NAME = "regn"
	_2P_ID   = 0x80000002
	_2P_NAME = "yaya"
)

var _ch = make(chan int)

type mockClient struct {
	clientID int
	t        *testing.T
	c        *websocket.Conn
	name     string
	gd       fwb.GameData
	pinfo    map[string]int
	pbi      map[int]string

	gameOver bool
	cch      chan string
}

type commandReaction func(me *mockClient, command sgs.Command) bool

var _cmdReaction = map[int]commandReaction{
	fwb.CMD_GAME_START:              onGameStart,
	fwb.CMD_SYNC_GAME_STATE:         onSyncGameState,
	fwb.CMD_ROUND_START:             onRoundStart,
	fwb.CMD_START_TURN:              onStartTurn,
	fwb.CMD_ACTION_COMMITTED:        onActionCommitted,
	fwb.CMD_ROUND_SETTLEMENT:        onRoundSettlement,
	fwb.CMD_ROUND_SETTLEMENT_UPDATE: onRoundSettlementUpdate,
	fwb.CMD_GAME_FINISH:             onGameFinish,
	fwb.CMD_GAME_OVER:               onGameOver,
	fwb.CMD_ACTION_REJECTED:         onActionRejected,
}

func (me *mockClient) connect() {

	me.gameOver = false
	me.cch = make(chan string)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	host := fmt.Sprintf("127.0.0.1:%v", _conf.Port)
	u := "ws://" + host + "/join_session?username=" + me.name + "&clientid=" + strconv.Itoa(me.clientID) + "&token=xxx"
	log.Printf("connecting to %s", u)

	var err error
	me.c, _, err = websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Printf("Failed to connect to FWB, %v", err)
	}
	defer me.c.Close()

	go func() {
		for {
			_, message, err := me.c.ReadMessage()
			if err != nil {
				log.Fatal("error read ws", err)
			}
			command := sgs.Command{}
			json.Unmarshal(message, &command)
			log.Printf("recv: 0x%v, %v, %v", command.HexID(), command.Who, command.Payload)
			if !_cmdReaction[command.ID](me, command) {
				log.Fatal("Error with command", command.ID)
			}
			if me.gameOver {
				break
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-interrupt:
			log.Printf("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := me.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Fatal("error write ws close", err)
			}
			select {
			case <-time.After(time.Second):
			}
			return
		case <-me.cch:
			_ch <- me.clientID
			return
		}

	}
}

func (me *mockClient) sendToGame(command sgs.Command) error {
	return me.c.WriteJSON(&command)
}

func onGameStart(me *mockClient, command sgs.Command) bool {
	type gameStartCommand struct {
		ID      int
		Source  int
		Payload map[string]int
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Printf("Failed to decode game start command, %v %v", command, err.Error())
		return false
	}

	gs := gameStartCommand{}
	err = json.Unmarshal(cb, &gs)
	if err != nil {
		log.Printf("Failed to decode game start command, %v %v", command, err.Error())
		return false
	}

	me.pinfo = gs.Payload
	me.pbi = make(map[int]string)
	for k, v := range me.pinfo {
		me.pbi[v] = k
	}

	if !verifyPInfo(me) {
		return false
	}

	log.Printf("%v send command %x", me.name, fwb.CMD_GAME_START_ACK)
	return me.sendToGame(sgs.Command{
		ID:  fwb.CMD_GAME_START_ACK,
		Who: me.clientID,
	}) == nil
}

func onRoundStart(me *mockClient, command sgs.Command) bool {
	type targetCommand struct {
		ID      int
		Source  int
		Payload []fwb.Card
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode command", command, err.Error())
	}

	log.Printf("%v received cards: %v", me.name, tc.Payload)
	me.gd.Cards = tc.Payload

	return me.sendToGame(sgs.Command{
		ID:  fwb.CMD_ROUND_START_ACK,
		Who: me.clientID,
	}) == nil
}

func onSyncGameState(me *mockClient, command sgs.Command) bool {
	type targetCommand struct {
		ID      int
		Source  int
		Payload fwb.GameData
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode Sync Game State command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode Sync Game State command", command, err.Error())
	}

	me.syncGameState(tc.Payload)
	return true
}

func onActionCommitted(me *mockClient, command sgs.Command) bool {
	type targetCommand struct {
		ID      int
		Source  int
		Payload fwb.GameData
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode Action Committed command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode Action Committed command", command, err.Error())
	}

	me.syncGameState(tc.Payload)
	return true
}

func onRoundSettlementUpdate(me *mockClient, command sgs.Command) bool {
	type targetCommand struct {
		ID      int
		Source  int
		Payload fwb.GameData
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode Round Settlement Update command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode Round Settlement Update command", command, err.Error())
	}

	me.syncGameState(tc.Payload)
	return true
}

func (me *mockClient) syncGameState(gd fwb.GameData) {
	me.gd = gd
	log.Printf("%v updated game state: %v", me.name, gd)
}

func onStartTurn(me *mockClient, command sgs.Command) bool {
	type targetCommand struct {
		ID      int
		Source  int
		Payload int
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode command", command, err.Error())
	}

	var actionCmd *sgs.Command
	if tc.Payload == me.clientID {
		log.Printf("%v received Start Turn command, my turn", me.name)
		actionCmd = me.randomAction()
	} else {
		log.Printf("%v received Start Turn command, %v turn", me.name, me.pbi[tc.Payload])
	}

	if actionCmd != nil {
		return me.sendToGame(*actionCmd) == nil
	}

	return true
}

func verifyPInfo(me *mockClient) bool {
	return me.pinfo[_1P_NAME] == _1P_ID && me.pinfo[_2P_NAME] == _2P_ID && len(me.pinfo) == 2
}

type actnCmd struct {
	ActionID int
	Payload  interface{}
}

func (me *mockClient) randomAction() *sgs.Command {

	ac := &actnCmd{}
	actions := make([]int, actn.ACTN_MAX-actn.ACTN_SKIP-1)
	for i := 0; i < actn.ACTN_MAX-actn.ACTN_SKIP-1; i++ {
		actions[i] = actn.ACTN_SKIP + i + 1
	}

	if rand.Intn(len(actions)+1) < 1 {
		goto __skip
	}

	actions = sutil.ShuffleInt(actions...)

	for _, action := range actions {
		if !actn.HasCardSlots(&me.gd, action) {
			continue
		}
		acpl, affordable := me.randomActionCmdPayload(action)
		if affordable {
			ac.ActionID = action
			ac.Payload = acpl
			break
		}
	}

__skip:
	if ac == nil {
		ac.ActionID = actn.ACTN_SKIP
	}

	log.Printf("%v request action %v %v", me.name, actn.ActionNames[ac.ActionID], ac)
	return &sgs.Command{
		ID:      fwb.CMD_ACTION,
		Who:     me.clientID,
		Payload: ac,
	}
}

func (me *mockClient) randomActionCmdPayload(action int) (interface{}, bool) {
	var payload interface{}
	sp := 0

	i := me.gd.GetPDIndex(me.clientID)
	data := me.gd.PData[i]

	inte := data[fwb.PD_SK_INTELLIGENCE]
	know := data[fwb.PD_SK_KNOWLEDGE]
	stre := data[fwb.PD_SK_STRENGTH]

	switch action {
	case actn.ACTN_FARM:
		fallthrough
	case actn.ACTN_TAKE_OFF:
		fallthrough
	case actn.ACTN_PARTTIME_WORK:
		fallthrough
	case actn.ACTN_HUNT:
		fallthrough
	case actn.ACTN_BEG:
		sp = 1
	case actn.ACTN_GOLD_MINING:
		if data[fwb.PD_SK_KNOWLEDGE] >= 1 && data[fwb.PD_SK_STRENGTH] >= 1 {
			sp = 1
		}

	case actn.ACTN_EMPLOY:
		if data[fwb.PD_PT_GOLD] >= 10 {
			sp = 1
		}

	case actn.ACTN_FEED_SHEEP:
		log.Print(actn.ActionNames[actn.ACTN_FEED_SHEEP], " ", data[fwb.PD_PT_CEREALS])
		if data[fwb.PD_PT_CEREALS] >= 2 {
			sp = 1
		}

	case actn.ACTN_TAKE_VACATION:
		if data[fwb.PD_PT_GOLD] >= 5 {
			sp = 1
		}

	case actn.ACTN_PARTY:
		log.Print(actn.ActionNames[actn.ACTN_PARTY], " ", data[fwb.PD_PT_WINE])
		if data[fwb.PD_PT_WINE] >= 2 && data[fwb.PD_PT_MEAT] >= 2 {
			sp = 1
		}

	case actn.ACTN_UPGRADE_HOUSE:
		if data[fwb.PD_PT_GOLD] < 10 {
			break
		}

		if data[fwb.PD_HOUSE_LV] == 0 {
			sp = 3
		} else if data[fwb.PD_HOUSE_LV] == 1 {
			if know >= 1 && inte >= 1 && stre >= 1 {
				sp = 3
			}
		} else if data[fwb.PD_HOUSE_LV] == 2 {
			if know >= 2 && inte >= 2 && stre >= 2 {
				sp = 3
			}
		}

	case actn.ACTN_MAKE_WINE:
		if data[fwb.PD_PT_CEREALS] >= 8 {
			sp = 1
		}

	case actn.ACTN_TRAIN:
		s := rand.Intn(3) + fwb.PD_SK_STRENGTH
		c := [3]int{2, 5, 10}
		for i := 0; i < 3; i++ {
			s = s + i
			if s > fwb.PD_SK_INTELLIGENCE {
				s = fwb.PD_SK_STRENGTH
			}
			if data[s] < 3 && data[fwb.PD_PT_GOLD] >= c[data[s]] {
				sp = 1
				payload = s
				break
			}
		}

	case actn.ACTN_TRADE:
		type Payload struct {
			Direction int
			Amount    fwb.PlayerData
		}
		pld := Payload{
			Amount: make(fwb.PlayerData, fwb.PD_MAX),
		}
		if rand.Intn(2) > 0 {
			if data[fwb.PD_PT_GOLD] < 2 {
				break
			}
			pld.Direction = -1
			for i := fwb.PD_PT_CEREALS; i <= fwb.PD_PT_WINE; i++ {
				if data[i] > 0 {
					pld.Amount[i] = 1
				}
			}
		} else {
			pld.Direction = 1
			ptp := actn.PtPrices()
			p := 0
			for i := fwb.PD_PT_CEREALS; i <= fwb.PD_PT_WINE; i++ {
				if rand.Intn(3) > 0 {
					continue
				}

				if data[fwb.PD_PT_GOLD] >= p+ptp[i]+2 {
					pld.Amount[i] = 1
					p += ptp[i]
				}
			}
			if data[fwb.PD_PT_GOLD] < p+2 {
				break
			}
		}
		sp = 1
		payload = pld

	case actn.ACTN_STEAL:
		var px int
		for px = 1; px < len(me.gd.PData); px++ {
			if me.clientID != data[fwb.PD_CLIENT_ID] &&
				data[fwb.PD_SK_INTELLIGENCE] > me.gd.PData[px][fwb.PD_SK_KNOWLEDGE] {
				payload = me.gd.PData[px][fwb.PD_CLIENT_ID]
				break
			}
		}
		if px != len(me.gd.PData) {
			sp = 1
		}

	case actn.ACTN_WEAVE:
		if data[fwb.PD_PT_WOOL] >= 2 {
			sp = 1
		}
	}

	if sp > 0 && data[fwb.PD_PAWNS] >= sp {
		return payload, true
	}

	return payload, false
}

func onRoundSettlement(me *mockClient, command sgs.Command) bool {

	type targetCommand struct {
		ID      int
		Source  int
		Payload interface{}
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode Round Settlement command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode Round Settlement command", command, err.Error())
	}

	if tc.ID != fwb.CMD_ROUND_SETTLEMENT {
		return false
	}

	type playerSData struct {
		Cereals int
		Meat    int
		Sweater int
	}

	sd := playerSData{}

	px := me.gd.GetPDIndex(me.clientID)
	data := me.gd.PData[px]
	if data[fwb.PD_PT_CEREALS] > 0 {
		sd.Cereals = 1
	}

	if data[fwb.PD_PT_MEAT] > 0 {
		sd.Meat = 1
	}

	if data[fwb.PD_PT_SWEATER] > 0 {
		sd.Sweater = 1
	}

	csc := sgs.Command{
		ID:      fwb.CMD_COMMIT_ROUND_SETTLEMENT,
		Who:     me.clientID,
		Payload: sd,
	}
	return me.sendToGame(csc) == nil
}

func onGameFinish(me *mockClient, command sgs.Command) bool {
	type gameFinish struct {
		ID      int
		Who     int
		Payload map[int]int
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Printf("Failed to decode Game Finish command, %v %v", command, err.Error())
		return false
	}

	gs := gameFinish{}
	err = json.Unmarshal(cb, &gs)
	if err != nil {
		log.Printf("Failed to decode Game Finish command, %v %v", command, err.Error())
		return false
	}

	rankMap := gs.Payload

	log.Printf("Game Finish, player rank:")
	for p, r := range rankMap {
		log.Printf("player %v, rank %v", me.pbi[p], r)
	}

	if rand.Intn(2) > 0 {
		log.Printf("player %v requests rematch", me.name)
		return me.sendToGame(sgs.Command{
			ID:  fwb.CMD_REMATCH,
			Who: me.clientID,
		}) == nil
	}
	return true
}

func onGameOver(me *mockClient, command sgs.Command) bool {
	type targetCommand struct {
		ID      int
		Source  int
		Payload interface{}
	}

	cb, err := json.Marshal(command)
	if err != nil {
		log.Fatal("Failed to decode Game Over command", command, err.Error())
	}

	tc := targetCommand{}
	err = json.Unmarshal(cb, &tc)
	if err != nil {
		log.Fatal("Failed to decode Game Over command", command, err.Error())
	}

	if tc.ID != fwb.CMD_GAME_OVER {
		return false
	}

	me.gameOver = true
	me.cch <- "quit"
	return true
}

func onActionRejected(me *mockClient, command sgs.Command) bool {
	log.Fatal("Action rejected")
	return false
}
