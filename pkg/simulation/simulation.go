package simulation

import (
	"fmt"
)

// Simulation holds all of the functionality
// to start a barnes hut simulation.
type Simulation struct {
	// grav is the gravitational constant used to calculate
	// the forces between bodies.
	grav float64
	// theta is used to determine what granularity
	// to which to calculate the forces for each Body.
	// Between 1 and 0, 1 being full granularity.
	theta float64
	// Bodies stores the list of bodies within the simulation
	Bodies []Body
}

// NewSimulation returns an instance of a Simulation
// struct. It initilises some simulation paramaters
// and can optionally set the bodies for the simulation.
func NewSimulation(grav, theta float64, bodies ...Body) Simulation {
	return Simulation{
		grav:   grav,
		theta:  theta,
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
	root.CalcForces(s.grav, s.theta)

	// Display the resulting oct tree
	fmt.Println(root.String(0))

	return root.GetBodies()
}

// XSteps simulates a number of steps in a simulation
func (s *Simulation) XSteps(steps int) []Body {
	for i := 0; i < steps; i++ {
		fmt.Printf("---------- Step %v ----------\n", i)
		s.Bodies = s.oneStep(s.Bodies)
	}
	return s.Bodies
}
