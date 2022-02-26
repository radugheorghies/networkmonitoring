package aggregator

import (
	"log"
	"networkmonitoring/pkg/core/env"
	"os"
	"time"

	"gopkg.in/jcelliott/turnpike.v2"
)

func (a *Aggregator) init(startListen bool) {
	// this will initiate everything
	a.initWampServer()

	if startListen {
		// listen for events
		a.listen()
	}
}

func (a *Aggregator) initWampServer() {
	for {
		client, err := turnpike.NewWebsocketClient(turnpike.JSON, env.Vars.WAMPURL, nil, nil, nil)
		if err != nil {
			log.Println("Error initiating the client:", err)
			time.Sleep(time.Second * 5)
			continue
		}

		if _, err := client.JoinRealm(env.Vars.WAMPRealm, nil); err != nil {
			log.Println("Error joining to the realm:", err)
			time.Sleep(time.Second * 5)
			continue
		}
		a.wamp = client
		break
	}
}

func (a *Aggregator) initFiles() (err error) {

	// open the file to write cpu data
	a.cpuFile, err = os.Create("/cpu.csv")
	if err != nil {
		log.Println("Error creating the cpu file:", err)
		return
	}

	// open the file to write mem data
	a.memFile, err = os.Create("/mem.csv")
	if err != nil {
		log.Println("Error creating the mem file:", err)
		return
	}

	// open the file to write network data
	a.networkFile, err = os.Create("/network.csv")
	if err != nil {
		log.Println("Error creating the network file:", err)
		return
	}

	return
}
