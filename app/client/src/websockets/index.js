const WS_HOST = 'localhost:8080/ws';
const WS_URL = "ws://" + WS_HOST;

var socket = new WebSocket(WS_URL);

let connect = onMessageReceive => {
  console.log("connecting");

  socket.onopen = (e) => {
    console.log("Successfully Connected");
  };

  socket.onmessage = e => {
    const message = JSON.parse(e.data);
    onMessageReceive(message);
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMessage = message => {
  console.log("sending msg: ", message);
  socket.send(message);
};

export { connect, sendMessage };