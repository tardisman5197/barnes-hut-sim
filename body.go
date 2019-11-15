package main

import "math"

// Body ...
type Body struct {
	name    string
	x, y, z float64
	radius  float64
	density float64
}

func (b *Body) mass() float64 {
	volume := (4 / 3) * math.Pi * math.Pow(b.radius, 3)

	mass := volume * b.density

	return mass
}
