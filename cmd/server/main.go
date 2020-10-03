package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tardisman5197/barnes-hut-sim/cmd/server/api"
)

func main() {
	fmt.Println("Creating API")
	// Create an API instance
	a := api.NewAPI()
	apiDone := a.Listen()
	fmt.Println("Listening")

	signalDone := make(chan bool, 1)
	go handleSignal(signalDone)

	// Wait for the api to finish
	// or a signal to be received
mainBreak:
	for {
		select {
		case <-apiDone:
			fmt.Println("API Done")
			break mainBreak
		case <-signalDone:
			fmt.Println("Signal Done")
			break mainBreak
		}
	}

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	a.Shutdown(ctx)
	fmt.Println("Shutdown")
}

// handleSignal captures a signal to kill the server.
// When a signal is received a value is placed in the
// done channel.
func handleSignal(done chan bool) {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		signal := <-signals
		fmt.Println(signal)
		done <- true
	}()
}
