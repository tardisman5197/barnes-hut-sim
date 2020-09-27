package simulation

import (
	"fmt"
	"math"
	"strings"
)

// OctNode represents a cube in 3D space
// organised in a tree.
type OctNode struct {
	// Body is struct containing the data about each object
	// in space
	body Body
	// childeren is a slice of nodes that are contained by
	// the node
	children []OctNode
	// parent is a pointer to the node of which this node
	// is a child to
	parent *OctNode
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

// NewOctNode child node of the parent at a given position
// the size of the node is half in every direction of the
// parent node. The new node does not contain any children
// and is empty when returned.
func NewOctNode(parent *OctNode, x, y, z float64) OctNode {
	return OctNode{
		parent:   parent,
		empty:    true,
		x:        x,
		y:        y,
		z:        z,
		dx:       (*parent).dx / 2,
		dy:       (*parent).dy / 2,
		dz:       (*parent).dz / 2,
		children: make([]OctNode, 0),
	}
}

// NewRootNode creates an empty node with a size and position
// which encompasses all of the bodies provided.
func NewRootNode(bodies []Body) OctNode {
	// find the lowest and highest coordinates
	lx, ly, lz := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64
	hx, hy, hz := -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for i := 0; i < len(bodies); i++ {
		// Check if the body is in a higher
		// position
		if bodies[i].X > hx {
			hx = bodies[i].X + 1
		}
		if bodies[i].Y > hy {
			hy = bodies[i].Y + 1
		}
		if bodies[i].Z > hz {
			hz = bodies[i].Z + 1
		}

		// Check if the body is in a lower
		// position
		if bodies[i].X < lx {
			lx = bodies[i].X - 1
		}
		if bodies[i].Y < ly {
			ly = bodies[i].Y - 1
		}
		if bodies[i].Z < lz {
			lz = bodies[i].Z - 1
		}
	}

	// create the root node
	return OctNode{
		children: make([]OctNode, 0),
		empty:    true,
		x:        lx,
		y:        ly,
		z:        lz,
		dx:       hx - lx,
		dy:       hy - ly,
		dz:       hz - lz,
	}
}

// BuildOcttree creates a new oct tree from a root node
// and a list of bodies.
func (n *OctNode) BuildOcttree(bodies []Body) {
	// Add points to tree
	for _, body := range bodies {
		n.OctInsert(body)
	}

	// Remove empty leafs
	n.removeEmpty()
}

// OctInsert inserts a body into a oct tree based on the new
// body's position. If the node that should contain the new
// body is already filled that node is split into 8 smaller nodes.
func (n *OctNode) OctInsert(body Body) {
	// Find the correct node to insert the data into
	if len(n.children) > 1 {
		// Check if the node has children
		// then find the correct leaf to insert the data into
		// Loop through the children of the node and
		// check which leaf the data should be contained
		for i := 0; i < len(n.children); i++ {
			if n.children[i].inside(body) {
				n.children[i].OctInsert(body)
			}
		}
	} else if !n.empty && len(n.children) == 0 {
		// if the node is not empty and does not have any children
		// split the node into 8 empty leafs.
		// Then insert the original node's data into a new leaf.
		// Find the leaf for the new data and insert it.

		// Create the new leafs
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				for k := 0; k < 2; k++ {
					// Create a new leaf with half the dx, dy and dz
					n.children = append(
						n.children,
						NewOctNode(n,
							n.x+(float64(i)*(n.dx/2.0)),
							n.y+(float64(j)*(n.dy/2.0)),
							n.z+(float64(k)*(n.dz/2.0)),
						),
					)
				}
			}
		}

		// Insert the original node's data into the correct leaf
		for i := 0; i < len(n.children); i++ {
			if n.children[i].inside(n.body) {
				n.children[i].OctInsert(n.body)
			}
		}

		// Find the child to insert the new data in
		for i := 0; i < len(n.children); i++ {
			if n.children[i].inside(body) {
				n.children[i].OctInsert(body)
			}
		}

		n.body = Body{}

	} else if n.empty {
		// if the node is empty then insert the data into
		// the leaf node
		n.body = body
		n.empty = false
	}
}

