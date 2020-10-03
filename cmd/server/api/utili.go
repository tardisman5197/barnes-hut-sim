package api

import "math/rand"

// randomID returns a string of the length specified
// with random capital chracters.
func randomID(length int) string {
	id := make([]rune, length)

	for i := 0; i < length; i++ {
		id[i] = rune(rand.Intn(25)) + 'A'
	}

	return string(id)
}
