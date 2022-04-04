import React from "react";
import { useRecoilValue } from 'recoil'
import { windState } from '../atoms/atoms'

export const Wind = () => { 
  const wind = useRecoilValue(windState)
  return (
    <div className = "wind-wrap">          
      Wind: {wind ? wind : "none"}
    </div>
  )
}  