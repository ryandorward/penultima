import { useEffect } from 'react'
import { useRecoilState, useSetRecoilState } from 'recoil'
import { connect } from "../api"
import { moonsState, windState, boardState, gemPeerState, messagesState } from '../recoil/atoms'

export const ServerListener = () => {    

  const setMoons = useSetRecoilState(moonsState)
  const setWind = useSetRecoilState(windState)
  const setBoard = useSetRecoilState(boardState)
  const setGemPeer = useSetRecoilState(gemPeerState)

  const [messages, setMessages] = useRecoilState(messagesState)
  
  // updates every time messages updates. Why? Because the messageCallback needs to be reinitialized with the serverMessageCallback having messages in the proper state
  // otherwise the messages variable is always stuck at the state it was when it was initialized.
  useEffect(() => connect({ messageCallback: serverMessageCallback}), []); // was [messages] 
     
  const serverMessageCallback = msg => { 
    if (msg.fov) { 
      // console.log('fov update', msg.fov);
      setBoard(msg.fov)
    }
    else if (msg.message) {      
      setMessages((msgs) => [...msgs, msg.message.message].slice(-22))     
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

  return null
} 