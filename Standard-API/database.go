package main

import "time"

type Player struct {
	ID      int       `json:"id"`
	Text    string    `json:"text"`
	Tags    []string  `json:"tags"`
	Created time.Time `json:"created"`
}

type playerServer struct {
	store *playerstore.PlayerStore
}

func NewPlayerServer() *playerServer {
	store := playerstore.New()
	return &playerServer{store: store}
}

func New() *PlayerStore

// Create a new Player
func (ps *PlayerStore) CreatePlayer(text string, tags []string, created time.Time) int

// Get a PLayer by ID. Returns nil if not found
func (ps *PlayerStore) GetPlayer(id int) (Player, error)

// Get all Players
func (ps *PlayerStore) GetAllPlayers() []Player

// Delete a Player with a given ID.
func (ps *PlayerStore) DeletePlayer(id int) error

// Delete all Players
func (ps *PlayerStore) DeletePlayers() error

// Get all Players with a given tag
func (ps *PlayerStore) GetPlayersByTag(tag string) []Player

// Get all Players by the time they were created
func (ps *PlayerStore) GetPlayersByCreated(created time.Time) []Player
