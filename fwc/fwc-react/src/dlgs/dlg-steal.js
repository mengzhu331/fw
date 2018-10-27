import React, { Component } from 'react';
import '../App.css';
import {IconCoin} from './icons'
import fwb from '../fwb'

class StealDialog extends Component{
    onClick(event){
        this.props.onOk("steal", parseInt(event.currentTarget.id))
    }


    render(){

        let inte= fwb.gd.pd.get(fwb.gd.clientID).int

        let playerButton=<div>没有可以偷窃的玩家</div>
        let playerbuttons= []
        for(let e of fwb.gd.pd){
            console.log("try steal", fwb.gd.clientID, e, inte)
            if(fwb.gd.clientID!== e[0] && inte> e[1].knw){
                let prop= fwb.getActionProperties("steal", e[1].name)
                playerbuttons.push(<div className= "button-option" id={e[0]} onClick= {this.onClick.bind(this)}>
                {e[1].name}可获得{prop.gain.gold}<IconCoin/>
                </div>)
            }
        }
        if (playerbuttons.length>0){
            playerButton= playerbuttons
        }

        return <div className="topdown">
            <div>您要偷窃哪位玩家的金币?（智力必须超过对方的知识）</div>
            {playerButton}
            <div className= "button-cancel" onClick={this.props.onCancel}>再想想</div>
        </div>
    }
}

export default StealDialog