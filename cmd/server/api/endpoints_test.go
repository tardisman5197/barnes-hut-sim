package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tardisman5197/barnes-hut-sim/pkg/simulation"
)

const (
	StartURL  string = "/start"
	BaseURL   string = "http://localhost:5000/simulation"
	NewURL    string = BaseURL + "/new"
	StatusURL string =  "/status"
)

var testsStart = []struct {
	expected    int
	given       string
	description string
}{
	{expected: http.StatusNotFound, given: "/t/10", description: "Non-existent simulation"},
}

var testsStatus = []struct {
	expected    int
	given       string
	description string
	body        string
}{
	{expected: http.StatusNotFound, given: "/1",
		description: "Non-existent simulation", body: "there is no simulation with the simID 1\n"},
}

var simul = simulation.Simulation{
	Grav:  9.81,
	Theta: 0.1,
	Bodies: []simulation.Body{
		{Name: "Hi", X: 10.0, Y: 10.0, Z: 10.0, Radius: 1, Density: 1},
		{Name: "stuff", X: 1.0, Y: 10.0, Z: 1.0, Radius: 1, Density: 1},
	},
}

func addSimul(baseUrl string) {
	res := createSimulation(baseUrl)

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


func statusSimul(baseURL string){
	res := createSimulation(baseURL)

	testsStatus = append(testsStatus, []struct {
		expected    int
		given       string
		description string
		body        string
	}{
		{expected: http.StatusOK, given: fmt.Sprintf("/%s", res.ID),
			description: "Existent simulation", body: fmt.Sprintf("{\"id\":\"%s\",\"step\":0}\n",
			res.ID)},
	}...)
}

func createSimulation(baseUrl string) struct {
	ID         string                `json:"id"`
	Simulation simulation.Simulation `json:"simulation"`
} {
	body, err := json.Marshal(simul)
	if err != nil {
		panic(err)
	}

	// Create a dummy simulation to test on an existing one
	r, err := http.Post(baseUrl+"/new", "application/json", bytes.NewBuffer(body))
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
	return res
}

func TestStart(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.router())
	defer srv.Close()

	simulationEndpoint := srv.URL + "/simulation"
	addSimul(simulationEndpoint)

	for _, test := range testsStart {
		t.Logf("Test case: %s", test.description)
		r, err := http.Get(simulationEndpoint + StartURL + test.given)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != test.expected {
			t.Fatalf("%s%s: expected %d, got %d", StartURL, test.given, test.expected, r.StatusCode)
		}

		if err := r.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestStatus(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.router())
	defer srv.Close()

	simulationEndpoint := srv.URL + "/simulation"
	statusSimul(simulationEndpoint)

	t.Logf("testsStatus : %v", testsStatus)


	for _, test := range testsStatus {
		t.Logf("Test case: %s", test.description)
		t.Logf("Get : %s", simulationEndpoint + StatusURL + test.given)
		r, err := http.Get(simulationEndpoint + StatusURL + test.given)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		if r.StatusCode != test.expected {
			t.Fatalf("%s%s: expected %d, got %d", StatusURL, test.given, test.expected, r.StatusCode)
		}
		if len(test.description) > 0 {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			bodyString := string(bodyBytes)

			if bodyString != test.body {
				t.Fatalf("%s%s: expected %s, got %s", StatusURL, test.given, test.body,
					bodyString)
			}
		}
	}
}

func TestStatusApi(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.router())
	defer srv.Close()

	inputSimulation := &simulation.Simulation{
		Grav:  9.81,
		Theta: 31,
		Bodies: []simulation.Body{
			{
				Name:    "test",
				X:       1,
				Y:       2,
				Z:       3,
				Radius:  4,
				Density: 5,
			},
		},
	}
	api.simulations["test_id"] = inputSimulation

	resp, err := http.Get(srv.URL + "/simulation/results/test_id")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", resp.StatusCode, http.StatusOK)
	}

	var resultResponse simulationResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&resultResponse); err != nil {
		t.Fatal(err)
	}

	if resultResponse.Simulation.Grav != inputSimulation.Grav {
		t.Fatalf("unexpected simulation Grav value %f != %f", resultResponse.Simulation.Grav, inputSimulation.Grav)
	}

	if resultResponse.Simulation.Theta != inputSimulation.Theta {
		t.Fatalf("unexpected simulation Theta value %f != %f", resultResponse.Simulation.Theta, inputSimulation.Theta)
	}

	if len(resultResponse.Simulation.Bodies) != len(inputSimulation.Bodies) {
		t.Fatalf("unexpected simulation Bodies len")
	}

	if resultResponse.Simulation.Bodies[0] != inputSimulation.Bodies[0] {
		t.Fatalf("unexpected simulation Bodies values")
	}

}

func TestStatusApiWithInvalidSimID(t *testing.T) {
	api := NewAPI()
	srv := httptest.NewServer(api.router())
	defer srv.Close()

	api.simulations["test_id"] = &simulation.Simulation{}

	resp, err := http.Get(srv.URL + "/simulation/results/invalid_test_id")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d != %d", resp.StatusCode, http.StatusBadRequest)
	}
}
