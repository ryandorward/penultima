import React from "react";

export const Board = ({board, player}) => {          

  console.log(board.animals, board.terrain)

  return (  
    <div className="board">
      {               
        board.terrain.map((row, y) => {   
          return (
          <div 
            className = "row"
            key={y}
          >
            {
              row.map((cell, x) => {     
                                
                let tileClass = "cell tile "    
                var animalTileClass = "avatar-tile animal-tile avatar-tile-"+ board.animals[x][y]            
                if ( x== 7 && y == 7 )
                  animalTileClass = "avatar-tile animal-tile avatar-tile-"+ player.avatar                 
                             
                tileClass += "tile-" + board.terrain[x][y]
                return (
                <div key={x} className ={tileClass}>
                  { ((board.animals[x][y] !== 0) || ( x== 7 && y == 7 ))  &&                                      
                    <div className={animalTileClass}>
                    </div>                  
                  }
                </div>
                )
              })
            }
          </div>
          )
        }) 
      }
    </div>
  )    

}  