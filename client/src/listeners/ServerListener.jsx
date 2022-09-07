import { useEffect } from 'react'
import { useRecoilState, useSetRecoilState } from 'recoil'
import { connect } from "../api"
import { 
  moonsState,
  windState,
  boardState, 
  gemPeerState, 
  messagesState,
  promptState,
  nameState,
  foodState,
  healthState,
  gemState,
  silverState
} from '../recoil/atoms'
import { messageDelay } from './listenerUtilities'

export const ServerListener = () => {    

  const setMoons = useSetRecoilState(moonsState)
  const setWind = useSetRecoilState(windState)
  const setBoard = useSetRecoilState(boardState)
  const setGemPeer = useSetRecoilState(gemPeerState)
  const setPrompt = useSetRecoilState(promptState)
  const setName = useSetRecoilState(nameState)
  const setFood = useSetRecoilState(foodState)
  const setHealth = useSetRecoilState(healthState)
  const setMessages = useSetRecoilState(messagesState)
  const setGems = useSetRecoilState(gemState)
  const setSilver = useSetRecoilState(silverState)
  
  // updates every time messages updates. Why? Because the messageCallback needs to be reinitialized with the serverMessageCallback having messages in the proper state
  // otherwise the messages variable is always stuck at the state it was when it was initialized.
  useEffect(() => connect({ messageCallback: serverMessageCallback}), []); // was [messages] 
     
  const serverMessageCallback = msg => { 
    if (msg.fov) { 
      // console.log('fov update', msg.fov);
      setBoard(msg.fov)
    }
    else if (msg.message) {      
      setMessages((msgs) => [...msgs, msg.message.message, ' '].slice(-22))
    }
    else if (msg.result) {     
      let message = msg.result.message 
      setPrompt('empty')
      msg.result.status == 'success' && (message = '> ' + message)                
      setMessages((msgs) => [...msgs, message].slice(-22))      
      setTimeout(function (){  
        setPrompt('default')                    
      }, messageDelay)

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
    if (msg.food)
      setFood(msg.food) 
    if (msg.health)
      setHealth(msg.health)
    if (msg.stats) {   
      // console.log("Stats message",msg)     
      setName(msg.stats)
    }
    if (msg.stat) {   
      // console.log("Stat message",msg)     
      if (msg.stat.name === 'food')
        setFood(msg.stat.value)
      if (msg.stat.name === 'health')
        setHealth(msg.stat.value)
      if (msg.stat.name === 'gems')
        setGems(msg.stat.value)
      if (msg.stat.name === 'silver')
        setSilver(msg.stat.value)
    }
    // else console.log("Other message",msg)    
  }

  return null
} 