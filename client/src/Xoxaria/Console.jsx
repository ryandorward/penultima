import React, {useState } from "react";
import { useRecoilValue } from 'recoil'
import { messagesState, promptState } from '../recoil/atoms'
import { getLastMessage } from '../recoil/selectors'
import { Message } from "./Message"
import { Stats} from "./Stats"

export const Console = () => {    
   
  const prompt = useRecoilValue(promptState)
  const messages = useRecoilValue(messagesState)
  const lastMessage = useRecoilValue(getLastMessage)

  return (    
    <div className="console">      
      <Stats />
      <div className = 'messages'>
        <div className='console-inner'>
          {        
            messages && messages.slice(0,-1).map((message, index) => {               
              if (message) {              
                return (                           
                  <Message 
                    message = {message} 
                    key = {index}
                    prompt = {null}
                  />    
                )
              }          
            })
          } 
          {        
            lastMessage && (   
              <>                             
                <Message 
                  message = {lastMessage} 
                  key = 'lastmessage'
                  prompt = {prompt}
                />            
              </>    
            )                         
          }       
        </div>
      </div>
    </div>
  )    
}