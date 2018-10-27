import React, { Component } from 'react';
import './App.css';
import fwb from './fwb'

class Round extends Component{
    render(){
        let round= fwb.getRound()

        return <div className="round">
            <div></div>
            <div className= "round-text">第</div>
            <div className="round-label">{round}</div>
            <div className= "round-text">回合</div>
            <div></div>

        </div>
    }
}

export default Round