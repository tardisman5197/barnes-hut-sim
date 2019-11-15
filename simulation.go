package main

import (
	"fmt"
	"math"
)

func oneStep(bodies []Body) []Body {
	root := createRootNode(bodies)
	buildOcttree(bodies, &root)
	calcMass(&root)
	calcForces(&root)
	printOctTree(&root, 0)
	return getBodies(&root)
}

func xSteps(steps int, bodies []Body) []Body {
	for i := 0; i < steps; i++ {
		fmt.Printf("---------- Step %v ----------\n", i)
		bodies = oneStep(bodies)
	}
	return bodies
}

func createRootNode(bodies []Body) octNode {
	lx, ly, lz := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	hx, hy, hz := -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < len(bodies); i++ {
		if bodies[i].x > hx {
			hx = bodies[i].x + 1
		}
		if bodies[i].y > hy {
			hy = bodies[i].y + 1
		}
		if bodies[i].z > hz {
			hz = bodies[i].z + 1
		}

		if bodies[i].x < lx {
			lx = bodies[i].x - 1
		}
		if bodies[i].y < ly {
			ly = bodies[i].y - 1
		}
		if bodies[i].z < lz {
			lz = bodies[i].z - 1
		}
	}
	root := octNode{
		children: make([]octNode, 0),
		empty:    true,
		x:        lx,
		y:        ly,
		z:        lz,
		dx:       hx - lx,
		dy:       hy - ly,
		dz:       hz - lz,
	}
	return root
}
