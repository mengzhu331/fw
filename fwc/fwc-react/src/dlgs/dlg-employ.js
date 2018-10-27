import React, { Component } from 'react';
import '../App.css';
import {IconCoin, IconToken} from './icons'
import fwb from '../fwb'

class EmployDialog extends Component{
    onOk(){
        this.props.onOk("employ")
    }

    render(){
        let prop= fwb.getActionProperties("employ")
        console.log("prop output", prop)

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">雇工</p> 通过支付金币临时获得人手吗？</div>
            <div>您将支付:</div>
            <div className= "highlight-info">{prop.cost.gold}<IconCoin/></div>
            <div>您将获得：</div>
            <div className= "highlight-info">{prop.gain.token}<IconToken/></div>
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

export default EmployDialog