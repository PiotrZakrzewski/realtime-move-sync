package main

import "math"
import "log"

func isColliding(x1 float64, y1 float64, pos *Position, boundingSquare float64) bool {
	var x2 = pos.X
	var w2 = boundingSquare
	var w1 = boundingSquare
	var y2 = pos.Y
	var h2 = boundingSquare
	var h1 = boundingSquare
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && h1+y1 > y2
}

func canMove(posFrom *Position, x float64, y float64, otherPositions []*Position, boundingSquare float64) bool {
	for _, pos := range otherPositions {
		if pos == nil {
			log.Print("nil pointer in positions array")
			continue
		}
		if posFrom == pos {
			continue
		}
		if isColliding(x, y, pos, boundingSquare) {
			log.Print("Truely collides")
			return false
		}
	}
	return true
}

const eps = 0.1

func floatEqual(val1 float64, val2 float64) bool {
	if math.Abs(val1-val2) < eps {
		return true
	} else {
		return false
	}
}
