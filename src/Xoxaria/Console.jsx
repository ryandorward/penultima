import React, { useState, useEffect } from "react";
import { Message } from "./Message";

export const Console= ({messages}) => {    
  return (  
    <div className="console">
      <div className='console-inner'>
        {        
          messages && messages.map((message, index) => {               
            if (message.message) {              
              return (                           
                <Message 
                  message = {message.message} 
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