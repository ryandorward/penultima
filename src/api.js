// api/index.js

var socket = new WebSocket("ws://192.168.0.151:8083/ws");

let connect = ({messageCallback, moveCallback, updateWorldViewCallback}) => {
  // console.log("connect");

  socket.onopen = () => {
    console.log("Successfully Connected");

    requestMove(13)       

  };

  socket.onmessage = msg => {            
    const message = JSON.parse(msg.data)    
     // what kind of message? 
    switch (message.type) {
      case 1: 
        messageCallback(message);
        break;
      case 2:
        // console.log("User message received:",message)
        break;
      case 3: 
        moveCallback(message)       
        break;
      case 4: 
        updateWorldViewCallback(message)       
        break;
      default:
        console.log("Didn't catch message type",message)
    }
        
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {  
  const sendable = JSON.stringify({
    "message": msg
  })
  socket.send(sendable);
};

let requestMove = keycode => {   
  // console.log("requesting move") ;
  socket.send(JSON.stringify({
    "move": keycode
  }));
};

export { connect, sendMsg, requestMove };