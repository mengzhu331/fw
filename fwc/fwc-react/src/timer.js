import React, { Component } from 'react';
import './App.css';

class Timer extends Component{
 
    render(){
       return <div className="timer">{this.props.secondToSkip}</div>
    }
}

export default Timer