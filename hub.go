// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
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
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	ticker := time.NewTicker(time.Duration(*updateRate) * time.Millisecond)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			playerMsg := PlayerMsg{}
			json.Unmarshal(message, playerMsg)
			setDirection(playerMsg.ID, playerMsg.Direction)
		case <-ticker.C:
			posUpdates := move()
			serializedMsg, _ := json.Marshal(posUpdates)
			for client := range h.clients {
				client.send <- serializedMsg
			}
		}
	}
}

type PlayerMsg struct {
	ID        string    `json: "id"`
	Direction Direction `json: "direction"`
}
