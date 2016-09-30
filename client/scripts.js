
var exampleSocket = new WebSocket("ws://127.0.0.1:3012");

exampleSocket.onopen = function (event) {
  //exampleSocket.send("Here's some text that the server is urgently awaiting!");
  console.log("opened websocket: " + exampleSocket);
  exampleSocket.send("hello world");
};

exampleSocket.onmessage = function (event) {
  console.log("Received message: " + event.data);
}
