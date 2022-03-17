import React, { useState, useEffect } from "react";
import { sendMsg, requestMove, requestPeerGem } from "../api";
import { Console } from "./Console";
import { Board } from "./Board";
import { Moons } from "./Moons";
import { Wind } from "./Wind";
import { GemPeer } from "./GemPeer" 

export const Xoxaria = ({messages, board, moons, wind, gemPeer}) => {  
  
  const [animationOn, setAnimationOn] = useState(false)
  
  useEffect(() => {
    document.addEventListener("keydown", catchKeys, false) // mount
    return () => { // unmount
      document.removeEventListener("keydown", catchKeys, false)
    };
  }, []);
  
  const send = (e) => {
    if(e.keyCode === 13) { // enter
      sendMsg(e.target.value) 
      e.target.value = ""
    }
  } 

  const catchMove = e => {
    let move
    switch (e.key) {
      case 'ArrowUp':
        move = {x: 0, y: -1}         
        break;
      case 'ArrowDown':
        move = {x: 0, y: 1}       
        break;
      case 'ArrowLeft':
        move = {x: -1, y: 0}       
        break;
      case 'ArrowRight':
        move = {x: 1, y: 0}       
        break;
    }
    if (move){
      requestMove(move)      
      e.preventDefault() 
      return true    
    }
  }

  const catchKeys = e => {    
    if (catchMove(e)) return
    if (e.code=='KeyP') requestPeerGem()  
    else console.log(e)
  } 

  return (
    <>               
      <div className='above-board-console-bar'>
        <div>
          <Moons moons={moons} />
          <Wind wind={wind} />
        </div>
        <div>          
          <input type="checkbox" id='animations' checked={animationOn} onChange={e=>setAnimationOn(!animationOn) }/>     
          <label htmlFor="animations">Animations?</label>
        </div>
      </div>
      <div className='wrap-board-console'>   
        { gemPeer ? 
          <GemPeer gemPeer={gemPeer} /> :       
          <Board 
            board={board}        
            wind={wind}
            animationOn={animationOn}
          /> 
        }
        <Console messages={messages} className="console"/>
      </div>
    </>
  );
}