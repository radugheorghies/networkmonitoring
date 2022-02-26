package aggregator

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func (a *Aggregator) listen() {
	// will listen for commands
	go a.listenForCPURecords()
	go a.listenForMemRecords()
	go a.listenForNetworkRecords()
	go a.listenStartAggregator()
	go a.listenStopAggregator()
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

				s := result.NodeName + "," + fmt.Sprintf("%.2f", result.Avg) + "," + strings.Trim(strings.Replace(fmt.Sprint(result.PerCore), " ", ",", -1), "[]")
				a.cpuFile.WriteString(s)

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

				s := result.NodeName + "," + fmt.Sprintf("%d", result.Total) + "," + fmt.Sprintf("%d", result.Used) + "," + fmt.Sprintf("%.2f", result.Usedpercent)
				a.memFile.WriteString(s)

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
				length := len(tmpData.NetworkRecords)
				if length > 1 {
					tmpData.AvgNetwork.BytesRecv = (tmpData.AvgNetwork.BytesRecv*(int64(length)-1) + result.BytesRecv - tmpData.NetworkRecords[length-2].BytesRecv) / int64(length)
					tmpData.AvgNetwork.BytesSent = (tmpData.AvgNetwork.BytesSent*(int64(length)-1) + result.BytesSent - tmpData.NetworkRecords[length-2].BytesSent) / int64(length)
					tmpData.AvgNetwork.PacketsRecv = (tmpData.AvgNetwork.PacketsRecv*(int64(length)-1) + result.PacketsRecv - tmpData.NetworkRecords[length-2].PacketsRecv) / int64(length)
					tmpData.AvgNetwork.PacketsSent = (tmpData.AvgNetwork.PacketsSent*(int64(length)-1) + result.PacketsSent - tmpData.NetworkRecords[length-2].PacketsSent) / int64(length)
				} else {
					tmpData.AvgNetwork = NetworkReport{}
				}

				tmpData.TotalNetwork.BytesRecv = result.BytesRecv - tmpData.NetworkRecords[0].BytesRecv
				tmpData.TotalNetwork.BytesSent = result.BytesSent - tmpData.NetworkRecords[0].BytesSent
				tmpData.TotalNetwork.PacketsRecv = result.PacketsRecv - tmpData.NetworkRecords[0].PacketsRecv
				tmpData.TotalNetwork.PacketsSent = result.PacketsSent - tmpData.NetworkRecords[0].PacketsSent

				a.nodesData.nodes[result.NodeName] = tmpData
				a.nodesData.Unlock()

				s := result.NodeName + "," + fmt.Sprintf("%d", tmpData.TotalNetwork.BytesSent) + "," + fmt.Sprintf("%d", tmpData.TotalNetwork.BytesRecv) + "," +
					fmt.Sprintf("%d", tmpData.TotalNetwork.PacketsSent) + "," + fmt.Sprintf("%d", tmpData.TotalNetwork.PacketsRecv)
				a.networkFile.WriteString(s)

			}
		}

	}); err != nil {
		log.Fatalln("Error subscribing to start monitoring network channel:", err)
	}
}

func (a *Aggregator) listenStartAggregator() {
	if err := a.wamp.Subscribe("startAggregator", nil, func(args []interface{}, kwargs map[string]interface{}) {
		a.initFiles()
		if err := a.wamp.Publish("stratMonitoring", nil, []interface{}{}, nil); err != nil {
			log.Println("Problem occurred while publishing start commands to nodes:", err)
		}
	}); err != nil {
		log.Fatalln("Error subscribing to start monitoring channel:", err)
	}
}

func (a *Aggregator) listenStopAggregator() {
	if err := a.wamp.Subscribe("stopAggregator", nil, func(args []interface{}, kwargs map[string]interface{}) {
		if err := a.wamp.Publish("stopMonitoring", nil, []interface{}{}, nil); err != nil {
			log.Println("Problem occurred while publishing start commands to nodes:", err)
		}
		if a.cpuFile != nil {
			defer a.cpuFile.Close()
		}
		if a.memFile != nil {
			defer a.memFile.Close()
		}
		if a.networkFile != nil {
			defer a.networkFile.Close()
		}
	}); err != nil {
		log.Fatalln("Error subscribing to stop monitoring channel:", err)
	}
}
