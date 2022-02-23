package node

import (
	"os"
	"os/signal"
)

// New is the factory function for this package
func New() *Node {
	return &Node{
		stop: make(chan struct{}),
		// init the cpu monitoring
		cpuChan:     make(chan string, 1),
		exitCPUChan: make(chan struct{}, 1),
		// init the memmory alocation monitoring
		memChan:     make(chan string, 1),
		exitMemChan: make(chan struct{}, 1),
		// init network traffic monitoring
		netChan:     make(chan string, 1),
		exitNetChan: make(chan struct{}, 1),
	}
}

// Run will start the **magic**
func (n *Node) Run() {
	n.init(true)

	n.listen()

	// Wait for SIGINT (CTRL-c), then close server and exit.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
