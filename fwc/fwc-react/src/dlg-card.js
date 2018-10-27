import React, { Component } from 'react';
import './App.css';
import FarmDialog from './dlgs/dlg-farm';
import fwb from './fwb'
import ParttimeWorkDialog from'./dlgs/dlg-parttime-work'
import HuntDialog from'./dlgs/dlg-hunt'
import FeedSheepDialog from './dlgs/dlg-feed-sheep'
import BegDialog from './dlgs/dlg-beg'
import EmployDialog from './dlgs/dlg-employ'
import UpgradeHomeDialog from './dlgs/dlg-upgrade-home';
import VacationDialog from './dlgs/dlg-vacation';
import StealDialog from './dlgs/dlg-steal';
import TrainDialog from './dlgs/dlg-train';
import WeaveDialog from './dlgs/dlg-weave';
import GoldMiningDialog from './dlgs/dlg-gold-mining';
import MakeWineDialog from './dlgs/dlg-make-wine';
import PartyDialog from './dlgs/dlg-party';
import TradeDialog from './dlgs/dlg-trade';

class CardDialog extends Component{
    onOk(action, param){
        console.log("ok clicked ", action)
        fwb.sendAction({
            playerID: fwb.gd.clientID,
            action: action,
            param: param,
        })
        this.props.onOk(this.props.card)
    }

    onCancel(){
        this.props.onCancel(this.props.card)
    }

    render(){
        let actionDlg=null
        let actions= fwb.getCardActions(this.props.card.id)
        if(actions[0]==="farm"){
            actionDlg= <FarmDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="parttime-work"){
            actionDlg= <ParttimeWorkDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="hunt"){
            actionDlg= <HuntDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="feed-sheep"){
            actionDlg= <FeedSheepDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]=== "beg"){
            actionDlg=<BegDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="employ"){
            actionDlg=<EmployDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="upgrade-home"){
            actionDlg=<UpgradeHomeDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="take-off"){
            actionDlg=<VacationDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="steal"){
            actionDlg=<StealDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="train"){
            actionDlg=<TrainDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="weave"){
            actionDlg=<WeaveDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="gold-mining"){
            actionDlg=<GoldMiningDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="make-wine"){
            actionDlg=<MakeWineDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="party"){
            actionDlg=<PartyDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }else if(actions[0]==="trade"){
            actionDlg=<TradeDialog onOk= {this.onOk.bind(this)} onCancel={this.onCancel.bind(this)}/>
        }

        return <div className= "dlg-card">
            {actionDlg}
        </div>
    }
}

export default CardDialog