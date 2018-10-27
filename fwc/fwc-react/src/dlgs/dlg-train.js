import React, { Component } from 'react';
import '../App.css';
import {IconCoin, IconInte, IconKnow, IconStre} from './icons'
import fwb from '../fwb'

class TrainDialog extends Component{
    onClick(event){
        if(event.target.id==="inte"){
            this.props.onOk("train", 10)
        }
        if(event.target.id==="know"){
            this.props.onOk("train", 9)
        }
        if(event.target.id==="stre"){
            this.props.onOk("train", 8)
        }
    }


    render(){

        let skillButton=<div>没有适合升级的技能</div>
        let skillbuttons= []
        let pd= fwb.getPd()

        let prop= fwb.getActionProperties("train", "inte")

        if(pd.int<3 && pd.gold>= prop.cost.gold){
            skillbuttons.push(<div className="button-option" id="inte" onClick={this.onClick.bind(this)}>智力<IconInte/>花费{prop.cost.gold}<IconCoin/></div>)
        }

        prop= fwb.getActionProperties("train", "know")

        if(pd.knw<3 && pd.gold>= prop.cost.gold){
            skillbuttons.push(<div className="button-option" id="know" onClick={this.onClick.bind(this)}>知识<IconKnow/>花费{prop.cost.gold}<IconCoin/></div>)
        }

        prop= fwb.getActionProperties("train", "stre")

        if(pd.str<3 && pd.gold>= prop.cost.gold){
            skillbuttons.push(<div className="button-option" id="stre" onClick={this.onClick.bind(this)}>体力<IconStre/>花费{prop.cost.gold}<IconCoin/></div>)
        }

        if(skillbuttons.length>0){
            skillButton= skillbuttons
        }
        return <div className="topdown">
            <div>您想要升级什么技能？（只能选择不到3级并且能付足够金币的技能）</div>
            {skillButton}
            <div className= "button-cancel" onClick={this.props.onCancel}>再想想</div>
        </div>
    }
}

export default TrainDialog