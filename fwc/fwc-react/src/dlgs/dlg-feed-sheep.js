import React, { Component } from 'react';
import '../App.css';
import fwb from '../fwb'
import {IconCereals, IconWool, IconMeat} from './icons'

class FeedSheepDialog extends Component{
    onOk(){
        this.props.onOk("feed-sheep")
    }

    render(){
        let prop= fwb.getActionProperties("feed-sheep")

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">养羊</p> 从而获得羊毛和肉吗？</div>
            <div>您将支付:</div>
            <div className= "highlight-info">{prop.cost.cereals}<IconCereals/></div>
            <div>您将获得：</div>
            <div className= "highlight-info">{prop.gain.meat}<IconMeat/>和{prop.gain.wool}<IconWool/></div>
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

export default FeedSheepDialog