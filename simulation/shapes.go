package simulation

import (
	"sync"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Ball struct {
	R     float64 `json:"r"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	VX    float64 `json:"vx"`
	VY    float64 `json:"vy"`
	Color string  `json:"color"`
}

type Square struct {
	TopLeft     Point   `json:"topLeft"`
	BottomRight Point   `json:"bottomRight"`
	Velocity    float64 `json:"velocity"`
	Weight      float64 `json:"weight"`
  Width       float64
}

type Simulation struct {
	Width             float64 `json:"width"`
	Height            float64 `json:"height"`
	Gravity           float64 `json:"gravity"`
	SizeMultiplier    float64 `json:"size"`
	Balls             []Ball  `json:"balls"`
	Paused            bool    `json:"paused"`
	War               bool    `json:"war"`
	FPS               float64 `json:"fps"`
	Mu                sync.Mutex
	OldSizeMultiplier float64
}

type Pillision struct {
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	WeightMultiplier float64 `json:"multiplier"`
	NumCollisions    int     `json:"numCollisions"`
	SmallSquare      Square  `json:"smallSquare"`
	BigSquare        Square  `json:"bigSquare"`
	Paused           bool    `json:"paused"`
	Mu               sync.Mutex
}
