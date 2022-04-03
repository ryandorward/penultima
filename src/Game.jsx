import './css/App.css';
import React, { useState, useEffect } from "react";
import { PlayerBuilder } from "./PlayerBuilder"
import { Xoxaria } from "./Xoxaria/Xoxaria"
import Cookies from 'universal-cookie'
import { connect } from "./api";
import { UtilWidgets } from "./UtilWidgets"

const cookies = new Cookies()

function Game() {
  const [player, setPlayer] = useState(
    cookies.get('player') || {
      name: ''
    }
  ); 
  const [play, setPlay] = useState(true)
  const [messages, setMessages] = useState([])
  const [moons, setMoons] = useState({trammel:0, felucca:0})
  const [wind, setWind] = useState()
  const [board, setBoard] = useState([])
  const [gemPeer, setGemPeer] = useState()
  
  const updatePlayer = player => {
    setPlayer(player)
    cookies.set('player',player)
  }

  const serverMessageCallback = msg => { 
    if (msg.fov) { 
      // console.log('fov update', msg.fov);
      setBoard(msg.fov)
    }
    else if (msg.message) {
     // console.log('setting messages', messages, msg.message, [...messages, msg.message].slice(-22) )   
      setMessages([...messages, msg.message].slice(-22)) 
    }
    else if (msg.gemPeer) {
      setGemPeer(msg.gemPeer)      
    }     
    else if (msg.zone) {
      const zone = msg.zone
      if (zone.moonPhase)
        setMoons(zone.moonPhase)
      if (zone.wind) {
        if (zone.wind.X === 0 && zone.wind.Y === 0)
          setWind(null)
        else if (zone.wind.X === 0 && zone.wind.Y === -1)
          setWind('north')
        else if (zone.wind.X === 1 && zone.wind.Y === 0)
          setWind('east')
        else if (zone.wind.X === 0 && zone.wind.Y === 1)
          setWind('south')
        else if (zone.wind.X === -1 && zone.wind.Y === 0)
          setWind('west')                   
      }
    }           
    else console.log("Other message",msg)    
  }
  
  // updates every time messages updates. Why? Because the messageCallback needs to be reinitialized with the serverMessageCallback having messages in the proper state
  // otherwise the messages variable is always stuck at the state it was when it was initialized.
  useEffect(() => connect({ messageCallback: serverMessageCallback}), [messages]);
  
  useEffect(() => {
    document.addEventListener("keydown", catchKeys, false) // mount
    return () => document.removeEventListener("keydown", catchKeys, false) // unmount    
  }, []);

  const catchKeys = e => {          
    if (e.key) setGemPeer() // any key will escape gemPeer state
  } 
 
  return (
    <div className="App">

      <header className="App-header" key = {player.name} >                      

        { play ? 
          <><UtilWidgets setPlay={setPlay} player={player} />
          <Xoxaria                                                    
            messages={messages}
            setMessages={setMessages}                   
            board={board}
            moons={moons}
            wind={wind}           
            gemPeer={gemPeer}
        /> </>:
          <PlayerBuilder 
            player={player}
            setPlayer={updatePlayer}
            setPlay={setPlay}
            key="1"                   
          />                                                       
        }
      </header>
    </div>
  )
}

export default Game