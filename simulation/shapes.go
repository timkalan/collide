package simulation

import (
	"sync"
)

type Ball struct {
	R     float64 `json:"r"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	VX    float64 `json:"vx"`
	VY    float64 `json:"vy"`
	Color string  `json:"color"`
}

type Simulation struct {
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Gravity float64 `json:"gravity"`
	Balls   []Ball  `json:"balls"`
	Paused  bool    `json:"paused"`
	War     bool    `json:"war"`
	FPS     float64 `json:"fps"`
	Mu      sync.Mutex
}
