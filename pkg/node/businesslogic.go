package node

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"networkmonitoring/pkg/core/env"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func (n *Node) listen() {
	// will listen for commands
	go n.listenStartMonitoring()
	go n.listenStopMonitoring()
}

func (n *Node) listenStartMonitoring() {
	if err := n.wamp.Subscribe("startMonitoring", nil, func(args []interface{}, kwargs map[string]interface{}) {
		n.initMonitoring()
	}); err != nil {
		log.Fatalln("Error subscribing to start monitoring channel:", err)
	}
}

func (n *Node) listenStopMonitoring() {
	if err := n.wamp.Subscribe("stopMonitoring", nil, func(args []interface{}, kwargs map[string]interface{}) {
		n.stop <- struct{}{}
	}); err != nil {
		log.Fatalln("Error subscribing to stop monitoring channel:", err)
	}
}

func (n *Node) initMonitoring() {

	go n.getCpuInfo()
	go n.getMemInfo()
	go n.getNetInfo()

loop:
	for {
		select {
		case c := <-n.cpuChan:
			if err := n.wamp.Publish("cpuData", nil, []interface{}{c}, nil); err != nil {
				log.Println("Problem occurred while publishing cpu data:", err)
				n.exitCPUChan <- struct{}{}
				n.exitMemChan <- struct{}{}
				n.exitNetChan <- struct{}{}
				break loop
			}
			log.Println(c)
		case m := <-n.memChan:
			if err := n.wamp.Publish("memData", nil, []interface{}{m}, nil); err != nil {
				log.Println("Problem occurred while publishing mem data:", err)
				n.exitCPUChan <- struct{}{}
				n.exitMemChan <- struct{}{}
				n.exitNetChan <- struct{}{}
				break loop
			}
			log.Println(m)
		case nt := <-n.netChan:
			if err := n.wamp.Publish("networkData", nil, []interface{}{nt}, nil); err != nil {
				log.Println("Problem occurred while publishing network data:", err)
				n.exitCPUChan <- struct{}{}
				n.exitMemChan <- struct{}{}
				n.exitNetChan <- struct{}{}
				break loop
			}
			log.Println(n)
		case <-n.stop:
			n.exitCPUChan <- struct{}{}
			n.exitMemChan <- struct{}{}
			n.exitNetChan <- struct{}{}
			return
		}
	}

	n.init(false)
}

// cpu info
func (n *Node) getCpuInfo() {
	for {
		select {
		case <-n.exitCPUChan:
			return
		default:
			percents, err := cpu.Percent(timeInterval, true)
			if err != nil {
				continue
			}

			var avg float64
			for _, v := range percents {
				avg += v
			}
			avg = avg / float64(len(percents))

			cpuReport := CPUReport{
				NodeName: env.Vars.NodeName,
				Avg:      avg,
				PerCore:  percents,
			}
			if rep, err := json.Marshal(cpuReport); err == nil {
				n.cpuChan <- string(rep)
			} else {
				log.Println(err)
			}

		}
	}
}

// memory
func (n *Node) getMemInfo() {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if memInfo, err := mem.VirtualMemory(); err == nil {
				resTxt := fmt.Sprintf("%v", memInfo)
				res := make(map[string]interface{})
				json.Unmarshal([]byte(resTxt), &res)

				responseObj := MemReport{}
				responseObj.NodeName = env.Vars.NodeName

				if responseObj.Total, err = strconv.ParseInt(fmt.Sprintf("%.0f", res["total"]), 10, 64); err != nil {
					log.Println("Error converting Total to int64:", res["total"])
					continue
				}

				if responseObj.Used, err = strconv.ParseInt(fmt.Sprintf("%.0f", res["used"]), 10, 64); err != nil {
					log.Println("Error converting Used to int64:", res["used"])
					continue
				}

				if responseObj.Usedpercent, err = strconv.ParseFloat(fmt.Sprintf("%v", res["usedPercent"]), 64); err != nil {
					log.Println("Error converting UsedPercent to int64:", res["usedPercent"])
					continue
				}

				response, _ := json.Marshal(responseObj)
				n.memChan <- string(response)
			}
		case <-n.exitMemChan:
			return
		}
	}
}

// net
func (n *Node) getNetInfo() {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if info, err := net.IOCounters(true); err == nil {
				resTxt := fmt.Sprintf("%v", info[0])
				res := make(map[string]interface{})
				json.Unmarshal([]byte(resTxt), &res)

				responseObj := NetworkReport{}
				responseObj.NodeName = env.Vars.NodeName

				if responseObj.BytesSent, err = strconv.ParseInt(fmt.Sprintf("%.0f", res["bytesSent"]), 10, 64); err != nil {
					log.Println("Error converting BytesSent to int64:", res["bytesSent"])
					continue
				}

				if responseObj.BytesRecv, err = strconv.ParseInt(fmt.Sprintf("%.0f", res["bytesRecv"]), 10, 64); err != nil {
					log.Println("Error converting BytesReceived to int64:", res["bytesRecv"])
					continue
				}

				if responseObj.PacketsSent, err = strconv.ParseInt(fmt.Sprintf("%.0f", res["packetsSent"]), 10, 64); err != nil {
					log.Println("Error converting PackageSent to int64:", res["packetsSent"])
					continue
				}

				if responseObj.PacketsRecv, err = strconv.ParseInt(fmt.Sprintf("%.0f", res["packetsRecv"]), 10, 64); err != nil {
					log.Println("Error converting PackageReceived to int64:", res["packetsRecv"])
					continue
				}

				response, _ := json.Marshal(responseObj)
				n.netChan <- string(response)
			}
		case <-n.exitNetChan:
			return
		}
	}
}
