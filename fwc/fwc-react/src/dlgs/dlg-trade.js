import React, { Component } from 'react';
import '../App.css';
import {IconCoin, IconCereals, IconBeer, IconMeat, IconWool, IconSweater} from './icons'
import fwb from '../fwb'
import plus from '../img/plus.png'
import minus from '../img/minus.png'

class TrainDialog extends Component{
    onClick(event){
        this.order.direction= event.currentTarget.id
        this.forceUpdate()
    }

    onOk(){
        fwb.sendAction({
            action:"trade",
            param:this.order
        })
    }

    constructor(props){
        super(props)
        this.order={}
        this.order.direction=""
        this.order.cereals= 0
        this.order.meat= 0
        this.order.wool= 0
        this.order.beer= 0
        this.order.sweater= 0
    }
    componentDidMount(){
        this.order={}
        this.order.direction= ""
        this.order.cereals= 0
        this.order.meat= 0
        this.order.wool= 0
        this.order.beer= 0
        this.order.sweater= 0
    }
    
    clone(obj) {
        if (null == obj || "object" != typeof obj) return obj;
        var copy = obj.constructor();
        for (var attr in obj) {
            if (obj.hasOwnProperty(attr)) copy[attr] = obj[attr];
        }
        return copy;
    }
    
    plus(type){
        let order= this.clone(this.order)
        let pd= fwb.getPd()
        order[type]++

        if(this.order.direction==="sell"){
            console.log("plus", order.direction, order[type], pd[type])
            if(order[type]> pd[type]){
                return
            }
        }else{
            let prop= fwb.getActionProperties("trade", order)
            console.log("plus", order.direction, prop.cost.gold, pd.gold)
            if(prop.cost.gold> pd.gold){
                return
            }
        }

        this.order[type]= order[type]
        this.forceUpdate()
    }

    minus(type){
        if(this.order[type]>0){
            this.order[type]--
        }
    }

    render(){
        let prop= fwb.getActionProperties("trade", this.order)

        let direction=(<div className="topdown">
        <div>您想进行买入还是卖出呢？</div>
        <div className="button-option" id="buy" onClick={this.onClick.bind(this)}>我想买进</div>
        <div className="button-option" id="sell" onClick={this.onClick.bind(this)}>我想卖出</div>
        </div>)

        if(this.order.direction!== ""){
            direction= null
        }

        let order= null
        if(this.order.direction!== ""){
            order=<div className= "topdown-order">
                <div className="button-order">
                    <IconCereals/>
                    <img alt="plus" src={plus} className="icon-dlg" type="cereals" onClick={this.plus.bind(this, "cereals")}/>
                    {this.order.cereals}
                    <img alt="minus" src={minus} className="icon-dlg" type="cereals" onClick={this.minus.bind(this, "cereals")}/>
                </div>
                <div className="button-order">
                    <IconMeat/>
                    <img alt="plus" src={plus} className="icon-dlg" type="meat" onClick={this.plus.bind(this, "meat")}/>
                    {this.order.meat}
                    <img alt="minus" src={minus} className="icon-dlg" type="meat" onClick={this.minus.bind(this, "meat")}/>
                </div>
                <div className="button-order">
                    <IconWool/>
                    <img alt="plus" src={plus} className="icon-dlg" type="wool" onClick={this.plus.bind(this, "wool")}/>
                    {this.order.wool}
                    <img alt="minus" src={minus} className="icon-dlg" type="wool" onClick={this.minus.bind(this, "wool")}/>
                </div>
                <div className="button-order">
                    <IconBeer/>
                    <img alt="plus" src={plus} className="icon-dlg" type="beer" onClick={this.plus.bind(this, "beer")}/>
                    {this.order.beer}
                    <img alt="minus" src={minus} className="icon-dlg" type="beer" onClick={this.minus.bind(this, "beer")}/>
                </div>
                <div className="button-order">
                    <IconSweater/>
                    <img alt="plus" src={plus} className="icon-dlg" type="sweater" onClick={this.plus.bind(this, "sweater")}/>
                    {this.order.sweater}
                    <img alt="minus" src={minus} className="icon-dlg" type="sweater" onClick={this.minus.bind(this, "sweater")}/>
                </div>
           </div>
        }

        let costgain= <div>
            <div>您需要支付</div>
            <div>{prop.cost.gold}<IconCoin/></div>
        </div>
        if(this.order.direction===""){
            costgain=null
        }else if(this.order.direction==="sell"){
            costgain= <div>
            <div>您能获得</div>
            <div>{prop.gain.gold-prop.cost.gold}<IconCoin/></div>
            </div>
        }

        let submit= null
        if(this.order.direction!== ""){
            submit= <div className= "button-ok" onClick={this.onOk.bind(this)}>好</div>
        }
    return <div className= "topdown">
        {direction}
        {order}
        {costgain}
        {submit}
        <div className= "button-cancel" onClick={this.props.onCancel}>再想想</div>
        </div>
    }
}

export default TrainDialog