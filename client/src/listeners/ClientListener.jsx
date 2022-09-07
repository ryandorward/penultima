import React, { useState, useEffect } from "react";
import { useRecoilState, useSetRecoilState, useRecoilValue } from 'recoil'
import { messagesState, gemPeerState, extraBoardClassesState, promptState } from '../recoil/atoms'
import { getLastMessage } from '../recoil/selectors'
import { 
  requestMove, requestPeerGem, requestMagicSpell, 
  requestLook, requestTalk, requestSimpleAction,
  requestDirectionalAction
} from "../api"
import { catchDirection, messageDelay } from './listenerUtilities'

export const ClientListener = () => {    
  
  const [directionContext, setDirectionContext] = useState()
  const [commandContext, setCommandContext] = useState()  
  const [messages, setMessages] = useRecoilState(messagesState)
  const [extraBoardClasses, setExtraBoardClasses] = useRecoilState(extraBoardClassesState)
  const setGemPeer = useSetRecoilState(gemPeerState)
  const setPrompt = useSetRecoilState(promptState)
  const lastMessage = useRecoilValue(getLastMessage)
        
  useEffect(() => {
    document.addEventListener("keydown", catchKeys, false) // mount
    return () => document.removeEventListener("keydown", catchKeys, false) // unmount  
  }, [commandContext, messages])
      
  const cancelPrevious = () => setMessages((msgs) => msgs.slice(0,-1))
  
  // wrapper fn to trim the messages array as we go
  const messagesSetter = (setter) => setMessages((msgs) => setter(msgs).slice(-22)) 
  
  // msgs can be a single message, or an array of messages
  const appendMsg = (msg) =>  messagesSetter((msgs) => [...msgs, ...(Array.isArray(msg) ? msg : [msg]) ]) 

  const queueMsg = (message, prompt, callback) => {
    setPrompt('empty')                
    appendMsg(message)
    setTimeout(function (){  
      setPrompt(prompt || 'default')  
      callback && callback()              
    }, messageDelay)
  }
 
  const catchKeys = e => {            
    if (e.key) {
      setGemPeer() // any key will escape gemPeer state
    }
    if (!commandContext) {      
      catchDirection({e, directionCallback: ({dir}) => {              
        requestMove(dir)      
        e.preventDefault()
      }})
      if (e.code=='KeyP') requestPeerGem()  
      if (e.code=='KeyC') {                                
        queueMsg(['> Cast spell. Which spell?', ' '], 'dialogue')       
        setCommandContext('cast spell')
      }
      if (e.code=='KeyL'){                  
        // queueMsg(['> Look.', 'Which direction?']) 
        setPrompt('empty')                
        appendMsg('> Look. Which direction?')
        setCommandContext('look')       
      }
      if (e.code=='KeyT'){              
        queueMsg(['> Talk. Which direction?']) 
        setCommandContext('talk')       
      }
      if (e.code=='KeyA'){              
        queueMsg(['> Attack. Which direction?']) 
        setCommandContext('attack')       
      }
      if (e.code=='KeyE') requestSimpleAction('enter') 
      if (e.code=='KeyG') requestSimpleAction('get')  
    } 
    else if (commandContext=='cast spell') {     
      dialogue({
        e,
        success: () => {
          setTimeout(() => requestMagicSpell(lastMessage), 850)          
          const orig = extraBoardClasses
          setExtraBoardClasses([...extraBoardClasses,'spell-cast'])
          // wait 0.5s before removing the spell-cast class
          setTimeout(() => setExtraBoardClasses(orig), 1000)
          setCommandContext()
          setPrompt('default') 
        },
        cancel: () => {
          appendMsg('No spell!') 
          setCommandContext()
        }
      })      
    }     
    else if (commandContext=='look')
      catchDirection({e, directionCallback: ({dir, dirLabel}) => {                   
        queueMsg('> Look ' + dirLabel, 'default', () => requestLook(dir) )       
        e.preventDefault()
        setCommandContext()
      }})
    else if (commandContext=='attack')
      catchDirection({e, directionCallback: ({dir, dirLabel}) => {                   
        queueMsg('> Attack ' + dirLabel, 'default', () => requestDirectionalAction({dir: dir, action: "attack"}) )       
        e.preventDefault()
        setCommandContext()
      }})
    else if (commandContext=='talk') {
      if (! directionContext)
        catchDirection({e, directionCallback: ({dir, dirLabel}) => {  
          setDirectionContext(dir)          
          queueMsg(['> Talk ' + dirLabel], 'dialogue')
          requestTalk({dir: dir})          
          e.preventDefault()
        }})
      else {
        const cancel = () => {
          appendMsg('Bye!') 
          setPrompt('default')  
          setDirectionContext()
          setCommandContext()
        }       
        dialogue({
          e,
          success: () => requestTalk({dir: directionContext, text: lastMessage}),          
          cancel: cancel,
          directionCallback: cancel,
        })
      }   
    }   
  }   

  const dialogue = ({e, success, cancel, directionCallback}) => {
   
    if (e.keyCode === 13) {// enter
      // messagesSetter((msgs) => [...msgs.slice(0,-1),lastMessage?.slice(0,-1)]) // get rid of dialogue prompt
      // setPrompt('default')
      success()            
    }
    else if (e.keyCode === 27) { // esc
      messagesSetter((msgs) => [...msgs.slice(0,-1),lastMessage?.slice(0,-1)])
      setCommandContext('')       
      cancelPrevious() // {message: 'No spell!'}]) 
      cancel()
    }
    else if (e.keyCode === 8) { // backspace                               
      messagesSetter((msgs) => [...msgs.slice(0,-1),lastMessage?.slice(0,-1)])
    }
    else if (catchDirection({e, directionCallback})) {}// what happens if a direction key is pressed during dialogue      
    //else if (String.fromCharCode(e.keyCode).match(/(\w|\s|\(|\))/g) || e.key === ':' || e.key === ',')  {// normal text/number input        
    else if (e.key.length === 1) { // this probably doesn't work with emojis or chinese?    
      const newMessage = lastMessage + e.key
      messagesSetter((msgs) => [...msgs.slice(0,-1),newMessage] )
      // ^ Depends on last message being initialized with a character when dialogue starts. Can be empty       
      // else pressed key is a non-char e.g. 'esc', 'backspace', 'up arrow'          
    }
   
    else console.log(e.keyCode, e.code, String.fromCharCode(e.keyCode), e.key)
  }

  return <></>    
}