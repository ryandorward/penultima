import React from "react"
import { useRecoilValue } from 'recoil'
import { windState, boardState, gemPeerState } from '../recoil/atoms'

export const Board = ({ animationOn }) => {    
  const wind = useRecoilValue(windState)      
  const board = useRecoilValue(boardState) 
  const gemPeer = useRecoilValue(gemPeerState)   

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
      <div className='center-ring'></div>
      <div className='center'></div>   
    </div>
  )    
} 