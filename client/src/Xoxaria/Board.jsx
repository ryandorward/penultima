import React from "react"
import { useRecoilValue } from 'recoil'
import { windState, boardState, gemPeerState, extraBoardClassesState } from '../recoil/atoms'

export const Board = ({ animationOn}) => {    
  const wind = useRecoilValue(windState)      
  const board = useRecoilValue(boardState) 
  const gemPeer = useRecoilValue(gemPeerState)  
  const extraBoardClasses = useRecoilValue(extraBoardClassesState)   

  let boardClass = "board" + (wind ? ' wind-' + wind : '') + (animationOn ? ' animate':'') + (gemPeer ? ' gemPeer' : '' ); 
  boardClass = extraBoardClasses ? boardClass + ' ' + extraBoardClasses.join(' ') : boardClass;
  let displayGrid = gemPeer ? gemPeer : board 
  return (
    <div className={boardClass}>
      {               
        displayGrid.map((row, y) => {            
          return (
            <div className = "row" key={y}>
              { row.map( (_, x) => {

                let tile = displayGrid[x][y]
                if (tile >= 1000) {
                  /*
                  if (tile == 1464) // headless                    
                    tile = 1464 + (x+y) % 4                  
                  */
                 
                  // spritesheet is 32x16, index starts at 1000. Calculate the x background position                 
                  const background_x = ((tile-1000) % 32 ) * -32
                  // spritesheet is 32x16, calculate the y background position
                  const background_y = Math.floor((tile-1000) / 32) * -32
                  const background = `${background_x}px ${background_y}px`                               
                  return <div 
                    key={x} 
                    className={"cell tile indexed-tile tile-" + tile} 
                    style={ {backgroundPosition: background, color: 'red'} }
                  /> 
                }                  
                return (
                  <div key={x} className={"cell tile tile-" + tile}/> 
                )
              }) }             
            </div>
          )
        })
      }  
      <div className='center-ring'></div>
      <div className='center'></div>   
    </div>
  )    
} 