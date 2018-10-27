import React, { Component } from 'react';
import './App.css';
import fwb from './fwb';

class Card extends Component{
    onClick(event){
        if(!fwb.shouldEnableCard(this.props.card)){
            return
        }
        this.props.clickHandler(this.props.card)
    }

    render(){
        let id= this.props.card.id
        if(id<1){
            return null
        }

        let bg= require("./img/card"+ id+ ".png")

        let token= this.props.card.token
        
        let icnTkn= []

        for(var i=0; i< token.length; i++){
            icnTkn.push(<img className="iconToken" alt="token" key={i} src={require("./img/token"+token[i]+".png")}/>)
        }

        let style= {
            backgroundImage: "url("+ bg+")",
        }

        let cardClass= "card-disabled"
        if(fwb.shouldEnableCard(this.props.card)){
            cardClass="card"
        }
        return <div className={cardClass} style={ style } onClick={this.onClick.bind(this)}>{icnTkn}</div>
    }
}

export default Card