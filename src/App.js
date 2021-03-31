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

  const updatePlayer = player => {
    setPlayer(player)
    cookies.set('player',player)
  }

  const resetPlayer = () => {
    player.name = ""     
    setPlayer({ ...player})
    cookies.set('player',player)
  }

  // cookies.remove('player')

  return (
    <div className="App">
      <header 
        className="App-header"
        key = {player.name}      
      > 
      { player.name }
        { ! player.name && 
          <PlayerBuilder 
            player={player}
            setPlayer={updatePlayer}
          />
        }
        {player.name && (   
          <Xoxaria           
            player = {player }           
            resetPlayer = {resetPlayer} 
          />                     
        )}
      </header>
    </div>
  );
}

export default App;
