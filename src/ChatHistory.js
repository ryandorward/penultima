import React, { useState, useEffect } from "react";
import { Message } from "./Message";

export const ChatHistory= ({messages}) => {  
 
  return (  
    <div className="ChatHistory">
      <h2>Chat History</h2>
      {
        messages && messages.map((message, index) => {   
          
          // message = JSON.parse(message.data)
                   
          if (message.type === 1) return (                 
            <Message 
              message = {message} 
              key = {index}
            />    
          )
          else console.log("Message was:", message)
        })
      }
    </div>
  )    

}  