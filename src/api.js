// api/index.js

// var socket = new WebSocket("ws://192.168.0.151:8084/ws");

import { host } from './settings';

var socket = new WebSocket('ws://'+host+':8084/ws'); 

let connect = ({messageCallback}) => {  

  // console.log('connect',socket.readyState)
  
  /*
  if ( socket.readyState === WebSocket.CLOSED)
    console.log('socket is closed')
  else if ( socket.readyState === WebSocket.CLOSING)
    console.log('socket is closing')
  else if ( socket.readyState === WebSocket.OPEN)
    console.log('socket is open')
  else if ( socket.readyState === WebSocket.CONNECTING)
    console.log('socket is connecting') 
  else 
    console.log('socket is unknown')
  */
 
  socket.onopen = () => {
    console.log("onopen","Successfully Connected");
    // requestMove(13)
  };

  socket.onmessage = msg => {            
    const message = JSON.parse(msg.data)        
    messageCallback(message)     
  };

  socket.onclose = event => {
    console.log('onclose',"Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log('onerror',"Socket Error: ", error);
  };
};

let sendMsg = msg => {  
  const sendable = JSON.stringify({
    "message": msg
  })
  socket.send(sendable);
};

let requestMove = move => {  
  if (socket.readyState === WebSocket.CLOSING || socket.readyState === WebSocket.CLOSED) {
    console.log('Socket is closed/closing. Cannot send move')   
    alert("Connection to the game has closed! Adventure awaits, try reloading ↻ this page.")
  }
  socket.send(JSON.stringify({
    "move": move
  }));
  return "move successful"
};

let requestUpdateAvatar = id => { 
  socket.send(JSON.stringify({
    "avatar": {id: parseInt(id)}
  }));
};

let requestPeerGem = move => {  
  socket.send(JSON.stringify({
    "peerGem": {}
  }));
};

export { connect, sendMsg, requestMove, requestUpdateAvatar, requestPeerGem  };