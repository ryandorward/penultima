import React, { useState, useEffect } from "react";
import { useRecoilState, useSetRecoilState, useRecoilValue } from 'recoil'
import { messagesState, gemPeerState } from '../recoil/atoms'
import { getLastMessage } from '../recoil/selectors'
import { requestMove, requestPeerGem, requestMagicSpell, requestLook, requestTalk } from "../api"
import { catchDirection } from './clientListenerUtils'


export const ClientListener = () => {    
  
  const [directionContext, setDirectionContext] = useState()
  const [commandContext, setCommandContext] = useState()  
  const [messages, setMessages] = useRecoilState(messagesState)
  const setGemPeer = useSetRecoilState(gemPeerState)
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

  const catchKeys = e => {        
    if (e.key) setGemPeer() // any key will escape gemPeer state
    if (!commandContext) {      
      catchDirection({e, directionCallback: ({dir}) => {              
        requestMove(dir)      
        e.preventDefault()               
      }})
      if (e.code=='KeyP') requestPeerGem()  
      if (e.code=='KeyC') {                        
        appendMsg(['> Cast spell.', 'Which spell?', ' '])    
        setCommandContext('cast spell')
      }
      if (e.code=='KeyL'){               
        appendMsg(['> Look.', 'Which direction?']) 
        setCommandContext('look')       
      }
      if (e.code=='KeyT'){              
        appendMsg(['> Talk.', 'Which direction?']) 
        setCommandContext('talk')       
      }
    } 
    else if (commandContext=='cast spell')
      dialogue({
        e,
        success: () => {
          requestMagicSpell(lastMessage)
          setCommandContext()
        },
        cancel: () => {
          appendMsg('No spell!') 
          setCommandContext()
        }
      })           
    else if (commandContext=='look')
      catchDirection({e, directionCallback: ({dir, dirLabel}) => {       
        requestLook(dir)    
        appendMsg('> Look ' + dirLabel)       
        e.preventDefault()       
        setCommandContext()        
      }})
    else if (commandContext=='talk') {
      if (! directionContext)
        catchDirection({e, directionCallback: ({dir, dirLabel}) => {  
          setDirectionContext(dir)          
          appendMsg('> Talk ' + dirLabel)          
          requestTalk({dir: dir})
          e.preventDefault()
        }})
      else {
        const cancel = () => {
          appendMsg('Bye!')       
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
    if (e.keyCode === 13) // enter
      success()            
    else if (e.keyCode === 27) { // esc
      setCommandContext('')       
      cancelPrevious() // {message: 'No spell!'}]) 
      cancel()      
    }
    else if (e.keyCode === 8)  // backspace                               
      messagesSetter((msgs) => [...msgs.slice(0,-1),lastMessage?.slice(0,-1)])
    else if (String.fromCharCode(e.keyCode).match(/(\w|\s)/g))  // normal text/number input    
      messagesSetter((msgs) => [...msgs.slice(0,-1),lastMessage + e.key] )
      // ^ Depends on last message being initialized with a character when dialogue starts. Can be empty       
      // else pressed key is a non-char e.g. 'esc', 'backspace', 'up arrow'          
    else if (directionCallback) // what happens if a direction key is pressed during dialogue
      catchDirection({e, directionCallback})
    else console.log(e.keyCode, e.code)
  }

  return <></>    
}