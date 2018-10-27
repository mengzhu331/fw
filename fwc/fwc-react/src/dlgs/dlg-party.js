import React, { Component } from 'react';
import '../App.css';
import fwb from '../fwb'
import {IconBeer, IconMeat, IconHeart} from './icons'

class PartyDialog extends Component{
    onOk(){
        this.props.onOk("party")
    }

    render(){
        let prop= fwb.getActionProperties("party")

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">宴会</p> 从而获得红心吗？</div>
            <div>您将支付:</div>
            <div className= "highlight-info">{prop.cost.meat}<IconMeat/></div>
            <div className= "highlight-info">{prop.cost.beer}<IconBeer/></div>
            <div>您将获得：</div>
            <div className= "highlight-info">{prop.gain.heart}<IconHeart/></div>
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

export default PartyDialog