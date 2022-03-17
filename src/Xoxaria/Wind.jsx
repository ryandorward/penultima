import React from "react";

export const Wind = ({ wind }) => { 
  return (
    <div className = "wind-wrap">          
      Wind: {wind ? wind : "none"}
    </div>
  )
}  