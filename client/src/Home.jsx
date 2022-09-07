import React, { useState } from 'react';
import { host } from './settings';

export const Home = () => {         
  
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginErrorMessage, setLoginErrorMessage] = useState('');
  const [registerErrorMessage, setRegisterErrorMessage] = useState('');
  const [tab, setTab] = useState('about');

  const handleSubmit = () => {
    const data = new FormData();
    data.append('username', username);
    data.append('password', password);

    if (tab === 'login') {
      console.log(host)
      fetch('//'+host+':8084/login', {
        body: data,
        method: 'post',
        credentials: 'include'
      }).then(res => {
        if (res.status === 200) {
          // redirect to game        
          window.location.href = '/game';
          
        } else if (res.status === 400) {
          setLoginErrorMessage('There was a problem logging in. Please use a correct username and password.');
        }
      });
    } else if (tab === 'register') {
      fetch('http://'+host+':8084/register', {
        body: data,
        method: 'post',  
        credentials: 'include'     
      }).then(res => {
        if (res.status === 200) {
          window.location.href = '/game';
        } else {
          setRegisterErrorMessage('There was a problem registering. Please pick a novel username.');
        }
      });  
    } else if (tab === 'logout') {      
      fetch('http://'+host+':8084/logout', {      
        method: 'post',  
        credentials: 'include'     
      }).then(res => {
        if (res.status === 200) {
          // window.location.href = '/game';
        } else {
          setRegisterErrorMessage('There was a problem logouting.');
        }
      });
    }
  
  };

  return (
    <div>
      <h1>Penultima</h1>

      <div id="wrapper">
        <div className="box">
          <span className="tab" onClick={() => setTab('about')}>about</span> |{' '}
          <span className="tab" onClick={() => setTab('register')}>register</span> |{' '}
          <span className="tab" onClick={() => setTab('login')}>login</span> |{' '}
          <span className="tab" onClick={() => setTab('logout')}>logout</span>
        </div>

        <div className="box-no-border">
          {tab === 'about' &&
            <div>
              Would you like to play?
            </div>
          }

          {tab === 'register' && <h2>Register</h2>}
          {tab === 'login' && <h2>Login</h2>}
          {(tab === 'register' || tab === 'login') &&
            <>
              <input type="text" placeholder="Username" onChange={e => setUsername(e.target.value)} value={username} />
              <input type="password" placeholder="Password" onChange={e => setPassword(e.target.value)} value={password} />
            </>
          }
          {tab === 'register' && <input type="submit" value="Register" onClick={handleSubmit} />}
          {(tab === 'register' && registerErrorMessage) && <div className="box">{registerErrorMessage}</div>}

          {tab === 'login' && <input type="submit" value="Login" onClick={handleSubmit} />}
          {(tab === 'login' && loginErrorMessage) && <div className="box">{loginErrorMessage}</div>}

          {tab === 'logout' && <input type="submit" value="Logout" onClick={handleSubmit} />}
                         
        </div>
      </div>
    </div>
  )
}  

