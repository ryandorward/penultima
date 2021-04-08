import React from "react";
import { RadioGroup, FormControl, FormControlLabel, FormLabel, Radio } from '@material-ui/core';

export const AvatarSelector = ({avatar,setAvatar}) => {  
  
  const handleChange = (event) => {
    setAvatar(event.target.value);
  };



  return (
    <>
    <h2>Pick your Guy</h2>
    <FormControl component="fieldset">      
      <RadioGroup className="avatarSelectorRadios" aria-label="adventurer" name="gender1" value={avatar} onChange={handleChange}>             
        <AvatarFormControlLabel index="1" /> 
        <AvatarFormControlLabel index="2" /> 
        <AvatarFormControlLabel index="3" />    
        <AvatarFormControlLabel index="4" /> 
        <AvatarFormControlLabel index="5" />    
        <AvatarFormControlLabel index="6" /> 
             
      </RadioGroup>
    </FormControl>
    </>
    ) 
  
}  


const AvatarFormControlLabel = ({index}) => {
  return (
    <FormControlLabel 
      className="avatar-formControlLabel"
      value={index}
      control={
       <Radio 
        icon={<AvatarTile index={index} checked={false} />}
        checkedIcon={<AvatarTile index={index} checked={true} />}
      />} 
    />   
  )
}

const AvatarTile = ({index, checked}) => {
  const className = "avatar-tile avatar-tile-" + index + (checked ? " checked" : "")
  return (
    <div className={className} />
  )

}