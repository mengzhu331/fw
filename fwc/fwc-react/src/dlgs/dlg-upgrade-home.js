import React, { Component } from 'react';
import '../App.css';
import fwb from '../fwb'
import {IconCoin, IconHome} from './icons'

class UpgradeHomeDialog extends Component{
    onOk(){
        this.props.onOk("upgrade-home")
    }

    render(){
        let prop= fwb.getActionProperties("upgrade-home")
        console.log(prop)

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">升级房屋</p> 从而每回合获得红心吗？</div>
            <div>您将支付:</div>
            <div className= "highlight-info">{prop.cost.gold}<IconCoin/></div>
            <div>您将获得：</div>
            <div className= "highlight-info"><IconHome/>+{prop.gain.hm}</div>
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

export default UpgradeHomeDialog