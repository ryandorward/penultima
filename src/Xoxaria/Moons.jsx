import React from "react";
import { useRecoilValue } from 'recoil'
import { moonsState } from '../atoms/atoms'
import '../css/moons.css'

const Moon = ({name, phase}) => {
  const className = "moon" + (name ? " " + name : "") + (phase ? " phase-" + phase : "")
  return (   
    <div className={className}>
      <div className="disc"></div>
    </div>  
  )
}

export const Moons = () => { 

  const moons = useRecoilValue(moonsState)

  return (
    <div className = "moons-wrap">          
      <Moon name="Trammel" phase={moons.trammel}/>
      <Moon name="Felucca" phase={moons.felucca}/>
    </div>
  )
}  