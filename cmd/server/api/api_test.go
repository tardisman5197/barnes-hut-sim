package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tardisman5197/barnes-hut-sim/pkg/simulation"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEndpointCreateSimulation(t *testing.T) {
	api := NewAPI()

	reqBody := NewSimulationRequest{
		Grav:   1,
		Theta:  0,
		Bodies: []simulation.Body{},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	request := &http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewBuffer(body)),
	}

	rr := httptest.NewRecorder()

	api.newSimulation(rr, request)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusOK)
	}

	var response NewSimulationResponse
	if err := json.NewDecoder(rr.Result().Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.ID == "" {
		t.Fatal("empty simulation ID")
	}

	if response.Simulation.Grav != reqBody.Grav {
		t.Fatalf("expected response Grav to be %f but was %f", response.Simulation.Grav, reqBody.Grav)
	}

	if response.Simulation.Theta != reqBody.Theta {
		t.Fatalf("expected response Theta to be %f but was %f", response.Simulation.Theta, reqBody.Theta)
	}

	currentSimulations := len(api.simulations)
	if currentSimulations != 1 {
		t.Fatalf("expected 1 simulation, found %d", currentSimulations)
	}

}

func TestEndpointRemoveSimulation(t *testing.T) {
	api := NewAPI()

	reqBody := NewSimulationRequest{
		Grav:   1,
		Theta:  0,
		Bodies: []simulation.Body{},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	request := &http.Request{
		Body: ioutil.NopCloser(bytes.NewBuffer(body)),
	}

	rr := httptest.NewRecorder()

	api.newSimulation(rr, request)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusOK)
	}

	var response NewSimulationResponse
	if err := json.NewDecoder(rr.Result().Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.ID == "" {
		t.Fatal("empty simulation ID")
	}

	currentSimulations := len(api.simulations)
	if currentSimulations != 1 {
		t.Fatalf("expected 1 simulation, found %d", currentSimulations)
	}

	rr = httptest.NewRecorder()

	request = mux.SetURLVars(&http.Request{
		Method: http.MethodGet,
	}, map[string]string{
		"simID": response.ID,
	})

	api.remove(rr, request)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusOK)
	}

	currentSimulations = len(api.simulations)
	if currentSimulations != 0 {
		t.Fatalf("expected 0 simulation, found %d", currentSimulations)
	}

}

func TestRemoveNotPresentSimulation(t *testing.T) {
	api := NewAPI()

	reqBody := NewSimulationRequest{
		Grav:   1,
		Theta:  0,
		Bodies: []simulation.Body{},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	request := &http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewBuffer(body)),
	}

	rr := httptest.NewRecorder()

	api.newSimulation(rr, request)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusOK)
	}

	rr = httptest.NewRecorder()

	request = mux.SetURLVars(&http.Request{
		Method: http.MethodGet,
	}, map[string]string{
		"simID": "test",
	})

	api.remove(rr, request)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestRemoveWithoutSimID(t *testing.T) {
	api := NewAPI()
	rr := httptest.NewRecorder()

	request := mux.SetURLVars(&http.Request{
		Method: http.MethodGet,
	}, map[string]string{})

	api.remove(rr, request)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestSimulationStatusEndpoint(t *testing.T) {
	api := NewAPI()

	reqBody := NewSimulationRequest{
		Grav:   1,
		Theta:  0,
		Bodies: []simulation.Body{},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	request := &http.Request{
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewBuffer(body)),
	}

	rr := httptest.NewRecorder()

	api.newSimulation(rr, request)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusOK)
	}

	var response NewSimulationResponse
	if err := json.NewDecoder(rr.Result().Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.ID == "" {
		t.Fatal("empty simulation ID")
	}

	rr = httptest.NewRecorder()

	request = mux.SetURLVars(&http.Request{
		Method: http.MethodGet,
	}, map[string]string{
		"simID": response.ID,
	})

	api.results(rr, request)

	if rr.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusOK)
	}

	var resultResponse simulationResultResponse
	if err := json.NewDecoder(rr.Result().Body).Decode(&resultResponse); err != nil {
		t.Fatal(err)
	}

	if resultResponse.Simulation.Grav != reqBody.Grav {
		t.Fatalf("unexpected simulation Grav value %f != %f", resultResponse.Simulation.Grav, reqBody.Grav)
	}

	if resultResponse.Simulation.Theta != reqBody.Theta {
		t.Fatalf("unexpected simulation Theta value %f != %f", resultResponse.Simulation.Theta, reqBody.Theta)
	}

	if reflect.DeepEqual(resultResponse.Simulation.Bodies, reqBody.Bodies) {
		t.Fatalf("unexpected simulation Bodies values")
	}

}

func TestSimulationStatusEndpointWithoutSimID(t *testing.T) {
	api := NewAPI()
	rr := httptest.NewRecorder()

	request := mux.SetURLVars(&http.Request{
		Method: http.MethodGet,
	}, map[string]string{})

	api.results(rr, request)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d != %d", rr.Result().StatusCode, http.StatusBadRequest)
	}
}
