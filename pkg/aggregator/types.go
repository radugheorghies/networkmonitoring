package aggregator

import (
	"sync"

	"gopkg.in/jcelliott/turnpike.v2"
)

// Aggregator is the default type of this package
type Aggregator struct {
	wamp      *turnpike.Client // connection to wamp server
	nodesData nodesData
}

type nodesData struct {
	nodes map[string]Node // map of node ID
	sync.Mutex
}

type Node struct {
	CPURecords     []CPUReport
	MemRecords     []MemReport
	NetworkRecords []NetworkReport
	AvgCPU         CPUReport
	AvgMem         MemReport
	AvgNetwork     NetworkReport
	TotalNetwork   NetworkReport
	FirstNetwork   NetworkReport
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
