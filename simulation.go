package main

import (
	"fmt"
)

// oneStep simulates on tick in the a simulation
func oneStep(bodies []Body) []Body {
	// Create a new Oct Tree based on the bodies
	root := createRootNode(bodies)
	buildOcttree(bodies, &root)
	calcMass(&root)

	// Calculate and apply the forces to all
	// the bodes in the simulation for a single
	// tick.
	calcForces(&root)

	// Display the resulting oct tree
	printOctTree(&root, 0)

	return getBodies(&root)
}

// xSteps simulates a number of steps in a simulation
func xSteps(steps int, bodies []Body) []Body {
	for i := 0; i < steps; i++ {
		fmt.Printf("---------- Step %v ----------\n", i)
		bodies = oneStep(bodies)
	}
	return bodies
}
