package node

import (
	"os"
	"os/signal"
)

// New is the factory function for this package
func New() *Node {
	return &Node{
		stop: make(chan struct{}),
	}
}

// Run will start the **magic**
func (n *Node) Run() {
	n.init()

	n.listen()

	// Wait for SIGINT (CTRL-c), then close server and exit.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
