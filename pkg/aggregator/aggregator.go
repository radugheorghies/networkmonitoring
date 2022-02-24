package aggregator

import (
	"os"
	"os/signal"
)

// New is the factory function for this package
func New() *Aggregator {
	return &Aggregator{
		nodesData: nodesData{
			nodes: make(map[string]Node),
		},
	}
}

// Run will start the **magic**
func (a *Aggregator) Run() {
	a.init(true)

	a.listen()

	// Wait for SIGINT (CTRL-c), then close server and exit.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
