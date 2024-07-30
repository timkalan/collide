package simulation

import (
	"math"
	"sort"
)

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

func (s *Simulation) dumbCollisionDetection() {
	for i := 0; i < len(s.Balls); i++ {
		for j := i + 1; j < len(s.Balls); j++ {
			dx := s.Balls[i].X - s.Balls[j].X
			dy := s.Balls[i].Y - s.Balls[j].Y
			dist := math.Sqrt(dx*dx + dy*dy)
			v12 := []float64{s.Balls[i].VX - s.Balls[j].VX, s.Balls[i].VY - s.Balls[j].VY}
			v21 := []float64{s.Balls[j].VX - s.Balls[i].VX, s.Balls[j].VY - s.Balls[i].VY}
			c12 := []float64{s.Balls[i].X - s.Balls[j].X, s.Balls[i].Y - s.Balls[j].Y}
			c21 := []float64{s.Balls[j].X - s.Balls[i].X, s.Balls[j].Y - s.Balls[i].Y}

			if dist < s.Balls[i].R+s.Balls[j].R {
				massConst1 := 2 * s.Balls[j].R / (s.Balls[i].R + s.Balls[j].R)
				massConst2 := 2 * s.Balls[i].R / (s.Balls[i].R + s.Balls[j].R)
				const1 := dotProduct(v12, c12) / dotProduct(c12, c12)
				const2 := dotProduct(v21, c21) / dotProduct(c21, c21)

				s.Balls[i].VX = s.Balls[i].VX - massConst1*const1*c12[0]
				s.Balls[i].VY = s.Balls[i].VY - massConst1*const1*c12[1]
				s.Balls[j].VX = s.Balls[j].VX - massConst2*const2*c21[0]
				s.Balls[j].VY = s.Balls[j].VY - massConst2*const2*c21[1]

				// Ensure there is no clipping
				s.Balls[i].X = s.Balls[j].X + (s.Balls[i].R+s.Balls[j].R)*(s.Balls[i].X-s.Balls[j].X)/dist
				s.Balls[i].Y = s.Balls[j].Y + (s.Balls[i].R+s.Balls[j].R)*(s.Balls[i].Y-s.Balls[j].Y)/dist

				if s.War {
					if s.Balls[i].R > s.Balls[j].R {
						s.Balls[j].Color = s.Balls[i].Color
					} else {
						s.Balls[i].Color = s.Balls[j].Color
					}
				}
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
				numActive++
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

	s.dumbCollisionDetection()
}
