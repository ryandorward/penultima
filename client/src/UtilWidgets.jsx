import React, { useState } from "react";
import { useRecoilValue } from 'recoil'
import { host } from './settings';
import { nameState } from './recoil/atoms'

export const UtilWidgets = ({setPlay}) => { 

  const updatePlay = (e) => {setPlay(false)}
  const [logoutErrorMessage, setLogoutErrorMessage] = useState([])

  const name = useRecoilValue(nameState)

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
        <span>Player: { name }</span> 
        <button onClick={updatePlay}>Change player!</button>
        <button onClick={handleLogout}>Log out</button>                   
      </div>
      <div>{logoutErrorMessage}</div>  
    </>
  )
}  