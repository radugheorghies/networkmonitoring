package node

import (
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

// Node is the default type of this package
type Node struct {
	wamp        *turnpike.Client // connection to wamp server
	stop        chan struct{}
	cpuChan     chan string
	exitCPUChan chan struct{}
	memChan     chan string
	exitMemChan chan struct{}
	netChan     chan string
	exitNetChan chan struct{}
}

type CPUReport struct {
	NodeName string
	Avg      float64
	PerCore  []float64
}

type MemReport struct {
	NodeName    string
	Total       int64
	Used        int64
	Usedpercent float64
}

type NetworkReport struct {
	NodeName    string
	BytesSent   int64
	BytesRecv   int64
	PacketsSent int64
	PacketsRecv int64
}

const (
	timeInterval = time.Second
)
