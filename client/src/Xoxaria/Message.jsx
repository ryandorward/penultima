import React from "react";
export const Message = ({message, prompt}) => {

  return (
    <>
    <div className="Message">
      {message}
      {prompt == 'dialogue' && (
        <span className='dialogue-prompt'></span>
      )}
      
    </div> 
    {prompt == 'default' && (
      <span className='default-prompt'>{'>'}</span>
    )}
    </>
  )

}