// inside is true if the Body is within the given
// node.
func (n *OctNode) inside(body Body) bool {
	if body.X >= n.x &&
		body.X < n.x+n.dx &&
		body.Y >= n.y &&
		body.Y < n.y+n.dy &&
		body.Z >= n.z &&
		body.Z < n.z+n.dz {
		return true
	}
	return false
}

// removeEmpty prunes the oct tree of any empty nodes
// with no children.
func (n *OctNode) removeEmpty() {
	// remove any emptyleaf nodes
	for i := len(n.children) - 1; i >= 0; i-- {
		if n.children[i].empty && len(n.children) == 0 {
			// if the leaf is empty and has no children remove
			n.children = append(n.children[:i], n.children[i+1:]...)
		} else if len(n.children[i].children) > 0 {
			// if the node has children check the children
			// for empty nodes
			n.children[i].removeEmpty()
		}
	}
}

// String returns a string respresentation of a OctNode
func (n *OctNode) String(indents int) string {
	var result string
	if !n.empty {
		result += fmt.Sprintf(
			"%vNode: %v mass=%v x=%v+%v y=%v+%v z=%v+%v - Body: mass=%v x=%v y=%v z=%v\n",
			strings.Repeat("\t", indents),
			n.body.Name,
			n.mass,
			n.x,
			n.dx,
			n.y,
			n.dy,
			n.z,
			n.dz,
			n.body.mass(),
			n.body.X,
			n.body.Y,
			n.body.Z,
		)
	} else {
		result += fmt.Sprintf(
			"%vNode: %v mass=%v x=%v+%v y=%v+%v z=%v+%v\n",
			strings.Repeat("\t", indents),
			n.body.Name,
			n.mass,
			n.x,
			n.dx,
			n.y,
			n.dy,
			n.z,
			n.dz,
		)
	}
	for i := 0; i < len(n.children); i++ {
		if len(n.children[i].children) > 0 {
			result += n.children[i].String(indents + 1)
		} else {
			if !n.children[i].empty {
				result += fmt.Sprintf(
					"%vNode: %v mass=%v x=%v+%v y=%v+%v z=%v+%v - Body: mass=%v x=%v y=%v z=%v\n",
					strings.Repeat("\t", indents+1),
					n.children[i].body.Name,
					n.children[i].mass,
					n.children[i].x,
					n.children[i].dx,
					n.children[i].y,
					n.children[i].dy,
					n.children[i].z,
					n.children[i].dz,
					n.children[i].body.mass(),
					n.children[i].body.X,
					n.children[i].body.Y,
					n.children[i].body.Z,
				)
			} else {
				result += fmt.Sprintf(
					"%vNode: EMPTY %v+%v, %v+%v, %v+%v\n",
					strings.Repeat("\t", indents+1),
					n.children[i].x,
					n.children[i].dx,
					n.children[i].y,
					n.children[i].dy,
					n.children[i].z,
					n.children[i].dz,
				)
			}
		}
	}
	return result
}

// CalcMass traverses an oct tree calculating the accumulative
// mass of each of a nodes children and calculates the position
// of the node's center of mass
func (n *OctNode) CalcMass() (totalMass, cmx, cmy, cmz float64) {
	if len(n.children) > 0 {
		// If the node has children work out the mass of all of the
		// children and calculate the center of mass of its children
		for i := 0; i < len(n.children); i++ {
			// find the mass and center of each child
			cMass, cx, cy, cz := n.children[i].CalcMass()
			totalMass += cMass

			cmx += cx * cMass
			cmy += cy * cMass
			cmz += cz * cMass
		}

		// The mass of the node is the sum of all
		// children's masses
		n.mass = totalMass

		// The center of mass for the node is
		// sum of the product of each child's mass
		// and center, divided by the total mass of
		// all children nodes.
		//
		// cm = sum(childMass * childCenter) / sum(childrenMass)
		n.cmx = cmx / totalMass
		n.cmy = cmx / totalMass
		n.cmz = cmx / totalMass

		return n.mass, n.cmx, n.cmy, n.cmz

	} else if len(n.children) == 0 && !n.empty {
		// If the node is a leaf node, a node with a body and no children,
		// the mass of the node is calculated by the mass of the body
		// and the center of mass is just the body's position
		n.mass = n.body.mass()
		n.cmx = n.body.X
		n.cmy = n.body.Y
		n.cmz = n.body.Z

		return n.mass, n.cmx, n.cmy, n.cmz
	}

	// In this case the node is empty and the mass and center
	// of mass is 0
	return n.mass, n.cmx, n.cmy, n.cmz
}

