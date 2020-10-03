package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/tardisman5197/barnes-hut-sim/pkg/simulation"
)

type NewSimulationRequest struct {
	Grav   float64           `json:"grav"`
	Theta  float64           `json:"theta"`
	Bodies []simulation.Body `json:"bodies,omitempty"`
}

type NewSimulationResponse struct {
	ID         string                `json:"id"`
	Simulation simulation.Simulation `json:"simulation"`
}

// newSimulation is called when a request is made to "/simulation/new".
// It creates a new simulation with a unique ID and then returns the
// details of the simulation to the requester.
func (a *API) newSimulation(w http.ResponseWriter, r *http.Request) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	fmt.Println("New Simulation Request")
	// Read in therequest body

	var req NewSimulationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find a new simulation ID/
	// I think this will timeout when the
	// WriteTimeout limit is reached.
	var id string
	for {
		id = randomID(SimulationIDLength)
		if _, exists := a.simulations[id]; !exists {
			break
		}
	}

	// Create a new simulation
	a.simulations[id] = simulation.NewSimulation(req.Grav, req.Theta, req.Bodies...)

	// Send simulation information back to the requester
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		NewSimulationResponse{
			ID:         id,
			Simulation: a.simulations[id],
		},
	)
}

// start is called when a request is made to "/simulation/start/{simID}/{steps}".
// This will start the simulation with the specified ID for
// a certain number of steps.
func (a *API) start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start Simulation Request")
	w.WriteHeader(http.StatusNotImplemented)
}

// status is called when a request is made to "/simulation/status/{simID}".
// This endpoint will return the status of the simulation with
// the specified simulation ID.
func (a *API) status(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status Simulation Request")
	w.WriteHeader(http.StatusNotImplemented)
}

// results is called when a request is made to "/simulation/results/{simID}".
// This endpoint will return the results of the simulation with
// the ID specified.
func (a *API) results(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Results Simulation Request")
	w.WriteHeader(http.StatusNotImplemented)
}

// remove is called when a request is made to "/simulation/remove/{simID}".
// This endpoint will remove the simulation with the ID requested.
func (a *API) remove(w http.ResponseWriter, r *http.Request) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	vars := mux.Vars(r)

	simID, hasSimID := vars["simID"]
	if !hasSimID {
		http.Error(w, "simulation id not provided", http.StatusBadRequest)
		return
	}

	_, present := a.simulations[simID]
	if !present {
		http.Error(w, fmt.Sprintf("simulation with id %s not present", simID), http.StatusBadRequest)
		return
	}

	delete(a.simulations, simID)
	w.WriteHeader(http.StatusOK)
}
