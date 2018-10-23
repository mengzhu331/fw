import React, { Component } from 'react';
import SignIn from './signin'
import fwb from './fwb'
import GameMain from './gamemain'
import './App.css';

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
    return true
  }

  componentDidMount(){
    this.setState({
      screen: "sign-in"
    })

    fwb.onGameStart= this.onGameStart.bind(this)    
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
