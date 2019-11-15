package main

import (
	"fmt"
	"math"
	"strings"
)

// octNode a node in the
type octNode struct {
	body     Body
	children []octNode
	parent   *octNode
	empty    bool
	// xyz is the top corner of the cube
	x, y, z float64
	// dx,dy,dz is the length width and depth of the cube
	dx, dy, dz float64
	// The center of mass of the cube
	cmx, cmy, cmz float64
	fx, fy, fz    float64
	mass          float64
}

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

func buildOcttree(bodys []Body, root *octNode) {
	// Add points to tree
	for _, body := range bodys {
		octInsert(body, root)
	}

	// Remove empty leafs
	removeEmpty(root)

}

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
		// Then insert the orginal node's data into a new leaf.
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

func calcMass(node *octNode) (totalMass, cmx, cmy, cmz float64) {
	// ... Compute the mass and center of mass (cm) of
	//    ... all the particles in the subtree rooted at n
	//    if n contains 1 particle
	//         ... the mass and cm of n are identical to
	//         ... the particle's mass and position
	//         store ( mass, cm ) at n
	//         return ( mass, cm )
	//    else
	//         for all four children c(i) of n (i=1,2,3,4)
	//             ( mass(i), cm(i) ) = Compute_Mass(c(i))
	//         end for
	//         mass = mass(1) + mass(2) + mass(3) + mass(4)
	//              ... the mass of a node is the sum of
	//              ... the masses of the children
	//         cm = (  mass(1)*cm(1) + mass(2)*cm(2)
	//               + mass(3)*cm(3) + mass(4)*cm(4)) / mass
	//              ... the cm of a node is a weighted sum of
	//              ... the cm's of the children
	//         store ( mass, cm ) a
	if len((*node).children) > 0 {
		for i := 0; i < len((*node).children); i++ {
			cMass, cx, cy, cz := calcMass(&(*node).children[i])
			totalMass += cMass
			// Calculate center of mass
			cmx += cx * cMass
			cmy += cy * cMass
			cmz += cz * cMass
		}

		(*node).mass = totalMass
		(*node).cmx = cmx / totalMass
		(*node).cmy = cmx / totalMass
		(*node).cmz = cmx / totalMass

		return (*node).mass, (*node).cmx, (*node).cmy, (*node).cmz

	} else if len((*node).children) == 0 && !(*node).empty {
		(*node).mass = (*node).body.mass()
		(*node).cmx = (*node).body.x
		(*node).cmy = (*node).body.y
		(*node).cmz = (*node).body.z

		return (*node).mass, (*node).cmx, (*node).cmy, (*node).cmz
	}

	return (*node).mass, (*node).cmx, (*node).cmy, (*node).cmz

}

func calcForces(node *octNode) {
	leafNodes := getLeafNodes(node)

	for i := 0; i < len(leafNodes); i++ {
		// Get force
		fx, fy, fz := treeForce(leafNodes[i], node)
		// apply force
		(*leafNodes[i]).body = applyForce((*leafNodes[i]).body, fx, fy, fz)
	}
}

func treeForce(particle, node *octNode) (fx, fy, fz float64) {
	// ... Compute gravitational force on particle i
	// ... due to all particles in the box at n
	// f = 0
	// if n contains one particle
	// 	f = force computed using formula (*) above
	// else
	// 	r = distance from particle i to
	// 		   center of mass of particles in n
	// 	D = size of box n
	// 	if D/r < theta
	// 		compute f using formula (*) above
	// 	else
	// 		for all children c of n
	// 			f = f + TreeForce(i,c)
	// 		end for
	// 	end if
	// end if

	// force = G * m * mcm *
	//             xcm - x       ycm - y         zcm - z
	//           ( ---------- , ---------- , ---------- )
	//                r3            r3            r3

	// r = sqrt(   ( xcm - x )2
	//         + ( ycm - y )2
	//         + ( zcm - z )2 )

	// Do not calculate the force on its self
	if particle.body == node.body {
		return 0, 0, 0
	}

	dx := (*node).cmx - (*particle).body.x
	dy := (*node).cmy - (*particle).body.y
	dz := (*node).cmz - (*particle).body.z

	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	size := (*node).dx * (*node).dy * (*node).dz

	// if the node is a leaf
	if len((*node).children) == 0 && !(*node).empty {
		// Calc force on particle

		fx = grav * (*particle).mass * (*node).mass * (((*node).cmx - (*particle).body.x) / (r * r * r))
		fy = grav * (*particle).mass * (*node).mass * (((*node).cmy - (*particle).body.y) / (r * r * r))
		fz = grav * (*particle).mass * (*node).mass * (((*node).cmz - (*particle).body.z) / (r * r * r))

		return fx, fy, fz
	}

	if size/r < theta {
		// Calc force

		fx = grav * (*particle).mass * (*node).mass * (((*node).cmx - (*particle).body.x) / (r * r * r))
		fy = grav * (*particle).mass * (*node).mass * (((*node).cmy - (*particle).body.y) / (r * r * r))
		fz = grav * (*particle).mass * (*node).mass * (((*node).cmz - (*particle).body.z) / (r * r * r))

		return fx, fy, fz
	}

	for i := 0; i < len((*node).children); i++ {
		ifx, ify, ifz := treeForce(particle, &(*node).children[i])
		fx += ifx
		fy += ify
		fz += ifz
	}

	return fx, fy, fz
}

func applyForce(body Body, fx, fy, fz float64) Body {
	body.x = body.x + fx
	body.y = body.y + fy
	body.z = body.z + fz
	return body
}

func getLeafNodes(node *octNode) (leafNodes []*octNode) {
	if len((*node).children) > 0 {
		for i := 0; i < len((*node).children); i++ {
			tmpLeafNodes := getLeafNodes(&(*node).children[i])
			leafNodes = append(leafNodes, tmpLeafNodes...)
		}
	} else if len((*node).children) == 0 && !(*node).empty {
		leafNodes = append(leafNodes, node)
	}

	return leafNodes
}

func getBodies(node *octNode) (bodies []Body) {
	leafNodes := getLeafNodes(node)
	for i := 0; i < len(leafNodes); i++ {
		bodies = append(bodies, leafNodes[i].body)
	}
	return bodies
}
