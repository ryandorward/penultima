import React from "react";
import { AvatarSelector } from "./AvatarSelector";

export const PlayerBuilder = ({player,setPlayer}) => {  
  
  const keydown = (e) => {
    if(e.keyCode === 13) { // enter  
      player.name = e.target.value     
      setPlayer({ ...player})
      e.target.value = "";
    }
  }

  const setAvatar = avatar => {
    player.avatar = avatar 
    setPlayer({ ...player})
  }
 
  return (
    <>
      <div className="userName">
        <h2>Please Enter your name to play: {player.name}</h2>
        <input onKeyDown={keydown} />
      </div>  
      <AvatarSelector 
        avatar={player.avatar}
        setAvatar={setAvatar}
      />
    </>

  )

}  