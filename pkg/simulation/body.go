package simulation

import "math"

// Body contains information about an object
// within the simulation
type Body struct {
	// Name stores the identifer of the body
	Name string `json:"name"`
	// X, Y, Z stores the position of the body
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	// Radius stores the radius of the body's sphere
	Radius float64 `json:"radius"`
	// Density stores the density of the material
	// of the body
	Density float64 `json:"density"`
}

// mass calculates the mass of the Body based
// on the radius and density.
func (b *Body) mass() float64 {
	volume := (4 / 3) * math.Pi * math.Pow(b.Radius, 3)

	mass := volume * b.Density

	return mass
}

// applyForce modifies the position of the Body
// based on a force provided.
func (b *Body) applyForce(fx, fy, fz float64) {
	b.X += fx
	b.Y += fy
	b.Z += fz
}
