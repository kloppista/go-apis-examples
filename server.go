package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Player struct {
	Name     string `json:"name"`
	Team     string `json:"team"`
	ID       string `json:"id"`
	Position string `json:"position"`
	Age      int    `json:"age"`
}

type playerHandlers struct {
	sync.Mutex
	store map[string]Player
}

func (h *playerHandlers) players(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.POST(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

func (h *playerHandlers) get(w http.ResponseWriter, r *http.Request) {
	players := make([]Player, len(h.store))

	h.Lock()
	i := 0
	for _, player := range h.store {
		players[i] = player
		i++
	}

	h.Unlock()
	jsonBytes, err := json.Marshal(players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Handler for random selection
func (h *playerHandlers) getRandomPlayer(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()

	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}

	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header.Add("location", fmt.Sprintf("/players/%s"), target)
	w.WriteHeader(http.StatusFound)
}

// Get a random player
func (h *playerHandlers) getPlayer(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomPlayer(w, r)
		return
	}

	h.Lock()
	player, ok := h.store[parts[2]]

	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *playerHandlers) POST(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Check for the content type
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type `application/json`, but got `%s`", ct)))
		return
	}

	// Unmarshal Body Bytes
	var player Player
	err = json.Unmarshal(bodyBytes, &player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Assign ID
	player.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.store[player.ID] = player
	defer h.Unlock()
}

func newPlayerHandlers() *playerHandlers {
	return &playerHandlers{
		store: map[string]Player{},
	}
}

// Defining admin struct
type adminPortal struct {
	password string
}

// Authentication
func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

// Handler for adminPortal
func (a adminPortal) handler(w http.ResponseWriter, r http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized"))
	}
	return

	w.Write([]byte("<html><h1>Admin Portal</h1></html>"))
}

func main() {
	// Admin
	admin := newAdminPortal()
	playerHandlers := newPlayerHandlers()
	// Handlers
	http.HandleFunc("/players", playerHandlers.players)
	http.HandleFunc("/players/", playerHandlers.getPlayer)
	http.HandleFunc("/admin/", admin.Handler)
	// Server Setup
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
