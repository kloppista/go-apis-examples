package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (ps *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/player/" {
		switch r.Method {
		case http.MethodPost:
			ps.processPost(w, r)
		case http.MethodGet:
			ps.processGet(w, r)
		case http.MethodDelete:
			ps.processDelete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		// Request has an ID in the path
		path := strings.Trim(r.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /player(<id> in player handler", http.StatusBadRequest)
			return
	}
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodDelete {
		ps.deletePlayerHandler(w, r, id)
	} else if r.Method == http.MethodGet {
		ps.getPlayerHandler(w, r, id)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (ps *PlayerServer) getAllPlayersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting all players at %s", r.URL.Path)

	allPlayers := ps.store.GetAllPlayers()
	js, err := json.Marshal(allPlayers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
