import React, { Component } from 'react';
import fwb from './fwb'
import './App.css';

class DlgInfo extends Component{
    
    render(){
        var buttons=[]
        if (this.props.buttons!= null){
            for(var i=0; i< this.props.buttons.length; i++){
                var key= "button"+i
                buttons.push(<button key={key}>{this.props.buttons[i]}</button>)
            }
        }

        var contenticon=""

        if(this.props.icon=== "wait"){
            contenticon= <img className= "dlg-info-icon" alt="wait" src={require("./img/wait.gif")}></img>
        }

        return <div className="dlg-info">
        <div className="dlg-info-header">{this.props.header}</div>
        <div className="dlg-info-content">
        {contenticon}
        {this.props.content}
        </div>
        <div className="dlg-info-buttons">{buttons}</div>       
        </div>
    }
}

export default DlgInfo