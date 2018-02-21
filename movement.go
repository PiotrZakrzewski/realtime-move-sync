package main

import (
	"log"
	"math"
	"time"
)

var movements = make(map[string]*Direction)
var positions = make(map[string]*Position)
var positionsArray = make([]*Position, 0)

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

const turningOffset = 0.3926991 // 22.5 degrees or 1/16th of a full turn

func move(scene string) []*Position {
	now := time.Now()
	timeStamp := now.UnixNano() / 1000000
	if !*interpolation {
		timeStamp = 0
	}
	updates := make([]*Position, 0)
	for id, dir := range movements {
		pos := positions[id]
		if pos.Scene != scene {
			continue
		}
		if pos == nil {
			log.Println("Tried to move id without position: " + id)
			continue
		}
		var xOffset float64
		var yOffset float64
		var angularOffset float64
		if dir.Angular == LEFT {
			angularOffset = -turningOffset
		} else if dir.Angular == RIGHT {
			angularOffset = turningOffset
		} else {
			angularOffset = 0.0
		}
		pos.Direction += angularOffset
		if dir.Forward == UP {
			yOffset = math.Cos(pos.Direction) * -*moveOffset
			xOffset = math.Sin(pos.Direction) * *moveOffset
		} else if dir.Forward == DOWN {
			yOffset = math.Cos(pos.Direction) * *moveOffset
			xOffset = math.Sin(pos.Direction) * -*moveOffset
		}
		if *collisionDetection && !canMove(pos, pos.X+xOffset, pos.Y+yOffset, positionsArray, float64(*boundingSquare)) {
			continue
		}
		pos.X += xOffset
		pos.Y += yOffset
		pos.Time = timeStamp
		updates = append(updates, pos)
	}
	return updates
}

type Position struct {
	ID        string  `json:"id"`
	Scene     string  `json:"scene"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Direction float64 `json:"direction"`
	Time      int64   `json:"time"`
}
