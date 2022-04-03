import React, { useState, useEffect } from "react"
import { sendMsg, requestMove, requestPeerGem, requestMagicSpell, requestLook, requestTalk } from "../api"
import { Console } from "./Console"
import { Board } from "./Board"
import { Moons } from "./Moons"
import { Wind } from "./Wind"

export const Xoxaria = ({messages, setMessages, board, moons, wind, gemPeer}) => {  
  
  const [animationOn, setAnimationOn] = useState(false)
  const [commandContext, setCommandContext] = useState()
  const [directionContext, setDirectionContext] = useState()
  const [textInput, setTextInput] = useState('')

  useEffect(() => {
    document.addEventListener("keydown", catchKeys, false) // mount
    return () => { // unmount
      document.removeEventListener("keydown", catchKeys, false)
    };
  }, [messages,commandContext])

  function updateMessages(messageSet) {
    if (! Array.isArray(messageSet))
      messageSet = [messageSet]
    for (let i=0;i<messageSet.length;i++) 
      messageSet[i] = {message: messageSet[i]}  
    setMessages( [...messages, ...messageSet].slice(-22))
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
    directionCallback && directionCallback({dir: move, dirLabel: dirLabel})     
  }

  const dialogue = ({e, success, cancel, directionCallback}) => {
    if (e.keyCode === 13) { // enter
      success()     
      setTextInput('')
    }
    else if (e.keyCode === 27) { // esc
      setCommandContext('')
      setTextInput('')
      setMessages([...messages.slice(0,-1)]) //  {message: 'No spell!'}]) 
      cancel()      
    }   
    else if (e.keyCode === 8) { // backspace 
      const newTextInput = textInput.slice(0,-1)
      setTextInput(newTextInput)  
      setMessages([...messages.slice(0,-1),{message: newTextInput}])
    }
    else if (String.fromCharCode(e.keyCode).match(/(\w|\s)/g)) { // normal text/number input
      const newTextInput = textInput + e.key
      setTextInput(newTextInput)
      setMessages([...messages.slice(0,-1),{message: newTextInput}]) // constantly trim off last message and update with new one. 
      // ^ Depends on last message being initialized with a character when dialogue starts. Can be empty       
    } // else pressed key is a non-char e.g. 'esc', 'backspace', 'up arrow'          
    else if (directionCallback) // what happens if a direction key is pressed during dialogue
      catchDirection({e, directionCallback})
    else console.log(e.keyCode, e.code)
  }

  const catchKeys = e => {        
    if (!commandContext) {      
      catchDirection({e, directionCallback: ({dir}) => {              
        requestMove(dir)      
        e.preventDefault()               
      }})
      if (e.code=='KeyP') requestPeerGem()  
      if (e.code=='KeyC') {           
        updateMessages(['Cast spell. Which spell?', ' '])
        //setMessages( [...messages, {message:'Cast spell. Which spell?'},{message:' '}].slice(-22))
        setCommandContext('cast spell')
      }
      if (e.code=='KeyL'){
        updateMessages(['Look. Which direction?', ' '])
        setCommandContext('look')       
      }
      if (e.code=='KeyT'){
        updateMessages(['Talk. Which direction?', ' '])
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
        cancel: () =>  setMessages([...messages,{message: 'No spell!'}])
      })           
    else if (commandContext=='look')
      catchDirection({e, directionCallback: ({dir, dirLabel}) => {       
        requestLook(dir)    
        updateMessages('Look ' +dirLabel)  
        e.preventDefault()       
        setCommandContext()        
      }})
    else if (commandContext=='talk')
      if (! directionContext)
        catchDirection({e, directionCallback: ({dir, dirLabel}) => {  
          setDirectionContext(dir)
          console.log('direction', dir, dirLabel)
          updateMessages('Talk ' +dirLabel)  
          requestTalk({dir: dir})
          e.preventDefault()
        }})
      else {
        const cancel = () => {
          setMessages([...messages,{message: 'Bye!'}])            
          setDirectionContext()          
          setCommandContext()
          setTextInput('')
        }
        dialogue({
          e,
          success: () => {
            requestTalk({dir: directionContext, text: textInput})
           //  updateMessages([textInput])            
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
          <Moons moons={moons} />
          <Wind wind={wind} />
        </div>
        <div>          
          <input type="checkbox" id='animations' checked={animationOn} onChange={e=>setAnimationOn(!animationOn) }/>     
          <label htmlFor="animations">Animations?</label>
        </div>
      </div>
      <div className='wrap-board-console'>   
                    
        <Board 
          board={board}        
          wind={wind}
          animationOn={animationOn}
          gemPeer={gemPeer}
        /> 
      
        <Console messages={messages} className="console"/>
      </div>     
    </>
  );
}