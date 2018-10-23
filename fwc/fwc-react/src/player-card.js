import React, { Component } from 'react';
import './App.css';

class PlayerCard extends Component{
    render(){
        let icnStr= <img className="icon" alt="str" src={require("./img/stre1.png")}/>
        let icnKlg= <img className="icon" alt="str" src={require("./img/know1.png")}/>
        let icnInt= <img className="icon" alt="str" src={require("./img/wisd1.png")}/>
        let icnHom= <img className="icon" alt="str" src={require("./img/house1.png")}/>

        let icnCrl= <img className="icon" alt="str" src={require("./img/cereals.png")}/>
        let icnMea= <img className="icon" alt="str" src={require("./img/meat.png")}/>
        let icnWoo= <img className="icon" alt="str" src={require("./img/wool.png")}/>
        let icnSwt= <img className="icon" alt="str" src={require("./img/sweater.png")}/>
        let icnBee= <img className="icon" alt="str" src={require("./img/beer.png")}/>

        let color= this.props.pd.color

        let token= []
        for(var i=0; i< this.props.pd.token; i++){
            let icnTkn= <img className="iconToken" alt="token" key={i} src={require("./img/token"+color+".png")}/>
            token.push(icnTkn)
        }

        return(
        <div className="player-card">
            <div className= "pc-username">{this.props.pd.name}</div>
            <div className= "pc-heart">{this.props.pd.heart}</div>
            <div className= "pc-gold">{this.props.pd.gold}</div>

            <div className= "pc-strength">{icnStr}{this.props.pd.str}</div>
            <div className= "pc-knowledge">{icnKlg}{this.props.pd.knw}</div>
            <div className= "pc-intelligence">{icnInt}{this.props.pd.int}</div>
            <div className= "pc-home-level">{icnHom}{this.props.pd.hm}</div>

            <div className= "pc-cereals">{icnCrl}{this.props.pd.cereals}</div>
            <div className= "pc-meat">{icnMea}{this.props.pd.meat}</div>
            <div className= "pc-wool">{icnWoo}{this.props.pd.wool}</div>
            <div className= "pc-sweater">{icnSwt}{this.props.pd.sweater}</div>
            <div className= "pc-beer">{icnBee}{this.props.pd.beer}</div>
            <div className= "pc-token">{token}</div>
        </div>
        )
    }
}

export default PlayerCard