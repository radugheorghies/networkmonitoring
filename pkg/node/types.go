package node

import (
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

// Node is the default type of this package
type Node struct {
	wamp *turnpike.Client // connection to wamp server
	stop chan struct{}
}

type CPUReport struct {
	Avg     float64
	PerCore []float64
}

const (
	wampURL      = "ws://localhost:8087/"
	wampRealm    = "dataGateway"
	timeInterval = time.Second
)
