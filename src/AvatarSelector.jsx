import React from "react";
import { RadioGroup, FormControl, FormControlLabel, FormLabel, Radio } from '@material-ui/core';

export const AvatarSelector = ({avatar,setAvatar}) => {  
  

  const handleChange = (event) => {
    setAvatar(event.target.value);
  };

  return (
    <FormControl component="fieldset">
      <FormLabel component="legend">Select your Adventurer</FormLabel>
      <RadioGroup aria-label="adventurer" name="gender1" value={avatar} onChange={handleChange}>
        <FormControlLabel value="1" control={<Radio />} label={<AvatarTile index="1" />} />
        <FormControlLabel value="2" control={<Radio />} label={<AvatarTile index="2" />} />
        <FormControlLabel value="3" control={<Radio />} label={<AvatarTile index="3" />} /> 
        <FormControlLabel value="4" control={<Radio />} label={<AvatarTile index="4" />} />   
        <FormControlLabel value="5" control={<Radio />} label={<AvatarTile index="5" />} /> 
        <FormControlLabel value="6" control={<Radio />} label={<AvatarTile index="6" />} />        
      </RadioGroup>
    </FormControl>

  )

}  

const AvatarTile = ({index}) => {
  const className = "avatar-tile avatar-tile-" + index
  return (
    <div className={className} />
  )

}