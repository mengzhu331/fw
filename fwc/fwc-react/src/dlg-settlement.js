import React, { Component } from 'react';
import fwb from './fwb'
import DlgInfo from './dlg-info'
import Timer from "./timer"
import './App.css';
import plus from './img/plus.png'
import minus from './img/minus.png'
import { IconCereals, IconMeat, IconSweater } from './dlgs/icons';

class SettlementDialog extends Component{
    constructor(props){
        super(props)
        this.order={
            cereals: 0,
            meat: 0,
            sweater: 0,
        }
    }

    plus(event){
        let pd= fwb.getPd()

        if(this.order[event.target.id]>= pd[event.target.id]){
            return
        }

        if(event.target.id==="cereals" || event.target.id==="meat"){
            if(this.order.cereals+ this.order.meat >=5){
                if(event.target.id==="cereals"){
                    if(this.order.cereals>= 5){
                        return
                    }
                    this.order.meat--
                }else{
                    if(this.order.meat>= 5){
                        return
                    }
                    this.order.cereals--
                }
            }
        }
        if(this.order[event.target.id]>= 5){
            return
        }
        this.order[event.target.id]++
        this.forceUpdate()
    }

    onOk(){
        fwb.sendSettlement(this.order)
        this.props.onOk()
    }

    minus(event){
        if(this.order[event.target.id]>0){
            this.order[event.target.id]--
        }
        this.forceUpdate()
    }

    render(){

        return (<div className="settlement-dlg">
            <div>请分配产品从而获得红心</div>
            <div className="button-option">
                <img alt="plus" src={plus} id="cereals" className="icon-dlg" onClick={this.plus.bind(this)}/>
                {this.order.cereals}
                <img alt="minus" src={minus} id="cereals" className="icon-dlg" onClick={this.minus.bind(this)}/>
                <IconCereals/>
            </div>
            <div className="button-option">
                <img alt="plus" src={plus} id="meat" className="icon-dlg" onClick={this.plus.bind(this)}/>
                {this.order.meat}
                <img alt="minus" src={minus} id="meat" className="icon-dlg" onClick={this.minus.bind(this)}/>
                <IconMeat/>
            </div>
            <div className="button-option">
                <img alt="plus" src={plus} id="sweater" className="icon-dlg" onClick={this.plus.bind(this)}/>
                {this.order.sweater}
                <img alt="minus" src={minus} id="sweater" className="icon-dlg" onClick={this.minus.bind(this)}/>
                <IconSweater/>
            </div>
            <div className= "button-ok" onClick={this.onOk.bind(this)}>好</div>
        </div>)

    }
}

export default SettlementDialog