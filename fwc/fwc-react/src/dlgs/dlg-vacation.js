import React, { Component } from 'react';
import '../App.css';
import {IconCoin, IconHeart} from './icons'
import fwb from '../fwb'

class VacationDialog extends Component{
    onClickTakeOff(){
        this.props.onOk("take-off")
    }

    onClickTakeVacation(){
        this.props.onOk("take-vacation")
    }


    render(){
        let prop= fwb.getActionProperties("take-off")
        let prop1= fwb.getActionProperties("take-vacation")
        console.log("prop output", prop, prop1)

        return <div className="topdown">
            <div className="card-info">您要进行 </div>
            <div className="button-option" onClick={this.onClickTakeOff.bind(this)}><p className="highlight-info">休闲</p>获得{prop.gain.heart}<IconHeart/></div>
            <div>或者</div>
            <div className= "button-option" onClick={this.onClickTakeVacation.bind(this)}>支付{prop1.cost.gold}<IconCoin/>进行<p className="highlight-info">度假</p>以获得
            {prop1.gain.heart}<IconHeart/>
            </div>
            <div className= "button-cancel" onClick={this.props.onCancel}>再想想</div>
        </div>
    }
}

export default VacationDialog