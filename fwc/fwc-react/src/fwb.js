const CMD_C_APP = 0x02000000

const CMD_C_CLIENT = 0x04000000

const CMD_C_CLIENT_TO_APP = CMD_C_CLIENT | 0x00010000

const CMD_C_APP_TO_CLIENT = CMD_C_APP | 0x00010000

const CMD_GAME_START = CMD_C_APP_TO_CLIENT + 1

//CMD_GAME_OVER command for notifying players that game is up
const CMD_GAME_OVER = CMD_C_APP_TO_CLIENT + 2

//CMD_ROUND_START command for notifying players that a new round started
const CMD_ROUND_START = CMD_C_APP_TO_CLIENT + 3

//CMD_START_TURN command for notifying players that a new turn started
const CMD_START_TURN = CMD_C_APP_TO_CLIENT + 4

//CMD_ACTION_REJECTED command for notifying players that the requested action is rejected
const CMD_ACTION_REJECTED = CMD_C_APP_TO_CLIENT + 5

//CMD_ACTION_COMMITTED command for notifying players that the action has been committed
const CMD_ACTION_COMMITTED = CMD_C_APP_TO_CLIENT + 6

//CMD_ROUND_SETTLEMENT command for notifying players to make settlement
const CMD_ROUND_SETTLEMENT = CMD_C_APP_TO_CLIENT + 7

//CMD_ROUND_SETTLEMENT_INVALID command for notifying a player that the data with the requested settlement is invalid
const CMD_ROUND_SETTLEMENT_INVALID = CMD_C_APP_TO_CLIENT + 8

//CMD_ROUND_SETTLEMENT_UPDATE command for notifying players that the settlement has updated the game data
const CMD_ROUND_SETTLEMENT_UPDATE = CMD_C_APP_TO_CLIENT + 9

//CMD_GAME_FINISH command for notifying players that game has finished with player rank ready
const CMD_GAME_FINISH = CMD_C_APP_TO_CLIENT + 10

//CMD_SYNC_GAME_STATE command for synchronize game state between components
const CMD_SYNC_GAME_STATE = CMD_C_APP_TO_CLIENT + 11

//CMD_GAME_START_ACK command for notifying the game system that the player acknowledged game having started
const CMD_GAME_START_ACK = CMD_C_CLIENT_TO_APP + 1

//CMD_ROUND_START_ACK command for notifying the game system that the player ackowledged a new round having started
const CMD_ROUND_START_ACK = CMD_C_CLIENT_TO_APP + 2

//CMD_ACTION command for requesting an action
const CMD_ACTION = CMD_C_CLIENT_TO_APP + 3

//CMD_COMMIT_ROUND_SETTLEMENT command for commiting round settlement
const CMD_COMMIT_ROUND_SETTLEMENT = CMD_C_CLIENT_TO_APP + 4

//CMD_REMATCH command for requesting next match with the same players
const CMD_REMATCH = CMD_C_CLIENT_TO_APP + 5

//PD_CLIENT_ID index in PlayerData for player client ID
const PD_CLIENT_ID = 0

//PD_PT_HEART index in PlayerData for property Heart
const PD_PT_HEART = 1

//PD_PT_GOLD index in PlayerData for property Gold
const PD_PT_GOLD = 2

//PD_PT_CEREALS index in PlayerData for property Cereals
const PD_PT_CEREALS = 3

//PD_PT_MEAT index in PlayerData for property Meat
const PD_PT_MEAT = 4

//PD_PT_WOOL index in PlayerData for property Woo
const PD_PT_WOOL = 5

//PD_PT_SWEATER index in PlayerData for property Sweater
const PD_PT_SWEATER = 6

//PD_PT_BEER index in PlayerData for property Wine
const PD_PT_BEER = 7

//PD_SK_STRENGTH index in PlayerData for skill Strength
const PD_SK_STRENGTH = 8

//PD_SK_KNOWLEDGE index in PlayerData for skill Knowledge
const PD_SK_KNOWLEDGE = 9

