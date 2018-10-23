import React, { Component } from 'react';
import fwb from './fwb'
import DlgInfo from './dlg-info'
import Round from "./round"
import Timer from "./timer"
import Card from './card'
import './App.css';
import PlayerCard from './player-card';

class GameMain extends Component{
    constructor(props){
        super(props)
        fwb.onSyncGameState= this.onSyncGameState.bind(this)
        fwb.onRoundStart= this.onRoundStart.bind(this)
        fwb.onStartTurn= this.onStartTurn.bind(this)
        this.timer= 0;
        this.state= {
            displayOwnTurn:false,
            displayOthersTurn:false,
        }
    }

    onSyncGameState(cmd){
        console.log("sync game state")
        this.forceUpdate()
    }

    onRoundStart(cmd){
        console.log("round start")
        this.forceUpdate()
        return true
    }

    onStartTurn(cmd, ownTurn, whoseTurn){
        if(ownTurn){
            this.setState({
                displayOwnTurn:true,
                displayOthersTurn:false,
                whoseTurn: fwb.gd.username,
            })

            setTimeout(function(){
                this.setState({
                    displayOwnTurn:false,
                })
            }.bind(this), 2000)
        }else{
            this.setState({
                displayOwnTurn: false,
                displayOthersTurn: true,
                whoseTurn: whoseTurn,
            })
        }

        this.setState({
            secondToSkip: 60,
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

    render(){

        let playerCards=[]
        let pcdata= fwb.gd.pdlist

        for(var i=0; i< pcdata.length; i++){
            let pcard= <PlayerCard key={i} pd={pcdata[i]}/>
            playerCards.push(pcard)
        }

        let cards= []
        let cdata= fwb.gd.cards
        for(i=0; i< cdata.length; i++){
            let card= <Card key={i} card={cdata[i]}/>
            cards.push(card)
        }

        let round= fwb.gd.round

        let info= null;
        if(this.state.displayOwnTurn){
            info=  <DlgInfo header="回合进行中" content="准备行动"></DlgInfo>
        }else if(this.state.displayOthersTurn){
            let name= this.state.whoseTurn;  
            let content="等待"+ name+ "的行动"
            info=  <DlgInfo header="回合进行中" content={content}></DlgInfo>
        }

        return (<div className="game-main-screen">
            <Round round={round} />
            <Timer secondToSkip={this.state.secondToSkip}/>
            {info}
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
        </div>)
    }
}

export default GameMain