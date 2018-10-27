import React, { Component } from 'react';
import '../App.css';
import {IconCereals} from './icons'
import fwb from '../fwb'

class FarmDialog extends Component{
    onOk(){
        this.props.onOk("farm")
    }

    render(){
        let prop= fwb.getActionProperties("farm")
        console.log("prop output", prop)

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">耕种</p> 从而获得粮食吗？</div>
            <div>您将获得：</div>
            <div className= "highlight-info">{prop.gain.cereals}<IconCereals/></div>
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

export default FarmDialog