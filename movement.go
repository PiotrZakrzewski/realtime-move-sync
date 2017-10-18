package main

import (
	"log"
	"math"
	"time"
)

var movements = make(map[string]*Direction)
var positions = make(map[string]*Position)

const UP = 1
const DOWN = 2
const LEFT = 3
const RIGHT = 4
const STOP = 5

type Direction struct {
	Forward int8 `json:"forward"`
	Angular int8 `json:"angular"`
}

func setDirection(character string, newDirection Direction) {
	movements[character] = &newDirection
}

func toRadians(degree float64) float64 {
	radians := (math.Pi * degree) / 180
	return radians
}

const moveOffset = 5
const turningOffset = 0.3926991 // 22.5 degrees or 1/16th of a full turn

func move() []*Position {
	now := time.Now()
	timeStamp := now.UnixNano() / 1000000
	updates := make([]*Position, len(movements))
	i := 0
	for id, dir := range movements {
		pos := positions[id]
		if pos == nil {
			log.Println("Tried to move id without position: " + id)
			continue
		}
		var xOffset float64
		var yOffset float64
		var angularOffset float64
		if dir.Forward == UP {
			xOffset = math.Cos(pos.Direction) * -moveOffset
			yOffset = math.Sin(pos.Direction) * moveOffset
		} else if dir.Forward == DOWN {
			xOffset = math.Cos(pos.Direction) * moveOffset
			yOffset = math.Sin(pos.Direction) * -moveOffset
		}
		if dir.Angular == LEFT {
			angularOffset = -turningOffset
		} else if dir.Angular == RIGHT {
			angularOffset = turningOffset
		}
		newX := pos.X + xOffset
		newY := pos.Y + yOffset
		pos.X = newX
		pos.Y = newY
		pos.Time = timeStamp
		pos.Direction += angularOffset
		updates[i] = pos
		i++
	}
	return updates
}

type Position struct {
	ID        string  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction float64 `json:"direction"`
	Time      int64   `json:"time"`
}