//PD_SK_INTELLIGENCE index in PlayerData for skill Intelligence
const PD_SK_INTELLIGENCE = 10

//PD_HOUSE_LV index in PlayerData for the level of the player house
const PD_HOUSE_LV = 11

//PD_MAX_PAWNS index in PlayerData for the max pawns
const PD_MAX_PAWNS = 12

//PD_PAWNS index in PlayerData for the pawns left
const PD_PAWNS = 13

//PD_MAX max index in PlayerData
const PD_MAX = 14

class FWB {
    constructor() {
        this.wsconn = null
        this.onError = null
        this.onOpen = null
        this.onGameStart = null
        this.onActionCommitted= null
        this.onActionRejected= null
        this.onGameFinished= null
        this.onGameOver= null
        this.onRoundSettlement= null
        this.onRoundSettlementUpdate= null
        this.onRoundSettlementInvalid= null
        this.onRoundStart= null
        this.onStartTurn= null
        this.onSyncGameState= null

        this.__onGameStart = this.__onGameStart.bind(this)
        this.__onActionCommitted = this.__onActionCommitted.bind(this)
        this.__onActionRejected = this.__onActionRejected.bind(this)
        this.__onDefaultCmd = this.__onDefaultCmd.bind(this)
        this.__onGameFinish = this.__onGameFinish.bind(this)
        this.__onGameOver = this.__onGameOver.bind(this)
        this.__onRoundSettlement = this.__onRoundSettlement.bind(this)
        this.__onRoundSettlementInvalid = this.__onRoundSettlementInvalid.bind(this)
        this.__onRoundSettlementUpdate = this.__onRoundSettlementUpdate.bind(this)
        this.__onRoundStart = this.__onRoundStart.bind(this)
        this.__onStartTurn = this.__onStartTurn.bind(this)
        this.__onSyncGameState = this.__onSyncGameState.bind(this)

        this.gd = {
            pd: new Map(),
            cards: [
                {
                    id:1,
                    token:[],
                },
                {
                    id:2,
                    token:[],
                },
                {
                    id:3,
                    token:[],
                },
                {
                    id:4,
                    token:[],
                },
                {
                    id:5,
                    token:[],
                },
                {
                    id:6,
                    token:[],
                },
                {
                    id:7,
                    token:[],
                },
                {
                    id:8,
                    token:[],
                },
                {
                    id:9,
                    token:[],
                }
            ],
            round: 0,
            clientID: 0x80000001,
            username:"",
            pdlist:[],
            currentPlayer: 0,
        }

        this.cmdMap = new Map([
            [CMD_GAME_START, this.__onGameStart.bind(this)],
            [CMD_GAME_OVER, this.__onGameOver.bind(this)],
            [CMD_ROUND_START, this.__onRoundStart.bind(this)],
            [CMD_START_TURN, this.__onStartTurn.bind(this)],
            [CMD_ACTION_REJECTED, this.__onActionRejected.bind(this)],
            [CMD_ACTION_COMMITTED, this.__onActionCommitted.bind(this)],
            [CMD_ROUND_SETTLEMENT, this.__onRoundSettlement.bind(this)],
            [CMD_ROUND_SETTLEMENT_INVALID, this.__onRoundSettlementInvalid.bind(this)],
            [CMD_ROUND_SETTLEMENT_UPDATE, this.__onRoundSettlementUpdate.bind(this)],
            [CMD_GAME_FINISH, this.__onGameFinish.bind(this)],
            [CMD_SYNC_GAME_STATE, this.__onSyncGameState.bind(this)]
        ])
    }

    updateGameData(gd) {
        console.log("update game data", gd)
        this.gd.round = gd.Round
        this.updatePlayerData(gd.PData)
        this.updateCards(gd.Cards)

    }

