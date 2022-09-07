import './css/App.css';
import React, { useState } from 'react'
import { RecoilRoot } from 'recoil'
import { PlayerBuilder } from "./PlayerBuilder"
import { Xoxaria } from "./Xoxaria/Xoxaria"
import Cookies from 'universal-cookie'
import { UtilWidgets } from "./UtilWidgets"
import { ServerListener } from "./listeners/ServerListener"
import { ClientListener } from "./listeners/ClientListener"

const cookies = new Cookies()

function Game() {
  const [player, setPlayer] = useState(
   cookies.get('player') || {
      name: ''
    } 
  ); 
  const [play, setPlay] = useState(true)
  
  const updatePlayer = player => {
    setPlayer(player)
    cookies.set('player',player)
  }

  return (
    <RecoilRoot>       
      <div className="App">
        <header className="App-header" key = {player.name} >                            
          { play ?                       
            <>
              <UtilWidgets setPlay={setPlay} player={player} />
              <Xoxaria />
              <ServerListener />
              <ClientListener />
            </>:
            <PlayerBuilder 
              player={player}
              setPlayer={updatePlayer}
              setPlay={setPlay}
              key="1"                   
            />                                                       
          }        
        </header>
      </div>
    </RecoilRoot>
  )
}

export default Game