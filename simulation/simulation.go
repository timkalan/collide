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

func (s *Simulation) collisionResponse(ball1, ball2 *Ball) {

	dx := ball1.X - ball2.X
	dy := ball1.Y - ball2.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	v12 := []float64{ball1.VX - ball2.VX, ball1.VY - ball2.VY}
	v21 := []float64{ball2.VX - ball1.VX, ball2.VY - ball1.VY}
	c12 := []float64{ball1.X - ball2.X, ball1.Y - ball2.Y}
	c21 := []float64{ball2.X - ball1.X, ball2.Y - ball1.Y}

	if dist < ball1.R+ball2.R {
		massConst1 := 2 * ball2.R / (ball1.R + ball2.R)
		massConst2 := 2 * ball1.R / (ball1.R + ball2.R)
		const1 := dotProduct(v12, c12) / dotProduct(c12, c12)
		const2 := dotProduct(v21, c21) / dotProduct(c21, c21)

		ball1.VX = ball1.VX - massConst1*const1*c12[0]
		ball1.VY = ball1.VY - massConst1*const1*c12[1]
		ball2.VX = ball2.VX - massConst2*const2*c21[0]
		ball2.VY = ball2.VY - massConst2*const2*c21[1]

		// Ensure there is no clipping
		ball1.X = ball2.X + (ball1.R+ball2.R)*(ball1.X-ball2.X)/dist
		ball1.Y = ball2.Y + (ball1.R+ball2.R)*(ball1.Y-ball2.Y)/dist

		if s.War {
			if ball1.R > ball2.R {
				ball2.Color = ball1.Color
			} else {
				ball1.Color = ball2.Color
			}
		}
	}
}

func (s *Simulation) dumbCollisionDetection() {
	for i := 0; i < len(s.Balls); i++ {
		for j := i + 1; j < len(s.Balls); j++ {
			s.collisionResponse(&s.Balls[i], &s.Balls[j])
		}
	}
}

func (s *Simulation) sweepAndPruneCollisionDetection() {
	sort.Slice(s.Balls, func(i, j int) bool {
		return s.Balls[i].X < s.Balls[j].X
	})

	var active []*Ball
	active = append(active, &s.Balls[0])

	for i := range s.Balls[1:] {
		active = append(active, &s.Balls[i+1])

		for j := range active {
			for l := range active[j+1:] {
				k := l + j + 1
				s.collisionResponse(active[j], active[k])
			}
		}
		// update active slice
		for j := range active {
			if active[j].X+active[j].R < s.Balls[i].X {
				active = append(active[:j], active[j+1:]...)
			} else {
				break
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
	s.sweepAndPruneCollisionDetection()
	s.dumbCollisionDetection()
}