    updateCards(cards) {
        this.gd.cards = []
        console.log("update cards", cards)
        if (cards== null){
            return
        }
        console.log("cards length", cards.length)
        for (var i = 0; i < cards.length; i++) {
            let card = {
                id: cards[i].ID,
                token: [],
            }
            if(cards[i].Pawns!= null){
                for (var j = 0; j < cards[i].Pawns.length; j++) {
                    let color = this.gd.pd.get(cards[i].Pawns[j]).color
                    card.token.push(color)
                }
            }
            console.log("add card", card)

            this.gd.cards.push(card)
        }
        console.log("fwb cards", this.gd.cards)
    }

    updatePlayerData(pd) {
        console.log("update player data", pd)
        let colors = ["red", "blue", "yellow", "green", "black", "brown", "orange"]
        let pdmap= this.gd.pd

        this.gd.pdlist=[]
        this.gd.pd= new Map()
        if(pd== null){
            return
        }

        for (var i = 0; i < pd.length; i++) {
            let clientID = pd[i][PD_CLIENT_ID]
            let p={
                clientID: clientID,
                cereals: pd[i][PD_PT_CEREALS],
                meat: pd[i][PD_PT_MEAT],
                wool: pd[i][PD_PT_WOOL],
                sweater: pd[i][PD_PT_SWEATER],
                beer: pd[i][PD_PT_BEER],
                heart: pd[i][PD_PT_HEART],
                gold: pd[i][PD_PT_GOLD],
    
                str: pd[i][PD_SK_STRENGTH],
                knw: pd[i][PD_SK_KNOWLEDGE],
                int:pd[i][PD_SK_INTELLIGENCE],
                hm: pd[i][PD_HOUSE_LV],
                token: pd[i][PD_PAWNS],
                name: pdmap.get(clientID).name,
            }
            this.gd.pd.set(clientID, p)
            this.gd.pdlist.push(p)
        }

        this.gd.pdlist.sort(function(a, b){
            return a[PD_CLIENT_ID]- b[PD_CLIENT_ID]
        })

        let pdlist= this.gd.pdlist
        for(i= 0; i<pdlist.length; i++){
            let p= this.gd.pdlist[i]
            p.color= colors[i]
            this.gd.pdlist[i]= p
            this.gd.pd.set(p.clientID, p)
        }
        
    }

    connError(event) {
        console.log(event)
    }

    fwbMessage(event) {
        console.log(event)
    }

    sendGameStartAck(){
        let gsAck = {
            "ID": CMD_GAME_START_ACK,
            "Who": this.gd.clientID,
            "Payload": null
        }
        this.wsconn.send(JSON.stringify(gsAck))
    }

    sendRoundStartAck(){
        let rsAck = {
            "ID": CMD_ROUND_START_ACK,
            "Who": this.gd.clientID,
            "Payload": null
        }
        this.wsconn.send(JSON.stringify(rsAck))
    }

    __onWSError(event) {
        console.log(event)
        if (this.onError != null) {
            this.onError(event)
        }
    }

    __onActionCommitted(cmd) {
        console.log("action committed", cmd.Payload)
        this.updateGameData(cmd.Payload)
        
        if(this.onActionCommitted!= null){
            this.onActionCommitted(cmd)
        }
    }

    __onActionRejected(cmd) {
        if(this.onActionRejected!= null){
            this.onActionRejected(cmd)
        }
    }

    __onGameFinish(cmd) {
        let playerList=[]
        for(var i=0; i< cmd.Payload.keys.length; i++){
            let clientID= cmd.Payload.keys[i]
            let playerInfo={
                name: this.gd.pd.get(clientID).name,
                rank: cmd.Payload[clientID],
            }
            playerList.push(playerInfo)
        }

        playerList.sort(function(a, b){
            return a.rank- b.rank
        })

        let youwin= 0
        if ((cmd.Payload[this.gd.clientID]=== 1) && (playerList[playerList.length-1].rank!== 1)){
            youwin= 1
        }else if (cmd.Payload[this.gd.clientID]!== 1){
            youwin= -1
        }
        if(this.onGameFinished!= null){
            this.onGameFinished(cmd, playerList, youwin)
        }
    }

