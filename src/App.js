import './App.css';
import React, { useState, useEffect } from "react";
import { PlayerBuilder } from "./PlayerBuilder";
import { Xoxaria } from "./Xoxaria";
import Cookies from 'universal-cookie';

const cookies = new Cookies();

function App() {

  const [player, setPlayer] = useState(
    cookies.get('player') || {
      name: ''
    }
  ); 
  const [play, setPlay] = useState(false);  

  const updatePlayer = player => {
    setPlayer(player)
    cookies.set('player',player)
  }

  const resetPlayer = () => {
    player.name = ""     
    setPlayer({ ...player})
    cookies.set('player',player)
  }

  
  return (
    <div className="App">
      <header 
        className="App-header"
        key = {player.name}      
      >        
        { ! play && 
          <PlayerBuilder 
            player={player}
            setPlayer={updatePlayer}
            setPlay={setPlay}
            key="1"
          />
        }
        {play && (   
          <Xoxaria           
            player = {player }           
            resetPlayer = {resetPlayer} 
            setPlay={setPlay}
          />                     
        )}
      </header>
    </div>
  );
}

export default App;
