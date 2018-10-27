import React, { Component } from 'react';
import '../App.css';
import {IconMeat} from './icons'
import fwb from '../fwb'

class HuntDialog extends Component{
    onOk(){
        this.props.onOk("hunt")
    }

    render(){
        let prop= fwb.getActionProperties("hunt")

        return <div>
            <div className="card-info">您要进行 <p className="highlight-info">打猎</p> 以获得肉吗？</div>
            <div>您将获得：</div>
            <div className= "highlight-info">{prop.gain.meat}<IconMeat/></div>
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

export default HuntDialog