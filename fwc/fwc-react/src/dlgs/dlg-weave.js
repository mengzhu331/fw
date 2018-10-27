import React, { Component } from 'react';
import '../App.css';
import {IconWool, IconSweater, IconCoin} from './icons'
import fwb from '../fwb'
import plus from '../img/plus.png'
import minus from '../img/minus.png'

class WeaveDialog extends Component{
    constructor(props){
        super(props)
        this.amount=0
    }
    
    onOk(){
        this.props.onOk("weave", this.amount)
    }

    plus(){
        let pd= fwb.getPd()
        let prop= fwb.getActionProperties("weave", this.amount+1)
        if(prop.cost.wool<= pd.wool){
            this.amount++
        }
        this.forceUpdate()
    }

    minus(){
        if(this.amount> 0){
            this.amount--
        }
        this.forceUpdate()
    }

    render(){
        let prop= fwb.getActionProperties("weave", this.amount)
        console.log("prop output", prop)

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">织布</p> 从而获得毛衣吗？</div>
            <div>您想制作:</div>
            <div className= "input-tool">
                <IconSweater/>
                <img alt="plus" src={plus} className="icon-dlg" onClick={this.plus.bind(this)}/>
                {this.amount}
                <img alt="plus" src={minus} className="icon-dlg" onClick={this.minus.bind(this)}/>
            </div>

            <div>您需要支付：</div>
            <div className= "highlight-info">{prop.cost.wool}<IconWool/></div>
            <div className= "highlight-info">{prop.cost.gold}<IconCoin/></div>
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

export default WeaveDialog