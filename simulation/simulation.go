package simulation

import (
	"math"
	"sort"
)

func (s *Simulation) wallCollisionDetection(i int) {
	if s.Balls[i].Center.X-s.Balls[i].R <= 0 {
		s.Balls[i].Velocity.X = -s.Balls[i].Velocity.X
		s.Balls[i].Center.X = s.Balls[i].R
	}
	if s.Balls[i].Center.X+s.Balls[i].R >= s.Width {
		s.Balls[i].Velocity.X = -s.Balls[i].Velocity.X
		s.Balls[i].Center.X = s.Width - s.Balls[i].R
	}
	if s.Balls[i].Center.Y-s.Balls[i].R <= 0 {
		s.Balls[i].Velocity.Y = -s.Balls[i].Velocity.Y
		s.Balls[i].Center.Y = s.Balls[i].R
	}
	if s.Balls[i].Center.Y+s.Balls[i].R >= s.Height {
		s.Balls[i].Velocity.Y = -s.Balls[i].Velocity.Y
		s.Balls[i].Center.Y = s.Height - s.Balls[i].R
	}
}

func (s *Simulation) collisionResponse(ball1, ball2 *Ball) {

	dx := ball1.Center.X - ball2.Center.X
	dy := ball1.Center.Y - ball2.Center.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	v12 := []float64{ball1.Velocity.X - ball2.Velocity.X, ball1.Velocity.Y - ball2.Velocity.Y}
	v21 := []float64{ball2.Velocity.X - ball1.Velocity.X, ball2.Velocity.Y - ball1.Velocity.Y}
	c12 := []float64{ball1.Center.X - ball2.Center.X, ball1.Center.Y - ball2.Center.Y}
	c21 := []float64{ball2.Center.X - ball1.Center.X, ball2.Center.Y - ball1.Center.Y}

	if dist < ball1.R+ball2.R {
		massConst1 := 2 * ball2.R / (ball1.R + ball2.R)
		massConst2 := 2 * ball1.R / (ball1.R + ball2.R)
		const1 := dotProduct(v12, c12) / dotProduct(c12, c12)
		const2 := dotProduct(v21, c21) / dotProduct(c21, c21)

		ball1.Velocity.X = ball1.Velocity.X - massConst1*const1*c12[0]
		ball1.Velocity.Y = ball1.Velocity.Y - massConst1*const1*c12[1]
		ball2.Velocity.X = ball2.Velocity.X - massConst2*const2*c21[0]
		ball2.Velocity.Y = ball2.Velocity.Y - massConst2*const2*c21[1]

		// Ensure there is no clipping
		ball1.Center.X = ball2.Center.X + (ball1.R+ball2.R)*(ball1.Center.X-ball2.Center.X)/dist
		ball1.Center.Y = ball2.Center.Y + (ball1.R+ball2.R)*(ball1.Center.Y-ball2.Center.Y)/dist

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
		return s.Balls[i].Center.X < s.Balls[j].Center.X
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
			if active[j].Center.X+active[j].R < s.Balls[i].Center.X {
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
		s.Balls[i].Center.X += s.Balls[i].Velocity.X * dt
		s.Balls[i].Center.Y += s.Balls[i].Velocity.Y * dt

		// Apply gravity
		if s.Gravity > 0 {
			s.Balls[i].Velocity.Y += s.Gravity
		}

		s.wallCollisionDetection(i)
	}
	s.sweepAndPruneCollisionDetection()
	// s.dumbCollisionDetection()
}
