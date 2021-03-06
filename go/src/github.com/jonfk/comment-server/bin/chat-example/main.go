package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"text/template"
)

var (
	addr         = flag.String("addr", ":8080", "http service address")
	homeTemplate = template.Must(template.ParseFiles(os.Getenv("PATH_DEBUG_PATH")))
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"context": "serveHome",
		"url":     r.URL,
		"method":  r.Method,
	}).Info("Received request")
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, r.Host)
}

func init() {
	flag.Parse()
}

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.WithFields(log.Fields{
		"context": "main",
		"addr":    *addr,
	}).Info("http Listening and Serving")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
