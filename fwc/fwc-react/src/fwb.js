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

//CARD_VOID used for undefined card
const CARD_VOID = 0

//CARD_FARM card to enable farm action
const CARD_FARM = 1

//CARD_FEED_SHEEP card to enable feed sheep action
const CARD_FEED_SHEEP = 2

//CARD_TAKE_OFF card to enable take off and take vacation action
const CARD_TAKE_OFF = 3

//CARD_PARTTIME_WORK card to enable parttime work action
const CARD_PARTTIME_WORK = 4

//CARD_UPGRADE_HOUSE card to enable upgrade house action
const CARD_UPGRADE_HOUSE = 5

//CARD_TRAIN card to enable train skills action
const CARD_TRAIN = 6

//CARD_EMPLOY card to enable employ action
const CARD_EMPLOY = 7

//CARD_TRADE card to enable trade action
const CARD_TRADE = 8

//CARD_HUNT card to enable hunt action
const CARD_HUNT = 9

//CARD_STEAL card to enable steal action
const CARD_STEAL = 10

//CARD_WEAVE card to enable weave wool action
const CARD_WEAVE = 11

//CARD_MAKE_WINE card to enable make wine action
const CARD_MAKE_WINE = 12

//CARD_PARTY card to enable party action
const CARD_PARTY = 13

//CARD_BEG card to enable beg action
const CARD_BEG = 14

//CARD_GOLD_MINING card to enable gold mining action
const CARD_GOLD_MINING = 15

//CARD_MAX boundary value of card index
const CARD_MAX = 16

//ACTN_SKIP skip one turn
const ACTN_SKIP = 0

//ACTN_FARM farm action
const ACTN_FARM = 1

//ACTN_FEED_SHEEP feed sheep action
const ACTN_FEED_SHEEP = 2

//ACTN_TAKE_OFF take off action
const ACTN_TAKE_OFF = 3

//ACTN_PARTTIME_WORK parttime work action
const ACTN_PARTTIME_WORK = 4

//ACTN_UPGRADE_HOUSE upgrade house action
const ACTN_UPGRADE_HOUSE = 5

//ACTN_TRAIN train skills action
const ACTN_TRAIN = 6

//ACTN_EMPLOY employ action
const ACTN_EMPLOY = 7

//ACTN_TRADE trade action
const ACTN_TRADE = 8

//ACTN_HUNT hunt action
const ACTN_HUNT = 9

//ACTN_STEAL steal action
const ACTN_STEAL = 10

//ACTN_WEAVE weave wool action
const ACTN_WEAVE = 11

//ACTN_MAKE_WINE make wine action
const ACTN_MAKE_WINE = 12

//ACTN_PARTY party action
const ACTN_PARTY = 13

//ACTN_BEG beg action
const ACTN_BEG = 14

//ACTN_GOLD_MINING gold mining action
const ACTN_GOLD_MINING = 15

//ACTN_TAKE_VACATION take vacation action
const ACTN_TAKE_VACATION = 16

//ACTN_MAX max action index
const ACTN_MAX = 17

const cardActionMap = new Map([
    [CARD_FARM, ["farm"]],
    [CARD_FEED_SHEEP, ["feed-sheep"]],
    [CARD_PARTTIME_WORK, ["parttime-work"]],
    [CARD_TRADE, ["trade"]],
    [CARD_TAKE_OFF, ["take-off", "take-vacation"]],
    [CARD_UPGRADE_HOUSE, ["upgrade-home"]],
    [CARD_TRAIN, ["train"]],
    [CARD_HUNT, ["hunt"]],
    [CARD_BEG, ["beg"]],
    [CARD_EMPLOY, ["employ"]],
    [CARD_STEAL, ["steal"]],
    [CARD_GOLD_MINING, ["gold-mining"]],
    [CARD_VOID, ["skip"]],
    [CARD_MAKE_WINE, ["make-wine"]],
    [CARD_WEAVE, ["weave"]],
    [CARD_PARTY, ["party"]],
])

