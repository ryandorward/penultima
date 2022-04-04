import React from "react";
import { useRecoilValue } from 'recoil'
import { messagesState } from '../atoms/atoms'
import { Message } from "./Message";

export const Console= () => {    
  const messages = useRecoilValue(messagesState)    
  return (  
    <div className="console">
      <div className='console-inner'>
        {        
          messages && messages.map((message, index) => {               
            if (message) {              
              return (                           
                <Message 
                  message = {message} 
                  key = {index}
                />    
              )
            }          
          })
        }       
      </div>
    </div>
  )    
}