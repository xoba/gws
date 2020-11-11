var exampleSocket = new WebSocket("ws://"+ location.host +"/t")

exampleSocket.onopen = function (event) {
    console.log("open");
};


exampleSocket.onmessage = function (event) {
    console.log("got: " + event.data);
    exampleSocket.send(JSON.stringify("thanks for " + event.data));
}
