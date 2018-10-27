import React, { Component } from 'react';
import SignIn from './signin'
import fwb from './fwb'
import GameMain from './gamemain'
import './App.css';
import fwbg from './sound/fwbg.mp3'

import WeaveDialog from './dlgs/dlg-weave'
import MakeWineDialog from './dlgs/dlg-make-wine';
import PartyDialog from './dlgs/dlg-party';
import GoldMiningDialog from './dlgs/dlg-gold-mining';
import TradeDialog from './dlgs/dlg-trade'
import SettlementDialog from './dlg-settlement';

class App extends Component {

  constructor(props){
    super(props)
    this.state={
      screen: ""
    }
  }

  onGameStart(cmd){
    this.setState({
      screen: "game-main"
    })
    return false
  }

  onGameOver(cmd){
    this.setState({
      screen:"sign-in"
    })
  }

  componentDidMount(){
    let myAudio = new Audio(fwbg); 
    myAudio.addEventListener('ended', function() {
      this.currentTime = 0;
      this.play();
    }, false);
    //myAudio.play();
    this.setState({
      screen: "sign-in"
    })

    fwb.onGameStart= this.onGameStart.bind(this)
    fwb.onGameOver= this.onGameOver.bind(this)
  }

  render() {
    let content= ""
    
    if(this.state.screen==="game-main"){
      content= <GameMain></GameMain>
    }else if (this.state.screen=== "sign-in"){
      content= <SignIn></SignIn>
    }

    //content= <GameMain></GameMain>

    return (
      <div className="App">
        {content}
      </div>
    );
  }
}

export default App;
