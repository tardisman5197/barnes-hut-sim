package simulation

import (
	"fmt"
)

// Simulation holds all of the functionality
// to start a barnes hut simulation.
type Simulation struct {
	// grav is the gravitational constant used to calculate
	// the forces between bodies.
	Grav float64 `json:"grav"`
	// theta is used to determine what granularity
	// to which to calculate the forces for each Body.
	// Between 1 and 0, 1 being full granularity.
	Theta float64 `json:"theta"`
	// Bodies stores the list of bodies within the simulation
	Bodies []Body `json:"bodies"`
	// Step the current number of steps that has taken place.
	Step int `json:"step"`
}

// NewSimulation returns an instance of a Simulation
// struct. It initilises some simulation paramaters
// and can optionally set the bodies for the simulation.
func NewSimulation(grav, theta float64, bodies ...Body) Simulation {
	return Simulation{
		Grav:   grav,
		Theta:  theta,
		Bodies: bodies,
	}
}

// oneStep simulates on tick in the a simulation
func (s *Simulation) oneStep(bodies []Body) []Body {
	// Create a new Oct Tree based on the bodies
	root := NewRootNode(bodies)
	root.BuildOcttree(bodies)
	root.CalcMass()

	// Calculate and apply the forces to all
	// the bodes in the simulation for a single
	// tick.
	root.CalcForces(s.Grav, s.Theta)

	// Display the resulting oct tree
	fmt.Println(root.String(0))

	return root.GetBodies()
}

// Steps simulates a number of steps in a simulation
func (s *Simulation) Steps(steps int) []Body {
	for i := 0; i < steps; i++ {
		s.Step++
		fmt.Printf("---------- Step %v ----------\n", s.Step)
		s.Bodies = s.oneStep(s.Bodies)
	}
	return s.Bodies
}
