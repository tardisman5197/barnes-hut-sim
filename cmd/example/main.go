package main

import "github.com/tardisman5197/barnes-hut-sim/pkg/simulation"

const (
	// theta is used to determine what granularity
	// to which to calculate the forces for each Body.
	// Between 1 and 0, 1 being full granularity.
	theta = 0.5
	// grav is the gravitational constant used to calculate
	// the forces between bodies.
	grav = 9.81
)

// main this is run when the program is executed.
func main() {
	bodies := []simulation.Body{
		simulation.Body{Name: "Hi", X: 10.0, Y: 10.0, Z: 10.0, Radius: 1, Density: 1},
		// Body{name: "world", x: 13.0, y: 13.0, z: 13.0, radius: 1, density: 2},
		// Body{name: "people", x: 5.0, y: 5.0, z: 5.0, radius: 3, density: 2},
		simulation.Body{Name: "stuff", X: 1.0, Y: 10.0, Z: 1.0, Radius: 1, Density: 1},
	}

	simulation.XSteps(3, bodies, grav, theta)
}
