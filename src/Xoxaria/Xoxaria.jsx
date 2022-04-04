import React, { useState, useEffect, useReducer } from 'react'
import { useRecoilState, useSetRecoilState } from 'recoil'
import { connect, requestMove, requestPeerGem, requestMagicSpell, requestLook, requestTalk } from "../api"
import { Console } from "./Console"
import { Board } from "./Board"
import { Moons } from "./Moons"
import { moonsState, windState, boardState, gemPeerState, messagesState } from '../atoms/atoms'
import { Wind } from "./Wind"
 
export const Xoxaria = () => {  

  const [animationOn, setAnimationOn] = useState(false)
  const [commandContext, setCommandContext] = useState()
  const [directionContext, setDirectionContext] = useState()
  const [textInput, setTextInput] = useState('')

  const setMoons = useSetRecoilState(moonsState)
  const setWind = useSetRecoilState(windState)
  const setBoard = useSetRecoilState(boardState)
  const setGemPeer = useSetRecoilState(gemPeerState)

  const [messages, setMessages] = useRecoilState(messagesState)

  const messagesDispatch = (action) => {
    switch (action.type) {
      case "append":
        if (! Array.isArray(action.messages))
          action.messages = [action.messages] 
        setMessages([...messages, ...action.messages].slice(-22))
        break             
      case "set": // updates the whole message set
        setMessages(action.messages.slice(-22))       
        break       
      case "cancelPrevious":
        setMessages(messages.slice(0,-1))
        break            
    }
  }  
  
  useEffect(() => {
    document.addEventListener("keydown", catchKeys, false) // mount
    return () => { // unmount
      document.removeEventListener("keydown", catchKeys, false)
    };
  }, [messages,commandContext])

  // updates every time messages updates. Why? Because the messageCallback needs to be reinitialized with the serverMessageCallback having messages in the proper state
  // otherwise the messages variable is always stuck at the state it was when it was initialized.
  useEffect(() => connect({ messageCallback: serverMessageCallback}), [messages]);
   
  const appendMsg = msg => messagesDispatch({ type: "append", messages: msg })

  const serverMessageCallback = msg => { 
    if (msg.fov) { 
      // console.log('fov update', msg.fov);
      setBoard(msg.fov)
    }
    else if (msg.message) {      
      appendMsg(msg.message.message)
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


  const catchDirection = ({e, directionCallback}) => {
    let move
    let dirLabel
    switch (e.key) {
      case 'ArrowUp':
        dirLabel = 'North'
        move = {x: 0, y: -1}         
        break;
      case 'ArrowDown':
        dirLabel = 'South'
        move = {x: 0, y: 1}       
        break;
      case 'ArrowLeft':
        dirLabel = 'West'
        move = {x: -1, y: 0}       
        break;
      case 'ArrowRight':
        dirLabel = 'East'
        move = {x: 1, y: 0}       
        break;
    }
    directionCallback && move && directionCallback({dir: move, dirLabel: dirLabel})     
  }

  const dialogue = ({e, success, cancel, directionCallback}) => {
    if (e.keyCode === 13) { // enter
      success()     
      setTextInput('')
    }
    else if (e.keyCode === 27) { // esc
      setCommandContext('')
      setTextInput('')     
      messagesDispatch({ type: "cancelPrevious"}) //  {message: 'No spell!'}]) 
      cancel()      
    }   
    else if (e.keyCode === 8) { // backspace 
      const newTextInput = textInput.slice(0,-1)
      setTextInput(newTextInput)        
      messagesDispatch({ type: "set", messages: [...messages.slice(0,-1),newTextInput] }); //  {message: 'No spell!'}]) 
    }
    else if (String.fromCharCode(e.keyCode).match(/(\w|\s)/g)) { // normal text/number input
      const newTextInput = textInput + e.key
      setTextInput(newTextInput)     
      messagesDispatch({ type: "set", messages:[...messages.slice(0,-1),newTextInput]}) // constantly trim off last message and update with new one.
      // ^ Depends on last message being initialized with a character when dialogue starts. Can be empty       
    } // else pressed key is a non-char e.g. 'esc', 'backspace', 'up arrow'          
    else if (directionCallback) // what happens if a direction key is pressed during dialogue
      catchDirection({e, directionCallback})
    else console.log(e.keyCode, e.code)
  }

  const catchKeys = e => {        
    if (e.key) setGemPeer() // any key will escape gemPeer state
    if (!commandContext) {      
      catchDirection({e, directionCallback: ({dir}) => {              
        requestMove(dir)      
        e.preventDefault()               
      }})
      if (e.code=='KeyP') requestPeerGem()  
      if (e.code=='KeyC') {                        
        appendMsg(['> Cast spell.', 'Which spell?'])    
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
      else console.log(e)
    } 

    // sketch for a dialogue framework
    // commandContext is a stack (array) that you append the command sequence onto.
    // so it might look like ['cast spell','alakazam teleport',[-12,0]]

    // there will be a dialoge mode, and you send callbacks to it:
    // so a callback for submit (triggered when they press enter)
    // it should pass the typed command back to the callback, let the callback decide what to do with it

    // callback will decide how to handle command, probably pop commands off the stack, or add one on

    else if (commandContext=='cast spell')
      dialogue({
        e,
        success: () => {
          requestMagicSpell(textInput)
          setCommandContext()
        },
        cancel: () => appendMsg('No spell!') 
      })           
    else if (commandContext=='look')
      catchDirection({e, directionCallback: ({dir, dirLabel}) => {       
        requestLook(dir)    
        appendMsg('> Look ' + dirLabel )       
        e.preventDefault()       
        setCommandContext()        
      }})
    else if (commandContext=='talk')
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
          setTextInput('')
        }
        dialogue({
          e,
          success: () => {
            requestTalk({dir: directionContext, text: textInput})                    
          },
          cancel: cancel,
          directionCallback: cancel,
        })
      }      
  } 

  return (
    <>               
      <div className='above-board-console-bar'>
        <div>
          <Moons />
          <Wind />
        </div>
        <div>          
          <input type="checkbox" id='animations' checked={animationOn} onChange={e=>setAnimationOn(!animationOn) }/>     
          <label htmlFor="animations">Animations?</label>
        </div>
      </div>
      <div className='wrap-board-console'>                       
        <Board animationOn={animationOn} />      
        <Console className="console"/>
      </div>     
    </>
  );
}