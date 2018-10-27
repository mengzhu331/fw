import React, { Component } from 'react';
import fwb from './fwb'
import DlgInfo from './dlg-info'
import Round from "./round"
import Timer from "./timer"
import Card from './card'
import './App.css';
import PlayerCard from './player-card';
import ChooseSkill from './chooseskill';
import CardDialog from './dlg-card';
import SettlementDialog from './dlg-settlement';
import GameFinishDialog from './gamefinish-dlg';
import yourturn from './sound/yourturn.wav'
import actionsfx from './sound/action.wav'
import victory from './sound/victory.mp3'
import defeat from './sound/defeat.mp3'

class GameMain extends Component{
    constructor(props){
        super(props)
        fwb.onSyncGameState= this.onSyncGameState.bind(this)
        fwb.onRoundStart= this.onRoundStart.bind(this)
        fwb.onStartTurn= this.onStartTurn.bind(this)
        fwb.onRoundSettlement= this.onRoundSettlement.bind(this)
        fwb.onGameFinished= this.onGameFinished.bind(this)
        this.timer= 0;
        this.state= {
            displayOwnTurn:false,
            displayOthersTurn:false,
            cardsEnabled:true,
            showCardDialog:false,
            displaySettlement: false,
            waitSettlement: false,
            displayGameFinish:false,
        }
    }

    componentDidMount(){
        this.yourturn_sfx = new Audio(yourturn); 
        this.action_sfx = new Audio(actionsfx); 
        this.victory_sfx = new Audio(victory); 
        this.defeat_sfx= new Audio(defeat)

        this.setState({
            displayChooseSkill: true,
            secondToSkip: 10,
        })

        if(this.timer!== 0){
            clearInterval(this.timer)
        }
        this.timer= setInterval(function(){
            let s= this.state.secondToSkip-1;
            if(s<0){
                s= 0;
            }
            this.setState({
                secondToSkip: s,
            })
        }.bind(this), 1000)
    }

    onSyncGameState(cmd){
        console.log("sync game state")
        this.forceUpdate()
    }

    onRoundSettlement(){
        this.setState({
            displaySettlement: true,
            waitSettlement: false,
            secondToSkip:30,
        })
    }

    onClickSettlement(){
        this.setState({
            displaySettlement: false,
            waitSettlement:true,
        })
    }

    onRoundStart(cmd){
        console.log("round start")
        this.setState({
            displayChooseSkill:false,
            displaySettlement: false,
            waitSettlement: false,
            displayGameFinish:false,
        })
        return true
    }

    onStartTurn(cmd, ownTurn, whoseTurn){
        if(ownTurn){
            //this.yourturn_sfx.play()
            this.setState({
                displayOwnTurn:true,
                displayOthersTurn:false,
                whoseTurn: fwb.gd.username,
                cardsEnabled:true,
                showCardDialog:false,
            })

            setTimeout(function(){
                this.setState({
                    displayOwnTurn:false,
                })
            }.bind(this), 1000)
        }else{
            this.setState({
                displayOwnTurn: false,
                displayOthersTurn: true,
                whoseTurn: whoseTurn,
                cardsEnabled:false,
                showCardDialog:false,
            })
        }

        this.setState({
            secondToSkip: 60,
        })
    }    

    onGameFinished(cmd, result, win){
        if(win>0){
            //this.victory_sfx.play()
        }else if(win<0){
            //this.defeat_sfx.play()
        }
        console.log(cmd, result, win)
        this.setState({
            displayOwnTurn:false,
            displayOthersTurn:false,
            cardsEnabled:false,
            showCardDialog:false,
            waitSettlement:false,
            displaySettlement:false,
            displayGameFinish: true,
            ranklist:result,
            win:win,
            secondToSkip: 10,
        })
    }

    waitClickHandler(card){
        console.log("clicked card ", card)
    }

    cardClickHandler(card){
        console.log("clicked card ", card)
        this.setState({
            cardsEnabled: false,
            showCardDialog: true,
            clickedCard: card,
        })
    }

    onCardDialogOk(card){
        console.log("clicked ok, card ", card)
        //this.action_sfx.play()
        this.setState({
            cardsEnabled: false,
            showCardDialog: false,
        })
    }

    onCardDialogCancel(card){
        console.log("clicked cancel, card ", card)
        this.setState({
            cardsEnabled: true,
            showCardDialog: false,
            clickedCard: null,
        })
    }

    render(){

        let playerCards=[]
        let pcdata= fwb.gd.pdlist

        for(var i=0; i< pcdata.length; i++){
            let pcard= <PlayerCard key={i} pd={pcdata[i]}/>
            playerCards.push(pcard)
        }

        let cards= []
        let cdata= fwb.gd.cards

        let clickHandler=this.waitClickHandler.bind(this)
        if(this.state.cardsEnabled){
            clickHandler= this.cardClickHandler.bind(this)
        }

        let cardDialog=null
        if(this.state.showCardDialog){
            cardDialog=<CardDialog card={this.state.clickedCard} onOk={this.onCardDialogOk.bind(this)} onCancel={this.onCardDialogCancel.bind(this)}/>
        }

        for(i=0; i< cdata.length; i++){
            let card= <Card key={i} card={cdata[i]} clickHandler={clickHandler}/>
            cards.push(card)
        }

        let info= null;
        if(this.state.displayOwnTurn){
            info=  <DlgInfo header="回合进行中" content="准备行动"></DlgInfo>
        }else if(this.state.displayOthersTurn){
            let name= this.state.whoseTurn;  
            let content="等待"+ name+ "的行动"
            info=  <DlgInfo header="回合进行中" content={content}></DlgInfo>
        }

        let chooseSkill= null;
        if(this.state.displayChooseSkill){
            chooseSkill=<ChooseSkill/>
        }

        let settlement= null;
        if(this.state.displaySettlement){
            settlement=<SettlementDialog onOk={this.onClickSettlement.bind(this)}/>
        }
        if(this.state.waitSettlement){
            settlement=<DlgInfo header="已分配好" content="等待其他玩家分配产品"/>
        }

        let gamefinishdialog= null
        if(this.state.displayGameFinish){
            gamefinishdialog=<GameFinishDialog playerlist={this.state.ranklist} youwin={this.state.win}/>
        }

        return (<div className="game-main-screen">
            <Round/>
            <div className="player-data-row">
                <div className="player-data-bg">
                    <div className= "player-data">
                        {playerCards}
                    </div>
                </div>
            </div>
            <div className="board-bg">
            <div className="board">
                {cards}
            </div>
            </div>
            {cardDialog}
            <Timer secondToSkip={this.state.secondToSkip}/>
            {info}
            {chooseSkill}
            {settlement}
            {gamefinishdialog}
        </div>)
    }
}

export default GameMain