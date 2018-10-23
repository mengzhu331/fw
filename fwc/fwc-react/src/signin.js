import React, { Component } from 'react';
import fwb from './fwb'
import DlgInfo from './dlg-info'
import './App.css';

class SignIn extends Component {
    constructor(props) {
        super(props)
        this.unInput = React.createRef();
        this.onClickLogin= this.onClickLogin.bind(this)
        fwb.onOpen= this.onConnected.bind(this)
        this.state={
            disableInput: false,
            connected: false
        }
    }

    onConnected(event){
        this.setState({
            connected: true
        })
    }

    onClickLogin(e){
        if(this.unInput.current.value=== ""){
            alert("必须先输入昵称")
            return
        }
        this.setState({
            disableInput: true
        })
        fwb.connect(this.unInput.current.value, "")
    }

    render(){
        var dlginfo=""

        if(this.state.connected){
            var content= this.unInput.current.value+ "已连接游戏服务器，请等待其他玩家"
            dlginfo=  <DlgInfo header="请稍候" content={content} icon="wait"></DlgInfo>
        }
        return (
        <div className="signin">
            <div></div>
            <div className="dlg-row">
                <div id="si-dialog">
                    <div id="si-info">你的妮称</div>
                    <input id= "si-username" type="text" disabled={this.state.disableInput} ref= {this.unInput}/>
                    <div id="si-login-container">
                            <button id= "si-login" disabled={this.state.disableInput} onClick={this.onClickLogin}>进入游戏</button>
                    </div>
                </div>
            </div>
            {dlginfo}
       </div>
       )
    }
}

export default SignIn