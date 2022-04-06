import React, { useState } from 'react'
import { Console } from "./Console"
import { Board } from "./Board"
import { Moons } from "./Moons"
import { Wind } from "./Wind"
 
export const Xoxaria = () => {  

  const [animationOn, setAnimationOn] = useState(false)
  
  return (
    <>               
      <div className='above-board-console-bar'>
        <div>
          <Moons />
          <Wind />
        </div>
        <div>          
          <input type="checkbox" id='animations' checked={animationOn} onChange={e=>setAnimationOn(!animationOn) }/>     
          <label htmlFor="animations">Animations?</label>
        </div>
      </div>
      <div className='wrap-board-console'>                       
        <Board animationOn={animationOn} />      
        <Console className="console"/>
      </div>     
    </>
  );
}