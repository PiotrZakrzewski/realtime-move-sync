package main

import (
	"flag"
	"log"
	random "math/rand"
	"net/http"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var updateRate = flag.Int("rate", 50, "Update rate for the game loop in milliseconds")
var collectionRate = flag.Int("collection", 1000, "Interval at which empty scenes are collected")
var moveOffset = flag.Float64("speed", 10, "Movement speed of characters")
var collisionDetection = flag.Bool("collisions", false, "Collisions detection on/off")
var boundingSquare = flag.Int("bound", 20, "Bounding square size")
var interpolation = flag.Bool("interpolation", false, "Use position interpolation and buffering")

func main() {
	flag.Parse()
	go http.HandleFunc("/", serveHome)
	hub := newHub()
	go hub.run()
	bot(float64(random.Int()%500), float64(random.Int()%500))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "index.html")
	} else if r.URL.Path == "/assets/player.png" {
		http.ServeFile(w, r, "assets/player.png")
	} else {
		http.Error(w, "Not found", 404)
		return
	}
}
func bot(x1 float64, y1 float64) {
	now := time.Now()
	epoch := now.UnixNano() / 1000000
	botUUID := pseudoUUID()
	pos := &Position{Scene: "a", ID: botUUID, X: x1, Y: y1, Direction: 0.0, Time: epoch}
	positions[botUUID] = pos
	positionsArray = append(positionsArray, pos)
	dir := Direction{Forward: UP, Angular: RIGHT}
	setDirection(botUUID, dir)
}
