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

func GenerateBalls(n int) []Ball {
	balls := make([]Ball, n)
	centers := make([][]float64, n)
	for i := range balls {
		center := []float64{math.Max(50, rand.Float64()*750), math.Max(50, rand.Float64()*550)}
		radius := math.Max(10, rand.Float64()*20)
		centers[i] = center
		balls[i] = Ball{
			R:     radius,
			X:     center[0],
			Y:     center[1],
			VX:    randomDirection() * rand.Float64() * 3,
			VY:    randomDirection() * rand.Float64() * 3,
			Color: randomColor(),
		}
	}
	return balls
}
