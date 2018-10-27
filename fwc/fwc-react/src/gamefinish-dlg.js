import React, { Component } from 'react';
import fwb from './fwb'
import victory from './img/victory.png'
import defeat from './img/defeat.png'
import draw from './img/draw.jpg'

class GameFinishDialog extends Component{
        render(){
        console.log(this.props)

        let icon=<img alt="victory" src={victory} className="icon-s"/>
        if(this.props.youwin=== -1){
            icon=<img alt="defeat" src={defeat} className="icon-s"/>
        }else if(this.props.youwin=== 0){
            icon=<img alt="draw" src={draw} className="icon-s"/>
        }

        let list=[]
        for(let i=0; i< this.props.playerlist.length; i++){
            let rank= <div key={i}>{this.props.playerlist[i].rank} - {this.props.playerlist[i].name}</div>
            list.push(rank)
        }
        return (<div className="gamefinish-screen">
        {icon}
        {list}
        </div>)
    }
}

export default GameFinishDialog