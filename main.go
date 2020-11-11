package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	s := &http.Server{
		Addr:    ":8080",
		Handler: Handler{},
	}
	fmt.Printf("starting server on %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type Handler struct {
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("%s %q", r.Method, r.RequestURI)
	SetCommonHeaders(w)
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <title> testing websockets </title>
  </head>
  <body>
    <h1> websocket test </h1>
  </body>
<script>


var exampleSocket = new WebSocket("ws://"+ location.host +"/t")

exampleSocket.onopen = function (event) {
console.log("open");
};


exampleSocket.onmessage = function (event) {
  console.log("got: " + event.data);
  exampleSocket.send(JSON.stringify("thanks for " + event.data));

}


</script>
</html>
`)

	case "/t":
		websocket.Handler(h.ServeWebsocket).ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h Handler) ServeWebsocket(ws *websocket.Conn) {
	for {
		t := time.Now()
		log.Printf("sending: %v\n", t)
		if err := websocket.JSON.Send(ws, t); err != nil {
			log.Printf("send error: %v", err)
			return
		}
		var data interface{}
		if err := websocket.JSON.Receive(ws, &data); err != nil {
			log.Printf("receive error: %v", err)
			return
		}
		log.Printf("got: %v\n", data)
		time.Sleep(time.Second)
	}
}

func SetCommonHeaders(w http.ResponseWriter) {
	h := w.Header()
	h.Add("Access-Control-Allow-Origin", "*")
	h.Add("Referrer-Policy", "no-referrer")
	h.Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	h.Add("X-Content-Type-Options", "nosniff")
	h.Add("X-Frame-Options", "SAMEORIGIN")
	h.Add("X-Permitted-Cross-Domain-Policies", "none")
	h.Add("X-XSS-Protection", "1; mode=block")
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Pragma", "no-cache")
	h.Set("Expires", "0")
}
