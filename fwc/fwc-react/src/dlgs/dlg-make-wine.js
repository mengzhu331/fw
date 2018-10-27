import React, { Component } from 'react';
import '../App.css';
import {IconBeer, IconCereals} from './icons'
import fwb from '../fwb'
import plus from '../img/plus.png'
import minus from '../img/minus.png'

class MakeWineDialog extends Component{
    constructor(props){
        super(props)
        this.amount=0
    }
    
    onOk(){
        this.props.onOk("make-wine", this.amount)
    }

    plus(){
        let pd= fwb.getPd()
        pd.cereals= 100
        let prop= fwb.getActionProperties("make-wine", this.amount+1)
        if(prop.cost.cereals<= pd.cereals){
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
        let prop= fwb.getActionProperties("make-wine", this.amount)
        console.log("prop output", prop)

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">酿酒</p> 从而获得啤酒吗？</div>
            <div>您想制作:</div>
            <div className= "input-tool">
                <IconBeer/>
                <img alt="plus" src={plus} className="icon-dlg" onClick={this.plus.bind(this)}/>
                {this.amount}
                <img alt="plus" src={minus} className="icon-dlg" onClick={this.minus.bind(this)}/>
            </div>

            <div>您需要支付：</div>
            <div className= "highlight-info">{prop.cost.cereals}<IconCereals/></div>
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

export default MakeWineDialog