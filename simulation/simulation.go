package simulation

import (
	"math"
	"math/rand"
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
}

func dotProduct(v1, v2 []float64) float64 {
	res := 0.0
	for i := range v1 {
		res += v1[i] * v2[i]
	}
	return res
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
        balls[i].X = balls[j].X + (balls[i].R + balls[j].R) * (balls[i].X - balls[j].X) / dist
        balls[i].Y = balls[j].Y + (balls[i].R + balls[j].R) * (balls[i].Y - balls[j].Y) / dist
        
				// if balls[i].R > balls[j].R {
				// 	balls[j].Color = balls[i].Color
				// } else {
				// 	balls[i].Color = balls[j].Color
				// }
			}
		}
	}
}

func (s *Simulation) wallCollisionDetection(ball Ball) {
		if ball.X-ball.R <= 0 {
			ball.VX = -ball.VX
			ball.X = ball.R
		}
		if ball.X+ball.R >= s.Width {
			ball.VX = -ball.VX
			ball.X = s.Width - ball.R
		}
		if ball.Y-ball.R <= 0 {
			ball.VY = -ball.VY
			ball.Y = ball.R
		}
		if ball.Y+ball.R >= s.Height {
			ball.VY = -ball.VY
			ball.Y = s.Height - ball.R
		}
}

func (s *Simulation) Update(dt float64) {
	for i := range s.Balls {
		s.Balls[i].X += s.Balls[i].VX * dt
		s.Balls[i].Y += s.Balls[i].VY * dt

		// Bounce off the walls
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

		// Apply gravity
		if s.Gravity > 0 {
			s.Balls[i].VY += s.Gravity
		}
	}

	dumbCollisionDetection(s.Balls)
}

func randomColor() string {
	colors := []string{"#fabd2f", "#b16286", "#458588", "#ebdbb2"}
	ind := rand.Intn(len(colors))
	return colors[ind]
}

func GenerateBalls(n int) []Ball {
	balls := make([]Ball, n)
	centers := make([][]float64, n)
	for i := range balls {
		center := []float64{math.Max(50, rand.Float64()*750), math.Max(50, rand.Float64()*550)}
		radius := math.Max(15, rand.Float64()*30)
		if i == 0 {
		} else {
			for j := range i {
				for math.Sqrt(math.Pow(center[0]-centers[j][0], 2)+math.Pow(center[1]-centers[j][1], 2)) < radius+balls[j].R {
					center[0] = math.Max(50, rand.Float64()*750)
					center[1] = math.Max(50, rand.Float64()*550)
				}
			}
		}
		centers[i] = center
		balls[i] = Ball{
			R:     radius,
			X:     center[0],
			Y:     center[1],
			VX:    rand.Float64() * 3,
			VY:    rand.Float64() * 3,
			Color: randomColor(),
		}
	}
	return balls
}
