package simulation

import (
	"math"
	"math/rand"
)

func dotProduct(v1, v2 []float64) float64 {
	res := 0.0
	for i := range v1 {
		res += v1[i] * v2[i]
	}
	return res
}

func randomColor() string {
	colors := []string{"#fabd2f", "#b16286", "#458588", "#ebdbb2"}
	ind := rand.Intn(len(colors))
	return colors[ind]
}

func (s *Simulation) RandomizeColors() {
	for i := range s.Balls {
		s.Balls[i].Color = randomColor()
	}
}

func randomDirection() float64 {
	if rand.Intn(2) == 0 {
		return -1.0
	}
	return 1.0
}

func (s *Simulation) GenerateBalls(n int) {
	balls := make([]Ball, n)
	for i := range balls {
		radius := math.Max(1, rand.Float64()*2) * s.SizeMultiplier
		center := Point{X: math.Max(50, rand.Float64()*750), Y: math.Max(50, rand.Float64()*550)}
		velocity := Point{X: randomDirection() * rand.Float64() * 3, Y: randomDirection() * rand.Float64() * 3}
		balls[i] = Ball{
			R:        radius,
			Center:   center,
			Velocity: velocity,
			Color:    randomColor(),
		}
	}
	s.Balls = balls
}

func (s *Simulation) UpdateBallSize() {
	for i := range s.Balls {
		s.Balls[i].R *= s.SizeMultiplier / s.OldSizeMultiplier
	}
}

func (s *Simulation) Reset() {
	s.Paused = true
	s.War = false
	s.GenerateBalls(100)
	s.Gravity = 0
}

func (pi *Pillision) Reset(velocity float64) {
	pi.Paused = true
	pi.NumCollisions = 0
	pi.SmallSquare.TopLeft.X = 100
	pi.SmallSquare.BottomRight.X = 200
	pi.SmallSquare.Velocity = 0
	pi.BigSquare.TopLeft.X = 300
	pi.BigSquare.BottomRight.X = 500
	pi.BigSquare.Velocity = velocity
}
