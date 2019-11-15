package main

import (
	"fmt"
	"math"
	"strings"
)

// octNode represents a cube in 3D space
// organised in a tree.
type octNode struct {
	// Body is struct containing the data about each object
	// in space
	body Body
	// childeren is a slice of nodes that are contained by
	// the node
	children []octNode
	// parent is a pointer to the node of which this node
	// is a child to
	parent *octNode
	// empty is true if it contains no body
	empty bool
	// xyz is the top corner of the cube
	x, y, z float64
	// dx, dy, dz is the length width and depth of the cube
	dx, dy, dz float64
	// The center of mass of the cube
	cmx, cmy, cmz float64
	// fx. fy, fz is the force that should be applied to
	// the cube
	fx, fy, fz float64
	// mass is the total mass of all its self and children
	// nodes
	mass float64
}

// newOctNode child node of the parent at a given position
// the size of the node is half in every direction of the
// parent node. The new node does not contain any children
// and is empty when returned.
func newOctNode(parent *octNode, x, y, z float64) octNode {
	newNode := octNode{}
	newNode.parent = parent
	newNode.empty = true
	newNode.x = x
	newNode.y = y
	newNode.z = z
	newNode.dx = (*parent).dx / 2
	newNode.dy = (*parent).dy / 2
	newNode.dz = (*parent).dz / 2
	newNode.children = make([]octNode, 0)
	return newNode
}

// createRootNode creates an empty node with a size and position
// which encompasses all of the bodies provided.
func createRootNode(bodies []Body) octNode {
	// find the lowest and highest coordinates
	lx, ly, lz := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	hx, hy, hz := -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < len(bodies); i++ {
		// Check if the body is in a higher
		// position
		if bodies[i].x > hx {
			hx = bodies[i].x + 1
		}
		if bodies[i].y > hy {
			hy = bodies[i].y + 1
		}
		if bodies[i].z > hz {
			hz = bodies[i].z + 1
		}

		// Check if the body is in a lower
		// position
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

	// create the root node
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

// buildOctTree creates a new oct tree from a root node
// and a list of bodies.
func buildOcttree(bodies []Body, root *octNode) {
	// Add points to tree
	for _, body := range bodies {
		octInsert(body, root)
	}

	// Remove empty leafs
	removeEmpty(root)

}

// octInsert inserts a body into a oct tree based on the new
// body's position. If the node that should contain the new
// body is already filled that node is split into 8 smaller nodes.
func octInsert(body Body, node *octNode) {
	// Find the correct node to insert the data into
	if len((*node).children) > 1 {
		// Check if the node has children
		// then find the correct leaf to insert the data into
		// Loop through the children of the node and
		// check which leaf the data should be contained
		for i := 0; i < len((*node).children); i++ {
			if inside(body, (*node).children[i]) {
				octInsert(body, &(*node).children[i])
			}
		}
	} else if !(*node).empty && len((*node).children) == 0 {
		// if the node is not empty and does not have any children
		// split the node into 8 empty leafs.
		// Then insert the original node's data into a new leaf.
		// Find the leaf for the new data and insert it.

		// Create the new leafs
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					// Create a new leaf with half the dx, dy and dz
					(*node).children = append(
						(*node).children,
						newOctNode(node,
							(*node).x+(float64(i)*((*node).dx/2.0)),
							(*node).y+(float64(j)*((*node).dy/2.0)),
							(*node).z+(float64(k)*((*node).dz/2.0))))
				}
			}
		}

		// Insert the original node's data into the correct leaf
		for i := 0; i < len((*node).children); i++ {
			if inside((*node).body, (*node).children[i]) {
				octInsert((*node).body, &(*node).children[i])
			}
		}

		// Find the child to insert the new data in
		for i := 0; i < len((*node).children); i++ {
			if inside(body, (*node).children[i]) {
				octInsert(body, &(*node).children[i])
			}
		}

		(*node).body = Body{}

	} else if (*node).empty {
		// if the node is empty then insert the data into
		// the leaf node
		(*node).body = body
		(*node).empty = false
	}
}

