import React, {useState} from "react";
import { AvatarSelector } from "./AvatarSelector";

export const PlayerBuilder = ({player,setPlayer, setPlay}) => {  
   
  const [name, setName] = useState(player.name || "");

  const handleChange = e => {
    setName(e.target.value)
  }

  const setAvatar = avatar => {
    player.avatar = avatar 
    setPlayer({ ...player})
  }

  const updatePlay = e => {
    player.name = name
    setPlayer({ ...player})
    setPlay(true)
  }
 
  return (
    <>
      <div className="userName" key="1">
        <h2>What's your name:</h2>
        <input    
          key="1.2"
          type="text"       
          value={name} 
          onChange={handleChange}
        />
      </div> 
      <AvatarSelector 
        avatar={player.avatar}
        setAvatar={setAvatar}
      />
      <br/><br/>
       
      <button onClick={updatePlay}>Let's Play!</button>   
      
    </>
  )

}  