const actionMap = new Map([
    ["farm", ACTN_FARM],
    ["feed-sheep", ACTN_FEED_SHEEP],
    ["parttime-work", ACTN_PARTTIME_WORK],
    ["trade", ACTN_TRADE],
    ["take-off", ACTN_TAKE_OFF],
    ["take-vacation", ACTN_TAKE_VACATION],
    ["upgrade-home", ACTN_UPGRADE_HOUSE],
    ["train", ACTN_TRAIN],
    ["hunt", ACTN_HUNT],
    ["beg", ACTN_BEG],
    ["employ", ACTN_EMPLOY],
    ["steal", ACTN_STEAL],
    ["gold-mining", ACTN_GOLD_MINING],
    ["skip", ACTN_SKIP],
    ["make-wine", ACTN_MAKE_WINE],
    ["weave", ACTN_WEAVE],
    ["party", ACTN_PARTY]
])

class FWB {
    constructor() {
        this.wsconn = null
        this.onError = null
        this.onOpen = null
        this.onGameStart = null
        this.onActionCommitted = null
        this.onActionRejected = null
        this.onGameFinished = null
        this.onGameOver = null
        this.onRoundSettlement = null
        this.onRoundSettlementUpdate = null
        this.onRoundSettlementInvalid = null
        this.onRoundStart = null
        this.onStartTurn = null
        this.onSyncGameState = null

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
            pd: new Map([
                [0x80000001, {
                    cereals: 0,
                    gold: 0,
                    heart:0,
                    meat: 0,
                    int: 0,
                    knw: 0,
                    str: 0,
                    wool: 0,
                    sweater: 0,
                    beer: 0,
                    token: 5,
                    color: "red"
                }]
            ]),
            cards: [{
                    id: 1,
                    token: [],
                },
                {
                    id: 2,
                    token: [],
                },
                {
                    id: 3,
                    token: [],
                },
                {
                    id: 4,
                    token: [],
                },
                {
                    id: 5,
                    token: [],
                },
                {
                    id: 6,
                    token: [],
                },
                {
                    id: 7,
                    token: [],
                },
                {
                    id: 8,
                    token: [],
                },
                {
                    id: 10,
                    token: [],
                }
            ],
            round: 0,
            clientID: 0x80000001,
            username: "",
            pdlist: [],
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

    getPd() {
        console.log(this.gd.pd)
        console.log(this.gd.clientID)
        console.log(this.gd.pd.get(this.gd.clientID))
        return this.gd.pd.get(this.gd.clientID)
    }

    sendSettlement(order){
        let command = {
            ID: CMD_COMMIT_ROUND_SETTLEMENT,
            Who: this.gd.clientID,
            Payload: order,
        }
        console.log("send settlement ", command)
        this.wsconn.send(JSON.stringify(command))
    }

    sendAction(rawAction) {

        let actionCmd={
            ActionID: 0,
            Payload:{}
        }

        if(rawAction.action==="trade"){
            if(rawAction.param.direction==="buy"){
                actionCmd.Payload.Direction= 1
            }else{
                actionCmd.Payload.Direction= -1
            }
            actionCmd.Payload.Amount=Array.apply(null, Array(PD_MAX)).map(Number.prototype.valueOf,0);
            actionCmd.Payload.Amount[PD_PT_CEREALS]= rawAction.param.cereals
            actionCmd.Payload.Amount[PD_PT_MEAT]= rawAction.param.meat
            actionCmd.Payload.Amount[PD_PT_WOOL]= rawAction.param.wool
            actionCmd.Payload.Amount[PD_PT_BEER]= rawAction.param.beer
            actionCmd.Payload.Amount[PD_PT_SWEATER]= rawAction.param.sweater
        }else{
            actionCmd.Payload= rawAction.param
        }

        let actionID= actionMap.get(rawAction.action)
        actionCmd.ActionID= actionID

        let command = {
            ID: CMD_ACTION,
            Who: this.gd.clientID,
            Payload: actionCmd,
        }
        console.log("send action ", command)
        this.wsconn.send(JSON.stringify(command))
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
        if (cards == null) {
            return
        }
        console.log("cards length", cards.length)
        for (var i = 0; i < cards.length; i++) {
            let card = {
                id: cards[i].ID,
                token: [],
            }
            if (cards[i].Pawns != null) {
                for (var j = 0; j < cards[i].Pawns.length; j++) {
                    let color = this.gd.pd.get(cards[i].Pawns[j]).color
                    card.token.push(color)
                }
            }
            let tokens= 0
            if(cards[i].Pawns!= null){
                tokens= cards[i].Pawns.length
            }
            card.full = cards[i].MaxSlot - tokens < cards[i].PawnPerTurn
            this.gd.cards.push(card)
        }
        console.log("fwb cards", this.gd.cards)
    }

    getCardActions(cardid){
        return cardActionMap.get(cardid)
    }

    updatePlayerData(pd) {
        console.log("update player data", pd)
        let colors = ["red", "blue", "yellow", "green", "black", "brown", "orange"]
        let pdmap = this.gd.pd

        this.gd.pdlist = []
        this.gd.pd = new Map()
        if (pd == null) {
            return
        }

        for (var i = 0; i < pd.length; i++) {
            let clientID = pd[i][PD_CLIENT_ID]
            let p = {
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
                int: pd[i][PD_SK_INTELLIGENCE],
                hm: pd[i][PD_HOUSE_LV],
                token: pd[i][PD_PAWNS],
                name: pdmap.get(clientID).name,
            }
            this.gd.pd.set(clientID, p)
            this.gd.pdlist.push(p)
        }

        this.gd.pdlist.sort(function (a, b) {
            return a[PD_CLIENT_ID] - b[PD_CLIENT_ID]
        })

        let pdlist = this.gd.pdlist
        for (i = 0; i < pdlist.length; i++) {
            let p = this.gd.pdlist[i]
            p.color = colors[i]
            this.gd.pdlist[i] = p
            this.gd.pd.set(p.clientID, p)
        }

    }

    connError(event) {
        console.log(event)
    }

    fwbMessage(event) {
        console.log(event)
    }

    sendGameStartAck(skill) {
        let gsAck = {
            "ID": CMD_GAME_START_ACK,
            "Who": this.gd.clientID,
            "Payload": skill
        }
        console.log("game ack", gsAck)
        this.wsconn.send(JSON.stringify(gsAck))
    }

    sendRoundStartAck() {
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

        if (this.onActionCommitted != null) {
            this.onActionCommitted(cmd)
        }
    }

    __onActionRejected(cmd) {
        if (this.onActionRejected != null) {
            this.onActionRejected(cmd)
        }
    }

    __onGameFinish(cmd) {
        let rank = new Map(Object.entries(cmd.Payload))

        let playerList = []
        for (const [key, value] of rank) {
            let ikey = Number(key)
            let playerInfo = {
                name: this.gd.pd.get(ikey).name,
                rank: value,
            }
            playerList.push(playerInfo)
        }

        console.log(playerList)

        playerList.sort(function (a, b) {
            return a.rank - b.rank
        })

        console.log(playerList)

        let youwin = 0
        if ((rank.get(this.gd.clientID.toString()) === 1) && (playerList[playerList.length - 1].rank !== 1)) {
            youwin = 1
        } else if (rank.get(this.gd.clientID.toString()) !== 1) {
            youwin = -1
        }

        console.log("game finished: ", playerList, youwin)
        if (this.onGameFinished != null) {
            this.onGameFinished(cmd, playerList, youwin)
        }
    }

    getActionProperties(action, parameter) {
        let prop = {
            cost: {
                nocost: true,
                gold: 0,
                cereals: 0,
                meat: 0,
                wool: 0,
                sweater: 0,
                beer: 0,
            },
            gain: {
                nogain: true,
                gold: 0,
                cereals: 0,
                meat: 0,
                wool: 0,
                sweater: 0,
                beer: 0,
                inte: 0,
                know: 0,
                stre: 0,
                hm:0,
                token: 0
            }
        }

        let pd = this.gd.pd.get(this.gd.clientID)

        if (action === "feed-sheep") {
            prop.cost.nocost= false
            prop.cost.cereals= 2
        } else if (action === "take-vacation") {
            prop.cost.nocost= false
            prop.cost.gold= 5
        } else if (action === "train") {
            let sklv= 0
            prop.cost.nocost= false
            if(parameter=== "inte"){
                sklv= pd.int
            }else if(parameter=== "know"){
                sklv= pd.knw
            }else if(parameter==="stre"){
                sklv= pd.str
            }
            if (sklv === 0) {
                prop.cost.gold= 2
            } else if (sklv === 1) {
                prop.cost.gold= 5
            } else if (sklv === 2) {
                prop.cost.gold=10
            }
        }if (action === "upgrade-home") {
            prop.cost.nocost= false
            prop.cost.gold= 10
        } else if (action === "weave") {
            prop.cost.nocost= false
            prop.cost.gold= 5
            prop.cost.wool= parameter*2
        }else if (action=== "trade"){
            prop.cost.nocost= false
            prop.cost.gold= 2
            if(parameter.direction=== "buy"){
                prop.cost.gold+= parameter.cereals*2+ parameter.meat*5 + parameter.wool*5 + parameter.sweater*10 + parameter.beer* 20
            }else if(parameter.direction=== "sell"){
                prop.cost.cereals= parameter.cereals
                prop.cost.meat= parameter.meat
                prop.cost.wool= parameter.wool
                prop.cost.sweater= parameter.sweater
                prop.cost.beer= parameter.beer
            }
        }else if (action === "make-wine") {
            prop.cost.nocost= false
            prop.cost.cereals= parameter*8
        } else if (action === "party") {
            prop.cost.nocost= false
            prop.cost.beer= 2
            prop.cost.meat= 2
        } else if (action === "employ") {
            prop.cost.nocost= false
            prop.cost.gold= 10
        } else if (action === "skip") {
            prop.cost.nocost= true
        }

        if (action === "farm"){
            prop.gain.nogain= false
            if(pd.str=== 0){
                prop.gain.cereals= 2
            }else if(pd.str===1){
                prop.gain.cereals= 4
            }else if(pd.str===2){
                prop.gain.cereals= 6
            }else if(pd.str===3){
                prop.gain.cereals= 10
            }
            console.log("farm prop", prop, pd.str)

        }else if(action === "take-off"){
            prop.gain.nogain= false
            prop.gain.heart= 1
        }else if(action === "hunt"){
            prop.gain.nogain= false
            if(pd.str===0){
                prop.gain.meat= 1
            }else if(pd.str>2){
                prop.gain.meat= 4
            }else{
                prop.gain.meat= 2
            }
        }else if(action === "beg"){
            prop.gain.nogain= false
            if(pd.int>2){
                prop.gain.meat= 3
                prop.gain.gold= 5
            }else{
                prop.gain.cereals= 2
            }
        }else if(action==="parttime-work"){
            prop.gain.nogain= false
            if(pd.int<1){
                prop.gain.gold= 2
            }else if(pd.int===1){
                prop.gain.gold= 4
            }else if(pd.int=== 2){
                prop.gain.gold= 6
            }else if(pd.int=== 3){
                prop.gain.gold= 10
            }
        }else if (action === "gold-mining") {
            prop.gain.nogain= false
            if(pd.str===3 && pd.knw===3){
                prop.gain.gold= 20
            }else if(pd.str> 1 && pd.knw> 1){
                prop.gain.gold= 10
            }else{
                prop.gain.gold= 5
            }
        } else if (action === "feed-sheep") {
            prop.gain.nogain= false
            if(pd.knw===0){
                prop.gain.meat= 1
                prop.gain.wool= 1
            }else if(pd.knw===1){
                prop.gain.meat= 2
                prop.gain.wool= 1
            }else if(pd.knw===2){
                prop.gain.meat= 2
                prop.gain.wool= 2
            }else if(pd.knw>2){
                prop.gain.meat= 3
                prop.gain.wool= 3
            }
        } else if (action === "take-vacation") {
            prop.gain.nogain= false
            prop.gain.heart= 2
        } else if (action === "train") {
            prop.gain.nogain= false
            prop.gain[parameter]= 1
        } else if (action === "steal") {
            prop.gain.nogain= false
            let opd = null
            if(typeof parameter!= undefined){
                opd= this.findPlayerData(parameter)
            }
            const lvdelta= pd.int- opd.knw
            if(lvdelta===3){
                prop.gain.gold= 30
            }else if(lvdelta===2){
                prop.gain.gold= 20
            }else if(lvdelta===1){
                prop.gain.gold= 10
            }
        } else if (action === "trade") {
            if(parameter.direction==="buy"){
                prop.gain.cereals= parameter.cereals
                prop.gain.meat= parameter.meat
                prop.gain.wool= parameter.wool
                prop.gain.sweater= parameter.sweater
                prop.gain.beer= parameter.beer
                if(parameter.cereals+ parameter.meat+ parameter.wool+ parameter.sweater+ parameter.beer> 0){
                    prop.gain.nogain= false
                }
            }else if(parameter.direction==="sell"){
                prop.gain.gold= parameter.cereals*2 + parameter.meat* 5 + parameter.wool* 5 + parameter.sweater* 10 + parameter.beer* 20
                if(prop.gain.gold-2>0){
                    prop.gain.nogain= false
                }
            }
        } else if (action === "upgrade-home") {
            prop.gain.nogain= false
            prop.gain.hm= 1
        } else if (action === "weave") {
            prop.gain.nogain= false
            prop.gain.sweater= parameter/2
        } else if (action === "make-wine") {
            prop.gain.nogain= false
            prop.gain.beer= parameter/8
        } else if (action === "party") {
            prop.gain.nogain= false
            prop.gain.heart= 10
        } else if (action === "employ") {
            prop.gain.nogain= false
            prop.gain.token= 3
        } else if (action === "skip") {
            prop.gain.nogain= true
        }

        return prop;
    }

    findPlayerData(name) {
        for (const entry of this.gd.pd) {
            if (entry[1].name === name) {
                return entry[1]
            }
        }
        return null
    }

    getRound(){
        return this.gd.round
    }

    shouldEnableCard(card){
        if(card.full){
            return false
        }

        const actions= cardActionMap.get(card.id)
        for(let i in actions){
            if(this.isActionApplicable(actions[i])){
                return true
            }
        }

        return false
    }

    isActionApplicable(action, parameter) {
        let pd = this.gd.pd.get(this.gd.clientID)
        let opd = null
        if(typeof parameter!= undefined){
            opd= this.findPlayerData(parameter)
        }

        if (action === "farm" ||
            action === "take-off" ||
            action === "hunt" ||
            action === "beg" ||
            action === "parttime-work"
        ) {
            return pd.token >= 1
        } else if (action === "gold-mining") {
            return pd.knw >= 1 && pd.str >= 1 && pd.token >= 1
        } else if (action === "feed-sheep") {
            return pd.cereals >= 2 && pd.token >= 1
        } else if (action === "take-vacation") {
            return pd.gold >= 5 && pd.token >= 1
        } else if (action === "train") {
            const minsk = Math.min(pd.int, pd.knw, pd.str)
            if (minsk === 0) {
                return pd.gold >= 2 && pd.token >= 1
            } else if (minsk === 1) {
                return pd.gold >= 5 && pd.token >= 1
            } else if (minsk === 2) {
                return pd.gold >= 10 && pd.token >= 1
            }
            return false
        } else if (action === "steal") {
            if(opd=== null){
                return pd.token>= 1
            }else{
                return  pd.token >= 1 && pd.int > opd.knw
            }
        } else if (action === "trade") {
            return pd.gold >= 2 && pd.token >= 1
        } else if (action === "upgrade-home") {
            if(pd.hm=== 1){
                if(pd.int<1 || pd.knw<1 || pd.str<1){
                    return false
                }
            }else if(pd.hm=== 2){
                if(pd.int< 2 || pd.knw< 2 || pd.str< 2){
                    return false
                }
            }else if(pd.hm>2){
                return false
            }

            return pd.gold >= 10 && pd.token >= 3
        } else if (action === "weave") {
            return pd.wool >= 2 && pd.gold >= 5
        } else if (action === "make-wine") {
            return pd.cereals >= 8 && pd.token >= 1
        } else if (action === "party") {
            return pd.beer >= 2 && pd.meat >= 2 && pd.token >= 1
        } else if (action === "employ") {
            return pd.gold >= 10 && pd.token >= 1
        } else if (action === "skip") {
            return true
        }
        return false
    }

    __onGameOver(cmd) {
        if (this.onGameOver != null) {
            this.onGameOver(cmd)
        }
    }

    __onRoundSettlement(cmd) {
        if (this.onRoundSettlement != null) {
            this.onRoundSettlement(cmd)
        }
    }

    __onRoundSettlementInvalid(cmd) {
        if (this.onRoundSettlementInvalid != null) {
            this.onRoundSettlementInvalid(cmd)
        }
    }

    __onRoundSettlementUpdate(cmd) {
        this.updateGameData(cmd.Payload)

        if (this.onRoundSettlementUpdate != null) {
            this.onRoundSettlementUpdate(cmd)
        }
    }

    __onRoundStart(cmd) {
        this.updateGameData(cmd.Payload)

        if (this.onRoundStart != null) {
            if (this.onRoundStart(cmd)) {
                this.sendRoundStartAck()
            }
        }
    }

    __onStartTurn(cmd) {
        let myturn = cmd.Payload === this.gd.clientID
        let whoseturn = this.gd.pd.get(cmd.Payload).name
        console.log("start turn", this.gd.pd, cmd.Payload, this.gd.pd.get(cmd.Payload))
        console.log("start turn ", this.gd.pd, this.gd.clientID, this.gd.pd.get(this.gd.clientID))
        this.gd.currentPlayer = cmd.Payload

        if (this.onStartTurn != null) {
            this.onStartTurn(cmd, myturn, whoseturn)
        }
    }

    __onSyncGameState(cmd) {
        this.updateGameData(cmd.Payload)

        if (this.onSyncGameState != null) {
            this.onSyncGameState(cmd)
        }
    }

    __onGameStart(cmd) {
        console.log(cmd)
        this.gd.pd = new Map()
        this.gd.pdlist = []
        let playerInfo = new Map(Object.entries(cmd.Payload))
        console.log(playerInfo)
        let keys = [...playerInfo.keys()]
        for (var i = 0; i < keys.length; i++) {
            let clientID = playerInfo.get(keys[i])
            let pn = keys[i]
            let p = {
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
            this.sendGameStartAck("stre")
        }
    }

    __onDefaultCmd(cmd) {
        console.log("No handler for command: ", cmd)
    }

    __onWSMessage(event) {

        console.log(event.data)
        var cmd = JSON.parse(event.data)

        console.log("received ", cmd)

        if (!this.cmdMap.has(cmd.ID)) {
            this.__onDefaultCmd(cmd)
        } else {
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
        this.gd.username = username
        this.retry = 0
        this.doConnect(username)
    }
}

export default new FWB()