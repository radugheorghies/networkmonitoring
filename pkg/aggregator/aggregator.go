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
	// prepare files

	a.listen()

	if a.cpuFile != nil {
		defer a.cpuFile.Close()
	}
	if a.memFile != nil {
		defer a.memFile.Close()
	}
	if a.networkFile != nil {
		defer a.networkFile.Close()
	}

	// Wait for SIGINT (CTRL-c), then close server and exit.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
