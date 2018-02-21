// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"time"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Scenes and clients who belong to it
	scenes map[string][]*Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		scenes:     make(map[string][]*Client),
	}
}

func (h *Hub) run() {
	ticker := time.NewTicker(time.Duration(*updateRate) * time.Millisecond)
	garbageTicker := time.NewTicker(time.Duration(*collectionRate) * time.Millisecond)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.scenes[client.scene] = append(h.scenes[client.scene], client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			playerMsg := PlayerMsg{}
			json.Unmarshal(message, &playerMsg)
			log.Println("Received msg from the client:", playerMsg)
			setDirection(playerMsg.ID, playerMsg.Direction)
		case <-ticker.C:
			for scene, clients := range h.scenes {
				posUpdates := move(scene)
				if len(posUpdates) > 0 {
					updateMsg := map[string][]*Position{
						"updates": posUpdates,
					}
					serializedMsg, _ := json.Marshal(updateMsg)
					for _, client := range clients {
						client.send <- serializedMsg
					}
				}
			}
		case <-garbageTicker.C:
			deadScenes := make([]string, 0)
			for scene, clients := range h.scenes {
				dead := true
				for _, client := range clients {
					if h.clients[client] {
						dead = false
						break
					}
				}
				if dead {
					deadScenes = append(deadScenes, scene)
				}
			}
			for _, deadScene := range deadScenes {
				delete(h.scenes, deadScene)

			}
		}
	}
}

type PlayerMsg struct {
	ID        string    `json: "id"`
	Direction Direction `json: "direction"`
}