// CalcForces calculates the force that would be applied
// to each body in the oct tree and then applies that force
// to each body.
func (n *OctNode) CalcForces(grav, theta float64) {
	leafNodes := n.GetLeafNodes()

	for i := 0; i < len(leafNodes); i++ {
		// Calculate the force applied to that Body
		fx, fy, fz := treeForce(leafNodes[i], n, grav, theta)
		// Apply force to the body
		(*leafNodes[i]).body.applyForce(fx, fy, fz)
	}
}

// treeForce calculates the force that should be applied to
// a particle based on a oct tree.
func treeForce(particle, node *OctNode, grav, theta float64) (fx, fy, fz float64) {
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

	dx := (*node).cmx - (*particle).body.X
	dy := (*node).cmy - (*particle).body.Y
	dz := (*node).cmz - (*particle).body.Z

	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	size := (*node).dx * (*node).dy * (*node).dz

	// If the node is a leaf containing a body
	if len((*node).children) == 0 && !(*node).empty {
		// Calculate the force on particle
		fx = grav * (*particle).mass * (*node).mass * (((*node).cmx - (*particle).body.X) / (r * r * r))
		fy = grav * (*particle).mass * (*node).mass * (((*node).cmy - (*particle).body.Y) / (r * r * r))
		fz = grav * (*particle).mass * (*node).mass * (((*node).cmz - (*particle).body.Z) / (r * r * r))

		return fx, fy, fz
	}

	if size/r < theta {
		// Calc the force on particle
		fx = grav * (*particle).mass * (*node).mass * (((*node).cmx - (*particle).body.Z) / (r * r * r))
		fy = grav * (*particle).mass * (*node).mass * (((*node).cmy - (*particle).body.Y) / (r * r * r))
		fz = grav * (*particle).mass * (*node).mass * (((*node).cmz - (*particle).body.Z) / (r * r * r))

		return fx, fy, fz
	}

	for i := 0; i < len((*node).children); i++ {
		// Calculate the resulting force of all
		// of the nodes children's forces
		ifx, ify, ifz := treeForce(particle, &(*node).children[i], grav, theta)
		fx += ifx
		fy += ify
		fz += ifz
	}

	return fx, fy, fz
}

// GetLeafNodes returns a list of all the nodes that contain
// a body.
func (n *OctNode) GetLeafNodes() (leafNodes []*OctNode) {
	// If the node has children, check each child for
	// leaf nodes
	if len(n.children) > 0 {
		for i := 0; i < len(n.children); i++ {
			tmpLeafNodes := n.children[i].GetLeafNodes()
			leafNodes = append(leafNodes, tmpLeafNodes...)
		}
	} else if len(n.children) == 0 && !n.empty {
		// Found a leaf node return add the node to be
		// returned
		leafNodes = append(leafNodes, n)
	}

	return leafNodes
}

// GetBodies returns all of the Body struts stored
// in a oct tree
func (n *OctNode) GetBodies() (bodies []Body) {
	// Get all the leaf nodes in the oct tree
	leafNodes := n.GetLeafNodes()
	// Retreive the body's attached to each
	// leaf node
	for i := 0; i < len(leafNodes); i++ {
		bodies = append(bodies, leafNodes[i].body)
	}
	return bodies
}
