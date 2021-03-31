import React from "react";

export const ChatInput = ({keydown}) => {  
 
  return (  
    <div className="ChatInput">
      <input onKeyDown={keydown} />
    </div>
  )    

}  