package main

func main() {

	unis := []Uni{
		Uni{name: "Hi", x: 10.0, y: 10.0, z: 10.0, noOfStudents: 2, greggDensity: 2},
		Uni{name: "world", x: 13.0, y: 13.0, z: 13.0, noOfStudents: 1, greggDensity: 2},
		Uni{name: "people", x: 5.0, y: 5.0, z: 5.0, noOfStudents: 3, greggDensity: 2},
		Uni{name: "stuff", x: 9.0, y: 9.0, z: 13.0, noOfStudents: 4, greggDensity: 2},
	}

	root := octNode{children: make([]octNode, 0), empty: true, dx: 50, dy: 50, dz: 50}
	buildOcttree(unis, &root)
	calcMass(&root)
	printOctTree(&root, 0)
}
