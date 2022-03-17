import React from "react";

export const GemPeer = ({gemPeer}) => {          
  const boardClass = "gemPeer board";
  return (
    <div className={boardClass}> 
      {               
        gemPeer.map((row, y) => {   
          return (
            <div className = "row" key={y}>
              { row.map( (_, x) => <div key={x} className={"cell tile tile-" + gemPeer[x][y]}/> ) }
            </div>
          )
        }) 
      }
      <div className='center-ring'></div>
      <div className='center'></div>
    </div>
  )    
}