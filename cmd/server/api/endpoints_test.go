package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/tardisman5197/barnes-hut-sim/pkg/simulation"
)

const (
	BaseURL  string = "http://localhost:5000/simulation"
	NewURL   string = BaseURL + "/new"
	StartURL string = BaseURL + "/start"
)

var testsStart = []struct {
	expected    int
	given       string
	description string
}{
	{expected: http.StatusNotFound, given: "/t/10", description: "Non-existent simulation"},
}

var simul = simulation.Simulation{
	Grav:  9.81,
	Theta: 0.1,
	Bodies: []simulation.Body{
		{Name: "Hi", X: 10.0, Y: 10.0, Z: 10.0, Radius: 1, Density: 1},
		{Name: "stuff", X: 1.0, Y: 10.0, Z: 1.0, Radius: 1, Density: 1},
	},
}

func addSimul() {
	body, err := json.Marshal(simul)
	if err != nil {
		panic(err)
	}

	// Create a dummy simulation to test on an existing one
	r, err := http.Post(NewURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var res struct {
		ID         string                `json:"id"`
		Simulation simulation.Simulation `json:"simulation"`
	}

	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		panic(err)
	}

	testsStart = append(testsStart, []struct {
		expected    int
		given       string
		description string
	}{
		{expected: http.StatusOK, given: fmt.Sprintf("/%s/10", res.ID), description: "Normal use"},
		{expected: http.StatusBadRequest, given: fmt.Sprintf("/%s/-1", res.ID), description: "Negative number of steps"},
		{expected: http.StatusBadRequest, given: fmt.Sprintf("/%s/t", res.ID), description: "Non digit steps"},
	}...)
}

func TestMain(m *testing.M) {
	a := NewAPI()
	a.Listen()
	addSimul()

	os.Exit(m.Run())
}

func TestStart(t *testing.T) {
	for _, test := range testsStart {
		r, err := http.Get(StartURL + test.given)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		if r.StatusCode != test.expected {
			t.Fatalf("%s%s: expected %d, got %d", StartURL, test.given, test.expected, r.StatusCode)
		}
	}
}
