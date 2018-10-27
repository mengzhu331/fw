import React, { Component } from 'react';
import './App.css';
import fwb from './fwb';

class ChooseSkill extends Component{
    
    constructor(props){
        super(props)
        this.skillMap={
            "inte":"智力",
            "know":"知识",
            "stre":"体力",
        }
        this.state={
            choosen:false,
            skill:"",
        }
    }
    componentDidMount(){
        console.log("did mount")
        this.setState({
            choosen:false,
            skill:"",
        })
    }

    onClickSkillButton(event){
        console.log("choose ", event.target, event.target.id)
        fwb.sendGameStartAck(event.target.id)

        this.setState({
            choosen: true,
            skill: this.skillMap[event.target.id],
        })
    }

    render(){
        console.log(this.state)
        let imgInte= require('./img/wisd1.png')
        let imgKnow= require('./img/know1.png')
        let imgStre= require('./img/stre1.png')
        let chooseContent= <div className="choose-skill-dialog">
            <div className= "choose-skill-info">请选择一种基本技能</div>
            <div className= "choose-skill-button" id="inte" onClick={this.onClickSkillButton.bind(this)}>
                <img alt="inte" src={imgInte} id="inte" className="choose-skill-icon"/>智力
            </div>
            <div className= "choose-skill-button" id="know" onClick={this.onClickSkillButton.bind(this)}>
                <img alt="know" src={imgKnow} id="know" className="choose-skill-icon"/>知识
            </div>
            <div className= "choose-skill-button" id= "stre" onClick={this.onClickSkillButton.bind(this)}>
                <img alt="stre" src={imgStre} id="stre" className="choose-skill-icon"/>体力
            </div>
        </div>
        let waitContent=<div className= "skill-chosen-info">已选择<div className="highlight-info">{this.state.skill}</div>请等待其他玩家</div>

        let content=""
        if(this.state.choosen){
            content= waitContent
        }else{
            content= chooseContent
        }
        return content
    }
}

export default ChooseSkill