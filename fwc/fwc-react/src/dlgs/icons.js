import React, { Component } from 'react';
import '../App.css';
import fwb from '../fwb';

function IconCereals(){
    const iconurl=require('../img/cereals.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconCoin(){
    const iconurl=require('../img/coin.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconMeat(){
    const iconurl=require('../img/meat.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconWool(){
    const iconurl=require('../img/wool.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconSweater(){
    const iconurl=require('../img/sweater.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconBeer(){
    const iconurl=require('../img/beer.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconInte(){
    const iconurl=require('../img/wisd1.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconKnow(){
    const iconurl=require('../img/know1.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconStre(){
    const iconurl=require('../img/stre1.png')
    return <img alt="coin" src={iconurl} className="icon-dlg"/>
}

function IconToken(){
    let color= fwb.gd.pd.get(fwb.gd.clientID).color
    const iconurl= require('../img/token'+ color+'.png')
    return <img alt="token" src={iconurl} className="icon-dlg"/>
}

function IconHome(){
    const iconurl=require('../img/house1.png')
    return <img alt="home" src={iconurl} className="icon-dlg"/>
}

function IconHeart(){
    const iconurl=require('../img/heart.png')
    return <img alt="heart" src={iconurl} className="icon-dlg"/>
}

export {
    IconCereals,
    IconCoin,
    IconMeat,
    IconWool,
    IconSweater,
    IconBeer,
    IconInte,
    IconKnow,
    IconStre,
    IconToken,
    IconHome,
    IconHeart
} 