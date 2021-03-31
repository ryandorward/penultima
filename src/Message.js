import React, { useState } from "react";

export const Message = ({message}) => {  
    
  return (  
    <div className="Message">
      {message.body}
    </div>
  )    

}