import React, { useState, useEffect } from "react";
import { connect, sendMsg, requestMove } from "./api";
import { ChatHistory } from "./ChatHistory";
import { ChatInput } from "./ChatInput";
import { Board } from "./Board";

export const Xoxaria = ({player, resetPlayer}) => {  

  const [messages, setMessages] = useState([]);   

  const [board, setBoard] = useState({
    terrain: [],
    animals: []
  })

  useEffect(() => {
    document.addEventListener("keydown", catchKeys, false);
    return () => {
      document.removeEventListener("keydown", catchKeys, false);
    };
  }, []);
  
  // updates every time messages updates
  useEffect(() => {        
    connect({
      messageCallback: msg => {setMessages([...messages, msg])},
      moveCallback: msg => { 
      },
      updateWorldViewCallback:  msg => {
        setBoard({
          terrain: msg.terrainView,
          animals: msg.animalView
        })
      }
    });
  },[messages]);

  const send = (e) => {
    if(e.keyCode === 13) { // enter
      sendMsg(e.target.value);
      e.target.value = "";
    }
  } 

  const catchKeys = e => {    
    switch (e.keyCode) {
      case 38: // up
      case 40: // down
      case 37: // left
      case 39: // right
        requestMove(e.keyCode);
        e.preventDefault();
        break;
    }
  }

  return (
    <>    
      <h2>Player: { player.name }</h2> 
      <button onClick={resetPlayer}>Start over!</button>      
      <Board 
        board={board}
        player={player} 
      />         
      <ChatHistory messages={messages} />
      <ChatInput keydown={send} />                 
    </>
  );
}