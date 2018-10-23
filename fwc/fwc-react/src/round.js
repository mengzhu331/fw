import React, { Component } from 'react';
import './App.css';

class Round extends Component{
    render(){
        return <div className="round">
            <div></div>
            <div className= "round-text">第</div>
            <div className="round-label">{this.props.round}</div>
            <div className= "round-text">回合</div>
            <div></div>

        </div>
    }
}

export default Round