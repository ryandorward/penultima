import React, { useRef, useEffect } from "react";
import { useRecoilValue } from 'recoil'
import { foodState, healthState, gemState, silverState } from '../recoil/atoms'

export const Stats = () => {
  return (
    <div className="stats">
      <Stat key="food" atom={foodState} label="Food" />
      <Stat key="health" atom={healthState} label="Health" /><br/> 
      <Stat key="gems" atom={gemState} label="Gems" />
      <Stat key="silver" atom={silverState} label="Silver" />
    </div>
  )
}

const Stat = ({atom, label}) => {
  const val = useRecoilValue(atom) 
  const now = useRef()
  const previous = useRef()
  useEffect(() => {
    if (now.current !== val) {            
      previous.current = now.current
      now.current = val      
    } 
  }, [val])

  return (
    <div className={'stats-item ' + (val > previous.current ? 'up':'') + (val < previous.current ? 'down':'') } key={val}>
      {label}: {val}
      &nbsp;&nbsp;
    </div>
  )
}