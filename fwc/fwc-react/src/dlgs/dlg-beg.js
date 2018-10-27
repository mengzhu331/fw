import React, { Component } from 'react';
import '../App.css';
import {IconMeat, IconCereals, IconCoin} from './icons'

import fwb from '../fwb'

class BegDialog extends Component{

    onOk(){
        this.props.onOk("beg")
    }

    render(){
        let prop= fwb.getActionProperties("beg")
        let meatAmount= null
        let meat=null
        let goldAmount= null
        let gold= null
        let cerealsAmount= null
        let cereals= null
        if(prop.gain.meat!== 0){
            meatAmount= prop.gain.meat
            meat= <IconMeat/>
        }

        if(prop.gain.gold!== 0){
            goldAmount= prop.gain.gold
            gold= <IconCoin/>
        }

        if(prop.gain.cereals!== 0){
            cerealsAmount=prop.gain.cereals
            cereals= <IconCereals/>
        }

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">乞讨</p> 从而获得物资吗？</div>
            <div>您将获得：</div>
            <div className= "highlight-info">{meatAmount}{meat}{cerealsAmount}{cereals}{goldAmount}{gold}</div>
            <div className= "tool-bar">
                <div/>
                <div className= "button-cancel" onClick={this.props.onCancel}>再想想</div>
                <div/>
                <div className= "button-ok" onClick={this.onOk.bind(this)}>好</div>
                <div/>
            </div>
        </div>
    }
}

export default BegDialog