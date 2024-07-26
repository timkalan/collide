package simulation

import (
	"math"
	"math/rand"
	"sort"
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
	Mu      sync.Mutex
}

func dotProduct(v1, v2 []float64) float64 {
	res := 0.0
	for i := range v1 {
		res += v1[i] * v2[i]
	}
	return res
}

func (s *Simulation) wallCollisionDetection(i int) {
	if s.Balls[i].X-s.Balls[i].R <= 0 {
		s.Balls[i].VX = -s.Balls[i].VX
		s.Balls[i].X = s.Balls[i].R
	}
	if s.Balls[i].X+s.Balls[i].R >= s.Width {
		s.Balls[i].VX = -s.Balls[i].VX
		s.Balls[i].X = s.Width - s.Balls[i].R
	}
	if s.Balls[i].Y-s.Balls[i].R <= 0 {
		s.Balls[i].VY = -s.Balls[i].VY
		s.Balls[i].Y = s.Balls[i].R
	}
	if s.Balls[i].Y+s.Balls[i].R >= s.Height {
		s.Balls[i].VY = -s.Balls[i].VY
		s.Balls[i].Y = s.Height - s.Balls[i].R
	}
}

func dumbCollisionDetection(balls []Ball) {
	for i := 0; i < len(balls); i++ {
		for j := i + 1; j < len(balls); j++ {
			dx := balls[i].X - balls[j].X
			dy := balls[i].Y - balls[j].Y
			dist := math.Sqrt(dx*dx + dy*dy)
			v12 := []float64{balls[i].VX - balls[j].VX, balls[i].VY - balls[j].VY}
			v21 := []float64{balls[j].VX - balls[i].VX, balls[j].VY - balls[i].VY}
			c12 := []float64{balls[i].X - balls[j].X, balls[i].Y - balls[j].Y}
			c21 := []float64{balls[j].X - balls[i].X, balls[j].Y - balls[i].Y}

			if dist < balls[i].R+balls[j].R {
				massConst1 := 2 * balls[j].R / (balls[i].R + balls[j].R)
				massConst2 := 2 * balls[i].R / (balls[i].R + balls[j].R)
				const1 := dotProduct(v12, c12) / dotProduct(c12, c12)
				const2 := dotProduct(v21, c21) / dotProduct(c21, c21)

				balls[i].VX = balls[i].VX - massConst1*const1*c12[0]
				balls[i].VY = balls[i].VY - massConst1*const1*c12[1]
				balls[j].VX = balls[j].VX - massConst2*const2*c21[0]
				balls[j].VY = balls[j].VY - massConst2*const2*c21[1]

				// Ensure there is no clipping
				balls[i].X = balls[j].X + (balls[i].R+balls[j].R)*(balls[i].X-balls[j].X)/dist
				balls[i].Y = balls[j].Y + (balls[i].R+balls[j].R)*(balls[i].Y-balls[j].Y)/dist

				// if balls[i].R > balls[j].R {
				// 	balls[j].Color = balls[i].Color
				// } else {
				// 	balls[i].Color = balls[j].Color
				// }
			}
		}
	}
}

func sweepAndPruneCollisionDetection(balls []Ball) {
	sort.Slice(balls, func(i, j int) bool {
		return balls[i].X < balls[j].X
	})

	active := make([]Ball, len(balls)/10)
	numActive := 0

	for i := range balls {
		for j := range active {
			dx := balls[i].X - active[j].X
			dy := balls[i].Y - active[j].Y
			dist := math.Sqrt(dx*dx + dy*dy)
			v12 := []float64{balls[i].VX - active[j].VX, balls[i].VY - active[j].VY}
			v21 := []float64{balls[j].VX - balls[i].VX, active[j].VY - balls[i].VY}
			c12 := []float64{balls[i].X - active[j].X, balls[i].Y - active[j].Y}
			c21 := []float64{balls[j].X - balls[i].X, active[j].Y - balls[i].Y}

			if dist < balls[i].R+active[j].R {
				massConst1 := 2 * active[j].R / (balls[i].R + active[j].R)
				massConst2 := 2 * balls[i].R / (balls[i].R + active[j].R)
				const1 := dotProduct(v12, c12) / dotProduct(c12, c12)
				const2 := dotProduct(v21, c21) / dotProduct(c21, c21)

				balls[i].VX = balls[i].VX - massConst1*const1*c12[0]
				balls[i].VY = balls[i].VY - massConst1*const1*c12[1]
				active[j].VX = active[j].VX - massConst2*const2*c21[0]
				active[j].VY = active[j].VY - massConst2*const2*c21[1]

				// Ensure there is no clipping
				balls[i].X = active[j].X + (balls[i].R+active[j].R)*(balls[i].X-active[j].X)/dist
				balls[i].Y = active[j].Y + (balls[i].R+active[j].R)*(balls[i].Y-active[j].Y)/dist

        // Update active
        numActive ++
        active[numActive] = balls[i]
			}
		}
	}
}

func (s *Simulation) Update(dt float64) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for i := range s.Balls {
		// Adjust position
		s.Balls[i].X += s.Balls[i].VX * dt
		s.Balls[i].Y += s.Balls[i].VY * dt

		// Apply gravity
		if s.Gravity > 0 {
			s.Balls[i].VY += s.Gravity
		}

		s.wallCollisionDetection(i)
	}

	dumbCollisionDetection(s.Balls)
}

func randomColor() string {
	colors := []string{"#fabd2f", "#b16286", "#458588", "#ebdbb2"}
	ind := rand.Intn(len(colors))
	return colors[ind]
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
