import React from "react";

export const Board = ({board, wind, animationOn}) => {          
  const boardClass = "board" + (wind ? ' wind-' + wind : '') + (animationOn ? ' animate':'');
  return (     
    <div className={boardClass}>
      {               
        board.map((row, y) => {            
          return (
            <div className = "row" key={y}>
              {
                row.map((_, x) => {                    
                  let tileClass = "cell tile tile-" + board[x][y]              
                  return <div key={x} className={tileClass}></div>                
                })
              }
            </div>
          )
        })
      }
    </div>
  )    
}  