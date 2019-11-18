package main

// main this is run when the program is executed.
func main() {

	bodies := []Body{
		Body{name: "Hi", x: 10.0, y: 10.0, z: 10.0, radius: 1, density: 1},
		// Body{name: "world", x: 13.0, y: 13.0, z: 13.0, radius: 1, density: 2},
		// Body{name: "people", x: 5.0, y: 5.0, z: 5.0, radius: 3, density: 2},
		Body{name: "stuff", x: 1.0, y: 10.0, z: 1.0, radius: 1, density: 1},
	}

	xSteps(3, bodies)
}