import React, { Component } from 'react';
import './App.css';

class Card extends Component{
    render(){
        let id= this.props.card.id
        if(id<1){
            return null
        }

        let bg= require("./img/card"+ id+ ".png")

        let token= this.props.card.token
        
        let icnTkn= []

        for(var i=0; i< token.length; i++){
            icnTkn.push(<img className="iconToken" alt="token" src={require("./img/token"+token[i]+".png")}/>)
        }

        let style= {
            backgroundImage: "url("+ bg+")",
        }
        return <div className="card" style={ style }>{icnTkn}</div>
    }
}

export default Card