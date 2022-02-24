package aggregator

import (
	"encoding/json"
	"fmt"
	"log"
)

func (a *Aggregator) listen() {
	// will listen for commands
	go a.listenForCPURecords()
	go a.listenForMemRecords()
	go a.listenForNetworkRecords()
}

func (a *Aggregator) listenForCPURecords() {
	if err := a.wamp.Subscribe("cpuData", nil, func(args []interface{}, kwargs map[string]interface{}) {
		if len(args) > 0 {
			for _, arg := range args {
				result := CPUReport{}
				if err := json.Unmarshal([]byte(fmt.Sprintf("%s", arg)), &result); err != nil {
					log.Println("ERROR UNMARSHALING CPU DATA:", err)
					continue
				}

				a.nodesData.Lock()
				if _, ok := a.nodesData.nodes[result.NodeName]; !ok {
					a.nodesData.nodes[result.NodeName] = Node{}
				}
				tmpData := a.nodesData.nodes[result.NodeName]

				tmpData.CPURecords = append(tmpData.CPURecords, result)

				// now we calculate the averages
				length := len(tmpData.CPURecords)
				if length > 1 {
					tmpData.AvgCPU.Avg = (tmpData.AvgCPU.Avg*(float64(length)-1) + result.Avg) / float64(length)

					// we will do this for each cpu core
					for i, v := range tmpData.AvgCPU.PerCore {
						tmpData.AvgCPU.PerCore[i] = (v*(float64(length)-1) + result.Avg) / float64(length)
					}
				} else {
					tmpData.AvgCPU = result
				}

				a.nodesData.nodes[result.NodeName] = tmpData
				a.nodesData.Unlock()

			}
		}

	}); err != nil {
		log.Fatalln("Error subscribing to start monitoring cpu channel:", err)
	}
}

func (a *Aggregator) listenForMemRecords() {
	if err := a.wamp.Subscribe("memData", nil, func(args []interface{}, kwargs map[string]interface{}) {
		if len(args) > 0 {
			for _, arg := range args {
				result := MemReport{}
				if err := json.Unmarshal([]byte(fmt.Sprintf("%s", arg)), &result); err != nil {
					log.Println("ERROR UNMARSHALING MEM DATA:", err)
					continue
				}

				a.nodesData.Lock()
				if _, ok := a.nodesData.nodes[result.NodeName]; !ok {
					a.nodesData.nodes[result.NodeName] = Node{}
				}
				tmpData := a.nodesData.nodes[result.NodeName]

				tmpData.MemRecords = append(tmpData.MemRecords, result)

				// now we calculate the averages
				length := len(tmpData.MemRecords)
				if length > 1 {
					tmpData.AvgMem.Total = (tmpData.AvgMem.Total*(int64(length)-1) + int64(result.Total)) / int64(length)
					tmpData.AvgMem.Used = (tmpData.AvgMem.Used*(int64(length)-1) + int64(result.Used)) / int64(length)
					tmpData.AvgMem.Usedpercent = (tmpData.AvgMem.Usedpercent*(float64(length)-1) + float64(result.Usedpercent)) / float64(length)
				} else {
					tmpData.AvgMem = result
				}

				a.nodesData.nodes[result.NodeName] = tmpData
				a.nodesData.Unlock()

			}
		}

	}); err != nil {
		log.Fatalln("Error subscribing to start monitoring mem channel:", err)
	}
}

func (a *Aggregator) listenForNetworkRecords() {
	if err := a.wamp.Subscribe("networkData", nil, func(args []interface{}, kwargs map[string]interface{}) {
		if len(args) > 0 {
			for _, arg := range args {
				result := NetworkReport{}
				if err := json.Unmarshal([]byte(fmt.Sprintf("%s", arg)), &result); err != nil {
					log.Println("ERROR UNMARSHALING NETWORK DATA:", err)
					continue
				}

				a.nodesData.Lock()
				if _, ok := a.nodesData.nodes[result.NodeName]; !ok {
					a.nodesData.nodes[result.NodeName] = Node{}
				}
				tmpData := a.nodesData.nodes[result.NodeName]

				tmpData.NetworkRecords = append(tmpData.NetworkRecords, result)

				// now we calculate the averages
				length := len(tmpData.MemRecords)
				if length > 1 {
					tmpData.AvgNetwork.BytesRecv = (tmpData.AvgNetwork.BytesRecv*(int64(length)-1) + result.BytesRecv - tmpData.NetworkRecords[length-2].BytesRecv) / int64(length)
					tmpData.AvgNetwork.BytesSent = (tmpData.AvgNetwork.BytesSent*(int64(length)-1) + result.BytesSent - tmpData.NetworkRecords[length-2].BytesSent) / int64(length)
					tmpData.AvgNetwork.PacketsRecv = (tmpData.AvgNetwork.PacketsRecv*(int64(length)-1) + result.PacketsRecv - tmpData.NetworkRecords[length-2].PacketsRecv) / int64(length)
					tmpData.AvgNetwork.PacketsSent = (tmpData.AvgNetwork.PacketsSent*(int64(length)-1) + result.PacketsSent - tmpData.NetworkRecords[length-2].PacketsSent) / int64(length)
				} else {
					tmpData.AvgNetwork = NetworkReport{}
					tmpData.FirstNetwork = result
				}

				a.nodesData.nodes[result.NodeName] = tmpData
				a.nodesData.Unlock()

			}
		}

	}); err != nil {
		log.Fatalln("Error subscribing to start monitoring network channel:", err)
	}
}
