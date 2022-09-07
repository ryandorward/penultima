import React from "react";
import { useRecoilState } from 'recoil'
import { AvatarSelector } from "./AvatarSelector";
import { requestUpdateName } from "./api";
import { nameState } from './recoil/atoms'

export const PlayerBuilder = ({player,setPlayer, setPlay}) => {  

  // const [name, setName] = useState(player?.name);
  const [name, setName] = useRecoilState(nameState)

  const setAvatar = avatar => {
    player.avatar = avatar 
    setPlayer({ ...player})
  }

  /*
  const setName = name => {    
    player.name = name
    setPlayer({ ...player})
  }*/

  const updatePlay = e => {
    setPlayer({ ...player, name: name})
    requestUpdateName(name)
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
          onChange={e => setName(e.target.value)}
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