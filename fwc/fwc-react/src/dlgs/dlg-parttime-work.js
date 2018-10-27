import React, { Component } from 'react';
import '../App.css';
import {IconCoin} from './icons'
import fwb from '../fwb'

class ParttimeWorkDialog extends Component{
    onOk(){
        this.props.onOk("parttime-work")
    }

    render(){
        let prop= fwb.getActionProperties("parttime-work")
        console.log("prop output", prop)

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">打工</p> 从而获得金币吗？</div>
            <div>您将获得：</div>
            <div className= "highlight-info">{prop.gain.gold}<IconCoin/></div>
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

export default ParttimeWorkDialog