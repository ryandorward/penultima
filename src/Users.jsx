import React from "react";

export const Users= ({messages}) => {   
  return (  
    <div className="UsersInfo">
      <h2>Users</h2>
      {
        messages && messages.map((message, index) => {                            
          if (message.type === 2) return (                 
            <span key={index}>Count: {message.body} </span>
          )         
        })
      }
    </div>
  )
}