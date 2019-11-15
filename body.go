package main

import "math"

// Body contains information about an object
// within the simulation
type Body struct {
	// name stores the identifer of the body
	name string
	// x, y, z stores the position of the body
	x, y, z float64
	// radius stores the radius of the body's sphere
	radius float64
	// density stores the density of the material
	// of the body
	density float64
}

// mass calculates the mass of the Body based
// on the radius and density.
func (b *Body) mass() float64 {
	volume := (4 / 3) * math.Pi * math.Pow(b.radius, 3)

	mass := volume * b.density

	return mass
}

// applyForce modifies the position of the Body
// based on a force provided.
func (b *Body) applyForce(fx, fy, fz float64) {
	b.x += fx
	b.y += fy
	b.z += fz
}