// inside is true if the Body is within the given
// node.
func inside(body Body, node octNode) bool {
	if body.x >= node.x &&
		body.x < node.x+node.dx &&
		body.y >= node.y &&
		body.y < node.y+node.dy &&
		body.z >= node.z &&
		body.z < node.z+node.dz {
		return true
	}
	return false
}

// removeEmpty prunes the oct tree of any empty nodes
// with no children.
func removeEmpty(node *octNode) {
	// remove any emptyleaf nodes
	for i := len((*node).children) - 1; i >= 0; i-- {
		if (*node).children[i].empty && len((*node).children) == 0 {
			// if the leaf is empty and has no children remove
			(*node).children = append((*node).children[:i], (*node).children[i+1:]...)
		} else if len((*node).children[i].children) > 0 {
			// if the node has children check the children
			// for empty nodes
			removeEmpty(&(*node).children[i])
		}
	}
}

// printOctTree prints out the nodes in the oct tree
// to stdout.
func printOctTree(node *octNode, indents int) {
	if !(*node).empty {
		fmt.Printf("%vNode: %v mass=%v x=%v+%v y=%v+%v z=%v+%v - Body: mass=%v x=%v y=%v z=%v\n",
			strings.Repeat("\t", indents),
			(*node).body.name,
			(*node).mass,
			(*node).x,
			(*node).dx,
			(*node).y,
			(*node).dy,
			(*node).z,
			(*node).dz,
			(*node).body.mass(),
			(*node).body.x,
			(*node).body.y,
			(*node).body.z,
		)
	} else {
		fmt.Printf("%vNode: %v mass=%v x=%v+%v y=%v+%v z=%v+%v\n",
			strings.Repeat("\t", indents),
			(*node).body.name,
			(*node).mass,
			(*node).x,
			(*node).dx,
			(*node).y,
			(*node).dy,
			(*node).z,
			(*node).dz,
		)
	}
	for i := 0; i < len((*node).children); i++ {
		if len((*node).children[i].children) > 0 {
			printOctTree(&(*node).children[i], indents+1)
		} else {
			if !(*node).children[i].empty {
				fmt.Printf("%vNode: %v mass=%v x=%v+%v y=%v+%v z=%v+%v - Body: mass=%v x=%v y=%v z=%v\n",
					strings.Repeat("\t", indents+1),
					(*node).children[i].body.name,
					(*node).children[i].mass,
					(*node).children[i].x,
					(*node).children[i].dx,
					(*node).children[i].y,
					(*node).children[i].dy,
					(*node).children[i].z,
					(*node).children[i].dz,
					(*node).children[i].body.mass(),
					(*node).children[i].body.x,
					(*node).children[i].body.y,
					(*node).children[i].body.z,
				)
			} else {
				fmt.Printf("%vNode: EMPTY %v+%v, %v+%v, %v+%v\n",
					strings.Repeat("\t", indents+1),
					(*node).children[i].x,
					(*node).children[i].dx,
					(*node).children[i].y,
					(*node).children[i].dy,
					(*node).children[i].z,
					(*node).children[i].dz,
				)
			}
		}
	}

}

// calcMass traverses an oct tree calculating the accumulative
// mass of each of a nodes children and calculates the position
// of the node's center of mass
func calcMass(node *octNode) (totalMass, cmx, cmy, cmz float64) {

	if len((*node).children) > 0 {
		// If the node has children work out the mass of all of the
		// children and calculate the center of mass of its children
		for i := 0; i < len((*node).children); i++ {
			// find the mass and center of each child
			cMass, cx, cy, cz := calcMass(&(*node).children[i])
			totalMass += cMass

			cmx += cx * cMass
			cmy += cy * cMass
			cmz += cz * cMass
		}

		// The mass of the node is the sum of all
		// children's masses
		(*node).mass = totalMass

		// The center of mass for the node is
		// sum of the product of each child's mass
		// and center, divided by the total mass of
		// all children nodes.
		//
		// cm = sum(childMass * childCenter) / sum(childrenMass)
		(*node).cmx = cmx / totalMass
		(*node).cmy = cmx / totalMass
		(*node).cmz = cmx / totalMass

		return (*node).mass, (*node).cmx, (*node).cmy, (*node).cmz

	} else if len((*node).children) == 0 && !(*node).empty {
		// If the node is a leaf node, a node with a body and no children,
		// the mass of the node is calculated by the mass of the body
		// and the center of mass is just the body's position
		(*node).mass = (*node).body.mass()
		(*node).cmx = (*node).body.x
		(*node).cmy = (*node).body.y
		(*node).cmz = (*node).body.z

		return (*node).mass, (*node).cmx, (*node).cmy, (*node).cmz
	}

	// In this case the node is empty and the mass and center
	// of mass is 0
	return (*node).mass, (*node).cmx, (*node).cmy, (*node).cmz

}

