package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a Mux
	mux := http.NewServeMux()
	server := NewPlayerServer()
	mux.HandleFunc("/player/", server.playerHandler)
	mux.HandleFunc("/tag/", server.tagHandler)
	mux.HandleFunc("/created/", server.createdHandler)

	log.Fatal(http.ListenandServe("localhost:"+os.GetEnv("SERVERPORT"), mux))
}
