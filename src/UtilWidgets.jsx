import React, { useState } from "react";
import { host } from './settings';

export const UtilWidgets = ({setPlay, player}) => { 

  const updatePlay = (e) => {setPlay(false)}
  const [logoutErrorMessage, setLogoutErrorMessage] = useState([])

  const handleLogout = () => {
    fetch('//'+host+':8084/logout', {    
        method: 'post'
      }).then(res => {
        if (res.status === 200) {
          // redirect to login
          window.location.href = '/';
        } else if (res.status === 400) {
          setLogoutErrorMessage('There was a problem logging out. Sorry.');
        }
    });
  }    

  return (
    <>
      <div className='top-widgets'>
        <span>Player: { player.name }</span> 
        <button onClick={updatePlay}>Change player!</button>
        <button onClick={handleLogout}>Log out</button>                   
      </div>
      <div>{logoutErrorMessage}</div>  
    </>
  )
}  