// calcForces calculates the force that would be applied
// to each body in the oct tree and then applies that force
// to each body.
func calcForces(node *octNode) {
	leafNodes := getLeafNodes(node)

	for i := 0; i < len(leafNodes); i++ {
		// Calculate the force applied to that Body
		fx, fy, fz := treeForce(leafNodes[i], node)
		// Apply force to the body
		(*leafNodes[i]).body.applyForce(fx, fy, fz)
	}
}

// treeForce calculates the force that should be applied to
// a particle based on a oct tree.
func treeForce(particle, node *octNode) (fx, fy, fz float64) {
	// force = G * m * mcm *
	//             xcm - x       ycm - y         zcm - z
	//           ( ---------- , ---------- , ---------- )
	//                r3            r3            r3

	// Do not calculate the force on its self
	if particle.body == node.body {
		return 0, 0, 0
	}

	// 	r = distance from particle i to
	// 		   center of mass of particles in n
	//    = sqrt(   ( xcm - x )2
	//         + ( ycm - y )2
	//         + ( zcm - z )2 )

	dx := (*node).cmx - (*particle).body.x
	dy := (*node).cmy - (*particle).body.y
	dz := (*node).cmz - (*particle).body.z

	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	size := (*node).dx * (*node).dy * (*node).dz

	// If the node is a leaf containing a body
	if len((*node).children) == 0 && !(*node).empty {
		// Calculate the force on particle
		fx = grav * (*particle).mass * (*node).mass * (((*node).cmx - (*particle).body.x) / (r * r * r))
		fy = grav * (*particle).mass * (*node).mass * (((*node).cmy - (*particle).body.y) / (r * r * r))
		fz = grav * (*particle).mass * (*node).mass * (((*node).cmz - (*particle).body.z) / (r * r * r))

		return fx, fy, fz
	}

	if size/r < theta {
		// Calc the force on particle
		fx = grav * (*particle).mass * (*node).mass * (((*node).cmx - (*particle).body.x) / (r * r * r))
		fy = grav * (*particle).mass * (*node).mass * (((*node).cmy - (*particle).body.y) / (r * r * r))
		fz = grav * (*particle).mass * (*node).mass * (((*node).cmz - (*particle).body.z) / (r * r * r))

		return fx, fy, fz
	}

	for i := 0; i < len((*node).children); i++ {
		// Calculate the resulting force of all
		// of the nodes children's forces
		ifx, ify, ifz := treeForce(particle, &(*node).children[i])
		fx += ifx
		fy += ify
		fz += ifz
	}

	return fx, fy, fz
}

// getLeafNodes returns a list of all the nodes that contain
// a body.
func getLeafNodes(node *octNode) (leafNodes []*octNode) {
	// If the node has children, check each child for
	// leaf nodes
	if len((*node).children) > 0 {
		for i := 0; i < len((*node).children); i++ {
			tmpLeafNodes := getLeafNodes(&(*node).children[i])
			leafNodes = append(leafNodes, tmpLeafNodes...)
		}
	} else if len((*node).children) == 0 && !(*node).empty {
		// Found a leaf node return add the node to be
		// returned
		leafNodes = append(leafNodes, node)
	}

	return leafNodes
}

// getBodies returns all of the Body struts stored
// in a oct tree
func getBodies(node *octNode) (bodies []Body) {
	// Get all the leaf nodes in the oct tree
	leafNodes := getLeafNodes(node)
	// Retreive the body's attached to each
	// leaf node
	for i := 0; i < len(leafNodes); i++ {
		bodies = append(bodies, leafNodes[i].body)
	}
	return bodies
}
