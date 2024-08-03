package simulation

func (pi *Pillision) wallCollisionDetection() {
	if pi.BigSquare.TopLeft.X <= 0 {
		pi.BigSquare.Velocity *= -1
		pi.BigSquare.TopLeft.X = 0
	}
 //  if pi.BigSquare.TopLeft.X <= pi.SmallSquare.Width {
	// 	pi.BigSquare.Velocity *= -1
	// 	pi.BigSquare.TopLeft.X = pi.SmallSquare.Width
	// }
	if pi.SmallSquare.TopLeft.X <= 0 {
		pi.SmallSquare.Velocity *= -1
		pi.SmallSquare.TopLeft.X = 0
		pi.SmallSquare.BottomRight.X = pi.SmallSquare.Width
		pi.NumCollisions++
	}
}

func (pi *Pillision) squareCollisionDetection() {
	dx := pi.BigSquare.TopLeft.X - pi.SmallSquare.BottomRight.X
	v12 := pi.BigSquare.Velocity - pi.SmallSquare.Velocity
	v21 := pi.SmallSquare.Velocity - pi.BigSquare.Velocity

	if dx <= 0 {
		pi.NumCollisions++
		massConst1 := 2 * pi.SmallSquare.Weight / (pi.BigSquare.Weight + pi.SmallSquare.Weight)
		massConst2 := 2 * pi.BigSquare.Weight / (pi.BigSquare.Weight + pi.SmallSquare.Weight)

		pi.BigSquare.Velocity = pi.BigSquare.Velocity - massConst1*v12
		pi.SmallSquare.Velocity = pi.SmallSquare.Velocity - massConst2*v21

		// Ensure there is no clipping
		// pi.BigSquare.TopLeft.X = pi.SmallSquare.BottomRight.X
		// pi.BigSquare.BottomRight.X = pi.BigSquare.TopLeft.X + pi.BigSquare.Width
	}
}

func (pi *Pillision) Update() {
	pi.wallCollisionDetection()
	pi.squareCollisionDetection()

	// move squares according to velocity
	pi.BigSquare.TopLeft.X += pi.BigSquare.Velocity
	pi.BigSquare.BottomRight.X += pi.BigSquare.Velocity
	pi.SmallSquare.TopLeft.X += pi.SmallSquare.Velocity
	pi.SmallSquare.BottomRight.X += pi.SmallSquare.Velocity
}
