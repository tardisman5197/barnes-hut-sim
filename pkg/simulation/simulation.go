package simulation

import (
	"fmt"
)

// oneStep simulates on tick in the a simulation
func oneStep(bodies []Body, grav, theta float64) []Body {
	// Create a new Oct Tree based on the bodies
	root := NewRootNode(bodies)
	root.BuildOcttree(bodies)
	root.CalcMass()

	// Calculate and apply the forces to all
	// the bodes in the simulation for a single
	// tick.
	root.CalcForces(grav, theta)

	// Display the resulting oct tree
	fmt.Println(root.String(0))

	return root.GetBodies()
}

// XSteps simulates a number of steps in a simulation
func XSteps(steps int, bodies []Body, grav, theta float64) []Body {
	for i := 0; i < steps; i++ {
		fmt.Printf("---------- Step %v ----------\n", i)
		bodies = oneStep(bodies, grav, theta)
	}
	return bodies
}
