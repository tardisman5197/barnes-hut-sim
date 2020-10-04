# Barnes Hut Simulator  

A Go implementation of a Barnes–Hut simulation, it allows you to start a server which has an HTTP API which allows you to create and run multiple simulations.

## Server
The server resides at `/cmd/server/main.go`

## Server Api
### Create new Sim
**POST** /simulation/new

This request should return a `simID`

### Start Sim
**GET** /simulation/start/**simID**/**steps**
- `simID`: the ID of the sim you want to start
- `steps`: the number of steps you want the sim to run for

### Sim Status
**GET** /simulation/status/**SimID**
- `simID`: the ID of the sim you want the status for

### Sim Results
**GET** /simulation/results/**SimID**
- `simID`: the ID of the sim you want results for

### Sim Remove
**GET** /simulation/remove/**SimID**
- `simID`: the ID of the sim you want to remove

## Code Examples 
Some examples can be found in `/cmd/examples`

<img align="right" alt="Github Stats" src="https://github-readme-stats.vercel.app/api?username=tardisman5197&show_icons=true&hide_border=true&count_private=true" />
