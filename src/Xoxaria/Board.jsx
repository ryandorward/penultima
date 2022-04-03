import React from "react";

export const Board = ({board, gemPeer, wind, animationOn }) => {          
  const boardClass = "board" + (wind ? ' wind-' + wind : '') + (animationOn ? ' animate':'') + (gemPeer ? ' gemPeer' : '' ); 
  let displayGrid = gemPeer ? gemPeer : board 
  return (     
    <div className={boardClass}>
      {               
        displayGrid.map((row, y) => {            
          return (
            <div className = "row" key={y}>
              { row.map( (_, x) => <div key={x} className={"cell tile tile-" + displayGrid[x][y]}/> ) }             
            </div>
          )
        })
      }
    </div>
  )    
} 