package node

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func (n *Node) listen() {
	// will listen for commands
}

func (n *Node) initMonitoring() {
	// init the cpu monitoring
	cpuChan := make(chan string, 1)
	exitCPUChan := make(chan struct{}, 1)
	go getCpuInfo(cpuChan, exitCPUChan)

	// init the memmory alocation monitoring
	memChan := make(chan string, 1)
	exitMemChan := make(chan struct{}, 1)
	go getMemInfo(memChan, exitMemChan)

	// init network traffic monitoring
	netChan := make(chan string, 1)
	exitNetChan := make(chan struct{}, 1)
	go getNetInfo(netChan, exitNetChan)

	for {
		select {
		case c := <-cpuChan:
			log.Println(c)
		case m := <-memChan:
			log.Println(m)
		case n := <-netChan:
			log.Println(n)
		case <-n.stop:
			exitCPUChan <- struct{}{}
			exitMemChan <- struct{}{}
			exitNetChan <- struct{}{}
			return
		}
	}
}

// cpu info
func getCpuInfo(cpuChan chan string, exitCPUChan chan struct{}) {
	for {
		select {
		case <-exitCPUChan:
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
				Avg:     avg,
				PerCore: percents,
			}
			if rep, err := json.Marshal(cpuReport); err == nil {
				cpuChan <- string(rep)
			} else {
				log.Println(err)
			}

		}
	}
}

// memory
func getMemInfo(memChan chan string, exitMemChan chan struct{}) {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if memInfo, err := mem.VirtualMemory(); err == nil {
				memChan <- fmt.Sprintf("%v", memInfo)
			}
		case <-exitMemChan:
			return
		}
	}
}

// net
func getNetInfo(netChan chan string, exitNetChan chan struct{}) {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if info, err := net.IOCounters(true); err == nil {
				netChan <- fmt.Sprintf("%v", info[0])
			}
		case <-exitNetChan:
			return
		}
	}
}
