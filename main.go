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
		Handler: http.HandlerFunc(ServeHTTP),
	}
	fmt.Printf("starting server on %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("%s %q", r.Method, r.RequestURI)
	SetCommonHeaders(w)
	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "index.html")
	case "/script.js":
		http.ServeFile(w, r, "script.js")
	case "/t":
		websocket.Handler(ServeWebsocket).ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

func ServeWebsocket(ws *websocket.Conn) {
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