    __onGameOver(cmd) {
        if(this.onGameOver!= null){
            this.onGameOver(cmd)
        }
    }

    __onRoundSettlement(cmd) {
        if(this.onRoundSettlement!= null){
            this.onRoundSettlement(cmd)
        }
    }

    __onRoundSettlementInvalid(cmd) {
        if(this.onRoundSettlementInvalid!= null){
            this.onRoundSettlementInvalid(cmd)
        }
    }

    __onRoundSettlementUpdate(cmd) {
        this.updateGameData(cmd.Payload)

        if(this.onRoundSettlementUpdate!= null){
            this.onRoundSettlementUpdate(cmd)
        }
    }

    __onRoundStart(cmd) {
        this.updateGameData(cmd.Payload)

        if(this.onRoundStart!= null){
            if(this.onRoundStart(cmd)){
                this.sendRoundStartAck()
            }
        }
    }

    __onStartTurn(cmd) {
        let myturn= cmd.Payload=== this.gd.clientID
        let whoseturn= this.gd.pd.get(cmd.Payload).name
        this.gd.currentPlayer= cmd.Payload

        if(this.onStartTurn!= null){
            this.onStartTurn(cmd, myturn, whoseturn)
        }
    }

    __onSyncGameState(cmd) {
        this.updateGameData(cmd.Payload)

        if(this.onSyncGameState!=null){
            this.onSyncGameState(cmd)
        }
    }

    __onGameStart(cmd) {
        console.log(cmd)
        this.gd.pd= new Map()
        this.gd.pdlist= []
        let playerInfo = new Map(Object.entries(cmd.Payload))
        console.log(playerInfo)
        let keys= [...playerInfo.keys()]
        for (var i = 0; i < keys.length; i++) {
            let clientID = playerInfo.get(keys[i])
            let pn = keys[i]
            let p={
                name: pn,
                clientID: clientID,
            }
            console.log(p)
            this.gd.pd.set(clientID, p)
            this.gd.pdlist.push(p)
        }
        console.log(this.gd)

        if (this.onGameStart == null) {
            console.log("No game start callback")
            return
        }

        if (this.onGameStart(cmd)) {
            this.sendGameStartAck()
        }
    }

    __onDefaultCmd(cmd) {
        console.log("No handler for command: ", cmd)
    }

    __onWSMessage(event) {

        console.log(event.data)
        var cmd = JSON.parse(event.data)
        console.log(cmd.ID)
        console.log(this.cmdMap)
        if(!this.cmdMap.has(cmd.ID)){
            this.__onDefaultCmd(cmd)
        }else{
            let handler = this.cmdMap.get(cmd.ID)
            console.log(handler)
            handler(cmd)
        }
    }

    __onOpen(event) {
        console.log(event)
        this.wsconn.onerror = this.__onWSError.bind(this)
        this.wsconn.onmessage = this.__onWSMessage.bind(this)
        if (this.onOpen != null) {
            this.onOpen(event)
        }
    }

    doConnect(username) {
        var url = "ws://10.0.1.5:9090"
        var endpoint = "join_session"
        var wsrequest = url + "/" + endpoint + "?" + "username=" + username + "&" + "clientid=" + this.gd.clientID + "&" + "token=" + this.token

        this.wsconn = new WebSocket(wsrequest);
        this.wsconn.onerror = function (event) {
            console.log(event)
            this.retry++
            if (this.retry >= 10) {
                alert("Failed to connect to FWB after 10 retrys")
                return
            }
            this.doConnect(username, ++this.gd.clientID, this.token)
        }.bind(this)

        this.wsconn.onopen = this.__onOpen.bind(this)
    }

    connect(username, password) {
        console.log(username, password, this.gd)
        this.gd.username= username
        this.retry = 0
        this.doConnect(username)
    }
}

export default new FWB()