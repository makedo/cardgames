import User from "../service/User";

const WS_HOST = 'localhost:8080/ws';
const WS_URL = "ws://" + WS_HOST;

var socket = null; 
function getSocket(params = {}) {

  let qs = Object.keys(params)
    .map(key => `${key}=${params[key] === null || typeof params[key] === 'undefined' ? '' : params[key]}`)
    .join('&');
  
  if (qs) {
    qs = '?' + qs
  }

  if (null === socket) {
    socket = new WebSocket(WS_URL + qs);
  }

  return socket;
}

export function connect(onMessage, onOpen, params = {}) {
  console.log("connecting");
  socket = getSocket(params)

  socket.onopen = (e) => {
    console.log("Successfully Connected");
    onOpen()
  };

  socket.onmessage = e => {
    const message = JSON.parse(e.data);
    onMessage(message);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

export function send(message) {
  
  if (typeof message !== 'string') {
    message = JSON.stringify(message)
  }
  console.log("MESSAGE_SENT");
  console.log(message);
  console.log(User.getId());
  console.log(typeof User.getId());

  return socket.send(message);
}